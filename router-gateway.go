// +build gateway

/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/rosbit/mgin"
	"go-wxpay-gateway/conf"
	"go-wxpay-gateway/rest"
	"net/http"
	"os"
	"fmt"
)

func StartService() error {
	serviceConf := &conf.ServiceConf

	api := mgin.NewMgin(mgin.WithLogger("wxpay-gateway"), mgin.CreateBodyDumpingHandler(os.Stderr, "raw body"))
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
	api.POST(paymentEndpoint,        rest.CreatePayment)
	api.POST(endpoints.CreateRefund, rest.CreateRefundment)
	api.POST(endpoints.QueryOrder,   rest.QueryOrder)
	api.POST(endpoints.CloseOrder,   rest.CloseOrder)
	api.POST(endpoints.Transfer,     rest.Transfer)
	api.POST(endpoints.QueryTransfer,rest.QueryTransfer)
	api.POST(endpoints.V3Transfer,           rest.V3Transfer)
	api.POST(endpoints.V3QueryTransfer,      rest.V3QueryTransfer)
	api.POST(endpoints.V3QueryTransferDetail,rest.V3QueryTransferDetail)
	api.POST(endpoints.VerifyNotifyPay,    rest.VerifyNotifyPayment)
	api.POST(endpoints.VerifyNotifyRefund, rest.VerifyNotifyRefundment)
	if len(endpoints.RealnameAuthRoot) > 0 {
		realnameAuthEndpoint := appendEndpointParam(endpoints.RealnameAuthRoot, rest.REALNAME_TYPE_NAME)
		api.POST(realnameAuthEndpoint, rest.AuthRealname)
	}
	if len(endpoints.HealthCheck) > 0 {
		api.GET(endpoints.HealthCheck, func(c *mgin.Context) {
			c.String(http.StatusOK, "OK\n")
		})
	}

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	fmt.Printf("I am listening at %s...\n", listenParam)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}

