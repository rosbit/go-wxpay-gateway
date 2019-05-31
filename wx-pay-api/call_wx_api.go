// +build gateway getsandbox

package wxpay

import (
	"github.com/rosbit/go-wget"
	"fmt"
)

func callWxAPI(httpFunc wget.HttpFunc, url string, method string, postData interface{}) ([]byte, error) {
	status, content, _, err := httpFunc(url, method, postData, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("status %d", status)
	}
	return content, nil
}

func _CallWxAPI(url string, method string, postData interface{}) ([]byte, error) {
	return callWxAPI(wget.Wget, url, method, postData)
}

func _CallSecureWxAPI(url string, method string, postData interface{}, certFile, keyFile string) ([]byte, error) {
	req, err := wget.NewHttpsRequestWithCerts(0, certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return callWxAPI(req.Run, url, method, postData)
}
