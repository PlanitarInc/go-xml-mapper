package mappers

import (
	"encoding/xml"
	"strings"
)

// Pruner removes all XML comments (xml.Comment) and
// any chardata (xml.CharData), consisting of whitespaces only.
type Pruner struct{}

func (p Pruner) Map(t xml.Token) (xml.Token, error) {
	switch token := t.(type) {
	case xml.Comment:
		return nil, nil
	case xml.CharData:
		if strings.TrimSpace(string(token)) == "" {
			return nil, nil
		}
		return t, nil
	default:
		return t, nil
	}
}
