package mappers

import (
	"encoding/xml"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNSCollection(t *testing.T) {
	RegisterTestingT(t)

	c := NSCollection{}
	Ω(c.FindPrefix("")).Should(BeNil())
	Ω(c.FindPrefix("qwe")).Should(BeNil())
	Ω(c.FindURI("")).Should(BeNil())
	Ω(c.FindURI("qwe:qwe")).Should(BeNil())

	c.Set("qwe", "qwe:qwe")
	Ω(c.FindPrefix("")).Should(BeNil())
	Ω(c.FindPrefix("qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))
	Ω(c.FindURI("")).Should(BeNil())
	Ω(c.FindURI("qwe:qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))

	c.Set("", "a-s-d")
	Ω(c.FindPrefix("")).Should(Equal(&NSPair{
		Prefix: "",
		URI:    "a-s-d",
	}))
	Ω(c.FindPrefix("qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))
	Ω(c.FindURI("")).Should(BeNil())
	Ω(c.FindURI("a-s-d")).Should(Equal(&NSPair{
		Prefix: "",
		URI:    "a-s-d",
	}))
}

func TestNSStack(t *testing.T) {
	RegisterTestingT(t)

	s := NSStack{}
	Ω(s.FindPrefix("qwe")).Should(BeNil())
	Ω(s.FindURI("qwe:qwe")).Should(BeNil())

	s.Push()
	Ω(s).Should(HaveLen(1))
	Ω(s[0]).Should(Equal(NSCollection{}))

	s.Set("qwe", "qwe:qwe")
	Ω(s.FindPrefix("qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))
	Ω(s.FindURI("qwe:qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))

	s.Push()
	Ω(s).Should(HaveLen(2))
	Ω(s[0]).Should(Equal(NSCollection{}))

	Ω(s.FindPrefix("qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))
	Ω(s.FindURI("qwe:qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))

	s.Set("", "a-s-d")
	Ω(s.FindPrefix("qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))
	Ω(s.FindURI("a-s-d")).Should(Equal(&NSPair{
		Prefix: "",
		URI:    "a-s-d",
	}))

	s.Set("qwe", "a-s-d")
	Ω(s.FindPrefix("qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "a-s-d",
	}))

	s.Pop()
	Ω(s).Should(HaveLen(1))

	Ω(s.FindPrefix("qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))
	Ω(s.FindURI("qwe:qwe")).Should(Equal(&NSPair{
		Prefix: "qwe",
		URI:    "qwe:qwe",
	}))
}

func TestNSNormalizerPushPop(t *testing.T) {
	RegisterTestingT(t)

	var err error

	ns := NSNormalizer{}

	Ω(ns.NS).Should(HaveLen(0))

	_, err = ns.Map(xml.StartElement{})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(ns.NS).Should(HaveLen(1))

	_, err = ns.Map(xml.StartElement{})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(ns.NS).Should(HaveLen(2))

	_, err = ns.Map(xml.EndElement{})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(ns.NS).Should(HaveLen(1))

	_, err = ns.Map(xml.StartElement{})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(ns.NS).Should(HaveLen(2))

	_, err = ns.Map(xml.EndElement{})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(ns.NS).Should(HaveLen(1))

	_, err = ns.Map(xml.EndElement{})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(ns.NS).Should(HaveLen(0))
}

func TestNSNormalizerEmptyStack(t *testing.T) {
	RegisterTestingT(t)

	var res xml.Token
	var err error

	ns := NSNormalizer{}

	res, err = ns.Map(xml.StartElement{
		Name: xml.Name{
			Space: "ns",
			Local: "tag",
		},
		Attr: []xml.Attr{},
	})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(res).Should(Equal(xml.StartElement{
		Name: xml.Name{
			Local: "ns:tag",
		},
		Attr: []xml.Attr{},
	}))

	res, err = ns.Map(xml.EndElement{
		Name: xml.Name{
			Space: "ns",
			Local: "tag",
		},
	})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(res).Should(Equal(xml.EndElement{
		Name: xml.Name{
			Local: "ns:tag",
		},
	}))

	res, err = ns.Map(xml.StartElement{
		Name: xml.Name{
			Space: "http://www.w3.org/2000/svg",
			Local: "g",
		},
		Attr: []xml.Attr{},
	})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(res).Should(Equal(xml.StartElement{
		Name: xml.Name{
			Space: "http://www.w3.org/2000/svg",
			Local: "g",
		},
		Attr: []xml.Attr{},
	}))

	res, err = ns.Map(xml.EndElement{
		Name: xml.Name{
			Space: "http://www.w3.org/2000/svg",
			Local: "g",
		},
	})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(res).Should(Equal(xml.EndElement{
		Name: xml.Name{
			Space: "http://www.w3.org/2000/svg",
			Local: "g",
		},
	}))

}

func TestNSNormalizerDefaultNS(t *testing.T) {
	RegisterTestingT(t)

	ns := NSNormalizer{}
	ns.NS.Push()
	ns.NS.Set("", "http://www.w3.org/2000/svg")

	res, err := ns.Map(xml.StartElement{
		Name: xml.Name{
			Space: "http://www.w3.org/2000/svg",
			Local: "g",
		},
	})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(res).Should(Equal(xml.StartElement{
		Name: xml.Name{
			Local: "g",
		},
		Attr: []xml.Attr{},
	}))

	res, err = ns.Map(xml.StartElement{
		Name: xml.Name{
			Local: "g",
		},
	})
	Ω(err).ShouldNot(HaveOccurred())
	Ω(res).Should(Equal(xml.StartElement{
		Name: xml.Name{
			Local: "g",
		},
		Attr: []xml.Attr{},
	}))
}

func TestNSNormalizerDefNSOverridesOthers(t *testing.T) {
	RegisterTestingT(t)

	ns := NSNormalizer{}
	ns.NS.Push()
	ns.NS.Set("", "http://www.w3.org/2000/svg")
	ns.NS.Push()
	ns.NS.Set("svg", "http://www.w3.org/2000/svg")

	res, err := ns.Map(xml.StartElement{
		Name: xml.Name{
			Space: "http://www.w3.org/2000/svg",
			Local: "g",
		},
	})

	Ω(err).ShouldNot(HaveOccurred())
	Ω(res).Should(Equal(xml.StartElement{
		Name: xml.Name{
			Local: "g",
		},
		Attr: []xml.Attr{},
	}))
}

func TestNSNormalizerNSAttributes(t *testing.T) {
	RegisterTestingT(t)

	ns := NSNormalizer{}
	ns.NS.Push()
	ns.NS.Set("dc", "http://purl.org/dc/elements/1.1/")
	ns.NS.Set("rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#")
	ns.NS.Set("mavrodi", "http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd")
	ns.NS.Push()
	ns.NS.Set("inkscape", "http://www.inkscape.org/namespaces/inkscape")
	ns.NS.Set("sodipodi", "http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd")
	ns.NS.Set("cc", "http://creativecommons.org/ns#")

	res, err := ns.Map(xml.StartElement{
		Name: xml.Name{
			Space: "http://creativecommons.org/ns#",
			Local: "Work",
		},
		Attr: []xml.Attr{
			xml.Attr{
				Name:  xml.Name{Space: "http://www.w3.org/1999/02/22-rdf-syntax-ns#", Local: "about"},
				Value: "",
			},
			xml.Attr{
				Name:  xml.Name{Space: "unknown", Local: "tag"},
				Value: "qwe",
			},
			xml.Attr{
				Name:  xml.Name{Space: "http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd", Local: "docname"},
				Value: "test.svg",
			},
			xml.Attr{
				Name:  xml.Name{Space: "http://www.inkscape.org/namespaces/inkscape", Local: "version"},
				Value: "0.91 r13725",
			},
		},
	})

	Ω(err).ShouldNot(HaveOccurred())
	Ω(res).Should(Equal(xml.StartElement{
		Name: xml.Name{
			Local: "cc:Work",
		},
		Attr: []xml.Attr{
			xml.Attr{
				Name:  xml.Name{Local: "rdf:about"},
				Value: "",
			},
			xml.Attr{
				Name:  xml.Name{Local: "unknown:tag"},
				Value: "qwe",
			},
			xml.Attr{
				Name:  xml.Name{Local: "sodipodi:docname"},
				Value: "test.svg",
			},
			xml.Attr{
				Name:  xml.Name{Local: "inkscape:version"},
				Value: "0.91 r13725",
			},
		},
	}))
}
