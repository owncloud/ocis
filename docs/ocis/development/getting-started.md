---
title: "Getting Started"
date: 2020-07-07T20:35:00+01:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Requirements

We want contribution to oCIS and the creation of extensions to be as easy as possible.
So we are trying to reflect this in the tooling. It should be kept simple and quick to be set up.

Besides standard development tools like git and a text editor, you need the following software for development:

- Go >= v1.15 ([install instructions](https://golang.org/doc/install))
- Yarn ([install instructions](https://classic.yarnpkg.com/en/docs/install))
- docker ([install instructions](https://docs.docker.com/get-docker/))
- docker-compose ([install instructions](https://docs.docker.com/compose/install/))

If you find tools needed besides the mentioned above, please feel free to open an issue or open a PR.

## Repository structure

oCIS consists of multiple micro services, also called extensions. We started by having standalone repositories for each of them, but quickly noticed that this adds a time consuming overhead for developers. So we ended up with a monorepo housing all the extensions in one repository.

Each extension lives in a subfolder (eg. `accounts` or `settings`) within this respository as an independent Go module, following the [golang-standard project-layout](https://github.com/golang-standards/project-layout). They have common Makefile targets and can be used to change, build and run individual extensions.

The `ocis` folder contains our [go-micro](https://github.com/asim/go-micro/) and [suture](https://github.com/thejerf/suture) based runtime. It is used to import all extensions and implements commands to manage them, similar to a small orchestrator. With the resulting oCIS binary you can start single extensions or even all extensions at the same time.

The `docs` folder contains the source for the [oCIS documentation]({{< ref "../" >}}).

The `deployments` folder contains documented deployment configurations and templates. On a single node, running a single ocis runtime is a resource efficient way to deploy ocis. For multiple nodes docker compose or helm charts for kubernetes examples can be found here.

The `scripts` folder contains scripts to perform various build, install, analysis, etc operations.

## Starting points

Depending on what you want to develop there are different starting points. These will be described below.

### Developing oCIS

If you want to contribute to oCIS:

- see [contribution guidelines](https://github.com/owncloud/ocis#contributing)
- make sure the tooling is set up by [building oCIS]({{< ref "build" >}}) and [building the docs]({{< ref "build-docs" >}})
- create or pick an [open issue](https://github.com/owncloud/ocis/issues) to develop on and mention in the issue that you are working on it
- open a PR and get things done

### Developing extensions

If you want to develop an extension, start here: [Extensions]({{< ref "extensions">}})
