---
title: "Identity Provider"
date: 2023-05-03T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/identity-provider
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Overview

oCIS provides out of the box a minimal OpenID Connect provider via the [IDP service](../../services/idp/) and a minimal LDAP service via the [IDM service](../../services/idm/). Both services are limited in the provided functionality, see the [admin documentation](https://doc.owncloud.com/ocis/next/deployment/services/s-list/idp.html) for details, and can be used for small environments like up to a few hundred users. For enterprise environments, it is highly recommended using enterprise grade external software like KeyCloak plus openLDAP or MS ADFS with Active Directory, which can be configured in the respective service. Entrada ID (formerly Azure AD) is in preparation, but not yet released or documented and might need some small fixes, and for certain functions a LDAP/AD connection.
