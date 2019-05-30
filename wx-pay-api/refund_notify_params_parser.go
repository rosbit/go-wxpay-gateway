// +build gateway notify

package wxpay

import (
	"fmt"
	"crypto/md5"
	"encoding/base64"
	"github.com/rosbit/go-aes"
	"io"
)

type IRefundNotifyParams struct {
	AppId   string    `json:"app_id"`
	MchId   string    `json:"mch_id"`
	TransactionId string `json:"transaction_id"`
	OrderId string     `json:"order_id"`
	WxRefundId string  `json:"wx_refund_id"`
	RefundId string    `json:"refund_id"`
	TotalFee int      `json:"total_fee"`
	SettlementTotalFee int `json:"settlement_total_fee"`
	RefundFee int      `json:"refund_fee"`
	SettlementRefundFee int `json:"settlement_refund_fee"`
	RefundStatus string     `json:"refund_status"`
	SuccessTime string     `json:"success_time"`
	RefundRecvAccout string `json:"refund_recv_accout"`
	RefundAccount string  `json:"refund_account"`
	RefundRequestSource string `json:"refund_request_source"`
}

type RefundNotifyParams struct {
	AppName string    `json:"app_name"`
	CbUrl   string    `json:"cb_url"`
	IRefundNotifyParams
}

func _DecryptRefundNotify(reqInfo string, appKey string) ([]byte, error) {
	oriReq, err := base64.StdEncoding.DecodeString(reqInfo)
	if err != nil {
		return nil, err
	}
	h := md5.New()
	io.WriteString(h, appKey)
	key := []byte(fmt.Sprintf("%x", h.Sum(nil))) // [32]byte
	return goaes.AesDecrypt(oriReq, key)
}

func ParseRefundNotifyBody(prompt string, body []byte, appKey string) (*RefundNotifyParams, error) {
	_paymentLog.Printf("[refund-notify] 1. *** %s received: %s\n", prompt, string(body))

	res, err := xml2map(body)
	if err != nil {
		_paymentLog.Printf("[refund-notify] 2. --- %s error: %v\n", prompt, err)
		return nil, err
	}
	_paymentLog.Printf("[refund-notify] 2. ### %s result: %v\n", prompt, res)
	req_info, ok := res["req_info"]
	if !ok {
		_paymentLog.Printf("[refund-notify] 3. ### %s no req_info found\n", prompt)
		return nil, fmt.Errorf("no req_info found in notify result")
	}
	reqInfoXml, err := _DecryptRefundNotify(req_info, appKey)
	if err != nil {
		return nil, err
	}
	_paymentLog.Printf("[refund-notify] 3. ### %s decrypted req_info: %s\n", prompt, string(reqInfoXml))
	reqInfo, err := xml2mapWithRoot(reqInfoXml, "root")
	if err != nil {
		return nil, err
	}

	m := _M(res)
	params := &RefundNotifyParams{}

	if err = m.getString(&params.AppId, "appid", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.MchId, "mch_id", true); err != nil {
		return nil, err
	}

	m = _M(reqInfo)
	if err = m.getString(&params.TransactionId, "transaction_id", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.OrderId, "out_trade_no", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.RefundId, "out_refund_no", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.WxRefundId, "refund_id", true); err != nil {
		return nil, err
	}
	if err = m.getInt(&params.TotalFee, "total_fee", true); err != nil {
		return nil, err
	}
	m.getInt(&params.SettlementTotalFee, "settlement_total_fee", false)
	if err = m.getInt(&params.RefundFee, "refund_fee", true); err != nil {
		return nil, err
	}
	m.getInt(&params.SettlementRefundFee, "settlement_refund_fee", false)
	if err = m.getString(&params.RefundStatus, "refund_status", true); err != nil {
		return nil, err
	}
	m.getString(&params.SuccessTime, "success_time", false)
	if err = m.getString(&params.RefundRecvAccout, "refund_recv_accout", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.RefundAccount, "refund_account", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.RefundRequestSource, "refund_request_source", true); err != nil {
		return nil, err
	}

	return params, nil
}
