package xmlmsg

import (
	"github.com/beevik/etree"
	"fmt"
)

func Xml2mapWithRoot(body []byte, rootName string) (map[string]string, error) {
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(body)
	if err != nil {
		return nil, err
	}
	root := doc.SelectElement(rootName)
	if root == nil {
		return nil, fmt.Errorf("no root named \"%s\" found", rootName)
	}
	tags := root.ChildElements()
	if tags == nil {
		return nil, fmt.Errorf("no xml children")
	}

	res := make(map[string]string, len(tags))
	for _, tag := range tags {
		res[tag.Tag] = tag.Text()
	}

	return res, nil
}

func Xml2map(body []byte) (map[string]string, error) {
	return Xml2mapWithRoot(body, "xml")
}
