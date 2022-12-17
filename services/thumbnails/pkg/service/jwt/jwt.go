package jwt

import "github.com/golang-jwt/jwt/v4"

// ThumbnailClaims
// FIXME: nolint
// nolint: revive
type ThumbnailClaims struct {
	jwt.RegisteredClaims
	Key string `json:"key"`
}
