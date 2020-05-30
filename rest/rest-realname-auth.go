package rest

import (
	"go-wxpay-gateway/wx-pay-auth"
	"go-wxpay-gateway/conf"
	"net/http"
	"fmt"
)

const (
	REALNAME_TYPE_NAME = "op"

	REALNAME_URL      = "url"
	REALNAME_IDENTITY = "identity"
	REALNAME_GETINFO  = "getinfo"
)

var (
	authOps = map[string]http.HandlerFunc {
		REALNAME_URL: authRealnameUrl,
		REALNAME_IDENTITY: authRealnameIdentity,
		REALNAME_GETINFO: authRealnameGetinfo,
	}

	scopes = map[string]string {
		REALNAME_IDENTITY: oauth.PAY_IDENTITY,
		REALNAME_GETINFO:  oauth.PAY_REALNAME,
	}
)

// POST /realname/auth/:op
func AuthRealname(w http.ResponseWriter, r *http.Request) {
	opName := _PathParam(r, REALNAME_TYPE_NAME)
	if op, ok := authOps[opName]; ok {
		op(w, r)
		return
	}
	_WriteError(w, http.StatusBadRequest, fmt.Sprintf("unknown realname auth op %s", opName))
}

// POST /realname/auth/url
// BODY:
// {
//     "payApp": "name-of-app-in-wxpay-gateway",
//     "appId": "appId of mp/mini-prog",
//     "type": "identity|getinfo",
//     "redirectUrl": "GET形式的回调url，外网可以访问，微信支付服务会加上参数'?code=xxx&state=xxx'"
// }
//
// 1. 用户在统一获取实名信心或在网站填写“姓名”、“身份证”，提交
// 2. 表单处理服务调用本接口得到一个服务微信支付服务的url
// 3. 把第2步得到的url通过 302 的返回码 让用户浏览器跳转过去
// 4. 用户授权，则浏览器会跳转到参数中的redirectUrl，结合AuthRealnameIdentity()/authRealnameGetinfo()完成余下的步骤
func authRealnameUrl(w http.ResponseWriter, r *http.Request) {
	var authRealnameParams struct {
		PayApp string
		AppId  string
		Type   string
		RedirectUrl string
	}
	if code, err := _ReadJson(r, &authRealnameParams); err != nil {
		_WriteError(w, code, err.Error())
		return
	}

	scope, ok := scopes[authRealnameParams.Type]
	if !ok {
		_WriteError(w, http.StatusBadRequest, fmt.Sprintf("unknown type name %s", authRealnameParams.Type))
		return
	}
	mchConf, ok := conf.GetAppAttrs(authRealnameParams.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}

	_WriteJson(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"url": oauth.MakeAuthindexUrl(
			mchConf.MchId,
			authRealnameParams.AppId,
			authRealnameParams.RedirectUrl,
			scope,
		),
	})
}

// POST /realname/auth/identity
// BODY:
// {
//    "payApp": "name-of-app-in-wxpay-gateway",
//    "appId": "appId of mp/mini-prog",
//    "openId": "要验证的用户的openId",
//    "realName": "要验证的真实姓名",
//    "cardId": "要验证的身份证号码",
//    "requstURI": "访问redirectUrl的完整的URI"
// }
//
// 实名验证: 查看openId的realName/cardId是否真实
func authRealnameIdentity(w http.ResponseWriter, r *http.Request) {
	var verifyParams struct {
		PayApp string
		AppId  string
		OpenId string
		RealName string
		CardId   string
		RequestURI string
	}
	if code, err := _ReadJson(r, &verifyParams); err != nil {
		_WriteError(w, code, err.Error())
		return
	}
	if len(verifyParams.OpenId) == 0 || len(verifyParams.RealName) == 0 || len(verifyParams.CardId) == 0 {
		_WriteError(w, http.StatusBadRequest, "openId, realname and cardId must be specified")
		return
	}
	if len(verifyParams.AppId) == 0 {
		_WriteError(w, http.StatusBadRequest, "appId expected")
		return
	}
	if len(verifyParams.RequestURI) == 0 {
		_WriteError(w, http.StatusBadRequest, "requestURI expected")
		return
	}

	mchConf, ok := conf.GetAppAttrs(verifyParams.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}
	code, err := oauth.VerifyStateAndGetCode(verifyParams.RequestURI)
	if err != nil {
		_WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	i := oauth.NewPayUserIdentity(mchConf.MchId, verifyParams.AppId, mchConf.MchApiKey)
	if err := i.RealNameAuth(verifyParams.OpenId, code, verifyParams.RealName, verifyParams.CardId); err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_WriteJson(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
	})
}

// POST /realname/auth/getinfo
// BODY:
// {
//    "payApp": "name-of-app-in-wxpay-gateway",
//    "appId": "appId of mp/mini-prog",
//    "openId": "要验证的用户的openId",
//    "requstURI": "访问redirectUrl的完整的URI"
// }
//
// 实名授权: 获取openId的realName/cardId
func authRealnameGetinfo(w http.ResponseWriter, r *http.Request) {
	var getinfoParams struct {
		PayApp string
		AppId  string
		OpenId string
		RequestURI string
	}
	if code, err := _ReadJson(r, &getinfoParams); err != nil {
		_WriteError(w, code, err.Error())
		return
	}
	if len(getinfoParams.OpenId) == 0 {
		_WriteError(w, http.StatusBadRequest, "openId must be specified")
		return
	}
	if len(getinfoParams.AppId) == 0 {
		_WriteError(w, http.StatusBadRequest, "appId expected")
		return
	}
	if len(getinfoParams.RequestURI) == 0 {
		_WriteError(w, http.StatusBadRequest, "requestURI expected")
		return
	}

	mchConf, ok := conf.GetAppAttrs(getinfoParams.PayApp)
	if !ok {
		_WriteError(w, http.StatusBadRequest, "Unknown pay-app name")
		return
	}
	code, err := oauth.VerifyStateAndGetCode(getinfoParams.RequestURI)
	if err != nil {
		_WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	a := oauth.NewPayUserGetinfo(mchConf.MchId, getinfoParams.AppId, mchConf.MchApiKey)
	timestamp, realName, cardid, err := a.Getinfo(getinfoParams.OpenId, code, mchConf.MchCertSerialNo, mchConf.MchKeyPemFile)
	if err != nil {
		_WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_WriteJson(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]interface{}{
			"timestamp": timestamp,
			"realName": realName,
			"cardid": cardid,
		},
	})
}
