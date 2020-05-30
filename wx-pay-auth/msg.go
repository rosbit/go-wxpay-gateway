// + build gateway

package oauth

import (
	"go-wxpay-gateway/xml-msg"
	"go-wxpay-gateway/sign"
	"fmt"
)

func xml2mapWithRoot(body []byte, rootName string) (res map[string]string, err error) {
	if res, err = xmlmsg.Xml2mapWithRoot(body, rootName); err != nil {
		return
	}
	err = checkXmlMsgResult(res)
	return
}

func parseXmlResult(body []byte, apiKey string, signType string) (res map[string]string, err error) {
	if res, err = xml2mapWithRoot(body, "xml"); err != nil {
		return
	}
	err = sign.CheckSignature(signType, res, apiKey)
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
