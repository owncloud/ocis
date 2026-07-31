---
title: 'Tooling'
date: 2022-01-28T00:00:00+00:00
weight: 40
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/development
geekdocFilePath: tooling.md
---

{{< toc >}}

## Packaging

Web is using [pnpm](https://pnpm.io/) as package manager and [vite](https://vitejs.dev/) as build tool. The latter is built on top of [rollup](https://rollupjs.org/) and brings some additional features such as instant hot-reloading.

## Development Setup

### Prerequisites

- docker
- docker-compose (if not already included in your docker installation)
- pnpm
- node

If you’re not using Docker Desktop, you might have to modify your `/etc/hosts` and add `127.0.0.1 host.docker.internal` to make `host.docker.internal` links work.

### Installing Dependencies

After cloning the source code, install the dependencies via `pnpm install`.

### Starting the Server

You can start the server by running `docker-compose up ocis`.

Note that the container needs a short while to start because it is waiting for `tika` to be initialized. This is the case as soon as the `tika-service` container has stopped running.

### Building and Accessing Web

After the docker containers are running (and `tika` is being initialized), run `pnpm build:w` to build Web. This also includes hot-reloading after changes you make, although it will take a while to rebuild the project. See down below for some details on how to enable instant hot-reloading.

Now you can access Web via https://host.docker.internal:9200.

### Using Instant Hot-Reload via Vite

To work with instant hot-reloading, you can also build Web by running `pnpm vite`. The port to access Web is slightly different then: https://host.docker.internal:9201. Also note that the initial page load may take a bit longer than usual. This is normal and to be expected.

### Running Web with oC10

Older versions of Web (< 7.1.0) also support running oC10 as server. The development setup is nearly the same as mentioned above, the only differences are:

* The server can be started via `docker-compose up oc10`
* The server port is `8080` (`8081` when running Web via `pnpm vite:oc10`)
