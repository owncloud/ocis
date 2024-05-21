---
title: "Backup Considerations"
date: 2024-05-07T10:31:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: backup.md
---

This small guide aims to shed some light on the internal ocis data structure. You can refer to it when you are trying to optimize your backups or if you just curious how ocis stores its data.

## Ocis data structure

Ocis stores its data in folder that can be configured via the envvar `OCIS_BASE_DATA_PATH`.

The default value for the `OCIS_BASE_DATA_PATH` variable is `$HOME/.ocis` (or `/etc/dev/ocis` when using the docker container.

Inside this folder ocis will store all its data. That includes metadata, configuration, queues and stores. The actual bytes of the uploaded files are also stored here by default. If an s3 store is used as blobstore, the blobs need to be backed up seperately.

### Base Data Path overview

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

### `config`

Contains basic ocis configuration created by `ocis init`(Note: The location of the configuration folder can be specified with the `OCIS_CONFIG_DIR` envvar but for this readme we will assume this envvar is not set)

```bash
    ~/.ocis/config/:tree
.
└── ocis.yaml

1 directory, 1 file
```

* `ocis.yaml`: BACKUP RECOMMENDED. Holds ocis configuration data. The contents can vary depending on your environment variables. In general most of this file can be recreated again by running `ocis init`. This will recreate secrets and certificates. However if not backed up completely, some fields MUST be copied over from the old config:

| Field Name | Envvar Name | Description | If not backed up |
| --- | --- | --- | --- |
| `idp.ldap.bind_password` | `OCIS_LDAP_BIND_PASSWORD` | Password for the idp | no logins possible |
| `idm.service_user_passwords.idp_password`| `IDM_IDPSVC_PASSWORD` | Same as above | no logins possible |
| `system_user_id` | `OCIS_SYSTEM_USER_ID` | The id of storage-system user | no logins possible |
| `graph.identity.ldap.bind_password` | `GRAPH_LDAP_BIND_PASSWORD` | The password for the idm | no logins possible |
| `idm.service_user_passwords.idm_password` | `IDM_SVC_PASSWORD` | Same as above | no logins possible |
| `idm.service_user_passwords.reva_password`| `IDM_REVASVC_PASSWORD` | The reva password | no logins possible |
| `auth_basic.auth_providers.ldap.bind_password` | `AUTH_BASIC_LDAP_BIND_PASSWORD` | Same as above | no logins possible |
| `users.drivers.ldap.bind_password` | `USERS_LDAP_BIND_PASSWORD` | Same as above | no logins possible |
| `groups.drivers.ldap.bind_password` | `GROUPS_LDAP_BIND_PASSWORD` | Same as above | no logins possible |
| `storage_users.mount_id` | `STORAGE_USERS_MOUNT_ID` | The mountid of the storage_users service | sharing data lost |
| `gateway.storage_registry.storage_users_mount_id` | `GATEWAY_STORAGE_USERS_MOUNT_ID` | Same as above | sharing data lost |

### `idm`

Note: this folder will not appear if you use an external idm. Refer to your idms documentation for backup details in this case.

Contains the data for the internal ocis identity management. See IDM README.

```bash
    ~/.ocis/idm/:tree
.
├── ldap.crt
├── ldap.key
└── ocis.boltdb

1 directory, 3 files
```

* `ocis.boltdb`: BACKUP REQUIRED. This is the boltdb database that stores user data. Use `IDM_DATABASE_PATH` to specify its path. If not backed up, ocis will have no users, therefore also all data is lost.
* `ldap.crt`: BACKUP OPTIONAL. This is the certificate for the idm. Use `IDM_LDAPS_CERT` to specify its path. Will be auto-generated if not backed up.
* `ldap.key`: BACKUP OPTIONAL. This is the certificate key for the idm. Use `IDM_LDAPS_KEY` to specify its path. Will be auto-generated if not backed up.


### `idp`

Note: this folder will not appear if you use an external idp. Refer to your idps documentation for backup details in this case.

Contains the data for the internal ocis identity provider. SEE IDP README.

```bash
    ~/.ocis/idp/:tree
.
├── encryption.key
├── private-key.pem
└── tmp
    └── identifier-registration.yaml

2 directories, 3 files
```

* `encryption.key`: BACKUP OPTIONAL. This is the encryption secret. Use `IDP_ENCRYPTION_SECRET_FILE` to specify its paths. Will be auto-generated if not backed up.
* `private-key.pem`: BACKUP OPTIONAL. This is the encryption key. Use `IDP_SIGNING_PRIVATE_KEY_FILES` to specify its paths. Will be auto-generated if not backed up.
* `identifier-registration.yml`: BACKUP RECOMMENDED. It holds temporary data of active sessions. Not backing this up will force users to relogin.

### `nats`

Note: this folder will not appear if you use an external nats installation

Contains nats data for streams and stores. SEE NATS README.

```bash
    ~/.ocis/nats/:tree -L 1
.
└── jetstream

```

* `jetstream`: BACKUP RECOMMENDED. This folder contains nats data about streams and key-value stores. Use `NATS_NATS_STORE_DIR` to specify its part. Not backing it up can break history for multiple (non-vital) features such as history or notifications. Ocis functionality is not impacted.

### `proxy`

Contains proxy service data. SEE PROXY README.

```bash
    ~/.ocis/proxy/:tree
.
├── server.crt
└── server.key

1 directory, 2 files
```

* `server.crt`: BACKUP OPTIONAL. This is the certificate for the http services. Use `PROXY_TRANSPORT_TLS_CERT` to specify its path.
* `server.key`: BACKUP OPTIONAL. This is the certificate key for the http services. Use `PROXY_TRANSPORT_TLS_KEY` to specify its path.

### `search`

Contains the search index.

```bash
    ~/.ocis/search/:tree -L 1
.
└── bleve

2 directories, 0 files
```

* `bleve`: BACKUP RECOMMENDED/OPTIONAL. This contains the search index. Can be specified via `SEARCH_ENGINE_BLEVE_DATA_PATH`. If not backing it up the search index needs to be recreated. This can take a long time depending on amount of files.

### `storage`

Contains ocis meta (and blob) data.

```bash
    ~/.ocis/storage/:tree -L 1
.
├── metadata
├── ocm
└── users

4 directories, 0 files
```

* `metadata`: BACKUP REQUIRED. Contains system data. Path can be specified via `STORAGE_SYSTEM_OCIS_ROOT`. Not backing it up will remove shares from the system and will also remove custom settings.
* `ocm`: BACKUP REQUIRED/OMITABLE. Contains ocm share data. When not using ocm sharing, this folder does not need to be backed up.
* `users`: BACKUP REQUIRED. Contains user data. Path can be specified via `STORAGE_USERS_OCIS_ROOT`. Not backing it up will remove all spaces and all files. Rendering ocis empty, but functionally.

### `thumbnails`

Contains thumbnails data.

```bash
    ~/.ocis/thumbnails/:tree -L 1
.
└── files
```

* `files`: OPTIONAL/RECOMMENDED. This folder contains prerendered thumbnails. Can be specified via `THUMBNAILS_FILESYSTEMSTORAGE_ROOT`. If not backed up thumbnails will be regenerated automatically which leads to some load on thumbnails service.

### `web`

Contains web assests such as custom logos.

```bash
    ~/.ocis/web/:tree -L 1
.
└── assets

2 directories, 0 files
```

* `assests`: BACKUP RECOMMENDED/OMITABLE. This folder contains custom web assests. Can be specified via `WEB_ASSET_CORE_PATH`. If no custom web assets are used, there is no need for a backup. If those exist but are not backed up, they need to be reuploaded.

### `external services`

When using an external idp/idm/nats or blobstore its data needs to be backed up separately. Refer to your idp/idm/nats/blobstore documentation for backup details.

