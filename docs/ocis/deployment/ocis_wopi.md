---
title: "oCIS with WOPI server"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_wopi.md
---

{{< toc >}}

## Overview

* oCIS, Wopi server and Collabora running behind Traefik as reverse proxy
* Collabora enables you to edit text documents in your browser
* Wopi server acts as a bridge to make the oCIS storage accessible to Collabora
* Traefik generating self signed certificates for local setup or obtaining valid SSL certificates for a server setup

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_wopi)

The docker stack consists 5 containers. One of them is Traefik, a proxy which is terminating SSL and forwards the requests to oCIS in the internal docker network.

The next container is oCIS itself in a configuration like the [oCIS with Traefik example]({{< ref "ocis_traefik" >}}), except that for this example a custom proxy and web UI configuration is used to enable the oCIS Wopi extension.

The oCIS WOPI server extension is running in another container and enables you to open files in Collabora from within ownCloud Web.

The last two containers are the WOPI server and Collabora.

## Server Deployment

### Requirements

* Linux server with docker and docker-compose installed
* Three domains set up and pointing to your server
  - ocis.* for serving oCIS
  - collabora.* for serving Collabora
  - wopiserver.* for serving the WOPI server
  - traefik.* for serving the Traefik dashboard

See also [example server setup]({{< ref "preparing_server" >}})


### Install oCIS and Traefik

* Clone oCIS repository

  `git clone https://github.com/owncloud/ocis.git`

* Go to the deployment example

  `cd ocis/deployments/examples/ocis_wopi`

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
    # IDP LDAP bind password. Must be changed in order to have a secure oCIS. Defaults to "idp".
    IDP_LDAP_BIND_PASSWORD=
    # Storage LDAP bind password. Must be changed in order to have a secure oCIS. Defaults to "reva".
    STORAGE_LDAP_BIND_PASSWORD=
    # JWT secret which is used for the storage provider. Must be changed in order to have a secure oCIS. Defaults to "Pive-Fumkiu4"
    OCIS_JWT_SECRET=
    # JWT secret which is used for uploads to create transfer tokens. Must be changed in order to have a secure oCIS. Defaults to "replace-me-with-a-transfer-secret"
    OCIS_TRANSFER_SECRET=

    ### Wopi server settings ###
    # oCIS Wopi server version. Defaults to "latest"
    OCIS_WOPISERVER_DOCKER_TAG=
    # cs3org wopi server version. Defaults to "latest"
    CS3ORG_WOPISERVER_DOCKER_TAG=
    # cs3org wopi server domain. Defaults to "wopiserver.owncloud.test"
    WOPISERVER_DOMAIN=
    # JWT secret which is used for the documents to be request by the Wopi client from the cs3org Wopi server. Must be change in order to have a secure Wopi server. Defaults to "LoremIpsum567"
    WOPI_JWT_SECRET=
    # JWT secret which is used for the documents to be request by the Wopi client from the cs3org Wopi server. Must be change in order to have a secure Wopi server. Defaults to "LoremIpsum123"
    WOPI_IOP_SECRET=

    ### Collabora settings ###
    # Domain of Collabora, where you can find the frontend. Defaults to "collabora.owncloud.test"
    COLLABORA_DOMAIN=
    # Admin user for Collabora. Defaults to blank, provide one to enable access
    COLLABORA_ADMIN_USER=
    # Admin password for COllabora. Defaults to blank, provide one to enable access
    COLLABORA_ADMIN_PASSWORD=

  ```

  You are installing oCIS on a server and Traefik will obtain valid certificates for you so please remove `INSECURE=true` or set it to `false`.

  If you want to use the Traefik dashboard, set TRAEFIK_DASHBOARD to `true` (default is `false` and therefore not active). If you activate it, you must set a domain for the Traefik dashboard in `TRAEFIK_DOMAIN=` eg. `TRAEFIK_DOMAIN=traefik.owncloud.test`.

  The Traefik dashboard is secured by basic auth. Default credentials are the user `admin` with the password `admin`. To set your own credentials, generate a htpasswd (eg. by using [an online tool](https://htpasswdgenerator.de/) or a cli tool).

  Traefik will issue certificates with LetsEncrypt and therefore you must set an email address in `TRAEFIK_ACME_MAIL=`.

  By default oCIS will be started in the `latest` version. If you want to start a specific version of oCIS set the version to `OCIS_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated).

  Set your domain for the oCIS frontend in `OCIS_DOMAIN=`, eg. `OCIS_DOMAIN=ocis.owncloud.test`.

  You also must override three default secrets in `IDP_LDAP_BIND_PASSWORD`, `STORAGE_LDAP_BIND_PASSWORD` and `OCIS_JWT_SECRET` in order to secure your oCIS instance. Choose some random strings eg. from the output of `openssl rand -base64 32`. For more information see [secure an oCIS instance]({{< ref "./#secure-an-ocis-instance" >}}).

  By default the oCIS WOPI server extension will be started in the `latest` version. If you want to start a specific version of oCIS WOPI server set the version to `OCIS_WOPISERVER_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis-wopiserver/tags?page=1&ordering=last_updated).

  By default the CS3Org WOPI server will also be started in the `latest` version. If you want to start a specific version of it, you can set the version to `CS3ORG_WOPISERVER_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/cs3org/wopiserver/tags?page=1&ordering=last_updated).

  Set your domain for the CS3Org WOPI server in `WOPISERVER_DOMAIN=`, where Collabora can download the files.

  You also must override the default WOPI JWT secret and the WOPI IOP secret, in order to have a secure setup. Do this by setting `WOPI_JWT_SECRET` and `WOPI_IOP_SECRET` to a long and random string.

  Now it's time to set up Collabora and you need to configure the Domain of Collabora in `COLLABORA_DOMAIN=`.

  If you want to use the Collabora admin panel you need to set user name and passwort for in `COLLABORA_ADMIN_USER=` and `COLLABORA_ADMIN_PASSWORD=`.

  Now you have configured everything and can save the file.

* Start the docker stack

  `docker-compose up -d`

* You now can visit oCIS and are able to open an office document in your browser.

## Local setup
For a more simple local ocis setup see [Getting started]({{< ref "../getting-started" >}})

This docker stack can also be run locally. One downside is that Traefik can not obtain valid SSL certificates and therefore will create self signed ones. This means that your browser will show scary warnings. Another downside is that you can not point DNS entries to your localhost. So you have to add static host entries to your computer.

On Linux and macOS you can add them to your `/etc/hosts` files like this:
```
127.0.0.1 ocis.owncloud.test
127.0.0.1 traefik.owncloud.test
127.0.0.1 collabora.owncloud.test
127.0.0.1 wopiserver.owncloud.test
```

After that you're ready to start the application stack:

`docker-compose up -d`

Open https://collabora.owncloud.test and https://wopiserver.owncloud.test  in your browser and accept the invalid certificate warning.

Open https://ocis.owncloud.test in your browser and accept the invalid certificate warning. You are now able to open an office document in your browser.
