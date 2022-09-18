package core

import (
	"encoding/binary"
	"fmt"
	"io"

	"jigentec.homework/utils"
)

// +--------------------------------+----------------+------------------------------------
// |   Sequence Number (32-bits)    |  Len (16-bits) | File data (Len number of bytes) ...
// +--------------------------------+----------------+------------------------------------

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

func (cs *ChunkStream) Read(r io.Reader, pool *utils.Pool) error {
	seqAndLenBytes := pool.Get(seqSize + lenSize)
	_, err := io.ReadFull(r, seqAndLenBytes)
	if err != nil {
		return fmt.Errorf("faild to seq and len: %w", err)
	}
	cs.Seq = binary.BigEndian.Uint32(seqAndLenBytes[seqOffset : seqOffset+seqSize])
	cs.Len = binary.BigEndian.Uint16(seqAndLenBytes[lenOffset : lenOffset+lenSize])

	dataBytes := pool.Get(int(cs.Len))
	_, err = io.ReadFull(r, dataBytes)
	if err != nil {
		return fmt.Errorf("faild to read data: %w", err)
	}
	cs.Data = dataBytes
	return nil
}
