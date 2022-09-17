package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"jigentec.homework/utils"
)

type fakeReader struct {
	packets []byte
	offset  int
}

func (r *fakeReader) Read(dst []byte) (int, error) {
	n := copy(dst, r.packets[r.offset:])
	r.offset += n
	return n, nil
}

func TestChunkStreamRead(t *testing.T) {
	packets := []byte{
		0x11, 0x22, 0x33, 0x44,
		0x00, 0x02,
		0x11, 0x12,
	}

	want := &ChunkStream{
		Seq:  287454020,
		Len:  2,
		Data: []byte{0x11, 0x12},
	}

	bytePool := utils.NewPool()
	fr := &fakeReader{packets: packets}
	cs := &ChunkStream{}
	cs.Read(fr, bytePool)

	assert.Equal(t, want, cs)

}
