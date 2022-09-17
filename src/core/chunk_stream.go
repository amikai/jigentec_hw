package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// each chunk size is 6 bytes
const seqOffset int = 0
const seqSize int = 4

const lenOffset int = 4
const lenSize int = 2

const dataOffset int = 6

type ChunkStream struct {
	Seq  uint32
	Len  uint16
	Data []byte
}

func (cs *ChunkStream) Read(r io.Reader) error {
	seqLenBytes := make([]byte, 6)
	_, err := r.Read(seqLenBytes)
	if err != nil {
		return fmt.Errorf("faild to read seq and len: %w", err)
	}

	buf := bytes.NewReader(seqLenBytes[seqOffset : seqOffset+seqSize])
	err = binary.Read(buf, binary.BigEndian, &cs.Seq)
	if err != nil {
		return fmt.Errorf("failed read seq num: %w", err)
	}

	buf = bytes.NewReader(seqLenBytes[lenOffset : lenOffset+lenSize])
	err = binary.Read(buf, binary.BigEndian, &cs.Len)
	if err != nil {
		return fmt.Errorf("failed read seq num: %w", err)
	}

	dataBytes := make([]byte, cs.Len)
	_, err = r.Read(dataBytes)
	if err != nil {
		return fmt.Errorf("faild to read seq and len: %w", err)
	}
	cs.Data = dataBytes
	return nil
}
