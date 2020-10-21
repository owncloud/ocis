---
title: "ocis with traefik deployment scenario"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_traefik.md
---

{{< toc >}}


# ocis traefik deployment scenario

## Overview
ocis running on a hcloud node behind traefik as reverse proxy
* Cloudflare DNS is resolving the domain
* Letsencrypt provides a ssl certificate for the domain
* Traefik docker container terminates ssl and forwards http requests to ocis

## Node

### Requirements
* Server running Ubuntu 20.04 is public availible with an static ip address
* An A-record for domain is pointing on the servers ip address
* Create user `$sudo adduser username`
* Add user to sudo group `$sudo usermod -aG sudo username`
* Add users pub key to `~/.ssh/authorized_keys`
* Setup sshd to forbid root access and permit authorisation only by ssh key
* Install docker `$sudo apt install docker.io`
* Add user to docker group `$sudo usermod -aG docker username`
* Install docker-compose via `$ sudo curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose` (docker compose version 1.27.4 as of today)
* Make docker-compose executable `$ sudo chmod +x /usr/local/bin/docker-compose`
* Environment variables for OCIS Stack are provided by .env file

### Setup on server

- Clone ocis repository

  ```git clone https://github.com/owncloud/ocis.git```

- Copy example folder to /opt
  ```cp deployment/examples/ocis_traefik /opt/```

- Overwrite OCIS_DOMAIN in .env with your.domain.com
  ```sed -i 's/ocis.domain.com/your.domain.com/g' /opt/ocis_traefik/.env```

- Overwrite redirect uri with your.domain.com in identifier-registration.yml
  ```sed -i 's/ocis.domain.com/your.domain.com/g' /opt/ocis_traefik/config/identifier-registration.yml```

- Change into deployment folder
  ```cd /opt/ocis_traefik```

- Start application stack
  ```docker-compose up -d```

### Stack
The application stack contains two containers. The first one is a traefik proxy which is terminating ssl and forwards the requests to the internal docker network. Additional, traefik is creating a certificate that is stored in `acme.json` in the folder `letsencrypt` inside the users home directory.
The second one is th ocis server which is exposing the webservice on port 9200 to traefic.

### Config
Edit docker-compose.yml file to fit your domain setup
```
...
  traefik:
    image: "traefik:v2.2"
    ...
    labels:
      ...
      # Email address is neccesary for certificate creation
      - "--certificatesresolvers.ocisresolver.acme.email=username@${OCIS_DOMAIN}"
...
```

```
  ocis:
    container_name: ocis
    ...
    labels:
      ...
      # This is the domain for which traefik is creating the certificate from letsencrypt
      - "traefik.http.routers.ocis.rule=Host(`${OCIS_DOMAIN}`)"
      ...
```

In this example, ssl is terminated from traefik while inside of the docker network the services are comunicating via http. For this `PROXY_TLS: "false"` as environment parameter for ocis has to be set.

For ocis to work properly it's neccesary to provide one config file.
Change identifier-registration.yml to match your domain.

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
```

To make it availible for ocis inside of the container, `config` hast to be mounted as volume.

```
    ...
    volumes:
      - ./config:/etc/ocis
    environment:
      ...
      KONNECTD_IDENTIFIER_REGISTRATION_CONF: "/etc/ocis/identifier-registration.yml"
      ...
```
