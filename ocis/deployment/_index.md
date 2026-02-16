---
title: "Deployment"
date: 2020-10-01T20:35:00+01:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

{{< toc >}}

## Deployments scenarios and examples
This section handles deployments and operations for admins and people who are interested in how versatile oCIS is. If you want to just try oCIS you may also follow [Getting started]({{< ref "../getting-started" >}}).

### Setup oCIS on your server
oCIS deployments are super simple, yet there are many configurations possible for advanced setups.

- [Basic oCIS setup]({{< ref "basic-remote-setup" >}}) - configure domain, certificates and port
- [oCIS setup with Keycloak as identity provider]({{< ref "ocis_keycloak" >}})
- [Flexible oCIS setup with WebOffice and Search capabilities]({{< ref "ocis_full" >}})
- [Parallel deployment of oC10 and oCIS]({{< ref "oc10_ocis_parallel" >}})
- [oCIS with the Hello extension example]({{< ref "ocis_hello" >}})


## Secure an oCIS instance

oCIS no longer has any default secrets in versions later than oCIS 1.20.0. Therefore you're no
longer able to start oCIS without generating / setting all needed secrets.

The recommended way is to use `ocis init` for that. It will generate a secure config file for you.
