# ocis

The ocis package contains the Infinite Scale runtime and the commands for the Infinite Scale cli.

## Service registry

This package also configures the service registry which will be used to look up the service addresses. It defaults to mDNS. Keep that in mind when using systems with mDNS disabled by default (i.e. SUSE).

Available registries are:

-   nats
-   kubernetes
-   etcd
-   consul
-   memory
-   mdns (default)

To configure which registry to use, you have to set the environment variable `MICRO_REGISTRY`, and for all except `memory` and `mdns` you also have to set the registry address via `MICRO_REGISTRY_ADDRESS`.

### etcd

To authenticate the connection to the etcd registry, you have to set `ETCD_USERNAME` and `ETCD_PASSWORD`.

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

### Trash CLI

The trash cli allows removing empty folders from the trashbin. This should be used to speed up trash bin operations.

```bash
ocis trash purge-empty-dirs -p /base/path/storage/users
```

This command provides additional options:

* `--dry-run` (default: `true`)\
Do not remove any empty folders but print the empty folders that would be removed.
