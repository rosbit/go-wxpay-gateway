// +build notify

/**
 * wxpay-notify implementation.
 */
package main

import (
	lm "github.com/rosbit/logmerger"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/utils"
)

func StartService() error {
	if err := utils.InitNotifyLog(conf.NotifyConf.NotifyLogFile); err != nil {
		return err
	}
	utils.StartNotifyThreads()
	m := lm.NewLogMerger(conf.NotifyConf.TimeInterval)
	m.Run(conf.NotifyConf.NotifyFile, utils.Notify) // block and call notify() if there's data in log-file
	return nil
}

