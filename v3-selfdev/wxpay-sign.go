package v3sd

import (
	"encoding/base64"
	"encoding/json"
	"crypto/rsa"
	"crypto"
	"bytes"
	"time"
	"io"
	"os"
	"fmt"
)

// 对请求参数进行签名
func MakeSignature(privateKey *rsa.PrivateKey, httpMethod string, uri string, body interface{}, dumpingStrToSign ...bool) (timestamp, nonce, bodyStr, signature string, err error) {
	b := &bytes.Buffer{}
	if body != nil {
		j := json.NewEncoder(b)
		j.SetEscapeHTML(false)
		err = j.Encode(body)
		if err != nil {
			return
		}
		b.Truncate(b.Len()-1) // 去掉最后一个换行
	}
	bodyStr = b.String()
	t := time.Now()
	timestamp = fmt.Sprintf("%d", t.Unix())
	nonce = fmt.Sprintf("%d", t.UnixNano())

	in := make(chan string)
	go func() {
		in <- httpMethod
		in <- "\n"
		in <- uri
		in <- "\n"
		in <- timestamp
		in <- "\n"
		in <- nonce
		in <- "\n"
		in <- bodyStr
		in <- "\n"
		close(in)
	}()

	h := crypto.Hash.New(crypto.SHA256)
	if len(dumpingStrToSign) == 0 || !dumpingStrToSign[0] {
		makeDataToSign(in, h)
	} else {
		mw, deferFunc := dumpStrToSign(h)
		makeDataToSign(in, mw)
		deferFunc()
	}

	hashStr := h.Sum(nil)
	sign, e := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashStr)
	if e != nil {
		err = e
		return
	}
	signature = base64.StdEncoding.EncodeToString(sign)
	return
}

func dumpStrToSign(w io.Writer) (mw io.Writer, deferFunc func()) {
	fmt.Fprintf(os.Stderr, "----- strToSign begin -----\n")
	deferFunc = func() {
		fmt.Fprintf(os.Stderr, "\n----- strToSign end -----\n")
	}
	mw = io.MultiWriter(w, os.Stderr)
	return
}

func makeDataToSign(in <-chan string, w io.Writer) {
	for s := range in {
		io.WriteString(w, s)
	}
}

