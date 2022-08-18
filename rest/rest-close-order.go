package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/wx-pay-api"
	"go-wxpay-gateway/conf"
	"net/http"
)

// POST /closeorder
// POST Body:
// {
//      "appId": "appId of mp/mini-prog",
//      "payApp": "name-of-app-in-wxpay-gateway",
//      "orderId": "unique-order-id",
//      "debug": false|true, default is false
// }
func CloseOrder(c *mgin.Context) {
	var closeParam struct {
		AppId   string
		PayApp  string
		OrderId string
		Debug bool
	}
	if code, err := c.ReadJSON(&closeParam); err != nil {
		c.Error(code, err.Error())
		return
	}

	isSandbox := _IsSandbox(closeParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(closeParam.PayApp)
	if !ok {
		c.Error(http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	if closeParam.OrderId == "" {
		c.Error(http.StatusBadRequest, "Please specify orderId")
		return
	}

	res, sent, recv, err := wxpay.CloseOrder(
		closeParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		closeParam.OrderId,
		isSandbox,
	)
	if err != nil {
		sendResultWithMsg(c, closeParam.Debug, sent, recv, err)
		return
	}
	sendResultWithMsg(c, closeParam.Debug, sent, recv, nil, map[string]interface{}{
		"result": res,
	})
}

