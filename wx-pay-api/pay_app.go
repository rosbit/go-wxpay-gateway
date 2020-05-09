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
	payBody   string,
	cbParams  string,
	orderId string,
	fee     int,
	ip      string,
	notifyUrl string,
	isSandbox bool,
) (prepayId string, reqApp map[string]string, err error) {
	if isSandbox {
		/*
		if mchApiKey, err = GetSandbox(appId, mchId, mchApiKey); err != nil {
			return
		}*/
		fee = SANDBOX_FEE
	}
	prepayId, _, err = payOrder(appId, mchId, mchApiKey, "APP", payBody, cbParams, orderId, fee, ip, notifyUrl, "APP", "", "", nil, isSandbox)
	if err != nil {
		_paymentLog.Printf("[App-payment] 3. --- %v\n", err)
		return "", nil, err
	}

	return prepayId, CreateAppParams(appId, mchApiKey, prepayId, mchId), nil
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
	_paymentLog.Printf("[App-payment] ### app payment params: %v\n", params)
	return params
}
