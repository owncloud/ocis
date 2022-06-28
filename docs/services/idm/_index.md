---
title: IDM
date: 2022-03-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/idm
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The IDM service provides a minimal LDAP Service (based on https://github.com/libregraph/idm) for oCIS. It is started as part of
the default configuration and serves as a central place for storing user and group information.

It is mainly targeted at small oCIS installations. For larger setups it is recommended to replace IDM with a "real" LDAP server
or to switch to an external Identity Management Solution.

IDM listens on port 9325 by default. In the default configuration it only accepts TLS protected connections (LDAPS). The BaseDN
of the LDAP tree is `o=libregraph-idm`. IDM gives LDAP write permissions to a single user 
(DN: `uid=libregraph,ou=sysusers,o=libregraph-idm`). Any other authenticated user has read-only access. IDM stores its data in a
[boltdb](https://github.com/etcd-io/bbolt) file `idm/ocis.boltdb` inside the oCIS base data directory.

Note: IDM is limited in its functionality. It only supports a subset of the LDAP operations (namely BIND, SEARCH, ADD, MODIFY, DELETE).
Also IDM currently does not do any schema verification (e.g. structural vs. auxiliary object classes, require and option attributes,
syntax checks, ...). So it's not meant as a general purpose LDAP server.

## Table of Contents

{{< toc-tree >}}
