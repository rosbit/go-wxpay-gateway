// +build getsandbox

// command tool getsandbox
package main

import (
	"go-wxpay-gateway/wx-pay-api"
	"go-wxpay-gateway/conf"
	"fmt"
)

func StartService() error {
	sandbox_signkey, err := wxpay.GetSandbox(conf.AppId, conf.MchId, conf.MchApiKey)
	if err != nil {
		fmt.Printf("failed: %v\n", err)
	} else {
		fmt.Printf("sandbox_signkey: %v\n", sandbox_signkey)
	}
	return nil
}
