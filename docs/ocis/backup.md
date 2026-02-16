---
title: "Backup Considerations"
date: 2024-05-07T10:31:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: backup.md
---

This small guide aims to shed some light on the internal Infinite Scale data structure. You can refer to it when you are trying to optimize your backups or if you are just curious about how Infinite Scale stores its data.

Note, as a prerequisite backing up Infinite Scale, the instance has to be fully shut down for the time being.

## Ocis Data Structure

Ocis stores its data in a folder that can be configured via the environment variable `OCIS_BASE_DATA_PATH`. Without further configuration, services derive from that path when they store data, though individual settings for certain data types can be configured.

The default value for the `OCIS_BASE_DATA_PATH` variable is `$HOME/.ocis` (or `/var/lib/ocis` when using the docker container. Note: Configuration data is by default stored in `/etc/ocis/` in the container.).

Inside this folder, Infinite Scale will store all its data in separate subdirectories. That includes metadata, configurations, queues and stores etc. The actual bytes of files (blobs) are handled by a so called blobstore, which also stores here by default. Depending on the used blobstore, the blobs need to be backed up separately, for example if S3 is used. Note: See special case for the `config` folder in a docker container.

### Base Data Path Overview

Listing the contents of the folder will return the following:
```bash
    ~/.ocis/:tree -L 1
.
├── config
├── idm
├── idp
├── nats
├── proxy
├── search
├── storage
├── thumbnails
└── web

10 directories, 0 files
```

The following sections describe the content and background of the subdirectories to decide if a backup is required or recommended and its effect when it is not backed up.

### `config`

Contains basic Infinite Scale configuration created by `ocis init`(Note: The location of the configuration folder can be specified with the `OCIS_CONFIG_DIR` environment variable, but for this document we will assume this variable is not set and the default is used.)

```bash
    ~/.ocis/config/:tree
.
└── ocis.yaml

1 directory, 1 file
```

* `ocis.yaml`:\
BACKUP RECOMMENDED. Holds Infinite Scale configuration data. The contents can vary depending on your environment variables. In general, most of this file can be recreated again by running `ocis init`. This will recreate secrets and certificates. However, if not backed up completely, some fields MUST be copied over from the old config manually to regain data access after a restore:

| Field Name | Envvar Name | Description | If not backed up |
| --- | --- | --- | --- |
| `idp.ldap.bind_password` | `OCIS_LDAP_BIND_PASSWORD` | Password for the idp | no logins possible |
| `idm.service_user_passwords.idp_password`| `IDM_IDPSVC_PASSWORD` | Same as above | no logins possible |
| `system_user_id` | `OCIS_SYSTEM_USER_ID` | The id of storage-system user | no logins possible |
| `idm.service_user_passwords.reva_password`| `IDM_REVASVC_PASSWORD` | The reva password | no logins possible |
| `auth_basic.auth_providers.ldap.bind_password` | `AUTH_BASIC_LDAP_BIND_PASSWORD` | Same as above | no logins possible |
| `users.drivers.ldap.bind_password` | `USERS_LDAP_BIND_PASSWORD` | Same as above | no logins possible |
| `groups.drivers.ldap.bind_password` | `GROUPS_LDAP_BIND_PASSWORD` | Same as above | no logins possible |
| `storage_users.mount_id` | `STORAGE_USERS_MOUNT_ID` | The mountid of the storage_users service | sharing data lost |
| `gateway.storage_registry.storage_users_mount_id` | `GATEWAY_STORAGE_USERS_MOUNT_ID` | Same as above | sharing data lost |

### `idm`

Note: This folder will not appear if you use an external idm. Refer to your idms documentation for backup details in this case.

Contains the data for the internal Infinite Scale identity management. See the [IDM README]({{< ref "../services/idm/_index.md" >}}) for more details.

```bash
    ~/.ocis/idm/:tree
.
├── ldap.crt
├── ldap.key
└── ocis.boltdb

1 directory, 3 files
```

* `ocis.boltdb`:\
BACKUP REQUIRED. This is the boltdb database that stores user data. Use `IDM_DATABASE_PATH` to specify its path. If not backed up, Infinite Scale will have no users, therefore also all data is lost.
* `ldap.crt`:\
BACKUP OPTIONAL. This is the certificate for the idm. Use `IDM_LDAPS_CERT` to specify its path. Will be auto-generated if not backed up.
* `ldap.key`:\
BACKUP OPTIONAL. This is the certificate key for the idm. Use `IDM_LDAPS_KEY` to specify its path. Will be auto-generated if not backed up.


### `idp`

Note: This folder will not appear if you use an external idp. Refer to your idp's documentation for backup details in this case.

Contains the data for the internal Infinite Scale identity provider. See the [IDP README]({{< ref "../services/idp/_index.md" >}}) for more details.

```bash
    ~/.ocis/idp/:tree
.
├── encryption.key
├── private-key.pem
└── tmp
    └── identifier-registration.yaml

2 directories, 3 files
```

