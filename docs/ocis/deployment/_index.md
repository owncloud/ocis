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
This section handles deployments and operations for admins and people who are interested in how versatile oCIS is. If you want to just try oCIS you may also follow [Getting started]({{< ref "../getting-started.md" >}}).

### Setup oCIS on your server
oCIS deployments are super simple, yet there are many configurations possible for advanced setups.

- [Basic oCIS setup]({{< ref "basic-remote-setup.md" >}}) - configure domain, certificates and port
- [oCIS setup with Traefik for SSL termination]({{< ref "ocis_traefik.md" >}})
- [oCIS setup with Keycloak as identity provider]({{< ref "ocis_keycloak.md" >}})

### Migrate an existing ownCloud 10
You can run ownCloud 10 and oCIS together. This allows you to use new parts of oCIS already with ownCloud 10 and also to have a smooth transition for users from ownCloud 10 to oCIS.

- [ownCloud 10 setup with oCIS serving ownCloud Web and acting as OIDC provider]({{< ref "owncloud10_with_oc_web.md" >}}) - This allows you to switch between the traditional ownCloud 10 frontend and the new ownCloud Web frontend
- Run ownCloud 10 and oCIS in parallel - together
- Migrate users from ownCloud 10 to oCIS
