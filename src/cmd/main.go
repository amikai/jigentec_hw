package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"jigentec.homework/core"
)

const (
	ServerHost     = "assignment.jigentec.com"
	ServerPort     = "49152"
	ServerProtocol = "tcp"
)

var (
	fileName = kingpin.Arg("file", "output file name").Default("./file.txt").String()
)

func main() {
	kingpin.Parse()

	//establish conn
	conn, err := net.Dial(ServerProtocol, fmt.Sprintf("%s:%s", ServerHost, ServerPort))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// collect packets and parse
	chunks := []*core.ChunkStream{}
	for {
		chunk := &core.ChunkStream{}
		if err := chunk.Read(conn); err != nil {
			if errors.Is(err, io.EOF) {
				log.Infof("Read done: total number of chunks: %d", len(chunks))
				break
			}
			log.Fatalf("chunk stream read failed: %s", err)
		}
		chunks = append(chunks, chunk)
	}

	// sort all chunk by sequence number
	sort.Slice(chunks, func(i, j int) bool { return chunks[i].Seq < chunks[j].Seq })

	// write all chunks to file
	file, err := os.Create(*fileName)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	defer file.Close()
	for _, b := range chunks {
		if _, err := file.Write(b.Data); err != nil {
			log.Fatalf("failed to write byte to file: %s", err)
		}
	}
	log.Infof("Output file: %s", *fileName)
}
