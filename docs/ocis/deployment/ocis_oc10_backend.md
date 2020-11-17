---
title: "ocis frontend with oc10 backend deployment scenario"
date: 2020-10-12T14:04:00+01:00
weight: 25
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_frontend_oc10_backend.md
---

{{< toc >}}

This deployment scenario shows how to use ocis as frontend for an existing ownCloud 10 production installation. It enables
ownCloud 10 users to log in and work with their files using the new ocis-web UI. While the scenario includes
an ownCloud 10 instance, it only exists to show the necessary configuration for your already existing ownCloud 10
installation.

The described setup can also be used to do a zero-downtime migration from ownCloud 10 to ocis.

## Overview

### Node Setup

* ocis and oc10 running as docker containers behind traefik as reverse proxy
* Cloudflare DNS is resolving one domain for ocis and one for oc10
* Letsencrypt is providing valid ssl certificate for both domains

## Node Deployment

### Requirements

* Server running Ubuntu 20.04 is publicly available with a static ip address
* Two A-records for both domains are pointing to the servers ip address
* Create user

  `$ sudo adduser username`

* Add user to sudo group

  `$ sudo usermod -aG sudo username`

* Add users pub key to `~/.ssh/authorized_keys`
* Setup ssh to permit authorisation only by ssh key
* Install docker

  `$ sudo apt install docker.io`

* Add user to docker group

  `$ sudo usermod -aG docker username`

* Install docker-compose via

  `$ sudo curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose`

  (docker compose version 1.27.4 as of today)
* Make docker-compose executable

  `$ sudo chmod +x /usr/local/bin/docker-compose`

* Environment variables for oCIS Stack are provided by .env file

### Setup on server

* Clone ocis repository

  `git clone https://github.com/owncloud/ocis.git`

* Copy example folder to /opt
  `cp -r deployments/examples/ocis_oc10_backend /opt/`

* Change into deployment folder

  `cd /opt/ocis_oc10_backend`

* Overwrite OCIS_DOMAIN and OC10_DOMAIN in .env with your-ocis.example.org and your-oc10.example.org

  `sed -i 's/ocis.example.org/your-ocis.example.org/g' /opt/ocis_oc10_backend/.env`

  `sed -i 's/oc10.example.org/your-oc10.example.org/g' /opt/ocis_oc10_backend/.env`

* Start application stack

  `docker-compose up -d`

  The domains from your `.env` will be used for building the configuration files during the docker start.

### Stack

The application stack is separated in docker containers. One is a traefik proxy which is terminating ssl and forwards the https requests to the internal docker network. Additionally, traefik is creating two certificates that are stored in the file `letsencrypt/acme.json` of the users home directory. In a local setup, this traefik is not included.
The next container is the ocis server which is exposing the webservice on port 9200 to traefik and provides the oidc provider `konnectd` to owncloud.
oc10 is running as a three container setup out of owncloud-server, a db container and a redis container as memcache storage.

### Config

#### Repository structure

```bash
ocis_oc10_backend  # rootfolder
│   .env
│   docker-compose.yml
│
└───ocis #ocis related config files
│   └───config
│   │   └───web
│   │   │   └───config.json
│   │   │   identifier-registration.yaml
│   │   │   proxy-config.json
│   └───Dockerfile
│
└───oc10 #owncloud 10 related files
    └───apps
    │   └───graphapi-0.1.0.tar.gz
    └───overlay
    │   └───etc
    │       └───templates
    │           └───config.php
    └───Dockerfile
```

#### Traefik

In this deployment scenario, traefik requests letsencrypt to issue 2 ssl certificates, so two certificate resolvers are needed. These are named according to the services, ocis for the ocis container and oc10 for the oc10 container.

