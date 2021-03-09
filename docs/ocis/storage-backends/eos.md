---
title: "EOS"
date: 2020-02-27T20:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: eos.md
---

{{< toc >}}

oCIS can be configured to run on top of [eos](https://eos.web.cern.ch/). While the [eos documentation](http://eos-docs.web.cern.ch/) does cover a lot of topics, it leaves out some details that you may have to either pull from various [docker containers](https://gitlab.cern.ch/eos/eos-docker), the [forums](https://eos-community.web.cern.ch/) or even the [source](https://github.com/cern-eos/eos) itself.

This document is a work in progress of the current setup.

## Docker dev environment for eos storage

We begin with the docker-compose.yml found in https://github.com/owncloud/ocis/tree/master/ocis/ and
switch it to eos-storage.

### 1. Start eos & ocis containers

Start the eos cluster and ocis via the compose stack.
```
docker-compose up -d
```

{{< hint info >}}
The first time the **ocis** container starts up, it will compile ocis from scratch which can take a while.
To follow progress, run `docker-compose logs -f --tail=10 ocis`
{{< /hint >}}

### 2. LDAP support

Configure the OS to resolve users and groups using ldap

```
docker-compose exec -d ocis /start-ldap
```

Check that the OS in the ocis container can now resolve einstein or the other demo users

```
$ docker-compose exec ocis id einstein
uid=20000(einstein) gid=30000(users) groups=30000(users),30001(sailing-lovers),30002(violin-haters),30007(physics-lovers)
```

{{< hint info >}}
If the user is not found at first you might need to wait a few more minutes in case the ocis container is still compiling.
{{< /hint >}}

We also need to restart the storage-userprovider service, so it picks up the changed environment. Without a restart it is not able to resolve users from LDAP.
```
docker-compose exec ocis ./bin/ocis kill storage-userprovider
docker-compose exec ocis ./bin/ocis run storage-userprovider
```

### 3. Home storage

Kill the home storage. By default it uses the `ocis` storage driver. We need to switch it to the `eoshome` driver:

```
docker-compose exec ocis ./bin/ocis kill storage-home
docker-compose exec -e STORAGE_HOME_DRIVER=eoshome ocis ./bin/ocis run storage-home
```

### 4. Users storage

Kill the users storage. By default it uses the `ocis` storage driver. We need to switch it to the `eos` driver:

```
docker-compose exec ocis ./bin/ocis kill storage-users
docker-compose exec -e STORAGE_USERS_DRIVER=eos ocis ./bin/ocis run storage-users
```

### 5. Metadata storage

First we need to create the metadata root in eos and set an owner:
```
docker-compose exec ocis eos mkdir -p /eos/dockertest/ocis/metadata
docker-compose exec ocis eos chown 2:2 /eos/dockertest/ocis/metadata
```

{{< hint info >}}
The uid and gid `2` are referencing the user `daemon` inside the ocis container. That user is also configured when restarting the accounts service later. For production systems you should create a dedicated user for the metadata storage.
{{< /hint >}}

Kill the metadata storage. By default it uses the `ocis` storage driver. We need to switch it to the `eos` driver:

```
docker-compose exec ocis ./bin/ocis kill storage-metadata
docker-compose exec -e STORAGE_METADATA_DRIVER=eos -e STORAGE_METADATA_ROOT=/eos/dockertest/ocis/metadata ocis ./bin/ocis run storage-metadata
```


### 6. Accounts service

Kill the accounts service. By default it uses the `ocis` storage driver. We need to switch it to the `eos` driver:

```
docker-compose exec ocis ./bin/ocis kill accounts
docker-compose exec -e ACCOUNTS_SERVICE_USER_USERNAME=daemon -e ACCOUNTS_SERVICE_USER_UID=2 -e ACCOUNTS_SERVICE_USER_GID=2 ocis ./bin/ocis run accounts
```

## Verification

Login with `einstein / relativity`, upload a file to einsteins home and verify the file is there using

```
docker-compose exec ocis eos ls -l /eos/dockertest/reva/users/4/4c510ada-c86b-4815-8820-42cdf82c3d51/
-rw-r--r--   1 einstein users              10 Jul  1 15:24 newfile.txt
```

If the problem persists, please check the [troubleshooting section about uploads](#creation-and-upload-of-files-does-not-work).

## Further exploration

EOS has a built in shell that you can enter using
```
$ docker-compose exec mgm-master eos
# ---------------------------------------------------------------------------
# EOS  Copyright (C) 2011-2019 CERN/Switzerland
# This program comes with ABSOLUTELY NO WARRANTY; for details type `license'.
# This is free software, and you are welcome to redistribute it
# under certain conditions; type `license' for details.
# ---------------------------------------------------------------------------
EOS_INSTANCE=eostest
EOS_SERVER_VERSION=4.6.5 EOS_SERVER_RELEASE=1
EOS_CLIENT_VERSION=4.6.5 EOS_CLIENT_RELEASE=1
EOS Console [root://localhost] |/> help
access               Access Interface
accounting           Accounting Interface
acl                  Acl Interface
archive              Archive Interface
attr                 Attribute Interface
backup               Backup Interface
clear                Clear the terminal
cd                   Change directory
chmod                Mode Interface
chown                Chown Interface
config               Configuration System
console              Run Error Console
cp                   Cp command
debug                Set debug level
exit                 Exit from EOS console
file                 File Handling
fileinfo             File Information
find                 Find files/directories
newfind              Find files/directories (new implementation)
fs                   File System configuration
fsck                 File System Consistency Checking
fuse                 Fuse Mounting
fusex                Fuse(x) Administration
geosched             Geoscheduler Interface
group                Group configuration
health               Health information about system
help                 Display this text
info                 Retrieve file or directory information
inspector            Interact with File Inspector
io                   IO Interface
json                 Toggle JSON output flag for stdout
license              Display Software License
ls                   List a directory
ln                   Create a symbolic link
map                  Path mapping interface
member               Check Egroup membership
mkdir                Create a directory
motd                 Message of the day
mv                   Rename file or directory
node                 Node configuration
ns                   Namespace Interface
pwd                  Print working directory
quit                 Exit from EOS console
quota                Quota System configuration
reconnect            Forces a re-authentication of the shell
recycle              Recycle Bin Functionality
rmdir                Remove a directory
rm                   Remove a file
role                 Set the client role
route                Routing interface
rtlog                Get realtime log output from mgm & fst servers
silent               Toggle silent flag for stdout
space                Space configuration
stagerrm             Remove disk replicas of a file if it has tape replicas
stat                 Run 'stat' on a file or directory
squash               Run 'squashfs' utility function
test                 Run performance test
timing               Toggle timing flag for execution time measurement
touch                Touch a file
token                Token interface
tracker              Interact with File Tracker
transfer             Transfer Interface
version              Verbose client/server version
vid                  Virtual ID System Configuration
whoami               Determine how we are mapped on server side
who                  Statistics about connected users
?                    Synonym for 'help'
.q                   Exit from EOS console
EOS Console [root://localhost] |/>
```

But this is a different adventure. See the links at the top of this page for other sources of information on eos.

## Cleaning up

To clean up and start completely from scratch, run `docker-compose down -v`.
Then delete the local "bin" folder as root which contains the ocis binaries compiled by the "ocis" docker.

## Troubleshooting

### Docker-compose exits right away

When running `docker-compose up -d` ocis exits right away.

You can check the error code using `docker-compose ps` and investigate further by running only ocis again using `docker-compose up ocis` (without `-d` so you can see what is going on in the foreground).
One reason might be that the binary was already built but does not match the container env. Try running `make clean` before running `docker-compose up ocis` so it gets built inside the container.

### Where are the logs ?

The ocis logs can be accessed using `docker-compose logs ocis`. Add `-f` for following.

### How do I update a service in the ocis container?

1. `docker-compose exec ocis make clean build` to update the binary
2. `docker-compose exec ocis ./bin/ocis kill <service>` to kill the service
3. `docker-compose exec ocis ./bin/ocis run <service>` to start the service. Do not forget to set any env vars, eg.
  `docker-compose exec -e STORAGE_HOME_DRIVER=eoshome -e STORAGE_DRIVER_EOS_LAYOUT="{{substr 0 1 .Id.OpaqueId}}/{{.Id.OpaqueId}}"  ocis ./bin/ocis run storage-home`

### Creation and upload of files does not work

If the upload did not work, please check the status of the eos space using the command `docker-compose exec mgm-master eos fs ls`.
In case the default space appears as offline, run `docker-compose exec mgm-master eos space set default on`.

### Uploading big files appears to hang

Please note that the uploads first go into the "ocis" docker and land in its "/tmp" folder, then gets copied over to the EOS docker using `xrdcopy`.
This is why uploading first transfers all bytes and then seem to hang for a while during the final copy.

### Running out of space quickly

The EOS dockers are configured with replication, so every file uploaded there will be replicated 4 times,
so make sure there is enough physical space on disk when testing.

Also please note that older failed uploads might still be present in the "/tmp" directory of the "ocis" container.

