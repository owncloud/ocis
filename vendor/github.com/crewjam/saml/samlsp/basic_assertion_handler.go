package samlsp

import (
	"github.com/crewjam/saml"
)

var _ AssertionHandler = NopAssertionHandler{}

// NopAssertionHandler is an implementation of AssertionHandler that does nothing.
type NopAssertionHandler struct{}

// HandleAssertion is called and passed a SAML assertion. This implementation does nothing.
func (as NopAssertionHandler) HandleAssertion(_ *saml.Assertion) error {
	return nil
}