```yaml
...
  traefik:
    image: "traefik:v2.2"
    container_name: "traefik"
    command:
      ...
      # Ocis certificate resolver
      - "--certificatesresolvers.ocis.acme.tlschallenge=true"
      - "--certificatesresolvers.ocis.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"
      - "--certificatesresolvers.ocis.acme.email=user@${OCIS_DOMAIN}"
      - "--certificatesresolvers.ocis.acme.storage=/letsencrypt/acme-ocis.json"
      # OC10 certificate resolver
      - "--certificatesresolvers.oc10.acme.tlschallenge=true"
      - "--certificatesresolvers.oc10.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"
      - "--certificatesresolvers.oc10.acme.email=user@${OCIS_DOMAIN}"
      - "--certificatesresolvers.oc10.acme.storage=/letsencrypt/acme-oc10.json"
...
```

Both containers' traefik labels have to match the correct resolvers and domains

```yaml
  ocis:
    ...
    labels:
      ...
      - "traefik.http.routers.ocis.rule=Host(`${OCIS_DOMAIN}`)"
      ...
```

```yaml
  oc10:
    ...
    labels:
      ...
      - "traefik.http.routers.oc10.rule=Host(`${OC10_DOMAIN}`)"
      ...
```

A folder for letsencypt to store the certificate needs to be created
`$ mkdir ~/letsencrypt`
This folder is bound to the docker container and the certificate is persisted into it.

#### ocis

We will make use of some services from the ocis server package:
- `konnectd` for OpenID Connect (oidc). Your ownCloud 10 will need to switch the login method to oidc (see oc10 section), but user credentials remain the same.
- `proxy` a reverse proxy which decides where to route your requests to.
- `ocis-phoenix` serves the new ownCloud Web frontend.
- `accounts` learns your oc10 users and groups and will allow us to handle migration on a per-user basis later on.

Three config file templates are provided for ocis. All of them contain placeholder URLs which are replaced with
the URLs from your `.env` file during the docker build step. This section describes the configuration in detail, so
that you can make changes for your environment if necessary.

```bash
│
└───ocis #ocis related config files
│   └───web
│   │   └───config.json
│   │   identifier-registration.yaml
│   │   proxy-config.json
```

##### web/config.json

This is the configuration file for the new ownCloud Web frontend. The *server* domain needs to point to your ocis container,
since the `proxy` will take care of routing all requests - including oc10 backend requests - to the correct endpoints.

The *openIdConnect* block contains information required for ownCloud Web for retrieving users from your Identity Provider (IdP, in this case konnectd).

With the *applications* block you can define URLs which appear in either the `application switcher` or the `user menu` in ownCloud Web. For this deployment
we preconfigured it with a link to the classic web frontend, if users need access to applications which have not been ported to the new ownCloud Web frontend, yet.

The *apps* block contains the list of built in ownCloud Web extensions that are supposed to be enabled. Please note that the *files* extension is required at all times.

