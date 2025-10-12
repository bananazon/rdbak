package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/gdanko/rdbak/pkg/globals"
)

func Encrypt(msg string) (ciphertext string, err error) {
	key, _ := hex.DecodeString(globals.GetHexKey())
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

func Decrypt(encrypted string, keyHex string) (plaintext string, err error) {
	key, _ := hex.DecodeString(globals.GetHexKey())
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
