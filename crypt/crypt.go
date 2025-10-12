package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/gdanko/rdbak/globals"
)

func Encrypt(msg string) string {
	key, _ := hex.DecodeString(globals.GetHexKey())
	plaintext := []byte(msg)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(111)
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(222)
		panic(err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(333)
		panic(err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func Decrypt(encrypted string, keyHex string) string {
	key, _ := hex.DecodeString(globals.GetHexKey())
	enc, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s", plaintext)
}
