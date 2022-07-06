package ldap_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/ldap"
)

var _ = Describe("ldap", func() {
	DescribeTable("EscapeDNAttributeValue should escape special characters",
		func(input, expected string) {
			escaped := ldap.EscapeDNAttributeValue(input)
			Expect(escaped).To(Equal(expected))
		},
		Entry("normal dn", "foobar", "foobar"),
		Entry("including comma", "foo,bar", "foo\\,bar"),
		Entry("including equals", "foo=bar", "foo\\=bar"),
		Entry("beginning with number sign", "#foobar", "\\#foobar"),
		Entry("beginning with space", " foobar", "\\ foobar"),
		Entry("only one space", " ", "\\ "),
		Entry("two spaces", "  ", "\\ \\ "),
		Entry("ending with space", "foobar ", "foobar\\ "),
		Entry("containing multiple special chars", "f+o>o,b<a;r=\"\000", `f\+o\>o\,b\<a\;r\=\\"\00`),
	)
})
