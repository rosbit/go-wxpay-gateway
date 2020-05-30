// +build gateway

package wxpay

func postQueryTransfer(tradeNo, apiKey, certFile, keyFile string, xml []byte, isSandbox bool) (res *QueryTransferResult,recv []byte, err error) {
	transfer_query_url := _GetApiUrl(UT_TRANSFER_QUERY, isSandbox)
	if recv, err = _CallSecureWxAPI(transfer_query_url, "POST", xml, certFile, keyFile); err != nil {
		return
	}

	res, err = ParseQueryTransferResult("query-transfer-result", recv, apiKey)
	return
}

func QueryTransfer(
	appId     string,
	mchId     string,
	mchApiKey string,
	tradeNo   string,
	certFile  string,
	keyFile   string,
	isSandbox bool,
) (res *QueryTransferResult, xmlstr, recv []byte, err error) {
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

	xmlstr = xml.ToXML()

	res, recv, err = postQueryTransfer(tradeNo, mchApiKey, certFile, keyFile, xmlstr, isSandbox)
	return
}
