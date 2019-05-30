package rest

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/rosbit/go-wxpay-gateway/conf"
	"github.com/rosbit/go-wxpay-gateway/wx-pay-api"
)

const (
	TRADE_TYPE_NAME = "trade_type"

	TYPE_WX     = "JSAPI"  // 微信内嵌浏览器或小程序支付
	TYPE_NATIVE = "NATIVE" // native支付
	TYPE_APP    = "APP"    // app支付
	TYPE_H5     = "H5"     // H5支付
)

// POST /create/:trade_type
func CreatePayment(w http.ResponseWriter, r *http.Request) {
	tradeType := _PathParam(r, TRADE_TYPE_NAME)
	switch tradeType {
	case TYPE_WX:
		createWxPay(w, r)
	case TYPE_NATIVE:
		createNativePay(w, r)
	case TYPE_APP:
		createAppPay(w, r)
	case TYPE_H5:
		createH5Pay(w, r)
	default:
		_WriteError(w, http.StatusBadRequest, fmt.Sprintf("Unknown trade type \"%s\"", tradeType))
	}
}

/**
 * JSAPI pay
 * {
      "appId": "appId of mp/mini-prog",
      "payApp": "name-of-app-in-wxpay-gateway",
      "goods": "XXXX-xxxx",
      "udd": "user defined data as parameters when calling callback",
      "orderId": "unique order id",
      "fee": xxx-in-fen,
      "ip": "ip to create order",
      "openId": "openId of mp service or mimi app"
 * }
 */
func createWxPay(w http.ResponseWriter, r *http.Request) {
	var jsapiParam struct {
		AppId    string
		PayApp   string
		Goods    string
		Udd      string
		OrderId  string
		Fee      int
		Ip       string
		OpenId   string
	}

	if code, err := _ReadJson(r, &jsapiParam); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	isSandbox := _IsSandbox(jsapiParam.PayApp)
	mchConf, _, ok := conf.GetAppAttrs(jsapiParam.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	prepayId, reqJSAPI, err := wxpay.JSAPIPay(
		jsapiParam.AppId,
		mchConf.MchId,
		mchConf.MchAppKey,
		jsapiParam.Goods,
		jsapiParam.Udd,
		jsapiParam.OrderId,
		jsapiParam.Fee,
		jsapiParam.Ip,
		_AppendAppName(conf.ServiceConf.NotifyPayUrl, jsapiParam.PayApp),
		jsapiParam.OpenId,
		isSandbox,
	)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_WriteJson(w, http.StatusOK, map[string]interface{} {
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]interface{} {
			"prepay_id": prepayId,
			"jsapi_params": reqJSAPI,
		},
	})
}

/**
 * Native pay
 * {
      "appId": "appId of mp/mini-prog",
      "payApp": "name-of-app-in-wxpay-gateway",
      "goods": "XXXX-xxxx",
      "udd": "user defined data as parameters when calling callback",
      "orderId": "unique order id",
      "fee": xxx-in-fen,
      "ip": "ip to create order",
      "productId": "productId"
 * }
 */
func createNativePay(w http.ResponseWriter, r *http.Request) {
	var nativeParam struct {
		AppId    string
		PayApp   string
		Goods    string
		Udd      string
		OrderId  string
		Fee      int
		Ip       string
		ProductId string
	}

	if code, err := _ReadJson(r, &nativeParam); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	isSandbox := _IsSandbox(nativeParam.PayApp)
	mchConf, _, ok := conf.GetAppAttrs(nativeParam.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	prepayId, codeUrl, err := wxpay.NativePay(
		nativeParam.AppId,
		mchConf.MchId,
		mchConf.MchAppKey,
		nativeParam.Goods,
		nativeParam.Udd,
		nativeParam.OrderId,
		nativeParam.Fee,
		nativeParam.Ip,
		_AppendAppName(conf.ServiceConf.NotifyPayUrl, nativeParam.PayApp),
		nativeParam.ProductId,
		isSandbox,
	)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_WriteJson(w, http.StatusOK, map[string]interface{} {
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]interface{} {
			"prepay_id": prepayId,
			"code_url": codeUrl,
		},
	})
}

