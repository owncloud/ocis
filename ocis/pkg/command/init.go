package command

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/generators"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	cli "github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	authbasic "github.com/owncloud/ocis/extensions/auth-basic/pkg/config"
	authbearer "github.com/owncloud/ocis/extensions/auth-bearer/pkg/config"
	frontend "github.com/owncloud/ocis/extensions/frontend/pkg/config"
	graph "github.com/owncloud/ocis/extensions/graph/pkg/config"
	group "github.com/owncloud/ocis/extensions/group/pkg/config"
	idm "github.com/owncloud/ocis/extensions/idm/pkg/config"
	idp "github.com/owncloud/ocis/extensions/idp/pkg/config"
	ocdav "github.com/owncloud/ocis/extensions/ocdav/pkg/config"
	proxy "github.com/owncloud/ocis/extensions/proxy/pkg/config"
	storagemetadata "github.com/owncloud/ocis/extensions/storage-metadata/pkg/config"
	storageusers "github.com/owncloud/ocis/extensions/storage-users/pkg/config"
	thumbnails "github.com/owncloud/ocis/extensions/thumbnails/pkg/config"
	user "github.com/owncloud/ocis/extensions/user/pkg/config"
)

const configFilename string = "ocis.yaml" // TODO: use also a constant for reading this file
const passwordLength int = 32

// InitCommand is the entrypoint for the init command
func InitCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialise an ocis config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "insecure",
				EnvVars: []string{"OCIS_INSECURE"},
				Value:   "ask",
			},
			&cli.BoolFlag{
				Name:    "force-overwrite",
				Aliases: []string{"f"},
				EnvVars: []string{"OCIS_FORCE_CONFIG_OVERWRITE"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:  "config-path",
				Value: defaults.BaseConfigPath(),
				Usage: "config path for the ocis runtime",
			},
		},
		Action: func(c *cli.Context) error {
			insecureFlag := c.String("insecure")
			insecure := false
			if insecureFlag == "ask" {
				answer := strings.ToLower(stringPrompt("Insecure Backends? [Yes|No]"))
				if answer == "yes" || answer == "y" {
					insecure = true
				}
			} else if insecureFlag == "true" {
				insecure = true
			}
			err := createConfig(insecure, c.Bool("force-overwrite"), c.String("config-path"))
			if err != nil {
				log.Fatalf("Could not create config: %s", err)
			}
			return nil
		},
	}
}

func init() {
	register.AddCommand(InitCommand)
}

func checkConfigPath(configPath string) error {
	targetPath := path.Join(configPath, configFilename)
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("config in %s already exists", targetPath)
	}
	return nil
}

func createConfig(insecure, forceOverwrite bool, configPath string) error {
	err := checkConfigPath(configPath)
	if err != nil && !forceOverwrite {
		return err
	}
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return err
	}
	cfg := config.Config{
		TokenManager: &shared.TokenManager{},
		IDM:          &idm.Config{},
		AuthBasic: &authbasic.Config{
			AuthProviders: authbasic.AuthProviders{
				LDAP: authbasic.LDAPProvider{},
			},
		},
		Group: &group.Config{
			Drivers: group.Drivers{
				LDAP: group.LDAPDriver{},
			},
		},
		User: &user.Config{
			Drivers: user.Drivers{
				LDAP: user.LDAPDriver{},
			},
		},
		IDP: &idp.Config{},
	}

	if insecure {
		cfg.AuthBearer = &authbearer.Config{
			AuthProviders: authbearer.AuthProviders{
				OIDC: authbearer.OIDCProvider{
					Insecure: true,
				},
			},
		}
		cfg.Frontend = &frontend.Config{
			AppProvider: frontend.AppProvider{
				Insecure: true,
			},
			Archiver: frontend.Archiver{
				Insecure: true,
			},
		}
		cfg.Graph = &graph.Config{
			Spaces: graph.Spaces{
				Insecure: true,
			},
		}
		cfg.OCDav = &ocdav.Config{
			Insecure: true,
		}
		cfg.Proxy = &proxy.Config{
			InsecureBackends: true,
		}

		cfg.StorageMetadata = &storagemetadata.Config{
			DataProviderInsecure: true,
		}
		cfg.StorageUsers = &storageusers.Config{
			DataProviderInsecure: true,
		}
		cfg.Thumbnails = &thumbnails.Config{
			Thumbnail: thumbnails.Thumbnail{
				WebdavAllowInsecure: true,
				CS3AllowInsecure:    true,
			},
		}

	}

	idmServicePassword, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for idm: %s", err)
	}
	idpServicePassword, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for idp: %s", err)
	}
	ocisAdminServicePassword, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for ocis admin: %s", err)
	}
	revaServicePassword, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for reva: %s", err)
	}
	tokenManagerJwtSecret, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for tokenmanager: %s", err)
	}
	machineAuthApiKey, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for machineauthsecret: %s", err)
	}
	revaTransferSecret, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for machineauthsecret: %s", err)
	}

	cfg.MachineAuthAPIKey = machineAuthApiKey
	cfg.TransferSecret = revaTransferSecret
	cfg.TokenManager.JWTSecret = tokenManagerJwtSecret

	cfg.IDM.ServiceUserPasswords.Idm = idmServicePassword
	cfg.Graph.Identity.LDAP.BindPassword = idmServicePassword

	cfg.IDM.ServiceUserPasswords.Idp = idpServicePassword
	cfg.IDP.Ldap.BindPassword = idpServicePassword

	cfg.IDM.ServiceUserPasswords.Reva = revaServicePassword
	cfg.AuthBasic.AuthProviders.LDAP.BindPassword = revaServicePassword
	cfg.Group.Drivers.LDAP.BindPassword = revaServicePassword
	cfg.User.Drivers.LDAP.BindPassword = revaServicePassword

	cfg.IDM.ServiceUserPasswords.OcisAdmin = ocisAdminServicePassword

	yamlOutput, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not marshall config into yaml: %s", err)
	}
	targetPath := path.Join(configPath, configFilename)
	err = ioutil.WriteFile(targetPath, yamlOutput, 0600)
	if err != nil {
		return err
	}
	fmt.Printf(
		"======================================\n"+
			" generated OCIS Config\n"+
			"======================================\n"+
			" configpath : %s\n"+
			" user       : admin\n"+
			" password   : %s\n",
		targetPath, ocisAdminServicePassword)
	return nil
}

func stringPrompt(label string) string {
	input := ""
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		input, _ = reader.ReadString('\n')
		if input != "" {
			break
		}
	}
	return strings.TrimSpace(input)
}
