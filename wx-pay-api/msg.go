// + build gateway

package wxpay

import (
	"go-wxpay-gateway/xml-msg"
	"go-wxpay-gateway/sign"
	"fmt"
)

func newXmlGenerator(rootTag string) *xmlmsg.XmlGenerator {
	return xmlmsg.NewXmlGenerator(rootTag)
}

func addTag(xml *xmlmsg.XmlGenerator, tags map[string]string, tagName string, value string, optionalIfNull bool) {
	xmlmsg.AddTag(xml, tags, tagName, value, optionalIfNull)
}

func createMd5Signature(tags map[string]string, apiKey string) string {
	return sign.CreateSignature(sign.MD5, tags, apiKey)
}

func xml2mapWithRoot(body []byte, rootName string) (res map[string]string, err error) {
	if res, err = xmlmsg.Xml2mapWithRoot(body, rootName); err != nil {
		return
	}
	err = checkXmlMsgResult(res)
	return
}

func xml2map(body []byte) (res map[string]string, err error) {
	return xml2mapWithRoot(body, "xml")
}

func parseXmlResult(body []byte, apiKey string) (res map[string]string, err error) {
	if res, err = xml2map(body); err != nil {
		return
	}
	err = sign.CheckSignature(sign.MD5, res, apiKey)
	return
}

func checkXmlMsgResult(res map[string]string) (error) {
	if return_code, ok := res["return_code"]; ok {
		if return_code != "SUCCESS" {
			// return_code is "FAIL"
			if return_msg, ok := res["return_msg"]; ok {
				return fmt.Errorf(return_msg)
			}
			return fmt.Errorf("no return_msg found")
		}
	} else {
		return fmt.Errorf("no return_code")
	}
	return nil
}
