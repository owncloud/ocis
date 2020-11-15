package middleware

import (
	"encoding/json"
	"fmt"
	revaUser "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	tokenPkg "github.com/cs3org/reva/pkg/token"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	"net/http"
	"strconv"
	"strings"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
)

func AccountResolver(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret":  options.TokenManagerConfig.JWTSecret,
			"expires": int64(60),
		})
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
		}

		return &accountResolver{
			next:                  next,
			logger:                logger,
			tokenManager:          tokenManager,
			accountsClient:        options.AccountsClient,
			oidcIss:               options.OIDCIss,
			autoprovisionAccounts: options.AutoprovisionAccounts,
			settingsRoleService:   options.SettingsRoleService,
		}
	}
}

type accountResolver struct {
	oidcIss               string
	autoprovisionAccounts bool
	next                  http.Handler
	logger                log.Logger
	tokenManager          tokenPkg.Manager
	accountsClient        accounts.AccountsService
	settingsRoleService   settings.RoleService
}

func (m accountResolver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var account *accounts.Account
	var status int

	claims := oidc.FromContext(req.Context())

	if claims == nil {
		m.next.ServeHTTP(w, req)
		return
	}

	switch {
	case claims.Email != "":
		account, status = getAccount(m.logger, m.accountsClient, fmt.Sprintf("mail eq '%s'", strings.ReplaceAll(claims.Email, "'", "''")))
	case claims.PreferredUsername != "":
		account, status = getAccount(m.logger, m.accountsClient, fmt.Sprintf("preferred_name eq '%s'", strings.ReplaceAll(claims.PreferredUsername, "'", "''")))
	case claims.OcisID != "":
		account, status = getAccount(m.logger, m.accountsClient, fmt.Sprintf("id eq '%s'", strings.ReplaceAll(claims.OcisID, "'", "''")))
	default:
		// TODO allow lookup by custom claim, eg an id ... or sub
		m.logger.Error().Msg("Could not lookup accountResolver, no mail or preferred_username claim set")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if m.autoprovisionAccounts && status == http.StatusNotFound {
		account, status = createAccount(m.logger, claims, m.accountsClient)
	}

	if status != 0 || account == nil {
		w.WriteHeader(status)
		return
	}

	if !account.AccountEnabled {
		m.logger.Debug().Interface("accountResolver", account).Msg("accountResolver is disabled")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groups := make([]string, len(account.MemberOf))
	for i := range account.MemberOf {
		// reva needs the unix group name
		groups[i] = account.MemberOf[i].OnPremisesSamAccountName
	}

	// fetch active roles from ocis-settings
	assignmentResponse, err := m.settingsRoleService.ListRoleAssignments(req.Context(), &settings.ListRoleAssignmentsRequest{AccountUuid: account.Id})
	roleIDs := make([]string, 0)
	if err != nil {
		m.logger.Err(err).Str("accountID", account.Id).Msg("failed to fetch role assignments")
	} else {
		for _, assignment := range assignmentResponse.Assignments {
			roleIDs = append(roleIDs, assignment.RoleId)
		}
	}

	m.logger.Debug().Interface("claims", claims).Interface("accountResolver", account).Msgf("Associated claims with uuid")

	user := &revaUser.User{
		Id: &revaUser.UserId{
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
		m.logger.Err(jsonErr).Str("accountID", account.Id).Msg("failed to marshal roleIDs into json")
	} else {
		user.Opaque.Map["roles"] = &types.OpaqueEntry{
			Decoder: "json",
			Value:   roleIDsJSON,
		}
	}

	token, err := m.tokenManager.MintToken(req.Context(), user)

	if err != nil {
		m.logger.Error().Err(err).Msgf("Could not mint token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set("x-access-token", token)

	m.next.ServeHTTP(w, req)
}
