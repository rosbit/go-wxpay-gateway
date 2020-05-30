package rest

import (
	"net/http"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/wx-pay-api"
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
func QueryOrder(w http.ResponseWriter, r *http.Request) {
	var queryParam struct {
		AppId         string
		PayApp        string
		TransactionId string
		OrderId       string
		Debug bool
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

	var queryFunc wxpay.FnQueryOrder
	var id string
	if queryParam.TransactionId != "" {
		queryFunc = wxpay.QueryByTransactionId
		id = queryParam.TransactionId
	} else if queryParam.OrderId != "" {
		queryFunc = wxpay.QueryByOrderId
		id = queryParam.OrderId
	} else {
		_WriteError(w, http.StatusBadRequest, "Please specify transactionId or orderId")
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
		sendResultWithMsg(queryParam.Debug, w, sent, recv, err)
		return
	}
	sendResultWithMsg(queryParam.Debug, w, sent, recv, nil, map[string]interface{} {
		"result": res,
	})
}

