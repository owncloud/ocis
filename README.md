# ownCloud Infinite Scale

[![Matrix](https://img.shields.io/matrix/ocis%3Amatrix.org?logo=matrix)](https://app.element.io/#/room/#ocis:matrix.org)
[![Build Status](https://drone.owncloud.com/api/badges/owncloud/ocis/status.svg)](https://drone.owncloud.com/owncloud/ocis)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=owncloud_ocis&metric=security_rating)](https://sonarcloud.io/dashboard?id=owncloud_ocis)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=owncloud_ocis&metric=coverage)](https://sonarcloud.io/dashboard?id=owncloud_ocis)
[![Acceptance Test Coverage](https://sonarcloud.io/api/project_badges/measure?project=owncloud-1_ocis_acceptance-tests&metric=coverage)](https://sonarcloud.io/summary/new_code?id=owncloud-1_ocis_acceptance-tests)
[![Go Report](https://goreportcard.com/badge/github.com/owncloud/ocis)](https://goreportcard.com/report/github.com/owncloud/ocis)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis?status.svg)](http://godoc.org/github.com/owncloud/ocis)
[![oCIS docker image](https://img.shields.io/docker/v/owncloud/ocis?label=oCIS%20docker%20image&logo=docker&sort=semver)](https://hub.docker.com/r/owncloud/ocis)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

- [ownCloud Infinite Scale](#owncloud-infinite-scale)
  - [Introduction](#introduction)
  - [Quickstart](#quickstart)
  - [Overview](#overview)
    - [Clients](#clients)
    - [Web Office Applications](#web-office-applications)
    - [Authentication](#authentication)
    - [Installation](#installation)
  - [Important Readings](#important-readings)
  - [Run ownCloud Infinite Scale](#run-owncloud-infinite-scale)
    - [Use the Official Documentation](#use-the-official-documentation)
    - [Use the ocis Repo as Source](#use-the-ocis-repo-as-source)
  - [Documentation](#documentation)
    - [Admin Documentation](#admin-documentation)
    - [Development Documentation](#development-documentation)
  - [Security](#security)
  - [Contributing](#contributing)
  - [End User License Agreement](#end-user-license-agreement)
  - [Copyright](#copyright)

## Introduction

ownCloud Infinite Scale (oCIS) is the new file sync & share platform that will be the foundation of your data management platform.

Make sure to download the [latest released version](https://download.owncloud.com/ocis/ocis/stable/?sort=time&order=desc) today!

## Quickstart

For details of the commands used see the [Minimalistic Evaluation Guide for oCIS with Docker](https://owncloud.dev/ocis/guides/ocis-mini-eval/).

```bash
mkdir -p $HOME/ocis/ocis-config \
mkdir -p $HOME/ocis/ocis-data
sudo chown -Rfv 1000:1000 $HOME/ocis/
docker pull owncloud/ocis
docker run --rm -it \
    --mount type=bind,source=$HOME/ocis/ocis-config,target=/etc/ocis \
    --mount type=bind,source=$HOME/ocis/ocis-data,target=/var/lib/ocis \
    owncloud/ocis init --insecure yes
docker run \
    --name ocis_runtime \
    --rm \
    -it \
    -p 9200:9200 \
    --mount type=bind,source=$HOME/ocis/ocis-config,target=/etc/ocis \
    --mount type=bind,source=$HOME/ocis/ocis-data,target=/var/lib/ocis \
    -e OCIS_INSECURE=true \
    -e PROXY_HTTP_ADDR=0.0.0.0:9200 \
    -e OCIS_URL=https://localhost:9200 \
    owncloud/ocis
```
Use as URL `localhost:9200` and the user/password printed.

## Overview

### Clients

Infinite Scale allows the following ownCloud clients:

*   [web](https://github.com/owncloud/web),
*   [Android](https://github.com/owncloud/android),
*   [iOS](https://github.com/owncloud/ios-app) and
*   [Desktop](https://github.com/owncloud/client/)

to synchronize and share file spaces with a scalable server backend based on [reva](https://reva.link/) using open and well-defined APIs like [WebDAV](http://www.webdav.org/) and [CS3](https://github.com/cs3org/cs3apis/).

### Web Office Applications

Infinite Scale can integrate web office applications such as:

*   [Collabora Online](https://github.com/CollaboraOnline/online),
*   [OnlyOffice Docs](https://github.com/ONLYOFFICE/DocumentServer) or
*   [Microsoft Office Online Server](https://owncloud.com/microsoft-office-online-integration-with-wopi/)

Collaborative editing is supported by the [WOPI application gateway](https://github.com/cs3org/wopiserver).

### Authentication

Users are authenticated via [OpenID Connect](https://openid.net/connect/) using either an external IdP like [Keycloak](https://www.keycloak.org/) or the embedded [LibreGraph Connect](https://github.com/libregraph/lico) identity provider.

### Installation

With focus on easy install and operation, Infinite Scale is delivered as a single binary or container that allows scaling from a Raspberry Pi to a Kubernetes cluster by changing the configuration and starting multiple services as needed. The multiservice architecture allows tailoring the functionality to your needs and reusing services that may already be in place like when using Keycloak. See the details below for various installation options.

## Important Readings

Before starting to set up an instance, we **highly** recommend reading the [Prerequisites](https://doc.owncloud.com/ocis/next/prerequisites/prerequisites.html), the [Deployment](https://doc.owncloud.com/ocis/next/deployment/) section and especially the [General Information](https://doc.owncloud.com/ocis/next/deployment/general/general-info.html) page describing and explaining information that is valid for all deployment types.

## Run ownCloud Infinite Scale

### Use the Official Documentation

See the [Install Infinite Scale on a Server](https://doc.owncloud.com/ocis/next/depl-examples/ubuntu-compose/ubuntu-compose-prod.html) for a production ready deployment starting with a Raspberry Pi, a single server or VM.

### Use the ocis Repo as Source

Use this method to build and run an instance with the latest code. This is only recommended for development purposes.

The minimum go version required is `1.24.13`.\
Note that you need a C compile environment installed as a prerequisite because some dependencies, like reva, have components that require C-Go libraries and toolchains. The command installing for debian based systems is: `sudo apt install build-essential`.

To build and run a local instance with demo users:

```console
# get the source
git clone git@github.com:owncloud/ocis.git

# enter the ocis dir
cd ocis

# generate assets
make generate

# build the binary
make -C ocis build

# initialize a minimal oCIS configuration
./ocis/bin/ocis init

# run with demo users
IDM_CREATE_DEMO_USERS=true ./ocis/bin/ocis server

# Open your browser on http://localhost:9200 to access the bundled web-ui
```

All batteries included: no external database, no external IDP needed!

## Documentation

### Admin Documentation
Refer to the [Admin Documentation - Introduction to Infinite Scale](https://doc.owncloud.com/ocis/next/) to get started with running oCIS in production.

### Development Documentation
See the [Development Documentation - Getting Started](https://owncloud.dev/ocis/development/getting-started/) to get an overview of [Requirements](https://owncloud.dev/ocis/development/getting-started/#requirements), the [repository structure](https://owncloud.dev/ocis/development/getting-started/#repository-structure) and [other starting points](https://owncloud.dev/ocis/development/getting-started/#starting-points).

## Security

See the [Security Aspects](https://doc.owncloud.com/ocis/next/security/security.html) for a general overview of security related topics.
If you find a security issue, please contact [security@owncloud.com](mailto:security@owncloud.com) first.

## Contributing

We are _very_ happy that oCIS does not require a Contributor License Agreement (CLA) as it is [Apache 2.0 licensed](LICENSE). We hope this will make it easier to contribute code. If you want to get in touch, most of the developers hang out in our [matrix channel](https://app.element.io/#/room/#ocis:matrix.org) or reach out to the [ownCloud central forum](https://central.owncloud.org/).

Infinite Scale is carefully internationalized so that everyone, no matter what language they speak, has a great experience. To achieve this, we rely on the help of volunteer translators. If you want to help, you can find the projects behind the following links:
 [Transifex for ownCloud web](https://app.transifex.com/owncloud-org/owncloud-web/translate/) and [Transifex for ownCloud](https://app.transifex.com/owncloud-org/owncloud/translate/) (Select the resource by filtering for `ocis-`).

Please always refer to our [Contribution Guidelines](https://github.com/owncloud/ocis/blob/master/CONTRIBUTING.md).

## End User License Agreement

Some builds of stable ownCloud Infinite Scale releases provided by ownCloud GmbH are subject to an [End User License Agreement](https://owncloud.com/license-owncloud-infinite-scale/).

## Copyright

```console
Copyright (c) 2020-2023 ownCloud GmbH <https://owncloud.com>
```
