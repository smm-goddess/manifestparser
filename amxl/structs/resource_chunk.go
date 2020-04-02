package structs

import (
	"bytes"
	"encoding/binary"
)

type ResourceIdChunk struct {
	Signature uint32
	Size      uint32
	Items     []uint32
}

func ParseResourceIdChunk(bs []byte, offset uint32) (resourceChunk ResourceIdChunk, readCount uint32) {
	params := make([]uint32, 2, 2)
	_ = binary.Read(bytes.NewBuffer(bs[offset:]), binary.LittleEndian, &params)
	resourceChunk.Signature = params[0]
	resourceChunk.Size = params[1]
	items := make([]uint32, resourceChunk.Size/4-2, resourceChunk.Size/4-2)
	_ = binary.Read(bytes.NewBuffer(bs[offset+8:]), binary.LittleEndian, &items)
	readCount = resourceChunk.Size
	resourceChunk.Items = items
	return
}
