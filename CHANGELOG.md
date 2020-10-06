# Changes in unreleased

## Summary

* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)
* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)
* Bugfix - Fix director selection in proxy: [#521](https://github.com/owncloud/ocis/pull/521)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)
* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)
* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)
* Change - Improve reva service descriptions: [#536](https://github.com/owncloud/ocis/pull/536)
* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)
* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)
* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)
* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Account management permissions for Admin role: [#124](https://github.com/owncloud/product/issues/124)
* Change - Update phoenix to v0.18.0: [#651](https://github.com/owncloud/ocis/pull/651)
* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)
* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)
* Enhancement - Add the accounts service: [#593](https://github.com/owncloud/ocis/pull/593)
* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)
* Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)
* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)

## Details

* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)

   Without setting `REVA_FRONTEND_URL` and `REVA_DATAGATEWAY_URL` uploads would default to
   locahost and fail if `OCIS_DOMAIN` was used to run ocis on a remote host.

   https://github.com/owncloud/ocis/pull/392

* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)

   The command for ocis-phoenix enforced an empty external apps configuration. This was
   removed, as it was blocking a new set of default external apps in ocis-phoenix.

   https://github.com/owncloud/ocis/pull/473

* Bugfix - Fix director selection in proxy: [#521](https://github.com/owncloud/ocis/pull/521)

   We fixed a bug in ocis-proxy where simultaneous requests could be executed on the wrong
   backend.

   https://github.com/owncloud/ocis/pull/521
   https://github.com/owncloud/ocis-proxy/pull/99

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis/pull/416

* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)

   Added the thumbnails command so that the thumbnails service can get started via ocis.

   https://github.com/owncloud/ocis/issues/156

* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)

   https://github.com/owncloud/ocis/pull/249
   https://github.com/owncloud/ocis-migration

* Change - Improve reva service descriptions: [#536](https://github.com/owncloud/ocis/pull/536)

   The descriptions make it clearer that the services actually represent a mount point in the
   combined storage. Each mount point can have a different driver.

   https://github.com/owncloud/ocis/pull/536

* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)

   Just prepared an initial basic version which simply embeds the minimum of required services in
   the context of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis/issues/2

* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)

   COMMANDS:

  * list, ls        List existing accounts
  * add, create     Create a new account
  * update          Make changes to an existing account
  * remove, rm      Removes an existing account
  * inspect         Show detailed data on an existing account
  * help, h         Shows a list of commands or help for one command

   https://github.com/owncloud/product/issues/115

* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)

   Starts ocis-accounts in single binary mode (./ocis server). This service stores the
   user-account information.

   https://github.com/owncloud/product/issues/25
   https://github.com/owncloud/ocis/pull/239/files

* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)

   We moved away from using the go-micro runtime and are now using [our own
   runtime](https://github.com/refs/pman). This allows us to spawn service processes even
   when they are using different versions of go-micro. On top of that we now have the commands `ocis
   list`, `ocis kill` and `ocis run` available for service runtime management.

   https://github.com/owncloud/ocis/pull/287

* Change - Account management permissions for Admin role: [#124](https://github.com/owncloud/product/issues/124)

   We created an `AccountManagement` permission and added it to the default admin role. There are
   permission checks in place to protected http endpoints in ocis-accounts against requests
   without the permission. All existing default users (einstein, marie, richard) have the
   default user role now (doesn't have the `AccountManagement` permission). Additionally,
   there is a new default Admin user with credentials `moss:vista`.

   Known issue: for users without the `AccountManagement` permission, the accounts UI
   extension is still available in the ocis-web app switcher, but the requests for loading the
   users will fail (as expected). We are working on a way to hide the accounts UI extension if the
   user doesn't have the `AccountManagement` permission.

   https://github.com/owncloud/product/issues/124
   https://github.com/owncloud/ocis/settings/pull/59
   https://github.com/owncloud/ocis/settings/pull/66
   https://github.com/owncloud/ocis/settings/pull/67
   https://github.com/owncloud/ocis/settings/pull/69
   https://github.com/owncloud/ocis-proxy/pull/95
   https://github.com/owncloud/ocis-pkg/pull/59
   https://github.com/owncloud/ocis-accounts/pull/95
   https://github.com/owncloud/ocis-accounts/pull/100
   https://github.com/owncloud/ocis-accounts/pull/102

* Change - Update phoenix to v0.18.0: [#651](https://github.com/owncloud/ocis/pull/651)

   We updated phoenix to v0.18.0. Please refer to the changelog (linked) for details on the
   phoenix release. With the ODS release brought in by phoenix we now have proper oc-checkbox and
   oc-radio components for the settings and accounts UI.

   https://github.com/owncloud/ocis/pull/651
   https://github.com/owncloud/phoenix/releases/tag/v0.18.0
   https://github.com/owncloud/owncloud-design-system/releases/tag/v1.12.1

* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)

   This version delivers `settings` as a new service. It is part of the array of services in the
   `server` command.

   https://github.com/owncloud/ocis/pull/287

* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)

   Starts the proxy in single binary mode (./ocis server) on port 9200. The proxy serves as a
   single-entry point for all http-clients.

   https://github.com/owncloud/ocis/issues/119
   https://github.com/owncloud/ocis/issues/136

* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)

  * EOS homes are not configured with an enable-flag anymore, but with a dedicated storage driver.
  * We're using it now and adapted default configs of storages

   https://github.com/owncloud/ocis/pull/336
   https://github.com/owncloud/ocis/pull/337
   https://github.com/owncloud/ocis/pull/338
   https://github.com/owncloud/ocis/ocis-reva/pull/891

* Enhancement - Add the accounts service: [#593](https://github.com/owncloud/ocis/pull/593)

  * Bugfix - Initialize roleService client in GRPC server: [#114](https://github.com/owncloud/ocis-accounts/pull/114)
  * Bugfix - Cleanup separated indices in memory: [#224](https://github.com/owncloud/product/issues/224)
  * Change - Set user role on builtin users: [#102](https://github.com/owncloud/ocis-accounts/pull/102)
  * Change - Add new builtin admin user: [#102](https://github.com/owncloud/ocis-accounts/pull/102)
  * Change - We make use of the roles cache to enforce permission checks: [#100](https://github.com/owncloud/ocis-accounts/pull/100)
  * Change - We make use of the roles manager to enforce permission checks: [#108](https://github.com/owncloud/ocis-accounts/pull/108)
  * Enhancement - Add create account form: [#148](https://github.com/owncloud/product/issues/148)
  * Enhancement - Add delete accounts action: [#148](https://github.com/owncloud/product/issues/148)
  * Enhancement - Add enable/disable capabilities to the WebUI: [#118](https://github.com/owncloud/product/issues/118)
  * Enhancement - Improve visual appearance of accounts UI: [#222](https://github.com/owncloud/product/issues/222)
  * Bugfix - Adapting to new settings API for fetching roles: [#96](https://github.com/owncloud/ocis-accounts/pull/96)
  * Change - Create account api-call implicitly adds "default-user" role: [#173](https://github.com/owncloud/product/issues/173)
  * Change - Add role selection to accounts UI: [#103](https://github.com/owncloud/product/issues/103)
  * Bugfix - Atomic Requests: [#82](https://github.com/owncloud/ocis-accounts/pull/82)
  * Bugfix - Unescape value for prefix query: [#76](https://github.com/owncloud/ocis-accounts/pull/76)
  * Change - Adapt to new ocis-settings data model: [#87](https://github.com/owncloud/ocis-accounts/pull/87)
  * Change - Add permissions for language to default roles: [#88](https://github.com/owncloud/ocis-accounts/pull/88)
  * Bugfix - Add write mutexes: [#71](https://github.com/owncloud/ocis-accounts/pull/71)
  * Bugfix - Fix the accountId and groupId mismatch in DeleteGroup Method: [#60](https://github.com/owncloud/ocis-accounts/pull/60)
  * Bugfix - Fix index mapping: [#73](https://github.com/owncloud/ocis-accounts/issues/73)
  * Bugfix - Use NewNumericRangeInclusiveQuery for numeric literals: [#28](https://github.com/owncloud/ocis-glauth/issues/28)
  * Bugfix - Prevent segfault when no password is set: [#65](https://github.com/owncloud/ocis-accounts/pull/65)
  * Bugfix - Update account return value not used: [#70](https://github.com/owncloud/ocis-accounts/pull/70)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#64](https://github.com/owncloud/ocis-accounts/pull/64)
  * Change - Align structure of this extension with other extensions: [#51](https://github.com/owncloud/ocis-accounts/pull/51)
  * Change - Change api errors: [#11](https://github.com/owncloud/ocis-accounts/issues/11)
  * Change - Enable accounts on creation: [#43](https://github.com/owncloud/ocis-accounts/issues/43)
  * Change - Fix index update on create/update: [#57](https://github.com/owncloud/ocis-accounts/issues/57)
  * Change - Pass around the correct logger throughout the code: [#41](https://github.com/owncloud/ocis-accounts/issues/41)
  * Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis-accounts/pull/33)
  * Change - Tighten screws on usernames and email addresses: [#65](https://github.com/owncloud/ocis-accounts/pull/65)
  * Enhancement - Add early version of cli tools for user-management: [#69](https://github.com/owncloud/ocis-accounts/pull/69)
  * Enhancement - Update accounts API: [#30](https://github.com/owncloud/ocis-accounts/pull/30)
  * Enhancement - Add simple user listing UI: [#51](https://github.com/owncloud/ocis-accounts/pull/51)
  * Enhancement - Logging is configurable: [#24](https://github.com/owncloud/ocis-accounts/pull/24)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-accounts/issues/1)
  * Enhancement - Configuration: [#15](https://github.com/owncloud/ocis-accounts/pull/15)

   https://github.com/owncloud/ocis/pull/593

* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)

   We have added rules to the Makefile that use the official [eos docker
   images](https://gitlab.cern.ch/eos/eos-docker) to boot an eos cluster and configure OCIS
   to use it.

   https://github.com/owncloud/ocis/pull/172

* Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)

   In the future accounts, settings etc. should be stored in a dedicated metadata storage. The
   services should talk to this storage directly, bypassing reva-gateway.

   https://github.com/owncloud/ocis/pull/602

* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)

   We now apply the oCIS tracing config to all services which have tracing. With this it is possible
   to set one tracing config for all services at the same time.

   https://github.com/owncloud/product/issues/92
   https://github.com/owncloud/ocis/pull/329
   https://github.com/owncloud/ocis/pull/409
