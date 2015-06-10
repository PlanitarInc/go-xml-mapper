package xmlproc

import "encoding/xml"

// Mapper is an interface encapsulating the actual transformation
// (or other work) to be performed to a XML document.
// A mapper is called for every token read from a document.
// If an error is returned, the processing of a document is aborted with the
// error.
// If the returned token is nil, it's being ignored and not present in
// a output XML document.
type Mapper interface {
	Map(xml.Token) (xml.Token, error)
}
