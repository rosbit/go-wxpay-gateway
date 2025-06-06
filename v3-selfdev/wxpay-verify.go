// 微信支付验签

package v3sd

import (
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"go-wxpay-gateway/conf"
	"encoding/base64"
	"encoding/json"
	"crypto/sha256"
	"crypto/rsa"
	"crypto"
	"time"
	"math"
	"fmt"
	"io"
	"os"
	"bytes"
	"strconv"
	"net/http"
)

func VerifyTransferBills(appName string, r *http.Request) (code int, respBody json.RawMessage, billBody json.RawMessage, err error) {
	defer func() {
		if code == http.StatusOK {
			respBody = json.RawMessage(ok_msg)
		} else {
			respBody = json.RawMessage(fail_msg)
		}
	}()
	defer r.Body.Close()

	code = http.StatusNotImplemented
	mchConf, ok := conf.GetAppAttrs(appName)
	if !ok {
		err = fmt.Errorf("conf not found for %s", appName)
		return
	}

	if len(mchConf.WxpayV3Cert) > 0 {
		// 平台证书模式
		code, billBody, err = verifyReqWithPlatformMode(appName, mchConf, r)
	} else {
		// 微信支付公钥模式
		code, billBody, err = verfiyReqWithPublicKeyMode(appName, mchConf, r)
	}
	return
}

// 获取头信息
const (
	HeaderNonce     = "Wechatpay-Nonce"
	HeaderSignature = "Wechatpay-Signature"
	HeaderTimestamp = "Wechatpay-Timestamp"
	HeaderSerial    = "Wechatpay-Serial"
)
func getVerifyHeader(r *http.Request) (serialNo, nonce, timestamp, signature string){
	h := r.Header
	serialNo = h.Get(HeaderSerial)
	nonce    = h.Get(HeaderNonce)
	timestamp= h.Get(HeaderTimestamp)
	signature= h.Get(HeaderSignature)
	fmt.Fprintf(os.Stderr, "Wechatpay-Serial: %s\nWechatpay-Nonce: %s\nWechatpay-Timestamp: %s\nWechatpay-Signature:%s\n", serialNo, nonce, timestamp, signature)
	return
}

var (
	ok_msg   = []byte(`{"code":"SUCCESS"}`)
	fail_msg = []byte(`{"code":"SUCCESS","message":"FAILED"}`)
)

func verifyAndDecryptBody(pubKey *rsa.PublicKey, mchConf *conf.MerchantConf, r *http.Request) (code int, billBody json.RawMessage, err error) {
	code = http.StatusNotAcceptable

	b := &bytes.Buffer{}
	tr := io.TeeReader(r.Body, b)

	serialNo, nonce, timestamp, signature := getVerifyHeader(r)
	if serialNo != mchConf.WxpayPubkeyId {
		err = fmt.Errorf("seaialNo %s not matched", serialNo)
		return
	}
	if err = checkTimestamp(timestamp); err != nil {
		return
	}
	sigBytes, e := base64.StdEncoding.DecodeString(signature)
	if e != nil {
		err = e
		return
	}

	h := &bytes.Buffer{}
	// 应答时间戳\n应答随机串\n应答报文主体\n
	fmt.Fprintf(h, "%s\n%s\n", timestamp, nonce)
	io.Copy(h, tr) // now b with content of r.Body
	fmt.Fprintf(h, "\n")

	message := h.Bytes()
	fmt.Fprintf(os.Stderr, "message: %s\n", message)
	hashed := sha256.Sum256(message)

	if err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sigBytes); err != nil {
		fmt.Fprintf(os.Stderr, "failed to rsa.VerifyPKCS1v15: %v\n", err)
		return
	}

	// 验签通过
	code = http.StatusOK

	// AEAD_AES_256_GCM解密报文
	body := b.Bytes()
	fmt.Fprintf(os.Stderr, "body: %s\n", body)
	var vb VerifyBody
	if err = json.Unmarshal(body, &vb); err != nil {
		fmt.Fprintf(os.Stderr, "failed to json.Unmarshal: %v\n", err)
		return
	}
	resource := &vb.Resource
	plainText, e := utils.DecryptAES256GCM(mchConf.WxpayV3Key, resource.AssociatedData, resource.Nonce, resource.Ciphertext)
	if e != nil {
		fmt.Fprintf(os.Stderr, "failed to utils.DecryptAES256GCM: %v\n", e)
		err = e
		return
	}
	billBody = json.RawMessage([]byte(plainText))
	return
}

// 微信支付公钥模式验签[新模式]
func verfiyReqWithPublicKeyMode(appName string, mchConf *conf.MerchantConf, r *http.Request) (code int, billBody json.RawMessage, err error) {
	code = http.StatusNotImplemented
	if len(mchConf.MchCertSerialNo) == 0 {
		err = fmt.Errorf("mch-cert-serialno specified for %s", appName)
		return
	}
	if len(mchConf.WxpayPubkeyId) == 0 {
		err = fmt.Errorf("wxpay-pubkey-id specified for %s", appName)
		return
	}
	if len(mchConf.WxpayPubkeyFile) == 0 {
		err = fmt.Errorf("wxpay-pubkey-file not specified for %s", appName)
		return
	}
	wechatpayPublicKey, e := utils.LoadPublicKeyWithPath(mchConf.WxpayPubkeyFile)
	if e != nil {
		err = e
		return
	}

	// 用微信支付公钥对验签名串和签名进行 SHA256 with RSA 签名验证
	return verifyAndDecryptBody(wechatpayPublicKey, mchConf, r)
}

// 平台证书模式验签[旧模式]
func verifyReqWithPlatformMode(appName string, mchConf *conf.MerchantConf, r *http.Request) (code int, billBody json.RawMessage, err error) {
	code = http.StatusNotImplemented
	if len(mchConf.MchCertSerialNo) == 0 {
		err = fmt.Errorf("mch-cert-serialno not specified for %s", appName)
		return
	}

	mchPublicKey, e := utils.LoadPublicKey(mchConf.MchCertPemFile)
	if e != nil {
		err = e
		return
	}

	// 所有应答，微信支付都会使用平台证书私钥签名，商户需要使用平台证书公钥验证签名
	return verifyAndDecryptBody(mchPublicKey, mchConf, r)
}

const max_timestamp_gap = 120.0 // 时间戳最大间隔，2分钟
func checkTimestamp(timestamp string) (err error) {
	ts, e := strconv.ParseInt(timestamp, 10, 64)
	if e != nil {
		err = e
		return
	}
	now := time.Now().Unix()
	gap := math.Abs(float64(now - ts))
	if gap < max_timestamp_gap {
		return
	}
	return fmt.Errorf("timestamp %s is too old or too new", timestamp)
}
