package rest

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/wx-pay-api"
	"go-wxpay-gateway/conf"
	"net/http"
	"fmt"
	"encoding/json"
)

const (
	TRADE_TYPE_NAME = "trade_type"

	TYPE_WX     = "JSAPI"  // 微信内嵌浏览器或小程序支付
	TYPE_NATIVE = "NATIVE" // native支付
	TYPE_APP    = "APP"    // app支付
	TYPE_H5     = "H5"     // H5支付
)

var (
	createWxpays = map[string]func(*mgin.Context) {
		TYPE_WX:     createWxPay,
		TYPE_NATIVE: createNativePay,
		TYPE_APP:    createAppPay,
		TYPE_H5:     createH5Pay,
	}
)

// POST /create/:trade_type
func CreatePayment(c *mgin.Context) {
	tradeType := c.Param(TRADE_TYPE_NAME)
	if createPay, ok := createWxpays[tradeType]; ok {
		createPay(c)
		return
	}
	c.Error(http.StatusBadRequest, fmt.Sprintf("Unknown trade type \"%s\"", tradeType))
}

// POST /create/JSAPI
// POST BODY:
// {
//    "appId": "appId of mp/mini-prog",
//    "payApp": "name-of-app-in-wxpay-gateway",
//    "goods": "XXXX-xxxx",
//    "udd": "user defined data as parameters when calling callback",
//    "orderId": "unique order id",
//    "fee": xxx-in-fen,
//    "ip": "ip to create order",
//    "openId": "openId of mp service or mimi app",
//    "notifyUrl": "your notify url, which can be accessed outside",
//    "debug": false|true, default is false
// }
func createWxPay(c *mgin.Context) {
	var jsapiParam struct {
		AppId    string
		PayApp   string
		Goods    string
		Udd      string
		OrderId  string
		Fee      int
		Ip       string
		OpenId   string
		NotifyUrl string
		Debug bool
	}

	if code, err := c.ReadJSON(&jsapiParam); err != nil {
		c.Error(code, err.Error())
		return
	}

	isSandbox := _IsSandbox(jsapiParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(jsapiParam.PayApp)
	if !ok {
		c.Error(http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	prepayId, reqJSAPI, sent, recv, err := wxpay.JSAPIPay(
		jsapiParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		mchConf.Receipt,
		jsapiParam.Goods,
		jsapiParam.Udd,
		jsapiParam.OrderId,
		jsapiParam.Fee,
		jsapiParam.Ip,
		jsapiParam.NotifyUrl,
		jsapiParam.OpenId,
		isSandbox,
	)
	if err != nil {
		sendResultWithMsg(c, jsapiParam.Debug, sent, recv, err)
		return
	}

	// 为了和其它接口统一和方便保存, reqJSAPI转换成字符串再返回
	// 在使用时：先对result做JSON解析，得到jsapi_params后再做一次JSON解析
	b, _ := json.Marshal(reqJSAPI)
	sendResultWithMsg(c, jsapiParam.Debug, sent, recv, nil, map[string]interface{} {
		"result": map[string]interface{} {
			"prepay_id": prepayId,
			"jsapi_params": string(b), //reqJSAPI,
		},
	})
}

// POST /create/NATIVE
// POST Body:
// {
//    "appId": "appId of mp/mini-prog",
//    "payApp": "name-of-app-in-wxpay-gateway",
//    "goods": "XXXX-xxxx",
//    "udd": "user defined data as parameters when calling callback",
//    "orderId": "unique order id",
//    "fee": xxx-in-fen,
//    "ip": "ip to create order",
//    "productId": "productId",
//    "notifyUrl": "your notify url, which can be accessed outside",
//    "debug": false|true, default is false
// }
func createNativePay(c *mgin.Context) {
	var nativeParam struct {
		AppId    string
		PayApp   string
		Goods    string
		Udd      string
		OrderId  string
		Fee      int
		Ip       string
		ProductId string
		NotifyUrl string
		Debug bool
	}

	if code, err := c.ReadJSON(&nativeParam); err != nil {
		c.Error(code, err.Error())
		return
	}

	isSandbox := _IsSandbox(nativeParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(nativeParam.PayApp)
	if !ok {
		c.Error(http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	prepayId, codeUrl, sent, recv, err := wxpay.NativePay(
		nativeParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		mchConf.Receipt,
		nativeParam.Goods,
		nativeParam.Udd,
		nativeParam.OrderId,
		nativeParam.Fee,
		nativeParam.Ip,
		nativeParam.NotifyUrl,
		nativeParam.ProductId,
		isSandbox,
	)
	if err != nil {
		sendResultWithMsg(c, nativeParam.Debug, sent, recv, err)
		return
	}
	sendResultWithMsg(c, nativeParam.Debug, sent, recv, nil, map[string]interface{} {
		"result": map[string]interface{} {
			"prepay_id": prepayId,
			"code_url": codeUrl,
		},
	})
}

// POST /create/APP
// POST Body:
// {
//    "appId": "appId of mp/mini-prog",
//    "payApp": "name-of-app-in-wxpay-gateway",
//    "goods": "XXXX-xxxx",
//    "udd": "user defined data as parameters when calling callback",
//    "orderId": "unique order id",
//    "fee": xxx-in-fen,
//    "ip": "ip to create order",
//    "notifyUrl": "your notify url, which can be accessed outside",
//    "debug": false|true, default is false
// }
func createAppPay(c *mgin.Context) {
	var appParam struct {
		AppId    string
		PayApp   string
		Goods    string
		Udd      string
		OrderId  string
		Fee      int
		Ip       string
		NotifyUrl string
		Debug bool
	}

	if code, err := c.ReadJSON(&appParam); err != nil {
		c.Error(code, err.Error())
		return
	}

	isSandbox := _IsSandbox(appParam.PayApp)
	mchConf, ok := conf.GetAppAttrs(appParam.PayApp)
	if !ok {
		c.Error(http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	prepayId, reqAppPay, sent, recv, err := wxpay.AppPay(
		appParam.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		mchConf.Receipt,
		appParam.Goods,
		appParam.Udd,
		appParam.OrderId,
		appParam.Fee,
		appParam.Ip,
		appParam.NotifyUrl,
		isSandbox,
	)
	if err != nil {
		sendResultWithMsg(c, appParam.Debug, sent, recv, err)
		return
	}
	sendResultWithMsg(c, appParam.Debug, sent, recv, nil, map[string]interface{} {
		"result": map[string]interface{} {
			"prepay_id": prepayId,
			"req_params": reqAppPay,
		},
	})
}

// POST /create/H5
// POST Body:
// {
//     "appId": "appId of mp/mini-prog",
//     "payApp": "name-of-app-in-wxpay-gateway",
//     "goods": "XXXX-xxxx",
//     "udd": "user defined data as parameters when calling callback",
//     "orderId": "unique order id",
//     "fee": xxx-in-fen,
//     "ip": "ip to create order",
//     "redirectUrl": "url to redirect after payment",
//     "notifyUrl": "your notify url, which can be accessed outside",
//     "sceneInfo": {
//        "h5_info": {
//           "type": "",
//           "wap_name": "",
//           "wap_url": "wap-site-url"
//        }
//        ---- OR ----
//        "h5_info": {
//           "type": "",
//           "app_name":"",
//           "bundle_id": "ios-bundle-id"
//        }
//        ---- OR ----
//        "h5_info": {
//           "type": "",
//           "app_name":"",
//           "package_name": "android-package-name"
//        }
//     },
//     "debug": false|true, default is false
// }
func createH5Pay(c *mgin.Context) {
	var h5Param struct {
		AppId    string
		PayApp   string
		Goods    string
		Udd      string
		OrderId  string
		Fee      int
		Ip       string
		RedirectUrl string
		NotifyUrl   string
		SceneInfo interface{}
		Debug bool
	}

	if code, err := c.ReadJSON(&h5Param); err != nil {
		c.Error(code, err.Error())
		return
	}
	if h5Param.SceneInfo == nil {
		c.Error(http.StatusBadRequest, "sceneInfo not found")
		return
	}
	sceneInfo, err := json.Marshal(h5Param.SceneInfo)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	isSandbox := _IsSandbox(h5Param.PayApp)
	mchConf, ok := conf.GetAppAttrs(h5Param.PayApp)
	if !ok {
		c.Error(http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	prepayId, payUrl, sent, recv, err := wxpay.H5Pay(
		h5Param.AppId,
		mchConf.MchId,
		mchConf.MchApiKey,
		mchConf.Receipt,
		h5Param.Goods,
		h5Param.Udd,
		h5Param.OrderId,
		h5Param.Fee,
		h5Param.Ip,
		h5Param.NotifyUrl,
		h5Param.RedirectUrl,
		sceneInfo,
		isSandbox,
	)
	if err != nil {
		sendResultWithMsg(c, h5Param.Debug, sent, recv, err)
		return
	}
	sendResultWithMsg(c, h5Param.Debug, sent, recv, nil, map[string]interface{} {
		"result": map[string]interface{} {
			"prepay_id": prepayId,
			"pay_url": payUrl,
		},
	})
}
