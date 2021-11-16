// +build gateway

package wxpay

import (
	"fmt"
	"time"
)

func AppPay(
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
	isSandbox bool,
) (prepayId string, reqApp map[string]string, sent, recv []byte, err error) {
	if isSandbox {
		/*
		if mchApiKey, err = GetSandbox(appId, mchId, mchApiKey); err != nil {
			return
		}*/
		fee = SANDBOX_FEE
	}
	prepayId, _, sent, recv, err = payOrder(appId, mchId, mchApiKey, receipt, "APP", payBody, cbParams, orderId, fee, ip, notifyUrl, "APP", "", "", nil, isSandbox)
	if err != nil {
		return
	}

	reqApp = CreateAppParams(appId, mchApiKey, prepayId, mchId)
	return
}

func CreateAppParams(appId string, mchApiKey string, prepayId string, partnerId string) map[string]string {
	params := make(map[string]string, 6)
	params["appid"]     = appId
	params["partnerid"] = partnerId
	params["prepayid"]  = prepayId
	params["package"]   = "Sign=WXPay"
	params["noncestr"]  = string(_GetRandomBytes(32))
	params["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	params["sign"]      = createMd5Signature(params, mchApiKey)
	return params
}
