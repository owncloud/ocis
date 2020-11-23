---
title: "oCIS with external IDP"
date: 2020-10-12T14:39:00+01:00
weight: 26
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_external_idp.md
---

{{< toc >}}

This scenario shows how to setup oCIS and konnectd as external IDP (identity provider). Both have separate domains and will be configured to work together.

## Overview

* Server 1: oCIS running behind traefik as reverse proxy
* Server 2: IDP running behind traefik as reverse proxy
* Valid ssl certificates for the domains for ssl termination

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_external_konnectd)



## Server Deployment

### Requirements

* 2 Linux servers, each with docker and docker-compose installed
* Two domains set up and pointing to the target server

See also [example server setup]({{< ref "preparing_server.md" >}})

### Install oCIS server

* Clone ocis repository

  `git clone https://github.com/owncloud/ocis.git`

* Copy example sub folder for ocisnode to /opt

  `cp deployment/examples/ocis_external_konnectd/ocisnode /opt/`

* Overwrite OCIS_DOMAIN and IDP_DOMAIN in .env with your-ocis.domain.com and your-idp.domain.com

  `sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/ocisnode/.env`

  `sed -i 's/idp.domain.com/your-idp.domain.com/g' /opt/ocisnode/.env`

* Change into deployment folder

  `cd /opt/ocisnode`

* Start application stack

  `docker-compose up -d`

### Install IDP server

* Clone ocis repository

  `git clone https://github.com/owncloud/ocis.git`

* Copy example sub folder for idpnode to /opt

  `cp deployment/examples/ocis_external_konnectd/idpnode /opt/`

* Overwrite OCIS_DOMAIN and IDP_DOMAIN in .env with your-ocis.domain.com and your-idp.domain.com

  `sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/idpnode/.env`

  `sed -i 's/idp.domain.com/your-idp.domain.com/g' /opt/idpnode/.env`

* Overwrite redirect uri with your-ocis.domain.com in identifier-registration.yml

  `sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/idpnode/config/identifier-registration.yml`

* Change into deployment folder

  `cd /opt/idpnode`

* Start application stack

  `docker-compose up -d`

### Configuration

#### Repository structure

```bash
ocis_external_konnectd  # rootfolder
└───ocisnode
│   │   docker-compose.yml
│   │   .env
│
└───idpnode
    │   docker-compose.yml
    │   .env
    └───config
        │   identifier-registration.yml
```

Both subfolders contain the dockr-compose files including additionaly conf files if required. The content of both folders has to be deployed on each node.

#### Traefik

Traefik is set up similar to the traefik example on both nodes.
The certificate resolvers are named similar to their services and behave exactly like in the other examples.

#### Konnectd

Konnectd as Openid provider needs the redirect url's to point to ocis.

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

Behind traefik, http is used to communicate between the services. Setting KONNECTD_TLS enforces it.

```yaml
      KONNECTD_TLS: '0'
```

In order to resolve users from glauth service on ocis node, Konnectd needs ldap settings to work properly.

```yaml
      LDAP_URI: ldap://${OCIS_DOMAIN}:9125
      LDAP_BINDDN: cn=konnectd,ou=sysusers,dc=example,dc=org
      LDAP_BINDPW: konnectd
      LDAP_BASEDN: ou=users,dc=example,dc=org
      LDAP_SCOPE: sub
      LDAP_LOGIN_ATTRIBUTE: cn
      LDAP_EMAIL_ATTRIBUTE: mail
      LDAP_NAME_ATTRIBUTE=: n
      LDAP_UUID_ATTRIBUTE: uid
      LDAP_UUID_ATTRIBUTE_TYPE: text
      LDAP_FILTER: (objectClass=posixaccount)
```

#### ocis

On the ocis node, the setting is following a standard scenario, except, that port 9125 needs to be exposed for the idp node to resolve ldap querries from Konnectd.

```yaml
ocis:
...
    ports:
      - 9200:9200
      - 9125:9125
...
```

## Local setup
For simple local ocis setup see [Getting started]({{< ref "../getting-started.md" >}})

Local setup coming soon