// +build gateway

/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/urfave/negroni"
	"github.com/gernest/alien"
	"net/http"
	"fmt"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/rest"
	"go-wxpay-gateway/wx-pay-api"
	"go-wxpay-gateway/utils"
)

func StartService() error {
	serviceConf := &conf.ServiceConf
	wxpay.InitPaymentLog(serviceConf.PayLogFile)
	utils.StartSaver(serviceConf.NotifyFile)

	api := negroni.New()
	api.Use(negroni.NewRecovery())
	api.Use(negroni.NewLogger())

	router := alien.New()
	endpoints := serviceConf.Endpoints

	// set router
	var appendEndpointParam = func(uri, param string) string {
		l := len(uri)
		if uri[l-1] == '/' {
			return fmt.Sprintf("%s:%s", uri, param)
		}
		return fmt.Sprintf("%s/:%s", uri, param)
	}
	paymentEndpoint := appendEndpointParam(endpoints.CreatePay, rest.TRADE_TYPE_NAME)
	notifyPayEndpoint := appendEndpointParam(endpoints.NotifyPay, rest.NOTIFY_APP_NAME)
	notifyRefundEndpoint := appendEndpointParam(endpoints.NotifyRefund, rest.NOTIFY_APP_NAME)
	router.Post(paymentEndpoint,        rest.CreatePayment)
	router.Post(notifyPayEndpoint,      rest.NotifyPayment)
	router.Post(endpoints.CreateRefund, rest.CreateRefundment)
	router.Post(notifyRefundEndpoint,   rest.NotifyRefundment)
	router.Post(endpoints.QueryOrder,   rest.QueryOrder)
	router.Post(endpoints.CloseOrder,   rest.CloseOrder)
	router.Post(endpoints.Transfer,     rest.Transfer)
	router.Post(endpoints.QueryTransfer,rest.QueryTransfer)

	api.UseHandler(router)

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	fmt.Printf("I am listening at %s...\n", listenParam)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}

