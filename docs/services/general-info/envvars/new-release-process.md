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
Discard any changes in `env_vars.yaml`!
* Check if there is a change in the `extended-envars.yaml` output.\
If so, process [Extended Envvars - Fixing Changed Items]({{< ref "./special-envvars.md#fixing-changed-items" >}}).
* When done, re-run `make docs-generate` and check if the output matches the expectations in `./docs/services/_includes/adoc/extended_configvars.adoc`.

## Ordinary Envvars

### Set the Correct IntroductionVersion

* Once the release is cut, **before** creating the first release candidate, replace them with the actual semantic version (e.g. `releaseX` → `8.1.0`). To find these placeholders in `introductionVersion` keys, you can run a helper script by issuing the following command:
  ```bash
  docs/ocis/helpers/identify_envvar_placeholder_names.sh
  ```

  {{< hint info >}}
  A new production version **MUST NOT** contain any alphabetic identifyers but the semantic version only, using **major, minor and a patch version, which is always 0!**.
  {{< /hint >}}

* Create a PR and merge it **before** taking the next step maintaining the `env_vars.yaml` file! Do not forget to rebase your local git repo.

### Maintain the 'env_vars.yaml' File

This is **mandatory for a new release** !

* From the ocis root run:\
`sudo make docs-clean`\
`make docs-generate`\
Any changes in `env_vars.yaml` are now considered.
* This file will most likely show changes and merging them is **essential** as base for **added/removed or deprecated envvars** (envvar deltas). Note that this file will get additions/updates only, but items never get deleted automatically !!\
{{< hint info >}}
Note that due to how the code is currently designed, **things may get shifted** around though no real changes have been introduced.
{{< /hint >}}
* With a new branch, remove all envvars from the `env_vars.yaml` file manually that have formerly been deprecated **and removed** from the code.
* Commit the changes and merge it.\
Now, `env_vars.yaml` is up to date in the repo in master. Next steps depend on this updated file!

### Create Envvar Delta Files

* Create [Envvar Deltas Between Versions]({{< ref "./env-var-deltas/" >}}) files according the linked description.
