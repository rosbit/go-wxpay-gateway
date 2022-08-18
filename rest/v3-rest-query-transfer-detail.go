package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/v3"
	"fmt"
	"net/http"
)

// POST /v3/query-transfer-detail
// POST Body:
// {
//  "payApp": "name-of-app-in-wxpay-gateway",
//
//  "wxBatchId": "wx-batchid",
//  "wxDetailId": "wx-detailid"
//      ---- or ---
//  "batchNo": "unique-bach-no",
//  "tradeNo": "unique-trade-no"
// }
func V3QueryTransferDetail(c *mgin.Context) {
	var params struct {
		PayApp     string

		WxBatchId  string
		WxDetailId string
		// ---- or ----
		BatchNo    string
		TradeNo    string
	}

	if code, err := c.ReadJSON(&params); err != nil {
		c.Error(code, err.Error())
		return
	}
	if len(params.WxBatchId) == 0 && len(params.BatchNo) == 0 {
		c.Error(http.StatusBadRequest, "wxBatchId or batchNo expected")
		return
	}

	var queryTransferDetail v3.FnQueryTransferDetail
	var batchNo, tradeNo string
	if len(params.WxBatchId) > 0 {
		if len(params.WxDetailId) == 0 {
			c.Error(http.StatusBadRequest, "wxDetailId expected")
			return
		}
		queryTransferDetail, batchNo, tradeNo = v3.QueryTransferDetailByWxBatchId, params.WxBatchId, params.WxDetailId
	} else {
		if len(params.TradeNo) == 0 {
			c.Error(http.StatusBadRequest, "tradeNo expected")
			return
		}
		queryTransferDetail, batchNo, tradeNo = v3.QueryTransferDetailByBatchNo, params.BatchNo, params.TradeNo
	}

	res, err := queryTransferDetail(params.PayApp, batchNo, tradeNo)
	if err != nil {
		sendResultWithMsg(c, false, nil, nil, err)
		return
	}
	fmt.Printf("res: %v\n", res)

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]interface{}{
			"status": *res.DetailStatus,
			"amount": *res.TransferAmount,
			"remark": *res.TransferRemark,
			"reason": func()string{
				if res.FailReason == nil {
					return ""
				}
				return string(*res.FailReason)
			}(),
		},
	})
}

