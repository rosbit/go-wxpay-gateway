package wxpay

import (
	"log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	_paymentLog *log.Logger
)

func InitPaymentLog(logFile string) error {
	_paymentLogFile := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   false, // disabled by default
	}
	_paymentLog = log.New(_paymentLogFile, "", log.LstdFlags)
	return nil
}