More options for ownCloud Web config can be found in the [developer documentation](https://owncloud.github.io/clients/web/).

```json
{
  "server": "https://ocis.example.org",
  "theme": "owncloud",
  "version": "0.1.0",
  "openIdConnect": {
    "metadata_url": "https://ocis.example.org/.well-known/openid-configuration",
    "authority": "https://ocis.example.org",
    "client_id": "phoenix",
    "response_type": "code",
    "scope": "openid profile email"
  },
  "applications": [
    {
      "title": {
        "en": "Classic Design",
        "de": "Klassisches ownCloud"
      },
      "icon": "switch_ui",
      "url": "https://ocis.example.org",
      "target": "_self"
    },
    {
      "title": {
        "en": "Settings",
        "de": "Einstellungen"
      },
      "icon": "application",
      "url": "https://ocis.example.org/index.php/settings/personal",
      "target": "_self",
      "menu": "user"
    }
  ],
  "apps": [
    "files",
    "draw-io",
    "markdown-editor",
    "media-viewer"
  ]
}
```

##### identifier-registration.yaml

The `identifier registration` configuration registers clients for oidc, namely phoenix (which is ownCloud Web) and
ownCloud 10. There is also dynamic client registration available if needed.

```yaml
---
# OpenID Connect client registry.
clients:
  - id: phoenix
    name: OCIS
    application_type: web
    insecure: yes
    trusted: yes
    redirect_uris:
      - http://ocis.example.org/
      - https://ocis.example.org/
      - http://ocis.example.org/oidc-callback.html
      - https://ocis.example.org/oidc-callback.html
      - http://ocis.example.org/oidc-silent-redirect.html
      - https://ocis.example.org/oidc-silent-redirect.html
    origins:
      - http://ocis.example.org
      - https://ocis.example.org

  - id: oc10
    name: OC10
    application_type: web
    secret: super
    insecure: yes
    trusted: yes
    redirect_uris:
      - http://oc10.example.org/apps/openidconnect/redirect
      - https://oc10.example.org/apps/openidconnect/redirect
    origins:
      - http://oc10.example.org
      - https://oc10.example.org
```

##### proxy-config.json

With the `proxy config` you can configure endpoints of internal services for the ocis reverse proxy. Since we only have
one backend without any migration so far, we can use a static proxy policy selector.

```yaml
{
  "HTTP": {
    "Namespace": "works.owncloud"
  },
  "policy_selector": {
      "static": {
        "policy": "oc10"
      }
    },
  "policies": [
    {
      "name": "oc10",
      "routes": [
        {
          "endpoint": "/",
          "backend": "http://localhost:9100"
        },
        {
        ....
```

##### Environment variables in docker-compose.yaml

There are some environment variables needed for the used ocis services. The most important part is that oidc connects
to the user backend of ownCloud 10. This is achieved by exposing the user backend with the `graph` api plugin
in ownCloud 10 and connecting to it with `glauth` in ocis.

Glauth needs to be configured to utilize oc10 as primary user backend:
```yaml
GLAUTH_BACKEND_DATASTORE: owncloud
GLAUTH_BACKEND_SERVERS: https://${OC10_DOMAIN}/apps/graphapi/v1.0
```

To allow konnectd to connect to glauth, ldap needs to be configured:

```yaml
# Konnectd ldap setup
LDAP_URI: ldap://localhost:9125
LDAP_BINDDN: "cn=admin,dc=example,dc=org"
LDAP_BINDPW: "admin"
LDAP_BASEDN: "dc=example,dc=org"
LDAP_SCOPE: sub
LDAP_LOGIN_ATTRIBUTE: uid
LDAP_EMAIL_ATTRIBUTE: mail
LDAP_NAME_ATTRIBUTE: givenName
LDAP_UUID_ATTRIBUTE: uid
LDAP_UUID_ATTRIBUTE_TYPE: text
LDAP_FILTER: "(objectClass=posixaccount)"
```

#### oc10

OwnCloud 10 needs the graph api extensions to work in this setup. This extension is needed for Glauth to get oc10 users. It's necessary to add an image build step which extends owncloud/server:latest docker image with the app. The app is provided as tarball in the folder oc10/apps.

```bash
└───oc10
│   │   Dockerfile
│   │
│   └───apps
│   │   │   graphapi-0.1.0.tar.gz
```

The docker file is pretty simple

```Dockerfile

# Take the latest owncloud/server image as base
FROM owncloud/server:latest

# Add the provided tarballs into oc10's apps folder
ADD apps/graphapi-0.1.0.tar.gz /var/www/owncloud/apps/
```

The build is triggered by the terminal command `docker-compose build` from the root folder.

Constraints: In this setup it's mandatory that the user has an email address set and is assigned to at least one group in oc10.
Especially the default admin user doesn't have an email assigned. If your admin user doesn't have an email address, yet, please
set one: `docker-compose exec owncloud occ user:modify admin email "admin@example.org"`

## Local deployment

If you want to start the bridge setup on your local development machine, there are a few steps necessary:

### Domains
Instead of replacing the domains in the config files you can add `ocis.example.org` and `oc10.example.org` as localhost
aliases to your `/etc/hosts` file:
```
127.0.0.1       oc10.example.org
127.0.0.1       ocis.example.org
```

### Disable certificate checks
The `docker-compose.yml` file contains some `*INSECURE` environment variables for enabling or disabling certificate checks.
To disable certificate checks, set `INSECURE=true` in your `.env` file.
