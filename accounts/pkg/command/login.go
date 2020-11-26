package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	revaUser "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/flagset"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	"strconv"
)

// Login is a cli tool to get an x-access-token
func Login(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "Get an authenticated backend token",
		Flags: flagset.LoginWithConfig(cfg),
		Action: func(c *cli.Context) error {
			accSvcID := cfg.GRPC.Namespace + "." + cfg.Server.Name
			secret := cfg.TokenManager.JWTSecret
			issuer := c.String("issuer")

			accSvc := accounts.NewAccountsService(accSvcID, grpc.DefaultClient)
			rolesSvc := settings.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient)
			tokenManager, err := jwt.New(map[string]interface{}{
				"secret":  secret,
				"expires": int64(3600),
			})

			if err != nil {
				return err
			}

			res, err := accSvc.ListAccounts(c.Context, &accounts.ListAccountsRequest{
				Query: fmt.Sprintf("login eq '%s' and password eq '%s'", c.String("username"), c.String("password")),
			})

			if err != nil {
				return err
			}

			if len(res.Accounts) == 0 {
				return errors.New("login failed: account not found")
			}

			if len(res.Accounts) > 1 {
				return errors.New("login failed: multiple accounts with this usenrame found")
			}

			account := res.Accounts[0]

			if !account.AccountEnabled {
				return errors.New("account is disabled")
			}

			groups := make([]string, len(account.MemberOf))
			for i := range account.MemberOf {
				// reva needs the unix group name
				groups[i] = account.MemberOf[i].OnPremisesSamAccountName
			}

			// fetch active roles from ocis-settings
			assignmentResponse, err := rolesSvc.ListRoleAssignments(context.Background(), &settings.ListRoleAssignmentsRequest{AccountUuid: account.Id})
			roleIDs := make([]string, 0)
			if err != nil {
				return err
			} else {
				for _, assignment := range assignmentResponse.Assignments {
					roleIDs = append(roleIDs, assignment.RoleId)
				}
			}

			user := &revaUser.User{
				Id: &revaUser.UserId{
					OpaqueId: account.Id,
					Idp:      issuer,
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
				return jsonErr
			} else {
				user.Opaque.Map["roles"] = &types.OpaqueEntry{
					Decoder: "json",
					Value:   roleIDsJSON,
				}
			}

			token, err := tokenManager.MintToken(context.Background(), user)

			if err != nil {
				return err
			}

			fmt.Print(token)

			return nil
		}}
}