/**
 * App pay
 * {
      "appId": "appId of mp/mini-prog",
      "payApp": "name-of-app-in-wxpay-gateway",
      "goods": "XXXX-xxxx",
      "udd": "user defined data as parameters when calling callback",
      "orderId": "unique order id",
      "fee": xxx-in-fen,
      "ip": "ip to create order"
 * }
 */
func createAppPay(w http.ResponseWriter, r *http.Request) {
	var appParam struct {
		AppId    string
		PayApp   string
		Goods    string
		Udd      string
		OrderId  string
		Fee      int
		Ip       string
	}

	if code, err := _ReadJson(r, &appParam); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	isSandbox := _IsSandbox(appParam.PayApp)
	mchConf, _, ok := conf.GetAppAttrs(appParam.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	prepayId, reqAppPay, err := wxpay.AppPay(
		appParam.AppId,
		mchConf.MchId,
		mchConf.MchAppKey,
		appParam.Goods,
		appParam.Udd,
		appParam.OrderId,
		appParam.Fee,
		appParam.Ip,
		_AppendAppName(conf.ServiceConf.NotifyPayUrl, appParam.PayApp),
		isSandbox,
	)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_WriteJson(w, http.StatusOK, map[string]interface{} {
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]interface{} {
			"prepay_id": prepayId,
			"req_params": reqAppPay,
		},
	})
}

/**
 * H5 pay
 * {
      "appId": "appId of mp/mini-prog",
      "payApp": "name-of-app-in-wxpay-gateway",
      "goods": "XXXX-xxxx",
      "udd": "user defined data as parameters when calling callback",
      "orderId": "unique order id",
      "fee": xxx-in-fen,
      "ip": "ip to create order",
      "redirectUrl": "url to redirect after payment",
      "sceneInfo": {
         "h5_info": {
            "type": "",
            "wap_name": "",
            "wap_url": "wap-site-url"
         }
            ---- OR ----
         "h5_info": {
            "type": "",
            "app_name":"",
            "bundle_id": "ios-bundle-id"
         }
            ---- OR ----
         "h5_info": {
            "type": "",
            "app_name":"",
            "package_name": "android-package-name"
         }
      }
 * }
 */
func createH5Pay(w http.ResponseWriter, r *http.Request) {
	var h5Param struct {
		AppId    string
		PayApp   string
		Goods    string
		Udd      string
		OrderId  string
		Fee      int
		Ip       string
		RedirectUrl string
		SceneInfo interface{}
	}

	if code, err := _ReadJson(r, &h5Param); err != nil {
		_WriteError(w, code, err.Error())
		return
	}
	if h5Param.SceneInfo == nil {
		_WriteError(w, http.StatusBadRequest, "sceneInfo not found")
		return
	}
	sceneInfo, err := json.Marshal(h5Param.SceneInfo)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	isSandbox := _IsSandbox(h5Param.PayApp)
	mchConf, _, ok := conf.GetAppAttrs(h5Param.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	prepayId, payUrl, err := wxpay.H5Pay(
		h5Param.AppId,
		mchConf.MchId,
		mchConf.MchAppKey,
		h5Param.Goods,
		h5Param.Udd,
		h5Param.OrderId,
		h5Param.Fee,
		h5Param.Ip,
		_AppendAppName(conf.ServiceConf.NotifyPayUrl, h5Param.PayApp),
		h5Param.RedirectUrl,
		sceneInfo,
		isSandbox,
	)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_WriteJson(w, http.StatusOK, map[string]interface{} {
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]interface{} {
			"prepay_id": prepayId,
			"pay_url": payUrl,
		},
	})
}
