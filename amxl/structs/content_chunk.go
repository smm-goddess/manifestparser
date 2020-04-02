package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	START_NAMESPACE_CHUNK = 0x00100100
	END_NAMESPACE_CHUNK   = 0x00100101
	START_TAG_CHUNK       = 0x00100102
	END_TAG_CHUNK         = 0x00100103
	TEXT_CHUNK            = 0x00100104
)

type StartNamespaceChunk struct {
	Type       uint32
	Size       uint32
	LineNumber uint32
	Comment    uint32
	Prefix     uint32
	Uri        uint32
}

type EndNamespaceChunk StartNamespaceChunk
type EndTagChunk StartNamespaceChunk

type Attribute struct {
	NamespaceUri uint32
	Name         uint32
	ValueString  uint32
	Type         uint32
	Data         uint32
}

type StartTagChunk struct {
	StartNamespaceChunk
	Flags          uint32
	AttributeCount uint32
	ClassAttribute uint32
	Attributes     []Attribute
}

type TextChunk struct {
	Signature  uint32
	Size       uint32
	LineNumber uint32
	Unknown1   uint32
	Name       uint32
	Unknown2   uint32
	Unknown3   uint32
}

func ReadContentChunks(bs []byte, offset uint32) []interface{} {
	var chunkType uint32
	buffer := bytes.NewBuffer(bs[offset:])
	contentArray := make([]interface{}, 0)
	for offset < uint32(len(bs)) {
		_ = binary.Read(buffer, binary.LittleEndian, &chunkType)
		switch chunkType {
		case START_NAMESPACE_CHUNK:
			fmt.Println("start namespace chunk")
			properties := make([]uint32, 5, 5)
			_ = binary.Read(buffer, binary.LittleEndian, &properties)
			startNamespaceChunk := StartNamespaceChunk{
				Type:       START_NAMESPACE_CHUNK,
				Size:       properties[0],
				LineNumber: properties[1],
				Comment:    properties[2],
				Prefix:     properties[3],
				Uri:        properties[4],
			}
			contentArray = append(contentArray, startNamespaceChunk)
			offset += 24
		case END_NAMESPACE_CHUNK:
			fmt.Println("end namespace chunk")
			properties := make([]uint32, 5, 5)
			_ = binary.Read(buffer, binary.LittleEndian, &properties)
			endNamespaceChunk := EndNamespaceChunk{
				Type:       END_NAMESPACE_CHUNK,
				Size:       properties[0],
				LineNumber: properties[1],
				Comment:    properties[2],
				Prefix:     properties[3],
				Uri:        properties[4],
			}
			contentArray = append(contentArray, endNamespaceChunk)
			offset += 24
		case START_TAG_CHUNK:
			fmt.Println("start tag chunk")
			properties := make([]uint32, 8, 8)
			_ = binary.Read(buffer, binary.LittleEndian, &properties)
			attributes := make([]Attribute, properties[6], properties[6])
			_ = binary.Read(buffer, binary.LittleEndian, &attributes)
			startTagChunk := StartTagChunk{
				StartNamespaceChunk: StartNamespaceChunk{
					Type:       START_TAG_CHUNK,
					Size:       properties[0],
					LineNumber: properties[1],
					Comment:    properties[2],
					Prefix:     properties[3],
					Uri:        properties[4],
				},
				Flags:          properties[5],
				AttributeCount: properties[6],
				ClassAttribute: properties[7],
				Attributes:     attributes,
			}
			contentArray = append(contentArray, startTagChunk)
			offset += 36
			offset += 20 * properties[6]
		case END_TAG_CHUNK:
			fmt.Println("end tag chunk")
			properties := make([]uint32, 5, 5)
			_ = binary.Read(buffer, binary.LittleEndian, &properties)
			endTagChunk := EndTagChunk{
				Type:       END_NAMESPACE_CHUNK,
				Size:       properties[0],
				LineNumber: properties[1],
				Comment:    properties[2],
				Prefix:     properties[3],
				Uri:        properties[4],
			}
			contentArray = append(contentArray, endTagChunk)
			offset += 24
		case TEXT_CHUNK:
			fmt.Println("text chunk")
			properties := make([]uint32, 6, 6)
			_ = binary.Read(buffer, binary.LittleEndian, &properties)
			textChunk := TextChunk{
				Signature:  TEXT_CHUNK,
				Size:       properties[0],
				LineNumber: properties[1],
				Unknown1:   properties[2],
				Name:       properties[3],
				Unknown2:   properties[4],
				Unknown3:   properties[5],
			}
			contentArray = append(contentArray, textChunk)
			offset += 28
		}
	}
	return contentArray
}
