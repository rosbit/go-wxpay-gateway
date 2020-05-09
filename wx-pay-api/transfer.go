// +build gateway

package wxpay

import (
	"fmt"
)

func Transfer(
	appId     string,
	mchId     string,
	mchApiKey string,
	tradeNo   string,
	openId    string,
	userName  string,
	amount    int,
	desc      string,
	ip        string,
	certFile  string,
	keyFile   string,
	isSandbox bool,
) (*TransferResult, error) {
	tags := make(map[string]string)
	xml := newXmlGenerator("xml")
	addTag(xml, tags, "mch_appid",   appId,      false)
	addTag(xml, tags, "mchid",       mchId,      false)
	addTag(xml, tags, "nonce_str",   string(_GetRandomBytes(32)), false)
	addTag(xml, tags, "partner_trade_no",   tradeNo,       false)
	addTag(xml, tags, "openid",      openId,  false)
	if userName == "" {
		addTag(xml, tags, "check_name",  "NO_CHECK",  false)
	} else {
		addTag(xml, tags, "check_name",  "FORCE_CHECK",  false)
		addTag(xml, tags, "re_user_name", userName,      false)
	}
	addTag(xml, tags, "amount",   fmt.Sprintf("%d", amount),  false)
	addTag(xml, tags, "desc",     desc,                       false)
	addTag(xml, tags, "spbill_create_ip",  ip,                false)
	// sign
	signature := createMd5Signature(tags, mchApiKey)
	addTag(xml, tags, "sign", signature, false)

	xmlstr := xml.toXML()
	// fmt.Printf("xml: %s\n", string(xmlstr))

	return postTransfer(tradeNo, mchApiKey, certFile, keyFile, xmlstr, isSandbox)
}

func postTransfer(tradeNo, apiKey, certFile, keyFile string, xml []byte, isSandbox bool) (*TransferResult, error) {
	_paymentLog.Printf("[transfer] 1. ### Before POSTing transfer #%s: %s\n", tradeNo, string(xml))
	transfer_url := _GetApiUrl(UT_TRANSFER, isSandbox)
	content, err := _CallSecureWxAPI(transfer_url, "POST", xml, certFile, keyFile)
	if err != nil {
		_paymentLog.Printf("[transfer] 2. --- POST transfer #%s failed: %v\n", tradeNo, err)
		return  nil, err
	}
	_paymentLog.Printf("[transfer] 2. +++ Result of POSTing transfer #%s: %s\n", tradeNo, string(content))

	return ParseTransferResultBody("transfer result", content, apiKey)
}
