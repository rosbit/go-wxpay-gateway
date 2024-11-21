package v3

import (
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/rosbit/gnet"
	"go-wxpay-gateway/conf"
	"fmt"
	"context"
	"crypto/x509"
)

func createClient(appName string) (ctx context.Context, client *core.Client, err error) {
	mchConf, ok := conf.GetAppAttrs(appName)
	if !ok {
		err = fmt.Errorf("conf not found for %s", appName)
		return
	}

	opts, e := createOptionsWithPlatformMode(appName, mchConf)
	if e != nil {
		err = e
		return
	}
	if len(opts) == 0 {
		opts, e = createOptionsWithPubkeyMode(appName, mchConf)
		if e != nil {
			err = e
			return
		}
		if len(opts) == 0 {
			err = fmt.Errorf("no platform mode / pubkey mode params fro %s", appName)
			return
		}
	}

	ctx = context.Background()
	client, err = core.NewClient(ctx, opts...)
	return
}

// 使用平台证书模式创建选项
// @return err  非nil，有参数错误
// @return opts 在err为nil时有效：非nil表示成功；nil表示没有相关选项
func createOptionsWithPlatformMode(appName string, mchConf *conf.MerchantConf) (opts []core.ClientOption, err error) {
	if len(mchConf.WxpayV3Cert) == 0 {
		// 无相关参数
		return
	}

	if len(mchConf.MchCertSerialNo) == 0 {
		err = fmt.Errorf("mch-cert-serialno not specified for %s", appName)
		return
	}
	/*
	if len(mchConf.WxpayV3Cert) == 0 {
		err = fmt.Errorf("wxpay-v3-cert not specified for %s", appName)
		return
	}*/
	cert, e := utils.LoadCertificateWithPath(mchConf.WxpayV3Cert)
	if e != nil {
		err = e
		return
	}

	mchPrivateKey, e := utils.LoadPrivateKeyWithPath(mchConf.MchKeyPemFile)
	if e != nil {
		err = e
		return
	}

	// opt := option.WithWechatPayAutoAuthCipher(mchConf.MchId, mchConf.MchCertSerialNo, mchPrivateKey, mchConf.MchApiKey)
	opts = []core.ClientOption{
		option.WithWechatPayAuthCipher(mchConf.MchId, mchConf.MchCertSerialNo, mchPrivateKey, []*x509.Certificate{cert}),
		option.WithHTTPClient(gnet.NewHttpsRequest().GetClient()),
	}
	return
}

// 使用微信支付公钥模式创建选项
// @return err  非nil，有参数错误
// @return opts 在err为nil时有效：非nil表示成功；nil表示没有相关选项
func createOptionsWithPubkeyMode(appName string, mchConf *conf.MerchantConf) (opts []core.ClientOption, err error) {
	if len(mchConf.WxpayPubkeyFile) == 0 && len(mchConf.WxpayPubkeyId) == 0 && len(mchConf.MchCertSerialNo) == 0 {
		// 无相关参数
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

	// opt := option.WithWechatPayAutoAuthCipher(mchConf.MchId, mchConf.MchCertSerialNo, mchPrivateKey, mchConf.MchApiKey)
	opts = []core.ClientOption{
		option.WithWechatPayPublicKeyAuthCipher(mchConf.MchId, mchConf.MchCertSerialNo, mchPrivateKey, mchConf.WxpayPubkeyId, wechatpayPublicKey),
		option.WithHTTPClient(gnet.NewHttpsRequest().GetClient()),
	}
	return
}
