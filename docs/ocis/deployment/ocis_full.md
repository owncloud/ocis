---
title: "Full modular oCIS with WebOffice"
date: 2024-06-25T00:00:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_full.md
---

{{< toc >}}

## Overview

* oCIS, the collaboration service, Collabora or OnlyOffice running behind Traefik as reverse proxy
* Collabora or OnlyOffice enable you to edit office documents in your browser
* The collaboration server acts as a bridge to make the oCIS storage accessible to Collabora and OnlyOffice
* Traefik generating self-signed certificates for local setup or obtaining valid SSL certificates for a server setup
* The whole deployment acts as a modular toolkit to use different flavors of office suites and ocis features

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_full)

## Easy Default

The Infinite Scale Team and product management are providing a default setup for oCIS.

### Goal:
  - provide a good starting point for a production deployment
  - minimal effort to get started with an opinionated setup
  - keep it adjustable it to your needs.

### Default Components

- Infinite Scale
- Full Text Search
- Collabora Online Web Office
- Prepared for LetsEncrypt SSL certificates via Traefik Reverse Proxy

### Optional Components

- ClamAV Virusscanner
- Cloud Importer (Experimental)
- OnlyOffice as an alternative to Collabora
- S3 Storage config to connect to an S3 storage backend
- S3 Minio Server as a local S3 storage backend for debugging and development

### Important Note

If you deviate from the configuration setup and let the `collaboration` service run in its own container, you MUST
ensure the ocis configuration is shared as shown in the example deployment. This is because secrets generated
must be accessible for all services.

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

  `git clone https://github.com/owncloud/ocis.git --depth 1`

* Go to the deployment example

  `cd ocis/deployments/examples/ocis_full`

