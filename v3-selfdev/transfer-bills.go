// 商家转账
package v3sd

import (
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/rosbit/gnet"
	"go-wxpay-gateway/conf"
	"encoding/json"
	"fmt"
	"net/http"
)

type createTransferBillsRequest struct {
	*CreateTransferBillsRequest
	EncryptedUserName string `json:"user_name,omitempty"`
}

// 发起转账: https://pay.weixin.qq.com/doc/v3/merchant/4012716434
func CreateTransferBills(appName string, req *CreateTransferBillsRequest) (resp json.RawMessage, err error) {
	mchConf, ok := conf.GetAppAttrs(appName)
	if !ok {
		err = fmt.Errorf("conf not found for %s", appName)
		return
	}

	if len(mchConf.MchCertSerialNo) == 0 {
		err = fmt.Errorf("mch-cert-serialno specified for %s", appName)
		return
	}
	if len(mchConf.WxpayPubkeyId) == 0 {
		err = fmt.Errorf("wxpay-pubkey-id specified for %s", appName)
		return
	}
	if len(mchConf.WxpayPubkeyFile) == 0 {
		err = fmt.Errorf("wxpay-pubkey-file not specified for %s", appName)
		return
	}

	wechatpayPublicKey, e := utils.LoadPublicKeyWithPath(mchConf.WxpayPubkeyFile)
	if e != nil {
		err = e
		return
	}
	mchPrivateKey, e := utils.LoadPrivateKeyWithPath(mchConf.MchKeyPemFile)
	if e != nil {
		err = e
		return
	}

	realReq := &createTransferBillsRequest{
		CreateTransferBillsRequest: req,
	}

	if len(req.UserName) > 0 {
		encryptedUser, e := utils.EncryptPKCS1v15WithPublicKey(req.UserName, wechatpayPublicKey)
		if e != nil {
			err = e
			return
		}
		realReq.EncryptedUserName = encryptedUser
	}

	timestamp, nonce, bodyStr, signature, e := MakeSignature(mchPrivateKey, "POST", "/v3/fund-app/mch-transfer/transfer-bills", realReq, true)
	if e != nil {
		err = e
		return
	}

	status, content, _, e := gnet.JSON("https://api.mch.weixin.qq.com/v3/fund-app/mch-transfer/transfer-bills",
		gnet.M("POST"),
		gnet.Params(bodyStr),
		gnet.Headers(map[string]string{
			"Content-Type":     "application/json",
			"Accept":           "application/json",
			// "Wechatpay-Serial": mchConf.WxpayPubkeyId,
			"Authorization": fmt.Sprintf(
				"WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%s\",serial_no=\"%s\",signature=\"%s\"",
				mchConf.MchId, nonce, timestamp, mchConf.MchCertSerialNo, signature),
		}),
	)
	if e != nil {
		err = e
		return
	}
	if status != http.StatusOK {
		err = fmt.Errorf(string(content))
		return
	}
	if err = json.Unmarshal(content, &resp); err != nil {
		return
	}

	return
}
