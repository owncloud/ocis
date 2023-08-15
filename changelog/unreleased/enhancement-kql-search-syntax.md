Enhancement: Keyword Query Language (KQL) search syntax support

Introduce support for [KQL](https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference) syntax for search.
Not every KQL function is currently supported, the supported syntax will grow over time.

The following queries / query elements are supported:

* Tag search `tag:golden tag:"silver"`
* Filename search `name:file.txt name:"file.docx"`
* Content search `content:ahab content:"captain aha*"`

queries can be combined as follows `tag:"book" tag:"bestseller" name:"whale-books-1851-20*" content:"captain aha*"`,
which then gives the following result:

* Resources with a `book` tag.
* `AND` the tag `bestseller`.
* `AND` the name `whale-books-1851-20023.docx`, `whale-books-1851-20023.pdf`, ... .
* `AND` the content contains `captain ahab`.

https://github.com/owncloud/ocis/pull/7043
https://github.com/owncloud/ocis/issues/7042
