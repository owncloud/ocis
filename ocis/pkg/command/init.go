package command

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/generators"
	"github.com/owncloud/ocis/ocis/pkg/register"
	cli "github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const configFilename string = "ocis.yaml" // TODO: use also a constant for reading this file
const passwordLength int = 32

type tokenManager struct {
	JWT_Secret string
}

type insecureExtension struct {
	Insecure bool
}

type insecureProxyExtension struct {
	Insecure_backends bool
}

type dataProviderInsecureSettings struct {
	Data_provider_insecure bool
}

type ldapSettings struct {
	Bind_password string
}
type ldapBasedExtension struct {
	Ldap ldapSettings
}

type graphExtension struct {
	Spaces   insecureExtension
	Identity ldapBasedExtension
}

type serviceUserPasswordsSettings struct {
	Admin_password string
	Idm_password   string
	Reva_password  string
	Idp_password   string
}
type idmExtension struct {
	Service_user_Passwords serviceUserPasswordsSettings
}

type frontendExtension struct {
	Archiver     insecureExtension
	App_provider insecureExtension
}

type authbasicExtension struct {
	Auth_providers ldapBasedExtension
}

type authProviderSettings struct {
	Oidc insecureExtension
}
type authbearerExtension struct {
	Auth_providers authProviderSettings
}

type userAndGroupExtension struct {
	Drivers ldapBasedExtension
}

type thumbnailSettings struct {
	Webdav_allow_insecure bool
	Cs3_allow_insecure    bool
}

type thumbNailExtension struct {
	Thumbnail thumbnailSettings
}

type ocisConfig struct {
	Token_manager        tokenManager
	Machine_auth_api_key string
	Transfer_secret      string
	Graph                graphExtension
	Idp                  ldapBasedExtension
	Idm                  idmExtension
	Proxy                insecureProxyExtension
	Frontend             frontendExtension
	Auth_basic           authbasicExtension
	Auth_bearer          authbearerExtension
	User                 userAndGroupExtension
	Group                userAndGroupExtension
	Storage_metadata     dataProviderInsecureSettings
	Storage_users        dataProviderInsecureSettings
	Ocdav                insecureExtension
	Thumbnails           thumbNailExtension
}

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

func backupOcisConfigFile(configPath string) (string, error) {
	sourceConfig := path.Join(configPath, configFilename)
	targetBackupConfig := path.Join(configPath, configFilename+"."+time.Now().Format("2006-01-02-15-04-05")+".backup")
	source, err := os.Open(sourceConfig)
	if err != nil {
		log.Fatalf("Could not read %s (%s)", sourceConfig, err)
	}
	defer source.Close()
	target, err := os.Create(targetBackupConfig)
	if err != nil {
		log.Fatalf("Could not generate backup %s (%s)", targetBackupConfig, err)
	}
	defer target.Close()
	_, err = io.Copy(target, source)
	if err != nil {
		log.Fatalf("Could not write backup %s (%s)", targetBackupConfig, err)
	}
	return targetBackupConfig, nil
}

func createConfig(insecure, forceOverwrite bool, configPath string) error {
	err := checkConfigPath(configPath)
	targetBackupConfig := ""
	if err != nil && !forceOverwrite {
		return err
	} else if forceOverwrite {
		targetBackupConfig, err = backupOcisConfigFile(configPath)
		if err != nil {
			return err
		} else {

		}
	}
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return err
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

	cfg := ocisConfig{
		Token_manager: tokenManager{
			JWT_Secret: tokenManagerJwtSecret,
		},
		Machine_auth_api_key: machineAuthApiKey,
		Transfer_secret:      revaTransferSecret,
		Idm: idmExtension{
			Service_user_Passwords: serviceUserPasswordsSettings{
				Admin_password: ocisAdminServicePassword,
				Idp_password:   idpServicePassword,
				Reva_password:  revaServicePassword,
				Idm_password:   idmServicePassword,
			},
		},
		Idp: ldapBasedExtension{
			Ldap: ldapSettings{
				Bind_password: idpServicePassword,
			},
		},
		Auth_basic: authbasicExtension{
			Auth_providers: ldapBasedExtension{
				Ldap: ldapSettings{
					Bind_password: revaServicePassword,
				},
			},
		},
		Group: userAndGroupExtension{
			Drivers: ldapBasedExtension{
				Ldap: ldapSettings{
					Bind_password: revaServicePassword,
				},
			},
		},
		User: userAndGroupExtension{
			Drivers: ldapBasedExtension{
				Ldap: ldapSettings{
					Bind_password: revaServicePassword,
				},
			},
		},
		Graph: graphExtension{
			Identity: ldapBasedExtension{
				Ldap: ldapSettings{
					Bind_password: idmServicePassword,
				},
			},
		},
	}

	if insecure {
		cfg.Auth_bearer = authbearerExtension{
			Auth_providers: authProviderSettings{
				Oidc: insecureExtension{
					Insecure: true,
				},
			},
		}
		cfg.Frontend = frontendExtension{
			App_provider: insecureExtension{
				Insecure: true,
			},
			Archiver: insecureExtension{
				Insecure: true,
			},
		}
		cfg.Graph.Spaces = insecureExtension{
			Insecure: true,
		}
		cfg.Ocdav = insecureExtension{
			Insecure: true,
		}
		cfg.Proxy = insecureProxyExtension{
			Insecure_backends: true,
		}
		cfg.Storage_metadata = dataProviderInsecureSettings{
			Data_provider_insecure: true,
		}
		cfg.Storage_users = dataProviderInsecureSettings{
			Data_provider_insecure: true,
		}
		cfg.Thumbnails = thumbNailExtension{
			Thumbnail: thumbnailSettings{
				Webdav_allow_insecure: true,
				Cs3_allow_insecure:    true,
			},
		}
	}

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
		"\n\n=========================================\n"+
			" generated OCIS Config\n"+
			"=========================================\n"+
			" configpath : %s\n"+
			" user       : admin\n"+
			" password   : %s\n\n",
		targetPath, ocisAdminServicePassword)
	if targetBackupConfig != "" {
		fmt.Printf("\n=========================================\n"+
			"An older config file has been backuped to\n %s\n\n",
			targetBackupConfig)
	}
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
