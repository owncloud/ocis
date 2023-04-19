---
title: IDP
date: 2023-04-19T15:29:23.072599054Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/idp
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

This service provides a builtin minimal OpenID Connect provider based on
[LibreGraph Connect (lico)](https://github.com/libregraph/lico) for oCIS.
It is mainly targeted at smaller installations. For larger setups it is
recommended to replace IDP with and external OpenID Connect Provider.
By default, it is configured to use the ocis IDM service as its LDAP backend for
looking up and authenticating users. Other backends like an external LDAP
server can be configured via a set of
[enviroment variables](https://owncloud.dev/services/idp/configuration/#environment-variables).

## Table of Contents

* [Example Yaml Config](#example-yaml-config)

## Example Yaml Config

{{< include file="services/_includes/idp-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/idp_configvars.md" >}}

