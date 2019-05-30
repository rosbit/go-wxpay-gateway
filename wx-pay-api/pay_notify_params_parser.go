// +build gateway notify

package wxpay

import (
	"fmt"
)

type PayCoupon struct {
	CouponType string `json:"coupon_type"`
	CouponId   string `json:"coupon_id"`
	CouponFee  int    `json:"coupon_fee"`
}

type IPayNotifyParams struct {
	AppId   string    `json:"app_id"`
	MchId   string    `json:"mch_id"`
	DeviceInfo string `json:"device_info"`
	ResultCode string `json:"result_code"`
	ErrCode    string `json:"err_code"`
	ErrCodeDes string `json:"err_code_des"`
	OpenId     string `json:"open_id"`
	IsSubscribe bool  `json:"is_subscribe"`
	TradeType string  `json:"trade_type"`
	BankType string   `json:"bank_type"`
	TotalFee int      `json:"total_fee"`
	SettlementTotalFee int `json:"settlement_total_fee"`
	FeeType string     `json:"fee_type"`
	CashFee int        `json:"cash_fee"`
	CashFeeType string `json:"cash_fee_type"`
	CouponFee   int    `json:"coupon_fee"`
	CouponCount int    `json:"coupon_count"`
	Coupons []PayCoupon  `json:"coupons"`
	TransactionId string `json:"transaction_id"`
	OrderId string     `json:"order_id"`
	Attach  string     `json:"attach"`
	TimeEnd string     `json:"time_end"`
}

type PayNotifyParams struct {
	AppName string    `json:"app_name"`
	CbUrl   string    `json:"cb_url"`
	IPayNotifyParams
}

func ParsePayNotifyBody(prompt string, body []byte, appKey string) (*PayNotifyParams, error) {
	_paymentLog.Printf("[pay-notify] 1. *** %s received: %s\n", prompt, string(body))

	res, err := parseXmlResult(body, appKey)
	if err != nil {
		_paymentLog.Printf("[pay-notify] 2. --- %s error: %v\n", prompt, err)
		return nil, err
	}
	_paymentLog.Printf("[pay-notify] 2. ### %s result: %v\n", prompt, res)

	m := _M(res)
	params := &PayNotifyParams{}

	if err = m.getString(&params.AppId, "appid", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.MchId, "mch_id", true); err != nil {
		return nil, err
	}
	m.getString(&params.DeviceInfo, "device_info", false)
	if err = m.getString(&params.ResultCode, "result_code", true); err != nil {
		return nil, err
	}
	m.getString(&params.ErrCode, "err_code", false)
	m.getString(&params.ErrCodeDes, "err_code_des", false)
	if err = m.getString(&params.OpenId, "openid", true); err != nil {
		return nil, err
	}
	if err = m.getBool(&params.IsSubscribe, "is_subscribe", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.TradeType, "trade_type", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.BankType, "bank_type", true); err != nil {
		return nil, err
	}
	if err = m.getInt(&params.TotalFee, "total_fee", true); err != nil {
		return nil, err
	}
	m.getInt(&params.SettlementTotalFee, "settlement_total_fee", false)
	m.getString(&params.FeeType, "fee_type", false)
	if err = m.getInt(&params.CashFee, "cash_fee", true); err != nil {
		return nil, err
	}
	m.getString(&params.CashFeeType, "cash_fee_type", false)
	m.getInt(&params.CouponFee, "coupon_fee", false)
	m.getInt(&params.CouponCount, "coupon_count", false)
	if params.CouponCount > 0 {
		params.Coupons = make([]PayCoupon, params.CouponCount)
		for i:=0; i<params.CouponCount; i++ {
			m.getString(&(params.Coupons[i].CouponType), fmt.Sprintf("coupon_type_%d", i), false)
			m.getString(&(params.Coupons[i].CouponId), fmt.Sprintf("coupon_id_%d", i), false)
			m.getInt(&(params.Coupons[i].CouponFee), fmt.Sprintf("coupon_fee_%d", i), false)
		}
	}
	if err = m.getString(&params.TransactionId, "transaction_id", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.OrderId, "out_trade_no", true); err != nil {
		return nil, err
	}
	m.getString(&params.Attach, "attach", false)
	if err = m.getString(&params.TimeEnd, "time_end", true); err != nil {
		return nil, err
	}
	return params, err
}
