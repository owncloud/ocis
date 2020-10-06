---
title: "Deployment"
date: 2020-10-01T20:35:00+01:00
weight: -10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/deployment
geekdocFilePath: _index.md
---

{{< toc >}}

## Deployments Scenarios and Examples
This section handles deployments and operations for admins. If you are looking for a development setup, start with 

### Setup oCIS
oCIS deployments are super simple, yet there are many configrations possible for advanced setups.

- Basic setup - download and run
- Pick services and manage them individually
- SSL offloading with Traefik
- Use an external IDP

### Migrate an existing ownCloud 10
You can run ownCloud 10 and oCIS together. This allows you to use new parts of oCIS already with ownCloud 10 and also to have a smooth transition for users from ownCloud 10 to oCIS.

- ownCloud 10 with oCIS IDP
- Switch on the new front end "oCIS web" with ownCloud 10
- Run ownCloud 10 and oCIS in parallel - together
- Migrate users from ownCloud 10 to oCIS
