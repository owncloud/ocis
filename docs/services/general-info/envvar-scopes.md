---
title: Envvar Naming Scope
date: 2023-03-23T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info
geekdocFilePath: envvar-scopes.md
geekdocCollapseSection: true
---

The scope of an environment variable can be derived from its name. Therefore it is important to follow the correct naming scheme to enable easy and proper identification. This is important when either:

-   a new local envvar is introduced.
-   a new global envvar is added to an existing local envvar.

## Envvar Definition

-   A variable that is only used in a particular service is a **local envvar**.
-   A variable that is used in more than one service is a **global envvar**.
-   Mandatory when used in a service, a global envvar must have a local counterpart.
-   Variables that do not belong to any service are by definition global.

## Name Scope

### Local Envvars

A local envvar always starts with the the service name like `POSTPROCESSING_LOG_FILE`.

### Global Envvars

A global envvar always starts with `OCIS_` like `OCIS_LOG_FILE`.

Note that this envvar is the global representation of the local example from above.

To get a list of global envvars used in all services, see the [Global Environment Variables](https://doc.owncloud.com/ocis/next/deployment/services/env-vars-special-scope.html#global-environment-variables) table in the ocis admin documentation.

### Lifecycle

In the struct tag values of our config data types, we are using three key/value pairs to document the lifecycle of a config variable: `introductionVersion`, `deprecationVersion` and `removalVersion`. During the development cycle, a new value should set to `%%NEXT%%` as long as no release is scheduled. During the release process, the palceholder will be replaced with the actual version number. Our docs helper scripts will then automatically generate the correct documentation based on the version number.

## Deprecations

All environment variable types that are used in a service follow the same [deprecation rules]({{< ref "ocis/development/deprecating-variables/_index.md" >}}) independent of their scope.

## Separating Envvars

When multiple envvars are defined for one purpose like a global and local one, use `;` (semicolon) to properly separate the envvars in go code. Though it is possible to separate with `,` (comma) according go rules, the current implementation of the docs generation process only recognizes semicolons as separator.
