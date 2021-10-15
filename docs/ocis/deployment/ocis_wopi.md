---
title: "oCIS with WOPI server"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_wopi.md
---

{{< toc >}}

{{< hint warning >}}
OnlyOffice and CodiMD are not yet fully integrated and there are known issues. For the current state please have a look at [owncloud/ocis#2595](https://github.com/owncloud/ocis/issues/2595)
{{< /hint >}}

## Overview

* oCIS, Wopi server, Collabora, OnlyOffice and CodiMD running behind Traefik as reverse proxy
* Collabora, OnlyOffice and CodiMD enable you to edit documents in your browser
* Wopi server acts as a bridge to make the oCIS storage accessible to Collabora, OnlyOffice and CodiMD
* Traefik generating self signed certificates for local setup or obtaining valid SSL certificates for a server setup

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_wopi)

The docker stack consists of 10 containers. One of them is Traefik, a proxy which is terminating SSL and forwards the requests to oCIS in the internal docker network.

The next container is oCIS itself in a configuration like the [oCIS with Traefik example]({{< ref "ocis_traefik" >}}), except that for this example a custom mimetype configuration is used.

There are three oCIS app driver containers that register Collabora, OnlyOffice and CodiMD at the app registry.

The last four containers are the WOPI server, Collabora, OnlyOffice and CodiMD.

## Server Deployment

### Requirements

* Linux server with docker and docker-compose installed
* Three domains set up and pointing to your server
  - ocis.* for serving oCIS
  - collabora.* for serving Collabora
  - onlyoffice.* for serving OnlyOffice
  - codimd.* for serving CodiMD
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
    STORAGE_TRANSFER_SECRET=
    # Machine auth api key secret. Must be changed in order to have a secure oCIS. Defaults to "change-me-please"
    OCIS_MACHINE_AUTH_API_KEY=

    ### Wopi server settings ###
    # cs3org wopi server version. Defaults to "latest"
    WOPISERVER_DOCKER_TAG=
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
    # Admin password for Collabora. Defaults to blank, provide one to enable access
    COLLABORA_ADMIN_PASSWORD=

    ### OnlyOffice settings ###
    # Domain of OnlyOffice, where you can find the frontend. Defaults to "onlyoffice.owncloud.test"
    ONLYOFFICE_DOMAIN=

    ### CodiMD settings ###
    # Domain of Collabora, where you can find the frontend. Defaults to "codimd.owncloud.test"
    CODIMD_DOMAIN=
    # Secret which is used for the communication with the WOPI server. Must be changed in order to have a secure CodiMD. Defaults to "LoremIpsum456"
    CODIMD_SECRET=
  ```

  You are installing oCIS on a server and Traefik will obtain valid certificates for you so please remove `INSECURE=true` or set it to `false`.

  If you want to use the Traefik dashboard, set TRAEFIK_DASHBOARD to `true` (default is `false` and therefore not active). If you activate it, you must set a domain for the Traefik dashboard in `TRAEFIK_DOMAIN=` eg. `TRAEFIK_DOMAIN=traefik.owncloud.test`.

  The Traefik dashboard is secured by basic auth. Default credentials are the user `admin` with the password `admin`. To set your own credentials, generate a htpasswd (eg. by using [an online tool](https://htpasswdgenerator.de/) or a cli tool).

  Traefik will issue certificates with LetsEncrypt and therefore you must set an email address in `TRAEFIK_ACME_MAIL=`.

  By default oCIS will be started in the `latest` version. If you want to start a specific version of oCIS set the version to `OCIS_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated).

  Set your domain for the oCIS frontend in `OCIS_DOMAIN=`, eg. `OCIS_DOMAIN=ocis.owncloud.test`.

  You also must override three default secrets in `IDP_LDAP_BIND_PASSWORD`, `STORAGE_LDAP_BIND_PASSWORD` and `OCIS_JWT_SECRET` in order to secure your oCIS instance. Choose some random strings eg. from the output of `openssl rand -base64 32`. For more information see [secure an oCIS instance]({{< ref "./#secure-an-ocis-instance" >}}).

  By default the CS3Org WOPI server will also be started in the `latest` version. If you want to start a specific version of it, you can set the version to `WOPISERVER_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/cs3org/wopiserver/tags?page=1&ordering=last_updated).

  Set your domain for the CS3Org WOPI server in `WOPISERVER_DOMAIN=`, where all office suites can download the files via the WOPI protocol.

  You also must override the default WOPI JWT secret and the WOPI IOP secret, in order to have a secure setup. Do this by setting `WOPI_JWT_SECRET` and `WOPI_IOP_SECRET` to a long and random string.

  Now it's time to set up Collabora and you need to configure the domain of Collabora in `COLLABORA_DOMAIN=`.

  If you want to use the Collabora admin panel you need to set user name and passwort for in `COLLABORA_ADMIN_USER=` and `COLLABORA_ADMIN_PASSWORD=`.

  Next up is OnlyOffice, which also needs a domain in `ONLYOFFICE_DOMAIN=`.

  The last configuration options are for CodiMD, which needs a domain in `CODIMD_DOMAIN=` and a random secret in `CODIMD_SECRET=`.

  Now you have configured everything and can save the file.

* Start the docker stack

  `docker-compose up -d`

* You now can visit oCIS and are able to open an office document in your browser. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.

## Local setup
For a more simple local ocis setup see [Getting started]({{< ref "../getting-started" >}})

This docker stack can also be run locally. One downside is that Traefik can not obtain valid SSL certificates and therefore will create self signed ones. This means that your browser will show scary warnings. Another downside is that you can not point DNS entries to your localhost. So you have to add static host entries to your computer.

On Linux and macOS you can add them to your `/etc/hosts` files like this:
```
127.0.0.1 ocis.owncloud.test
127.0.0.1 traefik.owncloud.test
127.0.0.1 collabora.owncloud.test
127.0.0.1 onlyoffice.owncloud.test
127.0.0.1 codimd.owncloud.test
127.0.0.1 wopiserver.owncloud.test
```

After that you're ready to start the application stack:

`docker-compose up -d`

Open https://collabora.owncloud.test, https://onlyoffice.owncloud.test, https://codimd.owncloud.test and https://wopiserver.owncloud.test  in your browser and accept the invalid certificate warning.

Open https://ocis.owncloud.test in your browser and accept the invalid certificate warning. You are now able to open an office document in your browser. You may need to wait some minutes until all services are fully ready, so make sure that you try to reload the pages from time to time.
