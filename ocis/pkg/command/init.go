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

	idm "github.com/owncloud/ocis/extensions/idm/pkg/config"
	proxy "github.com/owncloud/ocis/extensions/proxy/pkg/config"
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
		//Accounts: &accounts.Config{},
		//Audit:    &audit.Config{},
		//GLAuth:        &glauth.Config{},
		//GraphExplorer: &graphExplorer.Config{},
		//Graph: &graph.Config{},
		IDM: &idm.Config{},
		//IDP:           &idp.Config{},
		//Nats:          &nats.Config{},
		//Notifications: &notifications.Config{},
		//Proxy:         &proxy.Config{},
		//OCS:        &ocs.Config{},
		//Settings:   &settings.Config{},
		// TODO: fix storage
		//AuthBasic:         &authbasic.Config{},
		//AuthBearer:        &authbearer.Config{},
		//AppProvider:       &appprovider.Config{},
		//AuthMachine:       &authmachine.Config{},
		//Gateway:           &gateway.Config{},
		//Group:             &group.Config{},
		//Sharing:           &sharing.Config{},
		//StorageMetadata:   &storagemetadata.Config{},
		//StorageUsers:      &storageusers.Config{},
		//StorageShares:     &storageshares.Config{},
		//StoragePublicLink: &storagepublic.Config{},
		//User:              &user.Config{},
		//OCDav:             &ocdav.Config{},
		//Thumbnails: &thumbnails.Config{},
		//Web:           &web.Config{},
		//WebDAV:        &webdav.Config{},
	}

	if insecure {
		cfg.Proxy = &proxy.Config{}
		cfg.Proxy.InsecureBackends = insecure
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

	// TODO: IDP config is missing (LDAP + GROUP provider)
	// TODO: REVA config is missing (LDAP + GROUP provider)
	// TODO: graph needs IDM password configured
	// TODO: add missing insecure occurences

	cfg.MachineAuthAPIKey = machineAuthApiKey
	cfg.TransferSecret = revaTransferSecret
	cfg.TokenManager.JWTSecret = tokenManagerJwtSecret
	//cfg.Commons.TokenManager.JWTSecret = tokenManagerJwtSecret
	//cfg.Accounts.TokenManager.JWTSecret = tokenManagerJwtSecret
	//cfg.Graph.TokenManager.JWTSecret = tokenManagerJwtSecret
	//fmt.Printf("%v\n", cfg.Graph.TokenManager)
	cfg.IDM.ServiceUserPasswords.Idm = idmServicePassword
	cfg.IDM.ServiceUserPasswords.Idp = idpServicePassword
	cfg.IDM.ServiceUserPasswords.OcisAdmin = ocisAdminServicePassword
	cfg.IDM.ServiceUserPasswords.Reva = revaServicePassword
	//cfg.Notifications.Notifications.MachineAuthAPIKey = machineAuthSecret
	//cfg.OCS.MachineAuthAPIKey = machineAuthSecret
	//cfg.Proxy.TokenManager.JWTSecret = tokenManagerJwtSecret
	//fmt.Printf("%v\n", cfg.Proxy.TokenManager)
	//cfg.Proxy.MachineAuthAPIKey = machineAuthSecret
	//cfg.Settings.Metadata.MachineAuthAPIKey = machineAuthSecret
	//cfg.Settings.TokenManager.JWTSecret = tokenManagerJwtSecret

	//TODO: move all jwt secrets to shared.common
	//cfg.AppProvider.JWTSecret = tokenManagerJwtSecret
	//cfg.AuthBasic.JWTSecret = tokenManagerJwtSecret
	//cfg.AuthBearer.JWTSecret = tokenManagerJwtSecret
	//cfg.AuthMachine.JWTSecret = tokenManagerJwtSecret
	//cfg.Gateway.JWTSecret = tokenManagerJwtSecret
	//cfg.Group.JWTSecret = tokenManagerJwtSecret
	//cfg.Sharing.JWTSecret = tokenManagerJwtSecret
	//cfg.StorageMetadata.JWTSecret = tokenManagerJwtSecret
	//cfg.StoragePublicLink.JWTSecret = tokenManagerJwtSecret
	//cfg.StorageShares.JWTSecret = tokenManagerJwtSecret
	//cfg.StorageUsers.JWTSecret = tokenManagerJwtSecret
	//cfg.User.JWTSecret = tokenManagerJwtSecret
	//cfg.OCDav.JWTSecret = tokenManagerJwtSecret

	//TODO: following line is defunc, figure out why
	//cfg.Gateway.MachineAuthAPIKey = machineAuthApiKey

	//cfg.Thumbnails.Thumbnail.TransferSecret = revaTransferTokenSecret
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
