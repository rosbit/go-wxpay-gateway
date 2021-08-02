package oauth

// 实名验证
//  - 输入：某个openId的姓名、身份证
//  - 功能：验证openId的姓名、身份证是否匹配
// 文档: https://pay.weixin.qq.com/wiki/doc/api/realnameauth.php?chapter=60_1&index=2

import (
	"github.com/rosbit/gnet"
	"go-wxpay-gateway/xml-msg"
	"go-wxpay-gateway/sign"
	"fmt"
)

type PayUserIdentity struct {
	PayAuth
}

func NewPayUserIdentity(mchId, appId, apiKey string) *PayUserIdentity {
	a := NewPayAuth(mchId, appId, PAY_IDENTITY, apiKey)
	return &PayUserIdentity{*a}
}

func (i *PayUserIdentity) RealNameAuth(openId, code, realName, cardId string) error {
	if err := i.GetAccessToken(openId, code); err != nil {
		return err
	}

	state, _ := GenerateState()
	params := map[string]string {
		"version": "1.0",
		"mch_id": i.mchId,
		"appid": i.appId,
		"openid": openId,
		"real_name": realName,
		"cred_type": "1",
		"cred_id": cardId,
		"nonce_str": state,
		"access_token": i.accessToken,
		"sign_type": _MD5,
	}
	params["sign"] = sign.CreateSignature(sign.MD5, params, i.apiKey)

	return i.auth(params, sign.MD5)
}

func (i *PayUserIdentity) auth(params map[string]string, signType string) error {
	xml := xmlmsg.NewXmlGenerator("xml")
	for t, v := range params {
		xml.AddTag(t, v)
	}
	xmlstr := xml.ToXML()

	_, content, _, err := gnet.Http("https://fraud.mch.weixin.qq.com/secsvc/realnameauth",
		gnet.M("POST"), gnet.Params(xmlstr), gnet.Headers(map[string]string{"Content-Type": "text/xml"}))
	if err != nil {
		return err
	}
	res, err := parseXmlResult(content, i.apiKey, signType)
	if err != nil {
		return err
	}

	if verifyOpenid, ok := res["verify_openid"]; !ok {
		return fmt.Errorf("no verify_openid found")
	} else {
		switch verifyOpenid {
		case "V_OP_NA":
			return fmt.Errorf("openId has not been real name authenticated")
		case "V_OP_NM_UM":
			return fmt.Errorf("name not matched")
		case "V_OP_NM_MA":
		default:
		}
	}

	if verifyRealName, ok := res["verify_real_name"]; !ok {
		return fmt.Errorf("no verify_real_name found")
	} else {
		switch verifyRealName {
		case "V_NM_ID_UM":
			return fmt.Errorf("name and id not matched")
		case "V_NM_ID_MA":
		default:
		}
		return nil
	}
}
