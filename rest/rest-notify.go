package rest

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/rosbit/go-wxpay-gateway/conf"
	"github.com/rosbit/go-wxpay-gateway/wx-pay-api"
	"github.com/rosbit/go-wxpay-gateway/utils"
)

const (
	NOTIFY_APP_NAME = "app"
)

func notifyFromWx(w http.ResponseWriter, r *http.Request, fnParse wxpay.FnParseNotifyBody, prompt string, cbUrlIdx int) {
	returnCode, returnMsg := "SUCCESS", "OK"

	// 微信支付要求的返回结果格式，参见：https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_7&index=8
	// 微信支付退款通知要求的返回结果格式，参见：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_16&index=10
	defer func() {
		_WriteMessage(w, fmt.Sprintf(`<xml>
<return_code><![CDATA[%s]]></return_code>
<return_msg><![CDATA[%s]]></return_msg>
</xml>`, returnCode, returnMsg))
	}()

	appName := _PathParam(r, NOTIFY_APP_NAME)
	mchConf, cbUrls, ok := conf.GetAppAttrs(appName)
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnCode = "FAIL"
		returnMsg  = "failed to read body"
		log.Printf("[%s] %v\n", prompt, err)
		return
	}

	// parse xml params
	params := fnParse(prompt, body, mchConf.MchAppKey)
	params.AppName, params.CbUrl = appName, cbUrls[cbUrlIdx]
	utils.SaveResult(params)
}

// POST /notify-pay/:app
func NotifyPayment(w http.ResponseWriter, r *http.Request) {
	notifyFromWx(w, r, wxpay.ParsePayNotifyBody, "pay-notify", conf.NOTIFY_PAY_CB_IDX)
}

// POST /notify-refund/:app
func NotifyRefundment(w http.ResponseWriter, r *http.Request) {
	notifyFromWx(w, r, wxpay.ParseRefundNotifyBody, "refund-notify", conf.NOTIFY_REFUND_CB_IDX)
}
