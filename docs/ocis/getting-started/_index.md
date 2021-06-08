---
title: "Getting Started"
date: 2020-02-27T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/getting-started
geekdocFilePath: _index.md
---

{{< toc >}}

## oCIS online demo

We have an oCIS demo instance running on [ocis.owncloud.com](https://ocis.owncloud.com) where you can get a first impression of it.

We also have some more variations of oCIS running and [continuously deployed]({{< ref "../deployment/continuous_deployment" >}}) to reflect different scenarios in that oCIS might be used.

## Run oCIS

We are distributing oCIS as binaries and Docker images.

You can find more deployment examples in the [deployment section]({{< ref "../deployment" >}}).

### Binaries

You can find the latest official release of oCIS at [our download mirror](https://download.owncloud.com/ocis/ocis/) or on [GitHub](https://github.com/owncloud/ocis/releases).
The latest build from the master branch can be found at [our download mirrors testing section](https://download.owncloud.com/ocis/ocis/testing/).

To run oCIS as binary you need to download it first and then run the following commands.
For this example, assuming version 1.7.0 of oCIS running on a Linux AMD64 host:

```console
# download
curl https://download.owncloud.com/ocis/ocis/1.7.0/ocis-1.7.0-linux-amd64 --output ocis

# make binary executable
chmod +x ocis

# run
./ocis server
```

The default primary storage location is `/var/tmp/ocis`. You can change that value by configuration.

{{< hint warning >}}
oCIS by default relies on Multicast DNS (mDNS), usually via avahi-daemon. If your system has a firewall, make sure mDNS is allowed in your active zone.
{{< /hint >}}

### Docker

Docker images for oCIS are available on [Docker Hub](https://hub.docker.com/r/owncloud/ocis).

The `latest` tag always reflects the current master branch.

```console
docker pull owncloud/ocis
docker run --rm -ti -p 9200:9200 owncloud/ocis
```

## Usage

### Login to ownCloud Web

Open [https://localhost:9200](https://localhost:9200) and [login using one of the demo accounts]({{< ref "./demo-users" >}}).

### Basic Management Commands

The oCIS single binary contains multiple extensions and the `ocis` command helps you to manage them. You already used `ocis server` to run all available extensions in the [Run oCIS]({{< ref "#run-ocis" >}}) section. We now will show you some more management commands, which you may also explore by typing `ocis --help` or going to the [docs]({{< ref "../configuration" >}}).

To start oCIS server:

{{< highlight txt >}}
ocis server
{{< / highlight >}}

The list command prints all running oCIS extensions.
{{< highlight txt >}}
ocis list
{{< / highlight >}}

To stop a particular extension:
{{< highlight txt >}}
ocis kill web
{{< / highlight >}}

To start a particular extension:
{{< highlight txt >}}
ocis run web
{{< / highlight >}}

The version command prints the version of your installed oCIS.
{{< highlight txt >}}
ocis --version
{{< / highlight >}}

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exist code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
ocis health --help
{{< / highlight >}}
