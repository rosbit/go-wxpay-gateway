// +build gateway

package wxpay

import (
	"fmt"
)

type FnRefund func(appId,mchId,mchApiKey,id,refundId string,totalFee,refundFee int, refundReason,refundNotify,certFile,keyFile string, isSandbox bool) (res *RefundResultParams, sent, recv []byte, err error)

func RefundByTransactionId(
	appId     string,
	mchId     string,
	mchApiKey string,
	transactionId string, // 微信生成的订单号
	refundId      string, // 商户自己的退款编号 char(32)
	totalFee  int, // 订单总金额，单位为分
	refundFee int, // 退款总金额，订单总金额，单位为分
	refundReason string, // 退款原因
	refundNotify string, // 异步接收微信支付退款结果通知的回调地址
	certFile string,
	keyFile  string,
	isSandbox bool,
) (res *RefundResultParams, sent, recv []byte, err error) {
	return refund(appId, mchId, mchApiKey, transactionId, "", refundId, totalFee, refundFee, refundReason, refundNotify, certFile, keyFile, isSandbox)
}

func RefundByOrderId(
	appId     string,
	mchId     string,
	mchApiKey string,
	orderId   string, // 商户自己的唯一订单号
	refundId  string, // 商户自己的退款编号 char(32)
	totalFee  int, // 订单总金额，单位为分
	refundFee int, // 退款总金额，订单总金额，单位为分
	refundReason string, // 退款原因
	refundNotify string, // 异步接收微信支付退款结果通知的回调地址
	certFile string,
	keyFile  string,
	isSandbox bool,
) (res *RefundResultParams, sent, recv []byte, err error) {
	return refund(appId, mchId, mchApiKey, "", orderId, refundId, totalFee, refundFee, refundReason, refundNotify, certFile, keyFile, isSandbox)
}

func refund(
	appId     string,
	mchId     string,
	mchApiKey string,
	transactionId string, // 微信生成的订单号
	orderId       string, // 商户自己的唯一订单号
	refundId  string, // 商户自己的退款编号 char(32)
	totalFee  int, // 订单总金额，单位为分
	refundFee int, // 退款总金额，订单总金额，单位为分
	refundReason string, // 退款原因
	refundNotify string, // 异步接收微信支付退款结果通知的回调地址
	certFile string,
	keyFile  string,
	isSandbox bool,
) (res *RefundResultParams, xmlstr, recv []byte, err error) {
	tags := make(map[string]string)
	xml := newXmlGenerator("xml")
	addTag(xml, tags, "appid",       appId,      false)
	addTag(xml, tags, "mch_id",      mchId,      false)
	addTag(xml, tags, "nonce_str",   string(_GetRandomBytes(32)), false)
	if transactionId != "" {
		addTag(xml, tags, "transaction_id", transactionId, false)
	} else {
		addTag(xml, tags, "out_trade_no",   orderId,       false)
	}
	addTag(xml, tags, "out_refund_no", refundId,  false)
	addTag(xml, tags, "total_fee",   fmt.Sprintf("%d", totalFee),  false)
	addTag(xml, tags, "refund_fee",  fmt.Sprintf("%d", refundFee), false)
	addTag(xml, tags, "refund_desc", refundReason,true)
	//if !isSandbox {
		addTag(xml, tags, "notify_url",  refundNotify,   true)
	//}
	// sign
	signature := createMd5Signature(tags, mchApiKey)
	addTag(xml, tags, "sign", signature, false)

	xmlstr = xml.ToXML()
	// fmt.Printf("xml: %s\n", string(xmlstr))

	res, recv, err = postRefund(refundId, mchApiKey, certFile, keyFile, xmlstr, isSandbox)
	return
}

func postRefund(refundId, apiKey, certFile, keyFile string, xml []byte, isSandbox bool) (res *RefundResultParams, recv []byte, err error) {
	refund_url := _GetApiUrl(UT_REFUND, isSandbox)
	if recv, err = _CallSecureWxAPI(refund_url, "POST", xml, certFile, keyFile); err != nil {
		return
	}

	res, err = ParseRefundResultBody("refund result", recv, apiKey)
	return
}

