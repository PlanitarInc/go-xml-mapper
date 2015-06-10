package mappers

import (
	"encoding/xml"
	"fmt"
)

// Logger prints to stdout a token being processed.
type Logger struct{}

func (p Logger) Map(t xml.Token) (xml.Token, error) {
	fmt.Printf("  processing token: %#v\n", t)
	return t, nil
}
