package middleware

import "github.com/golang-jwt/jwt/v5"

// Claims contains the jwt registered claims plus the used WOPI context
type Claims struct {
	WopiContext WopiContext `json:"WopiContext"`
	jwt.RegisteredClaims
}
