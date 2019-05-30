// +build gateway

package wxpay

import (
	"fmt"
)

type RefundCoupon struct {
	CouponType      string `json:"coupon_type"`
	CouponRefundId  string `json:"coupon_refund_id"`
	CouponRefundFee int    `json:"coupon_refund_fee"`
}

type RefundResultParams struct {
	AppId   string    `json:"app_id"`
	MchId   string    `json:"mch_id"`
	ResultCode string `json:"result_code"`
	ErrCode    string `json:"err_code"`
	ErrCodeDes string `json:"err_code_des"`
	TransactionId string `json:"transaction_id"`
	OrderId string     `json:"order_id"`
	WxRefundId string  `json:"wx_refund_id"`
	RefundId string    `json:"refund_id"`
	TotalFee int      `json:"total_fee"`
	SettlementTotalFee int `json:"settlement_total_fee"`
	RefundFee int      `json:"refund_fee"`
	SettlementRefundFee int `json:"settlement_refund_fee"`
	FeeType string     `json:"fee_type"`
	CashFee int        `json:"cash_fee"`
	CashFeeType string `json:"cash_fee_type"`
	CashRefundFee int  `json:"cash_refund_fee"`
	CouponRefundFee   int   `json:"coupon_refund_fee"`
	CouponRefundCount int   `json:"coupon_refund_count"`
	RefundCoupons []RefundCoupon  `json:"refund_coupons"`
}

func ParseRefundResultBody(prompt string, body []byte, appKey string) (*RefundResultParams, error) {
	// _paymentLog.Printf("[refund-result] 1. *** %s received: %s\n", prompt, string(body))

	res, err := parseXmlResult(body, appKey)
	if err != nil {
		_paymentLog.Printf("[refund-result] --- %s error: %v\n", prompt, err)
		return nil, err
	}
	_paymentLog.Printf("[refund-result] ### %s result: %v\n", prompt, res)

	m := _M(res)
	params := &RefundResultParams{}

	if err = m.getString(&params.AppId, "appid", true); err != nil {
		return nil, err
	}
	if err = m.getString(&params.MchId, "mch_id", true); err != nil {
		return nil, err
	}
	m.getString(&params.ErrCode, "err_code", false)
	m.getString(&params.ErrCodeDes, "err_code_des", false)
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
	m.getString(&params.FeeType, "fee_type", false)
	if err = m.getInt(&params.CashFee, "cash_fee", true); err != nil {
		return nil, err
	}
	m.getString(&params.CashFeeType, "cash_fee_type", false)
	m.getInt(&params.CashRefundFee, "cash_refund_fee", false)
	m.getInt(&params.CouponRefundFee, "coupon_refund_fee", false)
	m.getInt(&params.CouponRefundCount, "coupon_refund_count", false)
	if params.CouponRefundCount > 0 {
		params.RefundCoupons = make([]RefundCoupon, params.CouponRefundCount)
		for i:=0; i<params.CouponRefundCount; i++ {
			m.getString(&(params.RefundCoupons[i].CouponType), fmt.Sprintf("coupon_type_%d", i), false)
			m.getString(&(params.RefundCoupons[i].CouponRefundId), fmt.Sprintf("coupon_refund_id_%d", i), false)
			m.getInt(&(params.RefundCoupons[i].CouponRefundFee), fmt.Sprintf("coupon_refund_fee_%d", i), false)
		}
	}

	return params, nil
}
