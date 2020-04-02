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

func ReadUtf16String(bs []byte, offset uint32) (string, uint32) {
	var cnt uint16
	_ = binary.Read(bytes.NewBuffer(bs[offset:]), binary.LittleEndian, &cnt)
	buffer := make([]byte, 2*(cnt+1), 2*(cnt+1))
	copy(buffer, bs[offset+2:])
	s, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), buffer)
	return string(s), uint32(cnt+2) * 2
}

func ParseStringChunk(bs []byte, offset uint32) (StringChunk, uint32) {
	sc, readCount := make([]uint32, 7, 7), uint32(28)
	_ = binary.Read(bytes.NewBuffer(bs[offset:]), binary.LittleEndian, &sc)
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
	readCount += 4 * stringChunk.StringCount
	strings := make([]string, stringChunk.StringCount, stringChunk.StringCount)
	_ = binary.Read(bytes.NewReader(bs[offset+28:]), binary.LittleEndian, &stringOffsets)
	var cnt uint32
	for index, relativeOffset := range stringOffsets {
		strings[index], cnt = ReadUtf16String(bs, stringChunk.StringPoolOffset+0x08+relativeOffset)
		readCount += cnt
	}
	stringChunk.StringOffsets = stringOffsets
	stringChunk.Utf8Strings = strings
	return stringChunk, readCount
}
