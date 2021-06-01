package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// DecryptAES256GCM 使用 AEAD_AES_256_GCM 算法进行解密
func DecryptAES256GCM(aesKey, ciphertext string, nonce string, additionalData string) ([]byte, error) {
	s, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, []byte(nonce), s, []byte(additionalData))
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
