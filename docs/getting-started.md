---
title: "Getting Started"
date: 2020-02-27T20:35:00+01:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Installation

So far we are offering two different variants for the installation. You can choose between [Docker](https://www.docker.com/) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

### Docker

Docker images for ocis are hosted on https://hub.docker.com/r/owncloud/ocis.

The `latest` tag always reflects the current master branch.

```console
docker pull owncloud/ocis
```

#### Dependencies

- Running ocis currently needs a working Redis caching server
- The default storage location in the container is `/var/tmp/reva/data`. You may want to create a volume to persist the files in the primary storage

#### Docker compose

You can use our docker-compose [playground example](https://github.com/owncloud-docker/compose-playground/tree/master/ocis) to run ocis with dependencies with a single command in a docker network.

```console
git clone git@github.com:owncloud-docker/compose-playground.git
cd compose-playground/ocis
docker-compose -f ocis.yml -f ../cache/redis-ocis.yml up
```

### Binaries

The pre-built binaries for different platforms are downloadable at https://download.owncloud.com/ocis/ocis/ . Specific releases are organized in separate folders. They are in sync which every release tag on GitHub. The binaries from the current master branch can be found in https://download.owncloud.com/ocis/ocis/testing/

```console
curl https://download.owncloud.com/ocis/ocis/1.0.0-beta1/ocis-1.0.0-beta1-darwin-amd64 --output ocis
chmod +x ocis
./ocis server
```

#### Dependencies

- Running ocis currently needs a working Redis caching server
- The default promary storage location is `/var/tmp/reva/data`. You can change that value by configuration.

## Quickstart for Developers

Following https://github.com/owncloud/ocis#development

```console
git clone https://github.com/owncloud/ocis.git
cd ocis
make generate build
```

Open https://localhost:9200 and login using one of the demo accounts:

```console
einstein:relativity
marie:radioactivty
richard:superfluidity
```

## Runtime

Included with the ocis binary is embedded a go-micro runtime that is in charge of starting services as a fork of the master process. This provides complete control over the services. Ocis extensions can be added as part of this runtime.

```console
./bin/ocis micro
```

This will currently boot:

```console
com.owncloud.api
com.owncloud.http.broker
com.owncloud.proxy
com.owncloud.registry
com.owncloud.router
com.owncloud.runtime
com.owncloud.web
go.micro.http.broker
```

Further ocis extensions can be added to the runtime via the ocis command like:

```console
./bin/ocis hello
```

Which will register:

```console
com.owncloud.web.hello
com.owncloud.api.hello
```

To the list of available services.

