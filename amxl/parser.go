package amxl

import (
	"bytes"
	"fmt"
	"github.com/smm-goddess/manifestparser/amxl/structs"
)

var head = `
<?xml version="1.0" encoding="utf-8"?>
`

var indent = "  "

var uriPrefixMap = make(map[string]string)
var prefixUriMap = make(map[string]string)

func Parse(bs []byte) {
	buffer := bytes.NewBuffer(bs)
	header := structs.ParseHeader(buffer)
	fmt.Println(header)
	stringChunk := structs.ParseStringChunk(buffer)
	for i, s := range stringChunk.Utf8Strings {
		fmt.Println(i, s)
	}
	resourceChunk := structs.ParseResourceIdChunk(buffer)
	fmt.Println("resource:", resourceChunk)

	contentChunks := structs.ReadContentChunks(buffer)

	buff := bytes.NewBufferString(head)
	writeNamespace := false

	currentLine := uint32(2)

	indentCount := 0

	for _, content := range contentChunks {
		switch content.(type) {
		case structs.StartNamespaceChunk:
			if c, ok := content.(structs.StartNamespaceChunk); ok {
				prefix := stringChunk.Utf8Strings[c.Prefix]
				uri := stringChunk.Utf8Strings[c.Uri]
				prefixUriMap[prefix] = uri
				uriPrefixMap[uri] = prefix
			}
		case structs.EndNamespaceChunk:

		case structs.EndTagChunk:
			indentCount--
			if c, ok := content.(structs.EndTagChunk); ok {
				for currentLine < c.LineNumber {
					buff.WriteByte(0x0A)
					for i := 0; i < indentCount; i++ {
						buff.WriteString(indent)
					}
					currentLine = c.LineNumber
				}
				buff.WriteByte('<')
				buff.WriteByte('/')
				buff.WriteString(stringChunk.Utf8Strings[c.Uri])
				buff.WriteByte('>')
			}
		case structs.StartTagChunk:
			if c, ok := content.(structs.StartTagChunk); ok {
				for currentLine < c.LineNumber {
					buff.WriteByte(0x0A)
					for i := 0; i < indentCount; i++ {
						buff.WriteString(indent)
					}
					currentLine = c.LineNumber
				}
				buff.WriteByte('<')
				buff.WriteString(stringChunk.Utf8Strings[c.Name])
				buff.WriteByte(' ')
				if !writeNamespace {
					for prefix, uri := range prefixUriMap {
						buff.WriteString("xmlns:")
						buff.WriteString(prefix)
						buff.WriteByte('=')
						buff.WriteByte('"')
						buff.WriteString(uri)
						buff.WriteByte('"')
						buff.WriteByte(' ')
					}
					writeNamespace = true
				}
				for _, attribute := range c.Attributes {
					if attribute.NamespaceUri < stringChunk.Size {
						buff.WriteString(uriPrefixMap[stringChunk.Utf8Strings[attribute.NamespaceUri]])
						buff.WriteByte(':')
					}
					buff.WriteString(stringChunk.Utf8Strings[attribute.Name])
					buff.WriteByte('=')
					switch attribute.Type {
					case 0x01: //resourceId
						buff.WriteString(fmt.Sprintf("\"@%x\"", attribute.Data))
					case 0x03: //string
						buff.WriteString(fmt.Sprint("\"", stringChunk.Utf8Strings[attribute.Data], "\""))
					case 0x10: // int
						buff.WriteString(fmt.Sprintf("\"%d\"", attribute.Data))
					case 0x12: // boolean
						buff.WriteString(fmt.Sprint("\"", attribute.Data != 0, "\""))
					case 0x11: // hex
						buff.WriteString(fmt.Sprintf("\"0x%x\"", attribute.Data))
					default:
						fmt.Printf("UNKNOWN %x", attribute.Type)
						buff.WriteString(fmt.Sprintf("\"0x%x\"", attribute.Data))
					}
				}
				buff.WriteByte('>')
				indentCount++
			}
		case structs.TextChunk:
			if _, ok := content.(structs.TextChunk); ok {
				//fmt.Println("---------- text chunk ----------")
			}
		}
	}

	fmt.Println(buff.String())
}
