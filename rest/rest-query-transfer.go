package rest

import (
	"net/http"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/wx-pay-api"
)

// POST /query-transfer
// POST Body
// {
//     "appId": "appId of mp/mini-prog",
//     "payApp": "name-of-app-in-wxpay-gateway",
//     "tradeNo": "unique-trade-no",
//     "debug": false|true, default is false
// }
func QueryTransfer(w http.ResponseWriter, r *http.Request) {
	var queryParam struct {
		AppId   string
		PayApp  string
		TradeNo string
		Debug   bool
	}
	if code, err := _ReadJson(r, &queryParam); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	isSandbox := _IsSandbox(queryParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(queryParam.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}
	res, sent, recv, err := wxpay.QueryTransfer(
		queryParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		queryParam.TradeNo,
		mchConf.MchCertPemFile,
		mchConf.MchKeyPemFile,
		isSandbox,
	)
	if err != nil {
		sendResultWithMsg(queryParam.Debug, w, sent, recv, err)
		return
	}
	sendResultWithMsg(queryParam.Debug, w, sent, recv, nil, map[string]interface{} {
		"result": res,
	})
}
