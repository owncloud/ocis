package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/libregraph/idm/pkg/ldbbolt"
	"github.com/libregraph/idm/server"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	ocisdefaults "github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/idm/pkg/config"
	"github.com/owncloud/ocis/v2/services/idm/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/idm/pkg/logging"
	"github.com/urfave/cli/v2"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/term"
)

// User account types accepted by the --user-type flag.
const (
	userTypeUser    = "user"
	userTypeService = "service"
)

// ResetPassword is the entrypoint for the resetpassword command
func ResetPassword(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "resetpassword",
		Usage:    "Reset user password",
		Category: "password reset",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "user-name",
				Aliases: []string{"u"},
				Usage:   "User name",
				Value:   "admin",
			},
			&cli.StringFlag{
				Name:  "user-type",
				Usage: "Type of user account: 'user' (ou=users) or 'service' (ou=sysusers)",
				Value: userTypeUser,
			},
		},
		Before: func(_ *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			ctx, cancel := context.WithCancel(c.Context)

			defer cancel()

			userType := c.String("user-type")
			if userType != userTypeUser && userType != userTypeService {
				return fmt.Errorf("invalid --user-type %q: must be %q or %q", userType, userTypeUser, userTypeService)
			}

			return resetPassword(ctx, logger, cfg, c.String("user-name"), userType)
		},
	}
}

func resetPassword(_ context.Context, logger log.Logger, cfg *config.Config, userName string, userType string) error {
	servercfg := server.Config{
		Logger:      log.LogrusWrap(logger.Logger),
		LDAPHandler: "boltdb",
		LDAPBaseDN:  "o=libregraph-idm",

		BoltDBFile: cfg.IDM.DatabasePath,
	}

	ou := "users"
	if userType == userTypeService {
		ou = "sysusers"
	}
	userDN := fmt.Sprintf("uid=%s,ou=%s,%s", userName, ou, servercfg.LDAPBaseDN)
	fmt.Printf("Resetting password for user '%s'.\n", userDN)
	if _, err := os.Stat(servercfg.BoltDBFile); errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "IDM database does not exist.\n")
		return err
	}

	newPw, err := getPassword()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		return err
	}

	bdb := &ldbbolt.LdbBolt{}

	opts := bolt.Options{
		Timeout: 1 * time.Millisecond,
	}
	if err := bdb.Configure(servercfg.Logger, servercfg.LDAPBaseDN, servercfg.BoltDBFile, &opts); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: '%s'. Please stop any running ocis/idm instance, as this tool requires exclusive access to the database.\n", err)
		return err
	}
	defer bdb.Close()

	if err := bdb.Initialize(); err != nil {
		return err
	}

	pwRequest := ldap.NewPasswordModifyRequest(userDN, "", newPw)
	if err := bdb.UpdatePassword(pwRequest); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update user password: %v\n", err)
	}
	fmt.Printf("Password for user '%s' updated.\n", userDN)

	if userType == userTypeService {
		syncServiceUserBindConfig(userName, newPw)
	}
	return nil
}

// syncServiceUserBindConfig keeps ocis.yaml in sync after a service user's
// password was changed in the IDM database. Services bind to LDAP as the
// service users (e.g. uid=reva,ou=sysusers) using the bind_password from their
// config; if only the directory entry is changed those binds start failing on
// the next restart, which manifests as admin login returning 401. It rewrites
// the matching keys in ocis.yaml when present and always prints the env vars
// the operator must ensure carry the new value (covering env-var and
// distributed deployments this tool cannot rewrite).
func syncServiceUserBindConfig(userName, newPassword string) {
	if _, known := serviceUserConfigKeys[userName]; !known {
		return
	}

	configFile := path.Join(ocisdefaults.BaseConfigPath(), "ocis.yaml")
	if data, err := os.ReadFile(configFile); err == nil {
		out, updated, err := syncServiceUserPasswordInConfig(data, userName, newPassword)
		switch {
		case err != nil:
			fmt.Fprintf(os.Stderr, "Warning: could not update bind passwords in %q: %v\n", configFile, err)
		case len(updated) > 0:
			if err := os.WriteFile(configFile, out, configFilePerm); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not write updated config %q: %v\n", configFile, err)
			} else {
				fmt.Printf("Updated bind passwords in %q:\n", configFile)
				for _, key := range updated {
					fmt.Printf("  - %s\n", key)
				}
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "Warning: could not read config %q: %v\n", configFile, err)
	}

	fmt.Printf("\nIMPORTANT: the '%s' service user is used by oCIS services to bind to LDAP.\n", userName)
	fmt.Printf("Make sure the following environment variable(s) are set to the new password\n")
	fmt.Printf("wherever they are configured (env, secrets, per-service yaml), then restart oCIS:\n")
	fmt.Printf("  - OCIS_LDAP_BIND_PASSWORD (shared) or the per-service variables below\n")
	for _, env := range serviceUserBindEnvVars[userName] {
		fmt.Printf("  - %s\n", env)
	}
}

// serviceUserBindEnvVars lists the per-service bind_password environment
// variables that carry the given service user's password, for the operator
// guidance message. OCIS_LDAP_BIND_PASSWORD is a single shared variable used by
// all of these services, so it is noted separately by the caller rather than
// per user.
var serviceUserBindEnvVars = map[string][]string{
	"reva":       {"AUTH_BASIC_LDAP_BIND_PASSWORD", "USERS_LDAP_BIND_PASSWORD", "GROUPS_LDAP_BIND_PASSWORD", "IDM_REVASVC_PASSWORD"},
	"libregraph": {"GRAPH_LDAP_BIND_PASSWORD", "IDM_SVC_PASSWORD"},
	"idp":        {"IDP_LDAP_BIND_PASSWORD", "IDM_IDPSVC_PASSWORD"},
}

func getPassword() (string, error) {
	fmt.Print("Enter new password: ")
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Println("")
	fmt.Print("Re-enter new password: ")
	bytePasswordVerify, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Println("")

	password := string(bytePassword)
	passwordVerify := string(bytePasswordVerify)

	if password != passwordVerify {
		return "", errors.New("Passwords do not match")
	}
	return password, nil
}
