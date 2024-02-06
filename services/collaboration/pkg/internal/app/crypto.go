package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func keyPadding(key []byte) []byte {
	switch length := len(key); {
	case length < 16:
		return append(key, make([]byte, 16-length)...)
	case length == 16:
		return key
	case length < 24:
		return append(key, make([]byte, 24-length)...)
	case length == 24:
		return key
	case length < 32:
		return append(key, make([]byte, 32-length)...)
	case length == 32:
		return key
	case length > 32:
		return key[:32]
	}
	return []byte{}
}

func EncryptAES(key []byte, plainText string) (string, error) {
	src := []byte(plainText)

	block, err := aes.NewCipher(keyPadding(key))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(src))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], src)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(key []byte, securemess string) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyPadding(key))
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("ciphertext block size is too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
