// +build gateway

package wxpay

type TransferResult struct {
	AppId   string    `json:"mch_appid"`
	MchId   string    `json:"mch_id"`
	DeviceInfo string `json:"device_info"`
	ErrCode    string `json:"err_code"`
	ErrCodeDes string `json:"err_code_des"`
	TradeNo    string `json:"partner_trade_no"`
	PaymentNo  string `json:"payment_no"`
	PaymentTime string `json:"payment_time"`
}

func ParseTransferResultBody(prompt string, body []byte, apiKey string) (*TransferResult, error) {
	res, err := parseXmlResult(body, apiKey)
	if err != nil {
		_paymentLog.Printf("[transfer-result] --- %s error: %v\n", prompt, err)
		return nil, err
	}
	_paymentLog.Printf("[transfer-result] ### %s result: %v\n", prompt, res)

	m := _M(res)
	params := &TransferResult{}

	if err = m.getString(&params.AppId, "mch_appid", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.MchId, "mchid", true); err != nil {
		return nil, err
	}
	m.getString(&params.DeviceInfo, "device_info", false)
	m.getString(&params.ErrCode, "err_code", false)
	m.getString(&params.ErrCodeDes, "err_code_des", false)
	if err = m.getString(&params.TradeNo, "partner_trade_no", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.PaymentNo, "payment_no", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.PaymentTime, "payment_time", true); err != nil {
		return nil, err
	}

	return params, nil
}
