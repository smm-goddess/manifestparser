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

func ParseHeader(bs []byte) (Header, uint32) {
	var header Header
	_ = binary.Read(bytes.NewBuffer(bs), binary.LittleEndian, &header)
	return header, 8
}

func (header Header) String() string {
	return fmt.Sprintf("%x %x", header.Magic, header.FileSize)
}
