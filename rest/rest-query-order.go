package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/wx-pay-api"
	"go-wxpay-gateway/conf"
	"net/http"
)

// POST /queryorder
// {
//      "appId": "appId of mp/mini-prog",
//      "payApp": "name-of-app-in-wxpay-gateway",
//      "transactionId": "transaction_id returned by wx",
//          - or -
//      "orderId": "unique-order-id",
//      "debug": false|true, default is false
// }
func QueryOrder(c *mgin.Context) {
	var queryParam struct {
		AppId         string
		PayApp        string
		TransactionId string
		OrderId       string
		Debug bool
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

	var queryFunc wxpay.FnQueryOrder
	var id string
	if queryParam.TransactionId != "" {
		queryFunc = wxpay.QueryByTransactionId
		id = queryParam.TransactionId
	} else if queryParam.OrderId != "" {
		queryFunc = wxpay.QueryByOrderId
		id = queryParam.OrderId
	} else {
		c.Error(http.StatusBadRequest, "Please specify transactionId or orderId")
		return
	}
	res, sent, recv, err := queryFunc(
		queryParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		id,
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

