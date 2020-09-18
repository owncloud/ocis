Change: Dummy index.html is not required anymore by upstream

The workaround was required as identifier webapp was mandatory, but
we serve it from memory. This also introduces --disable-identifier-webapp flag.

https://github.com/owncloud/ocis-konnectd/issues/25
