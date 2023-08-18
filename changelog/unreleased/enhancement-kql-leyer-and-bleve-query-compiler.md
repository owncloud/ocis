Enhancement: Keyword Query Language (KQL) search syntax support

Introduce support for [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference) syntax for search.

The functionality consists of a kql lexer and a rego query compiler

Supported field queries:

* `Tag` search `tag:golden tag:"silver"`
* `Filename` search `name:file.txt name:"file.docx"`
* `Content` search `content:ahab content:"captain aha*"`

Supported conjunctive normal form queries:

* `Boolean` operators `AND`, `OR`, `NOT`,
* `Nesting` `(` `SUB_QUERY` `)`

some examples are:

query: `(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`

* Resources with `name: moby di*` `OR` `tag: bestseller`.
* `AND` with `tag:book`.
* `NOT` with `tag:read`.

https://github.com/owncloud/ocis/pull/7043
https://github.com/owncloud/ocis/issues/7042
