package app

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	WopiContext WopiContext `json:"WopiContext"`
	jwt.StandardClaims
}
