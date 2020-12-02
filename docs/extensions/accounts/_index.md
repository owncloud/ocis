---
title: Accounts
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/accounts
geekdocFilePath: _index.md
---

## Abstract
oCIS needs to be able to identify users. Without a non reassignable and persistent account ID share metadata cannot be reliably persisted. `accounts` allows exchanging oidc claims for a uuid. Using a uuid allows users to change the login, mail or even openid connect provider without breaking any persisted metadata that might have been attached to it.

- persists accounts
- uses graph api properties
- ldap can be synced using the onpremise* attributes

## Table of Contents

{{< toc-tree >}}