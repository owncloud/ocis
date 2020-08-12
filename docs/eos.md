---
title: "EOS"
date: 2020-02-27T20:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: eos.md
---

{{< toc >}}

OCIS can be configured to run on top of [eos](https://eos.web.cern.ch/). While the [eos documentation](http://eos-docs.web.cern.ch/) does cover a lot of topics it leaves out some details that you may have to either pull from various [docker containers](https://gitlab.cern.ch/eos/eos-docker), the [forums](https://eos-community.web.cern.ch/) or even the [source](https://github.com/cern-eos/eos) itself.

This document is a work in progress of the current setup.

## Docker dev environment for eos storage

### 1. Start eos & ocis containers

Start the eos cluster and ocis via the compose stack

```
docker-compose up -d
```

### 2. LDAP support

Configure the os to resolve users and groups using ldap

```
docker-compose exec -d ocis /start-ldap
```

Check that the os in the ocis container can now resolve einstein or the other demo users

```
$ docker-compose exec ocis id einstein
uid=20000(einstein) gid=30000(users) groups=30000(users),30001(sailing-lovers),30002(violin-haters),30007(physics-lovers)
```

We also need to restart the reva-users service so it picks up the changed environment. Without a restart it is not able to resolve users from LDAP.
```
docker-compose exec ocis ./bin/ocis kill reva-users
docker-compose exec ocis ./bin/ocis run reva-users
```

### 3. Home storage

Kill the home storage. By default it uses the `owncloud` storage driver. We need to switch it to the `eoshome` driver and a new layout:

```
docker-compose exec ocis ./bin/ocis kill reva-storage-home
docker-compose exec -e REVA_STORAGE_EOS_LAYOUT="{{substr 0 1 .Id.OpaqueId}}/{{.Id.OpaqueId}}" -e REVA_STORAGE_HOME_DRIVER=eoshome ocis ./bin/ocis run reva-storage-home
```

### 4. Home data provider

Kill the home data provider. By default it uses the `owncloud` storage driver. We need to switch it to the `eoshome` driver and a new layout:

```
docker-compose exec ocis ./bin/ocis kill reva-storage-home-data
docker-compose exec -e REVA_STORAGE_EOS_LAYOUT="{{substr 0 1 .Id.OpaqueId}}/{{.Id.OpaqueId}}" -e REVA_STORAGE_HOME_DATA_DRIVER=eoshome ocis ./bin/ocis run reva-storage-home-data
```

{{< hint info >}}
The difference between the *home storage* and the *home data provider* are that the former is responsible for metadata changes while the latter is responsible for actual data transfer. The *home storage* uses the cs3 api to manage a folder hierarchy, while the *home data provider* is responsible for moving bytes to and from the storage.
{{< /hint >}}

### 4. Frontend files namespace

Restart the reva frontend with a new namespace (pointing to the eos storage provider) for the dav files endpoint

```
docker-compose exec ocis ./bin/ocis kill reva-frontend
docker-compose exec -e DAV_FILES_NAMESPACE="/eos/" ocis ./bin/ocis run reva-frontend
```

## Verification

Login with `einstein / relativity`, upload a file to einsteins home and verify the file is there using

```
docker-compose exec ocis eos ls -l /eos/dockertest/reva/users/4/4c510ada-c86b-4815-8820-42cdf82c3d51/
-rw-r--r--   1 einstein users              10 Jul  1 15:24 newfile.txt
```

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

