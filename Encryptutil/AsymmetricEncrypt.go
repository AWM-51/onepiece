package Encryptutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// AsymmetricEncrypt 非对称加密，使用RSA算法
func AsymmetricEncrypt(plainText []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	// 使用公钥进行加密
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plainText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with RSA: %v", err)
	}

	// 返回加密后的数据
	return cipherText, nil
}

// AsymmetricDecrypt 非对称解密，使用RSA算法
func AsymmetricDecrypt(cipherText []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// 使用私钥进行解密
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt with RSA: %v", err)
	}

	// 返回解密后的明文数据
	return plainText, nil
}
