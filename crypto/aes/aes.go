package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

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

//AesEncryptECB Aes(ECB)
func AesEncryptECB(origData, key []byte) (encrypted []byte, err error) {
	c, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return nil, err
	}
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, c.BlockSize(); bs <= len(origData); bs, be = bs+c.BlockSize(), be+c.BlockSize() {
		c.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted, nil
}

//AesDecryptECB Aes(ECB)
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte, err error) {
	c, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return nil, err
	}
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, c.BlockSize(); bs < len(encrypted); bs, be = bs+c.BlockSize(), be+c.BlockSize() {
		c.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim], nil
}

//padding
func padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

//unPadding
func unPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

//AesEncryptCBC Aes(CBC)
func AesEncryptCBC(plainText, key, iv []byte) ([]byte, error) {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
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
	return cipherText, nil
}

//AesDecryptCBC Aes(CBC)
func AesDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//指定初始化向量IV,和加密的一致
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	//删除填充
	plainText = unPadding(plainText)
	return plainText, nil
}
