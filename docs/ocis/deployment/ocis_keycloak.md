---
title: "oCIS with Keycloak"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_keycloak.md
---

{{< toc >}}

## Overview

* oCIS and Keycloak running behind Traefik as reverse proxy
* Keycloak acting as the IDP for oCIS
* Traefik generating self signed certificates for local setup or obtaining valid SSL certificates for a server setup

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_keycloak)

The docker stack consists 4 containers. One of them is Traefik, a proxy which is terminating ssl and forwards the requests to oCIS in the internal docker network.

Keykloak add two containers: Keycloak itself and a PostgreSQL as database. Keycloak will be configured as oCIS' IDP instead of the internal IDP [LibreGraph Connect]({{< ref "../../extensions/idp" >}})

The other container is oCIS itself running all extensions in one container. In this example oCIS uses [oCIS storage driver]({{< ref "../../extensions/storage/storages#storage-drivers" >}})

## Server Deployment

### Requirements

* Linux server with docker and docker-compose installed
* Three domains set up and pointing to your server
  - ocis.* for serving oCIS
  - keycloak.* for serving Keycloak
  - traefik.* for serving the Traefik dashboard

See also [example server setup]({{< ref "preparing_server" >}})


### Install oCIS and Traefik

* Clone oCIS repository

  `git clone https://github.com/owncloud/ocis.git`

* Go to the deployment example

  `cd ocis/deployment/examples/ocis_keycloak`

* Open the `.env` file in a text editor
  The file by default looks like this:
  ```bash
  # If you're on a internet facing server please comment out following line.
  # It skips certificate validation for various parts of oCIS and is needed if you use self signed certificates.
  INSECURE=true

  ### Traefik settings ###
  # Serve Treafik dashboard. Defaults to "false".
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
  # ownCloud Web openid connect client id. Defaults to "ocis-web"
  OCIS_OIDC_CLIENT_ID=
  # IDP LDAP bind password. Must be changed in order to have a secure oCIS. Defaults to "idp".
  IDP_LDAP_BIND_PASSWORD=
  # Storage LDAP bind password. Must be changed in order to have a secure oCIS. Defaults to "reva".
  STORAGE_LDAP_BIND_PASSWORD=
  # JWT secret which is used for the storage provider. Must be changed in order to have a secure oCIS. Defaults to "Pive-Fumkiu4"
  OCIS_JWT_SECRET=
  # JWT secret which is used for uploads to create transfer tokens. Must be changed in order to have a secure oCIS. Defaults to "replace-me-with-a-transfer-secret"
  OCIS_TRANSFER_SECRET=

  ### Keycloak ###
  # Domain of Keycloak, where you can find the management and authentication frontend. Defaults to "keycloak.owncloud.test"
  KEYCLOAK_DOMAIN=
  # Realm which to be used with oCIS. Defaults to "oCIS"
  KEYCLOAK_REALM=
  # Admin user login name. Defaults to "admin"
  KEYCLOAK_ADMIN_USER=
  # Admin user login password. Defaults to "admin"
  KEYCLOAK_ADMIN_PASSWORD=

  ```

  You are installing oCIS on a server and Traefik will obtain valid certificates for you so please remove `INSECURE=true` or set it to `false`.

  If you want to use the Traefik dashboard, set TRAEFIK_DASHBOARD to `true` (default is `false` and therefore not active). If you activate it, you must set a domain for the Traefik dashboard in `TRAEFIK_DOMAIN=` eg. `TRAEFIK_DOMAIN=traefik.owncloud.test`.

  The Traefik dashboard is secured by basic auth. Default credentials are the user `admin` with the password `admin`. To set your own credentials, generate a htpasswd (eg. by using [an online tool](https://htpasswdgenerator.de/) or a cli tool).

  Traefik will issue certificates with LetsEncrypt and therefore you must set an email address in `TRAEFIK_ACME_MAIL=`.

  By default oCIS will be started in the `latest` version. If you want to start a specific version of oCIS set the version to `OCIS_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated).

  Set your domain for the oCIS frontend in `OCIS_DOMAIN=`, eg. `OCIS_DOMAIN=ocis.owncloud.test`.

  If you want to change the OIDC client id of th ownCloud Web frontend, you can do this by setting the name to `OCIS_OIDC_CLIENT_ID=`.

  You also must override three default secrets in `IDP_LDAP_BIND_PASSWORD`, `STORAGE_LDAP_BIND_PASSWORD` and `OCIS_JWT_SECRET` in order to secure your oCIS instance. Choose some random strings eg. from the output of `openssl rand -base64 32`. For more information see [secure an oCIS instance]({{< ref "./#secure-an-ocis-instance" >}}).

  Set your domain for the Keycloak administration panel and authentication endpoints to `KEYCLOAK_DOMAIN=` eg. `KEYCLOAK_DOMAIN=keycloak.owncloud.test`.

  Changing the used Keycloak realm can be done by setting `KEYCLOAK_REALM=`. This defaults to the oCIS realm `KEYCLOAK_REALM=oCIS`. The oCIS realm will be automatically imported on startup and includes our demo users.

  You probably should secure your Keycloak admin account by setting `KEYCLOAK_ADMIN_USER=` and `KEYCLOAK_ADMIN_PASSWORD=` to values other than `admin`.

  Now you have configured everything and can save the file.

* Start the docker stack

  `docker-compose up -d`

* You now can visit oCIS, Keycloak and Traefik dashboard on your configured domains

## Local setup
For a more simple local ocis setup see [Getting started]({{< ref "../getting-started" >}})

This docker stack can also be run locally. One downside is that Traefik can not obtain valid SSL certificates and therefore will create self signed ones. This means that your browser will show scary warnings. Another downside is that you can not point DNS entries to your localhost. So you have to add static host entries to your computer.

On Linux and macOS you can add them to your `/etc/hosts` files like this:
```
127.0.0.1 ocis.owncloud.test
127.0.0.1 traefik.owncloud.test
127.0.0.1 keycloak.owncloud.test
```

After that you're ready to start the application stack:

`docker-compose up -d`

Open https://keycloak.owncloud.test in your browser and accept the invalid certificate warning.

Open https://ocis.owncloud.test in your browser and accept the invalid certificate warning. You now can login to oCIS with the demo users.
