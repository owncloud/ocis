package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/proxy/policy"
)

// SelectorCookie provides a middleware which
func SelectorCookie(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger
	policySelector := options.PolicySelector

	return func(next http.Handler) http.Handler {
		return &selectorCookie{
			next:           next,
			logger:         logger,
			policySelector: policySelector,
		}
	}
}

type selectorCookie struct {
	next           http.Handler
	logger         log.Logger
	policySelector config.PolicySelector
}

func (m selectorCookie) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if m.policySelector.Regex == nil && m.policySelector.Claims == nil {
		// only set selector cookie for regex and claim selectors
		m.next.ServeHTTP(w, req)
		return
	}

	ctx := req.Context()
	claims := oidc.FromContext(ctx)

	selectorCookieName := ""
	if m.policySelector.Regex != nil {
		selectorCookieName = m.policySelector.Regex.SelectorCookieName
	} else if m.policySelector.Claims != nil {
		selectorCookieName = m.policySelector.Claims.SelectorCookieName
	}

	_, err := req.Cookie(selectorCookieName)
	if err != nil {
		// no cookie there - try to add one
		if claims != nil {

			selectorFunc, err := policy.LoadSelector(&m.policySelector)
			if err != nil {
				m.logger.Err(err)
			}

			selector, err := selectorFunc(ctx, req)
			if err != nil {
				m.logger.Err(err)
			}

			cookie := http.Cookie{
				Name:     selectorCookieName,
				Value:    selector,
				Domain:   req.Host,
				Path:     "/",
				MaxAge:   60 * 60,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
		}
	}

	m.next.ServeHTTP(w, req)
}
