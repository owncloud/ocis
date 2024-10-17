package jwt

import "github.com/golang-jwt/jwt/v4"

// ThumbnailClaims defines the claims for thumb-nailing
type ThumbnailClaims struct {
	jwt.RegisteredClaims
	Key string `json:"key"`
}
