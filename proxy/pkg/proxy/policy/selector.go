package policy

import (
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/asim/go-micro/plugins/client/grpc/v3"
	revauser "github.com/cs3org/reva/pkg/user"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/config"
)

var (
	// ErrMultipleSelectors in case there is more then one selector configured.
	ErrMultipleSelectors = fmt.Errorf("only one type of policy-selector (static, migration, claim or regex) can be configured")
	// ErrSelectorConfigIncomplete if policy_selector conf is missing
	ErrSelectorConfigIncomplete = fmt.Errorf("missing either \"static\", \"migration\", \"claim\" or \"regex\" configuration in policy_selector config ")
	// ErrUnexpectedConfigError unexpected config error
	ErrUnexpectedConfigError = fmt.Errorf("could not initialize policy-selector for given config")
)

// Selector is a function which selects a proxy-policy based on the request.
//
// A policy is a random name which identifies a set of proxy-routes:
//{
//  "policies": [
//    {
//      "name": "us-east-1",
//      "routes": [
//        {
//          "endpoint": "/",
//          "backend": "https://backend.us.example.com:8080/app"
//        }
//      ]
//    },
//    {
//      "name": "eu-ams-1",
//      "routes": [
//        {
//          "endpoint": "/",
//          "backend": "https://backend.eu.example.com:8080/app"
//        }
//      ]
//    }
//  ]
//}
type Selector func(ctx context.Context) (string, error)

// LoadSelector constructs a specific policy-selector from a given configuration
func LoadSelector(cfg *config.PolicySelector) (Selector, error) {
	selCount := 0

	if cfg.Migration != nil {
		selCount++
	}
	if cfg.Static != nil {
		selCount++
	}
	if cfg.Claims != nil {
		selCount++
	}
	if cfg.Regex != nil {
		selCount++
	}
	if selCount > 1 {
		return nil, ErrMultipleSelectors
	}

	if cfg.Migration == nil && cfg.Static == nil && cfg.Claims == nil && cfg.Regex == nil {
		return nil, ErrSelectorConfigIncomplete
	}

	if cfg.Static != nil {
		return NewStaticSelector(cfg.Static), nil
	}

	if cfg.Migration != nil {
		return NewMigrationSelector(
			cfg.Migration,
			accounts.NewAccountsService("com.owncloud.accounts", grpc.NewClient())), nil
	}

	if cfg.Claims != nil {
		return NewClaimsSelector(cfg.Claims), nil
	}

	if cfg.Regex != nil {
		return NewRegexSelector(cfg.Regex), nil
	}

	return nil, ErrUnexpectedConfigError
}

// NewStaticSelector returns a selector which uses a pre-configured policy.
//
// Configuration:
//
// "policy_selector": {
//    "static": {"policy" : "ocis"}
//  },
func NewStaticSelector(cfg *config.StaticSelectorConf) Selector {
	return func(ctx context.Context) (s string, err error) {
		return cfg.Policy, nil
	}
}

// NewMigrationSelector selects the policy based on the existence of the oidc "preferred_username" claim in the accounts-service.
// The policy for each case is configurable:
// "policy_selector": {
//    "migration": {
//      "acc_found_policy" : "ocis",
//      "acc_not_found_policy": "oc10",
//      "unauthenticated_policy": "oc10"
//    }
//  },
//
// This selector can be used in migration-scenarios where some users have already migrated from ownCloud10 to OCIS and
// thus have an entry in ocis-accounts. All users without accounts entry are routed to the legacy ownCloud10 instance.
func NewMigrationSelector(cfg *config.MigrationSelectorConf, ss accounts.AccountsService) Selector {
	var acc = ss
	return func(ctx context.Context) (s string, err error) {
		var claims map[string]interface{}
		if claims = oidc.FromContext(ctx); claims == nil {
			return cfg.UnauthenticatedPolicy, nil
		}

		var userID string
		var ok bool
		if userID, ok = claims[oidc.PreferredUsername].(string); !ok {
			// TODO clarify: what if the user just has no username ...
			return cfg.AccNotFoundPolicy, nil
		}

		if _, err := acc.GetAccount(ctx, &accounts.GetAccountRequest{Id: userID}); err != nil {
			return cfg.AccNotFoundPolicy, nil
		}
		return cfg.AccFoundPolicy, nil

	}
}

// NewClaimsSelector selects the policy based on the "ocis.routing.policy" claim
// The policy for corner cases is configurable:
// "policy_selector": {
//    "migration": {
//      "default_policy" : "ocis",
//      "unauthenticated_policy": "oc10"
//    }
//  },
//
// This selector can be used in migration-scenarios where some users have already migrated from ownCloud10 to OCIS and
func NewClaimsSelector(cfg *config.ClaimsSelectorConf) Selector {
	return func(ctx context.Context) (s string, err error) {
		if claims := oidc.FromContext(ctx); claims != nil {
			if p, ok := claims[oidc.OcisRoutingPolicy].(string); ok && p != "" {
				// TODO check we know the routing policy?
				return p, nil
			}
			return cfg.DefaultPolicy, nil
		}

		return cfg.UnauthenticatedPolicy, nil
	}
}

// NewRegexSelector selects the policy based on a user property
// The policy for each case is configurable:
// "policy_selector": {
//    "regex": {
//      "matches_policies": [
//        {"priority": 10, "property": "mail", "match": "marie@example.org", "policy": "ocis"},
//        {"priority": 20, "property": "mail", "match": "[^@]+@example.org", "policy": "oc10"},
//        {"priority": 30, "property": "username", "match": "(einstein|feynman)", "policy": "ocis"},
//        {"priority": 40, "property": "username", "match": ".+", "policy": "oc10"},
//        {"priority": 50, "property": "id", "match": "4c510ada-c86b-4815-8820-42cdf82c3d51", "policy": "ocis"},
//        {"priority": 60, "property": "id", "match": "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c", "policy": "oc10"}
//      ],
//      "unauthenticated_policy": "oc10"
//    }
//  },
//
// This selector can be used in migration-scenarios where some users have already migrated from ownCloud10 to OCIS and
func NewRegexSelector(cfg *config.RegexSelectorConf) Selector {
	regexRules := []*regexRule{}
	sort.Slice(cfg.MatchesPolicies, func(i, j int) bool {
		return cfg.MatchesPolicies[i].Priority < cfg.MatchesPolicies[j].Priority
	})
	for i := range cfg.MatchesPolicies {
		regexRules = append(regexRules, &regexRule{
			property: cfg.MatchesPolicies[i].Property,
			rule:     regexp.MustCompile(cfg.MatchesPolicies[i].Match),
			policy:   cfg.MatchesPolicies[i].Policy,
		})
	}
	return func(ctx context.Context) (s string, err error) {
		if u, ok := revauser.ContextGetUser(ctx); ok {
			for i := range regexRules {
				switch regexRules[i].property {
				case "mail":
					if regexRules[i].rule.MatchString(u.Mail) {
						return regexRules[i].policy, nil
					}
				case "username":
					if regexRules[i].rule.MatchString(u.Username) {
						return regexRules[i].policy, nil
					}
				case "id":
					if u.Id != nil && regexRules[i].rule.MatchString(u.Id.OpaqueId) {
						return regexRules[i].policy, nil
					}
				}
			}
			return cfg.DefaultPolicy, nil
		}

		return cfg.UnauthenticatedPolicy, nil
	}
}

type regexRule struct {
	property string
	rule     *regexp.Regexp
	policy   string
}
