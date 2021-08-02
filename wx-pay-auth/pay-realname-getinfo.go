package oauth

// 实名授权
//  - 输入：某个openId
//  - 功能：获取openId的姓名、身份证
// 文档: https://pay.weixin.qq.com/wiki/doc/api/realnameauth.php?chapter=60_3&index=40

import (
	"github.com/rosbit/gnet"
	"go-wxpay-gateway/xml-msg"
	"go-wxpay-gateway/sign"
	"fmt"
)

type PayUserGetinfo struct {
	PayAuth
}

func NewPayUserGetinfo(mchId, appId, apiKey string) *PayUserGetinfo{
	a := NewPayAuth(mchId, appId, PAY_REALNAME, apiKey)
	return &PayUserGetinfo{*a}
}

func (r *PayUserGetinfo) Getinfo(openId, code, certSerialNo, keyFile string) (timestamp uint32, realName, cardid string, err error) {
	if err = r.GetAccessToken(openId, code); err != nil {
		return
	}
	var state, certSign string
	state, timestamp = GenerateState()
	if certSign, err = r.signWithCert(keyFile, certSerialNo, timestamp); err != nil {
		return
	}

	params := map[string]string {
		"version": "1.0",
		"mch_id": r.mchId,
		"appid": r.appId,
		"openid": openId,
		"cert_serialno": certSerialNo,
		"access_token": r.accessToken,
		"timestamp": fmt.Sprintf("%d", timestamp),
		"cert_sign": certSign,
		"nonce_str": state,
		"sign_type": _MD5,
	}
	params["sign"] = sign.CreateSignature(sign.MD5, params, r.apiKey)

	realName, cardid, err = r.getinfo(keyFile, params, timestamp, sign.MD5)
	return
}

func (r *PayUserGetinfo) signWithCert(keyFile, certSerialNo string, timestamp uint32) (string, error) {
	signer, err := sign.CreateSha256Base64Signer(keyFile)
	if err != nil {
		return "", err
	}
	signStr := fmt.Sprintf("%s&&timestamp=%d", certSerialNo, timestamp)
	return signer.Sign([]byte(signStr))
}

func (r *PayUserGetinfo) getinfo(keyFile string, params map[string]string, timestamp uint32, signType string) (realName, cardid string, err error) {
	xml := xmlmsg.NewXmlGenerator("xml")
	for t, v := range params {
		xml.AddTag(t, v)
	}
	xmlstr := xml.ToXML()

	var content []byte
	if _, content, _, err = gnet.Http("https://fraud.mch.weixin.qq.com/secsvc/getrealnameinfo",
		gnet.M("POST"), gnet.Params(xmlstr), gnet.Headers(map[string]string{"Content-Type": "text/xml"}),
	); err != nil {
		return
	}
	var res map[string]string
	if res, err = parseXmlResult(content, r.apiKey, signType); err != nil {
		return
	}

	eRealName, ok := res["encrypted_real_name"]
	if !ok {
		err = fmt.Errorf("encrypted_real_name not found")
		return
	}
	eCardId, ok := res["encrypted_credential_id"]
	if !ok {
		err = fmt.Errorf("encrypted_credential_id not found")
		return
	}
	return eRealName, eCardId, nil
}
