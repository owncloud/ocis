package xmlenc

import (
	"crypto/rsa"

	"github.com/beevik/etree"
)

// Fuzz is the go-fuzz fuzzing function
func Fuzz(testKey *rsa.PrivateKey, data []byte) int {
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
