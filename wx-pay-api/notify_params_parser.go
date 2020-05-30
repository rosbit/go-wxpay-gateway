// +build gateway notify

package wxpay

type INotifyParams interface {
	parse(map[string]string, error) error
}

type FnParseNotifyBody func(prompt string, body []byte, apiKey string) (INotifyParams, error)

