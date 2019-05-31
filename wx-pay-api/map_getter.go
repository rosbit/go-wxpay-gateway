// +build gateway

package wxpay

import (
	"fmt"
	"strconv"
)

type _M  map[string]string

func (m _M) getString(val *string, name string, must bool) error {
	if v, ok := m[name]; !ok {
		if must {
			return fmt.Errorf("param %s not found", name)
		}
	} else {
		*val = v
	}
	return nil
}

func (m _M) getInt(val *int, name string, must bool) error {
	if v, ok := m[name]; !ok {
		if must {
			return fmt.Errorf("param %s not found", name)
		}
	} else {
		if n, err := strconv.Atoi(v); err != nil {
			return fmt.Errorf("param %s(%s) is not an integer", name, v)
		} else {
			*val = n
		}
	}
	return nil
}

func (m _M) getBool(val *bool, name string, must bool) error {
	if v, ok := m[name]; !ok {
		if must {
			return fmt.Errorf("param %s not found", name)
		}
	} else {
		*val = (v == "Y")
	}
	return nil
}

