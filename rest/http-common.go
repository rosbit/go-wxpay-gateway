package rest

import (
	"github.com/gernest/alien"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func _PathParam(r *http.Request, n string) string {
    p := alien.GetParams(r)
    return p.Get(n)
}

func _WriteBytesJson(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func _WriteJson(w http.ResponseWriter, code int, data interface{}) {
	j, _ := json.Marshal(data)
	_WriteBytesJson(w, code, j)
}

func _WriteError(w http.ResponseWriter, code int, msg string) {
	_WriteJson(w, code, map[string]interface{}{"code": code, "msg": msg})
}

func _WriteMessage(w http.ResponseWriter, msg string) {
	w.Write([]byte(msg))
}

func _ReadJson(r *http.Request, res interface{}) (int, error) {
	if r.Body == nil {
		return http.StatusBadRequest, fmt.Errorf("bad request")
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// fmt.Printf("body: %s\n", string(b))
	if err = json.Unmarshal(b, res); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

