package command

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// configFilePerm is the permission mask used when rewriting ocis.yaml, matching
// the mask `ocis init` uses to create it.
const configFilePerm = 0o600

// serviceUserConfigKeys maps an IDM service user (the ou=sysusers entries plus
// the bootstrap password keys) to the ocis.yaml key paths that must hold the
// same secret so services can still bind to LDAP after a password reset.
//
// reva, libregraph and idp share OCIS_LDAP_BIND_PASSWORD across several
// services but bind as distinct DNs, so their bind_password lives under
// different service sections. The idm.service_user_passwords.* keys are what
// `ocis init` writes and what re-bootstrapping a fresh database would read.
var serviceUserConfigKeys = map[string][][]string{
	"reva": {
		{"auth_basic", "auth_providers", "ldap", "bind_password"},
		{"users", "drivers", "ldap", "bind_password"},
		{"groups", "drivers", "ldap", "bind_password"},
		{"idm", "service_user_passwords", "reva_password"},
	},
	"libregraph": {
		{"graph", "identity", "ldap", "bind_password"},
		{"idm", "service_user_passwords", "idm_password"},
	},
	"idp": {
		{"idp", "ldap", "bind_password"},
		{"idm", "service_user_passwords", "idp_password"},
	},
}

// syncServiceUserPasswordInConfig rewrites the bind_password / service user
// password keys associated with the given service user to newPassword,
// preserving the rest of the document (including comments and unrelated keys).
// It returns the updated document, the dotted paths that were actually changed,
// and any error. Paths that are not present in the document are skipped.
func syncServiceUserPasswordInConfig(data []byte, userName, newPassword string) ([]byte, []string, error) {
	paths, ok := serviceUserConfigKeys[userName]
	if !ok {
		// Not a known service user with bind config to sync.
		return data, nil, nil
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, nil, fmt.Errorf("failed to parse config: %w", err)
	}
	if len(root.Content) == 0 {
		return data, nil, nil
	}
	doc := root.Content[0]

	updated := make([]string, 0, len(paths))
	for _, path := range paths {
		if node := lookupYAMLNode(doc, path); node != nil {
			node.SetString(newPassword)
			updated = append(updated, strings.Join(path, "."))
		}
	}

	if len(updated) == 0 {
		return data, nil, nil
	}

	out, err := yaml.Marshal(&root)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize config: %w", err)
	}
	return out, updated, nil
}

// lookupYAMLNode walks a mapping node along path and returns the scalar value
// node at the end, or nil if any segment is missing.
func lookupYAMLNode(node *yaml.Node, path []string) *yaml.Node {
	for _, key := range path {
		if node.Kind != yaml.MappingNode {
			return nil
		}
		var next *yaml.Node
		for i := 0; i+1 < len(node.Content); i += 2 {
			if node.Content[i].Value == key {
				next = node.Content[i+1]
				break
			}
		}
		if next == nil {
			return nil
		}
		node = next
	}
	return node
}
