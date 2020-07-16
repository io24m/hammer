package hammer

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
)

const (
	base62      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	presetKey   = "0CoJUm6Qyw8W8jud"
	linuxapiKey = "rFgB&h#%2?^eDg:Q"
	ivParameter = "0102030405060708"
	publicKey   = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDgtQn2JZ34ZC28NWYpAUd98iZ37BUrX/aKzmFbt7clFSs6sXqHauqKWqdtLkF2KexO40H1YTX8z2lSgBBOAxLsvaklV8k4cBFK9snQXE9/DDaFt6Rr7iVZMldczhC0JNgTz+SHXT6CBHuX3e9SdB1Ua44oncaTWz7OBGLbCiK45wIDAQAB\n-----END PUBLIC KEY-----"
)

func key(len int) (res []byte) {
	res = make([]byte, 0)
	for i := 0; i < len; i++ {
		res = append(res, base62[r.Intn(61)])
	}
	return
}

func reverseKey(key []byte) []byte {
	a := make([]byte, len(key), len(key))
	copy(a, key)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func rsaNoPaddingEncrypt(data, key []byte) string {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return ""
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return ""
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	biText := new(big.Int).SetBytes(data)
	biE := big.NewInt(int64(pub.E))
	biN := pub.N
	c := new(big.Int)
	exp := c.Exp(biText, biE, biN)
	//padding:0
	biRet := fmt.Sprintf("%x", exp)
	for len(biRet) < 256 {
		biRet = "0" + biRet
	}
	return biRet
}

func weapiEncrypt(data interface{}) (res map[string]interface{}) {
	res = make(map[string]interface{})
	jsonStr, _ := json.Marshal(data)
	secretKey := key(16)
	rKey := reverseKey(secretKey)
	encrypt := aesCbcEncrypt(jsonStr, []byte(presetKey), []byte(ivParameter))
	b64 := base64Encode(encrypt)
	aes128Encrypt := aesCbcEncrypt([]byte(b64), secretKey, []byte(ivParameter))
	b64 = base64Encode(aes128Encrypt)
	res["params"] = string(b64)
	res["encSecKey"] = rsaNoPaddingEncrypt(rKey, []byte(publicKey))
	return
}

func linuxapiEncrypt(data interface{}) (res map[string]interface{}) {
	res = make(map[string]interface{})
	jsondata, _ := json.Marshal(data)
	ecb := aesEncryptECB(jsondata, []byte(linuxapiKey))
	res["eparams"] = hex.EncodeToString(ecb)
	return
}

func eapiEncrypt(url string, data interface{}) (res map[string]interface{}) {
	return
}

func decrypt(data interface{}) interface{} {
	return nil
}

func aesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}
func aesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

//对明文进行填充
func padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

//对密文删除填充
func unPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

//AEC加密（CBC模式）
func aesCbcEncrypt(plainText, key, iv []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行填充
	plainText = padding(plainText, block.BlockSize())
	//指定初始向量vi,长度和block的块尺寸一致
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密连续数据库
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText
}

//AEC解密（CBC模式）
func aesCbcDecrypt(cipherText, key, iv []byte) []byte {
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//指定初始化向量IV,和加密的一致
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	//删除填充
	plainText = unPadding(plainText)
	return plainText
}

func base64Encode(data []byte) (buf []byte) {
	stdEncoding := base64.StdEncoding
	buf = make([]byte, stdEncoding.EncodedLen(len(data)))
	stdEncoding.Encode(buf, data)
	return
}
