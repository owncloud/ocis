# Changes in unreleased

## Summary

* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)
* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)
* Bugfix - Fix director selection in proxy: [#521](https://github.com/owncloud/ocis/pull/521)
* Bugfix - Cleanup separated indices in memory: [#224](https://github.com/owncloud/product/issues/224)
* Bugfix - Update ocis-glauth for fixed single user search: [#214](https://github.com/owncloud/product/issues/214)
* Bugfix - Fix builtin config for external apps: [#218](https://github.com/owncloud/product/issues/218)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)
* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)
* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)
* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)
* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)
* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)
* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Account management permissions for Admin role: [#124](https://github.com/owncloud/product/issues/124)
* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Update ocis-settings to v0.2.0: [#467](https://github.com/owncloud/ocis/pull/467)
* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)
* Change - Update ocis-accounts to v0.4.0: [#479](https://github.com/owncloud/ocis/pull/479)
* Change - Create accounts in accounts UI: [#148](https://github.com/owncloud/product/issues/148)
* Change - Delete accounts in accounts UI: [#148](https://github.com/owncloud/product/issues/148)
* Change - Enable/disable accounts in accounts UI: [#118](https://github.com/owncloud/product/issues/118)
* Change - Update ocis-ocs to v0.3.0: [#500](https://github.com/owncloud/ocis/pull/500)
* Change - Update ocis-phoenix to v0.13.0: [#487](https://github.com/owncloud/ocis/pull/487)
* Change - Update ocis-proxy to v0.7.0: [#476](https://github.com/owncloud/ocis/pull/476)
* Change - Update ocis-reva to 0.13.0: [#496](https://github.com/owncloud/ocis/pull/496)
* Change - Update proxy with disabled accounts cache: [#525](https://github.com/owncloud/ocis/pull/525)
* Change - Update ocis-reva to v0.14.0: [#556](https://github.com/owncloud/ocis/pull/556)
* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)
* Change - Update ocis-settings to v0.3.0: [#490](https://github.com/owncloud/ocis/pull/490)
* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)
* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)
* Enhancement - Accounts UI improvements: [#222](https://github.com/owncloud/product/issues/222)
* Enhancement - Add new REVA config variables to docs: [#345](https://github.com/owncloud/ocis/pull/345)
* Enhancement - Update extensions: [#180](https://github.com/owncloud/ocis/pull/180)
* Enhancement - Update extensions 2020-07-01: [#357](https://github.com/owncloud/ocis/pull/357)
* Enhancement - Update extensions 2020-09-02: [#516](https://github.com/owncloud/ocis/pull/516)
* Enhancement - Update extensions: [#209](https://github.com/owncloud/ocis/pull/209)
* Enhancement - Update extensions: [#151](https://github.com/owncloud/ocis/pull/151)
* Enhancement - Update extensions 2020-07-10: [#376](https://github.com/owncloud/ocis/pull/376)
* Enhancement - Update extensions: [#290](https://github.com/owncloud/ocis/pull/290)
* Enhancement - Update ocis-reva to 0.4.0: [#295](https://github.com/owncloud/ocis/pull/295)
* Enhancement - Update extensions: [#209](https://github.com/owncloud/ocis/pull/209)
* Enhancement - Update extensions 2020-06-29: [#334](https://github.com/owncloud/ocis/pull/334)
* Enhancement - Update proxy and reva: [#466](https://github.com/owncloud/ocis/pull/466)
* Enhancement - Update proxy to v0.2.0: [#167](https://github.com/owncloud/ocis/pull/167)
* Enhancement - Update ocis-reva 2020-09-10: [#334](https://github.com/owncloud/ocis/pull/334)

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


* Bugfix - Cleanup separated indices in memory: [#224](https://github.com/owncloud/product/issues/224)

   The accounts service was creating a bleve index instance in the service handler, thus creating
   separate in memory indices for the http and grpc servers. We moved the service handler creation
   out of the server creation so that the service handler, thus also the bleve index, is a shared
   instance of the servers.

   This fixes a bug that accounts created through the web ui were not able to sign in until a service
   restart.

   https://github.com/owncloud/product/issues/224
   https://github.com/owncloud/ocis-accounts/pull/117
   https://github.com/owncloud/ocis-accounts/pull/118
   https://github.com/owncloud/ocis/pull/555


* Bugfix - Update ocis-glauth for fixed single user search: [#214](https://github.com/owncloud/product/issues/214)

   We updated ocis-glauth to a version that comes with a fix for searching a single user or group.
   ocis-glauth was dropping search context before by ignoring the searchBaseDN for filtering.
   This has been fixed.

   https://github.com/owncloud/product/issues/214
   https://github.com/owncloud/ocis/pull/535
   https://github.com/owncloud/ocis-glauth/pull/32


* Bugfix - Fix builtin config for external apps: [#218](https://github.com/owncloud/product/issues/218)

   We fixed a bug in the builtin config of ocis-phoenix, having hardcoded urls instead of just the
   path of external apps.

   https://github.com/owncloud/product/issues/218
   https://github.com/owncloud/ocis-phoenix/pull/83
   https://github.com/owncloud/ocis/pull/544


* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis/pull/416


* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)

   Added the thumbnails command so that the thumbnails service can get started via ocis.

   https://github.com/owncloud/ocis/issues/156


* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)

   https://github.com/owncloud/ocis/pull/249
   https://github.com/owncloud/ocis-migration


* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)

   Just prepared an initial basic version which simply embeds the minimum of required services in
   the context of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis/issues/2


* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)

   COMMANDS: - list, ls List existing accounts - add, create Create a new account - update Make
   changes to an existing account - remove, rm Removes an existing account - inspect Show detailed
   data on an existing account - help, h Shows a list of commands or help for one command

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


* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)

   This version delivers `settings` as a new service. It is part of the array of services in the
   `server` command.

   https://github.com/owncloud/ocis/pull/287


* Change - Update ocis-settings to v0.2.0: [#467](https://github.com/owncloud/ocis/pull/467)

   This version delivers `settings` v0.2.0 and versions of accounts (v0.3.0) and phoenix
   (v0.15.0) needed for it.

   https://github.com/owncloud/ocis/pull/467


* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)

   Starts the proxy in single binary mode (./ocis server) on port 9200. The proxy serves as a
   single-entry point for all http-clients.

   https://github.com/owncloud/ocis/issues/119
   https://github.com/owncloud/ocis/issues/136


* Change - Update ocis-accounts to v0.4.0: [#479](https://github.com/owncloud/ocis/pull/479)

   Provides a web UI for role assignment.

   https://github.com/owncloud/ocis/pull/479


* Change - Create accounts in accounts UI: [#148](https://github.com/owncloud/product/issues/148)

   We've added a form to create new users above the accounts list in the accounts UI.

   https://github.com/owncloud/product/issues/148
   https://github.com/owncloud/ocis-accounts/pull/115
   https://github.com/owncloud/ocis/pull/525


* Change - Delete accounts in accounts UI: [#148](https://github.com/owncloud/product/issues/148)

   We've added an action into the actions dropdown of the accounts UI to enable admins to delete
   users.

   https://github.com/owncloud/product/issues/148
   https://github.com/owncloud/ocis-accounts/pull/115
   https://github.com/owncloud/ocis/pull/525


* Change - Enable/disable accounts in accounts UI: [#118](https://github.com/owncloud/product/issues/118)

   We added a new feature in the ocis-accounts web extension to enable or disable accounts. This
   also introduces batch actions, where accounts can be selected and a batch action applied to
   them. The UI for this is the same as in the files extension of ocis-web.

   https://github.com/owncloud/product/issues/118
   https://github.com/owncloud/ocis-accounts/pull/109
   https://github.com/owncloud/ocis/pull/525


* Change - Update ocis-ocs to v0.3.0: [#500](https://github.com/owncloud/ocis/pull/500)

   This change updates ocis-ocs to version 0.3.0

   https://github.com/owncloud/ocis/pull/500


* Change - Update ocis-phoenix to v0.13.0: [#487](https://github.com/owncloud/ocis/pull/487)

   This version delivers ocis-phoenix v0.13.0.

   https://github.com/owncloud/ocis/pull/487


* Change - Update ocis-proxy to v0.7.0: [#476](https://github.com/owncloud/ocis/pull/476)

   This version delivers ocis-proxy v0.7.0.

   https://github.com/owncloud/ocis/pull/476


* Change - Update ocis-reva to 0.13.0: [#496](https://github.com/owncloud/ocis/pull/496)

   This version delivers ocis-reva v0.13.0

   https://github.com/owncloud/ocis/pull/496


* Change - Update proxy with disabled accounts cache: [#525](https://github.com/owncloud/ocis/pull/525)

   We removed the accounts cache in ocis-proxy in order to avoid problems with accounts that have
   been updated in ocis-accounts.

   https://github.com/owncloud/ocis/pull/525
   https://github.com/owncloud/ocis-proxy/pull/100
   https://github.com/owncloud/ocis-accounts/pull/114


* Change - Update ocis-reva to v0.14.0: [#556](https://github.com/owncloud/ocis/pull/556)

   - Update ocis-reva to v0.14.0 - Fix default configuration for accessing shares
   (ocis-reva/#461) - Allow configuring arbitrary storage registry rules (ocis-reva/#461) -
   Update reva to v1.2.1-0.20200911111727-51649e37df2d (reva/#454, reva/#466)

   https://github.com/owncloud/ocis/pull/556
   https://github.com/owncloud/ocis/ocis-reva/pull/461
   https://github.com/owncloud/ocis/ocis-reva/pull/454
   https://github.com/owncloud/ocis/ocis-reva/pull/466


* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)

   - EOS homes are not configured with an enable-flag anymore, but with a dedicated storage
   driver. - We're using it now and adapted default configs of storages

   https://github.com/owncloud/ocis/pull/336
   https://github.com/owncloud/ocis/pull/337
   https://github.com/owncloud/ocis/pull/338
   https://github.com/owncloud/ocis/ocis-reva/pull/891


* Change - Update ocis-settings to v0.3.0: [#490](https://github.com/owncloud/ocis/pull/490)

   This version delivers ocis-settings v0.3.0.

   https://github.com/owncloud/ocis/pull/490


* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)

   We have added rules to the Makefile that use the official [eos docker
   images](https://gitlab.cern.ch/eos/eos-docker) to boot an eos cluster and configure OCIS
   to use it.

   https://github.com/owncloud/ocis/pull/172


* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)

   We now apply the oCIS tracing config to all services which have tracing. With this it is possible
   to set one tracing config for all services at the same time.

   https://github.com/owncloud/product/issues/92
   https://github.com/owncloud/ocis/pull/329
   https://github.com/owncloud/ocis/pull/409


* Enhancement - Accounts UI improvements: [#222](https://github.com/owncloud/product/issues/222)

   We aligned the visual appearance of the accounts UI with default ocis-web apps (full width,
   style of batch actions), added icons to buttons, extracted the buttons from the batch actions
   dropdown into individual buttons, improved the wording added a confirmation widget for the
   user deletion and removed the uid and gid columns.

   https://github.com/owncloud/product/issues/222
   https://github.com/owncloud/ocis-accounts/pull/116
   https://github.com/owncloud/ocis/pull/549


* Enhancement - Add new REVA config variables to docs: [#345](https://github.com/owncloud/ocis/pull/345)

   With the default setup of running oCIS with ocis-proxy we need to set `REVA_DATAGATEWAY_URL`
   and `REVA_FRONTEND_URL` environment variables. We added those to the configuration
   instructions in the dev docs.

   https://github.com/owncloud/ocis/pull/345


* Enhancement - Update extensions: [#180](https://github.com/owncloud/ocis/pull/180)

   We've updated various extensions to a tagged release: - ocis-phoenix v0.4.0 (phoenix v0.7.0)
   - ocis-pkg v2.2.0 - ocis-proxy v0.3.1 - ocis-reva v0.1.1 - ocis-thumbnails v0.1.0 -
   ocis-webdav v0.1.0

   https://github.com/owncloud/ocis/pull/180


* Enhancement - Update extensions 2020-07-01: [#357](https://github.com/owncloud/ocis/pull/357)

   - ocis-reva 0.9.0

   https://github.com/owncloud/ocis/pull/357


* Enhancement - Update extensions 2020-09-02: [#516](https://github.com/owncloud/ocis/pull/516)

   - ocis-accounts 0.4.2-0.20200828150703-2ca83cf4ac20 - ocis-ocs 0.3.1 - ocis-settings
   0.3.2-0.20200828130413-0cc0f5bf26fe

   https://github.com/owncloud/ocis/pull/516


* Enhancement - Update extensions: [#209](https://github.com/owncloud/ocis/pull/209)

   We've updated various extensions: - ocis-konnectd v0.3.1 - ocis-phoenix v0.5.0 (phoenix
   v0.8.0) - ocis-reva v0.2.0

   https://github.com/owncloud/ocis/pull/209


* Enhancement - Update extensions: [#151](https://github.com/owncloud/ocis/pull/151)

   We've updated various extensions to a tagged release: - ocis-konnectd v0.2.0 - ocis-glauth
   v0.4.0 - ocis-phoenix v0.3.0 (phoenix v0.6.0) - ocis-pkg v2.1.0 - ocis-proxy v0.1.0 -
   ocis-reva v0.1.0

   https://github.com/owncloud/ocis/pull/151


* Enhancement - Update extensions 2020-07-10: [#376](https://github.com/owncloud/ocis/pull/376)

   - ocis-reva 0.10.0 - ocis-phoenix 0.9.0

   https://github.com/owncloud/ocis/pull/376


* Enhancement - Update extensions: [#290](https://github.com/owncloud/ocis/pull/290)

   We've updated various extensions: - ocis-thumbnails v0.1.2 (tag) - ocis-reva v0.3.0 (tag)

   https://github.com/owncloud/ocis/pull/290


* Enhancement - Update ocis-reva to 0.4.0: [#295](https://github.com/owncloud/ocis/pull/295)

   Brings in fixes for trashbin and TUS upload. Also adds partial implementation of public
   shares.

   https://github.com/owncloud/ocis/pull/295


* Enhancement - Update extensions: [#209](https://github.com/owncloud/ocis/pull/209)

   We've updated various extensions: - ocis-konnectd v0.3.1 - ocis-phoenix v0.6.0 - ocis-reva
   v0.2.1 - ocis-pkg v2.2.1 - ocis-thumbnails v0.1.2

   https://github.com/owncloud/ocis/pull/209


* Enhancement - Update extensions 2020-06-29: [#334](https://github.com/owncloud/ocis/pull/334)

   - ocis-proxy 0.4.0 - ocis-migration 0.2.0 - ocis-reva 0.8.0 - ocis-phoenix 0.8.1

   https://github.com/owncloud/ocis/pull/334


* Enhancement - Update proxy and reva: [#466](https://github.com/owncloud/ocis/pull/466)

   - ocis-reva contains a lot of sharing, eos and trash fixes - ocis-proxy contains fixes to use
   ocis on top of eos

   https://github.com/owncloud/ocis/pull/466


* Enhancement - Update proxy to v0.2.0: [#167](https://github.com/owncloud/ocis/pull/167)

   https://github.com/owncloud/ocis/pull/167


* Enhancement - Update ocis-reva 2020-09-10: [#334](https://github.com/owncloud/ocis/pull/334)

   - ocis-reva v0.13.1-0.20200910085648-26465bbdcf46 - fixes file operations for received
   shares by changing OC storage default config - adds ability to overwrite storage registry
   rules

   https://github.com/owncloud/ocis/pull/334
   https://github.com/owncloud/ocis/ocis-reva/pull/461

