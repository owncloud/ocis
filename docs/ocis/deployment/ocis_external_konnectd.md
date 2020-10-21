---
title: "ocis with konnectd on external node deployment scenario"
date: 2020-10-12T14:39:00+01:00
weight: 26
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_external_konnectd.md
---

{{< toc >}}


# ocis with konnectd on external node deployment scenario

This scenario shows how to setup ocis with konnectd as idp running on a separate node. Both node are having separate domains pointing on the servers.

# ocis traefik deployment scenario

## Overview
ocis and konnectd running on linux nodes behind traefik as reverse proxy
* Cloudflare DNS is resolving the domains
* Letsencrypt provides ssl certificates for the domains
* Traefik docker container terminates ssl and forwards http requests to the services

## Nodes

### Requirements
* Server running Ubuntu 20.04 is public availible with a static ip address
* Two A-records for both domains are pointing on the servers ip address
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
  `$ sudo curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose` (docker compose version 1.27.4 as of today)
* Make docker-compose executable
  `$ sudo chmod +x /usr/local/bin/docker-compose`
* Environment variables for OCIS Stack are provided by .env file

### Setup on ocis server

- Clone ocis repository

  ```git clone https://github.com/owncloud/ocis.git```

- Copy example sub folder for ocisnode to /opt
  ```cp deployment/examples/ocis_external_konnectd/ocisnode /opt/```

- Overwrite OCIS_DOMAIN and IDP_DOMAIN in .env with your-ocis.domain.com and your-idp.domain.com
  ```
  sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/ocisnode/.env
  sed -i 's/idp.domain.com/your-idp.domain.com/g' /opt/ocisnode/.env
  ```

- Change into deployment folder
  ```cd /opt/ocisnode```

- Start application stack
  ```docker-compose up -d```

### Setup on idp server

- Clone ocis repository

  ```git clone https://github.com/owncloud/ocis.git```

- Copy example sub folder for idpnode to /opt
  ```cp deployment/examples/ocis_external_konnectd/idpnode /opt/```

- Overwrite OCIS_DOMAIN and IDP_DOMAIN in .env with your-ocis.domain.com and your-idp.domain.com
  ```
  sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/idpnode/.env
  sed -i 's/idp.domain.com/your-idp.domain.com/g' /opt/idpnode/.env
  ```

- Overwrite redirect uri with your-ocis.domain.com in identifier-registration.yml
  ```
  sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/idpnode/config/identifier-registration.yml
  ```

- Change into deployment folder
  ```cd /opt/idpnode```

- Start application stack
  ```docker-compose up -d```

### Stack
On both nodes, a traefik dokcer container is terminating ssl and forwards the http requests to the services. The nodes are named according to their services.

### Config

#### Repository structure

```
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

Behind traefik, http is used to communicate between the services. Setting KONNECTD_TLS enforces it.

```
      KONNECTD_TLS: '0'
```

In order to resolve users from glauth service on ocis node, Konnectd needs ldap settings to work properly.

```
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

```
ocis:
...
    ports:
      - 9200:9200
      - 9125:9125
...
```