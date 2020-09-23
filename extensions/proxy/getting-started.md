---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis-proxy
geekdocEditPath: edit/master/docs
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Installation

So far we are offering two different variants for the installation. You can choose between [Docker](https://www.docker.com/) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

### Docker

Docker images for ocis-reva are hosted on https://hub.docker.com/r/owncloud/ocis-proxy.

The `latest` tag always reflects the current master branch.

```console
docker pull owncloud/ocis-proxy
```

### Binaries

The pre-built binaries for different platforms are downloadable at https://download.owncloud.com/ocis/ocis-proxy/ . Specific releases are organized in separate folders. They are in sync which every release tag on GitHub. The binaries from the current master branch can be found in https://download.owncloud.com/ocis/ocis-proxy/testing/

```console
curl https://download.owncloud.com/ocis/ocis-proxy/1.0.0-beta1/ocis-proxy-1.0.0-beta1-darwin-amd64 --output ocis-proxy
chmod +x ocis-proxy
./ocis-proxy server
```

## Usage

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formated help output if you execute the binary via `ocis-proxy --help`.

### Server

The server command is used to start the http server. For further help please execute:

{{< highlight txt >}}
ocis-proxy server --help
{{< / highlight >}}
