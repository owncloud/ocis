---
title: "Full oCIS with WebOffice"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_full.md
---

{{< toc >}}

## Overview

* oCIS, Wopi server, Collabora or OnlyOffice running behind Traefik as reverse proxy
* Collabora or OnlyOffice enable you to edit documents in your browser
* Wopi server acts as a bridge to make the oCIS storage accessible to Collabora and OnlyOffice
* Traefik generating self-signed certificates for local setup or obtaining valid SSL certificates for a server setup

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_full)

The docker stack consists of 10 containers. One of them is Traefik, a proxy which is terminating SSL and forwards the requests to oCIS in the internal docker network.

The next container is oCIS itself in a configuration like the [oCIS with Traefik example]({{< ref "ocis_traefik" >}}), except that for this example a custom mimetype configuration is used.

There are three oCIS app driver containers that register Collabora and OnlyOffice at the app registry.

The last four containers are the WOPI server, Collabora and OnlyOffice.

## Server Deployment

### Requirements

* Linux server with docker and docker-compose installed
* Three domains set up and pointing to your server
  * ocis.* for serving oCIS
  * collabora.* for serving Collabora
  * onlyoffice.* for serving OnlyOffice
  * wopiserver.* for serving the WOPI server
  * traefik.* for serving the Traefik dashboard
  * companion.* for serving the uppy companion app

See also [example server setup]({{< ref "preparing_server" >}})

### Install oCIS and Traefik

* Clone oCIS repository

  `git clone https://github.com/owncloud/ocis.git`

* Go to the deployment example

  `cd ocis/deployments/examples/ocis_wopi`

