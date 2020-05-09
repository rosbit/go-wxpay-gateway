// +build gateway

package wxpay

func postClose(orderId string, xml []byte, isSandbox bool, apiKey string) error {
	_paymentLog.Printf("[close-order] 1. ### Before Closing order %s: %s\n", orderId, string(xml))
	orderclose_url := _GetApiUrl(UT_ORDER_CLOSE, isSandbox)
	content, err := _CallWxAPI(orderclose_url, "POST", xml)
	if err != nil {
		_paymentLog.Printf("[close-order] 2. --- Close order %s failed: %v\n", orderId, err)
		return err
	}
	_paymentLog.Printf("[close-order] 2. +++ Result of closing order %s: %s\n", orderId, string(content))

	_, err = parseXmlResult(content, apiKey)
	return err
}

func CloseOrder(
	appId     string,
	mchId     string,
	mchApiKey string,
	orderId   string,
	isSandbox bool,
) error {
	/*
	if isSandbox {
		var err error
		if mchApiKey, err = GetSandbox(appId, mchId, mchApiKey); err != nil {
			return nil, err
		}
	}*/
	xml := newXmlGenerator("xml")
	params := make(map[string]string)

	addTag(xml, params, "appid", appId, false)
	addTag(xml, params, "mch_id", mchId, false)
	addTag(xml, params, "out_trade_no", orderId, false)
	addTag(xml, params, "nonce_str", string(_GetRandomBytes(32)), false)
	addTag(xml, params, "sign_type", "MD5", false)

	// sign
	signature := createMd5Signature(params, mchApiKey)
	addTag(xml, params, "sign", signature, false)

	xmlstr := xml.toXML()

	return postClose(orderId, xmlstr, isSandbox, mchApiKey)
}
