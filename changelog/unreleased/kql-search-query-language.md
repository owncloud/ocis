Enhancement: Keyword Query Language (KQL) search syntax

We've introduced support for [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference) as the default oCIS search query language.

Some examples of a valid KQL query are:

* `Tag`: `tag:golden tag:"silver"`
* `Filename`: `name:file.txt name:"file.docx"`
* `Content`: `content:ahab content:"captain aha*"`

Conjunctive normal form queries:

* `Boolean`: `tag:golden AND tag:"silver`, `tag:golden OR tag:"silver`, `tag:golden NOT tag:"silver`
* `Group`: `(tag:book content:ahab*)`, `tag:(book pdf)`

Complex queries:

* `(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`

https://github.com/owncloud/ocis/pull/7212
https://github.com/owncloud/ocis/pull/7043
https://github.com/owncloud/web/pull/9653
https://github.com/owncloud/ocis/issues/7042
https://github.com/owncloud/ocis/issues/7179
https://github.com/owncloud/ocis/issues/7114
https://github.com/owncloud/web/issues/9636
https://github.com/owncloud/web/issues/9646
