package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/v3"
	"net/http"
)

// POST /v3/transfer
// POST Body:
// {
//  "payApp": "name-of-app-in-wxpay-gateway",
//  "appId": "appId of mp/mini-prog",
//  "batchNo": "unique-bach-no",
//  "batchName": "batch name",
//  "batchRemark": "batch remark",
//  "details": [
//    {
//     "tradeNo": "unique-trade-no",
//     "amount": xxx-in-fen,
//     "desc": "description",
//     "openId": "openid to be transfered",
//     "userName": "real user name"
//   },
//   {...}
//  ]
// }
func V3Transfer(c *mgin.Context) {
	var params struct {
		PayApp   string
		AppId    string
		BatchNo  string
		BatchName string
		BatchRemark string
		Details []v3.TransferDetail
	}

	if code, err := c.ReadJSON(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	wxBatchId, err := v3.BatchTransfer(params.PayApp, params.AppId, params.BatchNo, params.BatchName, params.BatchRemark, params.Details)
	if err != nil {
		sendResultWithMsg(c, false, nil, nil, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]interface{}{
			"wxBatchId": wxBatchId,
		},
	})
}

