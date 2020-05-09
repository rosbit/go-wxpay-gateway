// +build notify

/**
 * global conf
 * ENV:
 *   CONF_FILE      --- 配置文件名
 *   TZ             --- 时区名称"Asia/Shanghai"
 *
 * JSON of <CONF_FILE>:
 * {
      "time-interval": "1000",
      "notify-file": "file-path-to-save-notify-message",
      "worker-num": 5,
      "retry-count": 3,
      "notify-log-file": "log-file-related-to-notify"
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
	_DEFAULT_RETRY_COUNT   = 3
	_DEFAULT_WORKER_NUM    = 5
	_DEFAULT_TIME_INTERVAL = 100
)

type WxNotifyConf struct {
	TimeInterval  int    `json:"time-interval"`
	NotifyFile    string `json:"notify-file"`
	WorkerNum     int    `json:"worker-num"`
	RetryCount    int    `json:"retry-count"`
	NotifyLogFile string `json:"notify-log-file"`
}

var (
	NotifyConf WxNotifyConf
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
	if err = json.NewDecoder(fp).Decode(&NotifyConf); err != nil {
		return err
	}

	if NotifyConf.TimeInterval <= 0 {
		NotifyConf.TimeInterval = _DEFAULT_TIME_INTERVAL
	}
	if NotifyConf.WorkerNum <= 0 {
		NotifyConf.WorkerNum = _DEFAULT_WORKER_NUM
	}
	if NotifyConf.RetryCount <= 0 {
		NotifyConf.RetryCount = _DEFAULT_RETRY_COUNT
	}
	return nil
}

func DumpConf() {
	fmt.Printf("conf: %#v\n", NotifyConf)
	fmt.Printf("TZ time location: %v\n", Loc)
}
