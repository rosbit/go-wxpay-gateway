package v3

import (
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"net/http"
	"fmt"
	"os"
	"context"
)

type TransferQueryItem struct {
	DetailId     string `json:"detail_id"`
	OutDetailNo  string `json:"out_detail_no"`
	DetailStatus string `json:"detail_status"`
}

const (
	queryPageSize = int64(100)

	ALL     = "ALL"
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
)

type FnQueryTransfer func(appName string, wxBatchId string, needDetail bool, status string) (total int64, batchStatus string, it <-chan *TransferQueryItem, err error)

type fnQueryTranferBatch func(ctx context.Context, a *transferbatch.TransferBatchApiService, pageNo int64) (resp *transferbatch.TransferBatchEntity, err error)

// 微信批次单号查询批次单API: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter4_3_2.shtml
func QueryTransferByWxBatchId(appName string, wxBatchId string, needDetail bool, status string) (total int64, batchStatus string, it <-chan *TransferQueryItem, err error) {
	getPage := func(ctx context.Context, a *transferbatch.TransferBatchApiService, pageNo int64) (resp *transferbatch.TransferBatchEntity, err error) {
		res, result, e := a.GetTransferBatchByNo(ctx, transferbatch.GetTransferBatchByNoRequest{
			BatchId: core.String(wxBatchId),
			NeedQueryDetail: core.Bool(needDetail),
			Offset: core.Int64(pageNo * queryPageSize),
			Limit: core.Int64(queryPageSize),
			DetailStatus: core.String(status),
		})
		if e != nil {
			err = e
			return
		}
		if result.Response.StatusCode != http.StatusOK {
			err = fmt.Errorf("%s", result.Response.Status)
			return
		}
		resp = res
		return
	}

	return queryTransfer(appName, getPage)
}

// 商家批次单号查询批次单API: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter4_3_5.shtml
func QueryTransferByBatchNo(appName string, batchNo string, needDetail bool, status string) (total int64, batchStatus string, it <-chan *TransferQueryItem, err error) {
	getPage := func(ctx context.Context, a *transferbatch.TransferBatchApiService, pageNo int64) (resp *transferbatch.TransferBatchEntity, err error) {
		res, result, e := a.GetTransferBatchByOutNo(ctx, transferbatch.GetTransferBatchByOutNoRequest{
			OutBatchNo: core.String(batchNo),
			NeedQueryDetail: core.Bool(needDetail),
			Offset: core.Int64(pageNo * queryPageSize),
			Limit: core.Int64(queryPageSize),
			DetailStatus: core.String(status),
		})
		if e != nil {
			err = e
			return
		}
		if result.Response.StatusCode != http.StatusOK {
			err = fmt.Errorf("%s", result.Response.Status)
			return
		}
		resp = res
		return
	}

	return queryTransfer(appName, getPage)
}

// 逐页查询明细
func queryTransfer(appName string, getPage fnQueryTranferBatch) (total int64, batchStatus string, it <-chan *TransferQueryItem, err error) {
	ctx, client, e := createClient(appName)
	if e != nil {
		err = e
		return
	}
	svc := &transferbatch.TransferBatchApiService{Client: client}

	pageNo := int64(0)
	res, e := getPage(ctx, svc, pageNo)
	if e != nil {
		err = e
		return
	}
	batchRes := res.TransferBatch
	batchStatus = *batchRes.BatchStatus
	if total = *batchRes.TotalNum; total <= 0 {
		return
	}

	resPage := make(chan []transferbatch.TransferDetailCompact)
	it = makeChanResult(resPage)

	go func() {
		count := int64(0)
		for {
			resPage <- res.TransferDetailList
			count += int64(len(res.TransferDetailList))
			if count >= total {
				break
			}
			pageNo += 1

			if res, e = getPage(ctx, svc, pageNo); e != nil {
				fmt.Fprintf(os.Stderr, "error occurs when calling getPage(pageNo:%d): %v\n", pageNo, e)
				break
			}
		}
		close(resPage)
	}()

	return
}

func makeChanResult(resPage <-chan []transferbatch.TransferDetailCompact) (<-chan *TransferQueryItem) {
	it := make(chan *TransferQueryItem)
	go func() {
		for res := range resPage {
			for i, _ := range res {
				r := &res[i]
				it <- &TransferQueryItem{
					DetailId: *r.DetailId,
					OutDetailNo: *r.OutDetailNo,
					DetailStatus: *r.DetailStatus,
				}
			}
		}

		close(it)
	}()

	return it
}
