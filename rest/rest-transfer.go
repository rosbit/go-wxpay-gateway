package rest

import (
	"net/http"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/wx-pay-api"
)

// POST /transfer
/*
 {
	 "appId": "appId of mp/mini-prog",
	 "payApp": "name-of-app-in-wxpay-gateway",
     "tradeNo": "unique-trade-no",
	 "openId": "openid to be transfered",
	 "userName": "real user name",
	 "amount": xxx-in-fen,
	 "desc": "description",
	 "ip": "ip to create order"
 }
*/
func Transfer(w http.ResponseWriter, r *http.Request) {
	var transferParam struct {
		AppId    string
		PayApp   string
		TradeNo  string
		OpenId   string
		UserName string
		Amount   int
		Desc     string
		Ip       string
	}

	if code, err := _ReadJson(r, &transferParam); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	isSandbox := _IsSandbox(transferParam.PayApp)
	mchConf, _, ok := conf.GetAppAttrs(transferParam.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	res, err := wxpay.Transfer(
		transferParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		transferParam.TradeNo,
		transferParam.OpenId,
		transferParam.UserName,
		transferParam.Amount,
		transferParam.Desc,
		transferParam.Ip,
		mchConf.MchCertPemFile,
		mchConf.MchKeyPemFile,
		isSandbox,
	)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_WriteJson(w, http.StatusOK, &res)
}

