package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/libregraph/idm/pkg/ldbbolt"
	"github.com/libregraph/idm/server"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/idm/pkg/config"
	"github.com/owncloud/ocis/v2/services/idm/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/idm/pkg/logging"
	"github.com/urfave/cli/v2"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/term"
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
		},
		Before: func(_ *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
			}()

			defer cancel()
			return resetPassword(ctx, logger, cfg, c.String("user-name"))
		},
	}
}

func resetPassword(_ context.Context, logger log.Logger, cfg *config.Config, userName string) error {
	servercfg := server.Config{
		Logger:      log.LogrusWrap(logger.Logger),
		LDAPHandler: "boltdb",
		LDAPBaseDN:  "o=libregraph-idm",

		BoltDBFile: cfg.IDM.DatabasePath,
	}

	userDN := fmt.Sprintf("uid=%s,ou=users,%s", userName, servercfg.LDAPBaseDN)
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
	return nil
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
