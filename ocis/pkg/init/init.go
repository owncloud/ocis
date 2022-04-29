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
	JWTSecret string `yaml:"jwt_secret"`
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
	Bind_password string
}
type LdapBasedExtension struct {
	Ldap LdapSettings
}

type GraphExtension struct {
	Spaces   InsecureExtension
	Identity LdapBasedExtension
}

type ServiceUserPasswordsSettings struct {
	AdminPassword string `yaml:"admin_password"`
	IdmPassword   string `yaml:"idm_password"`
	RevaPassword  string `yaml:"reva_password"`
	IdpPassword   string `yaml:"idp_password"`
}
type IdmExtension struct {
	ServiceUserPasswords ServiceUserPasswordsSettings `yaml:"service_user_passwords"`
}

type FrontendExtension struct {
	Archiver    InsecureExtension
	AppProvider InsecureExtension `yaml:"app_provider"`
}

type AuthbasicExtension struct {
	AuthProviders LdapBasedExtension `yaml:"auth_providers"`
}

type AuthProviderSettings struct {
	Oidc InsecureExtension
}
type AuthbearerExtension struct {
	AuthProviders AuthProviderSettings `yaml:"auth_providers"`
}

type UserAndGroupExtension struct {
	Drivers LdapBasedExtension
}

type ThumbnailSettings struct {
	WebdavAllowInsecure bool `yaml:"webdav_allow_insecure"`
	Cs3AllowInsecure    bool `yaml:"cs3_allow_insecure"`
}

type ThumbNailExtension struct {
	Thumbnail ThumbnailSettings
}

// TODO: use the oCIS config struct instead of this custom struct
// We can't use it right now, because it would need  "omitempty" on
// all elements, in order to produce a slim config file with `ocis init`.
// We can't just add these "omitempty" tags, since we want to generate
// full example configuration files with that struct, too.
// Proposed solution to  get rid of this temporary solution:
// - use the oCIS config struct
// - set the needed values like below
// - marshal it to yaml
// - unmarshal it into yaml.Node
// - recurse through the nodes and delete empty / default ones
// - marshal it to yaml
type OcisConfig struct {
	TokenManager      TokenManager `yaml:"token_manager"`
	MachineAuthApiKey string       `yaml:"machine_auth_api_key"`
	TransferSecret    string       `yaml:"transfer_secret"`
	Graph             GraphExtension
	Idp               LdapBasedExtension
	Idm               IdmExtension
	Proxy             InsecureProxyExtension
	Frontend          FrontendExtension
	AuthBasic         AuthbasicExtension  `yaml:"auth_basic"`
	AuthBearer        AuthbearerExtension `yaml:"auth_bearer"`
	User              UserAndGroupExtension
	Group             UserAndGroupExtension
	StorageMetadata   DataProviderInsecureSettings `yaml:"storage_metadata"`
	StorageUsers      DataProviderInsecureSettings `yaml:"storage_users"`
	Ocdav             InsecureExtension
	Thumbnails        ThumbNailExtension
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

func CreateConfig(insecure, forceOverwrite bool, configPath, adminPassword string) error {
	targetBackupConfig := ""

	err := checkConfigPath(configPath)
	if err != nil && !forceOverwrite {
		return err
	} else if forceOverwrite && err != nil {
		targetBackupConfig, err = backupOcisConfigFile(configPath)
		if err != nil {
			return err
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
	ocisAdminServicePassword := adminPassword
	if ocisAdminServicePassword == "" {
		ocisAdminServicePassword, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for ocis admin: %s", err)
		}
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
		TokenManager: TokenManager{
			JWTSecret: tokenManagerJwtSecret,
		},
		MachineAuthApiKey: machineAuthApiKey,
		TransferSecret:    revaTransferSecret,
		Idm: IdmExtension{
			ServiceUserPasswords: ServiceUserPasswordsSettings{
				AdminPassword: ocisAdminServicePassword,
				IdpPassword:   idpServicePassword,
				RevaPassword:  revaServicePassword,
				IdmPassword:   idmServicePassword,
			},
		},
		Idp: LdapBasedExtension{
			Ldap: LdapSettings{
				Bind_password: idpServicePassword,
			},
		},
		AuthBasic: AuthbasicExtension{
			AuthProviders: LdapBasedExtension{
				Ldap: LdapSettings{
					Bind_password: revaServicePassword,
				},
			},
		},
		Group: UserAndGroupExtension{
			Drivers: LdapBasedExtension{
				Ldap: LdapSettings{
					Bind_password: revaServicePassword,
				},
			},
		},
		User: UserAndGroupExtension{
			Drivers: LdapBasedExtension{
				Ldap: LdapSettings{
					Bind_password: revaServicePassword,
				},
			},
		},
		Graph: GraphExtension{
			Identity: LdapBasedExtension{
				Ldap: LdapSettings{
					Bind_password: idmServicePassword,
				},
			},
		},
	}

	if insecure {
		cfg.AuthBearer = AuthbearerExtension{
			AuthProviders: AuthProviderSettings{
				Oidc: InsecureExtension{
					Insecure: true,
				},
			},
		}
		cfg.Frontend = FrontendExtension{
			AppProvider: InsecureExtension{
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
		cfg.StorageMetadata = DataProviderInsecureSettings{
			Data_provider_insecure: true,
		}
		cfg.StorageUsers = DataProviderInsecureSettings{
			Data_provider_insecure: true,
		}
		cfg.Thumbnails = ThumbNailExtension{
			Thumbnail: ThumbnailSettings{
				WebdavAllowInsecure: true,
				Cs3AllowInsecure:    true,
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
		"\n=========================================\n"+
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
