package jwt

import "github.com/golang-jwt/jwt/v4"

type ThumbnailClaims struct {
	jwt.RegisteredClaims
	Key string `json:"key"`
}
