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

oCIS provides an internal identity provider which can be configured via the [IDP service](../../services/idp/), or connect to an external identity provider like Keycloak (in connection with openLDAP) or Microsoft Active Directory Federation Service (ADFS) (in connection with MS Active Directory). Entrada ID (formerly Azure AD) is in preperation, but not yet documented and might need some small fixes and for certain functions a LDAP/AD connection.
