package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"jigentec.homework/core"
	"jigentec.homework/utils"
)

const (
	ServerHost     = "assignment.jigentec.com"
	ServerPort     = "49152"
	ServerProtocol = "tcp"
)

var (
	filePath = kingpin.Flag("file", "output file path").Default("./download_file").String()
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
	bytePool := utils.NewPool()
	connBuf := bufio.NewReader(conn)
	var totalSize uint64
	for {
		chunk := &core.ChunkStream{}
		if err := chunk.Read(connBuf, bytePool); err != nil {
			if errors.Is(err, io.EOF) {
				log.Infof("Read done: total size of data: %d", totalSize)
				break
			}
			log.Fatalf("chunk stream read failed: %s", err)
		}
		totalSize += uint64(chunk.Len)
		chunks = append(chunks, chunk)
	}

	// sort all chunk by sequence number
	sort.Slice(chunks, func(i, j int) bool { return chunks[i].Seq < chunks[j].Seq })

	// write all chunks to file
	file, err := os.Create(*filePath)
	writer := bufio.NewWriter(file)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	defer file.Close()
	for _, chunk := range chunks {
		if _, err := writer.Write(chunk.Data); err != nil {
			log.Fatalf("failed to write byte to file: %s", err)
		}
	}
	writer.Flush()
	log.Infof("Output file: %s", *filePath)
}
