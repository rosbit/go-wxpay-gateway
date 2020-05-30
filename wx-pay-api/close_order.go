// +build gateway

package wxpay

func postClose(orderId string, xml []byte, isSandbox bool, apiKey string) (recv []byte, err error) {
	orderclose_url := _GetApiUrl(UT_ORDER_CLOSE, isSandbox)
	if recv, err = _CallWxAPI(orderclose_url, "POST", xml); err != nil {
		return
	}

	_, err = parseXmlResult(recv, apiKey)
	return
}

func CloseOrder(
	appId     string,
	mchId     string,
	mchApiKey string,
	orderId   string,
	isSandbox bool,
) (xmlstr, recv []byte, err error) {
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

	xmlstr = xml.ToXML()
	recv, err = postClose(orderId, xmlstr, isSandbox, mchApiKey)
	return
}
