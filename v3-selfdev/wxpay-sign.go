package v3sd

import (
	"encoding/base64"
	"encoding/json"
	"crypto/rand"
	"crypto/rsa"
	"crypto"
	"bytes"
	"time"
	"io"
	"os"
	"fmt"
)

// 对请求参数进行签名
func MakeSignature(privateKey *rsa.PrivateKey, httpMethod string, uri string, body interface{}, dumpingStrToSign bool, extra ...string) (timestamp, nonce, bodyStr, signature string, err error) {
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

	if len(extra) >= 2 {
		timestamp = extra[0]
		nonce     = extra[1]
	} else {
		t := time.Now()
		timestamp = fmt.Sprintf("%d", t.Unix())
		nonce, _ = generateNonce()
		// nonce = fmt.Sprintf("%d", t.UnixNano())
	}

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
	if !dumpingStrToSign {
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

// GenerateNonce 生成一个长度为 NonceLength 的随机字符串（只包含大小写字母与数字）
func generateNonce() (string, error) {
	const (
		// NonceSymbols 随机字符串可用字符集
		NonceSymbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		// NonceLength 随机字符串的长度
		NonceLength = 32
	)

	bytes := make([]byte, NonceLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	symbolsByteLength := byte(len(NonceSymbols))
	for i, b := range bytes {
		bytes[i] = NonceSymbols[b%symbolsByteLength]
	}
	return string(bytes), nil
}
