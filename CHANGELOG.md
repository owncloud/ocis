# Changelog for [0.3.1] (2020-04-14)

The following sections list the changes in ocis-konnectd 0.3.1.

[0.3.1]: https://github.com/owncloud/ocis-konnectd/compare/v0.3.0...v0.3.1

## Summary

* Bugfix - Include the assets for #62: [#64](https://github.com/owncloud/ocis-konnectd/pull/64)

## Details

* Bugfix - Include the assets for #62: [#64](https://github.com/owncloud/ocis-konnectd/pull/64)

   PR 62 introduced new client names. These assets needs to be generated in the embed.go file.

   https://github.com/owncloud/ocis-konnectd/pull/64

# Changelog for [0.3.0] (2020-04-14)

The following sections list the changes in ocis-konnectd 0.3.0.

[0.3.0]: https://github.com/owncloud/ocis-konnectd/compare/v0.1.0...v0.3.0

## Summary

* Bugfix - Redirect to the provided uri: [#26](https://github.com/owncloud/ocis-konnectd/issues/26)
* Change - Add a trailing slash to trusted redirect uris: [#26](https://github.com/owncloud/ocis-konnectd/issues/26)
* Change - Improve client identifiers for end users: [#62](https://github.com/owncloud/ocis-konnectd/pull/62)
* Enhancement - Use upstream version of konnect library: [#14](https://github.com/owncloud/product/issues/14)

## Details

* Bugfix - Redirect to the provided uri: [#26](https://github.com/owncloud/ocis-konnectd/issues/26)

   The phoenix client was not set as trusted therefore when logging out the user was redirected to a
   default page instead of the provided url.

   https://github.com/owncloud/ocis-konnectd/issues/26


* Change - Add a trailing slash to trusted redirect uris: [#26](https://github.com/owncloud/ocis-konnectd/issues/26)

   Phoenix changed the redirect uri to `<baseUrl>#/login` that means it will contain a trailing
   slash after the base url.

   https://github.com/owncloud/ocis-konnectd/issues/26


* Change - Improve client identifiers for end users: [#62](https://github.com/owncloud/ocis-konnectd/pull/62)

   Improved end user facing client names in default identifier-registration.yaml

   https://github.com/owncloud/ocis-konnectd/pull/62


* Enhancement - Use upstream version of konnect library: [#14](https://github.com/owncloud/product/issues/14)

   https://github.com/owncloud/product/issues/14

# Changelog for [0.1.0] (2020-03-18)

The following sections list the changes in ocis-konnectd 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-konnectd/compare/v0.2.0...v0.1.0

[0.1.0]: https://github.com/owncloud/ocis-konnectd/compare/66337bb4dad4a3202880323adf7a51a1a3bb4085...v0.1.0

## Summary

* Bugfix - Generate a random CSP-Nonce in the webapp: [#17](https://github.com/owncloud/ocis-konnectd/issues/17)
* Change - Dummy index.html is not required anymore by upstream: [#25](https://github.com/owncloud/ocis-konnectd/issues/25)
* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-konnectd/issues/1)
* Change - Use glauth as ldap backend, default to running behind ocis-proxy: [#52](https://github.com/owncloud/ocis-konnectd/pull/52)

## Details

* Bugfix - Generate a random CSP-Nonce in the webapp: [#17](https://github.com/owncloud/ocis-konnectd/issues/17)

   https://github.com/owncloud/ocis-konnectd/issues/17
   https://github.com/owncloud/ocis-konnectd/pull/29


* Change - Dummy index.html is not required anymore by upstream: [#25](https://github.com/owncloud/ocis-konnectd/issues/25)

   The workaround was required as identifier webapp was mandatory, but we serve it from memory.
   This also introduces --disable-identifier-webapp flag.

   https://github.com/owncloud/ocis-konnectd/issues/25


* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-konnectd/issues/1)

   Just prepare an initial basic version to serve konnectd embedded into our microservice
   infrastructure in the scope of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis-konnectd/issues/1


* Change - Use glauth as ldap backend, default to running behind ocis-proxy: [#52](https://github.com/owncloud/ocis-konnectd/pull/52)

   We changed the default configuration to integrate better with ocis.

   The default ldap port changes to 9125, which is used by ocis-glauth and we use ocis-proxy to do
   the tls offloading. Clients are supposed to use the ocis-proxy endpoint
   `https://localhost:9200`

   https://github.com/owncloud/ocis-konnectd/pull/52

# Changelog for [0.2.0] (2020-03-18)

The following sections list the changes in ocis-konnectd 0.2.0.

## Summary

* Enhancement - Change default config for single-binary: [#55](https://github.com/owncloud/ocis-konnectd/pull/55)

## Details

* Enhancement - Change default config for single-binary: [#55](https://github.com/owncloud/ocis-konnectd/pull/55)

   https://github.com/owncloud/ocis-konnectd/pull/55

