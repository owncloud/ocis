---
title: "oCIS with Traefik"
date: 2020-10-12T14:04:00+01:00
weight: 24
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_traefik.md
---

{{< toc >}}

## Overview

* oCIS running behind traefik as reverse proxy
* Valid ssl certificates for the domains for ssl termination

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_traefik)



## Server Deployment

### Requirements

* Linux server(s) with docker and docker-compose installed
* Two domains set up and pointing to your server(s)

See also [example server setup]({{< ref "preparing_server.md" >}})


### Install oCIS and Traefik

The application stack contains two containers. The first one is a traefik proxy which is terminating ssl and forwards the requests to the internal docker network. Additional, traefik is creating a certificate that is stored in `acme.json` in the folder `letsencrypt` inside the users home directory.
The second one is th ocis server which is exposing the webservice on port 9200 to traefik.

* Clone ocis repository

  `git clone https://github.com/owncloud/ocis.git`

* Copy example folder to /opt

  `cp deployment/examples/ocis_traefik /opt/`

* Overwrite OCIS_DOMAIN in .env with your.domain.com

  `sed -i 's/ocis.domain.com/your.domain.com/g' /opt/ocis_traefik/.env`

* Overwrite redirect uri with your.domain.com in identifier-registration.yml

  `sed -i 's/ocis.domain.com/your.domain.com/g' /opt/ocis_traefik/config/identifier-registration.yml`

* Change into deployment folder

  `cd /opt/ocis_traefik`

* Start application stack

  `docker-compose up -d`

### Configuration

Edit docker-compose.yml file to fit your domain setup

```yaml
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

```yaml
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

```yaml
    ...
    volumes:
      - ./config:/etc/ocis
    environment:
      ...
      KONNECTD_IDENTIFIER_REGISTRATION_CONF: "/etc/ocis/identifier-registration.yml"
      ...
```

## Local setup
For simple local ocis setup see [Getting started]({{< ref "../getting-started.md" >}})

Local setup with Traefik coming soon