package init

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/generators"
	"gopkg.in/yaml.v2"
)

const configFilename string = "ocis.yaml" // TODO: use also a constant for reading this file
const passwordLength int = 32

type TokenManager struct {
	JWT_Secret string
}

type InsecureExtension struct {
	Insecure bool
}

type InsecureProxyExtension struct {
	Insecure_backends bool
}

type DataProviderInsecureSettings struct {
	Data_provider_insecure bool
}

type LdapSettings struct {
	Bindpassword string
}
type LdapBasedExtension struct {
	Ldap LdapSettings
}

type GraphExtension struct {
	Spaces   InsecureExtension
	Identity LdapBasedExtension
}

type ServiceUserPasswordsSettings struct {
	Admin_password string
	Idm_password   string
	Reva_password  string
	Idp_password   string
}
type IdmExtension struct {
	Service_user_Passwords ServiceUserPasswordsSettings
}

type FrontendExtension struct {
	Archiver     InsecureExtension
	App_provider InsecureExtension
}

type AuthbasicExtension struct {
	Auth_providers LdapBasedExtension
}

type AuthProviderSettings struct {
	Oidc InsecureExtension
}
type AuthbearerExtension struct {
	Auth_providers AuthProviderSettings
}

type UserAndGroupExtension struct {
	Drivers LdapBasedExtension
}

type ThumbnailSettings struct {
	Webdav_allow_insecure bool
	Cs3_allow_insecure    bool
}

type ThumbNailExtension struct {
	Thumbnail ThumbnailSettings
}

type OcisConfig struct {
	Token_manager        TokenManager
	Machine_auth_api_key string
	Transfer_secret      string
	Graph                GraphExtension
	Idp                  LdapBasedExtension
	Idm                  IdmExtension
	Proxy                InsecureProxyExtension
	Frontend             FrontendExtension
	Auth_basic           AuthbasicExtension
	Auth_bearer          AuthbearerExtension
	User                 UserAndGroupExtension
	Group                UserAndGroupExtension
	Storage_metadata     DataProviderInsecureSettings
	Storage_users        DataProviderInsecureSettings
	Ocdav                InsecureExtension
	Thumbnails           ThumbNailExtension
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

func CreateConfig(insecure, forceOverwrite bool, configPath string) error {
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

	cfg := OcisConfig{
		Token_manager: TokenManager{
			JWT_Secret: tokenManagerJwtSecret,
		},
		Machine_auth_api_key: machineAuthApiKey,
		Transfer_secret:      revaTransferSecret,
		Idm: IdmExtension{
			Service_user_Passwords: ServiceUserPasswordsSettings{
				Admin_password: ocisAdminServicePassword,
				Idp_password:   idpServicePassword,
				Reva_password:  revaServicePassword,
				Idm_password:   idmServicePassword,
			},
		},
		Idp: LdapBasedExtension{
			Ldap: LdapSettings{
				Bindpassword: idpServicePassword,
			},
		},
		Auth_basic: AuthbasicExtension{
			Auth_providers: LdapBasedExtension{
				Ldap: LdapSettings{
					Bindpassword: revaServicePassword,
				},
			},
		},
		Group: UserAndGroupExtension{
			Drivers: LdapBasedExtension{
				Ldap: LdapSettings{
					Bindpassword: revaServicePassword,
				},
			},
		},
		User: UserAndGroupExtension{
			Drivers: LdapBasedExtension{
				Ldap: LdapSettings{
					Bindpassword: revaServicePassword,
				},
			},
		},
		Graph: GraphExtension{
			Identity: LdapBasedExtension{
				Ldap: LdapSettings{
					Bindpassword: idmServicePassword,
				},
			},
		},
	}

	if insecure {
		cfg.Auth_bearer = AuthbearerExtension{
			Auth_providers: AuthProviderSettings{
				Oidc: InsecureExtension{
					Insecure: true,
				},
			},
		}
		cfg.Frontend = FrontendExtension{
			App_provider: InsecureExtension{
				Insecure: true,
			},
			Archiver: InsecureExtension{
				Insecure: true,
			},
		}
		cfg.Graph.Spaces = InsecureExtension{
			Insecure: true,
		}
		cfg.Ocdav = InsecureExtension{
			Insecure: true,
		}
		cfg.Proxy = InsecureProxyExtension{
			Insecure_backends: true,
		}
		cfg.Storage_metadata = DataProviderInsecureSettings{
			Data_provider_insecure: true,
		}
		cfg.Storage_users = DataProviderInsecureSettings{
			Data_provider_insecure: true,
		}
		cfg.Thumbnails = ThumbNailExtension{
			Thumbnail: ThumbnailSettings{
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
