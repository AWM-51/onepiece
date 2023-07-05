package Encryptutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// SymmetricEncrypt 对称加密，使用AES算法
func SymmetricEncrypt(plainText []byte, key []byte) ([]byte, error) {
	// 生成AES块密码
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 生成随机的初始化向量
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate AES IV: %v", err)
	}

	// 使用AES块密码和初始化向量创建加密器
	stream := cipher.NewCFBEncrypter(block, iv)

	// 加密明文数据
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)

	// 返回加密后的数据，包括初始化向量
	return append(iv, cipherText...), nil
}

// SymmetricDecrypt 对称解密，使用AES算法
func SymmetricDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	// 从密文中提取初始化向量和加密数据
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// 生成AES块密码
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 使用AES块密码和初始化向量创建解密器
	stream := cipher.NewCFBDecrypter(block, iv)

	// 解密密文数据
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)

	// 返回解密后的明文数据
	return plainText, nil
}

// Base64Encode 使用Base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 使用Base64解码
func Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
