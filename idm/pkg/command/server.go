package command

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/go-ldap/ldif"
	"github.com/libregraph/idm/pkg/ldappassword"
	"github.com/libregraph/idm/pkg/ldbbolt"
	"github.com/libregraph/idm/server"
	"github.com/owncloud/ocis/idm/pkg/config"
	"github.com/owncloud/ocis/idm/pkg/config/parser"
	"github.com/owncloud/ocis/idm/pkg/logging"
	pkgcrypto "github.com/owncloud/ocis/ocis-pkg/crypto"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return parser.ParseConfig(cfg)
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
			start(ctx, logger, cfg)
			return nil
		},
	}
}

func start(ctx context.Context, logger log.Logger, cfg *config.Config) error {
	servercfg := server.Config{
		Logger:          log.LogrusWrap(logger.Logger),
		LDAPHandler:     "boltdb",
		LDAPSListenAddr: cfg.IDM.LDAPSAddr,
		TLSCertFile:     cfg.IDM.Cert,
		TLSKeyFile:      cfg.IDM.Key,
		LDAPBaseDN:      "o=libregraph-idm",
		LDAPAdminDN:     "uid=libregrah,o=libregraph-idm",

		BoltDBFile: cfg.IDM.DatabasePath,
	}

	if cfg.IDM.LDAPSAddr != "" {
		// Generate a self-signing cert if no certificate is present
		if err := pkgcrypto.GenCert(cfg.IDM.Cert, cfg.IDM.Key, logger); err != nil {
			logger.Fatal().Err(err).Msgf("Could not generate test-certificate")
		}
	}
	if _, err := os.Stat(servercfg.BoltDBFile); errors.Is(err, os.ErrNotExist) {
		logger.Debug().Msg("Bootstrapping IDM database")
		err = bootstrap(logger, cfg, servercfg)
		logger.Error().Err(err).Msg("failed")
	}

	svc, err := server.NewServer(&servercfg)
	if err != nil {
		return err
	}
	return svc.Serve(ctx)
}

func bootstrap(logger log.Logger, cfg *config.Config, srvcfg server.Config) error {
	// Hash password if the config does not supply a hash already
	var pwhash string
	var err error
	if strings.HasPrefix(cfg.IDM.AdminPassword, "$argon2id$") {
		// password is alread hashed
		pwhash = "{ARGON2}" + cfg.IDM.AdminPassword
	} else {
		if pwhash, err = ldappassword.Hash(cfg.IDM.AdminPassword, "{ARGON2}"); err != nil {
			return err
		}
	}

	bdb := &ldbbolt.LdbBolt{}

	if err := bdb.Configure(srvcfg.Logger, srvcfg.LDAPBaseDN, srvcfg.BoltDBFile, nil); err != nil {
		return err
	}
	defer bdb.Close()

	if err := bdb.Initialize(); err != nil {
		return err
	}

	// Prepare the initial Data from template. To be able to set the
	// supplied admin password
	tmpl, err := template.New("baseldif").Parse(baseldif)
	if err != nil {
		return err
	}

	var tmplWriter strings.Builder
	// We need to treat the hash as binary in the LDIF template to avoid
	// go-ldap/ldif to to any fancy escaping
	b64 := base64.StdEncoding.EncodeToString([]byte(pwhash))
	err = tmpl.Execute(&tmplWriter, b64)
	if err != nil {
		return err
	}

	s := strings.NewReader(tmplWriter.String())
	lf := &ldif.LDIF{}
	err = ldif.Unmarshal(s, lf)
	if err != nil {
		return err
	}

	for _, entry := range lf.AllEntries() {
		logger.Debug().Str("dn", entry.DN).Msg("Adding entry")
		if err := bdb.EntryPut(entry); err != nil {
			return fmt.Errorf("error adding Entry '%s': %w", entry.DN, err)
		}
	}
	return nil
}

var baseldif string = `dn: o=libregraph-idm
o: libregraph-idm
objectClass: organization

dn: ou=users,o=libregraph-idm
objectClass: organizationalUnit
ou: users

dn: ou=groups,o=libregraph-idm
objectClass: organizationalUnit
ou: groups

dn: uid=libregraph,o=libregraph-idm
objectClass: account
objectClass: simpleSecurityObject
uid: libregraph
userPassword:: {{.}}`
