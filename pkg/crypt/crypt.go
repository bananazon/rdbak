package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

const (
	hexKey string = "2a89f4811dae393b45e9d902388783b39a15ddd63a2cbf19b90f6a3e3dd9a06b"
)

func Encrypt(msg string) (ciphertext string, err error) {
	key, _ := hex.DecodeString(hexKey)
	plaintext := []byte(msg)

	block, err := aes.NewCipher(key)
	if err != nil {
		return ciphertext, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return ciphertext, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return ciphertext, err
	}

	bytes := aesGCM.Seal(nonce, nonce, plaintext, nil)
	ciphertext = fmt.Sprintf("%x", bytes)
	return ciphertext, nil
}

func Decrypt(encrypted string) (plaintext string, err error) {
	key, _ := hex.DecodeString(hexKey)
	enc, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherBytes := enc[:nonceSize], enc[nonceSize:]

	bytes, err := aesGCM.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return plaintext, err
	}

	plaintext = string(bytes)
	return plaintext, nil
}
