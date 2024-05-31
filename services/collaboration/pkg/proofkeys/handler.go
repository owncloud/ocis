package proofkeys

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/beevik/etree"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

type PubKeys struct {
	Key    *rsa.PublicKey
	OldKey *rsa.PublicKey
}

type Verifier interface {
	Verify(accessToken, url, timestamp, sig64, oldSig64 string) error
}

type VerifyHandler struct {
	discoveryURL string
	insecure     bool
	logger       log.Logger
}

func NewVerifyHandler(discoveryURL string, insecure bool, logger log.Logger) Verifier {
	return &VerifyHandler{
		discoveryURL: discoveryURL,
		insecure:     insecure,
		logger:       logger,
	}
}

// Verify the request comes from a trusted source
// All the provided parameters are strings:
// * accessToken: The access token used for this request (targeting this collaboration service)
// * url: The full url for this request, including scheme, host and all query parameters,
// something like "https://wopiserver.test.private/wopi/file/abcbcbd?access_token=oiuiu" or
// "http://wopiserver:8888/wopi/file/abcdef?access_token=zzxxyy"
// * timestamp: The timestamp provided by the WOPI app in the "X-WOPI-TimeStamp" header, as string
// * sig64: The base64-encoded signature, which should come directly from the "X-WOPI-Proof" header
// * olSig64: The base64-encoded previous signature, coming from the "X-WOPI-ProofOld" header
//
// The public keys will be obtained from the /hosting/discovery path of the target WOPI app.
// Note that the method will perform the following checks in that order:
// * current signature with the current key
// * old signature with the current key
// * current signature with the old key
// If all of those checks are wrong, the method will fail, and the request should be rejected.
//
// The method will return an error if something fails, or nil if everything is ok
func (vh *VerifyHandler) Verify(accessToken, url, timestamp, sig64, oldSig64 string) error {
	// need to decode the signatures
	signature, err := base64.StdEncoding.DecodeString(sig64)
	if err != nil {
		return err
	}

	var oldSignature []byte
	if oldSig64 != "" {
		if oldSig, err := base64.StdEncoding.DecodeString(oldSig64); err != nil {
			return nil
		} else {
			oldSignature = oldSig
		}
	}

	// fetch the public keys
	pubkeys, err := vh.fetchPublicKeys()
	if err != nil {
		return err
	}

	// build and hash the expected proof
	expectedProof := vh.generateProof(accessToken, url, timestamp)
	hashedProof := sha256.Sum256(expectedProof)

	// verify
	if err := rsa.VerifyPKCS1v15(pubkeys.Key, crypto.SHA256, hashedProof[:], signature); err != nil {
		if err := rsa.VerifyPKCS1v15(pubkeys.Key, crypto.SHA256, hashedProof[:], oldSignature); err != nil {
			if pubkeys.OldKey != nil {
				return rsa.VerifyPKCS1v15(pubkeys.OldKey, crypto.SHA256, hashedProof[:], signature)
			} else {
				return err
			}
		}
	}
	return nil
}

// generateProof will generated a expected proof to be verified later.
// The method will return a slice of bytes with the proof (consider it binary
// data).
// The bytes will need to be hashed later in order to perform the verification
func (vh *VerifyHandler) generateProof(accessToken, url, timestamp string) []byte {
	tokenBytes := []byte(accessToken)
	tokenLen := len(tokenBytes)
	tokenLenBytes := big.NewInt(int64(tokenLen)).FillBytes(make([]byte, 4))

	// url needs to be uppercase
	urlBytes := []byte(strings.ToUpper(url))
	urlLen := len(urlBytes)
	urlLenBytes := big.NewInt(int64(urlLen)).FillBytes(make([]byte, 4))

	stampBigInt, _ := new(big.Int).SetString(timestamp, 10)
	stampBytes := stampBigInt.FillBytes(make([]byte, 8))
	stampLen := len(stampBytes)
	stampLenBytes := big.NewInt(int64(stampLen)).FillBytes(make([]byte, 4))

	proof := new(bytes.Buffer)
	proof.Write(tokenLenBytes)
	proof.Write(tokenBytes)
	proof.Write(urlLenBytes)
	proof.Write(urlBytes)
	proof.Write(stampLenBytes)
	proof.Write(stampBytes)
	return proof.Bytes()
}

// fetchPublicKeys will fetch the public keys from the /hosting/discovery URL
// of the provided WOPI app.
// It will return a PubKeys struct to hold the public keys based on the modulus
// and exponent found.
// The PubKeys returned might be either nil (with the non-nil error), or might
// contain only a PubKeys.Key field (the PubKeys.OldKey might be nil)
func (vh *VerifyHandler) fetchPublicKeys() (*PubKeys, error) {
	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: vh.insecure,
			},
		},
	}

	httpResp, err := httpClient.Get(vh.discoveryURL)
	if err != nil {
		vh.logger.Error().
			Err(err).
			Str("WopiAppUrl", vh.discoveryURL).
			Msg("WopiDiscovery: failed to access wopi app url")
		return nil, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		vh.logger.Error().
			Str("WopiAppUrl", vh.discoveryURL).
			Int("HttpCode", httpResp.StatusCode).
			Msg("WopiDiscovery: wopi app url failed with unexpected code")
		return nil, err
	}

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(httpResp.Body); err != nil {
		return nil, err
	}

	root := doc.SelectElement("wopi-discovery")
	if root == nil {
		return nil, errors.New("wopi-discovery element not found in the XML body")
	}

	proofKey := root.SelectElement("proof-key")
	if proofKey == nil {
		return nil, errors.New("proof-key element not found in the XML body")
	}

	mod64 := proofKey.SelectAttrValue("modulus", "")
	exp64 := proofKey.SelectAttrValue("exponent", "")
	oldMod64 := proofKey.SelectAttrValue("oldmodulus", "")
	oldExp64 := proofKey.SelectAttrValue("oldexponent", "")

	if mod64 == "" || exp64 == "" {
		return nil, errors.New("modulus or exponent not found in the proof-key element")
	}

	keys := &PubKeys{
		Key: vh.keyFromBase64(mod64, exp64),
	}

	if oldMod64 != "" && oldExp64 != "" {
		keys.OldKey = vh.keyFromBase64(oldMod64, oldExp64)
	}

	return keys, nil
}

// keyFromBase64 will create a rsa public key from the provided modulus and
// exponent, both encoded with base64.
// If any of the provided strings can't be decoded, nil will be returned.
func (vh *VerifyHandler) keyFromBase64(mod64, exp64 string) *rsa.PublicKey {
	dataMod, err := base64.StdEncoding.DecodeString(mod64)
	if err != nil {
		return nil
	}

	dataE, err := base64.StdEncoding.DecodeString(exp64)
	if err != nil {
		return nil
	}

	pub := &rsa.PublicKey{
		N: new(big.Int).SetBytes(dataMod),
		E: int(new(big.Int).SetBytes(dataE).Int64()),
	}
	return pub
}
