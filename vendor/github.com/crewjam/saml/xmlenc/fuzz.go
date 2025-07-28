package xmlenc

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/beevik/etree"
)

var testKey = func() *rsa.PrivateKey {
	// Generate a new test key instead of using hardcoded private key
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	return key
}()

// Fuzz is the go-fuzz fuzzing function
func Fuzz(data []byte) int {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(data); err != nil {
		return 0
	}
	if doc.Root() == nil {
		return 0
	}

	if _, err := Decrypt(testKey, doc.Root()); err != nil {
		return 0
	}
	return 1
}
