---
title: "Getting Started"
date: 2020-02-27T20:35:00+01:00
weight: -15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Run oCIS

We are distributing oCIS as binaries and Docker images.

You can find more deployments examples in the [deployment section](https://owncloud.github.io/ocis/deployment/)

### Binaries

The binaries for different platforms are downloadable at [our download mirror](https://download.owncloud.com/ocis/ocis/) or on [GitHub](https://github.com/owncloud/ocis/releases). Latest binaries from the master branch can be found at [our download mirrors testing section](https://download.owncloud.com/ocis/ocis/testing/).

```console
# for mac
curl https://download.owncloud.com/ocis/ocis/testing/ocis-testing-darwin-amd64 --output ocis
# for linux
curl https://download.owncloud.com/ocis/ocis/testing/ocis-testing-linux-amd64 --output ocis
# make binary executable
chmod +x ocis
./ocis server
```

The default primary storage location is `/var/tmp/`. You can change that value by configuration.


### Docker

Docker images for oCIS are available on [Docker Hub](https://hub.docker.com/r/owncloud/ocis).

The `latest` tag always reflects the current master branch.

```console
docker pull owncloud/ocis
docker run --rm -ti -p 9200:9200 owncloud/ocis
```

## Usage

### Login to oCIS Web

Open https://localhost:9200 and login using one of the demo accounts:

```console
einstein:relativity
marie:radioactivity
richard:superfluidity
```

There are admin demo accounts:
```console
moss:vista
admin:admin
```

### Basic Commands

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formatted help output if you execute the binary via `ocis --help`.

The server command is used to start the http and debug server on two addresses within a single process. The http server is serving the general web service while the debug server is used for health check, readiness check and to server the metrics mentioned below.

To start oCIS server:

{{< highlight txt >}}
ocis server
{{< / highlight >}}

The list command prints all running oCIS services.
{{< highlight txt >}}
ocis list
{{< / highlight >}}

To stop a particular extension:
{{< highlight txt >}}
ocis server kill phoenix
{{< / highlight >}}

To start a particular extension:
{{< highlight txt >}}
ocis server run phoenix
{{< / highlight >}}

The version command prints the version of your installed oCIS.
{{< highlight txt >}}
ocis --version
{{< / highlight >}}

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exist code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
ocis health --help
{{< / highlight >}}