package command

import (
	"errors"
	"os"
	"time"

	"github.com/owncloud/ocis/extensions/storage/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

const caTimeout = 5

func ldapConfigFromString(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"uri":               cfg.Reva.LDAP.URI,
		"cacert":            cfg.Reva.LDAP.CACert,
		"insecure":          cfg.Reva.LDAP.Insecure,
		"bind_username":     cfg.Reva.LDAP.BindDN,
		"bind_password":     cfg.Reva.LDAP.BindPassword,
		"user_base_dn":      cfg.Reva.LDAP.UserBaseDN,
		"group_base_dn":     cfg.Reva.LDAP.GroupBaseDN,
		"user_filter":       cfg.Reva.LDAP.UserFilter,
		"group_filter":      cfg.Reva.LDAP.GroupFilter,
		"user_objectclass":  cfg.Reva.LDAP.UserObjectClass,
		"group_objectclass": cfg.Reva.LDAP.GroupObjectClass,
		"login_attributes":  cfg.Reva.LDAP.LoginAttributes,
		"idp":               cfg.Reva.LDAP.IDP,
		"gatewaysvc":        cfg.Reva.Gateway.Endpoint,
		"user_schema": map[string]interface{}{
			"id":              cfg.Reva.LDAP.UserSchema.ID,
			"idIsOctetString": cfg.Reva.LDAP.UserSchema.IDIsOctetString,
			"mail":            cfg.Reva.LDAP.UserSchema.Mail,
			"displayName":     cfg.Reva.LDAP.UserSchema.DisplayName,
			"userName":        cfg.Reva.LDAP.UserSchema.Username,
		},
		"group_schema": map[string]interface{}{
			"id":              cfg.Reva.LDAP.GroupSchema.ID,
			"idIsOctetString": cfg.Reva.LDAP.GroupSchema.IDIsOctetString,
			"mail":            cfg.Reva.LDAP.GroupSchema.Mail,
			"displayName":     cfg.Reva.LDAP.GroupSchema.DisplayName,
			"groupName":       cfg.Reva.LDAP.GroupSchema.Groupname,
			"member":          cfg.Reva.LDAP.GroupSchema.Member,
		},
	}
}

func waitForLDAPCA(log log.Logger, cfg *config.LDAP) error {
	if !cfg.Insecure && cfg.CACert != "" {
		if _, err := os.Stat(cfg.CACert); errors.Is(err, os.ErrNotExist) {
			log.Warn().Str("LDAP CACert", cfg.CACert).Msgf("File does not exist. Waiting %d seconds for it to appear.", caTimeout)
			time.Sleep(caTimeout * time.Second)
			if _, err := os.Stat(cfg.CACert); errors.Is(err, os.ErrNotExist) {
				log.Warn().Str("LDAP CACert", cfg.CACert).Msgf("File does still not exist after Timeout")
				return err
			}
		}
	}
	return nil
}
