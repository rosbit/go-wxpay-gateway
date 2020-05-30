package rest

import (
	"net/http"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/wx-pay-api"
)

// POST /queryorder
// POST Body:
// {
//      "appId": "appId of mp/mini-prog",
//      "payApp": "name-of-app-in-wxpay-gateway",
//      "orderId": "unique-order-id",
//      "debug": false|true, default is false
// }
func CloseOrder(w http.ResponseWriter, r *http.Request) {
	var closeParam struct {
		AppId   string
		PayApp  string
		OrderId string
		Debug bool
	}
	if code, err := _ReadJson(r, &closeParam); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	isSandbox := _IsSandbox(closeParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(closeParam.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	if closeParam.OrderId == "" {
		_WriteError(w, http.StatusBadRequest, "Please specify orderId")
		return
	}

	sent, recv, err := wxpay.CloseOrder(
		closeParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		closeParam.OrderId,
		isSandbox,
	)
	sendResultWithMsg(closeParam.Debug, w, sent, recv, err)
}

