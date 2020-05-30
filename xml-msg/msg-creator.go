package xmlmsg

import (
	"github.com/beevik/etree"
	"bytes"
)

type XmlGenerator struct {
	doc  *etree.Document
	root *etree.Element
}

func NewXmlGenerator(rootTag string) *XmlGenerator {
	doc := etree.NewDocument()
	root := doc.CreateElement(rootTag)
	return &XmlGenerator{doc, root}
}

func (xml *XmlGenerator) AddTag(tagName string, value string) {
	tag := xml.root.CreateElement(tagName)
	tag.SetText(value)
}

func (xml *XmlGenerator) ToXML() ([]byte) {
	buf := &bytes.Buffer{}
	xml.doc.Indent(2)
	xml.doc.WriteTo(buf)
	return buf.Bytes()
}

func AddTag(xml *XmlGenerator, tags map[string]string, tagName string, value string, optionalIfNull bool) {
	if optionalIfNull && value == "" {
		return
	}

	xml.AddTag(tagName, value)
	if value != "" {
		tags[tagName] = value
	}
}

