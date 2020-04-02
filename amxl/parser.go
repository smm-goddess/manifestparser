package amxl

import (
	"encoding/xml"
	"fmt"
	"github.com/smm-goddess/manifestparser/amxl/structs"
	"os"
)

func Parse(bs []byte) {
	var readCount, offset uint32
	header, readCount := structs.ParseHeader(bs)
	fmt.Println(header)
	offset += readCount
	stringChunk, readCount := structs.ParseStringChunk(bs, offset)
	offset += readCount
	for _, s := range stringChunk.Utf8Strings {
		fmt.Println(s)
	}
	resourceChunk, readCount := structs.ParseResourceIdChunk(bs, offset)
	offset += readCount
	fmt.Println("resource:", resourceChunk)

	encoder := xml.NewEncoder(os.Stdout)
	contentArray := structs.ReadContentChunks(bs, offset)
	c := contentArray[0]
	if c, ok := c.(structs.StartNamespaceChunk); ok {
		content := xml.StartElement{
			Name: xml.Name{
				Space: stringChunk.Utf8Strings[c.Prefix],
				Local: stringChunk.Utf8Strings[c.Uri],
			},
			Attr: nil,
		}
		encoder.Encode(content)
	}
	//
	//for index, content := range contentArray {
	//	switch content.(type) {
	//	case structs.StartNamespaceChunk:
	//		if c, ok := content.(structs.StartNamespaceChunk); ok {
	//			fmt.Println(stringChunk.Utf8Strings[c.Prefix], ":", stringChunk.Utf8Strings[c.Uri])
	//		}
	//	case structs.EndNamespaceChunk:
	//		if c, ok := content.(structs.EndNamespaceChunk); ok {
	//			fmt.Println(stringChunk.Utf8Strings[c.Prefix], ":", stringChunk.Utf8Strings[c.Uri])
	//		}
	//	case structs.EndTagChunk:
	//		if c, ok := content.(structs.EndTagChunk); ok {
	//			fmt.Println(stringChunk.Utf8Strings[c.Prefix], ":", stringChunk.Utf8Strings[c.Uri])
	//		}
	//	case structs.StartTagChunk:
	//		if c, ok := content.(structs.StartTagChunk); ok {
	//			fmt.Println(stringChunk.Utf8Strings[c.Prefix], ":", stringChunk.Utf8Strings[c.Uri])
	//		}
	//	case structs.TextChunk:
	//		if _, ok := content.(structs.TextChunk); ok {
	//
	//		}
	//	}
	//}

}
