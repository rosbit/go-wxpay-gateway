package wxpay

import (
	"log"
	"os"
)

var (
	_paymentLogFile *os.File
	_paymentLog *log.Logger
)

func InitPaymentLog(logFile string) error {
	var err error
	_paymentLogFile, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_paymentLog = log.New(_paymentLogFile, "", log.LstdFlags)
	return nil
}
