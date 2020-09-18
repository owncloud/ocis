Bugfix: Stop advertising unsupported chunking v2

Removed "chunking" attribute in the DAV capabilities.
Please note that chunking v2 is advertised as "chunking 1.0" while
chunking v1 is the attribute "bigfilechunking" which is already false.

https://github.com/owncloud/ocis/ocis-revapull/145
