Enhancement: Keyword Query Language (KQL) search syntax support

Introduce support for [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference) syntax for search.
Not every KQL function is currently supported, the supported syntax will grow over time.

The following queries / query elements are supported:

* Tag search `tag:golden tag:"silver"`
* Filename search `name:file.txt name:"file.docx"`
* Content search `content:ahab content:"captain aha*"`

queries can be combined as follows `tag:"book" tags:"bestseller" AND name:"whale-books-1851-20*" content:"captain aha*"`,
which then gives the following result:

* Resources with a `book` tag.
* `OR` resources with the tag `bestseller` and with name `whale-books-1851-20*`.
* `OR` the name `whale-books-1851-20023.pdf`.
* `OR` the resource content contains `captain ahab`.

https://github.com/owncloud/ocis/pull/7043
https://github.com/owncloud/ocis/issues/7042
