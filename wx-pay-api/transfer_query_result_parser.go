// +build gateway

package wxpay

import (
	"strconv"
)

type QueryTransferResult struct {
	AppId       string `json:"appid"`
	MchId       string `json:"mch_id"`
	ErrCode     string `json:"err_code"`
	ErrCodeDes  string `json:"err_code_des"`
	TradeNo     string `json:"partner_trade_no"`
	DetailId    string `json:"detail_id"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	OpenId      string `json:"openid"`
	TransferName string `json:"transfer_name"`
	Amount      int    `json:"payment_amount"`
	TransferTime string `json:"transfer_time"`
	PaymentTime string `json:"payment_time"`
	Desc        string
}

func ParseQueryTransferResult(prompt string, body []byte, apiKey string) (*QueryTransferResult, error) {
	res, err := parseXmlResult(body, apiKey)
	if err != nil {
		_paymentLog.Printf("[query-transfer-result] --- %s error: %v\n", prompt, err)
		return nil, err
	}
	_paymentLog.Printf("[query-transfer-result] ### %s result: %v\n", prompt, res)

	m := _M(res)
	params := &QueryTransferResult{}

	if err = m.getString(&params.AppId, "appid", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.MchId, "mch_id", true); err != nil {
		return nil, err
	}
	m.getString(&params.ErrCode, "err_code", false)
	m.getString(&params.ErrCodeDes, "err_code_des", false)
	if err = m.getString(&params.TradeNo, "partner_trade_no", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.DetailId, "detail_id", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.Status, "status", true); err != nil {
		return nil, err
	}
	m.getString(&params.Reason, "reason", false)
	if err = m.getString(&params.OpenId, "openid", true); err != nil {
		return nil, err
	}
	m.getString(&params.TransferName, "transfer_name", false)
	var sAmount string
	if err = m.getString(&sAmount, "payment_amount", true); err != nil {
		return nil, err
	}
	if params.Amount, err = strconv.Atoi(sAmount); err != nil {
		return nil, err
	}
	if err = m.getString(&params.TransferTime, "transfer_time", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.PaymentTime, "payment_time", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.Desc, "desc", true); err != nil {
		return nil, err
	}

	return params, nil
}
