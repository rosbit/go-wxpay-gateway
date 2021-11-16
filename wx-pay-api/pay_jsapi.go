// +build gateway

package wxpay

import (
	"fmt"
	"time"
)

func JSAPIPay(
	appId  string,
	mchId  string,
	mchApiKey string,
	receipt bool,
	payBody   string,
	cbParams  string,
	orderId string,
	fee     int,
	ip      string,
	notifyUrl string,
	openId    string,
	isSandbox bool,
) (prepay_id string, reqJSAPI map[string]string, sent, recv []byte, err error) {
	if isSandbox {
		/*
		if mchApiKey, err = GetSandbox(appId, mchId, mchApiKey); err != nil {
			return
		}*/
		fee = SANDBOX_FEE
	}
	if prepay_id, _, sent, recv, err = payOrder(appId, mchId, mchApiKey, receipt, "WEB", payBody, cbParams, orderId, fee, ip, notifyUrl, "JSAPI", "", openId, nil, isSandbox); err != nil {
		return
	}

	reqJSAPI = CreateJSAPIParams(appId, mchApiKey, prepay_id)
	return
}

func CreateJSAPIParams(appId string, mchApiKey string, prepay_id string) map[string]string {
	params := make(map[string]string, 6)
	params["appId"]     = appId
	params["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	params["nonceStr"]  = string(_GetRandomBytes(32))
	params["package"]   = fmt.Sprintf("prepay_id=%s", prepay_id)
	params["signType"]  = "MD5"

	paySign := createMd5Signature(params, mchApiKey)
	params["paySign"]   = paySign
	return params
}
