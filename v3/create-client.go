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
	if len(mchConf.MchCertSerialNo) == 0 {
		err = fmt.Errorf("mch-cert-serialno not specified for %s", appName)
		return
	}
	if len(mchConf.WxpayV3Cert) == 0 {
		err = fmt.Errorf("wxpay-v3-cert not specified for %s", appName)
		return
	}
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

	ctx = context.Background()
	// opt := option.WithWechatPayAutoAuthCipher(mchConf.MchId, mchConf.MchCertSerialNo, mchPrivateKey, mchConf.MchApiKey)
	opts := []core.ClientOption{
		option.WithWechatPayAuthCipher(mchConf.MchId, mchConf.MchCertSerialNo, mchPrivateKey, []*x509.Certificate{cert}),
		option.WithHTTPClient(gnet.NewHttpsRequest().GetClient()),
	}
	client, err = core.NewClient(ctx, opts...)
	return
}
