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
	mchApiKey string,
	receipt bool,
	payBody   string,
	cbParams  string,
	orderId string,
	fee     int,
	ip      string,
	notifyUrl string,
	redirectUrl string,
	sceneInfo []byte,
	isSandbox bool,
) (prepay_id string, pay_url string, sent, recv []byte, err error) {
	if isSandbox {
		/*
		if mchApiKey, err = GetSandbox(appId, mchId, mchApiKey); err != nil {
			return
		}*/
		fee = SANDBOX_FEE
	}

	var res map[string]string
	if prepay_id, res, sent, recv, err = payOrder(appId, mchId, mchApiKey, receipt, "H5", payBody, cbParams, orderId, fee, ip, notifyUrl, "MWEB", "", "", sceneInfo, isSandbox); err != nil {
		return
	}

	mweb_url, ok := res["mweb_url"]
	if !ok {
		err = fmt.Errorf("no mweb_url")
		return
	}

	pay_url = _createPayUrl(mweb_url, redirectUrl)
	return
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
