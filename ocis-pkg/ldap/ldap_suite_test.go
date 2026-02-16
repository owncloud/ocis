package ldap_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLdap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ldap Suite")
}
