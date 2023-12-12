# Storage-Users

Purpose and description to be added

## Deprecated Metadata Backend

Starting with ocis version 3.0.0, the default backend for metadata switched to messagepack. If the setting `STORAGE_USERS_OCIS_METADATA_BACKEND` has not been defined manually, the backend will be migrated to `messagepack` automatically. Though still possible to manually configure `xattrs`, this setting should not be used anymore as it will be removed in a later version.

## Graceful Shutdown

Starting with Infinite Scale version 3.1, you can define a graceful shutdown period for the `storage-users` service.

IMPORTANT: The graceful shutdown period is only applicable if the `storage-users` service runs as standalone service. It does not apply if the `storage-users` service runs as part of the single binary or as single Docker environment. To build an environment where the `storage-users` service runs as a standalone service, you must start two instances, one _without_ the `storage-users` service and one _only with_ the the `storage-users` service. Note that both instances must be able to communicate on the same network.

When hard-stopping Infinite Scale, for example with the `kill <pid>` command (SIGKILL), it is possible and likely that not all data from the decomposedfs (metadata) has been written to the storage which may result in an inconsistent decomposedfs. When gracefully shutting down Infinite Scale, using a command like SIGTERM, the process will no longer accept any write requests from _other_ services and will try to write the internal open  requests which can take an undefined duration based on many factors. To mitigate that situation, the following things have been implemented:

*   With the value of the environment variable `STORAGE_USERS_GRACEFUL_SHUTDOWN_TIMEOUT`, the `storage-users` service will delay its shutdown giving it time to finalize writing necessary data. This delay can be necessary if there is a lot of data to be saved and/or if storage access/thruput is slow. In such a case you would receive an error log entry informing you that not all data could be saved in time. To prevent such occurrences, you must increase the default value.

