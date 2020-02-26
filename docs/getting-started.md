---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis-reva
geekdocEditPath: edit/master/docs
geekdocFilePath: getting-started.md
---

### Installation

So far we are offering two different variants for the installation. You can choose between [Docker](https://www.docker.com/) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

#### Docker

TBD

#### Binaries

TBD

### Configuration

We provide overall three different variants of configuration. The variant based on environment variables and commandline flags are split up into global values and command-specific values.

The configuration tries to map different configuration options from reva into dedicated services. For now please run `bin/ocis-reva {command} -h` to see the list of available options or have a look at [the flagsets](https://github.com/owncloud/ocis-reva/tree/master/pkg/flagset) and the mapping to a reva config in the corresponding [commands](https://github.com/owncloud/ocis-reva/tree/master/pkg/command).

