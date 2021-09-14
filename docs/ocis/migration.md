---
title: "Migration"
date: 2021-03-16T16:17:00+01:00
weight: 41
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: migration.md
---

The migration happens in subsequent stages while the service is online. First all users need to migrate to the new architecture, then the global namespace needs to be introduced. Finally, the data on disk can be migrated user by user by switching the storage driver.

<div class="editpage">

{{< hint warning >}}
@jfd: It might be easier to introduce the spaces api in oc10 and then migrate to oCIS. We cannot migrate both at the same time, the architecture to oCIS (which will change fileids) and introduce a global namespace (which requires stable fileids to let clients handle moves without redownloading). Either we implement arbitrary mounting of shares in oCIS / reva or we make clients and oc10 spaces aware.
{{< /hint >}}

</div>

## Migration Stages

### Stage 0: pre migration
Is the pre-migration stage when having a functional ownCloud 10 instance.

<div class="editpage">

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage 1: introduce ownCloud Web
Install and introduce [ownCloud Web](https://github.com/owncloud/web/) and let users test it voluntarily to gain early feedback on the new UI.

#### Steps
Deploy web and enable switching to and from it.
For more details see: [ownCloud 10 with ownCloud Web](https://owncloud.dev/clients/web/deployments/oc10-app/)

<div class="editpage">

_TODO allow limiting the web ui switch to an 'early adopters' group_

</div>

#### Validation
Ensure switching back an forth between the classic ownCloud 10 web UI and ownCloud web works as at our https://demo.owncloud.com. 

#### Rollback
Should there be problems with ownCloud web at this point it can simply be removed from the menu and be undeployed. 

#### Notes
<div style="break-after: avoid"></div>
The ownCloud 10 demo instance uses OAuth to obtain a token for ownCloud web and currently always requires explicit consent. In oCIS the token is provided by the OpenID Connect Identity Provider, which may skip the consent step for trusted clients for a more seamless login experience. You may want to introduce OpenID Connect before enabling the new web UI.

<div class="editpage">

_TODO make oauth2 in oc10 trust the new web ui, based on `redirect_uri` and CSRF so no explicit consent is needed_

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage 2: introduce OpenID Connect

Basic auth requires us to properly store and manage user credentials. Something we would rather like to delegate to a tool specifically built for that task.
While SAML and Shibboleth are protocols that solve that problem, they are limited to web clients. Desktop and mobile clients were an afterthought and keep running into timeouts. For these reasons, we decided to move to [OpenID Connect as our primary authentication protocol](https://owncloud.com/news/openid-connect-oidc-app/). 

<div class="editpage">

_TODO @butonic add ADR for OpenID Connect_

</div>

#### User impact
When introducing OpenID Connect, the clients will detect the new authentication scheme when their current way of authenticating returns an error. Users will then have to
reauthorize at the OpenID Connecd IdP, which again, may be configured to skip the consent step for trusted clients.

#### Steps
1. There are multiple products that can be used as an OpenID Connect IdP. We test with [LibreGraph Connect](https://github.com/libregraph/lico), which is also [embedded in oCIS](https://github.com/owncloud/web/). Other alternatives include [Keycloak](https://www.keycloak.org/) or [Ping](https://www.pingidentity.com/). Please refer to the corresponding setup instructions for the product you intent to use.

<div class="editpage">

_TODO @butonic flesh out oCIS IDP documentation_

</div>

2. Add [Openid Connect (OIDC)](https://doc.owncloud.com/server/admin_manual/configuration/user/oidc/) support to ownCloud 10.

#### Validation
When OpenID Connect support is enabled verify that all clients can login:
- web classic
- ownCloud web
- desktop
- android
- iOS

#### Rollback
Should there be problems with OpenID Connect at this point you can disable the app. Users will have to reauthenticate in this case.

#### Notes
<div style="break-after: avoid"></div>
Legacy clients relying on Basic auth or app passwords need to be migrated to OpenId Connect to work with oCIS. For a transition period Basic auth in oCIS can be enabled with `PROXY_ENABLE_BASIC_AUTH=true`, but we strongly recommend adopting OpenID Connect for other tools as well.

While OpenID Connect providers will send an `iss` and `sub` claim that relying parties (services like oCIS or ownCloud 10) can use to identify users we recommend introducing a dedicated, globally unique, persistent, non-reassignable user identifier like a UUID for every user. This `ownclouduuid` shold be sent as an additional claim to save additional lookups on the server side. It will become the user id in oCIS, e.g. when searching for recipients the `ownclouduuid` will be used to persist permissions with the share manager. It has a different purpose than the ownCloud 10 username, which is used to login. Using UUIDs we can not only mitigate username collisions when merging multiple instances but also allow renaming usernames after the migration to oCIS has been completed.

<div class="editpage">

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage 3: introduce oCIS interally

Before letting oCIS handle end user requests we will first make it available in the internal network. By subsequently adding services we can add functionality and verify the services work as intended.

Start oCIS backend and make read only tests on existing data using the `owncloudsql` storage driver which will read (and write)
- blobs from the same data directory layout as in ownCloud 10
- metadata from the ownCloud 10 database: 
The oCIS share manager will read share information from the ownCloud database using an `owncloud` driver as well.

<div class="editpage">

_TODO @butonic add guide on how to configure `owncloudsql`_

_TODO we need a share manager that can read from the ownCloud 10 database as well as from whatever new backend will be used for a pure oCIS setup. Currently, that would be the json file. Or that is migrated after all users have switched to oCIS. -- jfd_

</div>

#### User impact
None, only administrators will be able to explore oCIS during this stage.

#### Steps and verifications

We are going to run and explore a series of services that will together handle the same requests as ownCloud 10. For initial exploration the oCIS binary is recommended. The services can later be deployed using a single oCIS runtime or in multiple containers.


##### Storage provider for file metadata
1. Deploy OCIS storage provider with the `owncloudsql` driver.
2. Set `read_only: true` in the storage provider config. <div class="editpage">_TODO @butonic add read only flag to storage drivers_</div>
3. Use cli tool to list files using the CS3 api

##### File ID alternatives
Multiple ownCloud instances can be merged into one oCIS instance. To prevent the numeric ids from colliding, the file ids will be prefixed with a new storage space id which is used by oCIS to route requests to the correct storage provider. See Stage 8 below.

<div class="editpage">

{{< hint warning >}}
**Alternative 1**
Add a routable prefix to fileids in oc10, and replicate the prefix in oCIS.
### Stage-3.1
Let oc10 render file ids with prefixes: `<instance name>$<numeric storageid>!<fileid>`. This will allow clients to handle moved files.

### Stage-3.2
Roll out new clients that understand the spaces API and know how to convert local sync pairs for legacy oc10 `/webdav` or `/dav/files/<username>` home folders into multiple sync pairs.
One pair for `/webdav/home` or `/dav/files/<username>/home` and another pair for every accepted share. The shares will be accessible at `/webdav/shares/` when the server side enables the spaces API.
Files can be identified using `<instance name>$<numeric storageid>!<fileid>` and moved to the correct sync pair.

### Stage-3.3
Enable spaces API in oc10:
- New clients will get a response from the spaces API and can set up new sync pairs.
- Legacy clients will still poll `/webdav` or `/dav/files/<username>` where they will see new subfolders instead of the users home. They will move down the users files into `/home` and shares into `/shares`. Custom sync pairs will no longer be available, causing the legacy client to leave local files in place. They can be picked up manually when installing a new client.

{{< /hint >}}

{{< hint warning >}}
**Alternative 2**
An additional `uuid` property used only to detect moves. A lookup by uuid is not necessary for this. The `/dav/meta` endpoint would still take the fileid. Clients  would use the `uuid` to detect moves and set up new sync pairs when migrating to a global namespace.
### Stage-3.1
Generate a `uuid` for every file as a file property. Clients can submit a `uuid` when creating files. The server will create a `uuid` if the client did not provide one.

### Stage-3.2
Roll out new clients that understand the spaces API and know how to convert local sync pairs for legacy oc10 `/webdav` or `/dav/files/<username>` home folders into multiple sync pairs.
One pair for `/webdav/home` or `/dav/files/<username>/home` and another pair for every accepted share. The shares will be accessible at `/webdav/shares/` when the server side enables the spaces API. Files can be identified using the `uuid`  and moved to the correct sync pair.

### Stage-4.1
When reading the files from oCIS return the same `uuid`. It can be migrated to an extended attribute or it can be read from oc10. If users change it the client will not be able to detect a move and maybe other weird stuff happens. *What if the uuid gets lost on the server side due to a partial restore?* 

{{< /hint >}}
</div>


<div style="break-after: page"></div>

##### graph API endpoint
1. Deploy graph api to list spaces
2. Use curl to list spaces using graph drives endpoint

##### owncloud flavoured WebDAV endpoint
1. Deploy ocdav
2. Use curl to send PROPFIND

##### data provider for up and download
1. Deploy dataprovider
2. Use curl to up and download files
3. Use tus to upload files

Deploy ...

##### share manager
Deploy share manager with ownCloud driver

##### reva gateway
1. Deploy gateway to authenticate requests? I guess we need that first... Or we need the to mint a token. Might be a good exercise.

##### automated deployment
Finally, deploy oCIS with a config to set up everything running in a single oCIS runtime or in multiple containers.

#### Rollback
You can stop the oCIS process at any time.

#### Notes
<div style="break-after: avoid"></div>
Multiple ownCloud instances can be merged into one oCIS instance. The file ids will be prefixed with a new storage space id which is used to route requests to the correct storage provider.

<div class="editpage">

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage 4: internal write access with oCIS
Test writing data with oCIS into the existing ownCloud 10 data directory using the `owncloudsql` storage driver.

#### User impact
Only administrators will be able to explore oCIS during this stage. End users should not be affected if the testing is limited to test users.

#### Steps
Set `read_only: false` in the storage provider config.

<div class="editpage">

_TODO @butonic add read only flag to storage drivers_

</div>

#### Verification
#### Rollback
Set `read_only: true` in the storage provider config.

<div class="editpage">

_TODO @butonic add read only flag to storage drivers_

</div>

#### Notes
<div style="break-after: avoid"></div>
With write access it becomes possible to manipulate existing files and shares.

<div class="editpage">

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage-5: introduce user aware proxy
In the previous stages oCIS was only accessible for administrators with access to the network. To expose only a single service to the internet, oCIS comes with a user aware proxy that can be used to route requests to the existing ownCloud 10 installation or oCIS, based on the authenticated user. The proxy uses OIDC to identify the logged in user and route them to the configured backend.

#### User impact
The IP address of the ownCloud host changes. There is no change for the file sync and share functionality when requests are handled by the oCIS codebase as it uses the same database and storage system as owncloud 10.

#### Steps and verifications 

##### Deploy oCIS proxy
1. Deploy the `ocis proxy` 
2. Verify the requests are routed based on the ownCloud 10 routing policy `oc10` by default

##### Test user based routing
1. Change the routing policy for a user or an early adopters group to `ocis` <div class="editpage">_TODO @butonic currently, the migration selector will use the `ocis` policy for users that have been added to the accounts service. IMO we need to evaluate a claim from the IdP._</div>
2. Verify the requests are routed based on the oCIS routing policy `oc10` for 'migrated' users.

At this point you are ready to rock & roll!

##### Let ownCloud domain point to proxy
1. Update the dns to use the oCIS proxy instead of the ownCloud application servers directly.
2. Let DNS propagate the change and monitor requests moving from the ownCloud application servers to the oCIS proxy.
3. Verify the DNS change has propagated sufficiently. All requests should now use the oCIS Proxy.

#### Rollback
Should there be a problem with the oCIS routes the user can be routed to ownCloud by changing his routing policy. In case of unfixable problems with the proxy the DNS needs to be updated to use the ownCloud 10 application servers directly. This could also be done in a load balancer.

#### Notes
<div style="break-after: avoid"></div>
The proxy is stateless, multiple instances can be deployed as needed.

<div class="editpage">

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage-6: parallel deployment
Running ownCloud 10 and oCIS in parallel is a crucial stage for the migration: it allows users access to group shares regardless of the system that is being used to to access the data. A user by user migration with multiple domains would technically break group shares when users vanish because they (and their data) are no longer available in the old system.

Depending on the amount of power users on an instance, the admin may want to allow users to voluntarily migrate to the oCIS backend. A monitoring system can be used to visualize the behavior for the two systems and gain trust in the overall stability and performance.

#### User impact
Since the underling data is still stored in the same systems, a similar or performance can be expected.
<div class="editpage">

See _TODO hmpf outdated didn't we want to run them nightly? ..._
_TODO @butonic update performance comparisons nightly_

</div>

#### Steps
There are several options to move users to the oCIS backend:
- Use a canary app to let users decide thamselves
- Use an early adoptors group with an opt in
- Force migrate users in batch or one by one at the administrators will

#### Verification
The same verification steps as for the internal testing stage apply. Just from the outside.

#### Rollback
Until now, the oCIS configuration mimics ownCloud 10 and uses the old data directory layout and the ownCloud 10 database. Users can seamlessly be switched from ownCloud 10 to oCIS and back again.
<div class="editpage">

_TODO @butonic we need a canary app that allows users to decide for themself which backend to use_

</div>

<div style="break-after: page"></div>

#### Notes
Running the two systems in parallel requires additional maintenance effort. Try to keep the duration of this stage short. Until now, we only added services and made the system more complex. oCIS aims to reduce the maintenance cost of an ownCloud instance. You will not get there if you keep both systems alive.

<div class="editpage">

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage-7: introduce spaces using ocis
To encourage users to switch you can promote the workspaces feature that is built into oCIS. The ownCloud 10 storage backend can be used for existing users. New users and group or project spaces can be provided by storage providers that better suit the underlying storage system.

#### Steps
First, the admin needs to 
- deploy a storage provider with the storage driver that best fits the underlying storage system and requirements. 
- register the storage in the storage registry with a new storage id (we recommend a uuid).

Then a user with the necessary create storage space role can create a storage space and assign Managers.

<div class="editpage">

_TODO @butonic a user with management permission needs to be presented with a list of storage spaces where he can see the amount of free space and decide on which storage provider the storage space should be created. For now a config option for the default storage provider for a specific type might be good enough._

</div>

#### Verification
The new storage space should show up in the `/graph/drives` endpoint for the managers and the creator of the space.

#### Notes
Depending on the requirements and acceptable tradeoffs, a database less deployment using the ocis or s3ng storage driver is possible. There is also a [cephfs driver](https://github.com/cs3org/reva/pull/1209) on the way, that directly works on the API level instead of POSIX.

### Stage-8: shut down ownCloud 10 
Disable ownCloud 10 in the proxy, all requests are now handled by oCIS, shut down oc10 web servers and redis (or keep for calendar & contacts only? rip out files from oCIS?)

#### User impact
All users are already sent to the oCIS backend. Shutting down ownCloud 10 will remove the old web UI, apps and functionality that is not yet present in ownCloud web. For example contacts and calendar.

<div class="editpage">

_TODO @butonic recommend alternatives_

</div>

#### Steps 
1. Shut down the apache servers that are running the ownCloud 10 PHP code.
2. DO NOT SHUT DOWN THE DATABASE, YET!

#### Verification
The ownCloud 10 classic web UI should no longer be reachable.

#### Rollback
Redeploy ownCloud 10.

#### Notes
<div style="break-after: avoid"></div>
The database needs to remain online until the storage layer and share metadata have been migrated as well. One thing at a time.

<div class="editpage">

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage 9: storage migration
To get rid of the database we will move the metadata from the old ownCloud 10 database into dedicated storage providers. This can happen in a user by user fashion. group drives can properly be migrated to group, project or workspaces in this stage.

#### User impact
Noticeable performance improvements because we effectively shard the storage logic and persistence layer.

#### Steps 
1. User by user storage migration from `owncloud` or `ownclouds3` driver to `ocis`/`s3ng`/`cephfs`... currently this means copying the metadata from one storage provider to another using the cs3 api.
2. Change the responsible storage provider for a storage space (e.g. a user home, a group or project space are a workspace) in the storage registry.

<div class="editpage">

_TODO @butonic implement `ownclouds3` based on `s3ng`_
_TODO @butonic implement tiered storage provider for seamless migration_
_TODO @butonic document how to manually do that until the storage registry can discover that on its own._

</div>

#### Verification
Start with a test user, then move to early adopters and finally migrate all users.

#### Rollback
To switch the storage provider again the same storage space migration can be performed again: copy medatata and blob data using the CS3 api, then change the responsible storage provider in the storage registry.

#### Notes
<div style="break-after: avoid"></div>
Multiple ownCloud instances can be merged into one oCIS instance. The file ids will be prefixed with a new storage space id which is used to route requests to the correct storage provider.

The storage space migration will become a seamless feature in the future that allows administrators to move users to storage systems with different capabilities, to implement premium features, deprovisioning strategies or archiving.

<div class="editpage">

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

### Stage-10: share metadata migration
Migrate share data to _yet to determine_ share manager backend and shut down ownCloud database.

The ownCloud 10 database still holds share information in the `oc_share` and `oc_share_external` tables. They are used to efficiently answer queries about who shared what with whom. In oCIS shares are persisted using a share manager and if desired these grants are also sent to the storage provider so it can set ACLs if possible. Only one system should be responsible for the shares, which in case of treating the storage as the primary source effectively turns the share manager into a cache.

#### User impact
Depending on chosen the share manager provider some sharing requests should be faster: listing incoming and outgoing shares is no longer bound to the ownCloud 10 database but to whatever technology is used by the share provider:
  - For non HA scenarios they can be served from memory, backed by a simple json file.
  - TODO: implement share manager with redis / nats / ... key value store backend: use the micro store interface please ...

#### Steps 
1. Start new share manager
2. Migrate metadata using the CS3 API (copy from old to new)
3. Shut down old share manager
4. Shut down ownCloud 10 database

<div class="editpage">

_TODO for HA implement share manager with redis / nats / ... key value store backend: use the micro store interface please ..._
_TODO for batch migration implement share data migration cli with progress that reads all shares via the cs3 api from one provider and writes them into another provider_
_TODO for seamless migration implement tiered/chained share provider that reads share data from the old provider and writes newc shares to the new one_
_TODO for storage provider as source of truth persist ALL share data in the storage provider. Currently, part is stored in the share manager, part is in the storage provider. We can keep both, but the the share manager should directly persist its metadata to the storage system used by the storage provider so metadata is kept in sync_

</div>

#### Verification
After copying all metadata start a dedicated gateway and change the configuration to use the new share manager. Route a test user, a test group and early adoptors to the new gateway. When no problems occur you can start the desired number of share managers and roll out the change to all gateways.

<div class="editpage">

_TODO let the gateway write updates to multiple share managers ... or rely on the tiered/chained share manager provider to persist to both providers_

</div>

#### Rollback
To switch the share manager to the database one revert routing users to the new share manager. If you already shut down the old share manager start it again. Use the tiered/chained share manager provider in reverse configuration (new share provider as read only, old as write) and migrate the shares again. You can alse restore a database backup if needed.

<div class="editpage">

### Stage-11
Profit! Well, on the one hand you do not need to maintain a clustered database setup and can rely on the storage system. On the other hand you are now in microservice wonderland and will have to relearn how to identify bottlenecks and scale oCIS accordingly. The good thing is that tools like jaeger and prometheus have evolved and will help you understand what is going on. But this is a different topic. See you on the other side!

#### FAQ
_Feel free to add your question as a PR to this document using the link at the top of this page!_ 

</div>

<div style="break-after: page"></div>

## Architectural differences

The fundamental difference between ownCloud 10 and oCIS is that the file metadata is moved from the database in the `oc_filecache` table (which is misnamed, as it actually is an index) to the storage provider who can place metadata as close to the underlying storage system as possible. In effect, the file metadata is sharded over multiple specialized services.


## Data that will be migrated

Currently, oCIS focuses on file sync and share use cases. 

### Blob data

In ownCloud 10 the files are laid out on disk in the *data directory* using the following layout:
```
data
├── einstein
│   ├── cache
│   ├── files
│   │   ├── Photos
│   │   │   └── Portugal.jpg
│   │   ├── Projects
│   │   │   └── Notes.md
│   │   └── ownCloud Manual.pdf
│   ├── files_external
│   ├── files_trashbin
│   │   ├── files
│   │   │   ├── Documents.d1564687985
│   │   │   ├── TODO.txt.d1565721976
│   │   │   └── welcome.txt.d1564775872
│   │   └── versions
│   │   │   ├── TODO.txt.v1564605543.d1565721976
│   │   │   └── TODO.txt.v1564775936.d1565721976
│   ├── files_versions
│   │   ├── Projects
│   │   │   ├── Notes.md.v1496912691
│   │   │   └── Notes.md.v1540305560
│   │   └── ownCloud Manual.pdf.v1396628249
│   ├── thumbnails
│   │   └── 123                  
│   │   │   ├── 2048-1536-max.png
│   │   │   └── 32-32.png                 // the file id, eg. of /Photos/Portugal.jpg
│   └── uploads
├── marie
│   ├── cache
│   ├── files
│   ├── files_external
│   ├── files_trashbin
│   ├── files_versions
│   └── thumbnails
│   …
├── moss
…
```

The *data directory* may also contain subfolders for ownCloud 10 applications like `avatars`, `gallery`, `files_external` and `cache`.

When an object storage is used as the primary storage all file blobs are stored by their file id and a prefix, eg.: `urn:oid:<fileid>`.

The three types of blobs we need to migrate are stored in 
- `files` for file blobs, the current file content,
- `files_trashbin` for trashed files (and their versions) and 
- `files_versions` for file blobs of older versions.

<div style="break-after: page"></div>

### Filecache table

In both cases the file metadata, including a full replication of the file tree, is stored in the `oc_filecache` table of an ownCloud 10 database. The primary key of a row is the file id. It is used to attach additional metadata like shares, tags, favorites or arbitrary file properties.

The `filecache` table itself has more metadata:

| Field              | Type          | Null | Key | Default | Extra          | Comment        | Migration         |
|--------------------|---------------|------|-----|---------|----------------|----------------|----------------|
| `fileid`           | bigint(20)    | NO   | PRI | NULL    | auto_increment |                | MUST become the oCIS `opaqueid` of a file reference. `ocis` driver stores it in extendet attributes and can use numbers as node ids on disk. for eos see note below table |
| `storage`          | int(11)       | NO   | MUL | 0       |                | *the filecache holds metadata for multiple storages* | corresponds to an oCIS *storage space* |
| `path`             | varchar(4000) | YES  |     | NULL    |                | *the path relative to the storages root* | MUST become the `path` relative to the storage root. `files` prefix needs to be trimmed. |
| `path_hash`        | varchar(32)   | NO   |     |         |                | *mysql once had problems indexing long paths, so we stored a hash for lookup by path. | - |
| `parent`           | bigint(20)    | NO   | MUL | 0       |                | *used to implement the hierarchy and listing children of a folder by id. redundant with `path`* | - |
| `name`             | varchar(250)  | YES  |     | NULL    |                | *basename of `path`*               | - |
| `mimetype`         | int(11)       | NO   |     | 0       |                | *joined with the `oc_mimetypes` table. only relevant for object storage deployments* | can be determined from blob / file extension |
| `mimepart`         | int(11)       | NO   |     | 0       |                | *"*               | can be determined from blob / file extension |
| `size`             | bigint(20)    | NO   |     | 0       |                | *same as blob size unless encryption is used*               | MAY become size, can be determined from blob |
| `mtime`            | bigint(20)    | NO   |     | 0       |                | *same as blob mtime*               | for files MAY become mtime (can be determined from blob as well), for directories MUST become tmtime |
| `encrypted`        | int(11)       | NO   |     | 0       |                | *encrypted flag*                | oCIS currently does not support encryption |
| `etag`             | varchar(40)   | YES  |     | NULL    |                | *used to propagate changes in a tree* | MUST be migrated (or calculated in the same way) to prevent clients from syncing unnecessarily |
| `unencrypted_size` | bigint(20)    | NO   |     | 0       |                | *same as blob size* | oCIS currently does not support encryption |
| `storage_mtime`    | bigint(20)    | NO   |     | 0       |                | *used to detect external storage changes* | oCIS delegates that to the storage providers and drivers |
| `permissions`      | int(11)       | YES  |     | 0       |                | *used as the basis for permissions. synced from disk when running a file scan. * | oCIS delegates that to the storage providers and drivers |
| `checksum`         | varchar(255)  | YES  |     | NULL    |                | *same as blob checksum* | SHOULD become the checksum in the storage provider. eos calculates it itself, `ocis` driver stores it in extended attributes |


> Note: for EOS a hot migration only works seamlessly if file ids in oc10 are already read from eos. otherwise either a mapping from the oc10 filecache file id to the new eos file id has to be created under the assumption that these id sets do not intersect or files and corresponding shares need to be exported and imported offline to generate a new set of ids. While this will preserve public links, user, group and even federated shares, old internal links may still point to different files because they contain the oc10 fileid 

<div style="break-after: page"></div>

### share table

used to store
- Public links
- Private shares with users and groups
- Federated shares *partly*
- Guest shares

| Field         | Type         | Null | Key | Default | Extra          | Comment | [CS3 API](https://cs3org.github.io/cs3apis/) |
|---------------|--------------|------|-----|---------|----------------|---------|-|
| `id`            | int(11)      | NO   | PRI | NULL    | auto_increment | | `ShareId.opaqueid` string |
| `share_type`    | smallint(6)  | NO   |     | 0       |                | *in CS3 every type is handled by a dedicated API. See below the table* | does NOT map to [`Share.ShareType`](https://cs3org.github.io/cs3apis/#cs3.sharing.ocm.v1beta1.Share.ShareType) *TODO clarify* |
| `share_with`    | varchar(255) | YES  | MUL | NULL    |                | | `Share.grantee` [`Grantee`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Grantee) |
| `uid_owner`     | varchar(64)  | NO   |     |         |                | | `ShareId.owner` [`UserID`](https://cs3org.github.io/cs3apis/#cs3.identity.user.v1beta1.UserId) |
| `parent`        | int(11)      | YES  |     | NULL    |                | | - |
| `item_type`     | varchar(64)  | NO   | MUL |         |                | | `Share.resource_id` [`ResourceId`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId) |
| `item_source`   | varchar(255) | YES  | MUL | NULL    |                | | `Share.resource_id` [`ResourceId`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId) |
| `item_target`   | varchar(255) | YES  |     | NULL    |                | | `Share.resource_id` [`ResourceId`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId) |
| `file_source`   | bigint(20)   | YES  | MUL | NULL    |  | *cannot store uuid style file ids from oCIS. when all users have migrated to oCIS the share manager needs to be updated / migrated to a version that does.* | `Share.resource_id` [`ResourceId`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId) |
| `file_target`   | varchar(512) | YES  |     | NULL    |                | | `Share.resource_id` [`ResourceId`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId) |
| `permissions`   | smallint(6)  | NO   |     | 0       |                | | `Share.Permissions` [`SharePermissions`](https://cs3org.github.io/cs3apis/#cs3.sharing.ocm.v1beta1.SharePermissions) |
| `stime`         | bigint(20)   | NO   |     | 0       |                | | `Share.ctime`, `Share.mtime` |
| `accepted`      | smallint(6)  | NO   |     | 0       |                | | `ReceivedShare.ShareState` [`ShareState`](https://cs3org.github.io/cs3apis/#cs3.sharing.collaboration.v1beta1.ShareState) |
| `expiration`    | datetime     | YES  |     | NULL    |                | *only used for the Link API and storage provider api, currently cannot be added using the Collaboration or OCM API* | [`Grant`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Grant)  |
| `token`         | varchar(32)  | YES  | MUL | NULL    |                | | [`PublicShare.token`](https://cs3org.github.io/cs3apis/#cs3.sharing.link.v1beta1.PublicShare) |
| `mail_send`     | smallint(6)  | NO   |     | 0       |                | | - |
| `uid_initiator` | varchar(64)  | YES  |     | NULL    |                | | `ShareId.creator` [`UserID`](https://cs3org.github.io/cs3apis/#cs3.identity.user.v1beta1.UserId) |
| `share_name`    | varchar(64)  | YES  |     | NULL    |                | *only exists for public shares* | [`PublicShare.display_name`](https://cs3org.github.io/cs3apis/#cs3.sharing.link.v1beta1.PublicShare)  |
| `attributes`    | longtext     | YES  |     | NULL    |                | *additional share attributes* | *could be implemented using opaque data, but should be added to the CS3 api* |

In the CS3 API
1. public links are handled by the PublicShareProvider using the [Link API](https://cs3org.github.io/cs3apis/#cs3.sharing.link.v1beta1.LinkAPI)
2. internal shares are handled by the UserShareProvider using the [Collaboration API](https://cs3org.github.io/cs3apis/#cs3.sharing.collaboration.v1beta1.CollaborationAPI). This covers user and group shares.
3. federated shares are handled by the OcmShareProvider using the [OCM Share Provider AP](https://cs3org.github.io/cs3apis/#cs3.sharing.ocm.v1beta1.OcmAPI) aka. Open Cloud Mesh.


<div style="break-after: page"></div>

### share_external

Used to store additional metadata for federated shares.

| Field           | Type          | Null | Key | Default | Extra          | Comment |
|-----------------|---------------|------|-----|---------|----------------|---------|
| `id`              | bigint(20)    | NO   | PRI | NULL    | auto_increment | |
| `remote`          | varchar(512)  | NO   |     | NULL    |                | Url of the remote owncloud instance |
| `share_token`     | varchar(64)   | NO   |     | NULL    |                | Public share token |
| `password`        | varchar(64)   | YES  |     | NULL    |                | Optional password for the public share |
| `name`            | varchar(64)   | NO   |     | NULL    |                | Original name on the remote server |
| `owner`           | varchar(64)   | NO   |     | NULL    |                | User that owns the public share on the remote server |
| `user`            | varchar(64)   | NO   | MUL | NULL    |                | Local user which added the external share |
| `mountpoint`      | varchar(4000) | NO   |     | NULL    |                | Full path where the share is mounted |
| `mountpoint_hash` | varchar(32)   | NO   |     | NULL    |                | md5 hash of the mountpoint |
| `remote_id`       | varchar(255)  | NO   |     | -1      |                | |
| `accepted`        | int(11)       | NO   |     | 0       |                | |

<div class="editpage">

_TODO document how the reva OCM service currently persists the data_

</div>

<div style="break-after: page"></div>

### trusted_servers

used to determine if federated shares can automatically be accepted

| Field         | Type         | Null | Key | Default | Extra          | Comment |
|---------------|--------------|------|-----|---------|----------------|---------|
| `id`            | int(11)      | NO   | PRI | NULL    | auto_increment | |
| `url`           | varchar(512) | NO   |     | NULL    |                | Url of trusted server |
| `url_hash`      | varchar(255) | NO   | UNI |         |                | sha1 hash of the url without the protocol |
| `token`         | varchar(128) | YES  |     | NULL    |                | token used to exchange the shared secret |
| `shared_secret` | varchar(256) | YES  |     | NULL    |                | shared secret used to authenticate |
| `status`        | int(11)      | NO   |     | 2       |                | current status of the connection |
| `sync_token`    | varchar(512) | YES  |     | NULL    |                | cardDav sync token |

<div class="editpage">

_TODO clarify how OCM handles this and where we store / configure this. It seems related to trusted IdPs_

</div>

<div style="break-after: page"></div>

### user data

Users are migrated in two steps:
1. They should all be authenticated using OpenID Connect, which already moves them to a common identity management system.
2. To search share recipients, both, ownCloud 10 and oCIS need access to the same user directory using eg. LDAP.

<div class="editpage">

_TODO add state to CS3 API, so we can 'disable' users_
_TODO how do we map (sub) admins? -> map to roles & permissions_

</div>

accounts:

| Field         | Type                | Null | Key | Default | Extra          | Comment |
|---------------|---------------------|------|-----|---------|----------------|---------|
| `id`            | bigint(20) unsigned | NO   | PRI | NULL    | auto_increment | |
| `email`         | varchar(255)        | YES  | MUL | NULL    |                | |
| `user_id`       | varchar(255)        | NO   | UNI | NULL    |                | |
| `lower_user_id` | varchar(255)        | NO   | UNI | NULL    |                | |
| `display_name`  | varchar(255)        | YES  | MUL | NULL    |                | |
| `quota`         | varchar(32)         | YES  |     | NULL    |                | |
| `last_login`    | int(11)             | NO   |     | 0       |                | |
| `backend`       | varchar(64)         | NO   |     | NULL    |                | |
| `home`          | varchar(1024)       | NO   |     | NULL    |                | |
| `state`         | smallint(6)         | NO   |     | 0       |                | |

users:

| Field       | Type         | Null | Key | Default | Extra | Comment |
|-------------|--------------|------|-----|---------|-------|---------|
| `uid`         | varchar(64)  | NO   | PRI |         |       |
| `password`    | varchar(255) | NO   |     |         |       |
| `displayname` | varchar(64)  | YES  |     | NULL    |       |

groups:

The groups table really only contains the group name.

| Field | Type        | Null | Key | Default | Extra |
|-------|-------------|------|-----|---------|-------|
| `gid`   | varchar(64) | NO   | PRI |         |       |

<div style="break-after: page"></div>

### LDAP

<div class="editpage">

_TODO clarify if metadata from ldap & user_shibboleth needs to be migrated_

</div>

The `dn` -> *owncloud internal username* mapping that currently lives in the `oc_ldap_user_mapping` table needs to move into a dedicated ownclouduuid attribute in the LDAP server. The idp should send it as a claim so the proxy does not have to look up the user using LDAP again. The username cannot be changed in ownCloud 10 and the oCIS provisioning API will not allow changing it as well. When we introduce the graph api we may allow changing usernames when all clients have moved to that api.

The problem is that the username in owncloud 10 and in oCIS also need to be the same, which might not be the case when the ldap mapping used a different column. In that case we should add another owncloudusername attribute to the ldap server.


<div class="editpage">

### activities

*dedicated service, not yet implemented, requires decisions about an event system -- jfd*

| Field         | Type          | Null | Key | Default | Extra          | Comment |
|---------------|---------------|------|-----|---------|----------------|---------|
| `activity_id`   | bigint(20)    | NO   | PRI | NULL    | auto_increment |
| `timestamp`     | int(11)       | NO   | MUL | 0       |                |
| `priority`      | int(11)       | NO   |     | 0       |                |
| `type`          | varchar(255)  | YES  |     | NULL    |                |
| `user`          | varchar(64)   | YES  |     | NULL    |                |
| `affecteduser`  | varchar(64)   | NO   | MUL | NULL    |                |
| `app`           | varchar(255)  | NO   |     | NULL    |                |
| `subject`       | varchar(255)  | NO   |     | NULL    |                |
| `subjectparams` | longtext      | NO   |     | NULL    |                |
| `message`       | varchar(255)  | YES  |     | NULL    |                |
| `messageparams` | longtext      | YES  |     | NULL    |                |
| `file`          | varchar(4000) | YES  |     | NULL    |                |
| `link`          | varchar(4000) | YES  |     | NULL    |                |
| `object_type`   | varchar(255)  | YES  | MUL | NULL    |                |
| `object_id`     | bigint(20)    | NO   |     | 0       |                |

## Links

The [data_exporter](https://github.com/owncloud/data_exporter) has logic that allows exporting and importing users, including shares. The [model classes](https://github.com/owncloud/data_exporter/tree/master/lib/Model) contain the exact mapping.

</div>
