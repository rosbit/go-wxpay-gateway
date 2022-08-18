package v3

import (
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"net/http"
	"fmt"
	"context"
)

type FnQueryTransferDetail func(appName string, wxBatchId string, detailId string) (detail *transferbatch.TransferDetailEntity, err error)

type fnQueryTranferDetail func(ctx context.Context, a *transferbatch.TransferDetailApiService) (detail *transferbatch.TransferDetailEntity, err error)

// 微信明细单号查询明细单API: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter4_3_3.shtml
func QueryTransferDetailByWxBatchId(appName string, wxBatchId string, detailId string) (detail *transferbatch.TransferDetailEntity, err error) {
	getDetail := func(ctx context.Context, a *transferbatch.TransferDetailApiService) (detail *transferbatch.TransferDetailEntity, err error) {
		res, result, e := a.GetTransferDetailByNo(ctx, transferbatch.GetTransferDetailByNoRequest{
			BatchId: core.String(wxBatchId),
			DetailId: core.String(detailId),
		})
		if e != nil {
			err = e
			return
		}
		if result.Response.StatusCode != http.StatusOK {
			err = fmt.Errorf("%s", result.Response.Status)
			return
		}
		detail = res
		return
	}

	return queryTransferDetail(appName, getDetail)
}

// 商家明细单号查询明细单API: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter4_3_6.shtml
func QueryTransferDetailByBatchNo(appName string, batchNo string, tradeNo string) (detail *transferbatch.TransferDetailEntity, err error) {
	getDetail := func(ctx context.Context, a *transferbatch.TransferDetailApiService) (detail *transferbatch.TransferDetailEntity, err error) {
		res, result, e := a.GetTransferDetailByOutNo(ctx, transferbatch.GetTransferDetailByOutNoRequest{
			OutBatchNo: core.String(batchNo),
			OutDetailNo: core.String(tradeNo),
		})
		if e != nil {
			err = e
			return
		}
		if result.Response.StatusCode != http.StatusOK {
			err = fmt.Errorf("%s", result.Response.Status)
			return
		}
		detail = res
		return
	}

	return queryTransferDetail(appName, getDetail)
}

func queryTransferDetail(appName string, getDetail fnQueryTranferDetail) (detail *transferbatch.TransferDetailEntity, err error) {
	ctx, client, e := createClient(appName)
	if e != nil {
		err = e
		return
	}
	svc := &transferbatch.TransferDetailApiService{Client: client}
	return getDetail(ctx, svc)
}
