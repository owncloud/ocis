---
title: "Release Process for Envvars"
date: 2025-07-04T00:00:00+01:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info/envvars
geekdocFilePath: new-release-process.md
---

{{< toc >}}

**IMPORTANT**\
For a new ocis release, some tasks are necessary to be done **before** releasing. Follow the steps carefully to avoid issues. Most of the docs related tasks are not part of the CI. With each step finished successfully, the next step can be started. Sometimes, due to last minute changes, steps need to be redone!

The following can be done at any time but it must be done *latest* when no envvar changes are made which is just before a new release gets finally tagged. The data generated **must** be part of the upcoming release and be merged before tagging/branching!

## Special Scope Envvars

Ask the developers if envvars of this type have been changed (added or removed). See the [Special Envvars]({{< ref "./special-envvars.md#special-scope-envvars" >}}) documentation for more details on how to manage such a change.

## Extended Envvars

* From the ocis root run:\
`sudo make docs-clean`\
`make docs-generate`\
Drop any changes in `env_vars.yaml`!
* Check if there is a change in the `extended-envars.yaml` output.\
If so, process [Extended Envvars - Fixing Changed Item]({{< ref "./special-envvars.md#fixing-changed-items" >}}).
* When done, re-run `make docs-generate` and check if the output matches the expectations in `./docs/services/_includes/adoc/extended_configvars.adoc`.

## Ordinary Envvars

### Maintain the 'env_vars.yaml' File

This is **mandatory for a new release** !

* From the ocis root run:\
`sudo make docs-clean`\
`make docs-generate`\
Any changes in `env_vars.yaml` are now considered.
* This file will most likely show changes and merging them is **essential** as base for **added/removed or deprecated envvars**. Note that this file will get additions/updates only, but items never get deleted automatically !!\
{{< hint info >}}
Note that due to how the code is currently designed, **things may get shifted** around though no real changes have been introduced.
{{< /hint >}}
* First, check if any **alphabetic code names** are present in the changes. See [Introduce new Envvars]({{< ref "./envvar-naming-scopes.md/#introduce-new-envvars" >}}).
  * If so, create a new branch and replace them in the **service containing the source** with the actual semantic version (e.g. `releaseX` → `7.2.0`) first. Note that ALL of major, minor and patch numbers must be present, including patch versions == `0`.
  * If all changes are applied, rerun `make docs-generate` and check if all changes are incorporated in the yaml file.
  * Create a PR and merge these changes, dont forget to do a local pull of master afterwards...
* With a new branch, remove all envvars from the `env_vars.yaml` file manually that have formerly been deprecated and removed from the code.
* Commit the changes and merge it.\
Now `env_vars.yaml` is up to date on the repo in master, next steps are based on this state!

### Create Envvar Delta Files

* Create [Envvar Deltas Between Versions]({{< ref "./env-var-deltas/" >}}) files according the linked description.
