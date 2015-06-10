// Package xmlprocessor is intended to process and modify XML files.
package xmlprocessor

import (
	"encoding/xml"
	"io"

	"github.com/PlanitarInc/go-xml-processor/mappers"
)

// Processor type encapsulates the main logic of processing an XML file.
type Processor struct {
	mappers []Mapper
}

// Create a Processor with predefined set of mappers:
// mapper.Pruner and mapper.NSNormalizer.
func NewDefaultProcessor() *Processor {
	return &Processor{
		mappers: []Mapper{
			&mappers.Pruner{},
			&mappers.NSNormalizer{},
		},
	}
}

// AddMapper appends a mapper to processor's map list.
// No check for duplicates is performed;
// if the mapper is already present in the list,
// it would be appended anyway.
func (p *Processor) AddMapper(m Mapper) {
	p.mappers = append(p.mappers, m)
}

// RemMapper removes a first occurence of a given mapper
// in processor's map list.
func (p *Processor) RemMapper(m Mapper) {
	for i := range p.mappers {
		if p.mappers[i] == m {
			p.mappers = append(p.mappers[:i], p.mappers[i+1:]...)
		}
	}
}

// ProcessStreams reads an XML file from src,
// processes it by applying the mappers, and
// writes the resulting XML file to dst.
func (p Processor) ProcessStreams(dst io.Writer, src io.Reader) error {
	e := xml.NewEncoder(dst)
	e.Indent("", "  ")
	defer e.Flush()
	d := xml.NewDecoder(src)
	return p.Process(e, d)
}

// Process reads XML tokens using the provided decoder,
// processes them by applying the mappers, and
// writes the resulting XML token using the provided encoder.
func (p Processor) Process(e *xml.Encoder, d *xml.Decoder) error {
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		for _, m := range p.mappers {
			if token, err := m.Map(t); err != nil {
				return err
			} else if t == nil {
				break
			} else {
				t = token
			}
		}

		if t == nil {
			continue
		}

		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}
	return nil
}
