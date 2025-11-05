# ocis

The ocis package includes the Infinite Scale runtime and commands for the Infinite Scale command-line interface (CLI), which are not bound to a service.

Table of Contents
=================

   * [Service Registry](#service-registry)
   * [Memory limits](#memory-limits)
   * [CLI Commands](#cli-commands)
      * [Backup CLI](#backup-cli)
      * [Cleanup Orphaned Shares](#cleanup-orphaned-shares)
      * [List Unified Roles](#list-unified-roles)
      * [Move Stuck Uploads](#move-stuck-uploads)
      * [Revisions CLI](#revisions-cli)
      * [Service Health](#service-health)
      * [Trash CLI](#trash-cli)

<!-- Created by https://github.com/ekalinin/github-markdown-toc -->

## Service Registry

This package also configures the service registry which will be used to look up the service addresses.

Available registries are:

-   nats-js-kv (default)
-   memory

To configure which registry to use, you have to set the environment variable `MICRO_REGISTRY`, and for all except `memory` you also have to set the registry address via `MICRO_REGISTRY_ADDRESS`.

## Memory limits

oCIS will automatically set the go native `GOMEMLIMIT` to `0.9`. To disable the limit set `AUTOMEMEMLIMIT=off`. For more information take a look at the official [Guide to the Go Garbage Collector](https://go.dev/doc/gc-guide).

## CLI Commands

The ocis package offers a variety of cli commands to monitor or repair ocis installations. All these commands have a common mandatory parameter: `--basePath` (or `-p`) which needs to point to a storage provider. Example paths are:

```bash
.ocis/storage/users          # bare metal installation
/var/tmp/ocis/storage/users  # docker installation
...
```

These paths can vary depending on your ocis installation.

All commands provide a `-h` / `--help` option. Use to print all available options.

### Backup CLI

The backup command allows inspecting the consistency of an ocis storage:

```bash
ocis backup consistency -p /base/path/storage/users
```

This will check the consistency of the storage and output a list of inconsistencies. Inconsistencies can be:

* **Orphaned Blobs**\
A blob in the blobstore that is not referenced by any file metadata.
* **Missing Blobs**\
A blob referenced by file metadata that is not present in the blobstore.
* **Missing Nodes**\
A node that is referenced by a symlink but doesn't exist.
* **Missing Link**\
A node that is not referenced by any symlink but should be.
* **Missing Files**\
A node that is missing essential files (such as the `.mpk` metadata file).
* **Missing/Malformed Metadata**\
A node that doesn't have any (or malformed) metadata.

This command provides additional options:

* `-b` / `--blobstore`\
Allows specifying the blobstore to use. Defaults to `ocis`. Empty blobs will not be checked. Can also be switched to `s3ng`, but needs addtional envvar configuration (see the `storage-users` service for more details).
* `--fail`\
Exits with non-zero exit code if inconsistencies are found. Useful for automation.

### Cleanup Orphaned Shares

When a shared space or directory got deleted, use the `shares cleanup` command to remove those share orphans. This can't be done automatically at the moment.

```bash
ocis shares cleanup
```

### List Unified Roles

This command simplifies the process of finding out which UID belongs to which role. The command using markdown as output format is:

```bash
ocis graph list-unified-roles --output-format md
```

The output of this command includes the following information for each role:

* `Name`\
  The human readable name of the role.
* `UID`\
  The unique identifier of the role.
* `Enabled`\
  Whether the role is enabled or not.
* `Description`\
  A short description of the role.
* `Condition`
* `Allowed Resource Actions`

**Example output (shortned)**

| #  |              LABEL               |                 UID                  | ENABLED  |                                     DESCRIPTION                                      |                         CONDITION                         |         ALLOWED RESOURCE ACTIONS         |
|:--:|:--------------------------------:|:------------------------------------:|:--------:|:------------------------------------------------------------------------------------:|:---------------------------------------------------------:|:----------------------------------------:|
| 1  |              Viewer              | b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5 | enabled  |                                  View and download.                                  |                   exists @Resource.File                   |     libre.graph/driveItem/path/read      |
|    |                                  |                                      |          |                                                                                      |                  exists @Resource.Folder                  |     libre.graph/driveItem/quota/read     |
|    |                                  |                                      |          |                                                                                      |  exists @Resource.File && @Subject.UserType=="Federated"  |    libre.graph/driveItem/content/read    |
|    |                                  |                                      |          |                                                                                      | exists @Resource.Folder && @Subject.UserType=="Federated" |   libre.graph/driveItem/children/read    |
|    |                                  |                                      |          |                                                                                      |                                                           |    libre.graph/driveItem/deleted/read    |
|    |                                  |                                      |          |                                                                                      |                                                           |     libre.graph/driveItem/basic/read     |
| 2  |         ViewerListGrants         | d5041006-ebb3-4b4a-b6a4-7c180ecfb17d | disabled |                     View, download and show all invite

### Move Stuck Uploads

In some cases of saturated disk usage, Infinite Scale metadata may become stuck. This can occur when file metadata is being moved to its final destination after file operations. This issue was primarily seen with shares, where uploaded files could not be accessed. The required filename parameter aligns with Infinite Scale's internal processes and is used to complete the formerly stuck move action.

```bash
ocis shares move-stuck-upload-blobs [--dry-run=false] -p /base/path/storage/users
```

This command provides additional options:

* `--dry-run` (default: `true`)\
Only print found files stuck in transition.\
Note: This is a safety measure. You must specify `--dry-run=false` for the command to be effective.

* `--filename` value (default: "received.json")\
File to move from `uploads/` to share manager metadata `blobs/`

### Revisions CLI

The revisions command allows removing the revisions of files in the storage.

```bash
ocis revisions purge -p /base/path/storage/users
```

It takes the `--resource-id` (or `--r`) parameter which specify the scope of the command:

* An empty string (default) removes all revisions from all spaces.
* A spaceID (like `d419032c-65b9-4f4e-b1e4-0c69a946181d\$44b5a63b-540c-4002-a674-0e9c833bbe49`) removes all revisions in that space.
* A resourceID (e.g. `d419032c-65b9-4f4e-b1e4-0c69a946181d\$44b5a63b-540c-4002-a674-0e9c833bbe49\!e8a73d49-2e00-4322-9f34-9d7f178577b2`) removes all revisions from that specific file.

This command provides additional options:

* `--dry-run` (default: `true`)\
Do not remove any revisions but print the revisions that would be removed.
* `-b` / `--blobstore`\
Allows specifying the blobstore to use. Defaults to `ocis`. Can be switched to `s3ng` but needs addtional envvar configuration (see the `storage-users` service for more details).
* `-v` / `--verbose`\
Prints additional information about the revisions that are removed.
* `--glob-mechanism` (default: `glob`\
(advanced) Allows specifying the mechanism to use for globbing. Can be `glob`, `list` or `workers`. In most cases the default `glob` does not need to be changed. If large spaces need to be purged, `list` or `workers` can be used to improve performance at the cost of higher cpu and ram usage. `list` will spawn 10 threads that list folder contents in parallel. `workers` will use a special globbing mechanism and multiple threads to achieve the best performance for the highest cost.

### Service Health

The service health CLI command allows checking the health status of a service. If there are no issues found, nothing health related will get printed.

```bash
ocis <service-name> health
```

**Examples**

* The `collaboration` service has been started but not configured and is therefore not in a healthy state:
  ```bash
  ocis collaboration health
  
  The WOPI secret has not been set properly in your config for collaboration. Make sure your /root/.ocis/config config contains the proper values (e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in the config/corresponding environment variable).
  ```

* The `antivirus` service has not been started, the health check responds accordingly:
  ```bash
  ocis antivirus health
  
  {"level":"fatal","service":"antivirus","error":"Get \"http://127.0.0.1:9277/healthz\": dial tcp 127.0.0.1:9277: connect: connection refused","time":"2024-10-28T17:47:54+01:00","message":"Failed to request health check"}
  ```

### Trash CLI

The trash cli allows removing empty folders from the trashbin. This should be used to speed up trash bin operations.

```bash
ocis trash purge-empty-dirs -p /base/path/storage/users
```

This command provides additional options:

* `--dry-run` (default: `true`)\
Do not remove any empty folders but print the empty folders that would be removed.
