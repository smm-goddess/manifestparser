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

func ParseResourceIdChunk(buffer *bytes.Buffer) (resourceChunk ResourceIdChunk) {
	params := make([]uint32, 2, 2)
	_ = binary.Read(buffer, binary.LittleEndian, &params)
	resourceChunk.Signature = params[0]
	resourceChunk.Size = params[1]
	items := make([]uint32, resourceChunk.Size/4-2, resourceChunk.Size/4-2)
	_ = binary.Read(buffer, binary.LittleEndian, &items)
	resourceChunk.Items = items
	return
}
