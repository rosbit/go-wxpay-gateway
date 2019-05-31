package rest

import (
	"net/http"
	"github.com/rosbit/go-wxpay-gateway/conf"
	"github.com/rosbit/go-wxpay-gateway/wx-pay-api"
)

// POST /queryorder
/* {
      "appId": "appId of mp/mini-prog",
      "payApp": "name-of-app-in-wxpay-gateway",
      "transactionId": "transaction_id returned by wx",
          - or -
      "orderId": "unique-order-id"
 * }
 */
func QueryOrder(w http.ResponseWriter, r *http.Request) {
	var queryParam struct {
		AppId         string
		PayApp        string
		TransactionId string
		OrderId       string
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
	res, err := queryFunc(
		queryParam.AppId,
		mchConf.MchId,
		mchConf.MchAppKey,
		id,
		isSandbox,
	)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_WriteJson(w, http.StatusOK, &res.INotifyParams)
}

