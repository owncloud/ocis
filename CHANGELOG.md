# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes for unreleased.

[unreleased]: https://github.com/owncloud/ocis/compare/v2.0.0...master

## Summary

* Bugfix - Return 425 on Thumbnails: [#5300](https://github.com/owncloud/ocis/pull/5300)
* Bugfix - Disassociate users from deleted school: [#5343](https://github.com/owncloud/ocis/pull/5343)
* Bugfix - Fix Postprocessing events: [#5269](https://github.com/owncloud/ocis/pull/5269)
* Enhancement - Add global env variable extractor: [#5164](https://github.com/owncloud/ocis/pull/5164)
* Enhancement - Async Postprocessing: [#5207](https://github.com/owncloud/ocis/pull/5207)
* Enhancement - Bump libre-graph-api-go: [#5309](https://github.com/owncloud/ocis/pull/5309)
* Enhancement - Bump reva version: [#5243](https://github.com/owncloud/ocis/pull/5243)
* Enhancement - Collect global envvars: [#5367](https://github.com/owncloud/ocis/pull/5367)
* Enhancement - Drive group permissions: [#5312](https://github.com/owncloud/ocis/pull/5312)
* Enhancement - Make the group members addition limit configurable: [#5357](https://github.com/owncloud/ocis/pull/5357)
* Enhancement - Graph Drives IdentitySet displayName: [#5347](https://github.com/owncloud/ocis/pull/5347)
* Enhancement - Display surname and givenName attributes: [#5388](https://github.com/owncloud/ocis/pull/5388)
* Enhancement - Extended search: [#5221](https://github.com/owncloud/ocis/pull/5221)
* Enhancement - Resource tags: [#5227](https://github.com/owncloud/ocis/pull/5227)

# Changelog for [2.0.0] (2022-11-30)

The following sections list the changes for 2.0.0.

[2.0.0]: https://github.com/owncloud/ocis/compare/v1.20.0...v2.0.0

## Summary

* Bugfix - Fix configuration of mimetypes for the app registry: [#4411](https://github.com/owncloud/ocis/pull/4411)
* Bugfix - Disable default expiration for public links: [#4445](https://github.com/owncloud/ocis/issues/4445)
* Bugfix - Show help for some commands when unconfigured: [#4405](https://github.com/owncloud/ocis/pull/4405)
* Bugfix - Translations on login page: [#7550](https://github.com/owncloud/web/issues/7550)
* Bugfix - Autocreate IDP private key also if file exists but is empty: [#4394](https://github.com/owncloud/ocis/pull/4394)
* Bugfix - Rename extensions to services (leftover occurences): [#4407](https://github.com/owncloud/ocis/pull/4407)
* Bugfix - Fix DN parsing issues and sizelimit handling in libregraph/idm: [#3631](https://github.com/owncloud/ocis/issues/3631)
* Bugfix - Lower IDP token lifespans: [#5077](https://github.com/owncloud/ocis/pull/5077)
* Bugfix - Remove runtime kill and run commands: [#3740](https://github.com/owncloud/ocis/pull/3740)
* Bugfix - Check permissions when deleting Space: [#3709](https://github.com/owncloud/ocis/pull/3709)
* Bugfix - Do not reindex a space twice at the same time: [#5001](https://github.com/owncloud/ocis/pull/5001)
* Bugfix - Disable federation capabilities: [#4864](https://github.com/owncloud/ocis/pull/4864)
* Bugfix - Decomposedfs increase filelock duration factor: [#5130](https://github.com/owncloud/ocis/pull/5130)
* Bugfix - Find spaces by their name: [#5044](https://github.com/owncloud/ocis/pull/5044)
* Bugfix - Logging in on the wrong account when an email address is not unique: [#4039](https://github.com/owncloud/ocis/issues/4039)
* Bugfix - Allow empty environment variables: [#3892](https://github.com/owncloud/ocis/pull/3892)
* Bugfix - Remove unused transfer secret from app provider: [#3798](https://github.com/owncloud/ocis/pull/3798)
* Bugfix - Fix authentication for autoprovisioned users: [#4616](https://github.com/owncloud/ocis/issues/4616)
* Bugfix - Bring back the settings UI in Web: [#4691](https://github.com/owncloud/ocis/pull/4691)
* Bugfix - Fix cache stat table config: [#4732](https://github.com/owncloud/ocis/pull/4732)
* Bugfix - Adjust cache related configuration options: [#5087](https://github.com/owncloud/ocis/pull/5087)
* Bugfix - Make IDP secrets configurable via environment variables: [#3744](https://github.com/owncloud/ocis/pull/3744)
* Bugfix - CSP rules for silent token refresh in iframe: [#4031](https://github.com/owncloud/ocis/pull/4031)
* Bugfix - Enable debug server by default: [#3827](https://github.com/owncloud/ocis/pull/3827)
* Bugfix - Rework default role provisioning: [#3900](https://github.com/owncloud/ocis/issues/3900)
* Bugfix - Fix search index getting out of sync: [#3851](https://github.com/owncloud/ocis/pull/3851)
* Bugfix - Change the default value for PROXY_OIDC_INSECURE to false: [#4601](https://github.com/owncloud/ocis/pull/4601)
* Bugfix - Fix sharing jsoncs3 driver options: [#4593](https://github.com/owncloud/ocis/pull/4593)
* Bugfix - Inconsistency env var naming for LDAP filter configuration: [#3890](https://github.com/owncloud/ocis/issues/3890)
* Bugfix - Fix LDAP insecure options: [#3897](https://github.com/owncloud/ocis/pull/3897)
* Bugfix - Fix handling of invalid LDAP users and groups: [#4274](https://github.com/owncloud/ocis/issues/4274)
* Bugfix - Fix logging levels: [#4102](https://github.com/owncloud/ocis/pull/4102)
* Bugfix - Don't run auth-bearer service by default: [#4692](https://github.com/owncloud/ocis/issues/4692)
* Bugfix - Fix notifications service settings: [#4652](https://github.com/owncloud/ocis/pull/4652)
* Bugfix - Fix notifications Web UI url: [#4998](https://github.com/owncloud/ocis/pull/4998)
* Bugfix - Fix `OCIS_RUN_SERVICES`: [#4133](https://github.com/owncloud/ocis/pull/4133)
* Bugfix - Fix the OIDC provider cache: [#4600](https://github.com/owncloud/ocis/pull/4600)
* Bugfix - Fix permissions in REPORT: [#4520](https://github.com/owncloud/ocis/pull/4520)
* Bugfix - Set default name for public link via capabilities: [#3834](https://github.com/owncloud/ocis/pull/3834)
* Bugfix - Remove legacy accounts proxy routes: [#3831](https://github.com/owncloud/ocis/pull/3831)
* Bugfix - Fix unused config option `GRAPH_SPACES_INSECURE`: [#55555](https://github.com/owncloud/ocis/pull/55555)
* Bugfix - Remove unused configuration options: [#3973](https://github.com/owncloud/ocis/pull/3973)
* Bugfix - Remove static ocs user backend config: [#4077](https://github.com/owncloud/ocis/pull/4077)
* Bugfix - Remove unused OCS storage configuration: [#3955](https://github.com/owncloud/ocis/pull/3955)
* Bugfix - Fix the `ocis search` command: [#3796](https://github.com/owncloud/ocis/pull/3796)
* Bugfix - Rename search env variable for the grpc server address: [#3800](https://github.com/owncloud/ocis/pull/3800)
* Bugfix - Fix search in received shares: [#4308](https://github.com/owncloud/ocis/issues/4308)
* Bugfix - Fix search report: [#7557](https://github.com/owncloud/web/issues/7557)
* Bugfix - Render webdav permissions as string in search report: [#4575](https://github.com/owncloud/ocis/issues/4575)
* Bugfix - Fix make sensitive config values in the proxy's debug server: [#4086](https://github.com/owncloud/ocis/pull/4086)
* Bugfix - Fix the idm and settings extensions' admin user id configuration option: [#3799](https://github.com/owncloud/ocis/pull/3799)
* Bugfix - Mail notifications for group shares: [#4714](https://github.com/owncloud/ocis/pull/4714)
* Bugfix - Substring search for sharees: [#547](https://github.com/owncloud/ocis/issues/547)
* Bugfix - Fix configuration validation for extensions' server commands: [#3911](https://github.com/owncloud/ocis/pull/3911)
* Bugfix - Fix startup error logging: [#4093](https://github.com/owncloud/ocis/pull/4093)
* Bugfix - Disable cache for selected static web assets: [#4809](https://github.com/owncloud/ocis/pull/4809)
* Bugfix - Fix multiple storage-users env variables: [#3802](https://github.com/owncloud/ocis/pull/3802)
* Bugfix - Thumbnails for `/dav/xxx?preview=1` requests: [#3567](https://github.com/owncloud/ocis/pull/3567)
* Bugfix - Fix unfindable entities from shares/publicshares: [#4651](https://github.com/owncloud/ocis/pull/4651)
* Bugfix - Fix unrestricted quota on the graphAPI: [#4363](https://github.com/owncloud/ocis/pull/4363)
* Bugfix - Fix user autoprovisioning: [#3893](https://github.com/owncloud/ocis/issues/3893)
* Bugfix - Fix version info: [#3953](https://github.com/owncloud/ocis/pull/3953)
* Bugfix - Fix version number in status page: [#3788](https://github.com/owncloud/ocis/issues/3788)
* Bugfix - Fix CORS in frontend service: [#4948](https://github.com/owncloud/ocis/pull/4948)
* Bugfix - Graph service now forwards trace context: [#4582](https://github.com/owncloud/ocis/pull/4582)
* Bugfix - Fix the webdav URL of drive roots: [#3706](https://github.com/owncloud/ocis/issues/3706)
* Bugfix - Idp: Check if CA certificate if present: [#3623](https://github.com/owncloud/ocis/issues/3623)
* Bugfix - Fix graph endpoint: [#3925](https://github.com/owncloud/ocis/issues/3925)
* Bugfix - Initial role assingment with external IDM: [#5045](https://github.com/owncloud/ocis/issues/5045)
* Bugfix - Escape DN attribute value: [#4117](https://github.com/owncloud/ocis/pull/4117)
* Bugfix - Make IDP only wait for certs when using LDAP: [#3965](https://github.com/owncloud/ocis/pull/3965)
* Bugfix - Make ocdav service behave properly: [#3957](https://github.com/owncloud/ocis/pull/3957)
* Bugfix - Make storage users mount ids unique by default: [#5091](https://github.com/owncloud/ocis/pull/5091)
* Bugfix - Return proper errors when ocs/cloud/users is using the cs3 backend: [#3483](https://github.com/owncloud/ocis/issues/3483)
* Bugfix - Polish search: [#4094](https://github.com/owncloud/ocis/pull/4094)
* Bugfix - Fix the shareroot path in REPORT responses: [#4859](https://github.com/owncloud/ocis/pull/4859)
* Bugfix - Remove the storage-users event configuration: [#4825](https://github.com/owncloud/ocis/pull/4825)
* Bugfix - Trigger a rescan of spaces in the search index when items have changed: [#4777](https://github.com/owncloud/ocis/pull/4777)
* Bugfix - Save Katherine: [#3823](https://github.com/owncloud/ocis/issues/3823)
* Bugfix - Fix permission check in settings service: [#4890](https://github.com/owncloud/ocis/pull/4890)
* Bugfix - Fix Thumbnails for IDs without a trailing path: [#3791](https://github.com/owncloud/ocis/pull/3791)
* Bugfix - Space Creators can hand over spaces: [#4244](https://github.com/owncloud/ocis/pull/4244)
* Bugfix - Make tokeninfo endpoint unprotected: [#4715](https://github.com/owncloud/ocis/pull/4715)
* Bugfix - Update reva to version 2.12.0: [#5092](https://github.com/owncloud/ocis/pull/5092)
* Bugfix - URL encode the webdav url in the graph API: [#3597](https://github.com/owncloud/ocis/pull/3597)
* Bugfix - Store user passwords hashed in idm: [#3778](https://github.com/owncloud/ocis/issues/3778)
* Bugfix - Fix wopi access to public shares: [#4631](https://github.com/owncloud/ocis/pull/4631)
* Change - Update ocis packages and imports to V2: [#3678](https://github.com/owncloud/ocis/pull/3678)
* Change - Build service frontends with pnpm instead of yarn: [#4878](https://github.com/owncloud/ocis/pull/4878)
* Change - Load configuration files just from one directory: [#3587](https://github.com/owncloud/ocis/pull/3587)
* Change - Reduce permissions on docker image predeclared volumes: [#3641](https://github.com/owncloud/ocis/pull/3641)
* Change - Introduce `ocis init` and remove all default secrets: [#3551](https://github.com/owncloud/ocis/pull/3551)
* Change - Rename "uploads purge" command to "uploads clean": [#4403](https://github.com/owncloud/ocis/pull/4403)
* Change - Enable privatelinks by default: [#4599](https://github.com/owncloud/ocis/pull/4599/)
* Change - The `glauth` and `accounts` services are removed: [#3685](https://github.com/owncloud/ocis/pull/3685)
* Change - Reduce drives in graph /me/drives API: [#3629](https://github.com/owncloud/ocis/pull/3629)
* Change - Switched default configuration to use libregraph/idm: [#3331](https://github.com/owncloud/ocis/pull/3331)
* Change - Rename MetadataUserID: [#3671](https://github.com/owncloud/ocis/pull/3671)
* Change - Use new space ID util functions: [#3648](https://github.com/owncloud/ocis/pull/3648)
* Change - Prevent access to disabled space: [#3779](https://github.com/owncloud/ocis/pull/3779)
* Change - Rename serviceUser to systemUser: [#3673](https://github.com/owncloud/ocis/pull/3673)
* Change - Use the spaceID on the cs3 resource: [#4748](https://github.com/owncloud/ocis/pull/4748)
* Change - Split MachineAuth from SystemUser: [#3672](https://github.com/owncloud/ocis/pull/3672)
* Enhancement - Add capability for alias links: [#3983](https://github.com/owncloud/ocis/issues/3983)
* Enhancement - Add curl to the oCIS OCI image: [#4751](https://github.com/owncloud/ocis/pull/4751)
* Enhancement - Add deprecation annotation: [#3917](https://github.com/owncloud/ocis/issues/3917)
* Enhancement - Add drives field to users endpoint: [#4072](https://github.com/owncloud/ocis/pull/4072)
* Enhancement - Add Email templating: [#4564](https://github.com/owncloud/ocis/pull/4564)
* Enhancement - Add FRONTEND_ENABLE_RESHARING env variable: [#4023](https://github.com/owncloud/ocis/pull/4023)
* Enhancement - We added e-mail subject templating: [#4799](https://github.com/owncloud/ocis/pull/4799)
* Enhancement - Add number of total matches to the search result: [#4189](https://github.com/owncloud/ocis/issues/4189)
* Enhancement - Add tracing to search: [#5113](https://github.com/owncloud/ocis/pull/5113)
* Enhancement - Add webURL to space root: [#4588](https://github.com/owncloud/ocis/pull/4588)
* Enhancement - Align service naming: [#3606](https://github.com/owncloud/ocis/pull/3606)
* Enhancement - Add acting user to the audit log: [#3753](https://github.com/owncloud/ocis/issues/3753)
* Enhancement - Configurable max lock cycles: [#4965](https://github.com/owncloud/ocis/pull/4965)
* Enhancement - Allow to configuring the reva cache store: [#4627](https://github.com/owncloud/ocis/pull/4627)
* Enhancement - Add audit events for created containers: [#3941](https://github.com/owncloud/ocis/pull/3941)
* Enhancement - Add support for REPORT requests to /dav/spaces URLs: [#4661](https://github.com/owncloud/ocis/pull/4661)
* Enhancement - Don't setup demo role assignments on default: [#3661](https://github.com/owncloud/ocis/issues/3661)
* Enhancement - Introduce "delete-all-spaces" permission: [#4196](https://github.com/owncloud/ocis/issues/4196)
* Enhancement - Deny access to resources: [#4903](https://github.com/owncloud/ocis/pull/4903)
* Enhancement - Improve validation of OIDC access tokens: [#3841](https://github.com/owncloud/ocis/issues/3841)
* Enhancement - Add /app/open-with-web endpoint: [#4376](https://github.com/owncloud/ocis/pull/4376)
* Enhancement - Add previewFileMimeTypes to web default config: [#4414](https://github.com/owncloud/ocis/pull/4414)
* Enhancement - Added language option to the app provider: [#4399](https://github.com/owncloud/ocis/pull/4399)
* Enhancement - Improve error log for "could not get user by claim" error: [#4227](https://github.com/owncloud/ocis/pull/4227)
* Enhancement - Improve login screen design: [#4500](https://github.com/owncloud/ocis/pull/4500)
* Enhancement - Add configuration options for mail authentication and encryption: [#4443](https://github.com/owncloud/ocis/pull/4443)
* Enhancement - Introduce service registry cache: [#3833](https://github.com/owncloud/ocis/pull/3833)
* Enhancement - Reintroduce user autoprovisioning in proxy: [#3860](https://github.com/owncloud/ocis/pull/3860)
* Enhancement - Allow to configure applications in Web: [#4578](https://github.com/owncloud/ocis/pull/4578)
* Enhancement - Added command to reset administrator password: [#4084](https://github.com/owncloud/ocis/issues/4084)
* Enhancement - Disable the color logging in docker compose examples: [#871](https://github.com/owncloud/ocis/issues/871)
* Enhancement - Allow providing list of services NOT to start: [#4254](https://github.com/owncloud/ocis/pull/4254)
* Enhancement - Introduce insecure flag for smtp email notifications: [#4279](https://github.com/owncloud/ocis/pull/4279)
* Enhancement - Optional events in graph service: [#55555](https://github.com/owncloud/ocis/pull/55555)
* Enhancement - Fix behavior for foobar (in present tense): [#4346](https://github.com/owncloud/ocis/pull/4346)
* Enhancement - Add the "hidden" state to the search index: [#5018](https://github.com/owncloud/ocis/pull/5018)
* Enhancement - Restrict admins from self-removal: [#3713](https://github.com/owncloud/ocis/issues/3713)
* Enhancement - OCS get share now also handle received shares: [#4322](https://github.com/owncloud/ocis/issues/4322)
* Enhancement - Add config option to provide TLS certificate: [#3818](https://github.com/owncloud/ocis/issues/3818)
* Enhancement - Add descriptions for graph-explorer config: [#3759](https://github.com/owncloud/ocis/pull/3759)
* Enhancement - Add /me/changePassword endpoint to GraphAPI: [#3063](https://github.com/owncloud/ocis/issues/3063)
* Enhancement - Allow to setup TLS for grpc services: [#4798](https://github.com/owncloud/ocis/pull/4798)
* Enhancement - Generate signing key and encryption secret: [#3909](https://github.com/owncloud/ocis/issues/3909)
* Enhancement - Update IdP UI: [#3493](https://github.com/owncloud/ocis/issues/3493)
* Enhancement - Logging improvements: [#4815](https://github.com/owncloud/ocis/pull/4815)
* Enhancement - Wrap metadata storage with dedicated reva gateway: [#3602](https://github.com/owncloud/ocis/pull/3602)
* Enhancement - New migrate command for migrating shares and public shares: [#3987](https://github.com/owncloud/ocis/pull/3987)
* Enhancement - Default to tls 1.2: [#4969](https://github.com/owncloud/ocis/pull/4969)
* Enhancement - Add missing unprotected paths: [#4454](https://github.com/owncloud/ocis/pull/4454)
* Enhancement - Secure the nats connection with TLS: [#4781](https://github.com/owncloud/ocis/pull/4781)
* Enhancement - Product field in OCS version: [#2918](https://github.com/owncloud/ocis/pull/2918)
* Enhancement - Automatically orientate photos when generating thumbnails: [#4477](https://github.com/owncloud/ocis/issues/4477)
* Enhancement - Refactor extensions to services: [#3980](https://github.com/owncloud/ocis/pull/3980)
* Enhancement - Refactor the proxy service: [#4401](https://github.com/owncloud/ocis/issues/4401)
* Enhancement - Remove windows from ci & release makefile: [#5026](https://github.com/owncloud/ocis/pull/5026)
* Enhancement - Rename AUTH_BASIC_AUTH_PROVIDER envvar: [#4966](https://github.com/owncloud/ocis/pull/4966)
* Enhancement - Report parent id: [#4757](https://github.com/owncloud/ocis/pull/4757)
* Enhancement - Allow resharing: [#3904](https://github.com/owncloud/ocis/pull/3904)
* Enhancement - Rewrite of the request authentication middleware: [#4374](https://github.com/owncloud/ocis/pull/4374)
* Enhancement - Add initial version of the search extensions: [#3635](https://github.com/owncloud/ocis/pull/3635)
* Enhancement - Prohibit users from setting or listing other user's values: [#4897](https://github.com/owncloud/ocis/pull/4897)
* Enhancement - Add capability for public link single file edit: [#6787](https://github.com/owncloud/web/pull/6787)
* Enhancement - Added `share_jail` and `projects` feature flags in spaces capability: [#3626](https://github.com/owncloud/ocis/pull/3626)
* Enhancement - Use storageID when requesting special items: [#4356](https://github.com/owncloud/ocis/pull/4356)
* Enhancement - Add description tags to the thumbnails config structs: [#3752](https://github.com/owncloud/ocis/pull/3752)
* Enhancement - Make thumbnails service log less noisy: [#3959](https://github.com/owncloud/ocis/pull/3959)
* Enhancement - Add thumbnails support for tiff and bmp files: [#4634](https://github.com/owncloud/ocis/pull/4634)
* Enhancement - Update linkshare capabilities: [#3579](https://github.com/owncloud/ocis/pull/3579)
* Enhancement - Update reva: [#3944](https://github.com/owncloud/ocis/pull/3944)
* Enhancement - Update reva to version 2.7.2: [#4115](https://github.com/owncloud/ocis/pull/4115)
* Enhancement - Update reva to v2.7.4: [#4294](https://github.com/owncloud/ocis/pull/4294)
* Enhancement - Update reva to v2.8.0: [#4444](https://github.com/owncloud/ocis/pull/4444)
* Enhancement - Update reva to version 2.4.1: [#3746](https://github.com/owncloud/ocis/pull/3746)
* Enhancement - Update reva to version 2.5.1: [#3932](https://github.com/owncloud/ocis/pull/3932)
* Enhancement - Update Reva to version 2.10.0: [#4522](https://github.com/owncloud/ocis/pull/4522)
* Enhancement - Update reva to version 2.11.0: [#4588](https://github.com/owncloud/ocis/pull/4588)
* Enhancement - Update reva to v2.3.1: [#3552](https://github.com/owncloud/ocis/pull/3552)
* Enhancement - Update ownCloud Web to v5.5.0-rc.8: [#6854](https://github.com/owncloud/web/pull/6854)
* Enhancement - Update ownCloud Web to v5.5.0-rc.9: [#6854](https://github.com/owncloud/web/pull/6854)
* Enhancement - Update ownCloud Web to v5.5.0-rc.6: [#6854](https://github.com/owncloud/web/pull/6854)
* Enhancement - Update ownCloud Web to v5.7.0-rc.1: [#4005](https://github.com/owncloud/ocis/pull/4005)
* Enhancement - Update ownCloud Web to v6.0.0: [#5153](https://github.com/owncloud/ocis/pull/5153)
* Enhancement - Update ownCloud Web to v5.7.0-rc.4: [#4140](https://github.com/owncloud/ocis/pull/4140)
* Enhancement - Update ownCloud Web to v5.7.0-rc.8: [#4314](https://github.com/owncloud/ocis/pull/4314)
* Enhancement - Update ownCloud Web to v5.7.0-rc.10: [#4439](https://github.com/owncloud/ocis/pull/4439)
* Enhancement - Update ownCloud Web to v5.7.0: [#4508](https://github.com/owncloud/ocis/pull/4508)
* Enhancement - Expand personal drive on the graph user: [#4357](https://github.com/owncloud/ocis/pull/4357)
* Enhancement - Validate space names: [#4955](https://github.com/owncloud/ocis/pull/4955)
* Enhancement - Add descriptions to webdav configuration: [#3755](https://github.com/owncloud/ocis/pull/3755)
* Enhancement - Search service at the old webdav endpoint: [#4118](https://github.com/owncloud/ocis/pull/4118)
* Enhancement - Make it possible to configure a WOPI folderurl: [#4716](https://github.com/owncloud/ocis/pull/4716)

# Changelog for [1.20.0] (2022-04-13)

The following sections list the changes for 1.20.0.

[1.20.0]: https://github.com/owncloud/ocis/compare/v1.19.1...v1.20.0

## Summary

* Bugfix - Add `owncloudsql` driver to authprovider config: [#3435](https://github.com/owncloud/ocis/pull/3435)
* Bugfix - Corrected documentation: [#3439](https://github.com/owncloud/ocis/pull/3439)
* Bugfix - Ensure the same data on /ocs/v?.php/config like oC10: [#3113](https://github.com/owncloud/ocis/pull/3113)
* Bugfix - Use the default server download protocol if spaces are not supported: [#3386](https://github.com/owncloud/ocis/pull/3386)
* Change - Fix keys with underscores in the config files: [#3412](https://github.com/owncloud/ocis/pull/3412)
* Change - Don't create demo users by default: [#3474](https://github.com/owncloud/ocis/pull/3474)
* Enhancement - Alias links: [#3454](https://github.com/owncloud/ocis/pull/3454)
* Enhancement - Replace deprecated String.prototype.substr(): [#3448](https://github.com/owncloud/ocis/pull/3448)
* Enhancement - Add sorting to GraphAPI users and groups: [#3360](https://github.com/owncloud/ocis/issues/3360)
* Enhancement - Unify LDAP config settings accross services: [#3476](https://github.com/owncloud/ocis/pull/3476)
* Enhancement - Make config dir configurable: [#3440](https://github.com/owncloud/ocis/pull/3440)
* Enhancement - Use embeddable ocdav go micro service: [#3397](https://github.com/owncloud/ocis/pull/3397)
* Enhancement - Update reva to v2.2.0: [#3397](https://github.com/owncloud/ocis/pull/3397)
* Enhancement - Update ownCloud Web to v5.4.0: [#6709](https://github.com/owncloud/web/pull/6709)
* Enhancement - Implement audit events for user and groups: [#3467](https://github.com/owncloud/ocis/pull/3467)

# Changelog for [1.19.1] (2022-03-29)

The following sections list the changes for 1.19.1.

[1.19.1]: https://github.com/owncloud/ocis/compare/v1.19.0...v1.19.1

## Summary

* Bugfix - Return correct special item urls: [#3419](https://github.com/owncloud/ocis/pull/3419)

# Changelog for [1.19.0] (2022-03-29)

The following sections list the changes for 1.19.0.

[1.19.0]: https://github.com/owncloud/ocis/compare/v1.18.0...v1.19.0

## Summary

* Bugfix - Network configuration in individiual_services example: [#3238](https://github.com/owncloud/ocis/pull/3238)
* Bugfix - Improve gif thumbnails: [#3305](https://github.com/owncloud/ocis/pull/3305)
* Bugfix - Fix error handling in GraphAPI GetUsers call: [#3357](https://github.com/owncloud/ocis/pull/3357)
* Bugfix - Fix request validation on GraphAPI User updates: [#3167](https://github.com/owncloud/ocis/issues/3167)
* Bugfix - Replace public mountpoint fileid with grant fileid: [#3349](https://github.com/owncloud/ocis/pull/3349)
* Change - Add remote item to mountpoint and fix spaceID: [#3365](https://github.com/owncloud/ocis/pull/3365)
* Change - Switch NATS backend: [#3192](https://github.com/owncloud/ocis/pull/3192)
* Change - Drop json config file support: [#3366](https://github.com/owncloud/ocis/pull/3366)
* Change - Settings service now stores its data via metadata service: [#3232](https://github.com/owncloud/ocis/pull/3232)
* Enhancement - Audit logger will now log file events: [#3332](https://github.com/owncloud/ocis/pull/3332)
* Enhancement - Add password reset link to login page: [#3329](https://github.com/owncloud/ocis/pull/3329)
* Enhancement - Log sharing events in audit service: [#3301](https://github.com/owncloud/ocis/pull/3301)
* Enhancement - Add space aliases: [#3283](https://github.com/owncloud/ocis/pull/3283)
* Enhancement - Include etags in drives listing: [#3267](https://github.com/owncloud/ocis/pull/3267)
* Enhancement - Improve thumbnails API: [#3272](https://github.com/owncloud/ocis/pull/3272)
* Enhancement - Update reva to v2.1.0: [#3330](https://github.com/owncloud/ocis/pull/3330)
* Enhancement - Update ownCloud Web to v5.3.0: [#6561](https://github.com/owncloud/web/pull/6561)

# Changelog for [1.18.0] (2022-03-03)

The following sections list the changes for 1.18.0.

[1.18.0]: https://github.com/owncloud/ocis/compare/v1.17.0...v1.18.0

## Summary

* Bugfix - Capabilities for password protected public links: [#3229](https://github.com/owncloud/ocis/pull/3229)
* Bugfix - Make events settings configurable: [#3214](https://github.com/owncloud/ocis/pull/3214)
* Bugfix - Align storage metadata GPRC bind port with other variable names: [#3169](https://github.com/owncloud/ocis/pull/3169)
* Change - Unify file IDs: [#3185](https://github.com/owncloud/ocis/pull/3185)
* Enhancement - Add sorting to list Spaces: [#3200](https://github.com/owncloud/ocis/issues/3200)
* Enhancement - Change NATS port: [#3210](https://github.com/owncloud/ocis/pull/3210)
* Enhancement - Re-Enabling web cache control: [#3109](https://github.com/owncloud/ocis/pull/3109)
* Enhancement - Add SPA conform fileserver for web: [#3109](https://github.com/owncloud/ocis/pull/3109)
* Enhancement - Implement notifications service: [#3217](https://github.com/owncloud/ocis/pull/3217)
* Enhancement - Thumbnails in spaces: [#3219](https://github.com/owncloud/ocis/pull/3219)
* Enhancement - Update reva to v2.0.0: [#3231](https://github.com/owncloud/ocis/pull/3231)
* Enhancement - Update ownCloud Web to v5.2.0: [#6506](https://github.com/owncloud/web/pull/6506)

# Changelog for [1.17.0] (2022-02-16)

The following sections list the changes for 1.17.0.

[1.17.0]: https://github.com/owncloud/ocis/compare/v1.16.0...v1.17.0

## Summary

* Bugfix - Add `ocis storage-auth-machine` subcommand: [#2910](https://github.com/owncloud/ocis/pull/2910)
* Bugfix - Use same jwt secret for accounts as for metadata storage: [#3081](https://github.com/owncloud/ocis/pull/3081)
* Bugfix - Make the default grpc client use the registry settings: [#3041](https://github.com/owncloud/ocis/pull/3041)
* Bugfix - Remove group memberships when deleting a user: [#3027](https://github.com/owncloud/ocis/issues/3027)
* Bugfix - Fix retry handling for LDAP connections: [#2974](https://github.com/owncloud/ocis/issues/2974)
* Bugfix - Fix the default tracing provider: [#2952](https://github.com/owncloud/ocis/pull/2952)
* Bugfix - Fix configuration for space membership endpoint: [#2893](https://github.com/owncloud/ocis/pull/2893)
* Change - Change log level default from debug to error: [#3071](https://github.com/owncloud/ocis/pull/3071)
* Change - Remove the ownCloud storage driver: [#3072](https://github.com/owncloud/ocis/pull/3072)
* Change - Unify configuration and commands: [#2818](https://github.com/owncloud/ocis/pull/2818)
* Change - Functionality to restore spaces: [#3092](https://github.com/owncloud/ocis/pull/3092)
* Change - Extended Space Properties: [#3141](https://github.com/owncloud/ocis/pull/3141)
* Change - Update the graph api: [#2885](https://github.com/owncloud/ocis/pull/2885)
* Change - Update libre-graph-api to v0.3.0: [#2858](https://github.com/owncloud/ocis/pull/2858)
* Change - Return not found when updating non existent space: [#2869](https://github.com/cs3org/reva/pull/2869)
* Enhancement - Provide Description when creating a space: [#3167](https://github.com/owncloud/ocis/pull/3167)
* Enhancement - Add graph endpoint to delete and purge spaces: [#2979](https://github.com/owncloud/ocis/pull/2979)
* Enhancement - Add permissions to graph drives: [#3095](https://github.com/owncloud/ocis/pull/3095)
* Enhancement - Add new file url of the app provider to the ocs capabilities: [#2884](https://github.com/owncloud/ocis/pull/2884)
* Enhancement - Add spaces capability: [#2931](https://github.com/owncloud/ocis/pull/2931)
* Enhancement - Consul as supported service registry: [#3133](https://github.com/owncloud/ocis/pull/3133)
* Enhancement - Introduce User and Group Management capabilities on GraphAPI: [#2947](https://github.com/owncloud/ocis/pull/2947)
* Enhancement - Support signature auth in the public share auth middleware: [#2831](https://github.com/owncloud/ocis/pull/2831)
* Enhancement - Update REVA to v1.16.1-0.20220112085026-07451f6cd806: [#2953](https://github.com/owncloud/ocis/pull/2953)
* Enhancement - Add endpoint to retrieve a single space: [#2978](https://github.com/owncloud/ocis/pull/2978)
* Enhancement - Add filter by driveType and id to /me/drives: [#2946](https://github.com/owncloud/ocis/pull/2946)
* Enhancement - Update REVA to v1.16.1-0.20220215130802-df1264deff58: [#2878](https://github.com/owncloud/ocis/pull/2878)
* Enhancement - Update ownCloud Web to v5.0.0: [#2895](https://github.com/owncloud/ocis/pull/2895)

# Changelog for [1.16.0] (2021-12-10)

The following sections list the changes for 1.16.0.

[1.16.0]: https://github.com/owncloud/ocis/compare/v1.15.0...v1.16.0

## Summary

* Bugfix - Fix claim selector based routing for basic auth: [#2779](https://github.com/owncloud/ocis/pull/2779)
* Bugfix - Disallow creation of a group with empty name via the OCS api: [#2825](https://github.com/owncloud/ocis/pull/2825)
* Bugfix - Fix using s3ng as the metadata storage backend: [#2807](https://github.com/owncloud/ocis/pull/2807)
* Bugfix - Use the CS3api up- and download workflow for the accounts service: [#2837](https://github.com/owncloud/ocis/pull/2837)
* Change - Rename `APP_PROVIDER_BASIC_*` environment variables: [#2812](https://github.com/owncloud/ocis/pull/2812)
* Change - Restructure Configuration Parsing: [#2708](https://github.com/owncloud/ocis/pull/2708)
* Change - OIDC: fallback if IDP doesn't provide "preferred_username" claim: [#2644](https://github.com/owncloud/ocis/issues/2644)
* Enhancement - Cleanup ocis-pkg config: [#2813](https://github.com/owncloud/ocis/pull/2813)
* Enhancement - Correct shutdown of services under runtime: [#2843](https://github.com/owncloud/ocis/pull/2843)
* Enhancement - Update REVA to v1.17.0: [#2849](https://github.com/owncloud/ocis/pull/2849)
* Enhancement - Update ownCloud Web to v4.6.1: [#2846](https://github.com/owncloud/ocis/pull/2846)

# Changelog for [1.15.0] (2021-11-19)

The following sections list the changes for 1.15.0.

[1.15.0]: https://github.com/owncloud/ocis/compare/v1.14.0...v1.15.0

## Summary

* Bugfix - Don't allow empty password: [#197](https://github.com/owncloud/product/issues/197)
* Bugfix - Fix basic auth config: [#2719](https://github.com/owncloud/ocis/pull/2719)
* Bugfix - Fix basic auth with custom user claim: [#2755](https://github.com/owncloud/ocis/pull/2755)
* Bugfix - Fix oCIS startup ony systems with IPv6: [#2698](https://github.com/owncloud/ocis/pull/2698)
* Bugfix - Fix opening images in media viewer for some usernames: [#2738](https://github.com/owncloud/ocis/pull/2738)
* Bugfix - Fix error logging when there is no thumbnail for a file: [#2702](https://github.com/owncloud/ocis/pull/2702)
* Bugfix - Don't announce resharing via capabilities: [#2690](https://github.com/owncloud/ocis/pull/2690)
* Change - Make all insecure options configurable and change the default to false: [#2700](https://github.com/owncloud/ocis/issues/2700)
* Change - Update ownCloud Web to v4.5.0: [#2780](https://github.com/owncloud/ocis/pull/2780)
* Enhancement - Add API to list all spaces: [#2692](https://github.com/owncloud/ocis/pull/2692)
* Enhancement - Update REVA to v1.16.0: [#2737](https://github.com/owncloud/ocis/pull/2737)

# Changelog for [1.14.0] (2021-10-27)

The following sections list the changes for 1.14.0.

[1.14.0]: https://github.com/owncloud/ocis/compare/v1.13.0...v1.14.0

## Summary

* Security - Don't expose services by default: [#2612](https://github.com/owncloud/ocis/issues/2612)
* Bugfix - Create parent directories for idp configuration: [#2667](https://github.com/owncloud/ocis/issues/2667)
* Change - Configurable default quota: [#2621](https://github.com/owncloud/ocis/issues/2621)
* Change - New default data paths and easier configuration of the data path: [#2590](https://github.com/owncloud/ocis/pull/2590)
* Change - Split spaces webdav url and graph url in base and path: [#2660](https://github.com/owncloud/ocis/pull/2660)
* Change - Update ownCloud Web to v4.4.0: [#2681](https://github.com/owncloud/ocis/pull/2681)
* Enhancement - Add user setting capability: [#2655](https://github.com/owncloud/ocis/pull/2655)
* Enhancement - Broaden bufbuild/Buf usage: [#2630](https://github.com/owncloud/ocis/pull/2630)
* Enhancement - Replace fileb0x with go-embed: [#1199](https://github.com/owncloud/ocis/issues/1199)
* Enhancement - Upgrade to go-micro v4.1.0: [#2616](https://github.com/owncloud/ocis/pull/2616)
* Enhancement - Review and correct http header: [#2666](https://github.com/owncloud/ocis/pull/2666)
* Enhancement - Lower TUS max chunk size: [#2584](https://github.com/owncloud/ocis/pull/2584)
* Enhancement - Add sharees additional info paramater config to ocs: [#2637](https://github.com/owncloud/ocis/pull/2637)
* Enhancement - Add a middleware to authenticate public share requests: [#2536](https://github.com/owncloud/ocis/pull/2536)
* Enhancement - Report quota states: [#2628](https://github.com/owncloud/ocis/pull/2628)
* Enhancement - Start up a new machine auth provider in the storage service: [#2528](https://github.com/owncloud/ocis/pull/2528)
* Enhancement - Enforce permission on update space quota: [#2650](https://github.com/owncloud/ocis/pull/2650)
* Enhancement - Update lico to v0.51.1: [#2654](https://github.com/owncloud/ocis/pull/2654)
* Enhancement - Update reva to v1.15: [#2658](https://github.com/owncloud/ocis/pull/2658)

# Changelog for [1.13.0] (2021-10-13)

The following sections list the changes for 1.13.0.

[1.13.0]: https://github.com/owncloud/ocis/compare/v1.12.0...v1.13.0

## Summary

* Bugfix - Fix the account resolver middleware: [#2557](https://github.com/owncloud/ocis/pull/2557)
* Bugfix - Fix version information for extensions: [#2575](https://github.com/owncloud/ocis/pull/2575)
* Bugfix - Add the gatewaysvc to all shared configuration in REVA services: [#2597](https://github.com/owncloud/ocis/pull/2597)
* Bugfix - Use proper url path decode on the username: [#2511](https://github.com/owncloud/ocis/pull/2511)
* Bugfix - Remove notifications placeholder: [#2514](https://github.com/owncloud/ocis/pull/2514)
* Bugfix - Remove asset path configuration option from proxy: [#2576](https://github.com/owncloud/ocis/pull/2576)
* Bugfix - Race condition in config parsing: [#2574](https://github.com/owncloud/ocis/pull/2574)
* Change - Configure users and metadata storage separately: [#2598](https://github.com/owncloud/ocis/pull/2598)
* Change - Make the drives create method odata compliant: [#2531](https://github.com/owncloud/ocis/pull/2531)
* Change - Unify Envvar names configuring REVA gateway address: [#2587](https://github.com/owncloud/ocis/pull/2587)
* Change - Update ownCloud Web to v4.3.0: [#2589](https://github.com/owncloud/ocis/pull/2589)
* Enhancement - Updated MimeTypes configuration for AppRegistry: [#2603](https://github.com/owncloud/ocis/pull/2603)
* Enhancement - Add maximum files and size to archiver capabilities: [#2544](https://github.com/owncloud/ocis/pull/2544)
* Enhancement - Reduced repository size: [#2579](https://github.com/owncloud/ocis/pull/2579)
* Enhancement - Return the newly created space: [#2610](https://github.com/owncloud/ocis/pull/2610)
* Enhancement - Expose the reva archiver in OCIS: [#2509](https://github.com/owncloud/ocis/pull/2509)
* Enhancement - Favorites capability: [#2599](https://github.com/owncloud/ocis/pull/2599)
* Enhancement - Upgrade to GO 1.17: [#2605](https://github.com/owncloud/ocis/pull/2605)
* Enhancement - Make mimetype allow list configurable for app provider: [#2553](https://github.com/owncloud/ocis/pull/2553)
* Enhancement - Add allow_creation parameter to mime type config: [#2591](https://github.com/owncloud/ocis/pull/2591)
* Enhancement - Add option to skip generation of demo users and groups: [#2495](https://github.com/owncloud/ocis/pull/2495)
* Enhancement - Allow overriding the cookie based route by claim: [#2508](https://github.com/owncloud/ocis/pull/2508)
* Enhancement - Redirect invalid links to oC Web: [#2493](https://github.com/owncloud/ocis/pull/2493)
* Enhancement - Use reva's Authenticate method instead of spawning token managers: [#2528](https://github.com/owncloud/ocis/pull/2528)
* Enhancement - TLS config options for ldap in reva: [#2492](https://github.com/owncloud/ocis/pull/2492)
* Enhancement - Set reva JWT token expiration time to 24 hours by default: [#2527](https://github.com/owncloud/ocis/pull/2527)
* Enhancement - Update reva to v1.14.0: [#2615](https://github.com/owncloud/ocis/pull/2615)

# Changelog for [1.12.0] (2021-09-14)

The following sections list the changes for 1.12.0.

[1.12.0]: https://github.com/owncloud/ocis/compare/v1.11.0...v1.12.0

## Summary

* Bugfix - Remove non working proxy route and fix cs3 users example: [#2474](https://github.com/owncloud/ocis/pull/2474)
* Bugfix - Set English as default language in the dropdown in the settings page: [#2465](https://github.com/owncloud/ocis/pull/2465)
* Change - Remove OnlyOffice extension: [#2433](https://github.com/owncloud/ocis/pull/2433)
* Change - Remove OnlyOffice extension: [#2433](https://github.com/owncloud/ocis/pull/2433)
* Change - Update ownCloud Web to v4.2.0: [#2501](https://github.com/owncloud/ocis/pull/2501)
* Enhancement - Add app provider and app provider registry: [#2204](https://github.com/owncloud/ocis/pull/2204)
* Enhancement - Add the create space permission: [#2461](https://github.com/owncloud/ocis/pull/2461)
* Enhancement - Add set space quota permission: [#2459](https://github.com/owncloud/ocis/pull/2459)
* Enhancement - Create a Space using the Graph API: [#2471](https://github.com/owncloud/ocis/pull/2471)
* Enhancement - Update go-chi/chi to version 5.0.3: [#2429](https://github.com/owncloud/ocis/pull/2429)
* Enhancement - Upgrade go micro to v3.6.0: [#2451](https://github.com/owncloud/ocis/pull/2451)
* Enhancement - Update reva to v1.13.0: [#2477](https://github.com/owncloud/ocis/pull/2477)

# Changelog for [1.11.0] (2021-08-24)

The following sections list the changes for 1.11.0.

[1.11.0]: https://github.com/owncloud/ocis/compare/v1.10.0...v1.11.0

## Summary

* Bugfix - Specify primary user type for all accounts: [#2364](https://github.com/owncloud/ocis/pull/2364)
* Bugfix - Fix naming of the user- and groupprovider services: [#2388](https://github.com/owncloud/ocis/pull/2388)
* Change - Update ownCloud Web to v4.1.0: [#2426](https://github.com/owncloud/ocis/pull/2426)
* Enhancement - Use non root user for the owncloud/ocis docker image: [#2380](https://github.com/owncloud/ocis/pull/2380)
* Enhancement - Replace unmaintained jwt library: [#2386](https://github.com/owncloud/ocis/pull/2386)
* Enhancement - Update bleve to version 2.1.0: [#2391](https://github.com/owncloud/ocis/pull/2391)
* Enhancement - Update github.com/coreos/go-oidc to v3.0.0: [#2393](https://github.com/owncloud/ocis/pull/2393)
* Enhancement - Update reva to v1.12: [#2423](https://github.com/owncloud/ocis/pull/2423)

# Changelog for [1.10.0] (2021-08-06)

The following sections list the changes for 1.10.0.

[1.10.0]: https://github.com/owncloud/ocis/compare/v1.9.0...v1.10.0

## Summary

* Bugfix - Improve IDP Login Accessibility: [#5376](https://github.com/owncloud/web/issues/5376)
* Bugfix - Forward basic auth to OpenID connect token authentication endpoint: [#2095](https://github.com/owncloud/ocis/issues/2095)
* Bugfix - Log all requests in the proxy access log: [#2301](https://github.com/owncloud/ocis/pull/2301)
* Bugfix - Update glauth to 20210729125545-b9aecdfcac31: [#2336](https://github.com/owncloud/ocis/pull/2336)
* Change - Update ownCloud Web to v4.0.0: [#2353](https://github.com/owncloud/ocis/pull/2353)
* Enhancement - Proxy: Add claims policy selector: [#2248](https://github.com/owncloud/ocis/pull/2248)
* Enhancement - Add ocs cache warmup config and warn on protobuf ns conflicts: [#2328](https://github.com/owncloud/ocis/pull/2328)
* Enhancement - Refactor graph API: [#2277](https://github.com/owncloud/ocis/pull/2277)
* Enhancement - Update REVA: [#2355](https://github.com/owncloud/ocis/pull/2355)
* Enhancement - Use only one go.mod file for project dependencies: [#2344](https://github.com/owncloud/ocis/pull/2344)

# Changelog for [1.9.0] (2021-07-13)

The following sections list the changes for 1.9.0.

[1.9.0]: https://github.com/owncloud/ocis/compare/v1.8.0...v1.9.0

## Summary

* Bugfix - Panic when service fails to start: [#2252](https://github.com/owncloud/ocis/pull/2252)
* Bugfix - Dont use port 80 as debug for GroupsProvider: [#2271](https://github.com/owncloud/ocis/pull/2271)
* Change - Update ownCloud Web to v3.4.0: [#2276](https://github.com/owncloud/ocis/pull/2276)
* Change - Update WEB to v3.4.1: [#2283](https://github.com/owncloud/ocis/pull/2283)
* Enhancement - Runtime support for cherry picking extensions: [#2229](https://github.com/owncloud/ocis/pull/2229)
* Enhancement - Add readonly mode for storagehome and storageusers: [#2230](https://github.com/owncloud/ocis/pull/2230)
* Enhancement - Remove unnecessary Service.Init(): [#1705](https://github.com/owncloud/ocis/pull/1705)
* Enhancement - Update REVA to v1.9.1-0.20210628143859-9d29c36c0c3f: [#2227](https://github.com/owncloud/ocis/pull/2227)
* Enhancement - Update REVA to v1.9.1: [#2280](https://github.com/owncloud/ocis/pull/2280)

# Changelog for [1.8.0] (2021-06-28)

The following sections list the changes for 1.8.0.

[1.8.0]: https://github.com/owncloud/ocis/compare/v1.7.0...v1.8.0

## Summary

* Bugfix - External storage registration used wrong config: [#2120](https://github.com/owncloud/ocis/pull/2120)
* Bugfix - Remove authentication from /status.php completely: [#2188](https://github.com/owncloud/ocis/pull/2188)
* Bugfix - Make webdav namespace configurable across services: [#2198](https://github.com/owncloud/ocis/pull/2198)
* Change - Update ownCloud Web to v3.3.0: [#2187](https://github.com/owncloud/ocis/pull/2187)
* Enhancement - Properly configure graph-explorer client registration: [#2118](https://github.com/owncloud/ocis/pull/2118)
* Enhancement - Use system default location to store TLS artefacts: [#2129](https://github.com/owncloud/ocis/pull/2129)
* Enhancement - Update REVA to v1.9: [#2205](https://github.com/owncloud/ocis/pull/2205)

# Changelog for [1.7.0] (2021-06-04)

The following sections list the changes for 1.7.0.

[1.7.0]: https://github.com/owncloud/ocis/compare/v1.6.0...v1.7.0

## Summary

* Bugfix - Change the groups index to be case sensitive: [#2109](https://github.com/owncloud/ocis/pull/2109)
* Change - Update ownCloud Web to v3.2.0: [#2096](https://github.com/owncloud/ocis/pull/2096)
* Enhancement - Enable the s3ng storage driver: [#1886](https://github.com/owncloud/ocis/pull/1886)
* Enhancement - Color contrasts on IDP/OIDC login pages: [#2088](https://github.com/owncloud/ocis/pull/2088)
* Enhancement - Announce user profile picture capability: [#2036](https://github.com/owncloud/ocis/pull/2036)
* Enhancement - Update reva to v1.7.1-0.20210531093513-b74a2b156af6: [#2104](https://github.com/owncloud/ocis/pull/2104)

# Changelog for [1.6.0] (2021-05-12)

The following sections list the changes for 1.6.0.

[1.6.0]: https://github.com/owncloud/ocis/compare/v1.5.0...v1.6.0

## Summary

* Bugfix - Fix STORAGE_METADATA_ROOT default value override: [#1956](https://github.com/owncloud/ocis/pull/1956)
* Bugfix - Stop the supervisor if a service fails to start: [#1963](https://github.com/owncloud/ocis/pull/1963)
* Change - Update ownCloud Web to v3.1.0: [#2045](https://github.com/owncloud/ocis/pull/2045)
* Enhancement - Added dictionary files: [#2003](https://github.com/owncloud/ocis/pull/2003)
* Enhancement - Introduce login form with h1 tag for screen readers only: [#1991](https://github.com/owncloud/ocis/pull/1991)
* Enhancement - User Deprovisioning for the OCS API: [#1962](https://github.com/owncloud/ocis/pull/1962)
* Enhancement - Support thumbnails for txt files: [#1988](https://github.com/owncloud/ocis/pull/1988)
* Enhancement - Update reva to v1.7.1-0.20210430154404-69bd21f2cc97: [#2010](https://github.com/owncloud/ocis/pull/2010)
* Enhancement - Update reva to v1.7.1-0.20210507160327-e2c3841d0dbc: [#2044](https://github.com/owncloud/ocis/pull/2044)
* Enhancement - Use oc-select: [#1979](https://github.com/owncloud/ocis/pull/1979)
* Enhancement - Set SameSite settings to Strict for Web: [#2019](https://github.com/owncloud/ocis/pull/2019)

# Changelog for [1.5.0] (2021-04-21)

The following sections list the changes for 1.5.0.

[1.5.0]: https://github.com/owncloud/ocis/compare/v1.4.0...v1.5.0

## Summary

* Bugfix - Fixes "unaligned 64-bit atomic operation" panic on 32-bit ARM: [#1888](https://github.com/owncloud/ocis/pull/1888)
* Change - Make Protobuf package names unique: [#1875](https://github.com/owncloud/ocis/pull/1875)
* Change - Update ownCloud Web to v3.0.0: [#1938](https://github.com/owncloud/ocis/pull/1938)
* Enhancement - Change default path for thumbnails: [#1892](https://github.com/owncloud/ocis/pull/1892)
* Enhancement - Parse config on supervised mode with run subcommand: [#1931](https://github.com/owncloud/ocis/pull/1931)
* Enhancement - Update ODS in accounts & settings extension: [#1934](https://github.com/owncloud/ocis/pull/1934)
* Enhancement - Add config for public share SQL driver: [#1916](https://github.com/owncloud/ocis/pull/1916)
* Enhancement - Remove dead runtime code: [#1923](https://github.com/owncloud/ocis/pull/1923)
* Enhancement - Add option to reading registry rules from json file: [#1917](https://github.com/owncloud/ocis/pull/1917)
* Enhancement - Update reva to v1.6.1-0.20210414111318-a4b5148cbfb2: [#1872](https://github.com/owncloud/ocis/pull/1872)

# Changelog for [1.4.0] (2021-03-30)

The following sections list the changes for 1.4.0.

[1.4.0]: https://github.com/owncloud/ocis/compare/v1.3.0...v1.4.0

## Summary

* Bugfix - Fix thumbnail generation for jpegs: [#1785](https://github.com/owncloud/ocis/pull/1785)
* Change - Update ownCloud Web to v2.1.0: [#1870](https://github.com/owncloud/ocis/pull/1870)
* Enhancement - Add focus to input elements on login page: [#1792](https://github.com/owncloud/ocis/pull/1792)
* Enhancement - Improve accessibility to input elements on login page: [#1794](https://github.com/owncloud/ocis/pull/1794)
* Enhancement - Add new build targets: [#1824](https://github.com/owncloud/ocis/pull/1824)
* Enhancement - Clarify expected failures: [#1790](https://github.com/owncloud/ocis/pull/1790)
* Enhancement - Replace special character in login page title with a regular minus: [#1813](https://github.com/owncloud/ocis/pull/1813)
* Enhancement - File Logging: [#1816](https://github.com/owncloud/ocis/pull/1816)
* Enhancement - Runtime Hostname and Port are now configurable: [#1822](https://github.com/owncloud/ocis/pull/1822)
* Enhancement - Generate thumbnails for .gif files: [#1791](https://github.com/owncloud/ocis/pull/1791)
* Enhancement - Tracing Refactor: [#1819](https://github.com/owncloud/ocis/pull/1819)
* Enhancement - Update reva to v1.6.1-0.20210326165326-e8a00d9b2368: [#1683](https://github.com/owncloud/ocis/pull/1683)

# Changelog for [1.3.0] (2021-03-09)

The following sections list the changes for 1.3.0.

[1.3.0]: https://github.com/owncloud/ocis/compare/v1.2.0...v1.3.0

## Summary

* Bugfix - Purposely delay accounts service startup: [#1734](https://github.com/owncloud/ocis/pull/1734)
* Bugfix - Add missing gateway config: [#1716](https://github.com/owncloud/ocis/pull/1716)
* Bugfix - Fix accounts initialization: [#1696](https://github.com/owncloud/ocis/pull/1696)
* Bugfix - Fix the ttl of the authentication middleware cache: [#1699](https://github.com/owncloud/ocis/pull/1699)
* Change - Update ownCloud Web to v2.0.1: [#1683](https://github.com/owncloud/ocis/pull/1683)
* Change - Update ownCloud Web to v2.0.2: [#1776](https://github.com/owncloud/ocis/pull/1776)
* Enhancement - Remove the JWT from the log: [#1758](https://github.com/owncloud/ocis/pull/1758)
* Enhancement - Update go-micro to v3.5.1-0.20210217182006-0f0ace1a44a9: [#1670](https://github.com/owncloud/ocis/pull/1670)
* Enhancement - Update reva to v1.6.1-0.20210223065028-53f39499762e: [#1683](https://github.com/owncloud/ocis/pull/1683)
* Enhancement - Add initial nats and kubernetes registry support: [#1697](https://github.com/owncloud/ocis/pull/1697)

# Changelog for [1.2.0] (2021-02-17)

The following sections list the changes for 1.2.0.

[1.2.0]: https://github.com/owncloud/ocis/compare/v1.1.0...v1.2.0

## Summary

* Bugfix - Check if roles are present in user object before looking those up: [#1388](https://github.com/owncloud/ocis/pull/1388)
* Bugfix - Fix etcd address configuration: [#1546](https://github.com/owncloud/ocis/pull/1546)
* Bugfix - Remove unimplemented config file option for oCIS root command: [#1636](https://github.com/owncloud/ocis/pull/1636)
* Bugfix - Fix thumbnail generation when using different idp: [#1624](https://github.com/owncloud/ocis/issues/1624)
* Change - Initial release of graph and graph explorer: [#1594](https://github.com/owncloud/ocis/pull/1594)
* Change - Move runtime code on refs/pman over to owncloud/ocis/ocis: [#1483](https://github.com/owncloud/ocis/pull/1483)
* Change - Update ownCloud Web to v2.0.0: [#1661](https://github.com/owncloud/ocis/pull/1661)
* Enhancement - Make use of new design-system oc-table: [#1597](https://github.com/owncloud/ocis/pull/1597)
* Enhancement - Use a default protocol parameter instead of explicitly disabling tus: [#1331](https://github.com/cs3org/reva/pull/1331)
* Enhancement - Functionality to map home directory to different storage providers: [#1186](https://github.com/owncloud/ocis/pull/1186)
* Enhancement - Introduce ADR: [#1042](https://github.com/owncloud/ocis/pull/1042)
* Enhancement - Switch to opencontainers annotation scheme: [#1381](https://github.com/owncloud/ocis/pull/1381)
* Enhancement - Migrate ocis-graph-explorer to ocis monorepo: [#1596](https://github.com/owncloud/ocis/pull/1596)
* Enhancement - Migrate ocis-graph to ocis monorepo: [#1594](https://github.com/owncloud/ocis/pull/1594)
* Enhancement - Enable group sharing and add config for sharing SQL driver: [#1626](https://github.com/owncloud/ocis/pull/1626)
* Enhancement - Update reva to v1.5.2-0.20210125114636-0c10b333ee69: [#1482](https://github.com/owncloud/ocis/pull/1482)

# Changelog for [1.1.0] (2021-01-22)

The following sections list the changes for 1.1.0.

[1.1.0]: https://github.com/owncloud/ocis/compare/v1.0.0...v1.1.0

## Summary

* Change - Disable pretty logging by default: [#1133](https://github.com/owncloud/ocis/pull/1133)
* Change - Add "volume" declaration to docker images: [#1375](https://github.com/owncloud/ocis/pull/1375)
* Change - Add "expose" information to docker images: [#1366](https://github.com/owncloud/ocis/pull/1366)
* Change - Generate cryptographically secure state token: [#1203](https://github.com/owncloud/ocis/pull/1203)
* Change - Move k6 to cdperf: [#1358](https://github.com/owncloud/ocis/pull/1358)
* Change - Update go version: [#1364](https://github.com/owncloud/ocis/pull/1364)
* Change - Update ownCloud Web to v1.0.1: [#1191](https://github.com/owncloud/ocis/pull/1191)
* Enhancement - Add OCIS_URL env var: [#1148](https://github.com/owncloud/ocis/pull/1148)
* Enhancement - Use sync.cache for roles cache: [#1367](https://github.com/owncloud/ocis/pull/1367)
* Enhancement - Add named locks and refactor cache: [#1212](https://github.com/owncloud/ocis/pull/1212)
* Enhancement - Update reva to v1.5.1: [#1372](https://github.com/owncloud/ocis/pull/1372)
* Enhancement - Update reva to v1.4.1-0.20210111080247-f2b63bfd6825: [#1194](https://github.com/owncloud/ocis/pull/1194)

# Changelog for [1.0.0] (2020-12-17)

The following sections list the changes for 1.0.0.

## Summary

* Bugfix - Enable scrolling in accounts list: [#909](https://github.com/owncloud/ocis/pull/909)
* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)
* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)
* Bugfix - Lower Bound was not working for the cs3 api index implementation: [#741](https://github.com/owncloud/ocis/pull/741)
* Bugfix - Accounts config sometimes being overwritten: [#808](https://github.com/owncloud/ocis/pull/808)
* Bugfix - Make settings service start without go coroutines: [#835](https://github.com/owncloud/ocis/pull/835)
* Bugfix - Fix button layout after phoenix update: [#625](https://github.com/owncloud/ocis/pull/625)
* Bugfix - Fix choose account dialogue: [#846](https://github.com/owncloud/ocis/pull/846)
* Bugfix - Fix id or username query handling: [#745](https://github.com/owncloud/ocis/pull/745)
* Bugfix - Fix konnectd build: [#809](https://github.com/owncloud/ocis/pull/809)
* Bugfix - Fix path of files shared with me in ocs api: [#204](https://github.com/owncloud/product/issues/204)
* Bugfix - Use micro default client: [#718](https://github.com/owncloud/ocis/pull/718)
* Bugfix - Allow consent-prompt with switch-account: [#788](https://github.com/owncloud/ocis/pull/788)
* Bugfix - Mint token with uid and gid: [#737](https://github.com/owncloud/ocis/pull/737)
* Bugfix - Serve index.html for directories: [#912](https://github.com/owncloud/ocis/pull/912)
* Bugfix - Don't create account if id/mail/username already taken: [#709](https://github.com/owncloud/ocis/pull/709)
* Bugfix - Fix director selection in proxy: [#521](https://github.com/owncloud/ocis/pull/521)
* Bugfix - Permission checks for settings write access: [#1092](https://github.com/owncloud/ocis/pull/1092)
* Bugfix - Fix minor ui bugs: [#1043](https://github.com/owncloud/ocis/issues/1043)
* Bugfix - Disable public link expiration by default: [#987](https://github.com/owncloud/ocis/issues/987)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)
* Change - Accounts UI shows message when no permissions: [#656](https://github.com/owncloud/ocis/pull/656)
* Change - Cache password validation: [#958](https://github.com/owncloud/ocis/pull/958)
* Change - Filesystem based index: [#709](https://github.com/owncloud/ocis/pull/709)
* Change - Rebuild index command for accounts: [#748](https://github.com/owncloud/ocis/pull/748)
* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)
* Change - CS3 can be used as accounts-backend: [#1020](https://github.com/owncloud/ocis/pull/1020)
* Change - Use bcrypt to hash the user passwords: [#510](https://github.com/owncloud/ocis/issues/510)
* Change - Replace the library which scales the images: [#910](https://github.com/owncloud/ocis/pull/910)
* Change - Choose disk or cs3 storage for accounts and groups: [#623](https://github.com/owncloud/ocis/pull/623)
* Change - Enable OpenID dynamic client registration: [#811](https://github.com/owncloud/ocis/issues/811)
* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)
* Change - Improve reva service descriptions: [#536](https://github.com/owncloud/ocis/pull/536)
* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)
* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)
* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)
* Change - Properly style konnectd consent page: [#754](https://github.com/owncloud/ocis/pull/754)
* Change - Make all paths configurable and default to a common temp dir: [#1080](https://github.com/owncloud/ocis/pull/1080)
* Change - Move the indexer package from ocis/accounts to ocis/ocis-pkg: [#794](https://github.com/owncloud/ocis/pull/794)
* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Move ocis default config to root level: [#842](https://github.com/owncloud/ocis/pull/842)
* Change - Remove username field in OCS: [#709](https://github.com/owncloud/ocis/pull/709)
* Change - Account management permissions for Admin role: [#124](https://github.com/owncloud/product/issues/124)
* Change - Update phoenix to v0.18.0: [#651](https://github.com/owncloud/ocis/pull/651)
* Change - Default apps in ownCloud Web: [#688](https://github.com/owncloud/ocis/pull/688)
* Change - Proxy allow insecure upstreams: [#1007](https://github.com/owncloud/ocis/pull/1007)
* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)
* Change - Theme welcome and choose account pages: [#887](https://github.com/owncloud/ocis/pull/887)
* Change - Bring oC theme: [#698](https://github.com/owncloud/ocis/pull/698)
* Change - Unify Configuration Parsing: [#675](https://github.com/owncloud/ocis/pull/675)
* Change - Update phoenix to v0.20.0: [#674](https://github.com/owncloud/ocis/pull/674)
* Change - Update phoenix to v0.21.0: [#728](https://github.com/owncloud/ocis/pull/728)
* Change - Update phoenix to v0.22.0: [#757](https://github.com/owncloud/ocis/pull/757)
* Change - Update phoenix to v0.23.0: [#785](https://github.com/owncloud/ocis/pull/785)
* Change - Update phoenix to v0.24.0: [#817](https://github.com/owncloud/ocis/pull/817)
* Change - Update phoenix to v0.25.0: [#868](https://github.com/owncloud/ocis/pull/868)
* Change - Update phoenix to v0.26.0: [#935](https://github.com/owncloud/ocis/pull/935)
* Change - Update phoenix to v0.27.0: [#943](https://github.com/owncloud/ocis/pull/943)
* Change - Update phoenix to v0.28.0: [#1027](https://github.com/owncloud/ocis/pull/1027)
* Change - Update phoenix to v0.29.0: [#1034](https://github.com/owncloud/ocis/pull/1034)
* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)
* Change - Update reva to v1.4.1-0.20201209113234-e791b5599a89: [#1089](https://github.com/owncloud/ocis/pull/1089)
* Change - Clarify storage driver env vars: [#729](https://github.com/owncloud/ocis/pull/729)
* Change - Update ownCloud Web to v1.0.0-beta3: [#1105](https://github.com/owncloud/ocis/pull/1105)
* Change - Update ownCloud Web to v1.0.0-beta4: [#1110](https://github.com/owncloud/ocis/pull/1110)
* Change - Settings and accounts appear in the user menu: [#656](https://github.com/owncloud/ocis/pull/656)
* Enhancement - Add tracing to the accounts service: [#1016](https://github.com/owncloud/ocis/issues/1016)
* Enhancement - Add the accounts service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add basic auth option: [#627](https://github.com/owncloud/ocis/pull/627)
* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)
* Enhancement - Add the glauth service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add k6: [#941](https://github.com/owncloud/ocis/pull/941)
* Enhancement - Add the konnectd service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocis-phoenix service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocis-pkg package: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocs service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the proxy service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the settings service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the storage service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the store service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the thumbnails service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add a command to list the versions of running instances: [#226](https://github.com/owncloud/product/issues/226)
* Enhancement - Add the webdav service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Better adopt Go-Micro: [#840](https://github.com/owncloud/ocis/pull/840)
* Enhancement - Add permission check when assigning and removing roles: [#879](https://github.com/owncloud/ocis/issues/879)
* Enhancement - Create OnlyOffice extension: [#857](https://github.com/owncloud/ocis/pull/857)
* Enhancement - Show basic-auth warning only once: [#886](https://github.com/owncloud/ocis/pull/886)
* Enhancement - Add glauth fallback backend: [#649](https://github.com/owncloud/ocis/pull/649)
* Enhancement - Tidy dependencies: [#845](https://github.com/owncloud/ocis/pull/845)
* Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)
* Enhancement - Add a version command to ocis: [#915](https://github.com/owncloud/ocis/pull/915)
* Enhancement - Create a proxy access-log: [#889](https://github.com/owncloud/ocis/pull/889)
* Enhancement - Cache userinfo in proxy: [#877](https://github.com/owncloud/ocis/pull/877)
* Enhancement - Update reva to v1.4.1-0.20201125144025-57da0c27434c: [#1320](https://github.com/cs3org/reva/pull/1320)
* Enhancement - Runtime Cleanup: [#1066](https://github.com/owncloud/ocis/pull/1066)
* Enhancement - Update OCIS Runtime: [#1108](https://github.com/owncloud/ocis/pull/1108)
* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)
* Enhancement - Update glauth to dev fd3ac7e4bbdc93578655d9a08d8e23f105aaa5b2: [#834](https://github.com/owncloud/ocis/pull/834)
* Enhancement - Update glauth to dev 4f029234b2308: [#786](https://github.com/owncloud/ocis/pull/786)
* Enhancement - Update konnectd to v0.33.8: [#744](https://github.com/owncloud/ocis/pull/744)
* Enhancement - Update reva to v1.4.1-0.20201123062044-b2c4af4e897d: [#823](https://github.com/owncloud/ocis/pull/823)
* Enhancement - Update reva to v1.4.1-0.20201130061320-ac85e68e0600: [#980](https://github.com/owncloud/ocis/pull/980)
* Enhancement - Update reva to cdb3d6688da5: [#748](https://github.com/owncloud/ocis/pull/748)
* Enhancement - Update reva to dd3a8c0f38: [#725](https://github.com/owncloud/ocis/pull/725)
* Enhancement - Update reva to v1.4.1-0.20201127111856-e6a6212c1b7b: [#971](https://github.com/owncloud/ocis/pull/971)
* Enhancement - Update reva to 063b3db9162b: [#1091](https://github.com/owncloud/ocis/pull/1091)
* Enhancement - Add www-authenticate based on user agent: [#1009](https://github.com/owncloud/ocis/pull/1009)

