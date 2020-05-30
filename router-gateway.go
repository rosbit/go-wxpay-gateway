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
)

func StartService() error {
	serviceConf := &conf.ServiceConf

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
	router.Post(paymentEndpoint,        rest.CreatePayment)
	router.Post(endpoints.CreateRefund, rest.CreateRefundment)
	router.Post(endpoints.QueryOrder,   rest.QueryOrder)
	router.Post(endpoints.CloseOrder,   rest.CloseOrder)
	router.Post(endpoints.Transfer,     rest.Transfer)
	router.Post(endpoints.QueryTransfer,rest.QueryTransfer)
	router.Post(endpoints.VerifyNotifyPay,    rest.VerifyNotifyPayment)
	router.Post(endpoints.VerifyNotifyRefund, rest.VerifyNotifyRefundment)
	if len(endpoints.RealnameAuthRoot) > 0 {
		realnameAuthEndpoint := appendEndpointParam(endpoints.RealnameAuthRoot, rest.REALNAME_TYPE_NAME)
		router.Post(realnameAuthEndpoint, rest.AuthRealname)
	}
	if len(endpoints.HealthCheck) > 0 {
		router.Get(endpoints.HealthCheck, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "OK\n")
		})
	}

	api.UseHandler(router)

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	fmt.Printf("I am listening at %s...\n", listenParam)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}

