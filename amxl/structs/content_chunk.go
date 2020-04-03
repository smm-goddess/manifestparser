package structs

import (
	"bytes"
	"encoding/binary"
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
	Type           uint32
	Size           uint32
	LineNumber     uint32
	UNKNOWN        uint32
	NameSpaceUri   uint32
	Name           uint32
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

func ReadContentChunks(buffer *bytes.Buffer) []interface{} {
	var chunkType uint32
	contentArray := make([]interface{}, 0)
	for {
		_ = binary.Read(buffer, binary.LittleEndian, &chunkType)
		switch chunkType {
		case START_NAMESPACE_CHUNK:
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
		case END_NAMESPACE_CHUNK:
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
			return contentArray
		case START_TAG_CHUNK:
			properties := make([]uint32, 8, 8)
			_ = binary.Read(buffer, binary.LittleEndian, &properties)
			attributes := make([]Attribute, properties[6], properties[6])
			_ = binary.Read(buffer, binary.LittleEndian, &attributes)
			for index := range attributes {
				attributes[index].Type = attributes[index].Type >> 24
			}
			startTagChunk := StartTagChunk{
				Type:           START_TAG_CHUNK,
				Size:           properties[0],
				LineNumber:     properties[1],
				UNKNOWN:        properties[2],
				NameSpaceUri:   properties[3],
				Name:           properties[4],
				Flags:          properties[5],
				AttributeCount: properties[6],
				ClassAttribute: properties[7],
				Attributes:     attributes,
			}
			contentArray = append(contentArray, startTagChunk)
		case END_TAG_CHUNK:
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
		case TEXT_CHUNK:
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
		}
	}
}
