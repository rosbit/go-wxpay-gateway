// +build gateway

package wxpay

func postQuery(transactionId, orderId string, xml []byte, isSandbox bool, apiKey string) (res INotifyParams, recv []byte, err error) {
	orderquery_url := _GetApiUrl(UT_ORDER_QUERY, isSandbox)
	if recv, err = _CallWxAPI(orderquery_url, "POST", xml); err != nil {
		return
	}

	res, err = ParsePayNotifyBody("query-order-result", recv, apiKey)
	return
}

func queryOrder(
	appId     string,
	mchId     string,
	mchApiKey string,
	transactionId string,
	orderId   string,
	isSandbox bool,
) (res INotifyParams, xmlstr, recv []byte, err error) {
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
	if transactionId != "" {
		addTag(xml, params, "transaction_id", transactionId, false)
	} else {
		addTag(xml, params, "out_trade_no", orderId, false)
	}
	addTag(xml, params, "nonce_str", string(_GetRandomBytes(32)), false)
	addTag(xml, params, "sign_type", "MD5", false)

	// sign
	signature := createMd5Signature(params, mchApiKey)
	addTag(xml, params, "sign", signature, false)

	xmlstr = xml.ToXML()

	res, recv, err = postQuery(transactionId, orderId, xmlstr, isSandbox, mchApiKey)
	return
}

type FnQueryOrder func(
	appId     string,
	mchId     string,
	mchApiKey string,
	id        string,
	isSandbox bool,
) (res INotifyParams, sent, recv []byte, err error)

func QueryByOrderId(
	appId string,
	mchId string,
	mchApiKey string,
	orderId string,
	isSandbox bool,
) (INotifyParams, []byte, []byte, error) {
	return queryOrder(appId, mchId, mchApiKey, "", orderId, isSandbox)
}

func QueryByTransactionId(
	appId string,
	mchId string,
	mchApiKey string,
	transactionId string,
	isSandbox bool,
) (INotifyParams, []byte, []byte, error) {
	return queryOrder(appId, mchId, mchApiKey, transactionId, "", isSandbox)
}
