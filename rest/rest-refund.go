package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/wx-pay-api"
	"go-wxpay-gateway/conf"
	"net/http"
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
func CreateRefundment(c *mgin.Context) {
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

	if code, err := c.ReadJSON(&refundParam); err != nil {
		c.Error(code, err.Error())
		return
	}

	isSandbox := _IsSandbox(refundParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(refundParam.PayApp)
	if !ok {
		c.Error(http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	var refundFn wxpay.FnRefund
	var id string
	if refundParam.TransactionId != "" {
		refundFn, id = wxpay.RefundByTransactionId, refundParam.TransactionId
	} else if refundParam.OrderId != "" {
		refundFn, id = wxpay.RefundByOrderId, refundParam.OrderId
	} else {
		c.Error(http.StatusBadRequest, "Please specify transactionId or orderId")
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
		sendResultWithMsg(c, refundParam.Debug, sent, recv, err)
		return
	}
	sendResultWithMsg(c, refundParam.Debug, sent, recv, nil, map[string]interface{} {
		"result": refundNotifyParams,
	})
}

