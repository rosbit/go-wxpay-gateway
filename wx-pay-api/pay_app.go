// +build gateway

package wxpay

import (
	"fmt"
	"time"
)

func AppPay(
	appId  string,
	mchId  string,
	mchAppKey string,
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
		if mchAppKey, err = GetSandbox(appId, mchId, mchAppKey); err != nil {
			return
		}*/
		fee = SANDBOX_FEE
	}
	prepayId, _, err = payOrder(appId, mchId, mchAppKey, "APP", payBody, cbParams, orderId, fee, ip, notifyUrl, "APP", "", "", nil, isSandbox)
	if err != nil {
		_paymentLog.Printf("[App-payment] 3. --- %v\n", err)
		return "", nil, err
	}

	return prepayId, CreateAppParams(appId, mchAppKey, prepayId, mchId), nil
}

func CreateAppParams(appId string, mchAppKey string, prepayId string, partnerId string) map[string]string {
	params := make(map[string]string, 6)
	params["appid"]     = appId
	params["partnerid"] = partnerId
	params["prepayid"]  = prepayId
	params["package"]   = "Sign=WXPay"
	params["noncestr"]  = string(_GetRandomBytes(32))
	params["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	params["sign"]      = createMd5Signature(params, mchAppKey)
	_paymentLog.Printf("[App-payment] ### app payment params: %v\n", params)
	return params
}
