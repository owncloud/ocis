package jwx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// SignClaims signs the given claims using a given key and a method
func SignClaims(claims jwt.Claims, key interface{}, method jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	return token.SignedString(key)
}

// DecodeAccessTokenHeader decodes the header of the accessToken
func DecodeAccessTokenHeader(token string) (*DecodedAccessTokenHeader, error) {
	const errMessage = "could not decode access token header"
	token = strings.Replace(token, "Bearer ", "", 1)
	headerString := strings.Split(token, ".")
	decodedData, err := base64.RawStdEncoding.DecodeString(headerString[0])
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	result := &DecodedAccessTokenHeader{}
	err = json.Unmarshal(decodedData, result)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	return result, nil
}

func toBigInt(v string) (*big.Int, error) {
	decRes, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		return nil, err
	}

	res := big.NewInt(0)
	res.SetBytes(decRes)
	return res, nil
}

var (
	curves = map[string]elliptic.Curve{
		"P-224": elliptic.P224(),
		"P-256": elliptic.P256(),
		"P-384": elliptic.P384(),
		"P-521": elliptic.P521(),
	}
)

func decodeECDSAPublicKey(x, y, crv *string) (*ecdsa.PublicKey, error) {
	const errMessage = "could not decode public key"

	xInt, err := toBigInt(*x)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	yInt, err := toBigInt(*y)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}
	var c elliptic.Curve
	var ok bool
	if c, ok = curves[*crv]; !ok {
		return nil, errors.Wrap(fmt.Errorf("unknown curve alg: %s", *crv), errMessage)
	}
	return &ecdsa.PublicKey{X: xInt, Y: yInt, Curve: c}, nil
}

func decodeRSAPublicKey(e, n *string) (*rsa.PublicKey, error) {
	const errMessage = "could not decode public key"

	nInt, err := toBigInt(*n)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	decE, err := base64.RawURLEncoding.DecodeString(*e)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	var eBytes []byte
	if len(decE) < 8 {
		eBytes = make([]byte, 8-len(decE), 8)
		eBytes = append(eBytes, decE...)
	} else {
		eBytes = decE
	}

	eReader := bytes.NewReader(eBytes)
	var eInt uint64
	err = binary.Read(eReader, binary.BigEndian, &eInt)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	pKey := rsa.PublicKey{N: nInt, E: int(eInt)}
	return &pKey, nil
}

// DecodeAccessTokenRSACustomClaims decodes string access token into jwt.Token
func DecodeAccessTokenRSACustomClaims(accessToken string, e, n *string, customClaims jwt.Claims) (*jwt.Token, error) {
	const errMessage = "could not decode accessToken with custom claims"
	accessToken = strings.Replace(accessToken, "Bearer ", "", 1)

	rsaPublicKey, err := decodeRSAPublicKey(e, n)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	token2, err := jwt.ParseWithClaims(accessToken, customClaims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rsaPublicKey, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}
	return token2, nil
}

// DecodeAccessTokenECDSACustomClaims decodes string access token into jwt.Token
func DecodeAccessTokenECDSACustomClaims(accessToken string, x, y, crv *string, customClaims jwt.Claims) (*jwt.Token, error) {
	const errMessage = "could not decode accessToken"
	accessToken = strings.Replace(accessToken, "Bearer ", "", 1)

	publicKey, err := decodeECDSAPublicKey(x, y, crv)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	token2, err := jwt.ParseWithClaims(accessToken, customClaims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}
	return token2, nil
}
