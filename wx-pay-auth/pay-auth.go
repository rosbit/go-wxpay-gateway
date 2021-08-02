package oauth

// 实名验证 和 授权 的公共部分
//  - 获取token、刷新token

import (
	"github.com/rosbit/gnet"
	"go-wxpay-gateway/sign"
	"net/url"
	"fmt"
	"time"
)

const (
	PAY_IDENTITY = "pay_identity"
	PAY_REALNAME = "pay_realname"

	_GT_AUTH_CODE     = "authorization_code"
	_GT_REFRESH_TOKEN = "refresh_token"

	_HMAC_SHA256  = "HMAC-SHA256"
	_MD5          = "MD5"
)

type PayAuth struct {
	mchId string
	appId string
	scope  string
	apiKey string
	openId string
	accessToken string
	expireTime int64
	refreshToken string
	refreshExpireTime int64
}

func MakeAuthindexUrl(mchId, appId, rdUrl, scope string) string {
	state, _ := GenerateState()
	v := url.Values{}
	v.Set("mch_id", mchId)
	v.Set("appid",  appId)
	v.Set("redirect_uri", rdUrl)
	v.Set("state", state)
	v.Set("scope", scope)

	return fmt.Sprintf("https://payapp.weixin.qq.com/appauth/authindex?%s&response_type=code#wechat_redirect", v.Encode())
}

func VerifyStateAndGetCode(requestURI string) (code string, err error) {
	var req *url.URL
	if req, err = url.ParseRequestURI(requestURI); err != nil {
		return
	}
	q := req.Query()
	code = q.Get("code")
	state := q.Get("state")
	err = VerifyState(state)
	return
}

func NewPayAuth(mchId, appId, scope, apiKey string) *PayAuth {
	return &PayAuth{mchId:mchId, appId:appId, scope:scope, apiKey:apiKey}
}

func (a *PayAuth) GetAccessToken(openId, code string) error {
	if openId == a.openId {
		if err := a.refresh(); err == nil {
			return nil
		}
	}

	params := map[string]string {
		"mch_id": a.mchId,
		"appid": a.appId,
		"openid": openId,
		"code": code,
		"grant_type": _GT_AUTH_CODE,
		"scope": a.scope,
		"sign_type": _HMAC_SHA256,
	}
	params["sign"] = sign.CreateSignature(sign.HMAC_SHA256, params, a.apiKey)
	err := a.getAccessToken("https://api.mch.weixin.qq.com/appauth/getaccesstoken", params)
	if err != nil {
		return err
	}
	a.openId = openId
	return nil
}

func (a *PayAuth) getAccessToken(url string, params map[string]string) error {
	var token struct {
		Retcode      int
		Retmsg       string
		AccessToken  string  `json:"access_token"`
		ExpiresIn    int64   `json:"access_token_expires_in"`
		RefreshToken string  `json:"refresh_token"`
		RefreshExpiresIn int64 `json:"refresh_token_expires_in"`
	}
	if _, err := gnet.HttpCallJ(url, &token, gnet.Params(params)); err != nil {
		return err
	}
	if token.Retcode != 0 {
		return fmt.Errorf("%d: %s", token.Retcode, token.Retmsg)
	}

	a.accessToken  = token.AccessToken
	a.expireTime   = time.Now().Unix() + token.ExpiresIn - 10
	a.refreshToken = token.RefreshToken
	a.refreshExpireTime = time.Now().Unix() + token.RefreshExpiresIn - 10
	return nil
}

func (a *PayAuth) tokenExpired() bool {
	return a.expireTime == 0 || time.Now().Unix() >= a.expireTime
}

func (a *PayAuth) refreshTokenExpired() bool {
	return a.refreshExpireTime == 0 || time.Now().Unix() >= a.refreshExpireTime
}

func (a *PayAuth) refresh() error {
	if !a.tokenExpired() {
		return nil
	}
	if a.refreshTokenExpired() {
		return fmt.Errorf("refresh token expired")
	}

	params := map[string]string {
		"mch_id": a.mchId,
		"appid":  a.appId,
		"openid": a.openId,
		"scope":  a.scope,
		"grant_type": _GT_REFRESH_TOKEN,
		"refresh_token": a.refreshToken,
		"sign_type": _HMAC_SHA256,
	}
	params["sign"] = sign.CreateSignature(sign.HMAC_SHA256, params, a.apiKey)

	return a.getAccessToken("https://api.mch.weixin.qq.com/appauth/refreshtoken", params)
}
