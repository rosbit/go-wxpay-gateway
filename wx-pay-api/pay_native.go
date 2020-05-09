// +build gateway

package wxpay

import (
	"fmt"
)

func NativePay(
	appId  string,
	mchId  string,
	mchApiKey string,
	payBody   string,
	cbParams  string,
	orderId string,
	fee     int,
	ip      string,
	notifyUrl string,
	productId string,
	isSandbox bool,
) (prepay_id string, code_url string, err error) {
	if isSandbox {
		/*
		if mchApiKey, err = GetSandbox(appId, mchId, mchApiKey); err != nil {
			return
		}*/
		fee = SANDBOX_FEE
	}
	var res map[string]string
	prepay_id, res, err = payOrder(appId, mchId, mchApiKey, "WEB", payBody, cbParams, orderId, fee, ip, notifyUrl, "NATIVE", productId, "", nil, isSandbox)
	if err != nil {
		_paymentLog.Printf("[NATIVE-payment] 3. --- %v\n", err)
		return "", "", err
	}
	var ok bool
	if code_url, ok = res["code_url"]; !ok {
		return "", "", fmt.Errorf("code_url not found")
	}
	return
}
