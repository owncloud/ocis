package xmlenc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des" // nolint: gosec
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/beevik/etree"
)

// CBC implements Decrypter and Encrypter for block ciphers in CBC mode
type CBC struct {
	keySize   int
	algorithm string
	cipher    func([]byte) (cipher.Block, error)
}

// KeySize returns the length of the key required.
func (e CBC) KeySize() int {
	return e.keySize
}

// Algorithm returns the name of the algorithm, as will be found
// in an xenc:EncryptionMethod element.
func (e CBC) Algorithm() string {
	return e.algorithm
}

// Encrypt encrypts plaintext with key, which should be a []byte of length KeySize().
// It returns an xenc:EncryptedData element.
func (e CBC) Encrypt(key interface{}, plaintext []byte, _ []byte) (*etree.Element, error) {
	keyBuf, ok := key.([]byte)
	if !ok {
		return nil, ErrIncorrectKeyType("[]byte")
	}
	if len(keyBuf) != e.keySize {
		return nil, ErrIncorrectKeyLength(e.keySize)
	}

	block, err := e.cipher(keyBuf)
	if err != nil {
		return nil, err
	}

	encryptedDataEl := etree.NewElement("xenc:EncryptedData")
	encryptedDataEl.CreateAttr("xmlns:xenc", "http://www.w3.org/2001/04/xmlenc#")
	{
		randBuf := make([]byte, 16)
		if _, err := RandReader.Read(randBuf); err != nil {
			return nil, err
		}
		encryptedDataEl.CreateAttr("Id", fmt.Sprintf("_%x", randBuf))
	}

	em := encryptedDataEl.CreateElement("xenc:EncryptionMethod")
	em.CreateAttr("Algorithm", e.algorithm)
	em.CreateAttr("xmlns:xenc", "http://www.w3.org/2001/04/xmlenc#")

	plaintext = appendPadding(plaintext, block.BlockSize())

	iv := make([]byte, block.BlockSize(), block.BlockSize()+len(plaintext))
	if _, err := RandReader.Read(iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	ciphertext = append(iv, ciphertext...)

	cd := encryptedDataEl.CreateElement("xenc:CipherData")
	cd.CreateAttr("xmlns:xenc", "http://www.w3.org/2001/04/xmlenc#")
	cd.CreateElement("xenc:CipherValue").SetText(base64.StdEncoding.EncodeToString(ciphertext))
	return encryptedDataEl, nil
}

// Decrypt decrypts an encrypted element with key. If the ciphertext contains an
// EncryptedKey element, then the type of `key` is determined by the registered
// Decryptor for the EncryptedKey element. Otherwise, `key` must be a []byte of
// length KeySize().
func (e CBC) Decrypt(key interface{}, ciphertextEl *etree.Element) ([]byte, error) {
	// If the key is encrypted, decrypt it.
	if encryptedKeyEl := ciphertextEl.FindElement("./KeyInfo/EncryptedKey"); encryptedKeyEl != nil {
		var err error
		key, err = Decrypt(key, encryptedKeyEl)
		if err != nil {
			return nil, err
		}
	}

	keyBuf, ok := key.([]byte)
	if !ok {
		return nil, ErrIncorrectKeyType("[]byte")
	}
	if len(keyBuf) != e.KeySize() {
		return nil, ErrIncorrectKeyLength(e.KeySize())
	}

	block, err := e.cipher(keyBuf)
	if err != nil {
		return nil, err
	}

	ciphertext, err := getCiphertext(ciphertextEl)
	if err != nil {
		return nil, err
	}

	// The IV is prefixed, and at least one whole block must follow it.
	// CryptBlocks panics on input that is not a whole number of blocks.
	if len(ciphertext) < 2*block.BlockSize() {
		return nil, errors.New("ciphertext too short")
	}
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, errors.New("ciphertext is not a whole number of blocks")
	}

	iv := ciphertext[:block.BlockSize()]
	ciphertext = ciphertext[block.BlockSize():]

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext) // decrypt in place

	plaintext, err = stripPadding(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

var (
	// AES128CBC implements AES128-CBC symetric key mode for encryption and decryption
	AES128CBC BlockCipher = CBC{
		keySize:   16,
		algorithm: "http://www.w3.org/2001/04/xmlenc#aes128-cbc",
		cipher:    aes.NewCipher,
	}

	// AES192CBC implements AES192-CBC symetric key mode for encryption and decryption
	AES192CBC BlockCipher = CBC{
		keySize:   24,
		algorithm: "http://www.w3.org/2001/04/xmlenc#aes192-cbc",
		cipher:    aes.NewCipher,
	}

	// AES256CBC implements AES256-CBC symetric key mode for encryption and decryption
	AES256CBC BlockCipher = CBC{
		keySize:   32,
		algorithm: "http://www.w3.org/2001/04/xmlenc#aes256-cbc",
		cipher:    aes.NewCipher,
	}

	// TripleDES implements DES in CBC mode for encryption and decryption. Despite
	// its name and algorithm URI it takes an 8 byte single DES key, not a 24 byte
	// 3DES one, so it is not registered as a decrypter: a document naming
	// tripledes-cbc gets ErrAlgorithmNotImplemented instead of weak decryption.
	TripleDES BlockCipher = CBC{
		keySize:   8,
		algorithm: "http://www.w3.org/2001/04/xmlenc#tripledes-cbc",
		cipher:    des.NewCipher,
	}
)

func init() {
	RegisterDecrypter(AES128CBC)
	RegisterDecrypter(AES192CBC)
	RegisterDecrypter(AES256CBC)
}

func appendPadding(buf []byte, blockSize int) []byte {
	paddingBytes := blockSize - (len(buf) % blockSize)
	padding := make([]byte, paddingBytes)
	padding[len(padding)-1] = byte(paddingBytes) //nolint:gosec // paddingBytes <= blockSize, which is at most 16 for AES
	return append(buf, padding...)
}

func stripPadding(buf []byte) ([]byte, error) {
	if len(buf) < 1 {
		return nil, errors.New("buffer is too short for padding")
	}
	paddingBytes := int(buf[len(buf)-1])
	if paddingBytes > len(buf)-1 {
		return nil, errors.New("buffer is too short for padding")
	}
	if paddingBytes < 1 {
		return nil, errors.New("padding must be at least one byte")
	}
	return buf[:len(buf)-paddingBytes], nil
}
