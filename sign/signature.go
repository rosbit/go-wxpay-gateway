package sign

import (
	"hash"
	"crypto/md5"
	"crypto/sha256"
	"crypto/hmac"
	"sort"
	"fmt"
	"io"
)

const (
	MD5         = "MD5"
	HMAC_SHA256 = "HMAC-SHA256"
)

// 用不同的签名方法实现微信相关的签名方法
func CreateSignature(signType string, params map[string]string, key string) string {
	var h hash.Hash
	switch signType {
	case HMAC_SHA256:
		h = hmac.New(sha256.New, []byte(key))
	default:
		h = md5.New()
	}

	return createSignature(h, params, key)
}

// 校验签名，已有的签名值来自signField，缺省为"sign"
func CheckSignature(signType string, res map[string]string, key string, signField ...string) error {
	signFieldName := func()string {
		if len(signField) > 0 {
			return signField[0]
		}
		return "sign"
	}()

	if signature, ok := res[signFieldName]; !ok {
		return fmt.Errorf("no signature in result")
	} else {
		delete(res, signFieldName)
		createdSign := CreateSignature(signType, res, key)
		if signature != createdSign {
			return fmt.Errorf("signature not matched: %s != %s", signature, createdSign)
		}
	}
	return nil
}

func createSignature(h hash.Hash, params map[string]string, key string) string {
	// sort keys
	keys := make([]string, len(params))
	i := 0
	for k := range params {
		keys[i] = k
		i += 1
	}
	sort.Strings(keys)

	// create stringA="key1=val1&"
	stringA := h
	first := true // 是否第一个非空值参数
	for _, key := range keys {
		v, _ := params[key]
		if v != "" {
			if first {
				first = false
			} else {
				io.WriteString(stringA, "&")
			}
			io.WriteString(stringA, key)
			io.WriteString(stringA, "=")
			io.WriteString(stringA, v)
		}
	}

	// stringSignTemp=stringA+'&key='+apiKey
	io.WriteString(stringA, "&key=")
	io.WriteString(stringA, key)

	// sign(stringSignTemp) -> Hex -> toUpper
	return fmt.Sprintf("%X", stringA.Sum(nil))
}

