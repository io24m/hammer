package hammer

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
)

const (
	//sKey        = "dde4b1f8a9e6b814"
	presetKey   = "0CoJUm6Qyw8W8jud"
	ivParameter = "0102030405060708"
	base62      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	publicKey   = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDgtQn2JZ34ZC28NWYpAUd98iZ37BUrX/aKzmFbt7clFSs6sXqHauqKWqdtLkF2KexO40H1YTX8z2lSgBBOAxLsvaklV8k4cBFK9snQXE9/DDaFt6Rr7iVZMldczhC0JNgTz+SHXT6CBHuX3e9SdB1Ua44oncaTWz7OBGLbCiK45wIDAQAB\n-----END PUBLIC KEY-----"
)

func key(len int) (res []byte) {
	res = make([]byte, 0)
	for i := 0; i < len; i++ {
		res = append(res, base62[r.Intn(63)])
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
func rsaEncrypt(data, key []byte) string {
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
	//加密
	v15, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	toString := hex.EncodeToString(v15)
	return toString
}

func weapiEncrypt(data interface{}) (res map[string]interface{}) {
	res = make(map[string]interface{})
	jsonStr, _ := json.Marshal(data)
	secretKey := key(16)
	rKey := reverseKey(secretKey)
	encrypt := AES_CBC_Encrypt(jsonStr, []byte(presetKey), []byte(ivParameter))
	b64 := base64.StdEncoding.EncodeToString(encrypt)
	aes128Encrypt := AES_CBC_Encrypt([]byte(b64), secretKey, []byte(ivParameter))
	b64 = base64.StdEncoding.EncodeToString(aes128Encrypt)
	res["params"] = b64
	res["encSecKey"] = rsaEncrypt([]byte(rKey), []byte(publicKey))
	return
}

func linuxapiEncrypt(data interface{}) (res map[string]interface{}) {
	return
}

func eapiEncrypt(url string, data interface{}) (res map[string]interface{}) {
	return
}

func decrypt(data interface{}) interface{} {
	return nil
}

////加密
//func PswEncrypt(src string) string {
//	key := []byte(sKey)
//	iv := []byte(ivParameter)
//	result, err := Aes128Encrypt([]byte(src), key, iv)
//	if err != nil {
//		panic(err)
//	}
//	return base64.RawStdEncoding.EncodeToString(result)
//}
//
////解密
//func PswDecrypt(src string) string {
//	key := []byte(sKey)
//	iv := []byte(ivParameter)
//	var result []byte
//	var err error
//	result, err = base64.RawStdEncoding.DecodeString(src)
//	if err != nil {
//		panic(err)
//	}
//	origData, err := Aes128Decrypt(result, key, iv)
//	if err != nil {
//		panic(err)
//	}
//	return string(origData)
//
//}

func Aes128Encrypt(origData, key, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, IV[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func Aes128Decrypt(crypted, key []byte, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, IV[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//对明文进行填充
func Padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

//对密文删除填充
func UnPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

//AEC加密（CBC模式）
func AES_CBC_Encrypt(plainText, key, iv []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行填充
	plainText = Padding(plainText, block.BlockSize())
	//指定初始向量vi,长度和block的块尺寸一致
	//iv := []byte("12345678abcdefgh")
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密连续数据库
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText
}

//AEC解密（CBC模式）
func AES_CBC_Decrypt(cipherText, key, iv []byte) []byte {
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//指定初始化向量IV,和加密的一致
	//iv := []byte("12345678abcdefgh")
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	//删除填充
	plainText = UnPadding(plainText)
	return plainText
}
