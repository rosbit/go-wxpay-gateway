// +build gateway

package wxpay

import (
	"fmt"
)

func postQuery(transactionId, orderId string, xml []byte, isSandbox bool, apiKey string) (*NotifyParams, error) {
	var q string
	if transactionId != "" {
		q = fmt.Sprintf("transacton #%s", transactionId)
	} else {
		q = fmt.Sprintf("order #%s", orderId)
	}
	_paymentLog.Printf("[query] 1. ### Before Querying %s: %s\n", q, string(xml))
	orderquery_url := _GetApiUrl(UT_ORDER_QUERY, isSandbox)
	content, err := _CallWxAPI(orderquery_url, "POST", xml)
	if err != nil {
		_paymentLog.Printf("[query] 2. --- Query %s failed: %v\n", q, err)
		return nil, err
	}
	_paymentLog.Printf("[query] 2. +++ Result of querying %s: %s\n", q, string(content))

	return ParsePayNotifyBody("query-order-result", content, apiKey), nil
}

func queryOrder(
	appId     string,
	mchId     string,
	mchApiKey string,
	transactionId string,
	orderId   string,
	isSandbox bool,
) (*NotifyParams, error) {
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

	xmlstr := xml.toXML()

	return postQuery(transactionId, orderId, xmlstr, isSandbox, mchApiKey)
}

type FnQueryOrder func(
	appId     string,
	mchId     string,
	mchApiKey string,
	id        string,
	isSandbox bool,
) (res *NotifyParams, err error)

func QueryByOrderId(
	appId string,
	mchId string,
	mchApiKey string,
	orderId string,
	isSandbox bool,
) (*NotifyParams, error) {
	return queryOrder(appId, mchId, mchApiKey, "", orderId, isSandbox)
}

func QueryByTransactionId(
	appId string,
	mchId string,
	mchApiKey string,
	transactionId string,
	isSandbox bool,
) (*NotifyParams, error) {
	return queryOrder(appId, mchId, mchApiKey, transactionId, "", isSandbox)
}
