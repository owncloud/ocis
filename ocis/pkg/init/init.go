package init

import (
	"fmt"
	"github.com/gofrs/uuid"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"gopkg.in/yaml.v2"
)

const (
	configFilename = "ocis.yaml" // TODO: use also a constant for reading this file
	passwordLength = 32
)

var (
	_insecureService = InsecureService{Insecure: true}
	_insecureEvents  = Events{TLSInsecure: true}
)

func checkConfigPath(configPath string) error {
	targetPath := path.Join(configPath, configFilename)
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("config in %s already exists", targetPath)
	}
	return nil
}

func configExists(configPath string) bool {
	targetPath := path.Join(configPath, configFilename)
	if _, err := os.Stat(targetPath); err == nil {
		return true
	}
	return false
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

// CreateConfig creates a config file with random passwords at configPath
func CreateConfig(insecure, forceOverwrite, diff bool, configPath, adminPassword string) error {
	if diff {
		if forceOverwrite {
			return fmt.Errorf("diff and force-overwrite flags are mutually exclusive")
		}
		if adminPassword != "" {
			return fmt.Errorf("diff and admin-password flags are mutually exclusive")
		}
	}

	if configExists(configPath) {
		if !forceOverwrite && !diff {
			return fmt.Errorf("config file already exists, use --force-overwrite to overwrite or --diff to show diff")
		}
	}

	err := checkConfigPath(configPath)
	if err != nil && (!forceOverwrite && !diff) {
		fmt.Println("off")
		return err
	}
	targetBackupConfig := ""
	if err != nil {
		targetBackupConfig, err = backupOcisConfigFile(configPath)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return err
	}

	// Load old config
	var oldCfg OcisConfig
	if diff {
		fp, err := os.ReadFile(path.Join(configPath, configFilename))
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(fp, &oldCfg)
	}

	var (
		systemUserID, adminUserID, graphApplicationID, storageUsersMountID, serviceAccountID  string
		idmServicePassword, idpServicePassword, ocisAdminServicePassword, revaServicePassword string
		tokenManagerJwtSecret, collaborationWOPISecret, machineAuthAPIKey, systemUserAPIKey   string
		revaTransferSecret, thumbnailsTransferSecret, serviceAccountSecret                    string
	)

	if diff {
		systemUserID = oldCfg.SystemUserID
		adminUserID = oldCfg.AdminUserID
		graphApplicationID = oldCfg.Graph.Application.ID
		storageUsersMountID = oldCfg.Gateway.StorageRegistry.StorageUsersMountID
		serviceAccountID = oldCfg.Graph.ServiceAccount.ServiceAccountID

		idmServicePassword = oldCfg.Idm.ServiceUserPasswords.IdmPassword
		idpServicePassword = oldCfg.Idm.ServiceUserPasswords.IdpPassword
		ocisAdminServicePassword = oldCfg.Idm.ServiceUserPasswords.AdminPassword
		revaServicePassword = oldCfg.Idm.ServiceUserPasswords.RevaPassword
		tokenManagerJwtSecret = oldCfg.TokenManager.JWTSecret
		collaborationWOPISecret = oldCfg.Collaboration.WopiApp.Secret
		machineAuthAPIKey = oldCfg.MachineAuthAPIKey
		systemUserAPIKey = oldCfg.SystemUserAPIKey
		revaTransferSecret = oldCfg.TransferSecret
		thumbnailsTransferSecret = oldCfg.Thumbnails.Thumbnail.TransferSecret
		serviceAccountSecret = oldCfg.Graph.ServiceAccount.ServiceAccountSecret
	} else {
		systemUserID = uuid.Must(uuid.NewV4()).String()
		adminUserID = uuid.Must(uuid.NewV4()).String()
		graphApplicationID = uuid.Must(uuid.NewV4()).String()
		storageUsersMountID = uuid.Must(uuid.NewV4()).String()
		serviceAccountID = uuid.Must(uuid.NewV4()).String()

		idmServicePassword, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for idm: %s", err)
		}
		idpServicePassword, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for idp: %s", err)
		}
		ocisAdminServicePassword = adminPassword
		if ocisAdminServicePassword == "" {
			ocisAdminServicePassword, err = generators.GenerateRandomPassword(passwordLength)
			if err != nil {
				return fmt.Errorf("could not generate random password for ocis admin: %s", err)
			}
		}

		revaServicePassword, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for reva: %s", err)
		}
		tokenManagerJwtSecret, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for tokenmanager: %s", err)
		}
		collaborationWOPISecret, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random wopi secret for collaboration service: %s", err)
		}
		machineAuthAPIKey, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for machineauthsecret: %s", err)
		}
		systemUserAPIKey, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random system user API key: %s", err)
		}
		revaTransferSecret, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for revaTransferSecret: %s", err)
		}
		thumbnailsTransferSecret, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for thumbnailsTransferSecret: %s", err)
		}
		serviceAccountSecret, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for thumbnailsTransferSecret: %s", err)
		}
	}

	serviceAccount := ServiceAccount{
		ServiceAccountID:     serviceAccountID,
		ServiceAccountSecret: serviceAccountSecret,
	}

	cfg := OcisConfig{
		TokenManager: TokenManager{
			JWTSecret: tokenManagerJwtSecret,
		},
		MachineAuthAPIKey: machineAuthAPIKey,
		SystemUserAPIKey:  systemUserAPIKey,
		TransferSecret:    revaTransferSecret,
		SystemUserID:      systemUserID,
		AdminUserID:       adminUserID,
		Idm: IdmService{
			ServiceUserPasswords: ServiceUserPasswordsSettings{
				AdminPassword: ocisAdminServicePassword,
				IdpPassword:   idpServicePassword,
				RevaPassword:  revaServicePassword,
				IdmPassword:   idmServicePassword,
			},
		},
		Idp: LdapBasedService{
			Ldap: LdapSettings{
				BindPassword: idpServicePassword,
			},
		},
		AuthBasic: AuthbasicService{
			AuthProviders: LdapBasedService{
				Ldap: LdapSettings{
					BindPassword: revaServicePassword,
				},
			},
		},
		Collaboration: Collaboration{
			WopiApp: WopiApp{
				Secret: collaborationWOPISecret,
			},
		},
		Groups: UsersAndGroupsService{
			Drivers: LdapBasedService{
				Ldap: LdapSettings{
					BindPassword: revaServicePassword,
				},
			},
		},
		Users: UsersAndGroupsService{
			Drivers: LdapBasedService{
				Ldap: LdapSettings{
					BindPassword: revaServicePassword,
				},
			},
		},
		Graph: GraphService{
			Application: GraphApplication{
				ID: graphApplicationID,
			},
			Identity: LdapBasedService{
				Ldap: LdapSettings{
					BindPassword: idmServicePassword,
				},
			},
			ServiceAccount: serviceAccount,
		},
		Thumbnails: ThumbnailService{
			Thumbnail: ThumbnailSettings{
				TransferSecret: thumbnailsTransferSecret,
			},
		},
		Gateway: Gateway{
			StorageRegistry: StorageRegistry{
				StorageUsersMountID: storageUsersMountID,
			},
		},
		StorageUsers: StorageUsers{
			MountID:        storageUsersMountID,
			ServiceAccount: serviceAccount,
		},
		Userlog: Userlog{
			ServiceAccount: serviceAccount,
		},
		AuthService: AuthService{
			ServiceAccount: serviceAccount,
		},
		Search: Search{
			ServiceAccount: serviceAccount,
		},
		Notifications: Notifications{
			ServiceAccount: serviceAccount,
		},
		Frontend: FrontendService{
			ServiceAccount: serviceAccount,
		},
		Ocm: OcmService{
			ServiceAccount: serviceAccount,
		},
		Clientlog: Clientlog{
			ServiceAccount: serviceAccount,
		},
		Proxy: ProxyService{
			ServiceAccount: serviceAccount,
		},
		Settings: SettingsService{
			ServiceAccountIDs: []string{serviceAccount.ServiceAccountID},
		},
		Activitylog: Activitylog{
			ServiceAccount: serviceAccount,
		},
	}

	if insecure {
		cfg.AuthBearer = AuthbearerService{
			AuthProviders: AuthProviderSettings{Oidc: _insecureService},
		}
		cfg.Collaboration.App.Insecure = true
		cfg.Frontend.AppHandler = _insecureService
		cfg.Frontend.Archiver = _insecureService
		cfg.Graph.Spaces = _insecureService
		cfg.Graph.Events = _insecureEvents
		cfg.Notifications.Notifications.Events = _insecureEvents
		cfg.Search.Events = _insecureEvents
		cfg.Audit.Events = _insecureEvents
		cfg.Sharing.Events = _insecureEvents
		cfg.StorageUsers.Events = _insecureEvents
		cfg.Nats.Nats.TLSSkipVerifyClientCert = true
		cfg.Ocdav = _insecureService
		cfg.Proxy = ProxyService{
			InsecureBackends: true,
			OIDC: InsecureProxyOIDC{
				Insecure: true,
			},
			ServiceAccount: serviceAccount,
		}

		cfg.Thumbnails.Thumbnail.WebdavAllowInsecure = true
		cfg.Thumbnails.Thumbnail.Cs3AllowInsecure = true
	}
	yamlOutput, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not marshall config into yaml: %s", err)
	}
	if diff {
		fmt.Println("running in diff mode")
		tmpFile := path.Join(configPath, "ocis.yaml.tmp")
		err = os.WriteFile(tmpFile, yamlOutput, 0600)
		if err != nil {
			return err
		}
		fmt.Println("diff -u " + path.Join(configPath, configFilename) + " " + tmpFile)
		cmd := exec.Command("diff", "-u", path.Join(configPath, configFilename), tmpFile)
		stdout, _ := cmd.Output()
		fmt.Println(string(stdout))
		err = os.Remove(tmpFile)
		patchPath := path.Join(configPath, "ocis.config.patch")
		err = os.WriteFile(patchPath, stdout, 0600)
		if err != nil {
			return err
		}
		fmt.Printf("diff written to %s\n", patchPath)
	} else {
		targetPath := path.Join(configPath, configFilename)
		err = os.WriteFile(targetPath, yamlOutput, 0600)
		if err != nil {
			return err
		}
		printBanner(targetPath, ocisAdminServicePassword, targetBackupConfig)
	}
	return nil
}

func printBanner(targetPath, ocisAdminServicePassword, targetBackupConfig string) {
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
}