* `encryption.key`:\
BACKUP RECOMMENDED. This is the encryption secret. Use `IDP_ENCRYPTION_SECRET_FILE` to specify its paths. Not backing this up will force users to relogin.
* `private-key.pem`:\
BACKUP RECOMMENDED. This is the encryption key. Use `IDP_SIGNING_PRIVATE_KEY_FILES` to specify its paths. Not backing this up will force users to relogin.
* `identifier-registration.yml`:\
BACKUP OPTIONAL. It holds configuration for oidc clients (web, desktop, ios, android). Will be recreated if not backed up.

### `nats`

Note: This folder will not appear if you use an external nats installation. In that case, data has to secured in alignment with the external installation.

Contains nats data for streams and stores. See the [NATS README]({{< ref "../services/nats/_index.md" >}}) for more details.

```bash
    ~/.ocis/nats/:tree -L 1
.
└── jetstream

```

* `jetstream`:\
BACKUP RECOMMENDED. This folder contains nats data about streams and key-value stores. Use `NATS_NATS_STORE_DIR` to specify its path. Not backing it up can break history for multiple (non-vital) features such as history or notifications. The Infinite Scale functionality is not impacted if omitted.

### `proxy`

Contains proxy service data. See the [PROXY README]({{< ref "../services/proxy/_index.md" >}}) for more details.

```bash
    ~/.ocis/proxy/:tree
.
├── server.crt
└── server.key

1 directory, 2 files
```

* `server.crt`:\
BACKUP OPTIONAL. This is the certificate for the http services. Use `PROXY_TRANSPORT_TLS_CERT` to specify its path.
* `server.key`:\
BACKUP OPTIONAL. This is the certificate key for the http services. Use `PROXY_TRANSPORT_TLS_KEY` to specify its path.

### `search`

Contains the search index. See the [SEARCH README]({{< ref "../services/search/_index.md" >}}) for more details.

```bash
    ~/.ocis/search/:tree -L 1
.
└── bleve

2 directories, 0 files
```

* `bleve`:\
BACKUP RECOMMENDED/OPTIONAL. This contains the search index. Can be specified via `SEARCH_ENGINE_BLEVE_DATA_PATH`. If not backed up, the search index needs to be recreated. This can take a long time depending on the amount of files.

### `storage`

Contains Infinite Scale meta (and blob) data, depending on the blobstore. See the [STORAGE-USERS README]({{< ref "../services/storage-users/_index.md" >}}) for more details.

```bash
    ~/.ocis/storage/:tree -L 1
.
├── metadata
├── ocm
└── users

4 directories, 0 files
```

* `metadata`:\
BACKUP REQUIRED. Contains system data. Path can be specified via `STORAGE_SYSTEM_OCIS_ROOT`. Not backing it up will remove shares from the system and will also remove custom settings.
* `ocm`:\
BACKUP REQUIRED/OMITABLE. Contains ocm share data. When not using ocm sharing, this folder does not need to be backed up.
* `users`:\
BACKUP REQUIRED. Contains user data. Path can be specified via `STORAGE_USERS_OCIS_ROOT`. Not backing it up will remove all spaces and all files. As result, you will have a configured but empty Infinite Scale instance, which is fully functional accepting new data. Old data is lost.

### `thumbnails`

Contains thumbnails data. See the [THUMBNAILS README]({{< ref "../services/thumbnails/_index.md" >}}) for more details.

```bash
    ~/.ocis/thumbnails/:tree -L 1
.
└── files
```

* `files`:\
OPTIONAL/RECOMMENDED. This folder contains prerendered thumbnails. Can be specified via `THUMBNAILS_FILESYSTEMSTORAGE_ROOT`. If not backed up, thumbnails will be regenerated automatically on access which leads to some load on the thumbnails service.

### `web`

Contains web assets such as custom logos, themes etc. See the [WEB README]({{< ref "../services/web/_index.md" >}}) for more details.

```bash
    ~/.ocis/web/:tree -L 1
.
└── assets

2 directories, 0 files
```

* `assets`:\
BACKUP RECOMMENDED/OMITABLE. This folder contains custom web assets. Can be specified via `WEB_ASSET_CORE_PATH`. If no custom web assets are used, there is no need for a backup. If those exist but are not backed up, they need to be reuploaded.

### `external services`

When using an external idp/idm/nats or blobstore, its data needs to be backed up separately. Refer to your idp/idm/nats/blobstore documentation for backup details.

## Backup Consistency Command

Infinite Scale now allows checking an existing backup for consistency. Use the command:
```bash
ocis backup consistency -p "<path-to-base-folder>"
```

`path-to-base-folder` needs to be replaced with the path to the storage providers base path. Should be same as the `STORAGE_USERS_OCIS_ROOT`

Use the `-b s3ng` option when using an external (s3) blobstore. Note: When using this flag, the path to the blobstore must be configured via envvars or a yaml file to match the configuration of the original instance. Consistency checks for other blobstores than `ocis` and `s3ng` are not supported at the moment.
