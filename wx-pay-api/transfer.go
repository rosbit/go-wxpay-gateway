// +build gateway

package wxpay

import (
	"go-wxpay-gateway/sign"
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
) (res *TransferResult, xmlstr, recv []byte, err error) {
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
	signature := sign.CreateSignature(sign.MD5, tags, mchApiKey)
	addTag(xml, tags, "sign", signature, false)

	xmlstr = xml.ToXML()
	// fmt.Printf("xml: %s\n", string(xmlstr))

	res, recv, err = postTransfer(tradeNo, mchApiKey, certFile, keyFile, xmlstr, isSandbox)
	return
}

func postTransfer(tradeNo, apiKey, certFile, keyFile string, xml []byte, isSandbox bool) (res *TransferResult, recv []byte, err error) {
	transfer_url := _GetApiUrl(UT_TRANSFER, isSandbox)
	if recv, err = _CallSecureWxAPI(transfer_url, "POST", xml, certFile, keyFile); err != nil {
		return
	}

	res, err = ParseTransferResultBody("transfer result", recv, apiKey)
	return
}