*   If a shutdown error has been logged, the command-line maintenance tool [Inspect and Manipulate Node Metadata](https://doc.owncloud.com/ocis/next/maintenance/commands/commands.html#inspect-and-manipulate-node-metadata) can help to fix the issue. Please contact support for details.

## CLI Commands

### Manage Unfinished Uploads

<!-- referencing: [oCIS FS] clean up aborted uploads https://github.com/owncloud/ocis/issues/2622 -->

When using Infinite Scale as user storage, a directory named `storage/users/uploads` can be found in the Infinite Scale data folder. This is an intermediate directory based on [TUS](https://tus.io) which is an open protocol for resumable uploads. Each upload consists of a _blob_ and a _blob.info_ file. Note that the term _blob_ is just a placeholder.

*   **If an upload succeeds**, the _blob_ file will be moved to the target and the _blob.info_ file will be deleted.

*   **In case of incomplete uploads**, the _blob_ and _blob.info_ files will continue to receive data until either the upload succeeds in time or the upload expires based on the `STORAGE_USERS_UPLOAD_EXPIRATION` variable, see the table below for details.

*   **In case of expired uploads**, the _blob_ and _blob.info_ files will _not_ be removed automatically. Thus a lot of data can pile up over time wasting storage space.

*   **In the rare case of a failure**, after the upload succeeded but the file was not moved to its target location, which can happen when postprocessing fails, the situation is the same as with expired uploads.

Example cases for expired uploads

*   When a user uploads a big file but the file exceeds the user-quota, the upload can't be moved to the target after it has finished. The file stays at the upload location until it is manually cleared.
*   If the bandwidth is limited and the file to transfer can't be transferred completely before the upload expiration time is reached, the file expires and can't be processed.

There are two commands available to manage unfinished uploads

```bash
ocis storage-users uploads <command>
```

```plaintext
COMMANDS:
   list     Print a list of all incomplete uploads
   clean    Clean up leftovers from expired uploads
```

#### Command Examples

Command to identify incomplete uploads

```bash
ocis storage-users uploads list
```

```plaintext
Incomplete uploads:
 - 455bd640-cd08-46e8-a5a0-9304908bd40a (file_example_PPT_1MB.ppt, Size: 1028608, Expires: 2022-08-17T12:35:34+02:00)
```

Command to clear expired uploads
```bash
ocis storage-users uploads clean
```

```plaintext
Cleaned uploads:
- 455bd640-cd08-46e8-a5a0-9304908bd40a (Filename: file_example_PPT_1MB.ppt, Size: 1028608, Expires: 2022-08-17T12:35:34+02:00)
```

### Purge Expired Space Trash-Bins Items

<!-- referencing: https://github.com/owncloud/ocis/pull/5500 -->

This command is about the trash-bin to get an overview of items, restore items and purging old items of `project` spaces (spaces that have been created manually) and `personal` spaces.

```bash
ocis storage-users trash-bin <command>
```

#### Purge-expired
```plaintext
COMMANDS:
   purge-expired     Purge all expired items from the trashbin
```

The configuration for the `purge-expired` command is done by using the following environment variables.

*   `STORAGE_USERS_PURGE_TRASH_BIN_USER_ID` is used to obtain space trash-bin information and takes the system admin user as the default which is the `OCIS_ADMIN_USER_ID` but can be set individually. It should be noted, that the `OCIS_ADMIN_USER_ID` is only assigned automatically when using the single binary deployment and must be manually assigned in all other deployments. The command only considers spaces to which the assigned user has access and delete permission.

*   `STORAGE_USERS_PURGE_TRASH_BIN_PERSONAL_DELETE_BEFORE` has a default value of `30 days`, which means the command will delete all files older than `30 days`. The value is human-readable, valid values are `24h`, `60m`, `60s` etc. `0` is equivalent to disable and prevents the deletion of `personal space` trash-bin files.

*   `STORAGE_USERS_PURGE_TRASH_BIN_PROJECT_DELETE_BEFORE` has a default value of `30 days`, which means the command will delete all files older than `30 days`. The value is human-readable, valid values are `24h`, `60m`, `60s` etc. `0` is equivalent to disable and prevents the deletion of `project space` trash-bin files.

#### List and Restore Trash-Bins Items

To authenticate the cli command use `OCIS_MACHINE_AUTH_API_KEY=<some-ocis-machine-auth-api-key>`. The `storage-users` cli tool uses the default address to establish the connection to the `gateway` service. If the connection is failed check your custom `gateway`
service `GATEWAY_GRPC_ADDR` configuration and set the same address to `storage-users` variable `OCIS_GATEWAY_GRPC_ADDR` or `STORAGE_USERS_GATEWAY_GRPC_ADDR`.

The ID sources:
-   'userID' in a `https://{host}/graph/v1.0/me`
-   personal 'spaceID' in a `https://{host}/graph/v1.0/me/drives?$filter=driveType+eq+personal`
-   project 'spaceID' in a `https://{host}/graph/v1.0/me/drives?$filter=driveType+eq+project`

```bash
NAME:
   ocis storage-users trash-bin list - Print a list of all trash-bin items for a space.

USAGE:
   ocis storage-users trash-bin list command [command options] ['userID' required] ['spaceID' required]
```

```bash
NAME:
   ocis storage-users trash-bin restore-all - Restore all trash-bin items for a space.

USAGE:
   ocis storage-users trash-bin restore-all command [command options] ['userID' required] ['spaceID' required]

COMMANDS:
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --option value, -o value  The restore option defines the behavior for a file to be restored, where the file name already already exists in the target space. Supported values are: 'skip', 'replace' and 'keep-both'. The default value is 'skip' overwriting an existing file.
```

```bash
NAME:
   ocis storage-users trash-bin restore - Restore a trash-bin item by ID.

USAGE:
   ocis storage-users trash-bin restore command [command options] ['userID' required] ['spaceID' required] ['itemID' required]

COMMANDS:
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --option value, -o value  The restore option defines the behavior for a file to be restored, where the file name already already exists in the target space. Supported values are: 'skip', 'replace' and 'keep-both'. The default value is 'skip' overwriting an existing file.
```

## Caching

The `storage-users` service caches stat, metadata and uuids of files and folders via the configured store in `STORAGE_USERS_STAT_CACHE_STORE`, `STORAGE_USERS_FILEMETADATA_CACHE_STORE` and `STORAGE_USERS_ID_CACHE_STORE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `redis`: Stores metadata in a configured Redis cluster.
  -   `redis-sentinel`: Stores metadata in a configured Redis Sentinel cluster.
  -   `etcd`: Stores metadata in a configured etcd cluster.
  -   `nats-js`: Stores metadata using the key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in production environments.

1.  Note that in-memory stores are by nature not reboot-persistent.
2.  Though usually not necessary, a database name can be configured for event stores if the event store supports this. Generally not applicable for stores of type `in-memory`, `redis` and `redis-sentinel`. These settings are blank by default which means that the standard settings of the configured store apply.
3.  The `storage-users` service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
4.  When using `redis-sentinel`, the Redis master to use is configured via `STORAGE_USERS_STAT_CACHE_STORE_NODES`, `STORAGE_USERS_FILEMETADATA_CACHE_STORE_NODES` and `STORAGE_USERS_ID_CACHE_STORE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.
