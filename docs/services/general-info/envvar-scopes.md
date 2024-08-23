---
title: Envvar Naming Scopes
date: 2023-03-23T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info
geekdocFilePath: envvar-scopes.md
geekdocCollapseSection: true
---

{{< toc >}}

The scope of an environment variable can be derived from its name. Therefore, it is important to follow the correct naming scheme to enable easy and proper identification. This is important when either:

-   a new local envvar is introduced.
-   a new global envvar is added to an existing local envvar.

## Envvar Definition

-   A variable that is only used in a particular service is a **local envvar**.
-   A variable that is used in more than one service is a **global envvar**.
-   Mandatory when used in a service, a global envvar must have a local counterpart.
-   Variables that do not belong to any service are by definition global.

## Naming Scope

### Local Envvars

A local envvar always starts with the service name like `POSTPROCESSING_LOG_FILE`.

### Global Envvars

A global envvar always starts with `OCIS_` like `OCIS_LOG_FILE`.

Note that this envvar is the global representation of the local example from above.

To get a list of global envvars used in all services, see the [Global Environment Variables](https://doc.owncloud.com/ocis/next/deployment/services/env-vars-special-scope.html#global-environment-variables) table in the ocis admin documentation.

## Lifecycle of Envvars

The envvar struct tag contains at maximum the following key/value pairs to document the lifecycle of a config variable:

* `introductionVersion`
* `deprecationVersion`
* `removalVersion`
* `deprecationInfo`
* `deprecationReplacement`

### Introduce new Envvars

If a new envvar is introduced, only the `introductionVersion` is required.

{{< hint warning >}}
During the development cycle, the value for the `introductionVersion` must be set to `%%NEXT%%`. This placeholder will be removed by the real version number during the production releasing process. 
{{< /hint >}}

For the documentation to show the correct value for the `IV` (introduction version), our docs helper scripts will automatically generate the correct version to be printed in the documentation. If `%%NEXT%%` is found in the query, it will be replaced with `next`, else the value found is used.

During the releasing process for a production release, the placeholder `%%NEXT%%` has to be replaced with the new production version number like `%%NEXT%%` → `7.0.0`.

### Deprecate Existing Envvars

See the [deprecation rules]({{< ref "./deprecating-variables.md" >}}) documentation for more details.

## Separating Multiple Envvars

When multiple envvars are defined for one purpose like a global and local one, use `;` (semicolon) to properly separate the envvars in go code. Though it is possible to separate with `,` (comma) according go rules, the current implementation of the docs generation process only recognizes semicolons as separator.
