package kql_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/services/search/pkg/kql"
	"github.com/owncloud/ocis/v2/services/search/pkg/kql/grammar"
)

var _ = Describe("kql", func() {
	Describe("Parse", func() {
		parserTest := func(q string, e []*grammar.Token) {
			p, err := kql.Parse(q)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(e))
		}

		DescribeTable("TagToken", parserTest,
			Entry("", `tags:foo`, []*grammar.Token{
				{grammar.TagToken, "foo"},
			}),
			Entry("", `tags:"foo bar"`, []*grammar.Token{
				{grammar.TagToken, "foo bar"},
			}),
			Entry("", `tag:foo`, []*grammar.Token{
				{grammar.TagToken, "foo"},
			}),
			Entry("", `tag:"foo bar"`, []*grammar.Token{
				{grammar.TagToken, "foo bar"},
			}),
		)

		DescribeTable("ContentToken", parserTest,
			Entry("", `content:foo`, []*grammar.Token{
				{grammar.ContentToken, "foo"},
			}),
			Entry("", `content:"foo bar"`, []*grammar.Token{
				{grammar.ContentToken, "foo bar"},
			}),
		)

		DescribeTable("NameToken", parserTest,
			Entry("", `name:annual-accounts.xlsx`, []*grammar.Token{
				{grammar.NameToken, "annual-accounts.xlsx"},
			}),
			Entry("", `name:annual-account*`, []*grammar.Token{
				{grammar.NameToken, "annual-account*"},
			}),
			Entry("", `name:"annual accounts.docx"`, []*grammar.Token{
				{grammar.NameToken, "annual accounts.docx"},
			}),
		)

		DescribeTable("KnownTokens", parserTest,
			Entry("", `tags:foo tags:"foo bar" tag:foo tag:"foo bar" name:foo-_-!?12-test.xslx name:"foo bar" content:foo content:"foo bar"`, []*grammar.Token{
				{grammar.TagToken, "foo"},
				{grammar.TagToken, "foo bar"},
				{grammar.TagToken, "foo"},
				{grammar.TagToken, "foo bar"},
				{grammar.NameToken, "foo-_-!?12-test.xslx"},
				{grammar.NameToken, "foo bar"},
				{grammar.ContentToken, "foo"},
				{grammar.ContentToken, "foo bar"},
			}),
		)

		// https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference#specifying-property-restrictions
		// the white space causes the query to return content items containing
		// the terms "author" and "John Smith", instead of content items authored by John Smith
		DescribeTable("FallbackTokens", parserTest,
			Entry("", `author: "John Smith"`, []*grammar.Token{
				{grammar.FallbackToken, "author"},
				{grammar.FallbackToken, "John Smith"},
			}),
			Entry("", `author :"John Smith"`, []*grammar.Token{
				{grammar.FallbackToken, "author"},
				{grammar.FallbackToken, "John Smith"},
			}),
			Entry("", `author : "John Smith"`, []*grammar.Token{
				{grammar.FallbackToken, "author"},
				{grammar.FallbackToken, "John Smith"},
			}),
		)
	})
})
