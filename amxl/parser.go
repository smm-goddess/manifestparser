package amxl

import (
	"bytes"
	"fmt"
	"github.com/smm-goddess/manifestparser/amxl/structs"
)

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

	for _, content := range contentChunks {
		switch content.(type) {
		case structs.StartNamespaceChunk:
			if c, ok := content.(structs.StartNamespaceChunk); ok {
				prefix := stringChunk.Utf8Strings[c.Prefix]
				uri := stringChunk.Utf8Strings[c.Uri]
				fmt.Println("---------- start namespace chunk ----------")
				fmt.Println("chunk size:", c.Size)
				fmt.Println("line number:", c.LineNumber)
				fmt.Println("prefix:", prefix)
				fmt.Println("uri:", uri)
				prefixUriMap[prefix] = uri
				uriPrefixMap[uri] = prefix
			}
		case structs.EndNamespaceChunk:
			if c, ok := content.(structs.EndNamespaceChunk); ok {
				fmt.Println("---------- end namespace chunk ----------")
				fmt.Println("chunk size:", c.Size)
				fmt.Println("line number:", c.LineNumber)
				fmt.Println("prefix:", stringChunk.Utf8Strings[c.Prefix])
				fmt.Println("uri:", stringChunk.Utf8Strings[c.Uri])
			}
		case structs.EndTagChunk:
			if c, ok := content.(structs.EndTagChunk); ok {
				fmt.Println("---------- end tag chunk ----------")
				fmt.Println("chunk size:", c.Size)
				fmt.Println("line number:", c.LineNumber)
				if c.Prefix < stringChunk.Size {
					fmt.Println("prefix:", stringChunk.Utf8Strings[c.Prefix])
				}
				fmt.Println("uri:", stringChunk.Utf8Strings[c.Uri])
			}
		case structs.StartTagChunk:
			if c, ok := content.(structs.StartTagChunk); ok {
				fmt.Println("---------- start tag chunk ----------")
				fmt.Println("tag name:", stringChunk.Utf8Strings[c.Name])
				if c.NameSpaceUri < stringChunk.Size {
					fmt.Println("tag namespace uri:", stringChunk.Utf8Strings[c.NameSpaceUri])
				}
				for _, attribute := range c.Attributes {
					fmt.Println("---------- attribute ----------")
					fmt.Println("name:", stringChunk.Utf8Strings[attribute.Name])
					if attribute.NamespaceUri < stringChunk.Size {
						fmt.Println("prefix:", stringChunk.Utf8Strings[attribute.NamespaceUri])
					}
					switch attribute.Type {
					case 0x01:
						fmt.Printf("data:@%x\n", attribute.Data)
					case 0x03:
						fmt.Println("data:", stringChunk.Utf8Strings[attribute.Data])
					case 0x10:
						fmt.Printf("data:%d\n", attribute.Data)
					case 0x12:
						fmt.Println("data:", attribute.Data == 0)
					default:
						fmt.Printf("type:%x\n", attribute.Type)
					}
				}
			}
		case structs.TextChunk:
			if _, ok := content.(structs.TextChunk); ok {
				fmt.Println("---------- text chunk ----------")
			}
		}
	}

}
