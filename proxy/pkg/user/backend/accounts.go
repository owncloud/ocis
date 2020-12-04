package backend

import (
	"context"
	"fmt"
	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	"net/http"
	"strconv"
	"strings"
)

// NewAccountsServiceUserBackend creates a user-provider which fetches users from the ocis accounts-service
func NewAccountsServiceUserBackend(ac accounts.AccountsService, rs settings.RoleService, oidcISS string, logger log.Logger) UserBackend {
	return &accountsServiceBackend{
		accountsClient:      ac,
		settingsRoleService: rs,
		OIDCIss:             oidcISS,
		logger:              logger,
	}
}

type accountsServiceBackend struct {
	accountsClient      accounts.AccountsService
	settingsRoleService settings.RoleService
	OIDCIss             string
	logger              log.Logger
}

func (a accountsServiceBackend) GetUserByClaims(ctx context.Context, claim, value string, withRoles bool) (*cs3.User, error) {
	var account *accounts.Account
	var status int
	var query string

	switch claim {
	case "mail":
		query = fmt.Sprintf("mail eq '%s'", strings.ReplaceAll(value, "'", "''"))
	case "username":
		query = fmt.Sprintf("preferred_name eq '%s'", strings.ReplaceAll(value, "'", "''"))
	case "id":
		query = fmt.Sprintf("id eq '%s'", strings.ReplaceAll(value, "'", "''"))
	default:
		return nil, fmt.Errorf("invalid user by claim lookup must be  'mail', 'username' or 'id")
	}

	account, status = a.getAccount(ctx, query)

	if status != 0 || account == nil {
		return nil, fmt.Errorf("could not get account, got status: %d", status)
	}

	if !account.AccountEnabled {
		// TODO: handle not found
		return nil, nil
	}

	user := &cs3.User{
		Id: &cs3.UserId{
			OpaqueId: account.Id,
			Idp:      a.OIDCIss,
		},
		Username:     account.OnPremisesSamAccountName,
		DisplayName:  account.DisplayName,
		Mail:         account.Mail,
		MailVerified: account.ExternalUserState == "" || account.ExternalUserState == "Accepted",
		Groups:       expandGroups(account),
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

	if !withRoles {
		return user, nil
	}

	if err := lazyLoadRoles(ctx, user, a.settingsRoleService); err != nil {
		a.logger.Warn().Err(err).Msgf("Could not load roles... continuing without")
	}

	return user, nil

}

// Authenticate authenticates against the accounts services and returns the user on success
func (a *accountsServiceBackend) Authenticate(ctx context.Context, username string, password string) (*cs3.User, error) {
	query := fmt.Sprintf(
		"login eq '%s' and password eq '%s'",
		strings.ReplaceAll(username, "'", "''"),
		strings.ReplaceAll(password, "'", "''"),
	)
	account, status := a.getAccount(ctx, query)

	if status != 0 {
		return nil, fmt.Errorf("could not authenticate with username, password for user %s. Status: %d", username, status)
	}

	user := &cs3.User{
		Id: &cs3.UserId{
			OpaqueId: account.Id,
			Idp:      a.OIDCIss,
		},
		Username:     account.OnPremisesSamAccountName,
		DisplayName:  account.DisplayName,
		Mail:         account.Mail,
		MailVerified: account.ExternalUserState == "" || account.ExternalUserState == "Accepted",
		Groups:       expandGroups(account),
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

	if err := lazyLoadRoles(ctx, user, a.settingsRoleService); err != nil {
		a.logger.Warn().Err(err).Msgf("Could not load roles... continuing without")
	}

	return user, nil
}

func (a accountsServiceBackend) GetUserGroups(ctx context.Context, userID string) {
	panic("implement me")
}

func (a *accountsServiceBackend) getAccount(ctx context.Context, query string) (account *accounts.Account, status int) {
	resp, err := a.accountsClient.ListAccounts(ctx, &accounts.ListAccountsRequest{
		Query:    query,
		PageSize: 2,
	})

	if err != nil {
		a.logger.Error().Err(err).Str("query", query).Msgf("error fetching from accounts-service")
		status = http.StatusInternalServerError
		return
	}

	if len(resp.Accounts) <= 0 {
		a.logger.Error().Str("query", query).Msgf("account not found")
		status = http.StatusNotFound
		return
	}

	if len(resp.Accounts) > 1 {
		a.logger.Error().Str("query", query).Msgf("more than one account found, aborting")
		status = http.StatusForbidden
		return
	}

	account = resp.Accounts[0]
	return
}

func expandGroups(account *accounts.Account) []string {
	groups := make([]string, len(account.MemberOf))
	for i := range account.MemberOf {
		// reva needs the unix group name
		groups[i] = account.MemberOf[i].OnPremisesSamAccountName
	}
	return groups
}

// lazyLoadRoles adds roles from the roles-service to the user-struct by mutating an existing struct
func lazyLoadRoles(ctx context.Context, u *cs3.User, ss settings.RoleService) error {
	roleIDs, err := loadRolesIDs(ctx, u.Id.OpaqueId, ss)
	if err != nil {
		return err
	}

	if len(roleIDs) == 0 {
		return nil

	}

	enc, err := encodeRoleIDs(roleIDs)
	if err != nil {
		return err
	}

	u.Opaque = &types.Opaque{
		Map: map[string]*types.OpaqueEntry{
			"roles": enc,
		},
	}

	return nil
}
