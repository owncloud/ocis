Enhancement: Keyword Query Language (KQL) search syntax

We've introduced support for [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference) as the default oCIS search query language.

Simple queries:

* `tag:golden tag:"silver"`
* `name:file.txt name:"file.docx"`
* `content:ahab content:"captain aha*"`

Date/-range queries

* `Mtime:"2023-09-05T08:42:11.23554+02:00"`
* `Mtime>"2023-09-05T08:42:11.23554+02:00"`
* `Mtime>="2023-09-05T08:42:11.23554+02:00"`
* `Mtime<"2023-09-05T08:42:11.23554+02:00"`
* `Mtime<="2023-09-05T08:42:11.23554+02:00"`
* `Mtime:today` - range: start of today till end of today
* `Mtime:yesterday` - range: start of yesterday till end of yesterday
* `Mtime:"this week"` - range: start of this week till end of this week
* `Mtime:"this month"` - range: start of this month till end of this month
* `Mtime:"last month"` - range: start of last month till end of last month
* `Mtime:"this year"` - range: start of this year till end of this year
* `Mtime:"last year"` - range: start of last year till end of last year

Conjunctive normal form queries:

* `tag:golden AND tag:"silver`, `tag:golden OR tag:"silver`, `tag:golden NOT tag:"silver`
* `(tag:book content:ahab*)`, `tag:(book pdf)`

Complex queries:

* `(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`

https://github.com/owncloud/ocis/pull/7212
https://github.com/owncloud/ocis/pull/7043
https://github.com/owncloud/ocis/pull/7247
https://github.com/owncloud/ocis/pull/7248
https://github.com/owncloud/ocis/pull/7254
https://github.com/owncloud/ocis/pull/7262
https://github.com/owncloud/web/pull/9653
https://github.com/owncloud/web/pull/9672
https://github.com/owncloud/ocis/issues/7042
https://github.com/owncloud/ocis/issues/7179
https://github.com/owncloud/ocis/issues/7114
https://github.com/owncloud/web/issues/9636
https://github.com/owncloud/web/issues/9646
