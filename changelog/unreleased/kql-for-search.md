Enhancement: Keyword Query Language (KQL) search syntax support

Introduce support for [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference) search syntax.

The functionality consists of a kql lexer and a bleve query compiler

Supported field queries:

* `Tag` search `tag:golden tag:"silver"`
* `Filename` search `name:file.txt name:"file.docx"`
* `Content` search `content:ahab content:"captain aha*"`

Supported conjunctive normal form queries:

* `Boolean`: `AND`, `OR`, `NOT`,
* `Group`: `(tag:book content:ahab*)`, `tag:(book pdf)`

some examples are:

query: `(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`

* Resources with `name: moby di*` `OR` `tag: bestseller`.
* `AND` with `tag:book`.
* `NOT` with `tag:read`.

https://github.com/owncloud/ocis/pull/7043
https://github.com/owncloud/ocis/pull/7196
https://github.com/owncloud/ocis/issues/7042
https://github.com/owncloud/ocis/issues/7179
