A library for modifying XML documents.

[![GoDoc](http://godoc.org/github.com/PlanitarInc/go-xmlproc?status.svg)](http://godoc.org/github.com/PlanitarInc/go-xmlproc)

# Why?

2 main reasons:
  - I needed to modify slightly XML documents and did not find any simple and nice golang library to do that
  - I needed the library to be friendly to XML namespaces

# Examples

Changing attribute value:

```go
package main

import (
	"bytes"
	"encoding/xml"
	"os"

	"github.com/PlanitarInc/go-xmlproc"
	"github.com/PlanitarInc/go-xmlproc/mappers"
)

const example = `<taxii_11:Discovery_Response 
  xmlns:taxii="http://taxii.mitre.org/messages/taxii_xml_binding-1"
  xmlns:taxii_11="http://taxii.mitre.org/messages/taxii_xml_binding-1.1"
  xmlns:tdq="http://taxii.mitre.org/query/taxii_default_query-1" 
  message_id="32898" 
  in_response_to="1">
</taxii_11:Discovery_Response>`

type IncInResponseTo struct{}

func (m IncInResponseTo) Map(t xml.Token) (xml.Token, error) {
	switch token := t.(type) {
	default:
		return t, nil

	case xml.StartElement:
		if token.Name.Local != "Discovery_Response" {
			return t, nil
		}
		for i := range token.Attr {
			if token.Attr[i].Name.Local == "in_response_to" {
				token.Attr[i].Value = "2"
				break
			}
		}
		return t, nil
	}
}

func main() {
	src := bytes.NewBufferString(example)
	dec := xml.NewDecoder(src)

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	defer enc.Flush()

	p := xmlproc.Processor{
		Mappers: []xmlproc.Mapper{
			&mappers.Pruner{},
			&IncInResponseTo{},
			&mappers.NSNormalizer{},
		},
	}
	if err := p.Process(enc, dec); err != nil {
		panic(err)
	}
}
```

The produced output:

```
<taxii_11:Discovery_Response 
  xmlns:taxii="http://taxii.mitre.org/messages/taxii_xml_binding-1"
  xmlns:taxii_11="http://taxii.mitre.org/messages/taxii_xml_binding-1.1"
  xmlns:tdq="http://taxii.mitre.org/query/taxii_default_query-1" 
  message_id="32898" 
  in_response_to="2">
</taxii_11:Discovery_Response>
```
