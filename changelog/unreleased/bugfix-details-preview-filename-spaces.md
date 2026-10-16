Bugfix: Show image preview for files with spaces in the name on public links

The preview in the details panel was not shown for files opened via a public
link when the file name contained spaces. The preview URL was used unquoted in a
CSS `url()`, so the browser dropped it and never requested the thumbnail. The
URL is now quoted.

https://github.com/owncloud/ocis/issues/13137
https://github.com/owncloud/ocis/pull/13150
