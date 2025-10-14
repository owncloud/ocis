---
title: "Getting Started"
date: 2020-02-27T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/getting-started
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

{{< toc >}}

## oCIS online demo

We have an oCIS demo instance running on [ocis.owncloud.com](https://ocis.owncloud.com) where you can get a first impression of it.

We also have some more variations of oCIS running and [continuously deployed]({{< ref "../deployment/continuous_deployment" >}}) to reflect different scenarios in that oCIS might be used.

## Run oCIS

We are distributing oCIS as binaries and Docker images.

{{< hint warning >}}
The examples in this document assume that oCIS is accessed from the same host as it is running on (`localhost`). If you would like
to access oCIS remotely please refer to the [Basic Remote Setup]({{< ref "../deployment/basic-remote-setup" >}}) section. Especially
to the notes about setting the `PROXY_HTTP_ADDR` and `OCIS_URL` environment variables.
{{< /hint >}}

You can find more deployment examples in the [deployment section]({{< ref "../deployment" >}}).

### Binaries

You can find the latest official release of oCIS at [our download mirror](https://download.owncloud.com/ocis/ocis/stable/) or on [GitHub](https://github.com/owncloud/ocis/releases).
The latest build from the master branch can be found at [our download mirrors daily section](https://download.owncloud.com/ocis/ocis/daily/). Pre-Releases are available at [our download mirrors testing section](https://download.owncloud.com/ocis/ocis/testing/).

To run oCIS as binary you need to download it first and then run the following commands.
For this example, assuming version 2.0.0-beta.5 of oCIS running on a Linux AMD64 host:

```console
# download
curl https://download.owncloud.com/ocis/ocis/testing/2.0.0-beta.5/ocis-2.0.0-beta.5-linux-amd64 --output ocis

# make binary executable
chmod +x ocis

# initialize a minimal oCIS configuration
./ocis init

# run with demo users
IDM_CREATE_DEMO_USERS=true ./ocis server
```

The default primary storage location is `~/.ocis` or `/var/lib/ocis` depending on the packaging format and your operating system user. You can change that value by configuration.

{{< hint info >}}
When you're using oCIS with self-signed certificates, you need to answer the question for certificate checking with "yes" or set the environment variable `OCIS_INSECURE=true`, in order to make oCIS work.
{{< /hint >}}

{{< hint warning >}}
oCIS by default relies on Multicast DNS (mDNS), usually via avahi-daemon. If your system has a firewall, make sure mDNS is allowed in your active zone.
{{< /hint >}}

{{< hint warning >}}

#### Open Files on macOS

The start command `./ocis server` starts a runtime which runs all oCIS services in one process. On MacOS we have very low limits for open files. oCIS needs more than the default 256. Please raise the limit to 1024 by typing `ulimit -n 1024` within the same cli session where you start ocis from.
{{< /hint >}}

### Docker

Docker images for oCIS are available on [Docker Hub](https://hub.docker.com/r/owncloud/ocis).

The `latest` tag always reflects the current master branch.

```console
docker pull owncloud/ocis
docker run --rm -it -v ocis-config:/etc/ocis owncloud/ocis init
docker run --rm -p 9200:9200 -v ocis-config:/etc/ocis -v ocis-data:/var/lib/ocis -e IDM_CREATE_DEMO_USERS=true owncloud/ocis
```

{{< hint info >}}
When you're using oCIS with self-signed certificates, you need to set the environment variable `OCIS_INSECURE=true`, in order to make oCIS work.
{{< /hint >}}

{{< hint warming >}}
When you're creating the [demo users]({{< ref "./demo-users" >}}) by setting `IDM_CREATE_DEMO_USERS=true`, you need to be sure that this instance is not used in production because the passwords are public.
{{< /hint >}}

{{< hint warning >}}
We are using named volumes for the oCIS configuration and oCIS data in the above example (`-v ocis-config:/etc/ocis -v ocis-data:/var/lib/ocis`). You could instead also use host bind-mounts instead, eg. `-v /some/host/dir:/var/lib/ocis`.

You cannot use bind mounts on MacOS, since extended attributes are not supported ([owncloud/ocis#182](https://github.com/owncloud/ocis/issues/182), [moby/moby#1070](https://github.com/moby/moby/issues/1070)).
{{< /hint >}}

## Usage

### Login to ownCloud Web

Open [https://localhost:9200](https://localhost:9200) and [login using one of the demo accounts]({{< ref "./demo-users" >}}).

### Basic Management Commands

The oCIS single binary contains multiple extensions and the `ocis` command helps you to manage them. You already used `ocis server` to run all available extensions in the [Run oCIS]({{< ref "#run-ocis" >}}) section. We now will show you some more management commands, which you may also explore by typing `ocis --help` or going to the [docs]({{< ref "../config" >}}).

To initialize the oCIS configuration:

{{< highlight txt >}}
ocis init
{{< / highlight >}}

To start oCIS server:

{{< highlight txt >}}
ocis server
{{< / highlight >}}

The list command prints all running oCIS services.
{{< highlight txt >}}
ocis list
{{< / highlight >}}

The version command prints the version of your installed oCIS.
{{< highlight txt >}}
ocis --version
{{< / highlight >}}

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exit code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
ocis health --help
{{< / highlight >}}
