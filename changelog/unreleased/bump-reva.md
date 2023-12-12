Enhancement: Update reva to v2.17.0

Changelog for reva 2.17.0 (2023-12-12)
=======================================

The following sections list the changes in reva 2.17.0 relevant to
reva users. The changes are ordered by importance.

*   Bugfix [cs3org/reva#4278](https://github.com/cs3org/reva/pull/4278): Disable DEPTH infinity in PROPFIND
*   Bugfix [cs3org/reva#4318](https://github.com/cs3org/reva/pull/4318): Do not allow moves between shares
*   Bugfix [cs3org/reva#4290](https://github.com/cs3org/reva/pull/4290): Prevent panic when trying to move a non-existent file
*   Bugfix [cs3org/reva#4241](https://github.com/cs3org/reva/pull/4241): Allow an empty credentials chain in the auth middleware
*   Bugfix [cs3org/reva#4216](https://github.com/cs3org/reva/pull/4216): Fix an error message
*   Bugfix [cs3org/reva#4324](https://github.com/cs3org/reva/pull/4324): Fix capabilities decoding
*   Bugfix [cs3org/reva#4267](https://github.com/cs3org/reva/pull/4267): Fix concurrency issue
*   Bugfix [cs3org/reva#4362](https://github.com/cs3org/reva/pull/4362): Fix concurrent lookup
*   Bugfix [cs3org/reva#4336](https://github.com/cs3org/reva/pull/4336): Fix definition of "file-editor" role
*   Bugfix [cs3org/reva#4302](https://github.com/cs3org/reva/pull/4302): Fix checking of filename length
*   Bugfix [cs3org/reva#4366](https://github.com/cs3org/reva/pull/4366): Fix CS3 status code when looking up non existing share
*   Bugfix [cs3org/reva#4299](https://github.com/cs3org/reva/pull/4299): Fix HTTP verb of the generate-invite endpoint
*   Bugfix [cs3org/reva#4249](https://github.com/cs3org/reva/pull/4249): GetUserByClaim not working with MSAD for claim "userid"
*   Bugfix [cs3org/reva#4217](https://github.com/cs3org/reva/pull/4217): Fix missing case for "hide" in UpdateShares
*   Bugfix [cs3org/reva#4140](https://github.com/cs3org/reva/pull/4140): Fix missing etag in shares jail
*   Bugfix [cs3org/reva#4229](https://github.com/cs3org/reva/pull/4229): Fix destroying the Personal and Project spaces data
*   Bugfix [cs3org/reva#4193](https://github.com/cs3org/reva/pull/4193): Fix overwrite a file with an empty file
*   Bugfix [cs3org/reva#4365](https://github.com/cs3org/reva/pull/4365): Fix create public share
*   Bugfix [cs3org/reva#4380](https://github.com/cs3org/reva/pull/4380): Fix the public link update
*   Bugfix [cs3org/reva#4250](https://github.com/cs3org/reva/pull/4250): Fix race condition
*   Bugfix [cs3org/reva#4345](https://github.com/cs3org/reva/pull/4345): Fix conversion of custom ocs permissions to roles
*   Bugfix [cs3org/reva#4134](https://github.com/cs3org/reva/pull/4134): Fix share jail
*   Bugfix [cs3org/reva#4335](https://github.com/cs3org/reva/pull/4335): Fix public shares cleanup config
*   Bugfix [cs3org/reva#4338](https://github.com/cs3org/reva/pull/4338): Fix unlock via space API
*   Bugfix [cs3org/reva#4341](https://github.com/cs3org/reva/pull/4341): Fix spaceID in meta endpoint response
*   Bugfix [cs3org/reva#4351](https://github.com/cs3org/reva/pull/4351): Fix 500 when open public link
*   Bugfix [cs3org/reva#4352](https://github.com/cs3org/reva/pull/4352): Fix the tgz mime type
*   Bugfix [cs3org/reva#4388](https://github.com/cs3org/reva/pull/4388): Allow UpdateUserShare() to update just the expiration date
*   Bugfix [cs3org/reva#4214](https://github.com/cs3org/reva/pull/4214): Always pass adjusted default nats options
*   Bugfix [cs3org/reva#4291](https://github.com/cs3org/reva/pull/4291): Release lock when expired
*   Bugfix [cs3org/reva#4386](https://github.com/cs3org/reva/pull/4386): Remove dead enable_home config
*   Bugfix [cs3org/reva#4292](https://github.com/cs3org/reva/pull/4292): Return 403 when user is not permitted to lock
*   Enhancement [cs3org/reva#4389](https://github.com/cs3org/reva/pull/4389): Add audio and location props
*   Enhancement [cs3org/reva#4337](https://github.com/cs3org/reva/pull/4337): Check permissions before creating shares
*   Enhancement [cs3org/reva#4326](https://github.com/cs3org/reva/pull/4326): Add search mediatype filter
*   Enhancement [cs3org/reva#4367](https://github.com/cs3org/reva/pull/4367): Add GGS mime type
*   Enhancement [cs3org/reva#4194](https://github.com/cs3org/reva/pull/4194): Add hide flag to shares
*   Enhancement [cs3org/reva#4358](https://github.com/cs3org/reva/pull/4358): Add default permissions capability for links
*   Enhancement [cs3org/reva#4133](https://github.com/cs3org/reva/pull/4133): Add more metadata to locks
*   Enhancement [cs3org/reva#4353](https://github.com/cs3org/reva/pull/4353): Add support for .docxf files
*   Enhancement [cs3org/reva#4363](https://github.com/cs3org/reva/pull/4363): Add nats-js-kv store
*   Enhancement [cs3org/reva#4197](https://github.com/cs3org/reva/pull/4197): Add the Banned-Passwords List
*   Enhancement [cs3org/reva#4190](https://github.com/cs3org/reva/pull/4190): Add the password policies
*   Enhancement [cs3org/reva#4384](https://github.com/cs3org/reva/pull/4384): Add a retry postprocessing outcome and event
*   Enhancement [cs3org/reva#4271](https://github.com/cs3org/reva/pull/4271): Add search capability
*   Enhancement [cs3org/reva#4119](https://github.com/cs3org/reva/pull/4119): Add sse event
*   Enhancement [cs3org/reva#4392](https://github.com/cs3org/reva/pull/4392): Add additional permissions to service accounts
*   Enhancement [cs3org/reva#4344](https://github.com/cs3org/reva/pull/4344): Add url extension to mime type list
*   Enhancement [cs3org/reva#4372](https://github.com/cs3org/reva/pull/4372): Add validation to the public share provider
*   Enhancement [cs3org/reva#4244](https://github.com/cs3org/reva/pull/4244): Allow listing reveived shares by service accounts
*   Enhancement [cs3org/reva#4129](https://github.com/cs3org/reva/pull/4129): Auto-Accept Shares through ServiceAccounts
*   Enhancement [cs3org/reva#4374](https://github.com/cs3org/reva/pull/4374): Handle trashbin file listings concurrently
*   Enhancement [cs3org/reva#4325](https://github.com/cs3org/reva/pull/4325): Enforce Permissions
*   Enhancement [cs3org/reva#4368](https://github.com/cs3org/reva/pull/4368): Extract log initialization
*   Enhancement [cs3org/reva#4375](https://github.com/cs3org/reva/pull/4375): Introduce UploadSessionLister interface
*   Enhancement [cs3org/reva#4268](https://github.com/cs3org/reva/pull/4268): Implement sharing roles
*   Enhancement [cs3org/reva#4160](https://github.com/cs3org/reva/pull/4160): Improve utils pkg
*   Enhancement [cs3org/reva#4335](https://github.com/cs3org/reva/pull/4335): Add sufficient permissions check function
*   Enhancement [cs3org/reva#4281](https://github.com/cs3org/reva/pull/4281): Port OCM changes from master
*   Enhancement [cs3org/reva#4270](https://github.com/cs3org/reva/pull/4270): Opt out of public link password enforcement
*   Enhancement [cs3org/reva#4181](https://github.com/cs3org/reva/pull/4181): Change the variable names for the password policy
*   Enhancement [cs3org/reva#4256](https://github.com/cs3org/reva/pull/4256): Rename hidden share variable name
*   Enhancement [cs3org/reva#3926](https://github.com/cs3org/reva/pull/3926): Service Accounts
*   Enhancement [cs3org/reva#4359](https://github.com/cs3org/reva/pull/4359): Update go-ldap to v3.4.6
*   Enhancement [cs3org/reva#4170](https://github.com/cs3org/reva/pull/4170): Update password policies
*   Enhancement [cs3org/reva#4232](https://github.com/cs3org/reva/pull/4232): Improve error handling in utils package

https://github.com/owncloud/ocis/pull/7949
