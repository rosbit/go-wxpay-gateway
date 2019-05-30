// +build gateway

package wxpay

import (
	"fmt"
	"time"
)

func JSAPIPay(
	appId  string,
	mchId  string,
	mchAppKey string,
	payBody   string,
	cbParams  string,
	orderId string,
	fee     int,
	ip      string,
	notifyUrl string,
	openId    string,
	isSandbox bool,
) (prepay_id string, reqJSAPI map[string]string, err error) {
	if isSandbox {
		/*
		if mchAppKey, err = GetSandbox(appId, mchId, mchAppKey); err != nil {
			return
		}*/
		fee = SANDBOX_FEE
	}
	prepay_id, _, err = payOrder(appId, mchId, mchAppKey, "WEB", payBody, cbParams, orderId, fee, ip, notifyUrl, "JSAPI", "", openId, nil, isSandbox)
	if err != nil {
		_paymentLog.Printf("[JSAPI-payment] 3. --- %v\n", err)
		return "", nil, err
	}

	return prepay_id, CreateJSAPIParams(appId, mchAppKey, prepay_id), nil
}

func CreateJSAPIParams(appId string, mchAppKey string, prepay_id string) map[string]string {
	params := make(map[string]string, 6)
	params["appId"]     = appId
	params["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	params["nonceStr"]  = string(_GetRandomBytes(32))
	params["package"]   = fmt.Sprintf("prepay_id=%s", prepay_id)
	params["signType"]  = "MD5"

	paySign := createMd5Signature(params, mchAppKey)
	params["paySign"]   = paySign
	_paymentLog.Printf("[payment] ### JSAPI payment params: %v\n", params)
	return params
}
