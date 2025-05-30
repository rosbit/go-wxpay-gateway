package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/v3-selfdev"
	"net/http"
)

// POST /v3/query-transfer-bills
// POST Body:
// {
//  "payApp": "name-of-app-in-wxpay-gateway",
//  "out_bills_no": "商户系统内部的商家单号，只能由数字、大小写字母组成，在商户系统内部唯一"
// }
func V3QueryTransferBills(c *mgin.Context) {
	var params struct {
		PayApp   string `json:"payApp"`
		OutBillNo string `json:"out_bills_no"`
	}

	if code, err := c.ReadJSON(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	resp, err := v3sd.QueryTransferBills(params.PayApp, params.OutBillNo)
	if err != nil {
		sendResultWithMsg(c, true, nil, nil, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"result": resp,
	})
}

