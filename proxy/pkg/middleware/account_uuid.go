package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	revauser "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	acc "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
)

func getAccount(l log.Logger, ac acc.AccountsService, query string) (account *acc.Account, status int) {
	resp, err := ac.ListAccounts(context.Background(), &acc.ListAccountsRequest{
		Query:    query,
		PageSize: 2,
	})

	if err != nil {
		l.Error().Err(err).Str("query", query).Msgf("Error fetching from accounts-service")
		status = http.StatusInternalServerError
		return
	}

	if len(resp.Accounts) <= 0 {
		l.Error().Str("query", query).Msgf("Account not found")
		status = http.StatusNotFound
		return
	}

	if len(resp.Accounts) > 1 {
		l.Error().Str("query", query).Msgf("More than one account found. Not logging user in.")
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
		l.Error().Err(err).Interface("account", req.Account).Msg("could not create account")
		return nil, http.StatusInternalServerError
	}

	return created, 0
}

// AccountUUID provides a middleware which mints a jwt and adds it to the proxied request based
// on the oidc-claims
func AccountUUID(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	return func(next http.Handler) http.Handler {
		// TODO: handle error
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret":  opt.TokenManagerConfig.JWTSecret,
			"expires": int64(60),
		})
		if err != nil {
			opt.Logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := opt.Logger
			claims := oidc.FromContext(r.Context())
			var account *acc.Account
			var status int
			switch {
			case claims == nil:
				login, password, ok := r.BasicAuth()
				if opt.EnableBasicAuth && ok {
					l.Warn().Msg("basic auth enabled, use only for testing or development")
					account, status = getAccount(l, opt.AccountsClient, fmt.Sprintf("login eq '%s' and password eq '%s'", strings.ReplaceAll(login, "'", "''"), strings.ReplaceAll(password, "'", "''")))
					if status == 0 {
						// fake claims for the subsequent code flow
						claims = &oidc.StandardClaims{
							Iss: opt.OIDCIss,
						}
					} else {
						// we are still forwarding the request, the service is responsible for rendering the error
						// TODO or render a full blown ocs xml / json error response
						next.ServeHTTP(w, r)
						// tell client to reauthenticate
						//w.WriteHeader(http.StatusUnauthorized)
						return
					}
				} else {
					next.ServeHTTP(w, r)
					return
				}
			case claims.Email != "":
				account, status = getAccount(l, opt.AccountsClient, fmt.Sprintf("mail eq '%s'", strings.ReplaceAll(claims.Email, "'", "''")))
			case claims.PreferredUsername != "":
				account, status = getAccount(l, opt.AccountsClient, fmt.Sprintf("preferred_name eq '%s'", strings.ReplaceAll(claims.PreferredUsername, "'", "''")))
			case claims.OcisID != "":
				account, status = getAccount(l, opt.AccountsClient, fmt.Sprintf("id eq '%s'", strings.ReplaceAll(claims.OcisID, "'", "''")))
			default:
				// TODO allow lookup by custom claim, eg an id ... or sub
				l.Error().Err(err).Msg("Could not lookup account, no mail or preferred_username claim set")
				w.WriteHeader(http.StatusInternalServerError)
			}
			if status != 0 || account == nil {
				if opt.AutoprovisionAccounts && status == http.StatusNotFound {
					account, status = createAccount(l, claims, opt.AccountsClient)
					if status != 0 {
						w.WriteHeader(status)
						return
					}
				} else {
					w.WriteHeader(status)
					return
				}
			}
			if !account.AccountEnabled {
				l.Debug().Interface("account", account).Msg("account is disabled")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			groups := make([]string, len(account.MemberOf))
			for i := range account.MemberOf {
				// reva needs the unix group name
				groups[i] = account.MemberOf[i].OnPremisesSamAccountName
			}

			// fetch active roles from ocis-settings
			assignmentResponse, err := opt.SettingsRoleService.ListRoleAssignments(r.Context(), &settings.ListRoleAssignmentsRequest{AccountUuid: account.Id})
			roleIDs := make([]string, 0)
			if err != nil {
				l.Err(err).Str("accountID", account.Id).Msg("failed to fetch role assignments")
			} else {
				for _, assignment := range assignmentResponse.Assignments {
					roleIDs = append(roleIDs, assignment.RoleId)
				}
			}

			l.Debug().Interface("claims", claims).Interface("account", account).Msgf("Associated claims with uuid")
			user := &revauser.User{
				Id: &revauser.UserId{
					OpaqueId: account.Id,
					Idp:      claims.Iss,
				},
				Username:     account.OnPremisesSamAccountName,
				DisplayName:  account.DisplayName,
				Mail:         account.Mail,
				MailVerified: account.ExternalUserState == "" || account.ExternalUserState == "Accepted",
				Groups:       groups,
				Opaque: &types.Opaque{
					Map: map[string]*types.OpaqueEntry{},
				},
			}

			user.Opaque.Map["uid"] = &types.OpaqueEntry{
				Decoder: "plain",
				Value:   []byte(strconv.FormatInt(account.UidNumber, 10)),
			}
			user.Opaque.Map["gid"] = &types.OpaqueEntry{
				Decoder: "plain",
				Value:   []byte(strconv.FormatInt(account.GidNumber, 10)),
			}

			// encode roleIDs as json string
			roleIDsJSON, jsonErr := json.Marshal(roleIDs)
			if jsonErr != nil {
				l.Err(jsonErr).Str("accountID", account.Id).Msg("failed to marshal roleIDs into json")
			} else {
				user.Opaque.Map["roles"] = &types.OpaqueEntry{
					Decoder: "json",
					Value:   roleIDsJSON,
				}
			}

			token, err := tokenManager.MintToken(r.Context(), user)

			if err != nil {
				l.Error().Err(err).Msgf("Could not mint token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Header.Set("x-access-token", token)
			next.ServeHTTP(w, r)
		})
	}
}
