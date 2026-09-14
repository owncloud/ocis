# ownCloud Infinite Scale

[![Matrix](https://img.shields.io/matrix/ocis%3Amatrix.org?logo=matrix)](https://app.element.io/#/room/#ocis:matrix.org)
[![Build Status](https://github.com/owncloud/ocis/actions/workflows/acceptance-tests.yml/badge.svg)](https://github.com/owncloud/ocis/actions/workflows/acceptance-tests.yml)
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

See the [Install Infinite Scale on a Server](https://doc.owncloud.com/ocis/next/admin/depl-examples/ubuntu-compose/ubuntu-compose-prod.html) for a production ready deployment starting with a Raspberry Pi, a single server or VM.

### Use the ocis Repo as Source

Use this method to build and run an instance with the latest code. This is only recommended for development purposes.

The minimum go version required is `1.25.10`.\
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


## 🌐 Web Resources & Aesthetic Symbols Index
- [SYM 1F618](https://aestheticsymbols.io/symbol/sym-1f618/)
- [SYM 1D478](https://aestheticsymbols.io/symbol/sym-1d478/)
- [AQUARIUS ZODIAC WATER BEARER](https://aestheticsymbols.io/symbol/aquarius-zodiac-water-bearer/)
- [SYM 1D444](https://aestheticsymbols.io/symbol/sym-1d444/)
- [SKULL AND CROSSBONES](https://aestheticsymbols.io/symbol/skull-and-crossbones/)
- [SYM 1F609](https://aestheticsymbols.io/symbol/sym-1f609/)
- [AESTHETIC MINIMAL CLOUD](https://aestheticsymbols.io/symbol/aesthetic-minimal-cloud/)
- [SYM 26EF](https://aestheticsymbols.io/symbol/sym-26ef/)
- [SYM 1D495](https://aestheticsymbols.io/symbol/sym-1d495/)
- [SYM 26F9](https://aestheticsymbols.io/symbol/sym-26f9/)
- [SYM 1D41C](https://aestheticsymbols.io/symbol/sym-1d41c/)
- [SYM 1F605](https://aestheticsymbols.io/symbol/sym-1f605/)
- [SYM 1F92C](https://aestheticsymbols.io/symbol/sym-1f92c/)
- [SYM 1D409](https://aestheticsymbols.io/symbol/sym-1d409/)
- [LEFT HEAVY BRACKET BOX](https://aestheticsymbols.io/symbol/left-heavy-bracket-box/)
- [SYM 1D472](https://aestheticsymbols.io/symbol/sym-1d472/)
- [LITTLE CAT PAWS KAOMOJI](https://aestheticsymbols.io/symbol/little-cat-paws-kaomoji/)
- [ROYAL GOLD CROWN](https://aestheticsymbols.io/symbol/royal-gold-crown/)
- [SYM 1F61F](https://aestheticsymbols.io/symbol/sym-1f61f/)
- [MUSIC SHARP SIGN](https://aestheticsymbols.io/symbol/music-sharp-sign/)
- [SYM 26D4](https://aestheticsymbols.io/symbol/sym-26d4/)
- [SYM 2657](https://aestheticsymbols.io/symbol/sym-2657/)
- [SYM 1D420](https://aestheticsymbols.io/symbol/sym-1d420/)
- [INSTAGRAM BIO](https://aestheticsymbols.io/vi/instagram-bio/)
- [SYM 1D4A3](https://aestheticsymbols.io/symbol/sym-1d4a3/)
- [SYM 1F978](https://aestheticsymbols.io/symbol/sym-1f978/)
- [ARROWS LINES](https://aestheticsymbols.io/pt/arrows-lines/)
- [SYM 2680](https://aestheticsymbols.io/symbol/sym-2680/)
- [SYM 1F49A](https://aestheticsymbols.io/symbol/sym-1f49a/)
- [SYM 2663](https://aestheticsymbols.io/symbol/sym-2663/)
- [SYM 1D44A](https://aestheticsymbols.io/symbol/sym-1d44a/)
- [SYM 2636](https://aestheticsymbols.io/symbol/sym-2636/)
- [SYM 2721](https://aestheticsymbols.io/symbol/sym-2721/)
- [SYM 1F62D](https://aestheticsymbols.io/symbol/sym-1f62d/)
- [SYM 1D452](https://aestheticsymbols.io/symbol/sym-1d452/)
- [SYM 1D453](https://aestheticsymbols.io/symbol/sym-1d453/)
- [TAURUS ZODIAC BULL](https://aestheticsymbols.io/symbol/taurus-zodiac-bull/)
- [AESTHETICSYMBOLS.IO](https://aestheticsymbols.io/)
- [SYM 1D47D](https://aestheticsymbols.io/symbol/sym-1d47d/)
- [SYM 2639 FE0F](https://aestheticsymbols.io/symbol/sym-2639-fe0f/)
- [SYM 26FF](https://aestheticsymbols.io/symbol/sym-26ff/)
- [SYM 2683](https://aestheticsymbols.io/symbol/sym-2683/)
- [SYM 2645](https://aestheticsymbols.io/symbol/sym-2645/)
- [SYM 1D432](https://aestheticsymbols.io/symbol/sym-1d432/)
- [SYM 2633](https://aestheticsymbols.io/symbol/sym-2633/)
- [SYM 2638](https://aestheticsymbols.io/symbol/sym-2638/)
- [FIRST QUARTER WAXING MOON](https://aestheticsymbols.io/symbol/first-quarter-waxing-moon/)
- [SYM 2662](https://aestheticsymbols.io/symbol/sym-2662/)
- [SYM 1F643](https://aestheticsymbols.io/symbol/sym-1f643/)
- [SYM 1D467](https://aestheticsymbols.io/symbol/sym-1d467/)
- [SYM 2659](https://aestheticsymbols.io/symbol/sym-2659/)
- [RIGHTWARDS PAIRED HARPOON](https://aestheticsymbols.io/symbol/rightwards-paired-harpoon/)
- [SYM 1F640](https://aestheticsymbols.io/symbol/sym-1f640/)
- [GOTHIC OBSIDIAN SKULL CREST](https://aestheticsymbols.io/symbol/gothic-obsidian-skull-crest/)
- [SYM 26BA](https://aestheticsymbols.io/symbol/sym-26ba/)
- [SYM 1D401](https://aestheticsymbols.io/symbol/sym-1d401/)
- [SYM 1F479](https://aestheticsymbols.io/symbol/sym-1f479/)
- [TRENDING](https://aestheticsymbols.io/trending/)
- [SYM 268F](https://aestheticsymbols.io/symbol/sym-268f/)
- [SYM 267C](https://aestheticsymbols.io/symbol/sym-267c/)
- [RIGHT POINTING DOUBLE ANGLE QUOTATION](https://aestheticsymbols.io/symbol/right-pointing-double-angle-quotation/)
- [SYM 1F921](https://aestheticsymbols.io/symbol/sym-1f921/)
- [SYM 265C](https://aestheticsymbols.io/symbol/sym-265c/)
- [SYM 1D400](https://aestheticsymbols.io/symbol/sym-1d400/)
- [SYM 26C9](https://aestheticsymbols.io/symbol/sym-26c9/)
- [SYM 1F49B](https://aestheticsymbols.io/symbol/sym-1f49b/)
- [SYM 1D45A](https://aestheticsymbols.io/symbol/sym-1d45a/)
- [SYM 26FB](https://aestheticsymbols.io/symbol/sym-26fb/)
- [SYM 2658](https://aestheticsymbols.io/symbol/sym-2658/)
- [FLOWER GIRL SMILE KAOMOJI](https://aestheticsymbols.io/symbol/flower-girl-smile-kaomoji/)
- [SYM 1D485](https://aestheticsymbols.io/symbol/sym-1d485/)
- [SYM 1F497](https://aestheticsymbols.io/symbol/sym-1f497/)
- [SYM 26F7](https://aestheticsymbols.io/symbol/sym-26f7/)
- [SYM 1F63F](https://aestheticsymbols.io/symbol/sym-1f63f/)
- [TIKTOK CAPTIONS](https://aestheticsymbols.io/tiktok-captions/)
- [SYM 1D40D](https://aestheticsymbols.io/symbol/sym-1d40d/)
- [CLOUD WEATHER SYMBOL](https://aestheticsymbols.io/symbol/cloud-weather-symbol/)
- [SYM 1D45B](https://aestheticsymbols.io/symbol/sym-1d45b/)
- [SYM 1F92D](https://aestheticsymbols.io/symbol/sym-1f92d/)
- [SYM 2642](https://aestheticsymbols.io/symbol/sym-2642/)
- [ZODIAC CELESTIAL](https://aestheticsymbols.io/pt/zodiac-celestial/)
- [SYM 1F636](https://aestheticsymbols.io/symbol/sym-1f636/)
- [SYM 2640](https://aestheticsymbols.io/symbol/sym-2640/)
- [NATURE FLOWERS](https://aestheticsymbols.io/ru/nature-flowers/)
- [STARS](https://aestheticsymbols.io/es/stars/)
- [ZODIAC CELESTIAL](https://aestheticsymbols.io/ru/zodiac-celestial/)
- [SYM 1F60D](https://aestheticsymbols.io/symbol/sym-1f60d/)
- [SYM 2644](https://aestheticsymbols.io/symbol/sym-2644/)
- [SYM 26E8](https://aestheticsymbols.io/symbol/sym-26e8/)
- [CANCER ZODIAC CRAB](https://aestheticsymbols.io/symbol/cancer-zodiac-crab/)
- [SYM 2654](https://aestheticsymbols.io/symbol/sym-2654/)
- [SYM 1D480](https://aestheticsymbols.io/symbol/sym-1d480/)
- [SYM 2725](https://aestheticsymbols.io/symbol/sym-2725/)
- [SYM 2647](https://aestheticsymbols.io/symbol/sym-2647/)
- [SYM 265B](https://aestheticsymbols.io/symbol/sym-265b/)
- [MUSIC WEATHER](https://aestheticsymbols.io/es/music-weather/)
- [ANTICLOCKWISE OPEN CIRCLE ARROW](https://aestheticsymbols.io/symbol/anticlockwise-open-circle-arrow/)
- [SYM 1F47E](https://aestheticsymbols.io/symbol/sym-1f47e/)
- [RADIOACTIVE SYMBOL](https://aestheticsymbols.io/symbol/radioactive-symbol/)
- [ES](https://aestheticsymbols.io/es/)
- [SYM 26DC](https://aestheticsymbols.io/symbol/sym-26dc/)
- [LEFT MATHEMATICAL WHITE SQUARE BRACKET](https://aestheticsymbols.io/symbol/left-mathematical-white-square-bracket/)
- [SYM 2656](https://aestheticsymbols.io/symbol/sym-2656/)
- [STARRY LOVE AURA](https://aestheticsymbols.io/symbol/starry-love-aura/)
- [SYM 262B](https://aestheticsymbols.io/symbol/sym-262b/)
- [SYM 268C](https://aestheticsymbols.io/symbol/sym-268c/)
- [SYM 1D455](https://aestheticsymbols.io/symbol/sym-1d455/)
- [SYM 265A](https://aestheticsymbols.io/symbol/sym-265a/)
- [SYM 1D476](https://aestheticsymbols.io/symbol/sym-1d476/)
- [SYM 26BB](https://aestheticsymbols.io/symbol/sym-26bb/)
- [SYM 1F604](https://aestheticsymbols.io/symbol/sym-1f604/)
- [SYM 1F914](https://aestheticsymbols.io/symbol/sym-1f914/)
- [ZODIAC CELESTIAL](https://aestheticsymbols.io/es/zodiac-celestial/)
- [SYM 1D466](https://aestheticsymbols.io/symbol/sym-1d466/)
- [SYM 1F92B](https://aestheticsymbols.io/symbol/sym-1f92b/)
- [SYM 1F60F](https://aestheticsymbols.io/symbol/sym-1f60f/)
- [LEFT WING CLAN FLARE](https://aestheticsymbols.io/symbol/left-wing-clan-flare/)
- [SYM 265E](https://aestheticsymbols.io/symbol/sym-265e/)
- [QUARTER MUSICAL NOTE](https://aestheticsymbols.io/symbol/quarter-musical-note/)
- [FLUTTERING BUTTERFLY](https://aestheticsymbols.io/symbol/fluttering-butterfly/)
- [SWIMMING FISH LEFT](https://aestheticsymbols.io/symbol/swimming-fish-left/)
- [SYM 1F60B](https://aestheticsymbols.io/symbol/sym-1f60b/)
- [SYM 1F614](https://aestheticsymbols.io/symbol/sym-1f614/)
- [SYM 1F499](https://aestheticsymbols.io/symbol/sym-1f499/)
- [SYM 1F911](https://aestheticsymbols.io/symbol/sym-1f911/)
- [SWIMMING FISH RIGHT](https://aestheticsymbols.io/symbol/swimming-fish-right/)
- [SYM 26AF](https://aestheticsymbols.io/symbol/sym-26af/)
- [SYM 1F622](https://aestheticsymbols.io/symbol/sym-1f622/)
- [BRACKETS](https://aestheticsymbols.io/brackets/)
- [SYM 26F3](https://aestheticsymbols.io/symbol/sym-26f3/)
