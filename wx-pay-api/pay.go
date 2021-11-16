// +build gateway

package wxpay

import (
	"go-wxpay-gateway/conf"
	"fmt"
	"time"
)

const (
	SANDBOX_FEE = 101
)

func postOrder(orderId string, apiKey string, xml []byte, isSandbox bool) (res map[string]string, recv []byte, err error) {
	unifiedorder_url := _GetApiUrl(UT_UNIFIED_ORDER, isSandbox)
	if recv, err = _CallWxAPI(unifiedorder_url, "POST", xml); err != nil {
		return
	}

	res, err = parseXmlResult(recv, apiKey)
	return
}

func payOrder(
	appId string,
	mchId string,
	mchApiKey string,
	receipt bool,
	deviceInfo string,
	payBody string,
	cbParams string,
	orderId string,
	fee int,
	ip string,
	notifyUrl string,
	tradeType string,
	productId string,
	openId string,
	sceneInfo []byte,
	isSandbox bool,
) (prepay_id string, res map[string]string, xmlstr, recv []byte, err error) {
	tags := make(map[string]string)
	xml := newXmlGenerator("xml")
	addTag(xml, tags, "appid",       appId,      false)
	addTag(xml, tags, "mch_id",      mchId,      false)
	addTag(xml, tags, "device_info", deviceInfo, true)
	addTag(xml, tags, "nonce_str",   string(_GetRandomBytes(32)), false)
	addTag(xml, tags, "body",        payBody, true)
	addTag(xml, tags, "attach",      cbParams,    true)
	addTag(xml, tags, "out_trade_no",orderId,     false)
	addTag(xml, tags, "total_fee",   fmt.Sprintf("%d", fee), false)
	addTag(xml, tags, "spbill_create_ip", ip,     false)
	addTag(xml, tags, "notify_url",  notifyUrl,   false)
	addTag(xml, tags, "trade_type",  tradeType,   false)
	addTag(xml, tags, "product_id",  productId,   tradeType != "NATIVE")
	addTag(xml, tags, "openid",      openId,      tradeType != "JSAPI")
	addTag(xml, tags, "time_expire", genExpire(), false)
	if sceneInfo != nil {
		addTag(xml, tags, "scene_info", string(sceneInfo), false)
	}
	if receipt {
		addTag(xml, tags, "receipt", "Y", false)
	}
	// sign
	signature := createMd5Signature(tags, mchApiKey)
	addTag(xml, tags, "sign", signature, false)

	xmlstr = xml.ToXML()
	// fmt.Printf("xml: %s\n", string(xmlstr))

	if res, recv, err = postOrder(orderId, mchApiKey, xmlstr, isSandbox); err != nil {
		return
	}

	// return_code is "SUCCESS", then check result_code
	if result_code, ok := res["result_code"]; !ok {
		err = fmt.Errorf("no result_code")
		return
	} else {
		if result_code != "SUCCESS" {
			err_code_des, ok := res["err_code_des"]
			if !ok {
				err = fmt.Errorf("no err_code_des")
			} else {
				err = fmt.Errorf("err_code_des: %s", err_code_des)
			}
			return
		}
	}

	// get prepay_id
	var ok bool
	if prepay_id, ok = res["prepay_id"]; !ok {
		err = fmt.Errorf("no prepay_id")
	}
	return
}

func genExpire() string {
	expire := time.Now().Add(time.Duration(conf.ServiceConf.OrderExpireMinutes)*time.Minute)
	return expire.Format("20060102150405")
}
