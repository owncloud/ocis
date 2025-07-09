---
title: Envvar Naming Scopes
date: 2023-03-23T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info/envvars
geekdocFilePath: envvar-naming-scopes.md
geekdocCollapseSection: true
---

{{< toc >}}

The scope of an environment variable can be derived from its name. Therefore, it is important to follow the correct naming scheme to enable easy and proper identification. This is important when either:

-   a new local envvar is introduced.
-   a new global envvar is added to an existing local envvar.

## Envvar Definition

-   A variable that is only used in a particular service is a **local envvar**.
-   A variable that is used in more than one service is a **global envvar**.
-   If applicapable, a global envvar has a local counterpart.
-   Variables that are not limited to any service are by definition global.

## Naming Scope

### Local Envvars

A local envvar always starts with the service name such as `POSTPROCESSING_LOG_FILE`.

### Global Envvars

A global envvar always starts with `OCIS_` like `OCIS_LOG_FILE`.

Note that this envvar is the global representation of the local example from above.

To get a list of global envvars used in all services, see the [Global Environment Variables](https://doc.owncloud.com/ocis/next/deployment/services/env-vars-special-scope.html#global-environment-variables) table in the ocis admin documentation.

### Reserved Envvar Names

Services and their local envvars **MUST NOT** be named `extended` or `global`. These are reserved names for the automated documentation process.

## Lifecycle of Envvars

The envvar struct tag contains at maximum the following key/value pairs to document the lifecycle of a config variable:

* `introductionVersion`
* `deprecationVersion`
* `removalVersion`
* `deprecationInfo`
* `deprecationReplacement`

### Introduce new Envvars

If a new envvar is introduced, the complete struct needs to be added, but only the `introductionVersion` requires data.

{{< hint info >}}
* During development, set `introductionVersion` to a short, **alphabetic code name** that represents the upcoming release (e.g. `releaseX`).
* This identifier stays constant until the release receives its final production semantic-version number.
{{< /hint >}}

The docs helper scripts render these alphabetic identifiers verbatim. They appear in the next (master) branch of the admin docs as rendered here.

Once the release is cut, before tagging, replace them with the actual semantic version (e.g. `releaseX` → `7.2.0`). 

{{< hint info >}}
A new production version **MUST NOT** contain any alphabetic identifyers but the semantic version only.
{{< /hint >}}

### Adding Envvars to Existing Ones

If an envvar has been introduced with a particular version, the `introductionVersion` gets a value accordingly. If an additional envvar like a global one is added to this existing envvar later on, the introduction version will *not* be changed.

### Deprecate Existing Envvars

See the [deprecation rules]({{< ref "./deprecating-variables.md" >}}) documentation for more details.

## Separating Multiple Envvars

When multiple envvars are defined for one purpose like a global and local one, use `;` (semicolon) to properly separate the envvars in go code. Though it is possible to separate with `,` (comma) according go rules, the current implementation of the docs generation process only recognizes semicolons as separator.
