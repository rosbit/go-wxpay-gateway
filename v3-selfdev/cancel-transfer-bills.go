// 撤销转账
// 参考文档: https://pay.weixin.qq.com/doc/v3/merchant/4012716458
package v3sd

import (
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/rosbit/gnet"
	"go-wxpay-gateway/conf"
	"encoding/json"
	"fmt"
	"net/http"
)

func CancelTransferBills(appName string, outBillNo string) (resp json.RawMessage, err error) {
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

	mchPrivateKey, e := utils.LoadPrivateKeyWithPath(mchConf.MchKeyPemFile)
	if e != nil {
		err = e
		return
	}

	uri := fmt.Sprintf("/v3/fund-app/mch-transfer/transfer-bills/out-bill-no/%s/cancel", outBillNo)
	timestamp, nonce, _, signature, e := MakeSignature(mchPrivateKey, "POST", uri, nil, true)
	if e != nil {
		err = e
		return
	}

	status, content, _, e := gnet.Http(fmt.Sprintf("https://api.mch.weixin.qq.com%s", uri),
		gnet.M("POST"),
		gnet.Params(nil),
		gnet.Headers(map[string]string{
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
