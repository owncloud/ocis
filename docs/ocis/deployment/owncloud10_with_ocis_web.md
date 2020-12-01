---
title: "ownCloud 10 with oCIS Web"
date: 2020-10-12T14:04:00+01:00
weight: 25
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: owncloud10_with_ocis_web.md
---

{{< toc >}}

This deployment scenario shows how to use oCIS Web as frontend for an existing ownCloud 10 production installation. It enables ownCloud 10 users to log in and work with their files using the new ownCloud Web. While the scenario includes an ownCloud 10 instance, it only exists to show the necessary configuration for your already existing ownCloud 10 installation.

## Overview

* oCIS setup serving ownCloud Web
* oCIS acting as OIDC IDP on the ownCloud 10 user database
* ownCloud 10 setup connected to oCIS
* DNS is resolving one domain for ocis and one for oc10
* Valid ssl certificates for the domains for ssl termination

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_external_konnectd)

{{< hint info >}}
In this setup it's mandatory that the users in ownCloud 10 are assigned to at least one group.
{{< /hint >}}

{{< hint info >}}
In this setup relies on graph-api app to be installed in ownCloud 10. This app is included by default beginning with ownCloud 10.6. If you are on a lower version, please install it manually. 
{{< /hint >}}

## Server Deployment

### Requirements

* Linux server with docker and docker-compose installed
* Three domains set up and pointing to your server
  - ocis.* for serving oCIS
  - oc10.* for serving 
  - traefik.* for serving the Traefik dashboard

See also [example server setup]({{< ref "preparing_server.md" >}})

### Install oCIS and Traefik

* Clone oCIS repository

  `git clone https://github.com/owncloud/ocis.git`

* Go to the deployment example

  `cd ocis/deployment/examples/ocis_oc10_backend`

* Open the `.env` file in a text editor
  The file by default looks like this:
  ```bash
  # If you're on a internet facing server please comment out following line.
  # It skips certificate validation for various parts of oCIS and is needed if you use self signed certificates.
  INSECURE=true

  ### Traefik settings ###
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

  ### oC10 ###
  # Domain of ownCloud 10, where you can find the frontend. Defaults to "oc10.owncloud.test"
  #OC10_DOMAIN=
  ```

  You are installing oCIS on a server and Traefik will obtain valid certificates for you so please remove `INSECURE=true` or set it to `false`.

  Set your domain for the Traefik dasboard in `TRAEFIK_DOMAIN=` eg. `TRAEFIK_DOMAIN=traefik.owncloud.test`.

  The Traefik dasboard is secured by basic auth. Default credentials are the user `admin` with the password `admin`. To set your own credentials, generate a htpasswd (eg. by using [an online tool](https://htpasswdgenerator.de/) or a cli tool).

  Traefik will issue certificates with LetsEncrypt and therefore you must set an email address in `TRAEFIK_ACME_MAIL=`.

  oCIS will by default started in the `latest` version. If you want to start a specific version of oCIS set the version to `OCIS_DOCKER_TAG=`. Available versions can be found on [Docker Hub](https://hub.docker.com/r/owncloud/ocis/tags?page=1&ordering=last_updated).

  Set your domain for the oCIS frontend in `OCIS_DOMAIN=`, eg. `OCIS_DOMAIN=ocis.owncloud.test`.

  Set your domain for the ownCloud 10 frontend in `OC10_DOMAIN=` eg. `OC10_DOMAIN=oc10.owncloud.test`.

  Now you have configured everything and can save the file.

* Start the docker stack

  `docker-compose up -d`

* You now can visit oCIS and Traefik dashboard on your configured domains


## Local setup
For a more simple local ocis setup see [Getting started]({{< ref "../getting-started.md" >}})

This docker stack can also be run locally. One downside is that Traefik can not obtain valid SSL certificates and therefore will create self signed ones. This means that your browser will show scary warnings. Another downside is that you can not point DNS entries to your localhost. So you have to add static host entries to your computer.

On Linux you can add them to your `/etc/hosts` files like this:
```
127.0.0.1 ocis.owncloud.test
127.0.0.1 oc10.owncloud.test
127.0.0.1 traefik.owncloud.test
```

After that you're ready to start the application stack:

`docker-compose up -d`

Open https://oc10.owncloud.test in your browser and accept the invalid certificate warning. You now can login with the ownCloud 10 default user "admin" and password "admin". As you might have noticed, you did not see the login prompt of ownCloud 10. This was the login prompt of oCIS. When you go to application you can both in oCIS web and ownCloud 10 see a switch to switch vice versa.
