package policy

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/extensions/proxy/pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
)

var (
	// ErrMultipleSelectors in case there is more then one selector configured.
	ErrMultipleSelectors = fmt.Errorf("only one type of policy-selector (static, migration, claim or regex) can be configured")
	// ErrSelectorConfigIncomplete if policy_selector conf is missing
	ErrSelectorConfigIncomplete = fmt.Errorf("missing either \"static\", \"migration\", \"claim\" or \"regex\" configuration in policy_selector config ")
	// ErrUnexpectedConfigError unexpected config error
	ErrUnexpectedConfigError = fmt.Errorf("could not initialize policy-selector for given config")
)

const (
	SelectorCookieName = "owncloud-selector"
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
type Selector func(r *http.Request) (string, error)

// LoadSelector constructs a specific policy-selector from a given configuration
func LoadSelector(cfg *config.PolicySelector) (Selector, error) {
	selCount := 0

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

	if cfg.Static == nil && cfg.Claims == nil && cfg.Regex == nil {
		return nil, ErrSelectorConfigIncomplete
	}

	if cfg.Static != nil {
		return NewStaticSelector(cfg.Static), nil
	}

	if cfg.Claims != nil {
		if cfg.Claims.SelectorCookieName == "" {
			cfg.Claims.SelectorCookieName = SelectorCookieName
		}
		return NewClaimsSelector(cfg.Claims), nil
	}

	if cfg.Regex != nil {
		if cfg.Regex.SelectorCookieName == "" {
			cfg.Regex.SelectorCookieName = SelectorCookieName
		}
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
	return func(r *http.Request) (s string, err error) {
		return cfg.Policy, nil
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
	return func(r *http.Request) (s string, err error) {

		selectorCookie := func(r *http.Request) string {
			selectorCookie, err := r.Cookie(cfg.SelectorCookieName)
			if err == nil {
				// TODO check we know the routing policy?
				return selectorCookie.Value
			}
			return ""
		}

		// first, try to route by selector
		if claims := oidc.FromContext(r.Context()); claims != nil {
			if p, ok := claims[oidc.OcisRoutingPolicy].(string); ok && p != "" {
				// TODO check we know the routing policy?
				return p, nil
			}

			// basic auth requests don't have a routing claim, so check for the cookie
			if s := selectorCookie(r); s != "" {
				return s, nil
			}

			return cfg.DefaultPolicy, nil
		}

		// use cookie if provided
		if s := selectorCookie(r); s != "" {
			return s, nil
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
	return func(r *http.Request) (s string, err error) {
		// use cookie first if provided
		selectorCookie, err := r.Cookie(cfg.SelectorCookieName)
		if err == nil {
			return selectorCookie.Value, nil
		}

		// if no cookie is present, try to route by selector
		if u, ok := revactx.ContextGetUser(r.Context()); ok {
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
