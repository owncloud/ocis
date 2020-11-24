---
title: "Deployment"
date: 2020-10-01T20:35:00+01:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: _index.md
---

{{< toc >}}

## Deployments scenarios and examples
This section handles deployments and operations for admins. If want to just try oCIS you may also follow [Getting started]({{< ref "../getting-started.md" >}}).

### Setup oCIS on your server
oCIS deployments are super simple, yet there are many configurations possible for advanced setups.

- [Basic oCIS setup]({{< ref "basic-remote-setup.md" >}}) - configure domain, certificates and port
- [oCIS setup with Traefik for ssl termination]({{< ref "ocis_traefik.md" >}})
- [oCIS setup with external OIDC IDP]({{< ref "ocis_external_idp.md" >}})

### Migrate an existing ownCloud 10
You can run ownCloud 10 and oCIS together. This allows you to use new parts of oCIS already with ownCloud 10 and also to have a smooth transition for users from ownCloud 10 to oCIS.

- ownCloud 10 with oCIS IDP
- Switch on the new front end "oCIS web" with ownCloud 10
- Run ownCloud 10 and oCIS in parallel - together
- Migrate users from ownCloud 10 to oCIS
