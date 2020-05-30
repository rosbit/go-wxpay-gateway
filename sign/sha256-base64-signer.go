// +build gateway

package sign

import (
	"crypto"
	"crypto/x509"
	"crypto/rand"
	"crypto/sha256"
	"crypto/rsa"
	"encoding/pem"
	"encoding/base64"
	"io/ioutil"
	"fmt"
)

type Sha256Base64Signer struct {
	priv *rsa.PrivateKey
}

func CreateSha256Base64Signer(keyFile string) (*Sha256Base64Signer, error) {
	b, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, fmt.Errorf("no block read from keyFile")
	}
	p, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := p.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unknown private key type")
	}
	return &Sha256Base64Signer{priv}, nil
}

func (s *Sha256Base64Signer) Sign(plainText []byte) (cipherText string, err error) {
	hashed := sha256.Sum256(plainText)
	var signature []byte
	if signature, err = s.priv.Sign(rand.Reader, hashed[:], crypto.SHA256); err != nil {
		return
	}
	cipherText = base64.StdEncoding.EncodeToString(signature)
	return
}

