---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/reva
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Installation

So far we are offering two different variants for the installation. You can choose between [Docker](https://www.docker.com/) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

### Docker

Docker images for ocis-reva are hosted on https://hub.docker.com/r/owncloud/ocis-reva.

The `latest` tag always reflects the current master branch.

```console
docker pull owncloud/ocis-reva
```

#### Dependencies

- Running ocis-reva currently needs a working Redis caching server
- The default storage location in the container is `/var/tmp/reva/data`. You may want to create a volume to persist the files in the primary storage

### Binaries

The pre-built binaries for different platforms are downloadable at https://download.owncloud.com/ocis/ocis-reva/ . Specific releases are organized in separate folders. They are in sync which every release tag on GitHub. The binaries from the current master branch can be found in https://download.owncloud.com/ocis/ocis-reva/testing/

```console
curl https://download.owncloud.com/ocis/ocis/1.0.0-beta1/ocis-reva-1.0.0-beta1-darwin-amd64 --output ocis-reva
chmod +x ocis
./ocis-reva sharing
```

#### Dependencies

- Running ocis currently needs a working Redis caching server
- The default promary storage location is `/var/tmp/reva/data`. You can change that value by configuration.

## Usage

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formated help output if you execute the binary via `ocis-reva --help`.

### Health

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exist code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
ocis-reva health --help
{{< / highlight >}}
