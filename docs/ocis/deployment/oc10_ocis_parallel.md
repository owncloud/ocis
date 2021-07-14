---
title: "Parallel deployment of oC10 and oCIS"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: oc10_ocis_parallel.md
---

{{< toc >}}

## Overview

- This setup reflects [stage 6 of the oC10 to oCIS migration plan]({{< ref "migration#stage-6-parallel-deployment" >}})
- Traefik generating self signed certificates for local setup or obtaining valid SSL certificates for a server setup
- OpenLDAP server with demo users
- LDAP admin interface to edit users
- Keycloak as OpenID Connect provider in federation with the LDAP server
- ownCloud 10 with MariaDB and Redis
  - ownCloud 10 is configured to synchronize users from the LDAP server
  - ownCloud 10 is used to use OpenID Connect for authentication with Keycloak
- oCIS running behind Traefik as reverse proxy
  - oCIS is using the ownCloud storage driver on the same files and same database as ownCloud 10
  - oCIS is using Keycloak as OpenID Connect provider
  - oCIS is using the LDAP server as user backend
- All requests to both oCIS and oC10 are routed through the oCIS proxy and will be routed based on an OIDC claim to one of them. Therefore admins can change on a user basis in the LDAP which backend is used.

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/oc10_ocis_parallel)

## Server Deployment

### Requirements

- Linux server with docker and docker-compose installed
- four domains set up and pointing to your server
  - cloud.\* for serving oCIS
  - keycloak.\* for serving Keycloak
  - ldap .\* for serving the LDAP managment UI
  - traefik.\* for serving the Traefik dashboard

See also [example server setup]({{< ref "preparing_server" >}})

### Install this example

- Clone oCIS repository

  `git clone https://github.com/owncloud/ocis.git`

- Go to the deployment example

  `cd ocis/deployment/examples/oc10_ocis_parallel`

