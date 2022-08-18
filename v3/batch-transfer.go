// 发起商家转账：https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter4_3_1.shtml

package v3

import (
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"net/http"
	"fmt"
)

type TransferDetail struct {
	SearialNo string `json:"tradeNo"`  // 转账明细单的唯一标识
	Amount    int64  `json:"amount"`   // 转账金额，单位分
	Remark    string `json:"desc"`     // 备注，最多32个字符
	OpenId    string `json:"openId"`   // appid下的openid
	UserName  string `json:"userName"` // 收款方姓名
}

func BatchTransfer(appName string, appId string, batchNo string, name string, remark string, details []TransferDetail) (wxBatchId string, err error) {
	fmt.Printf("appName: %s\n", appName)
	fmt.Printf("appId: %s\n", appId)
	fmt.Printf("batchNo: %s\n", batchNo)
	fmt.Printf("name: %s\n", name)
	fmt.Printf("remark: %s\n", remark)
	// fmt.Printf("details: %v\n", details)

	if len(details) == 0 {
		err = fmt.Errorf("details expected")
		return
	}

	var totalAmount int64
	totalNum := int64(len(details))
	transDetails := make([]transferbatch.TransferDetailInput, len(details))
	for i, _ := range details {
		d := &details[i]
		totalAmount += d.Amount
		transDetails[i] = transferbatch.TransferDetailInput{
			OutDetailNo: core.String(d.SearialNo),
			TransferAmount: core.Int64(d.Amount),
			TransferRemark: core.String(d.Remark),
			Openid: core.String(d.OpenId),
			UserName: func()*string{
				if len(d.UserName) > 0 {
					return core.String(d.UserName)
				}
				return nil
			}(),
		}
	}

	ctx, client, e := createClient(appName)
	if e != nil {
		err = e
		return
	}
	svc := transferbatch.TransferBatchApiService{Client: client}
	resp, result, e := svc.InitiateBatchTransfer(ctx, transferbatch.InitiateBatchTransferRequest {
		Appid: core.String(appId),
		OutBatchNo: core.String(batchNo),
		BatchName: core.String(name),
		BatchRemark: core.String(remark),
		TotalAmount: core.Int64(totalAmount),
		TotalNum: core.Int64(totalNum),
		TransferDetailList: transDetails,
	})
	if e != nil {
		err = e
		return
	}
	if result.Response.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s", result.Response.Status)
		return
	}
	wxBatchId = *resp.BatchId
	return
}
