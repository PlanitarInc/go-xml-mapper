package mappers

import (
	"encoding/xml"
	"regexp"
)

type NSPair struct {
	Prefix string
	URI    string
}

type NSCollection struct {
	Pairs []NSPair
}

func (m NSCollection) FindPrefix(name string) *NSPair {
	for i := range m.Pairs {
		if m.Pairs[i].Prefix == name {
			return &m.Pairs[i]
		}
	}
	return nil
}

func (m NSCollection) FindURI(name string) *NSPair {
	for i := range m.Pairs {
		if m.Pairs[i].URI == name {
			return &m.Pairs[i]
		}
	}
	return nil
}

func (m *NSCollection) Set(alias, space string) {
	m.Pairs = append(m.Pairs, NSPair{Prefix: alias, URI: space})
}

type NSStack []NSCollection

func (s NSStack) FindPrefix(name string) *NSPair {
	for i := range s {
		if p := s[i].FindPrefix(name); p != nil {
			return p
		}
	}
	return nil
}

func (s NSStack) FindURI(name string) *NSPair {
	for i := range s {
		if p := s[i].FindURI(name); p != nil {
			return p
		}
	}
	return nil
}

func (s *NSStack) Set(alias, space string) {
	(*s)[0].Set(alias, space)
}

func (s *NSStack) Push() {
	*s = append(NSStack{NSCollection{}}, (*s)...)
}

func (s *NSStack) Pop() {
	*s = (*s)[1:]
}

// NSNormalizer manages the namespace stack for a processed document,
// and normalizes all the namespaces.
// The mapper comes to fix the issues with XML namespaces in a standard
// golang xml library.
// See https://github.com/golang/go/search?q=namespace&type=Issues&utf8=%E2%9C%93
type NSNormalizer struct {
	NS NSStack
}

func (p NSNormalizer) SetNSAlias(name *xml.Name) {
	if name.Space == "" {
		return
	}

	if ns := p.NS.FindPrefix(""); ns != nil && ns.URI == name.Space {
		name.Space = ""
		return
	}

	if ns := p.NS.FindURI(name.Space); ns != nil {
		name.Local = ns.Prefix + ":" + name.Local
		name.Space = ""
		return
	}

	if ok, _ := regexp.MatchString(`[^a-zA-Z]`, name.Space); !ok {
		name.Local = name.Space + ":" + name.Local
		name.Space = ""
		return
	}
}

func (p *NSNormalizer) Map(t xml.Token) (xml.Token, error) {
	switch token := t.(type) {
	case xml.StartElement:
		token = token.Copy()

		p.NS.Push()
		for _, a := range token.Attr {
			if a.Name.Space == "xmlns" {
				p.NS.Set(a.Name.Local, a.Value)
			} else if a.Name.Space == "" && a.Name.Local == "xmlns" {
				p.NS.Set("", a.Value)
			}
		}

		p.SetNSAlias(&token.Name)
		for i := range token.Attr {
			if token.Attr[i].Name.Space == "xmlns" {
				token.Attr[i].Name.Space = ""
				token.Attr[i].Name.Local = "xmlns:" + token.Attr[i].Name.Local
			} else if token.Attr[i].Name.Space != "" {
				p.SetNSAlias(&token.Attr[i].Name)
			}
		}

		//          fmt.Printf("    ns %#v\n", p.NS)
		return token, nil

	case xml.EndElement:

		p.SetNSAlias(&token.Name)
		p.NS.Pop()

		//          fmt.Printf("    end %#v\n", token)
		return token, nil

	default:
		return t, nil
	}
}
