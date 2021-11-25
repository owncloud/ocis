---
title: "oCIS with LDAP"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_ldap.md
---


{{< toc >}}

## Overview

- Traefik generating self signed certificates for local setup or obtaining valid SSL certificates for a server setup
- OpenLDAP server with demo users
- LDAP admin interface to edit users
- oCIS running behind Traefik as reverse proxy
  - oCIS is using the LDAP server as user backend

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_ldap)

## Server Deployment

### Requirements

- Linux server with docker and docker-compose installed
- four domains set up and pointing to your server
  - ocis.\* for serving oCIS
  - ldap .\* for serving the LDAP managment UI
  - traefik.\* for serving the Traefik dashboard

See also [example server setup]({{< ref "preparing_server" >}})

### Install this example

- Clone oCIS repository

  `git clone https://github.com/owncloud/ocis.git`

- Go to the deployment example

  `cd ocis/deployment/examples/ocis_ldap`

- Open the `.env` file in a text editor
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
  # JWT secret which is used for the storage provider. Must be changed in order to have a secure oCIS. Defaults to "Pive-Fumkiu4"
  OCIS_JWT_SECRET=
  # JWT secret which is used for uploads to create transfer tokens. Must be changed in order to have a secure oCIS. Defaults to "replace-me-with-a-transfer-secret"
  STORAGE_TRANSFER_SECRET=
  # Machine auth api key secret. Must be changed in order to have a secure oCIS. Defaults to "change-me-please"
  OCIS_MACHINE_AUTH_API_KEY=

  ### LDAP server settings ###
  # Password of LDAP user "cn=admin,dc=owncloud,dc=com". Defaults to "admin"
  LDAP_ADMIN_PASSWORD=

  ### LDAP manager settings ###
  # Domain of LDAP manager. Defaults to "ldap.owncloud.test"
  LDAP_MANAGER_DOMAIN=
  ```

  You are installing oCIS on a server and Traefik will obtain valid certificates for you so please remove `INSECURE=true` or set it to `false`.

  If you want to use the Traefik dashboard, set TRAEFIK_DASHBOARD to `true` (default is `false` and therefore not active). If you activate it, you must set a domain for the Traefik dashboard in `TRAEFIK_DOMAIN=` eg. `TRAEFIK_DOMAIN=traefik.owncloud.test`.

  The Traefik dashboard is secured by basic auth. Default credentials are the user `admin` with the password `admin`. To set your own credentials, generate a htpasswd (eg. by using [an online tool](https://htpasswdgenerator.de/) or a cli tool).

  Traefik will issue certificates with LetsEncrypt and therefore you must set an email address in `TRAEFIK_ACME_MAIL=`.

  By default oCIS will be started in the `latest` version. If you want to start a specific version of oCIS set the version to `OCIS_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated).

  Set your domain for the oCIS frontend in `OCIS_DOMAIN=`, eg. `OCIS_DOMAIN=cloud.owncloud.test`.

  You also must override the default secrets in `STORAGE_TRANSFER_SECRET` and `OCIS_JWT_SECRET` in order to secure your oCIS instance. Choose some random strings eg. from the output of `openssl rand -base64 32`. For more information see [secure an oCIS instance]({{< ref "./#secure-an-ocis-instance" >}}).

  The OpenLDAP server in this example deployment has an admin users, which is also used as bind user in order to keep theses examples simple. You can change the default password "admin" to a different one by setting it to `LDAP_ADMIN_PASSWORD=...`.

  Set your domain for the LDAP manager UI in `LDAP_MANAGER_DOMAIN=`, eg. `ldap.owncloud.test`.

  Now you have configured everything and can save the file.

- Start the docker stack

  `docker-compose up -d`

- You now can visit oCIS and Traefik dashboard on your configured domains. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.

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

Open https://ocis.owncloud.test in your browser and accept the invalid certificate warning. You now can login to oCIS with the default users, which also can be found here: [Getting started]({{< ref "../getting-started#login-to-ocis-web" >}}). You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.
