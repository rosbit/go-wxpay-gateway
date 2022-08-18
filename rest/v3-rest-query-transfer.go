package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/v3"
	"strings"
	"fmt"
	"net/http"
	"encoding/json"
)

// POST /v3/query-transfer
// POST Body:
// {
//  "payApp": "name-of-app-in-wxpay-gateway",
//  "status": "" | "ALL" | "SUCCESS" | "FAIL"
//  "needDetail": true|false,
//  "wxBatchId": "wx-batchid"
//      ---- or ---
//  "batchNo": "unique-bach-no"
// }
func V3QueryTransfer(c *mgin.Context) {
	var params struct {
		PayApp     string
		Status     string
		NeedDetail bool
		WxBatchId  string
		BatchNo    string
	}

	if code, err := c.ReadJSON(&params); err != nil {
		c.Error(code, err.Error())
		return
	}
	if len(params.WxBatchId) == 0 && len(params.BatchNo) == 0 {
		c.Error(http.StatusBadRequest, "wxBatchId or batchNo expected")
		return
	}
	var status string = v3.ALL
	switch strings.ToUpper(params.Status) {
	case v3.SUCCESS:
		status = v3.SUCCESS
	case v3.FAIL:
		status = v3.FAIL
	default:
		status = v3.ALL
	}

	var queryTransfer v3.FnQueryTransfer
	var batchNo string
	if len(params.WxBatchId) > 0 {
		queryTransfer, batchNo = v3.QueryTransferByWxBatchId, params.WxBatchId
	} else {
		queryTransfer, batchNo = v3.QueryTransferByBatchNo, params.BatchNo
	}

	total, batchStatus, it, err := queryTransfer(params.PayApp, batchNo, params.NeedDetail, status)
	if err != nil {
		sendResultWithMsg(c, false, nil, nil, err)
		return
	}

	out := c.Response()
	c.SetHeader("Content-Type", "application/json")
	fmt.Fprintf(out, `{"code":%d,"msg":"OK","result":{"status":"%s","total":%d`, http.StatusOK, batchStatus, total)
	if total == 0 {
		fmt.Fprintf(out, "}}")
		return
	}
	jOut := json.NewEncoder(out)
	jOut.SetEscapeHTML(false)
	i := 0
	fmt.Fprintf(out, `,"details":`)
	for res := range it {
		if i == 0 {
			fmt.Fprintf(out, "[")
		} else {
			fmt.Fprintf(out, ",")
		}
		i += 1
		jOut.Encode(res)
	}
	if i == 0 {
		fmt.Fprintf(out, "null}}")
	} else {
		fmt.Fprintf(out, "]}}")
	}
}

