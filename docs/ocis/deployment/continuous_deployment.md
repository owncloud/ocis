---
title: "Continuous Deployment"
date: 2020-10-12T14:04:00+01:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: continuous_deployment.md
---

{{< toc >}}

We are continuously deploying the following deployment examples. Every example is deployed in two flavors:

- Latest: reflects the current master branch state of oCIS and will be updated with every commit to master
- Released: reflects the newest release state (currently latest release of version 1) and will be updated with every release

The configuration for the continuous deployment can be found in the [oCIS repository](https://github.com/owncloud/ocis/tree/master/deployments/continuous-deployment-config).

# oCIS with Traefik

Credentials:

- oCIS: see [default demo users]({{< ref "../getting-started#login-to-owncloud-web" >}})

## Latest

- oCIS: [ocis.ocis-traefik.latest.owncloud.works](https://ocis.ocis-traefik.latest.owncloud.works)

## Released

- oCIS: [ocis.ocis-traefik.released.owncloud.works](https://ocis.ocis-traefik.released.owncloud.works)

# oCIS with WOPI server

Credentials:

- oCIS: see [default demo users]({{< ref "../getting-started#login-to-owncloud-web" >}})

## Latest

- oCIS: [ocis.ocis-wopi.latest.owncloud.works](https://ocis.ocis-wopi.latest.owncloud.works)

## Released

- oCIS: [ocis.ocis-wopi.released.owncloud.works](https://ocis.ocis-wopi.released.owncloud.works)

# oCIS with latest ownCloud Web

Credentials:

- oCIS: see [default demo users]({{< ref "../getting-started#login-to-owncloud-web" >}})

## Latest

- oCIS: [ocis.ocis-web.latest.owncloud.works](https://ocis.ocis-web.latest.owncloud.works)

# oCIS with Keycloak

Credentials:

- oCIS: see [default demo users]({{< ref "../getting-started#login-to-owncloud-web" >}})
- Keycloak:
  - username: admin
  - password: admin

## Latest

- oCIS: [ocis.ocis-keycloak.latest.owncloud.works](https://ocis.ocis-keycloak.latest.owncloud.works)
- Keycloak: [keycloak.ocis-keycloak.latest.owncloud.works](https://keycloak.ocis-keycloak.latest.owncloud.works)

## Released

- oCIS: [ocis.ocis-keycloak.released.owncloud.works](https://ocis.ocis-keycloak.released.owncloud.works)
- Keycloak: [keycloak.ocis-keycloak.released.owncloud.works](https://keycloak.ocis-keycloak.released.owncloud.works)

# Parallel deployment of oC10 and oCIS

Credentials:

- oC10 / oCIS: see [default demo users]({{< ref "../getting-started#login-to-owncloud-web" >}})
- Keycloak:
  - username: admin
  - password: admin
- LDAP management:
  - username: cn=admin,dc=owncloud,dc=com
  - password: admin

## Latest

- oC10 / oCIS: [cloud.oc10-ocis-parallel.latest.owncloud.works](https://cloud.oc10-ocis-parallel.latest.owncloud.works)
- LDAP management: [ldap.oc10-ocis-parallel.latest.owncloud.works](https://ldap.oc10-ocis-parallel.latest.owncloud.works)
- Keycloak: [keycloak.oc10-ocis-parallel.latest.owncloud.works](https://keycloak.oc10-ocis-parallel.latest.owncloud.works)

# oCIS with Hello extension

Credentials:

- oCIS: see [default demo users]({{< ref "../getting-started#login-to-owncloud-web" >}})

## Latest

- oCIS: [ocis.ocis-hello.latest.owncloud.works](https://ocis.ocis-hello.latest.owncloud.works)

# oCIS with S3 storage backend (MinIO)

Credentials:

- oCIS: see [default demo users]({{< ref "../getting-started#login-to-owncloud-web" >}})
- MinIO:
  - access key: ocis
  - secret access key: ocis-secret-key

## Latest

- oCIS: [ocis.ocis-s3.latest.owncloud.works](https://ocis.ocis-s3.latest.owncloud.works)
- MinIO: [minio.ocis-s3.latest.owncloud.works](https://minio.ocis-s3.latest.owncloud.works)

# oCIS with LDAP for users and groups

Credentials:

- oCIS: see [default demo users]({{< ref "../getting-started#login-to-owncloud-web" >}})
- LDAP admin:
  - username: cn=admin,dc=owncloud,dc=com
  - password: admin

## Latest

- oCIS: [ocis.ocis-ldap.latest.owncloud.works](https://ocis.ocis-ldap.latest.owncloud.works)
- LDAP admin: [ldap.ocis-ldap.latest.owncloud.works](https://ldap.ocis-ldap.latest.owncloud.works)

## Released

- oCIS: [ocis.ocis-ldap.released.owncloud.works](https://ocis.ocis-ldap.released.owncloud.works)
- LDAP admin: [ldap.ocis-ldap.released.owncloud.works](https://ldap.ocis-ldap.released.owncloud.works)
