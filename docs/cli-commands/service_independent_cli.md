---
title: Service Independent CLI
date: 2025-11-13T00:00:00+00:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/cli-commands/
geekdocFilePath: service_independent_cli.md
---

This document describes ocis CLI commands that are **service independent**.

{{< toc >}}

For **service dependent** CLI commands, see the following services:

* [Auth-App]({{< ref "../services/auth-app/" >}})
* [Graph]({{< ref "../services/graph/" >}})
* [Postprocessing]({{< ref "../services/postprocessing/" >}})
* [Storage-Users]({{< ref "../services/storage-users/" >}})

## Common Parameters

The ocis package offers a variety of CLI commands for monitoring or repairing ocis installations. Most of these commands have common parameters such as:

* `--help` (or `-h`)\
  Use to print all available options.

* `--basePath` (or `-p`)\
  Needs to point to a storage provider, paths can vary depending on your ocis installation. Example paths are:
  ```bash
  .ocis/storage/users          # bare metal installation
  /var/tmp/ocis/storage/users  # docker installation
  ...
  ```

* `--dry-run`\
  This parameter, when available, defaults to `true` and must explicitly set to `false`.

* `--verbose` (or `-v`)\
  Get a more verbose output.

## List of CLI Commands

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

### Cleanup Orphaned Grants

Detect and optionally delete storage grants that have no corresponding share-manager entry.

Sharing in ocis relies on two truths. The share manager and the grants. When a share is created, ocis will 

1. Create a grant for the specific file or folder.\
This grant is _checked when access to the file is requested_.

2. Create an entry in the `created.json`/`received.json` files of the specific user.\
These files are _checked whenever shares are listed_.

The process for creating a share is as follows: first, ocis creates the grant, and then adds the share entry. The reverse order is followed when deleting a share. This means that if the second step fails, the grant will still be present. This can be visually confirmed in the webUI. The webUI details of the "share" section will show an error fetching information for orphan grants.

The following command fixes the problem of orhaned grants.

Usage:
```bash
ocis shares clean-orphaned-grants \
  --service-account-id "<id>" \
  --service-account-secret "<secret>" \
  [--force] \
  [--space-id "<space-opaque-id>"] \
  [--dry-run=false]
```

Notes:
- `--dry-run`\
Defaults to `true` (no deletions). Set to `false` to remove orphaned grants.
- `--space-id`\
Limit the scan to a specific storage space (opaque ID).
- `--force`\
Force removal of suspected orphans even when listing shares fails.
- Public links are not touched.

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
