---
title: "Migration"
date: 2021-03-16T16:17:00+01:00
weight: 41
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: migration.md
---


## Migration

The migration happens in subsequent stages while the service is online. First all users need to migrate to the new architecture, then the data on disk can be migrated user by user by switching the storage driver.

### User Stories
As an admin I need to avoid downtime.
As an admin I want to migrate certain groups of users before others.
As a user, I need a seamless migration and not lose data by any chance.

### Stage-0
Is the pre-migration stage having a functional ownCloud 10 instance.

### Stage-1
Introduce OpenID Connect to server and Clients

### Stage-2
Install and introduce ownCloud Web and let users test it voluntarily.


{{< hint warning >}}
**Alternative 1**
Add a routable prefix to fileids in oc10, and replicate the prefix in ocis.
### Stage-2.1
Let oc10 render file ids with prefixes: `<instance name>$<numeric storageid>!<fileid>`. This will allow clients to handle moved files.

### Stage-2.2
Roll out new clients that understand the spaces API and know how to convert local sync pairs for legacy oc10 `/webdav` or `/dav/files/<username>` home folders into multiple sync pairs.
One pair for `/webdav/home` or `/dav/files/<username>/home` and another pair for every accepted share. The shares will be accessible at `/webdav/shares/` when the server side enables the spaces API.
Files can be identified using `<instance name>$<numeric storageid>!<fileid>` and moved to the correct sync pair.

### Stage-2.3
Enable spaces API in oc10:
- New clients will get a response from the spaces API and can set up new sync pairs.
- Legacy clients will still poll `/webdav` or `/dav/files/<username>` where they will see new subfolders instead of the users home. They will move down the users files into `/home` and shares into `/shares`. Custom sync pairs will no longer be available, causing the legacy client to leave local files in place. They can be picked up manually when installing a new client.

{{< /hint >}}

{{< hint warning >}}
**Alternative 2**
An additional `uuid` property used only to detect moves. A lookup by uuid is not necessary for this. The `/dav/meta` endpoint would still take the fileid. Clients  would use the `uuid` to detect moves and set up new sync pairs when migrating to a global namespace.
### Stage-2.1
Generate a `uuid` for every file as a file property. Clients can submit a `uuid` when creating files. The server will create a `uuid` if the client did not provide one.

### Stage-2.2
Roll out new clients that understand the spaces API and know how to convert local sync pairs for legacy oc10 `/webdav` or `/dav/files/<username>` home folders into multiple sync pairs.
One pair for `/webdav/home` or `/dav/files/<username>/home` and another pair for every accepted share. The shares will be accessible at `/webdav/shares/` when the server side enables the spaces API. Files can be identified using the `uuid`  and moved to the correct sync pair.

### Stage-3.1
When reading the files from ocis return the same `uuid`. It can be migrated to an extended attribute or it can be read from oc10. If users change it the client will not be able to detect a move and maybe other weird stuff happens. *What if the uuid gets lost on the server side due to a partial restore?* 

{{< /hint >}}

### Stage-3
Start oCIS backend and make read only tests on existing data using the `owncloud` storage driver which will read (and write)
- blobs from the same datadirectory layout as in ownCloud 10 and
- metadata from the ownCloud 10 databas
The oCIS share manager will read share infomation from the owncloud database as well.
- [ ] *we need a share manager that can read from the oc 10 db as well as from whatever new backend will be used for a pure oCIS setup. Currently, that would be the json file. Or that is migrated after all users have switched to oCIS. -- jfd*

### Stage-4
Test writing data with oCIS into the existing ownCloud 10 datafolder using the `owncloud` storage driver.

### Stage-5
Introduce reverse proxy and switch over early adoptors, let admins gain trust in the new backend by comparing metrics of the two systems and having it running in parallel.

### Stage-6
Voluntary transition period and subsequent hard deadline for all users

### Stage-7
Disable oc10 in the proxy, all requests are now handled by oCIS, shut down oc10 web servers and redis (or keep for calendar & contacts only? rip out files from oCIS?)

### Stage-8
User by user storage migration from owncloud driver to `ocis`/`s3ng`/`cephfs`...

### Stage-9
Migrate share data to &lt;yet to determine&gt; share manager backend and shut down owncloud database

