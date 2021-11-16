package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/wx-pay-api"
	"go-wxpay-gateway/conf"
	"net/http"
)

// POST /query-transfer
// POST Body
// {
//     "appId": "appId of mp/mini-prog",
//     "payApp": "name-of-app-in-wxpay-gateway",
//     "tradeNo": "unique-trade-no",
//     "debug": false|true, default is false
// }
func QueryTransfer(c *mgin.Context) {
	var queryParam struct {
		AppId   string
		PayApp  string
		TradeNo string
		Debug   bool
	}
	if code, err := c.ReadJSON(&queryParam); err != nil {
		c.Error(code, err.Error())
		return
	}

	isSandbox := _IsSandbox(queryParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(queryParam.PayApp)
	if !ok {
		c.Error(http.StatusBadRequest, "Unknown pay-app name")
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
		sendResultWithMsg(c, queryParam.Debug, sent, recv, err)
		return
	}
	sendResultWithMsg(c, queryParam.Debug, sent, recv, nil, map[string]interface{} {
		"result": res,
	})
}
