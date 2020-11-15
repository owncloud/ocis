package middleware

import (
	"context"
	"errors"
	gOidc "github.com/coreos/go-oidc"
	acc "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"golang.org/x/oauth2"
	"net/http"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
)

var (
	ErrInvalidToken = errors.New("invalid or missing token")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
)

type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*gOidc.UserInfo, error)
}

func getAccount(logger log.Logger, ac acc.AccountsService, query string) (account *acc.Account, status int) {
	resp, err := ac.ListAccounts(context.Background(), &acc.ListAccountsRequest{
		Query:    query,
		PageSize: 2,
	})

	if err != nil {
		logger.Error().Err(err).Str("query", query).Msgf("Error fetching from accounts-service")
		status = http.StatusInternalServerError
		return
	}

	if len(resp.Accounts) <= 0 {
		logger.Error().Str("query", query).Msgf("AccountResolver not found")
		status = http.StatusNotFound
		return
	}

	if len(resp.Accounts) > 1 {
		logger.Error().Str("query", query).Msgf("More than one accountResolver found. Not logging user in.")
		status = http.StatusForbidden
		return
	}

	account = resp.Accounts[0]
	return
}

func createAccount(l log.Logger, claims *oidc.StandardClaims, ac acc.AccountsService) (*acc.Account, int) {
	// TODO check if fields are missing.
	req := &acc.CreateAccountRequest{
		Account: &acc.Account{
			DisplayName:              claims.DisplayName,
			PreferredName:            claims.PreferredUsername,
			OnPremisesSamAccountName: claims.PreferredUsername,
			Mail:                     claims.Email,
			CreationType:             "LocalAccount",
			AccountEnabled:           true,
			// TODO assign uidnumber and gidnumber? better do that in ocis-accounts as it can keep track of the next numbers
		},
	}
	created, err := ac.CreateAccount(context.Background(), req)
	if err != nil {
		l.Error().Err(err).Interface("accountResolver", req.Account).Msg("could not create accountResolver")
		return nil, http.StatusInternalServerError
	}

	return created, 0
}