### Stage-10
Profit! (db for file metadata no longer necessary, less maintenance effort)


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
├── einstein
│   ├── files
│   ├── files_trash
│   └── files_versions
│   …
├── marie
…
```

The *data directory* may also contain subfolders for owncloud 10 applications like `avatars`, `gallery`, `files_external` and `cache`.

When an objectstorage is used as the primary storage all file blobs are stored by their file id and a prefix, eg.: `urn:oid:<fileid>`.

The three types of blobs we need to migrate are stored in 
- `files` for file blobs, the current file content,
- `files_trashbin` for trashed files (and their versions) and 
- `files_versions` for file blobs of older versions.

### Filecache table

In both cases the file metadata, including a full replication of the file tree, is stored in the `oc_filecache`  table of an ownCloud 10 database. The primary key of a row is the file id. It is used to attach additional metadata like shares, tags, favorites or arbitrary file properties.

The `filecache` table itself has more metadata:

| Field              | Type          | Null | Key | Default | Extra          | Comment        | Migration         |
|--------------------|---------------|------|-----|---------|----------------|----------------|----------------|
| `fileid`           | bigint(20)    | NO   | PRI | NULL    | auto_increment |                | MUST become the oCIS `opaqueid` of a file reference. `ocis` driver stores it in extendet attributes and can use numbers as node ids on disk. for eos see note below table |
| `storage`          | int(11)       | NO   | MUL | 0       |                | *the filecache holds metadata for multiple storages* | corresponds to an oCIS *storage space* |
| `path`             | varchar(4000) | YES  |     | NULL    |                | *the path relative to the storages root* | MUST become the `path` relative to the storage root. `files` prefix needs to be trimmed. |
| `path_hash`        | varchar(32)   | NO   |     |         |                | *mysql once had problems indexing long paths, so we stored a hash for lookup by path. | - |
| `parent`           | bigint(20)    | NO   | MUL | 0       |                | *used to implement the hierarchy and listing children of a folder by id. redundant with `path`* | - |
| `name`             | varchar(250)  | YES  |     | NULL    |                | *basename of `path`*               | - |
| `mimetype`         | int(11)       | NO   |     | 0       |                | *joined with the `oc_mimetypes` table. only relevant for objectstorage deployments* | can be determined from blob / file extension |
| `mimepart`         | int(11)       | NO   |     | 0       |                | *"*               | can be determined from blob / file extension |
| `size`             | bigint(20)    | NO   |     | 0       |                | *same as blob size unless encryption is used*               | MAY become size, can be determined from blob |
| `mtime`            | bigint(20)    | NO   |     | 0       |                | *same as blob mtime*               | for files MAY become mtime (can be determined from blob as well), for directories MUST become tmtime |
| `encrypted`        | int(11)       | NO   |     | 0       |                | *encrypted flag*                | oCIS currently does not support encryption |
| `etag`             | varchar(40)   | YES  |     | NULL    |                | *used to propagate changes in a tree* | MUST be migrated (or calculated in the same way) to prevent clients from syncing unnecessarily |
| `unencrypted_size` | bigint(20)    | NO   |     | 0       |                | *same as blob size* | oCIS currently does not support encryption |
| `storage_mtime`    | bigint(20)    | NO   |     | 0       |                | *used to detect external storage changes* | oCIS delegates that to the storage providers and drivers |
| `permissions`      | int(11)       | YES  |     | 0       |                | *used as the basis for permissions. synced from disk when running a file scan. * | oCIS delegates that to the storage providers and drivers |
| `checksum`         | varchar(255)  | YES  |     | NULL    |                | *same as blob checksum* | SHOULD become the checksum in the storage provider. eos calculates it itself, `ocis` driver stores it in extendetd attributes |


> Note: for EOS a hot migration only works seamlessly if file ids in oc10 are already read from eos. otherwise either a mapping from the oc10 filecache file id to the new eos file id has to be created under the assumption that these id sets do not intersect or files and corresponding shares need to be exported and imported offline to generate a new set of ids. While this will preserve public links, user, group and even federated shares, old internal links may still point to different files because they contain the oc10 fileid 

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
| `file_source`   | bigint(20)   | YES  | MUL | NULL    |  | *cannot store uuid style file ids from ocis. when all users have migrated to ocis the share manager needs to be updated / migrated to a version that does.* | `Share.resource_id` [`ResourceId`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId) |
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

*TODO document how the reva OCM service currently persists the data*

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

*TODO clarify how OCM handles this andt where we store / configure this. It seems related to trusted IdPs*

### user data

Users are migrated in two steps:
1. They should all be authenticated using openid connect, which already moves them to a common identity management system.
2. To search share recipients, both, ownCloud 10 and oCIS need acces to the same user directory using eg. LDAP.

*TODO: add state to CS3 API, so we can 'disable' users*
*TODO: how do we map (sub) admins? -> map to roles & permissions*

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

### LDAP

TODO clarify if metadata from ldap & user_shibboleth needs to be migrated
-  the `dn` -> *owncloud internal username* mapping that currently lives in the `oc_ldap_user_mapping` table needs to move into a dedicated ownclouduuid attribute in the LDAP server. The idp should send it as a claim so the proxy does not have to look up the user using LDAP again. The username cannot be changed in ownCloud 10 and the oCIS provisioning API will not allow changing it as well. When we introduce the graph api we may allow changing usernames when all clients have moved to that api.

The problem is that the username in owncloud 10 and in oCIS also need to be the same, which might not be the case when the ldap mapping used a different column. In that case we should add another owncloudusername attribute to the ldap server ...


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