* Open the `.env` file in a text editor.

  The file by default looks like this:

  ```shell {linenos=table,hl_lines=[8,24,48,50,135,138]}
  ## Basic Settings ##
  # Define the docker compose log driver used.
  # Defaults to local
  LOG_DRIVER=
  # If you're on an internet facing server, comment out following line.
  # It skips certificate validation for various parts of Infinite Scale and is
  # needed when self signed certificates are used.
  INSECURE=true


  ## Traefik Settings ##
  # Note: Traefik is always enabled and can't be disabled.
  # Serve Traefik dashboard.
  # Defaults to "false".
  TRAEFIK_DASHBOARD=
  # Domain of Traefik, where you can find the dashboard.
  # Defaults to "traefik.owncloud.test"
  TRAEFIK_DOMAIN=
  # Basic authentication for the traefik dashboard.
  # Defaults to user "admin" and password "admin" (written as: "admin:admin").
  TRAEFIK_BASIC_AUTH_USERS=
  # Email address for obtaining LetsEncrypt certificates.
  # Needs only be changed if this is a public facing server.
  TRAEFIK_ACME_MAIL=
  # Set to the following for testing to check the certificate process:
  # "https://acme-staging-v02.api.letsencrypt.org/directory"
  # With staging configured, there will be an SSL error in the browser.
  # When certificates are displayed and are emitted by # "Fake LE Intermediate X1",
  # the process went well and the envvar can be reset to empty to get valid certificates.
  TRAEFIK_ACME_CASERVER=


  ## Infinite Scale Settings ##
  # Beside Traefik, this service must stay enabled.
  # Disable only for testing purposes.
  # Note: the leading colon is required to enable the service.
  OCIS=:ocis.yml
  # The oCIS container image.
  # For production releases: "owncloud/ocis"
  # For rolling releases:    "owncloud/ocis-rolling"
  # Defaults to production if not set otherwise
  OCIS_DOCKER_IMAGE=owncloud/ocis-rolling
  # The oCIS container version.
  # Defaults to "latest" and points to the latest stable tag.
  OCIS_DOCKER_TAG=
  # Domain of oCIS, where you can find the frontend.
  # Defaults to "ocis.owncloud.test"
  OCIS_DOMAIN=
  # oCIS admin user password. Defaults to "admin".
  ADMIN_PASSWORD=
  # Demo users should not be created on a production instance,
  # because their passwords are public. Defaults to "false".
  # Also see: https://doc.owncloud.com/ocis/latest/deployment/general/general-info.html#demo-users-and-groups
  DEMO_USERS=
  # Define the oCIS loglevel used.
  # For more details see:
  # https://doc.owncloud.com/ocis/latest/deployment/services/env-vars-special-scope.html
  LOG_LEVEL=
  # Define the kind of logging.
  # The default log can be read by machines.
  # Set this to true to make the log human readable.
  # LOG_PRETTY=true
  #
  # Define the oCIS storage location. Set the paths for config and data to a local path.
  # Note that especially the data directory can grow big.
  # Leaving it default stores data in docker internal volumes.
  # For more details see:
  # https://doc.owncloud.com/ocis/next/deployment/general/general-info.html#default-paths
  # OCIS_CONFIG_DIR=/your/local/ocis/config
  # OCIS_DATA_DIR=/your/local/ocis/data

  # S3 Storage configuration - optional
  # Infinite Scale supports S3 storage as primary storage.
  # Per default, S3 storage is disabled and the local filesystem is used.
  # To enable S3 storage, uncomment the following line and configure the S3 storage.
  # For more details see:
  # https://doc.owncloud.com/ocis/next/deployment/storage/s3.html
  # Note: the leading colon is required to enable the service.
  #S3NG=:s3ng.yml
  # Configure the S3 storage endpoint. Defaults to "http://minio:9000" for testing purposes.
  S3NG_ENDPOINT=
  # S3 region. Defaults to "default".
  S3NG_REGION=
  # S3 access key. Defaults to "ocis"
  S3NG_ACCESS_KEY=
  # S3 secret. Defaults to "ocis-secret-key"
  S3NG_SECRET_KEY=
  # S3 bucket. Defaults to "ocis"
  S3NG_BUCKET=
  #
  # For testing purposes, add local minio S3 storage to the docker-compose file.
  # The leading colon is required to enable the service.
  #S3NG_MINIO=:minio.yml
  # Minio domain. Defaults to "minio.owncloud.test".
  MINIO_DOMAIN=

  # Define SMPT settings if you would like to send Infinite Scale email notifications.
  # For more details see:
  # https://doc.owncloud.com/ocis/latest/deployment/services/s-list/notifications.html
  # NOTE: when configuring mail server, these settings have no effect, see mailserver.yml for details.
  # SMTP host to connect to.
  SMTP_HOST=
  # Port of the SMTP host to connect to.
  SMTP_PORT=
  # An eMail address that is used for sending Infinite Scale notification eMails
  # like "ocis notifications <noreply@yourdomain.com>".
  SMTP_SENDER=
  # Username for the SMTP host to connect to.
  SMTP_USERNAME=
  # Password for the SMTP host to connect to.
  SMTP_PASSWORD=
  # Authentication method for the SMTP communication.
  SMTP_AUTHENTICATION=
  # Allow insecure connections to the SMTP server. Defaults to false.
  SMTP_INSECURE=


  ## Default Enabled Services ##

  ### Apache Tika Content Analysis Toolkit ###
  # Tika (search) is enabled by default, comment if not required.
  # Note: the leading colon is required to enable the service.
  TIKA=:tika.yml
  # Set the desired docker image tag or digest.
  # Defaults to "latest"
  TIKA_IMAGE=


  ### Collabora Settings ###
  # Collabora web office is default enabled, comment if not required.
  # Note: the leading colon is required to enable the service.
  COLLABORA=:collabora.yml
  # Domain of Collabora, where you can find the frontend.
  # Defaults to "collabora.owncloud.test"
  COLLABORA_DOMAIN=
  # Domain of the wopiserver which handles OnlyOffice.
  # Defaults to "wopiserver.owncloud.test"
  WOPISERVER_DOMAIN=
  # Admin user for Collabora.
  # Defaults to "admin".
  # Collabora Admin Panel URL:
  # https://{COLLABORA_DOMAIN}/browser/dist/admin/admin.html
  COLLABORA_ADMIN_USER=
  # Admin password for Collabora.
  # Defaults to "admin".
  COLLABORA_ADMIN_PASSWORD=
  # Set to true to enable SSL for Collabora Online. Default is true if not specified.
  COLLABORA_SSL_ENABLE=false
  # If you're on an internet-facing server, enable SSL verification for Collabora Online.
  # Please comment out the following line:
  COLLABORA_SSL_VERIFICATION=false
  ...
  ```
  #### Reverse Proxy and SSL

  {{< hint type=important >}}
  **Domains and SSL**\
  Though it may sound strange, most of the setups are failing due to a misconfiguration regarding domains and SSL. Please make sure that you have set up the domains correctly and that they are pointing to your server. Also, make sure that you have set up the email address for the LetsEncrypt certificates in `TRAEFIK_ACME_MAIL=`.
  {{< /hint >}}

  You are installing oCIS on a server and Traefik will obtain valid certificates for you so please remove `INSECURE=true` or set it to `false`.

  Traefik will issue certificates with LetsEncrypt and therefore you must set an email address in `TRAEFIK_ACME_MAIL=`.

  #### Infinite Scale Release and Version
  By default oCIS will be started in the `latest` rolling version. Please note that this deployment does currently not work with the 5.x productions releases.
  The oCIS "collaboration" service, which is required by this deployment, is not part of the 5.x releases.

  If you want to use a specific version of oCIS, set the version to a dedicated tag like `OCIS_DOCKER_TAG=6.3.0`. The minimal required oCIS Version to run this deployment is 6.3.0. Available  production versions can be found on [Docker Hub Production](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated) and available rolling releases can be found on [Docker Hub Rolling](https://hub.docker.com/r/owncloud/ocis-rolling/tags?page=1&ordering=last_updated)

  {{< hint type=info title="oCIS Releases" >}}
  You can read more about the different oCIS releases in the [oCIS Release Lifecycle](../release_roadmap.md).
  {{< /hint >}}

  Set your domain for the oCIS frontend in `OCIS_DOMAIN=`, e.g. `OCIS_DOMAIN=ocis.owncloud.test`.

  Set the initial admin user password in `ADMIN_PASSWORD=`, it defaults to `admin`.

  Web Office needs a public domain for the WOPI server to be set in `WOPISERVER_DOMAIN=`, where the office suite can work on the files via the WOPI protocol.

  Now it's time to set up Collabora and you need to configure the domain of Collabora in `COLLABORA_DOMAIN=`.

  If you want to use the Collabora admin panel you need to set the username and password for the administrator in `COLLABORA_ADMIN_USER=` and `COLLABORA_ADMIN_PASSWORD=`.

* Start the docker stack

  `docker-compose up -d`

* You now can visit oCIS and are able to open an office document in your browser. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.

## Local Setup

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
127.0.0.1 minio.owncloud.test
```

After that, you're ready to start the application stack:

`docker-compose pull && docker-compose up -d`

Open https://collabora.owncloud.test in your browser and accept the invalid certificate warning.

Open https://ocis.owncloud.test in your browser and accept the invalid certificate warning. You are now able to open an office document in your browser. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.

## Additional services

### Clamav Virusscanner

You can add a Clamav Virusscanner to the stack. The service is disabled by default. To enable it, uncomment the `CLAMAV` line in the `.env` file.

```shell {linenos=table,hl_lines=[3]}
## Clamav Settings ##
# The leading colon is required to enable the service.
CLAMAV=:clamav.yml
```

After enabling that service, you can add the service to the stack with `docker-compose up -d` again.

### Traefik Dashboard

If you want to use the Traefik dashboard, set TRAEFIK_DASHBOARD to `true` (default is `false` and therefore not active). If you activate it, you must set a domain for the Traefik dashboard in `TRAEFIK_DOMAIN=` e.g. `TRAEFIK_DOMAIN=traefik.owncloud.test`.

The Traefik dashboard is secured by basic auth. Default credentials are the user `admin` with the password `admin`. To set your own credentials, generate a htpasswd (e.g. by using [an online tool](https://htpasswdgenerator.de/) or a cli tool).

```shell {linenos=table,hl_lines=[4,7,10]}
### Traefik Settings ###
# Serve Traefik dashboard.
# Defaults to "false".
TRAEFIK_DASHBOARD=true
# Domain of Traefik, where you can find the dashboard.
# Defaults to "traefik.owncloud.test"
TRAEFIK_DOMAIN=
# Basic authentication for the traefik dashboard.
# Defaults to user "admin" and password "admin" (written as: "admin:admin").
TRAEFIK_BASIC_AUTH_USERS=
```
### Cloud Importer

Cloud importer can provide an Upload Interface to your oCIS instance. It is a separate service that can be enabled in the `.env` file.

```shell {linenos=table,hl_lines=[3]}
## Uppy Companion Settings ##
# The leading colon is required to enable the service.
CLOUD_IMPORTER=:cloudimporter.yml
## The docker image to be used for uppy companion.
# owncloud has built a container with public link import support.
COMPANION_IMAGE=
# Domain of Uppy Companion. Defaults to "companion.owncloud.test".
COMPANION_DOMAIN=
# Provider settings, see https://uppy.io/docs/companion/#provideroptions for reference.
# Empty by default, which disables providers.
COMPANION_ONEDRIVE_KEY=
COMPANION_ONEDRIVE_SECRET=
```

After Enabling that servive by uncommenting the `CLOUD_IMPORTER` line, you can add the service to the stack with `docker-compose up -d` again.

### S3 Storage

You can use an S3 compatible Storage as the primary data store. The metadatata of the files will still be stored on the local filesystem.

{{<hint type="info">}}
The endpoint, region and keys for your S3 Server need to be provided by the service or company who operates it. Normally you can get these via web portal.
{{</hint>}}

```shell {linenos=table,hl_lines=[8,10,12,14,16,18]}
# S3 Storage configuration - optional
# Infinite Scale supports S3 storage as primary storage.
# Per default, S3 storage is disabled and the local filesystem is used.
# To enable S3 storage, uncomment the following line and configure the S3 storage.
# For more details see:
# https://doc.owncloud.com/ocis/next/deployment/storage/s3.html
# Note: the leading colon is required to enable the service.
# S3NG=:s3ng.yml
# Configure the S3 storage endpoint. Defaults to "http://minio:9000" for testing purposes.
S3NG_ENDPOINT=
# S3 region. Defaults to "default".
S3NG_REGION=
# S3 access key. Defaults to "ocis"
S3NG_ACCESS_KEY=
# S3 secret. Defaults to "ocis-secret-key"
S3NG_SECRET_KEY=
# S3 bucket. Defaults to "ocis"
S3NG_BUCKET=
```

#### Use a Local Minio S3 Storage Backend

For testing purposes, you can use a local minio S3 storage backend. To enable it, uncomment the `S3NG_MINIO` line in the `.env` file.

The frontend for the minio server is available at `http://minio.owncloud.test` and the access key is `ocis` and the secret key is `ocis-secret`.

## Local Setup for Web Development

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