* Open the `.env` file in a text editor.

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
  # Basic authentication for the dashboard. Defaults to user "admin" and password "admin" (written as: "admin:admin").
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
  # Log level for oCIS. Defaults to "info".
  OCIS_LOG_LEVEL=

  ### Wopi server settings ###
  # cs3org wopi server version. Defaults to "v8.3.3"
  WOPISERVER_DOCKER_TAG=
  # cs3org wopi server domain. Defaults to "wopiserver.owncloud.test"
  WOPISERVER_DOMAIN=
  # JWT secret which is used for the documents to be request by the Wopi client from the cs3org Wopi server. Must be change in order to have a secure Wopi server. Defaults to "LoremIpsum567"
  WOPI_JWT_SECRET=

  ### Collabora settings ###
  # Domain of Collabora, where you can find the frontend. Defaults to "collabora.owncloud.test"
  COLLABORA_DOMAIN=
  # Admin user for Collabora. Defaults to blank, provide one to enable access. Collabora Admin Panel URL: https://{COLLABORA_DOMAIN}/browser/dist/admin/admin.html
  COLLABORA_ADMIN_USER=
  # Admin password for Collabora. Defaults to blank, provide one to enable access
  COLLABORA_ADMIN_PASSWORD=

  ### OnlyOffice settings ###
  # Domain of OnlyOffice, where you can find the frontend. Defaults to "onlyoffice.owncloud.test"
  ONLYOFFICE_DOMAIN=

  ### Email / Inbucket settings ###
  # Inbucket / Mail domain. Defaults to "mail.owncloud.test"
  INBUCKET_DOMAIN=

  ### Apache Tika Content analysis toolkit ###
  # Set the desired docker image tag or digest, defaults to "latest"
  TIKA_IMAGE=

  # If you want to use debugging and tracing with this stack,
  # you need uncomment following line. Please see documentation at
  # https://owncloud.dev/ocis/deployment/monitoring-tracing/
  #COMPOSE_FILE=docker-compose.yml:monitoring_tracing/docker-compose-additions.yml

  ### Uppy Companion settings ###
  # Domain of Uppy Companion. Defaults to "companion.owncloud.test"
  COMPANION_IMAGE=
  COMPANION_DOMAIN=
  # Provider settings, see https://uppy.io/docs/companion/#provideroptions for reference. Empty by default, which disables providers.
  COMPANION_ONEDRIVE_KEY=
  COMPANION_ONEDRIVE_SECRET=
  ```

  You are installing oCIS on a server and Traefik will obtain valid certificates for you so please remove `INSECURE=true` or set it to `false`.

  If you want to use the Traefik dashboard, set TRAEFIK_DASHBOARD to `true` (default is `false` and therefore not active). If you activate it, you must set a domain for the Traefik dashboard in `TRAEFIK_DOMAIN=` e.g. `TRAEFIK_DOMAIN=traefik.owncloud.test`.

  The Traefik dashboard is secured by basic auth. Default credentials are the user `admin` with the password `admin`. To set your own credentials, generate a htpasswd (e.g. by using [an online tool](https://htpasswdgenerator.de/) or a cli tool).

  Traefik will issue certificates with LetsEncrypt and therefore you must set an email address in `TRAEFIK_ACME_MAIL=`.

  By default oCIS will be started in the `latest` version. If you want to start a specific version of oCIS set the version to `OCIS_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated).

  Set your domain for the oCIS frontend in `OCIS_DOMAIN=`, e.g. `OCIS_DOMAIN=ocis.owncloud.test`.

  Set the initial admin user password in `ADMIN_PASSWORD=`, it defaults to `admin`.

  By default the CS3Org WOPI server will also be started in the `latest` version. If you want to start a specific version of it, you can set the version to `WOPISERVER_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/cs3org/wopiserver/tags?page=1&ordering=last_updated).

  Set your domain for the CS3Org WOPI server in `WOPISERVER_DOMAIN=`, where all office suites can download the files via the WOPI protocol.

  You also must override the default WOPI JWT secret in order to have a secure setup. Do this by setting `WOPI_JWT_SECRET` to a long and random string.

  Now it's time to set up Collabora and you need to configure the domain of Collabora in `COLLABORA_DOMAIN=`.

  If you want to use the Collabora admin panel you need to set the username and password for the administrator in `COLLABORA_ADMIN_USER=` and `COLLABORA_ADMIN_PASSWORD=`.

  Next up is OnlyOffice, which also needs a domain in `ONLYOFFICE_DOMAIN=`.

  Now you have configured everything and can save the file.

* Start the docker stack

  `docker-compose up -d`

* You now can visit oCIS and are able to open an office document in your browser. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.

## Local setup
For a more simple local ocis setup see [Getting started]({{< ref "../getting-started" >}})

This docker stack can also be run locally. One downside is that Traefik can not obtain valid SSL certificates and therefore will create self-signed ones. This means that your browser will show scary warnings. Another downside is that you can not point DNS entries to your localhost. So you have to add static host entries to your computer.

On Linux and macOS you can add them to your `/etc/hosts` file and on Windows to `C:\Windows\System32\Drivers\etc\hosts` file like this:

```
127.0.0.1 ocis.owncloud.test
127.0.0.1 traefik.owncloud.test
127.0.0.1 collabora.owncloud.test
127.0.0.1 onlyoffice.owncloud.test
127.0.0.1 wopiserver.owncloud.test
127.0.0.1 mail.owncloud.test
127.0.0.1 companion.owncloud.test
```

After that you're ready to start the application stack:

`docker-compose up -d`

Open https://collabora.owncloud.test, https://onlyoffice.owncloud.test and https://wopiserver.owncloud.test  in your browser and accept the invalid certificate warning.

Open https://ocis.owncloud.test in your browser and accept the invalid certificate warning. You are now able to open an office document in your browser. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.

## Local setup for web development

In case you want to run ownCloud Web from a development branch together with this deployment example (e.g. for feature development for the app provider frontend) you can use this deployment example with the local setup and some additional steps as described below.

1. Clone the [ownCloud Web repository](https://github.com/owncloud/web) on your development machine.
2. Run `pnpm i && pnpm build:w` for `web`, so that it creates and continuously updates the `dist` folder for web.
3. Add the dist folder as read only volume to `volumes` section of the `ocis` service in the `docker-compose.yml` file:
   ```yaml
   - /your/local/path/to/web/dist/:/web/dist:ro
   ```
   Make sure to point to the `dist` folder inside your local copy of the web repository.
4. Set the oCIS environment variables `WEB_ASSET_CORE_PATH` and `WEB_ASSET_APPS_PATH` in the `environment` section of the `ocis` service, so that it uses your mounted dist folder for the web assets, instead of the assets that are embedded into oCIS.
   ```yaml
   WEB_ASSET_CORE_PATH: "/web/dist"
   WEB_ASSET_APPS_PATH: "/web/dist"
   ```
5. Start the deployment example as described above in the `Local setup` section.

For app provider frontend development in `web` you can find the source code in `web/packages/web-app-external`. Some parts of the integration live in `web/packages/web-app-files`.

## Using Podman

Podman doesn't have a "local" log driver. Also it's docker-compatibility socket does live in a different location, especially when running a rootless podman.

Using the following settings you can run the deployment with a recent podman version:

```bash
LOG_DRIVER=journald \
DOCKER_SOCKET_PATH=/run/user/1000/podman/podman.sock \
podman compose start
```
