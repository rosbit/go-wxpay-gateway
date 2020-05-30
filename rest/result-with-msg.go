package rest

import (
	"net/http"
)

func sendResultWithMsg(debug bool, w http.ResponseWriter, sent, recv []byte, err error, extra ...map[string]interface{}) {
	status, msg := func() (int, string) {
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusOK, "OK"
	}()

	res := map[string]interface{}{
		"code": status,
		"msg":  msg,
	}
	if debug {
		res["sent"], res["recv"] = string(sent), string(recv)
	}
	if len(extra) > 0 {
		for k, v := range extra[0] {
			res[k] = v
		}
	}
	_WriteJson(w, status, res)
}
