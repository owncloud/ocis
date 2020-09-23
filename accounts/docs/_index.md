---
title: Accounts
date: 2018-05-02T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis-accounts
geekdocEditPath: edit/master/docs
geekdocFilePath: _index.md
---

[![GitHub](https://img.shields.io/github/license/owncloud/ocis-hello)](https://github.com/owncloud/ocis-hello/blob/master/LICENSE)

## Abstract
OCIS needs to be able to identify users. Whithout a non reassignable and persistend account ID share metadata cannot be reliably persisted. `ocis-accounts` allows exchanging oidc claims for a uuid. Using a uuid allows users to change the login, mail or even openid connect provider without breaking any persisted metadata that might have been attached to it.

- persists accounts
- uses graph api properties
  -ldap can be synced using the onpremise* attributes

## Table of Contents

{{< toc-tree >}}