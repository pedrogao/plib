package main

/*
	reference: https://mojotv.cn/go/golang-crypto
*/

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {
	assert := assert.New(t)

	orig := "hello world"
	key := "123456781234567812345678"
	t.Logf("原文：%s", orig)

	encryptCode := aesEncrypt(orig, key)
	t.Logf("密文：%s", encryptCode)

	decryptCode := aesDecrypt(encryptCode, key)
	t.Logf("解密结果：%s", decryptCode)

	assert.Equal(orig, decryptCode)
}

/*
 AES 加密
 @param org original string
 @param key encrypt key
*/
func aesEncrypt(org string, key string) string {
	orgData := []byte(org)
	k := []byte(key)

	// 分组密钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}
	// 密钥块长度
	blockSize := block.BlockSize()
	// 补全码
	orgData = PKCS7Padding(orgData, blockSize)
	// 加密模式
	// The length of iv must be the same as the Block's block size.
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(orgData))
	// 加密
	blockMode.CryptBlocks(cryted, orgData)
	// base64编码，for 可读性
	// 使用RawURLEncoding 不要使用StdEncoding
	// 不要使用StdEncoding  放在url参数中回导致错误
	return base64.RawURLEncoding.EncodeToString(cryted)
}

/*
 AES解密
 @param crypted string
 @param key
*/
func aesDecrypt(crypted string, key string) string {
	// base64解码
	cryptedBytes, err := base64.RawURLEncoding.DecodeString(crypted)
	if err != nil {
		panic(err)
	}
	k := []byte(key)
	// 分组密钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 原始数据
	orgData := make([]byte, len(cryptedBytes))
	// 解密
	blockMode.CryptBlocks(orgData, cryptedBytes)
	// 除去补全码
	org := PKCS7UnPadding(orgData)
	return string(org)
}

/*
 补码，将加密字符串补全至blockSize的倍数
*/
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

/*
 去补码，得到补全的个数，然后去掉
*/
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