- Open the `.env` file in a text editor
  The file by default looks like this:

  ```bash
      # If you're on a internet facing server please comment out following line.
      # It skips certificate validation for various parts of oCIS and is needed if you use self signed certificates.
      INSECURE=true

      ### Traefik settings ###
      TRAEFIK_LOG_LEVEL=
      # Serve Treafik dashboard. Defaults to "false".
      TRAEFIK_DASHBOARD=
      # Domain of Traefik, where you can find the dashboard. Defaults to "traefik.owncloud.test"
      TRAEFIK_DOMAIN=
      # Basic authentication for the dashboard. Defaults to user "admin" and password "admin"
      TRAEFIK_BASIC_AUTH_USERS=
      # Email address for obtaining LetsEncrypt certificates, needs only be changed if this is a public facing server
      TRAEFIK_ACME_MAIL=

      ### shared oCIS / oC10 settings ###
      # Domain of oCIS / oC10, where you can find the frontend. Defaults to "cloud.owncloud.test"
      CLOUD_DOMAIN=

      ### oCIS settings ###
      # oCIS version. Defaults to "latest"
      OCIS_DOCKER_TAG=
      # JWT secret which is used for the storage provider. Must be changed in order to have a secure oCIS. Defaults to "Pive-Fumkiu4"
      OCIS_JWT_SECRET=
      # JWT secret which is used for uploads to create transfer tokens. Must be changed in order to have a secure oCIS. Defaults to "replace-me-with-a-transfer-secret"
      STORAGE_TRANSFER_SECRET=

      ### oCIS settings ###
      # oC10 version. Defaults to "latest"
      OC10_DOCKER_TAG=
      # client secret which the openidconnect app uses to authenticate to Keycloak. Defaults to "oc10-oidc-secret"
      OC10_OIDC_CLIENT_SECRET=
      # app which will be shown when opening the ownCloud 10 UI. Defaults to "files" but also could be set to "web"
      OWNCLOUD_DEFAULT_APP=
      # if set to "false" (default) links will be opened in the classic UI, if set to "true" ownCloud Web is used
      OWNCLOUD_WEB_REWRITE_LINKS=

      ### LDAP settings ###
      # password for the LDAP admin user "cn=admin,dc=owncloud,dc=com", defaults to "admin"
      LDAP_ADMIN_PASSWORD=
      # Domain of the LDAP management frontend. Defaults to "ldap.owncloud.test"
      LDAP_MANAGER_DOMAIN=

      ### Keycloak ###
      # Domain of Keycloak, where you can find the managment and authentication frontend. Defaults to "keycloak.owncloud.test"
      KEYCLOAK_DOMAIN=
      # Realm which to be used with oC10 and oCIS. Defaults to "owncloud"
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

  Set your domain for the oC10 and oCIS frontend in `CLOUD_DOMAIN=`, eg. `CLOUD_DOMAIN=cloud.owncloud.test`.

  You also must override the default secrets in `STORAGE_TRANSFER_SECRET` and `OCIS_JWT_SECRET` in order to secure your oCIS instance. Choose some random strings eg. from the output of `openssl rand -base64 32`. For more information see [secure an oCIS instance]({{< ref "./#secure-an-ocis-instance" >}}).

  By default ownCloud 10 will be started in the `latest` version. If you want to start a specific version of oCIS set the version to `OC10_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated).

  You can switch the default application of ownCloud 10 by setting`OWNCLOUD_DEFAULT_APP=files` in oder to have the classic UI as frontend, which is also the default. If you prefer ownCloud Web as the default application in ownCloud 10 just set `OWNCLOUD_DEFAULT_APP=web`.

  In oder to change the default link open action which defaults to the classic UI (`OWNCLOUD_WEB_REWRITE_LINKS=false`) you can set it to `OWNCLOUD_WEB_REWRITE_LINKS=true`. This will lead to links being opened in ownCloud Web.

  The OpenLDAP server in this example deployment has an admin users, which is also used as bind user in order to keep theses examples simple. You can change the default password "admin" to a different one by setting it to `LDAP_ADMIN_PASSWORD=...`.

  Set your domain for the LDAP manager UI in `LDAP_MANAGER_DOMAIN=`, eg. `ldap.owncloud.test`.

  Set your domain for the Keycloak administration panel and authentication endpoints to `KEYCLOAK_DOMAIN=` eg. `KEYCLOAK_DOMAIN=keycloak.owncloud.test`.

  Changing the used Keycloak realm can be done by setting `KEYCLOAK_REALM=`. This defaults to the ownCloud realm `KEYCLOAK_REALM=owncloud`. The ownCloud realm will be automatically imported on startup and includes our demo users.

  You probably should secure your Keycloak admin account by setting `KEYCLOAK_ADMIN_USER=` and `KEYCLOAK_ADMIN_PASSWORD=` to values other than `admin`.

  Now you have configured everything and can save the file.

- Start the docker stack

  `docker-compose up -d`

- You now can visit the cloud, oC10 or oCIS depending on the user configuration. Marie defaults to oC10 and Richard and Einstein default to oCIS, but you can change the ownCloud selector at any time in the LDAP management UI.

## Local setup

For a more simple local ocis setup see [Getting started]({{< ref "../getting-started" >}})

This docker stack can also be run locally. One downside is that Traefik can not obtain valid SSL certificates and therefore will create self signed ones. This means that your browser will show scary warnings. Another downside is that you can not point DNS entries to your localhost. So you have to add static host entries to your computer.

On Linux and macOS you can add them to your `/etc/hosts` files like this:

```
127.0.0.1 cloud.owncloud.test
127.0.0.1 keycloak.owncloud.test
127.0.0.1 ldap.owncloud.test
127.0.0.1 traefik.owncloud.test
```

After that you're ready to start the application stack:

`docker-compose up -d`

You now can visit the cloud, oC10 or oCIS depending on the user configuration. Marie defaults to oC10 and Richard and Einstein default to oCIS, but you can change the ownCloud selector at any time in the LDAP management UI.
