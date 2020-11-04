package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
)

func RsaEncryptNoPadding(data, key []byte) (string, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(key)
	if block == nil {
		return "", errors.New("key error")
	}
	// 解析公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// 类型断言
	p, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("publicKey error")
	}
	biText := new(big.Int).SetBytes(data)
	biE := big.NewInt(int64(p.E))
	biN := p.N
	c := new(big.Int)
	exp := c.Exp(biText, biE, biN)
	//padding:0
	biRet := fmt.Sprintf("%x", exp)
	for len(biRet) < 256 {
		biRet = "0" + biRet
	}
	return biRet, nil
}
