package helpers

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
)

// ParseWopiFileID extracts the file id from a wopi path
//
// If the file id is a jwt, it will be decoded and the file id will be extracted from the jwt claims.
// If the file id is not a jwt, it will be returned as is.
func ParseWopiFileID(cfg *config.Config, path string) string {
	s := strings.Split(path, "/")
	if len(s) < 4 || (s[1] != "wopi" && s[2] != "files") {
		return path
	}
	// check if the fileid is a jwt
	if strings.Contains(s[3], ".") {
		token, err := jwt.Parse(s[3], func(_ *jwt.Token) (interface{}, error) {
			return []byte(cfg.Wopi.ProxySecret), nil
		})
		if err != nil {
			return s[3]
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return s[3]
		}

		f, ok := claims["f"].(string)
		if !ok {
			return s[3]
		}
		return f
	}
	// fileid is not a jwt
	return s[3]
}
