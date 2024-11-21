// +build gateway

/**
 * global conf
 * ENV:
 *   CONF_FILE      --- 配置文件名
 *   TZ             --- 时区名称"Asia/Shanghai"
 *
 * JSON:
 * {
      "listen-host": "",
      "listen-port": 7080,
      "worker-num": 5,
      "timeout": 0,
      "order-expire-minutes": 120,
      "endpoints": {
         "health-check": "/health",
         "create-pay": "/create-pay",
         "create-refund": "/create-refund",
         "query-order": "/query-order",
         "close-order": "/close-order",
         "transfer": "/transfer",
         "query-transfer": "/query-transfer",
         "v3-transfer": "/v3/transfer",
         "v3-query-transfer": "/v3/query-transfer",
         "v3-query-transfer-detail": "/v3/query-transfer-detail",
         "realname-auth-root": "/realname/auth -- will be replaced with ${realname-auth-root}/:op {url|identity|getinfo}",
         "verify-notify-pay": "/verify-notify-pay",
         "verify-notify-refund": "/verify-notify-refund"
      },
      "merchants": [
         {
             "name": "mch1",
             "mch-id": "", 
             "mch-api-key": "",
             "mch-cert-pem-file": "your-cert-pem-file-name, only used when refunding",
             "mch-key-pem-file": "your-key-pem-file-name, only used when refunding",
             "mch-cert-serialno": "optional, only used when real-name getinfo",
			 "wxpay-v3-cert": "optional, only used for batch-transfer",
             "receipt": true
         }
      ],
      "apps": [
         {
             "name": "app1",
             "merchant": "mch1, the name in merchants",
         },
         {
             "name": "app1-dev, name with suffix -dev will be call wxpay in sandbox",
             "merchant": "mch1, the name in merchants",
         }
      ]
   }
 *
 * Rosbit Xu
 */
package conf

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
)

const (
	NOTIFY_PAY_CB_IDX    = 0
	NOTIFY_REFUND_CB_IDX = 1
	CERT_PEM_IDX = 0
	KEY_PEM_IDX  = 1
)

type EndpointConf struct {
	HealthCheck   string `json:"health"`
	CreatePay     string `json:"create-pay"`
	CreateRefund  string `json:"create-refund"`
	QueryOrder    string `json:"query-order"`
	CloseOrder    string `json:"close-order"`
	Transfer      string `json:"transfer"`
	V3Transfer    string `json:"v3-transfer"`
	V3QueryTransfer string `json:"v3-query-transfer"`
	V3QueryTransferDetail string `json:"v3-query-transfer-detail"`
	QueryTransfer string `json:"query-transfer"`
	RealnameAuthRoot   string `json:"realname-auth-root"`
	VerifyNotifyPay    string `json:"verify-notify-pay"`
	VerifyNotifyRefund string `json:"verify-notify-refund"`
}

type MerchantConf struct {
	Name      string
	MchId     string `json:"mch-id"`                   // 微信商户号
	MchApiKey string `json:"mch-api-key"`              // 商户APIv2密钥
	MchCertPemFile  string `json:"mch-cert-pem-file"`  // 商户API证书: 公钥
	MchKeyPemFile   string `json:"mch-key-pem-file"`   // 商户API证书: 私钥
	MchCertSerialNo string `json:"mch-cert-serialno"`  // 商户证书序列号。
	WxpayV3Cert     string `json:"wxpay-v3-cert"`      // APIv3平台证书。[该证书5年有效期，平台证书模式逐渐弃用]
	WxpayV3Key      string `json:"wxpay-v3-key"`       // APIv3密钥。[当WxpayV3Cert为空时有效，表示使用微信支付公钥模式]
	WxpayPubkeyFile string `json:"wxpay-pubkey-file"`  // 微信支付公钥[当WxpayV3Cert为空时有效，表示使用微信支付公钥模式]
	WxpayPubkeyId   string `json:"wxpay-pubkey-id"`    // 微信支付公钥ID[当WxpayV3Cert为空时有效，表示使用微信支付公钥模式]
	Receipt bool `json:"receipt"`
}

type PayApps struct {
	Name      string
	Merchant  string
}

type WxPayServiceConf struct {
	ListenHost string `json:"listen-host"`
	ListenPort int    `json:"listen-port"`
	WorkerNum  int    `json:"worker-num"`
	Timeout    int
	OrderExpireMinutes int64 `json:"order-expire-minutes"`
	Endpoints  EndpointConf
	Merchants  []MerchantConf
	Apps       []PayApps
}

var (
	ServiceConf WxPayServiceConf
	Merchants map[string]*MerchantConf
	Apps map[string]*PayApps
	Loc = time.FixedZone("UTC+8", 8*60*60)
)

func getEnv(name string, result *string, must bool) error {
	s := os.Getenv(name)
	if s == "" {
		if must {
			return fmt.Errorf("env \"%s\" not set", name)
		}
	}
	*result = s
	return nil
}

func CheckGlobalConf() error {
	var p string
	getEnv("TZ", &p, false)
	if p != "" {
		if loc, err := time.LoadLocation(p); err == nil {
			Loc = loc
		}
	}

	var confFile string
	if err := getEnv("CONF_FILE", &confFile, true); err != nil {
		return err
	}

	fp, err := os.Open(confFile)
	if err != nil {
		return err
	}
	defer fp.Close()
	if err = json.NewDecoder(fp).Decode(&ServiceConf); err != nil {
		return err
	}

	Merchants = make(map[string]*MerchantConf, len(ServiceConf.Merchants))
	for i, _ := range ServiceConf.Merchants {
		merchant := &ServiceConf.Merchants[i]
		Merchants[merchant.Name] = merchant
	}

	Apps = make(map[string]*PayApps, len(ServiceConf.Apps))
	for i, _ := range ServiceConf.Apps {
		app := &ServiceConf.Apps[i]
		Apps[app.Name] = app
	}

	if ServiceConf.OrderExpireMinutes <= 0 {
		ServiceConf.OrderExpireMinutes = 120
	}

	return nil
}

func GetAppAttrs(appName string) (mchConf *MerchantConf, ok bool) {
	var app *PayApps
	if app, ok = Apps[appName]; !ok {
		return
	}
	if mchConf, ok = Merchants[app.Merchant]; !ok {
		return
	}
	return
}

func DumpConf() {
	fmt.Printf("conf: %#v\n", ServiceConf)
	fmt.Printf("TZ time location: %v\n", Loc)
}
