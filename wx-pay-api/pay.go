// +build gateway

package wxpay

import (
	"fmt"
)

const (
	SANDBOX_FEE = 101
)

func postOrder(orderId string, apiKey string, xml []byte, isSandbox bool) (map[string]string, error) {
	_paymentLog.Printf("[payment] 1. ### Before POSTing order #%s: %s\n", orderId, string(xml))
	unifiedorder_url := _GetApiUrl(UT_UNIFIED_ORDER, isSandbox)
	content, err := _CallWxAPI(unifiedorder_url, "POST", xml)
	if err != nil {
		_paymentLog.Printf("[payment] 2. --- POST order #%s failed: %v\n", orderId, err)
		return  nil, err
	}
	_paymentLog.Printf("[payment] 2. +++ Result of POSTing order #%s: %s\n", orderId, string(content))

	return parseXmlResult(content, apiKey)
}

func payOrder(
	appId string,
	mchId string,
	mchApiKey string,
	deviceInfo string,
	payBody string,
	cbParams string,
	orderId string,
	fee int,
	ip string,
	notifyUrl string,
	tradeType string,
	productId string,
	openId string,
	sceneInfo []byte,
	isSandbox bool,
) (prepay_id string, res map[string]string, err error) {
	tags := make(map[string]string)
	xml := newXmlGenerator("xml")
	addTag(xml, tags, "appid",       appId,      false)
	addTag(xml, tags, "mch_id",      mchId,      false)
	addTag(xml, tags, "device_info", deviceInfo, true)
	addTag(xml, tags, "nonce_str",   string(_GetRandomBytes(32)), false)
	addTag(xml, tags, "body",        payBody, true)
	addTag(xml, tags, "attach",      cbParams,    true)
	addTag(xml, tags, "out_trade_no",orderId,     false)
	addTag(xml, tags, "total_fee",   fmt.Sprintf("%d", fee), false)
	addTag(xml, tags, "spbill_create_ip", ip,     false)
	addTag(xml, tags, "notify_url",  notifyUrl,   false)
	addTag(xml, tags, "trade_type",  tradeType,   false)
	addTag(xml, tags, "product_id",  productId,   tradeType != "NATIVE")
	addTag(xml, tags, "openid",      openId,      tradeType != "JSAPI")
	if sceneInfo != nil {
		addTag(xml, tags, "scene_info", string(sceneInfo), false)
	}
	// sign
	signature := createMd5Signature(tags, mchApiKey)
	addTag(xml, tags, "sign", signature, false)

	xmlstr := xml.toXML()
	// fmt.Printf("xml: %s\n", string(xmlstr))

	if res, err = postOrder(orderId, mchApiKey, xmlstr, isSandbox); err != nil {
		return "", nil, err
	}

	// return_code is "SUCCESS", then check result_code
	if result_code, ok := res["result_code"]; !ok {
		return "", nil, fmt.Errorf("no result_code")
	} else {
		if result_code != "SUCCESS" {
			err_code_des, ok := res["err_code_des"]
			if !ok {
				return "", nil, fmt.Errorf("no err_code_des")
			} else {
				return "", nil, fmt.Errorf("err_code_des: %s", err_code_des)
			}
		}
	}

	// get prepay_id
	var ok bool
	if prepay_id, ok = res["prepay_id"]; !ok {
		return "", nil, fmt.Errorf("no prepay_id")
	}
	return
}

