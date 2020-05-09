// +build gateway

package wxpay

func postQueryTransfer(tradeNo, apiKey, certFile, keyFile string, xml []byte, isSandbox bool) (*QueryTransferResult, error) {
	_paymentLog.Printf("[query-transfer] 1. ### Before Querying tradeNo #%s: %s\n", tradeNo, string(xml))
	transfer_query_url := _GetApiUrl(UT_TRANSFER_QUERY, isSandbox)
	content, err := _CallSecureWxAPI(transfer_query_url, "POST", xml, certFile, keyFile)
	if err != nil {
		_paymentLog.Printf("[query-transfer] 2. --- Query tradeNo #%s failed: %v\n", tradeNo, err)
		return nil, err
	}
	_paymentLog.Printf("[query-transfer] 2. +++ Result of querying tradeNo #%s: %s\n", tradeNo, string(content))

	return ParseQueryTransferResult("query-transfer-result", content, apiKey)
}

func QueryTransfer(
	appId     string,
	mchId     string,
	mchApiKey string,
	tradeNo   string,
	certFile  string,
	keyFile   string,
	isSandbox bool,
) (*QueryTransferResult, error) {
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
	addTag(xml, params, "partner_trade_no", tradeNo, false)
	addTag(xml, params, "nonce_str", string(_GetRandomBytes(32)), false)
	addTag(xml, params, "sign_type", "MD5", false)

	// sign
	signature := createMd5Signature(params, mchApiKey)
	addTag(xml, params, "sign", signature, false)

	xmlstr := xml.toXML()

	return postQueryTransfer(tradeNo, mchApiKey, certFile, keyFile, xmlstr, isSandbox)
}
