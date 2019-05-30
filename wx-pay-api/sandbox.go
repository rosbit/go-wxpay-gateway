// +build getsandbox

package wxpay

import (
	"fmt"
	"log"
)

const getsignkey_url = "https://api.mch.weixin.qq.com/sandboxnew/pay/getsignkey"

func GetSandbox(appId string, mchId string, appKey string) (string, error) {
	tags := make(map[string]string)
	xml := newXmlGenerator("xml")
	addTag(xml, tags, "appid",  appId, false)
	addTag(xml, tags, "mch_id", mchId, false)
	addTag(xml, tags, "sign_type", "MD5", false)
	addTag(xml, tags, "nonce_str", string(_GetRandomBytes(32)), false)

	// sign
	signature := createMd5Signature(tags, appKey)
	addTag(xml, tags, "sign", signature, false)

	xmlstr := xml.toXML()
	log.Printf("[sandbox] try to get sandbox appKey: %s\n", string(xmlstr))

	res, err := _CallWxAPI(getsignkey_url, "POST", xmlstr)
	if err != nil {
		return "", err
	}
	log.Printf("[sandbox] restult: %s\n", string(res))
	return parseSandboxKey(res)
}

func parseSandboxKey(body []byte) (string, error) {
	res, err := xml2map(body)
	if err != nil {
		return "", err
	}

	if sandbox_signkey, ok := res["sandbox_signkey"]; ok {
		return sandbox_signkey, nil
	}
	return "", fmt.Errorf("no sandbox_signkey")
}
