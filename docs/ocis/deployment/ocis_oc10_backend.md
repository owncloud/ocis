---
title: "ocis frontend with oc10 backend deployment scenario"
date: 2020-10-12T14:04:00+01:00
weight: 25
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_frontend_oc10_backend.md
---

{{< toc >}}


# ocis frontend with oc10 backend deployment scenario

This deployment scenario shows how to use ocis as frontend for a existing owncloud 10 installation.
ocis will allow owncloud 10 users to log in and work with their files.

## Overview
### Node Setup
ocis and oc10 running as docker containers behind traefik as reverse proxy
* Cloudflare DNS is resolving one domain for ocis and one for oc10
* Letsencrypt is providing valid ssl certificate for both domains

## Node Deployment

### Requirements
* Server running Ubuntu 20.04 is public availible with a static ip address
* Two A-records for both domains are pointing on the servers ip address
* Create user `$sudo adduser username`
* Add user to sudo group `$sudo usermod -aG sudo username`
* Add users pub key to `~/.ssh/authorized_keys`
* Setup ssh to permit authorisation only by ssh key
* Install docker `$sudo apt install docker.io`
* Add user to docker group `$sudo usermod -aG docker username`
* Install docker-compose via `$ sudo curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose` (docker compose version 1.27.4 as of today)
* Make docker-compose executable `$ sudo chmod +x /usr/local/bin/docker-compose`
* Environment variables for OCIS Stack are provided by .env file
* Change in `.env`

```
  OCIS_DOMAIN=ocis.domain.org
  OC10_DOMAIN=oc10.domain.org
```


### Stack
The application stack is separated in docker containers. One is a traefik proxy which is terminating ssl and forwards the https requests to the internal docker network. Additional, traefik is creating two certificates that are stored in the file `letsencrypt/acme.json` of the users home directory. In a local setup, this traefik is not included.
The next container is the ocis server which is exposing the webservice on port 9200 to traefic and provides the oidc provider konnectd to owncloud.
oc10 is running as a three container setup out of owncloud-server, a db container and a redis container as memcache storage.

### Config

#### Repository structure

```
ocis_oc10_backend  # rootfolder
│   .env
│   docker-compose.yml
│
└───ocis #ocis related config files
│   │   identifier-registration.yml
│   │   proxy-config.json
│
└───oc10 #owncloud 10 related files
    │   Dockerfile
    │
    └───apps
        │   graphapi-0.1.0.tar.gz
```

#### Traefik

In this deployment scenario, traefik requests letsencrypt to issue 2 ssl certificates, so two certificate resolver are needed. These are named according to the services, ocis for the ocis container and oc10 for the oc10 container.


```
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
Both container's traefik labels have to match with the correct resolvers and domains
```
  ocis:
    ...
    labels:
      ...
      - "traefik.http.routers.ocis.rule=Host(`${OCIS_DOMAIN}`)"
      ...
```
```
  oc10:
    ...
    labels:
      ...
      - "traefik.http.routers.oc10.rule=Host(`${OC10_DOMAIN}`)"
      ...
```

A folder for letsencypt to store the certificate needs to be created
`$ mkdir ~/letsencrypt`
This folder is bind to the docker container and the certificate is persistently stored into it.

#### ocis

Since ssl shall be terminated from traefik and inside of the docker network the services shall comunicate via http, `PROXY_TLS: "false"` as environment parameter for ocis has to be set.

For ocis 2 config files are provided.

```
│
└───ocis #ocis related config files
│   │   identifier-registration.yml
│   │   proxy-config.json
```

Changes need to be done in identifier-registration.yml to match the domains
Phoenix client needs the redirects uri's set to the ocis domain while oc10 client needs them to point on the owncloud domain

```
---
# OpenID Connect client registry.
clients:
  - id: phoenix
    name: OCIS
    application_type: web
    insecure: yes
    trusted: yes
    redirect_uris:
      - http://ocis.domain.com/
      - https://ocis.domain.com/
      - http://ocis.domain.com/oidc-callback.html
      - https://ocis.domain.com/oidc-callback.html
      - http://ocis.domain.com/oidc-silent-redirect.html
      - https://ocis.domain.com/oidc-silent-redirect.html
    origins:
      - http://ocis.domain.com
      - https://ocis.domain.com

  - id: oc10
    name: OC10
    application_type: web
    secret: super
    insecure: yes
    trusted: yes
    redirect_uris:
      - https://oc10.domain.com/apps/openidconnect/redirect/
      - https://oc10.domain.com/apps/openidconnect/redirect
    origins:
      - http://oc10.domain.com
      - https://oc10.domain.com
```

The second file is proxy-config.json which configures the ocis internal service proxy routes. The policy_selector selector needs to be changed to forward to the related backend. ocis proxy makes the decision in this scenario to which backend the request needs to be forwarded based on the user storage.

```
{
  "HTTP": {
    "Namespace": "works.owncloud"
  },
  "policy_selector": {
    "migration": {
      "acc_found_policy" : "reva",
      "acc_not_found_policy": "oc10",
      "unauthenticated_policy": "oc10"
  }
  "policies": [
    {
      "name": "reva",
      "routes": [
        {
          "endpoint": "/",
          "backend": "http://localhost:9100"
        },
        {
        ....
```

Glauth needs to be configured to utilize oc10 as primary user backend.

```
GLAUTH_BACKEND_DATASTORE: owncloud
GLAUTH_BACKEND_SERVERS: https://${OC10_DOMAIN}/apps/graphapi/v1.0
GLAUTH_BACKEND_BASEDN: dc=example,dc=org
STORAGE_STORAGE_METADATA_PROVIDER_DRIVER: owncloud
STORAGE_STORAGE_METADATA_DATA_PROVIDER_DRIVER: owncloud
ACCOUNTS_STORAGE_DISK_PATH: /var/tmp/ocis-accounts # Accounts fails to start when cs3 backend is used atm
```

To allow konnectd to glauth, ldap needs to be configured have to be set.

```
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

Owncloud 10 needs the graph api extensions to work in this setup. This extension is needed for Glauth to get oc10 users. It's necessary to add a image build step which extends owncloud/server:latest docker image with the app. The app is provided as tarball in the folder oc10/apps

```
└───oc10
│   │   Dockerfile
│   │
│   └───apps
│   │   │   graphapi-0.1.0.tar.gz
```

The docker files is pretty simple

```
# Take the latest owncloud/server image as base
FROM owncloud/server:latest

# Add the provided tarballs into oc10's apps folder
ADD apps/graphapi-0.1.0.tar.gz /var/www/owncloud/apps/
```

The build is triggered by the terminal command `docker-compose build` from the root folder.


Constraints: In this setup it's mandatory that the user has an email adress set in oc10