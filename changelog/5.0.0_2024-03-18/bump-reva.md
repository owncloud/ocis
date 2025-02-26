Enhancement: Update reva to 2.19.2

We update reva to the version 2.19.2

*   Bugfix [cs3org/reva#4557](https://github.com/cs3org/reva/pull/4557): Fix ceph build
*   Bugfix [cs3org/reva#4570](https://github.com/cs3org/reva/pull/4570): Fix sharing invite on virtual drive
*   Bugfix [cs3org/reva#4559](https://github.com/cs3org/reva/pull/4559): Fix graph drive invite
*   Bugfix [cs3org/reva#4518](https://github.com/cs3org/reva/pull/4518): Fix an error when lock/unlock a file
*   Bugfix [cs3org/reva#4566](https://github.com/cs3org/reva/pull/4566): Fix public link previews
*   Bugfix [cs3org/reva#4561](https://github.com/cs3org/reva/pull/4561): Fix Stat() by Path on re-created resource
*   Enhancement [cs3org/reva#4556](https://github.com/cs3org/reva/pull/4556): Allow tracing requests by giving util functions a context
*   Enhancement [cs3org/reva#4545](https://github.com/cs3org/reva/pull/4545): Extend service account permissions
*   Enhancement [cs3org/reva#4564](https://github.com/cs3org/reva/pull/4564): Send file locked/unlocked events

We update reva to the version 2.19.1

*   Bugfix [cs3org/reva#4534](https://github.com/cs3org/reva/pull/4534): Fix remove/update share permissions
*   Bugfix [cs3org/reva#4539](https://github.com/cs3org/reva/pull/4539): Fix a typo

We update reva to the version 2.19.0

*   Bugfix [cs3org/reva#4464](https://github.com/cs3org/reva/pull/4464): Don't check lock grants
*   Bugfix [cs3org/reva#4516](https://github.com/cs3org/reva/pull/4516): The sharemanager can now reject grants with resharing permissions
*   Bugfix [cs3org/reva#4512](https://github.com/cs3org/reva/pull/4512): Bump dependencies
*   Bugfix [cs3org/reva#4481](https://github.com/cs3org/reva/pull/4481): Distinguish failure and node metadata reversal
*   Bugfix [cs3org/reva#4456](https://github.com/cs3org/reva/pull/4456): Do not lose revisions when restoring the first revision
*   Bugfix [cs3org/reva#4472](https://github.com/cs3org/reva/pull/4472): Fix concurrent access to a map
*   Bugfix [cs3org/reva#4457](https://github.com/cs3org/reva/pull/4457): Fix concurrent map access in sharecache
*   Bugfix [cs3org/reva#4498](https://github.com/cs3org/reva/pull/4498): Fix Content-Disposition header in dav
*   Bugfix [cs3org/reva#4461](https://github.com/cs3org/reva/pull/4461): CORS handling for WebDAV requests fixed
*   Bugfix [cs3org/reva#4462](https://github.com/cs3org/reva/pull/4462): Prevent setting container specific permissions on files
*   Bugfix [cs3org/reva#4479](https://github.com/cs3org/reva/pull/4479): Fix creating documents in the app provider
*   Bugfix [cs3org/reva#4474](https://github.com/cs3org/reva/pull/4474): Make /dav/meta consistent
*   Bugfix [cs3org/reva#4446](https://github.com/cs3org/reva/pull/4446): Disallow to delete a file during the processing
*   Bugfix [cs3org/reva#4517](https://github.com/cs3org/reva/pull/4517): Fix duplicated items in the sharejail root
*   Bugfix [cs3org/reva#4473](https://github.com/cs3org/reva/pull/4473): Decomposedfs now correctly lists sessions
*   Bugfix [cs3org/reva#4528](https://github.com/cs3org/reva/pull/4528): Respect IfNotExist option when uploading in cs3 metadata storage
*   Bugfix [cs3org/reva#4503](https://github.com/cs3org/reva/pull/4503): Fix an error when move
*   Bugfix [cs3org/reva#4466](https://github.com/cs3org/reva/pull/4466): Fix natsjskv store
*   Bugfix [cs3org/reva#4533](https://github.com/cs3org/reva/pull/4533): Fix recursive trashcan purge
*   Bugfix [cs3org/reva#4492](https://github.com/cs3org/reva/pull/4492): Fix the resource name
*   Bugfix [cs3org/reva#4463](https://github.com/cs3org/reva/pull/4463): Fix the resource name
*   Bugfix [cs3org/reva#4448](https://github.com/cs3org/reva/pull/4448): Fix truncating existing files
*   Bugfix [cs3org/reva#4434](https://github.com/cs3org/reva/pull/4434): Fix the upload postprocessing
*   Bugfix [cs3org/reva#4469](https://github.com/cs3org/reva/pull/4469): Handle interrupted uploads
*   Bugfix [cs3org/reva#4532](https://github.com/cs3org/reva/pull/4532): Jsoncs3 cache fixes
*   Bugfix [cs3org/reva#4449](https://github.com/cs3org/reva/pull/4449): Keep failed processing status
*   Bugfix [cs3org/reva#4529](https://github.com/cs3org/reva/pull/4529): We aligned some OCS return codes with oc10
*   Bugfix [cs3org/reva#4507](https://github.com/cs3org/reva/pull/4507): Make tusd CORS headers configurable
*   Bugfix [cs3org/reva#4452](https://github.com/cs3org/reva/pull/4452): More efficient share jail
*   Bugfix [cs3org/reva#4476](https://github.com/cs3org/reva/pull/4476): No need to unmark postprocessing when it was not started
*   Bugfix [cs3org/reva#4454](https://github.com/cs3org/reva/pull/4454): Skip unnecessary share retrieval
*   Bugfix [cs3org/reva#4527](https://github.com/cs3org/reva/pull/4527): Unify datagateway method handling
*   Bugfix [cs3org/reva#4530](https://github.com/cs3org/reva/pull/4530): Drop unnecessary grant exists check
*   Bugfix [cs3org/reva#4475](https://github.com/cs3org/reva/pull/4475): Upload session specific processing flag
*   Enhancement [cs3org/reva#4501](https://github.com/cs3org/reva/pull/4501): Allow sending multiple user ids in one sse event
*   Enhancement [cs3org/reva#4485](https://github.com/cs3org/reva/pull/4485): Modify the concurrency default
*   Enhancement [cs3org/reva#4526](https://github.com/cs3org/reva/pull/4526): Configurable s3 put options
*   Enhancement [cs3org/reva#4453](https://github.com/cs3org/reva/pull/4453): Disable the password policy
*   Enhancement [cs3org/reva#4477](https://github.com/cs3org/reva/pull/4477): Extend ResumePostprocessing event
*   Enhancement [cs3org/reva#4491](https://github.com/cs3org/reva/pull/4491): Add filename incrementor for secret filedrops
*   Enhancement [cs3org/reva#4490](https://github.com/cs3org/reva/pull/4490): Lazy initialize public share manager
*   Enhancement [cs3org/reva#4494](https://github.com/cs3org/reva/pull/4494): Start implementation of a plain posix storage driver
*   Enhancement [cs3org/reva#4502](https://github.com/cs3org/reva/pull/4502): Add spaceindex.AddAll()

## Changelog for reva 2.18.0 (2023-12-22)

The following sections list the changes in reva 2.18.0 relevant to
reva users. The changes are ordered by importance.

*   Bugfix [cs3org/reva#4424](https://github.com/cs3org/reva/pull/4424): Fixed panic in receivedsharecache pkg
*   Bugfix [cs3org/reva#4425](https://github.com/cs3org/reva/pull/4425): Fix overwriting files with empty files
*   Bugfix [cs3org/reva#4432](https://github.com/cs3org/reva/pull/4432): Fix /dav/meta endpoint for shares
*   Bugfix [cs3org/reva#4422](https://github.com/cs3org/reva/pull/4422): Fix disconnected traces
*   Bugfix [cs3org/reva#4429](https://github.com/cs3org/reva/pull/4429): Internal link creation
*   Bugfix [cs3org/reva#4407](https://github.com/cs3org/reva/pull/4407): Make ocdav return correct oc:spaceid
*   Bugfix [cs3org/reva#4410](https://github.com/cs3org/reva/pull/4410): Improve OCM support
*   Bugfix [cs3org/reva#4402](https://github.com/cs3org/reva/pull/4402): Refactor upload session
*   Enhancement [cs3org/reva#4421](https://github.com/cs3org/reva/pull/4421): Check permissions before adding, deleting or updating shares
*   Enhancement [cs3org/reva#4403](https://github.com/cs3org/reva/pull/4403): Add validation to update public share
*   Enhancement [cs3org/reva#4409](https://github.com/cs3org/reva/pull/4409): Disable the password policy
*   Enhancement [cs3org/reva#4412](https://github.com/cs3org/reva/pull/4412): Allow authentication for nats connections
*   Enhancement [cs3org/reva#4411](https://github.com/cs3org/reva/pull/4411): Add option to configure streams non durable
*   Enhancement [cs3org/reva#4406](https://github.com/cs3org/reva/pull/4406): Rework cache configuration
*   Enhancement [cs3org/reva#4414](https://github.com/cs3org/reva/pull/4414): Track more upload session metrics

## Changelog for reva 2.17.0 (2023-12-12)

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

https://github.com/owncloud/ocis/pull/8638
https://github.com/owncloud/ocis/pull/8519
https://github.com/owncloud/ocis/pull/8502
https://github.com/owncloud/ocis/pull/8340
https://github.com/owncloud/ocis/pull/8381
https://github.com/owncloud/ocis/pull/8287
https://github.com/owncloud/ocis/pull/8278
https://github.com/owncloud/ocis/pull/8264
https://github.com/owncloud/ocis/pull/8100
https://github.com/owncloud/ocis/pull/8100
https://github.com/owncloud/ocis/pull/8038
https://github.com/owncloud/ocis/pull/8056
https://github.com/owncloud/ocis/pull/7949
https://github.com/owncloud/ocis/pull/7793
https://github.com/owncloud/ocis/pull/7978
https://github.com/owncloud/ocis/pull/7979
https://github.com/owncloud/ocis/pull/7963
https://github.com/owncloud/ocis/pull/7986
https://github.com/owncloud/ocis/pull/7721
https://github.com/owncloud/ocis/pull/7727
https://github.com/owncloud/ocis/pull/7752
