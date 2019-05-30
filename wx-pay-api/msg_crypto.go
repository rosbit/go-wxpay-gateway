/**
 * 微信信息加解密/签名/随机数
 * 1. HashStrings([]string)    -- 微信常用的签名算法 []string -> sort -> sha1 -> hex
 * 2. GetRandomBytes(int)      -- 获取指定长度的随机串，随机字符为 数字/小写字母/大写字母
 */
package wxpay

import (
	"sort"
	"fmt"
	"crypto/sha1"
	"math/rand"
	"time"
	"io"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func _HashStrings(sl []string) string {
	sort.Strings(sl)
	h := sha1.New()
	for _, s := range sl {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

var _rule = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func _GetRandomBytes(n int) []byte {
	b := make([]byte, n)
	rc := len(_rule)
	for i:=0; i<n; i++ {
		b[i] = _rule[rand.Intn(rc)]
	}
	return b
}

