# Changes in unreleased

## Summary

* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)
* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)
* Enhancement - Update extensions: [#151](https://github.com/owncloud/ocis/pull/151)

## Details

* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)

   Just prepared an initial basic version which simply embeds the minimum of required services in
   the context of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis/issues/2


* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)

   Starts the proxy in single binary mode (./ocis server) on port 9200. The proxy serves as a
   single-entry point for all http-clients.

   https://github.com/owncloud/ocis/issues/119
   https://github.com/owncloud/ocis/issues/136


* Enhancement - Update extensions: [#151](https://github.com/owncloud/ocis/pull/151)

   We've updated various extensions to a tagged release: - ocis-konnectd v0.2.0 - ocis-glauth
   v0.4.0 - ocis-phoenix v0.3.0 (phoenix v0.6.0) - ocis-pkg v2.1.0 - ocis-proxy -v0.1.0

   We also updated ocis-reva to a PR commit that brings the latest reva to ocis. Work on the PR is
   ongoing because some acceptance tests fail.

   https://github.com/owncloud/ocis/pull/151

