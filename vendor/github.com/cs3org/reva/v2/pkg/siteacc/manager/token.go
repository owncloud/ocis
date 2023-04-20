// Copyright 2018-2020 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package manager

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
)

type userToken struct {
	jwt.StandardClaims

	User  string `json:"user"`
	Scope string `json:"scope"`
}

const (
	tokenKeyLength = 16
	tokenIssuer    = "sciencemesh_siteacc"
)

var (
	tokenSecret string
)

func generateUserToken(user string, scope string, timeout int) (string, error) {
	// Create a JWT as the user token
	claims := userToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(timeout) * time.Second).Unix(),
			Issuer:    tokenIssuer,
			IssuedAt:  time.Now().Unix(),
		},
		User:  user,
		Scope: scope,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", errors.Wrapf(err, "error signing token with claims %+v", claims)
	}

	return signedToken, nil
}

func extractUserToken(token string) (*userToken, error) {
	// Parse the token and try to extract the claims
	parsedToken, err := jwt.ParseWithClaims(token, &userToken{}, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return nil, errors.Wrap(err, "error parsing token")
	}

	if claims, ok := parsedToken.Claims.(*userToken); ok && parsedToken.Valid {
		if claims.Issuer != tokenIssuer {
			return nil, errors.Errorf("invalid token issuer")
		}

		return claims, nil
	}

	return nil, errors.Errorf("invalid token")
}

func init() {
	// Generate the token secret randomly
	tokenSecret = password.MustGenerate(tokenKeyLength, tokenKeyLength/4, tokenKeyLength/4, false, true)
}
