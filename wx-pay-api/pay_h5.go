// +build gateway

package wxpay

import (
	"fmt"
	"net/url"
	"strings"
)

func H5Pay(
	appId  string,
	mchId  string,
	mchAppKey string,
	payBody   string,
	cbParams  string,
	orderId string,
	fee     int,
	ip      string,
	notifyUrl string,
	redirectUrl string,
	sceneInfo []byte,
	isSandbox bool,
) (prepay_id string, pay_url string, err error) {
	if isSandbox {
		/*
		if mchAppKey, err = GetSandbox(appId, mchId, mchAppKey); err != nil {
			return
		}*/
		fee = SANDBOX_FEE
	}

	var res map[string]string
	if prepay_id, res, err = payOrder(appId, mchId, mchAppKey, "H5", payBody, cbParams, orderId, fee, ip, notifyUrl, "MWEB", "", "", sceneInfo, isSandbox); err != nil {
		_paymentLog.Printf("[H5-payment] 3. --- %v\n", err)
		return "", "", err
	}

	mweb_url, ok := res["mweb_url"]
	if !ok {
		return "", "", fmt.Errorf("no mweb_url")
	}

	return prepay_id, _createPayUrl(mweb_url, redirectUrl), nil
}

func _createPayUrl(mweb_url, redirectUrl string) string {
	var b byte
	if strings.Index(mweb_url, "?") >= 0 {
		b = byte('&')
	} else {
		b = byte('?')
	}
	return fmt.Sprintf("%s%credirect_url=%s", mweb_url, b, url.QueryEscape(redirectUrl))
}
