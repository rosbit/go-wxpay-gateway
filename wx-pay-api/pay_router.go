// +build gateway

package wxpay

import (
	"fmt"
)

var _payUrlRouters = map[bool]string{
	false: "",
	true: "sandboxnew/",
}

const (
	UT_UNIFIED_ORDER = iota
	UT_ORDER_QUERY
	UT_REFUND

	unifiedorder_url_fmt = "https://api.mch.weixin.qq.com/%spay/unifiedorder"
	orderquery_url_fmt   = "https://api.mch.weixin.qq.com/%spay/orderquery"
	refund_url_fmt       = "https://api.mch.weixin.qq.com/%spay/refund"

	SECURE_API_ROUTER    = "secapi/"
)

func _GetApiUrl(urlType int, isSandbox bool) string {
	urlRouter := _payUrlRouters[isSandbox]

	switch (urlType) {
	case UT_UNIFIED_ORDER:
		return fmt.Sprintf(unifiedorder_url_fmt, urlRouter)
	case UT_ORDER_QUERY:
		return fmt.Sprintf(orderquery_url_fmt, urlRouter)
	case UT_REFUND:
		if urlRouter == "" {
			urlRouter = SECURE_API_ROUTER
		}
		return fmt.Sprintf(refund_url_fmt, urlRouter)
	default:
		return ""
	}
}
