// +build gateway notify

package wxpay

type INotifyParams interface {
	parse(map[string]string, error) error
}

type NotifyError struct {
	Err string `json:"err"`
}

// implementation of interface INotifyParams
func (n *NotifyError) parse(_ map[string]string, err error) error {
	n.Err = err.Error()
	return nil
}

type NotifyParams struct {
	AppName string    `json:"app_name"`
	CbUrl   string    `json:"cb_url"`
	INotifyParams     `json:"params"`
}

func _NewNotifyError(err error) *NotifyParams {
	e := &NotifyError{}
	e.parse(nil, err)
	return &NotifyParams{INotifyParams:e}
}

type FnParseNotifyBody func(prompt string, body []byte, appKey string) *NotifyParams

