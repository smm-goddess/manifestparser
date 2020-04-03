package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Header struct {
	Magic    uint32
	FileSize uint32
}

func ParseHeader(buffer *bytes.Buffer) Header {
	var header Header
	_ = binary.Read(buffer, binary.LittleEndian, &header)
	return header
}

func (header Header) String() string {
	return fmt.Sprintf("%x %x", header.Magic, header.FileSize)
}
