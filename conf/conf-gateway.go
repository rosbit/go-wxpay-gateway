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
      "pay-log-file": "",
      "notify-pay-url": "outerside URL mapped to notify-pay-callback",
      "notify-refund-url": "outerside URL mapped to notify-refund-callback",
      "notify-file": "file-path-to-save-notify-message",
      "endpoints": {
         "create-pay": "/create-pay",
         "notify-pay": "/notify-pay",
         "create-refund": "/create-refund",
         "notify-refund": "/notify-refund",
         "query-order": "/query-order",
         "close-order": "/close-order"
      },
      "merchants": [
         {
             "name": "mch1",
             "mch-id": "", 
             "mch-app-key": "",
             "mch-cert-pem-file": "your-cert-pem-file-name, only used when refunding",
             "mch-key-pem-file": "your-key-pem-file-name, only used when refunding"
         }
      ],
      "apps": [
         {
             "name": "app1",
             "merchant": "mch1, the name in merchants",
             "notify-pay-callback": "a-url-to-receive-payment-callback-with-post-method",
             "notify-refund-callback": "a-url-to-receive-refundment-callback-with-post-method",
         },
         {
             "name": "app1-dev, name with suffix -dev will be call wxpay in sandbox",
             "merchant": "mch1, the name in merchants",
             "notify-pay-callback": "a-url-to-receive-payment-callback-with-post-method",
             "notify-refund-callback": "a-url-to-receive-refundment-callback-with-post-method",
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
	"io/ioutil"
	"encoding/json"
)

const (
	NOTIFY_PAY_CB_IDX    = 0
	NOTIFY_REFUND_CB_IDX = 1
	CERT_PEM_IDX = 0
	KEY_PEM_IDX  = 1
)

type EndpointConf struct {
	CreatePay     string `json:"create-pay"`
	NotifyPay     string `json:"notify-pay"`
	CreateRefund  string `json:"create-refund"`
	NotifyRefund  string `json:"notify-refund"`
	QueryOrder    string `json:"query-order"`
	CloseOrder    string `json:"close-order"`
}

type MerchantConf struct {
	Name      string
	MchId     string `json:"mch-id"`
	MchAppKey string `json:"mch-app-key"`
	MchCertPemFile string `json:"mch-cert-pem-file"`
	MchKeyPemFile  string `json:"mch-key-pem-file"`
}

type PayApps struct {
	Name      string
	Merchant  string
	NotifyPayCallback    string `json:"notify-pay-callback"`
	NotifyRefundCallback string `json:"notify-refund-callback"`
}

type WxPayServiceConf struct {
	ListenHost string `json:"listen-host"`
	ListenPort int    `json:"listen-port"`
	WorkerNum  int    `json:"worker-num"`
	Timeout    int
	PayLogFile      string `json:"pay-log-file"`
	NotifyPayUrl    string `json:"notify-pay-url"`
	NotifyRefundUrl string `json:"notify-refund-url"`
	NotifyFile      string `json:"notify-file"`
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

	b, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &ServiceConf); err != nil {
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

	return nil
}

func GetAppAttrs(appName string) (mchConf *MerchantConf, cbUrls []string, ok bool) {
	if app, ok := Apps[appName]; !ok {
		return nil, nil, false
	} else {
		merchant, ok := Merchants[app.Merchant]
		if !ok {
			return nil, nil, false
		}
		return merchant,
			[]string{app.NotifyPayCallback, app.NotifyRefundCallback},
			true
	}
}

func DumpConf() {
	fmt.Printf("conf: %#v\n", ServiceConf)
	fmt.Printf("TZ time location: %v\n", Loc)
}
