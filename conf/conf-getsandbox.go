// +build getsandbox

package conf

import (
	"os"
	"fmt"
)

var (
	AppId string
	MchId string
	MchApiKey string
)

func CheckGlobalConf() error {
	if len(os.Args) < 4 {
		return fmt.Errorf("Usage: %s <appId> <mchId> <mchApiKey>", os.Args[0])
	}

	AppId, MchId, MchApiKey = os.Args[1], os.Args[2], os.Args[3]
	return nil
}

func DumpConf() {
}
