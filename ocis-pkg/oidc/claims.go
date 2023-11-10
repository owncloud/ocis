package oidc

import (
	"fmt"
	"strings"
)

const (
	Iss               = "iss"
	Sub               = "sub"
	Email             = "email"
	Name              = "name"
	PreferredUsername = "preferred_username"
	UIDNumber         = "uidnumber"
	GIDNumber         = "gidnumber"
	Groups            = "groups"
	OwncloudUUID      = "ownclouduuid"
	OcisRoutingPolicy = "ocis.routing.policy"
)

// SplitWithEscaping splits s into segments using separator which can be escaped using the escape string
// See https://codereview.stackexchange.com/a/280193
func SplitWithEscaping(s string, separator string, escapeString string) []string {
	a := strings.Split(s, separator)

	for i := len(a) - 2; i >= 0; i-- {
		if strings.HasSuffix(a[i], escapeString) {
			a[i] = a[i][:len(a[i])-len(escapeString)] + separator + a[i+1]
			a = append(a[:i+1], a[i+2:]...)
		}
	}
	return a
}

func ReadStringClaim(path string, claims map[string]interface{}) (string, error) {
	// happy path
	value, _ := claims[path].(string)
	if value != "" {
		return value, nil
	}

	// try splitting path at .
	segments := SplitWithEscaping(path, ".", "\\")
	subclaims := claims
	lastSegment := len(segments) - 1
	for i := range segments {
		if i < lastSegment {
			if castedClaims, ok := subclaims[segments[i]].(map[string]interface{}); ok {
				subclaims = castedClaims
			} else if castedClaims, ok := subclaims[segments[i]].(map[interface{}]interface{}); ok {
				subclaims = make(map[string]interface{}, len(castedClaims))
				for k, v := range castedClaims {
					if s, ok := k.(string); ok {
						subclaims[s] = v
					} else {
						return "", fmt.Errorf("could not walk claims path, key '%v' is not a string", k)
					}
				}
			}
		} else {
			if value, _ = subclaims[segments[i]].(string); value != "" {
				return value, nil
			}
		}
	}

	return value, fmt.Errorf("claim path '%s' not set or empty", path)
}
