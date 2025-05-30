package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/v3-selfdev"
	"net/http"
)

// POST /v3/verify-transfer-bills/:payApp
// POST Body:
// {JSON-Body-From-WxPay}
func V3VerifyTransferBills(c *mgin.Context) {
	var params struct {
		PayApp   string `path:"payApp"`
	}

	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	respCode, respBody, billBody, err := v3sd.VerifyTransferBills(params.PayApp, c.Request())
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": http.StatusInternalServerError,
			"msg": err.Error(),
			"respToWxpay": map[string]interface{}{
				"respCode": respCode,
				"respBody": respBody,
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"respToWxpay": map[string]interface{}{
			"respCode": respCode,
			"respBody": respBody,
		},
		"result": map[string]interface{} {
			"bill": billBody,
		},
	})
}

