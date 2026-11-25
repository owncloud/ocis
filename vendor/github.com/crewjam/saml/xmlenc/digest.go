package xmlenc

import (
	"crypto/sha1" //nolint:gosec // required for protocol support
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	//nolint:staticcheck,gosec // We should support this for legacy reasons.
	"golang.org/x/crypto/ripemd160"
)

type digestMethod struct {
	algorithm string
	hash      func() hash.Hash
}

func (dm digestMethod) Algorithm() string {
	return dm.algorithm
}

func (dm digestMethod) Hash() hash.Hash {
	return dm.hash()
}

var (
	// SHA1 implements the SHA-1 digest method (which is considered insecure)
	SHA1 = digestMethod{
		algorithm: "http://www.w3.org/2000/09/xmldsig#sha1",
		hash:      sha1.New,
	}

	// SHA256 implements the SHA-256 digest method. The algorithm URI is the one
	// assigned by RFC 4051 section 2.3.2, which lives in the XML Encryption
	// namespace rather than the XML Signature one.
	SHA256 = digestMethod{
		algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
		hash:      sha256.New,
	}

	// SHA512 implements the SHA-512 digest method (RFC 4051 section 2.3.4).
	SHA512 = digestMethod{
		algorithm: "http://www.w3.org/2001/04/xmlenc#sha512",
		hash:      sha512.New,
	}

	// RIPEMD160 implements the RIPEMD160 digest method (RFC 4051 section 2.3.5).
	RIPEMD160 = digestMethod{
		algorithm: "http://www.w3.org/2001/04/xmlenc#ripemd160",
		hash:      ripemd160.New,
	}

	// The URIs below were emitted by earlier releases of this package, which
	// mistakenly placed these digests in the XML Signature namespace. They are
	// registered so that documents produced by those releases keep decrypting.
	legacySHA256 = digestMethod{
		algorithm: "http://www.w3.org/2000/09/xmldsig#sha256",
		hash:      sha256.New,
	}

	legacySHA512 = digestMethod{
		algorithm: "http://www.w3.org/2000/09/xmldsig#sha512",
		hash:      sha512.New,
	}

	legacyRIPEMD160 = digestMethod{
		algorithm: "http://www.w3.org/2000/09/xmldsig#ripemd160",
		hash:      ripemd160.New,
	}
)

func init() {
	RegisterDigestMethod(SHA1)
	RegisterDigestMethod(SHA256)
	RegisterDigestMethod(SHA512)
	RegisterDigestMethod(RIPEMD160)
	RegisterDigestMethod(legacySHA256)
	RegisterDigestMethod(legacySHA512)
	RegisterDigestMethod(legacyRIPEMD160)
}
