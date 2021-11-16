package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/wx-pay-api"
	"go-wxpay-gateway/conf"
	"net/http"
)

// POST /transfer
// POST Body:
// {
//  "appId": "appId of mp/mini-prog",
//  "payApp": "name-of-app-in-wxpay-gateway",
//  "tradeNo": "unique-trade-no",
//  "openId": "openid to be transfered",
//  "userName": "real user name",
//  "amount": xxx-in-fen,
//  "desc": "description",
//  "ip": "ip to create order",
//  "debug": false|true, default is false
// }
func Transfer(c *mgin.Context) {
	var transferParam struct {
		AppId    string
		PayApp   string
		TradeNo  string
		OpenId   string
		UserName string
		Amount   int
		Desc     string
		Ip       string
		Debug    bool
	}

	if code, err := c.ReadJSON(&transferParam); err != nil {
		c.Error(code, err.Error())
		return
	}

	isSandbox := _IsSandbox(transferParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(transferParam.PayApp)
	if !ok {
		c.Error(http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	res, sent, recv, err := wxpay.Transfer(
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
		sendResultWithMsg(c, transferParam.Debug, sent, recv, err)
		return
	}
	sendResultWithMsg(c, transferParam.Debug, sent, recv, nil, map[string]interface{} {
		"result": res,
	})
}

