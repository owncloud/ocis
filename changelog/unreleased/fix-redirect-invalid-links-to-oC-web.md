Bugfix: redirect invalid links to oC Web

Invalid links ending with a slash(eg. https://foo.bar/index.php/apps/pdfviewer/) have not been redirected to ownCloud Web. Instead the former 404 not found status page was displayed.

https://github.com/owncloud/ocis/pull/2493
https://github.com/owncloud/ocis/pull/2512
