package structs

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type StringChunk struct {
	Signature        uint32
	Size             uint32
	StringCount      uint32
	StyleCount       uint32
	Flags            uint32
	StringPoolOffset uint32
	StylePoolOffset  uint32
	StringOffsets    []uint32
	Utf8Strings      []string
}

func ReadUtf16String(buffer *bytes.Buffer, offset uint32) string {
	var cnt uint16
	_ = binary.Read(buffer, binary.LittleEndian, &cnt)
	buf := make([]byte, 2*cnt, 2*cnt)
	_, _ = buffer.Read(buf)
	s, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), buf)
	_ = binary.Read(buffer, binary.LittleEndian, &cnt)
	return string(s)
}

func ParseStringChunk(buffer *bytes.Buffer) StringChunk {
	sc := make([]uint32, 7, 7)
	_ = binary.Read(buffer, binary.LittleEndian, &sc)
	stringChunk := StringChunk{
		Signature:        sc[0],
		Size:             sc[1],
		StringCount:      sc[2],
		StyleCount:       sc[3],
		Flags:            sc[4],
		StringPoolOffset: sc[5],
		StylePoolOffset:  sc[6],
		StringOffsets:    nil,
		Utf8Strings:      nil,
	}
	stringOffsets := make([]uint32, stringChunk.StringCount, stringChunk.StringCount)
	strings := make([]string, stringChunk.StringCount, stringChunk.StringCount)
	_ = binary.Read(buffer, binary.LittleEndian, &stringOffsets)
	for index, relativeOffset := range stringOffsets {
		strings[index] = ReadUtf16String(buffer, stringChunk.StringPoolOffset+0x08+relativeOffset)
	}
	stringChunk.StringOffsets = stringOffsets
	stringChunk.Utf8Strings = strings
	return stringChunk
}
