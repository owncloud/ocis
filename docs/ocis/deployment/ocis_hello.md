---
title: "oCIS with Hello extension"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_hello.md
---

{{< toc >}}

## Overview

- oCIS running behind Traefik as reverse proxy
- oCIS Hello extension runs beside the main oCIS stack and providing the Hello functionality
- Traefik generating self-signed certificates for local setup or obtaining valid SSL certificates for a server setup

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_hello)

The docker stack consists of 3 containers. One of them is Traefik, a proxy which is terminating SSL and forwards the requests to oCIS in the internal docker network.

The oCIS Hello extension is running in another container and enables you to use its functionality from within ownCloud Web.

## Server Deployment

### Requirements

- Linux server with docker and docker-compose installed
- two domains set up and pointing to your server
  - ocis.\* for serving oCIS
  - traefik.\* for serving the Traefik dashboard

See also [example server setup]({{< ref "preparing_server" >}})

### Install oCIS and Traefik

- Clone oCIS repository

  `git clone https://github.com/owncloud/ocis.git`

- Go to the deployment example

  `cd ocis/deployment/examples/ocis_hello`

- Open the `.env` file in a text editor.

  The file by default looks like this:

  ```bash
  # If you're on a internet facing server please comment out following line.
  # It skips certificate validation for various parts of oCIS and is needed if you use self signed certificates.
  INSECURE=true

  ### Traefik settings ###
  # Serve Traefik dashboard. Defaults to "false".
  TRAEFIK_DASHBOARD=
  # Domain of Traefik, where you can find the dashboard. Defaults to "traefik.owncloud.test"
  TRAEFIK_DOMAIN=
  # Basic authentication for the dashboard. Defaults to user "admin" and password "admin"
  TRAEFIK_BASIC_AUTH_USERS=
  # Email address for obtaining LetsEncrypt certificates, needs only be changed if this is a public facing server
  TRAEFIK_ACME_MAIL=

  ### oCIS settings ###
  # oCIS version. Defaults to "latest"
  OCIS_DOCKER_TAG=
  # Domain of oCIS, where you can find the frontend. Defaults to "ocis.owncloud.test"
  OCIS_DOMAIN=
  # oCIS admin user password. Defaults to "admin".
  ADMIN_PASSWORD=
  # The demo users should not be created on a production instance
  # because their passwords are public. Defaults to "false".
  DEMO_USERS=

  ### oCIS Hello settings ###
  # oCIS Hello version. Defaults to "latest"
  OCIS_HELLO_DOCKER_TAG=
  ```

  You are installing oCIS on a server and Traefik will obtain valid certificates for you so please remove `INSECURE=true` or set it to `false`.

  If you want to use the Traefik dashboard, set TRAEFIK_DASHBOARD to `true` (default is `false` and therefore not active). If you activate it, you must set a domain for the Traefik dashboard in `TRAEFIK_DOMAIN=` e.g. `TRAEFIK_DOMAIN=traefik.owncloud.test`.

  The Traefik dashboard is secured by basic auth. Default credentials are the user `admin` with the password `admin`. To set your own credentials, generate a htpasswd (e.g. by using [an online tool](https://htpasswdgenerator.de/) or a cli tool).

  Traefik will issue certificates with LetsEncrypt and therefore you must set an email address in `TRAEFIK_ACME_MAIL=`.

  By default oCIS will be started in the `latest` version. If you want to start a specific version of oCIS set the version to `OCIS_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated).

  Set your domain for the oCIS frontend in `OCIS_DOMAIN=`, e.g. `OCIS_DOMAIN=ocis.owncloud.test`.

  Set the initial admin user password in `ADMIN_PASSWORD=`, it defaults to `admin`.

  By default the oCIS Hello extension will be started in the `latest` version. If you want to start a specific version of oCIS Hello set the version to `OCIS_HELLO_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis-hello/tags?page=1&ordering=last_updated).

  Now you have configured everything and can save the file.

- Start the docker stack

  `docker-compose up -d`

- You now can visit oCIS and are able to switch to the Hello extension by using the application switcher on the top right corner of ownCloud Web. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.

## Local setup

For a more simple local ocis setup see [Getting started]({{< ref "../getting-started" >}})

This docker stack can also be run locally. One downside is that Traefik can not obtain valid SSL certificates and therefore will create self-signed ones. This means that your browser will show scary warnings. Another downside is that you can not point DNS entries to your localhost. So you have to add static host entries to your computer.

On Linux and macOS you can add them to your `/etc/hosts` file and on Windows to `C:\Windows\System32\Drivers\etc\hosts` file like this:

```
127.0.0.1 ocis.owncloud.test
127.0.0.1 traefik.owncloud.test
```

After that you're ready to start the application stack:

`docker-compose up -d`

Open https://ocis.owncloud.test in your browser and accept the invalid certificate warning. You are now able to switch to the Hello extension by using the application switcher on the top right corner of ownCloud Web. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.
