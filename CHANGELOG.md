# Changes in unreleased

## Summary

* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)
* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)
* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)
* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)
* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)
* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)
* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)
* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Update ocis-settings to v0.2.0: [#467](https://github.com/owncloud/ocis/pull/467)
* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)
* Change - Update ocis-accounts to v0.4.0: [#479](https://github.com/owncloud/ocis/pull/479)
* Change - Update ocis-proxy to v0.7.0: [#476](https://github.com/owncloud/ocis/pull/476)
* Change - Update ocis-reva to 0.13.0: [#496](https://github.com/owncloud/ocis/pull/496)
* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)
* Change - Update ocis-settings to v0.3.0: [#490](https://github.com/owncloud/ocis/pull/490)
* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)
* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)
* Enhancement - Add new REVA config variables to docs: [#345](https://github.com/owncloud/ocis/pull/345)
* Enhancement - Update extensions: [#180](https://github.com/owncloud/ocis/pull/180)
* Enhancement - Update extensions 2020-07-01: [#357](https://github.com/owncloud/ocis/pull/357)
* Enhancement - Update extensions: [#209](https://github.com/owncloud/ocis/pull/209)
* Enhancement - Update extensions: [#151](https://github.com/owncloud/ocis/pull/151)
* Enhancement - Update extensions 2020-07-10: [#376](https://github.com/owncloud/ocis/pull/376)
* Enhancement - Update extensions: [#290](https://github.com/owncloud/ocis/pull/290)
* Enhancement - Update ocis-reva to 0.4.0: [#295](https://github.com/owncloud/ocis/pull/295)
* Enhancement - Update extensions: [#209](https://github.com/owncloud/ocis/pull/209)
* Enhancement - Update extensions 2020-06-29: [#334](https://github.com/owncloud/ocis/pull/334)
* Enhancement - Update proxy and reva: [#466](https://github.com/owncloud/ocis/pull/466)
* Enhancement - Update proxy to v0.2.0: [#167](https://github.com/owncloud/ocis/pull/167)

## Details

* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)

   Without setting `REVA_FRONTEND_URL` and `REVA_DATAGATEWAY_URL` uploads would default to
   locahost and fail if `OCIS_DOMAIN` was used to run ocis on a remote host.

   https://github.com/owncloud/ocis/pull/392


* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)

   The command for ocis-phoenix enforced an empty external apps configuration. This was
   removed, as it was blocking a new set of default external apps in ocis-phoenix.

   https://github.com/owncloud/ocis/pull/473


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


* Change - Update ocis-proxy to v0.7.0: [#476](https://github.com/owncloud/ocis/pull/476)

   This version delivers ocis-proxy v0.7.0.

   https://github.com/owncloud/ocis/pull/476


* Change - Update ocis-reva to 0.13.0: [#496](https://github.com/owncloud/ocis/pull/496)

   This version delivers ocis-reva v0.13.0

   https://github.com/owncloud/ocis/pull/496


* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)

   - EOS homes are not configured with an enable-flag anymore, but with a dedicated storage
   driver. - We're using it now and adapted default configs of storages

   https://github.com/owncloud/ocis/pull/336
   https://github.com/owncloud/ocis/pull/337
   https://github.com/owncloud/ocis/pull/338
   https://github.com/owncloud/ocis-reva/pull/891


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

