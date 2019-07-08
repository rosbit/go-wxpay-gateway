package rest

import (
	"net/http"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/wx-pay-api"
)

// POST /query-transfer
/*
 {
      "appId": "appId of mp/mini-prog",
      "payApp": "name-of-app-in-wxpay-gateway",
      "tradeNo": "unique-trade-no",
 }
*/
func QueryTransfer(w http.ResponseWriter, r *http.Request) {
	var queryParam struct {
		AppId   string
		PayApp  string
		TradeNo string
	}
	if code, err := _ReadJson(r, &queryParam); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	isSandbox := _IsSandbox(queryParam.PayApp)
	mchConf, _, ok := conf.GetAppAttrs(queryParam.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}
	res, err := wxpay.QueryTransfer(
		queryParam.AppId,
		mchConf.MchId,
		mchConf.MchAppKey,
		queryParam.TradeNo,
		mchConf.MchCertPemFile,
		mchConf.MchKeyPemFile,
		isSandbox,
	)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_WriteJson(w, http.StatusOK, res)
}
