package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// keyPadding will add the required zero padding to the provided key.
// The resulting key will have a length of either 16, 24 or 32 bytes.
// If the key has more than 32 bytes, only the first 32 bytes will be returned.
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

// EncryptAES encrypts the provided plainText using the provided key.
// AES CFB will be used as cryptographic method.
// Use DecryptAES to decrypt the resulting string
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

// DecryptAES decrypts the provided string using the provided key.
// The provided string must have been encrypted with AES CFB.
// This method will decrypt the result from the EncryptAES method
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
