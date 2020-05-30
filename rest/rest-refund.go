package rest

import (
	"net/http"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/wx-pay-api"
)

// POST /refundment
// {
//      "appId": "appId of mp/mini-prog",
//      "payApp": "name-of-app-in-wxpay-gateway",
//      "transactionId": "transaction_id returned by wx",
//          - or -
//      "orderId": "unique-order-id",
//      "refundId": "unique-refund-id",
//      "totalFee": xxx-in-fen,
//      "refundFee": xxx-in-fen,
//      "refundReason": "reason",
//	    "notifyUrl": "your notify url, which can be accessed outside",
//      "debug": false|true, default is false
// }
func CreateRefundment(w http.ResponseWriter, r *http.Request) {
	var refundParam struct {
		AppId         string
		PayApp        string
		TransactionId string
		OrderId       string
		RefundId      string
		TotalFee      int
		RefundFee     int
		RefundReason  string
		NotifyUrl     string
		Debug         bool
	}

	if code, err := _ReadJson(r, &refundParam); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	isSandbox := _IsSandbox(refundParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(refundParam.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	var refundFn wxpay.FnRefund
	var id string
	if refundParam.TransactionId != "" {
		refundFn, id = wxpay.RefundByTransactionId, refundParam.TransactionId
	} else if refundParam.OrderId != "" {
		refundFn, id = wxpay.RefundByOrderId, refundParam.OrderId
	} else {
		_WriteError(w, http.StatusBadRequest, "Please specify transactionId or orderId")
		return
	}

	refundNotifyParams, sent, recv, err := refundFn(
		refundParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		id,
		refundParam.RefundId,
		refundParam.TotalFee,
		refundParam.RefundFee,
		refundParam.RefundReason,
		refundParam.NotifyUrl,
		mchConf.MchCertPemFile,
		mchConf.MchKeyPemFile,
		isSandbox,
	)
	if err != nil {
		sendResultWithMsg(refundParam.Debug, w, sent, recv, err)
		return
	}
	sendResultWithMsg(refundParam.Debug, w, sent, recv, nil, map[string]interface{} {
		"result": refundNotifyParams,
	})
}

