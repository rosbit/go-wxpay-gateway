package wxpay

import (
	"github.com/beevik/etree"
	"crypto/md5"
	"bytes"
	"sort"
	"io"
	"fmt"
)

type xmlGenerator struct {
	doc  *etree.Document
	root *etree.Element
}

func newXmlGenerator(rootTag string) *xmlGenerator {
	doc := etree.NewDocument()
	root := doc.CreateElement(rootTag)
	return &xmlGenerator{doc, root}
}

func (xml *xmlGenerator) addTag(tagName string, value string) {
	tag := xml.root.CreateElement(tagName)
	tag.SetText(value)
}

func (xml *xmlGenerator) toXML() ([]byte) {
	buf := &bytes.Buffer{}
	xml.doc.Indent(2)
	xml.doc.WriteTo(buf)
	return buf.Bytes()
}

func addTag(xml *xmlGenerator, tags map[string]string, tagName string, value string, optionalIfNull bool) {
	if optionalIfNull && value == "" {
		return
	}

	xml.addTag(tagName, value)
	if value != "" {
		tags[tagName] = value
	}
}

func createMd5Signature(params map[string]string, apiKey string) string {
	// sort keys
	keys := make([]string, len(params))
	i := 0
	for k := range params {
		keys[i] = k
		i += 1
	}
	sort.Strings(keys)

	stringA := md5.New()
	// create stringA="key1=val1&"
	first := true // 是否第一个非空值参数
	for _, key := range keys {
		v, _ := params[key]
		if v != "" {
			if first {
				first = false
			} else {
				io.WriteString(stringA, "&")
			}
			io.WriteString(stringA, key)
			io.WriteString(stringA, "=")
			io.WriteString(stringA, v)
		}
	}

	io.WriteString(stringA, "&key=")
	io.WriteString(stringA, apiKey)

	// MD5
	sign := stringA.Sum(nil)
	return fmt.Sprintf("%X", sign)
}

func xml2mapWithRoot(body []byte, rootName string) (map[string]string, error) {
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

	if return_code, ok := res["return_code"]; ok {
		if return_code != "SUCCESS" {
			// return_code is "FAIL"
			if return_msg, ok := res["return_msg"]; ok {
				err = fmt.Errorf(return_msg)
			} else {
				err = fmt.Errorf("no return_msg found")
			}
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no return_code")
	}

	return res, nil
}

func xml2map(body []byte) (map[string]string, error) {
	return xml2mapWithRoot(body, "xml")
}

func parseXmlResult(body []byte, apiKey string) (map[string]string, error) {
	res, err := xml2map(body)
	if err != nil {
		return nil, err
	}

	// check signautre
	if sign, ok := res["sign"]; !ok {
		return nil, fmt.Errorf("no signature in result")
	} else {
		delete(res, "sign")
		createdSign := createMd5Signature(res, apiKey)
		if sign != createdSign {
			return nil, fmt.Errorf("signature not matched: %s != %s", sign, createdSign)
		}
	}

	return res, nil
}
