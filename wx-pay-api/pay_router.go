// +build gateway

package wxpay

import (
	"fmt"
)

const (
	UT_UNIFIED_ORDER = iota
	UT_ORDER_QUERY
	UT_REFUND
	UT_ORDER_CLOSE

	UT_TRANSFER
	UT_TRANSFER_QUERY

	SECURE_API_ROUTER    = "secapi/"
)

var (
	_payUrlRouters = map[bool]string{
		false: "",
		true: "sandboxnew/",
	}

	_payUrls = map[int]string{
		// payment
		UT_UNIFIED_ORDER: "https://api.mch.weixin.qq.com/%spay/unifiedorder",
		UT_ORDER_QUERY:   "https://api.mch.weixin.qq.com/%spay/orderquery",
		UT_REFUND:        "https://api.mch.weixin.qq.com/%spay/refund",
		UT_ORDER_CLOSE:   "https://api.mch.weixin.qq.com/%spay/closeorder",
		// transfer
		UT_TRANSFER:      "https://api.mch.weixin.qq.com/%smmpaymkttransfers/promotion/transfers",
		UT_TRANSFER_QUERY:"https://api.mch.weixin.qq.com/%smmpaymkttransfers/gettransferinfo",
	}
)

func _GetApiUrl(urlType int, isSandbox bool) string {
	urlFmt, ok := _payUrls[urlType]
	if !ok {
		return ""
	}

	urlRouter := _payUrlRouters[isSandbox]
	if urlType == UT_REFUND {
		if urlRouter == "" {
			urlRouter = SECURE_API_ROUTER
		}
	}
	return fmt.Sprintf(urlFmt, urlRouter)
}
