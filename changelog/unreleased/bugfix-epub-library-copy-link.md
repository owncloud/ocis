Bugfix: fixed copy link failing in the EPUB library book details

Clicking "Copy link" in the details of a book in the EPUB library showed
a "The link could not be copied" error whenever the book had been found
through the server-side search, because those search results carry no
link to the file.

The link is now built from the id of the file when the search result does
not provide one, so copying the link works for all books.

https://github.com/owncloud/ocis/pull/13239
