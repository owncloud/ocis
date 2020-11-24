---
title: "Continuous Integration"
date: 2020-10-01T20:35:00+01:00
weight: 100
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: continuous-integration.md
---

{{< toc >}}

oCIS uses [DRONE](https://www.drone.io/) as CI system. You can find the pipeline logs [here](https://drone.owncloud.com/owncloud/ocis) or in your PR.

## Concepts

The pipeline is defined in [Starlark](https://github.com/bazelbuild/starlark) and transformed to YAML upon pipeline run. This enables us to do a highly dynamic and non repeating pipeline configuration.

Upon running the pipeline, your branch gets merged to the master branch. This ensures that we always test your changeset if as it was applied to the master of oCIS. Please note that this does not apply to the pipeline definition (`.drone.star`).

## Things done in CI

- static code analysis
- linting
- running UI tests
- running ownCloud 10 test suite against oCIS
- build and release docker images
- build and release binaries
- build and release documentation

## Flags in commit message and PR title

You may add flags to your commit message or PR title in order to speed up pipeline runs and take load from the CI runners.

- `[CI SKIP]`: no ci is run on the commit or PR

- `[docs-only]`: please add this flag, if you only changed documentation. This will only trigger documentation related CI steps.

### Knowledge base

- My pipeline fails because some CI related files or commands are missing.

  Please make sure to rebase your branch onto the lastest master of oCIS. It could be that the pipeline definition (`.drone.star`) was changed on the master branch. This is is the only file, that will not be auto merged to master upon pipeline run. So things could be out of sync.

- How can I see the YAML drone pipeline definition?

  In order to see the Yaml pipeline definition you can use the drone-cli to convert the Starlark file.

  ```
  drone starlark
  ```

  {{< hint info >}}
  If you experience a `"build" struct has no .title attribute` error you need a patched drone-cli binary.
  You need to build it yourself from this [source code](https://github.com/micbar/drone-cli/tree/add-more-flags).
  (There is also an open [PR](https://github.com/drone/drone-cli/pull/175) for that on drone-cli)
  {{< /hint >}}
