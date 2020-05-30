package rest

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/wx-pay-api"
)

const (
	NOTIFY_APP_NAME = "app"
)

func verifyNotify(w http.ResponseWriter, r *http.Request, fnParse wxpay.FnParseNotifyBody, prompt string) {
	returnCode, returnMsg := "SUCCESS", "OK"
	var params wxpay.INotifyParams
	var err error

	// 微信支付要求的返回结果格式，参见：https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_7&index=8
	// 微信支付退款通知要求的返回结果格式，参见：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_16&index=10
	defer func() {
		msgForWxpay := fmt.Sprintf(`<xml>
<return_code><![CDATA[%s]]></return_code>
<return_msg><![CDATA[%s]]></return_msg>
</xml>`, returnCode, returnMsg)

		code := http.StatusOK
		errMsg := ""
		if err != nil {
			code = http.StatusNotAcceptable
			errMsg = err.Error()
		}
		_WriteJson(w, code, map[string]interface{}{
			"code": code,
			"msg": errMsg,
			"params": params,
			"msgForWxpay": msgForWxpay,
		})
	}()

	q := r.URL.Query()
	appName := q.Get(NOTIFY_APP_NAME)
	mchConf, ok := conf.GetAppAttrs(appName)
	if !ok {
		returnCode = "FAIL"
		returnMsg  = "bad request"
		log.Printf("[%s] unknow app-name %s in uri: %s\n", prompt, appName, r.RequestURI)
		return
	}

	// read post body
	if r.Body == nil {
		returnCode = "FAIL"
		returnMsg  = "bad request"
		log.Printf("[%s] no body given\n", prompt)
		return
	}
	defer r.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		returnCode = "FAIL"
		returnMsg  = "failed to read body"
		log.Printf("[%s] %v\n", prompt, err)
		return
	}

	// parse xml params
	params, err = fnParse(prompt, body, mchConf.MchApiKey)
}

// POST /verify-notify-pay?app=<appName>
func VerifyNotifyPayment(w http.ResponseWriter, r *http.Request) {
	verifyNotify(w, r, wxpay.ParsePayNotifyBody, "pay-notify")
}

// POST /verify-notify-refund?app=<appName>
func VerifyNotifyRefundment(w http.ResponseWriter, r *http.Request) {
	verifyNotify(w, r, wxpay.ParseRefundNotifyBody, "refund-notify")
}
