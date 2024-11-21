// +build gateway getsandbox

package wxpay

import (
	"github.com/rosbit/gnet"
	"os"
	"fmt"
)

func parseHttpRes(status int, content []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("status %d", status)
	}
	return content, nil
}

func callWxAPI(httpFunc gnet.HttpFunc, url string, method string, postData interface{}) ([]byte, error) {
	status, content, _, err := httpFunc(url, gnet.M(method), gnet.Params(postData), gnet.BodyLogger(os.Stderr))
	return parseHttpRes(status, content, err)
}

func _CallWxAPI(url string, method string, postData interface{}) ([]byte, error) {
	return callWxAPI(gnet.Http, url, method, postData)
}

func _CallSecureWxAPI(url string, method string, postData interface{}, certFile, keyFile string) ([]byte, error) {
	req, err := gnet.NewHttpsRequestWithCerts(certFile, keyFile, gnet.BodyLogger(os.Stderr))
	if err != nil {
		return nil, err
	}
	status, content, _, err := req.Http(url, method, postData, nil)
	return parseHttpRes(status, content, err)
}
