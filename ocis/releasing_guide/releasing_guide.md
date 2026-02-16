---
title: "Releasing Guide"
date: 2020-12-16T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis/releasing_guide
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: releasing_guide.md
---

{{< toc >}}

To prepare an oCIS production release, you must follow a specific sequence of steps in the correct order.

## New ocis Version

Use this detailed copy/paste checklists for Jira when starting a new release.

{{< hint warning >}}
* The **examples** use the following versions which will need to be **adapted** for the planned release:
  * For Major and Minor versions: `7.2.0`
  * For Patch releases: `7.2.1`
* The process differs slighly when creating a patch release only.
* Changes applied to a `stable-7.2` branch valid for documentation will be processed to the `docs-stable-x.y` branch automatically after merging.
{{< /hint >}}

### Steps for Major or Minor Releases

{{< hint warning >}}
* A new production release is being prepared and the implementation of new features has finished in master.
* `web`, `reva` and `dependencies` have already been bumped.
* Translations
  * For oCIS, all required translations changes are included.
  * For Web, all translation changes have been applied and are part of the web version provided.
* **Only** bug fixes to release candidates are allowed.
{{< /hint >}}

#### Overview

1. Start the Releasing Process
1. Start the Release Candidates Phase
1. Prepare Release\
Integrate all RC changes and set the final version
1. Sign-off the Releasing Process
1. Release the new Version
1. Postprocessing and Finalization\
Backport changes to `master` and create other required stable branches, update cronjobs and clean up

#### Details

```
https://owncloud.dev/ocis/releasing_guide/

* [ ] Start the Releasing Process
  * [ ] Create a `stable-7.2` branch based on current `master`
  * [ ] Inform documentation that the releasing process has started

* [ ] Start the Release Candidates Phase
  * [ ] Create a new feature branch based on `stable-7.2`
  * [ ] Move all changelog items from `unreleased` to `7.2.0-rc.1_2025-06-12`
  * [ ] Bump oCIS version in `ocis-pkg/version/version.go`
  * [ ] Bump oCIS version in `sonar-project.properties`
  * [ ] Create PR with `[full-ci][k6-test]` against `stable-7.2`
  * [ ] Get PR approved. **DO NOT MERGE YET**
  * [ ] Wait for pipeline to be green
  * [ ] Create tag: `git tag -s v7.2.0-rc.1 -m "Release 7.2.0-rc.1"`
  * [ ] Push tag: `git push origin v7.2.0-rc.1`
  * [ ] Watch the PR to see pipeline succeed (can be restarted)
  * [ ] Merge PR
  * [ ] Sync with DevOps and Product\
  Repeat process with `rc.2`,`rc.3`, ...
  * [ ] All required translations from oCIS are included
  * [ ] All issues are fixed, RC phase has finished

* [ ] Prepare Release
  * [ ] Create a new feature branch based on `stable-7.2`
  * [ ] Move all changelog items from `7.2.0-rc.*` to `7.2.0_2025-04-01`
  * [ ] Bump oCIS version in `ocis-pkg/version/version.go`
  * [ ] Bump oCIS version in `sonar-project.properties`
  * [ ] Create PR with `[full-ci][k6-test]` against `stable-7.2`
  * [ ] Mark PR as **Draft** to avoid accidentially merging
  * [ ] Get PR approved. **DO NOT MERGE YET**\
  Info: merging will be done in step *Release the new Version*
  * [ ] Wait for pipeline to be green

* [ ] Get new Release Sign-off (jira)
  * [ ] **EITHER (preferred):** Find someone who wants the release more than you do, and have them sign-off
  * [ ] **OR (not recommended):** Have the appropriate people sign the *release sign-off* document

* [ ] Release the new Version
  * [ ] Create tag: `git tag -s v7.2.0 -m "Release 7.2.0"`\
  Note the tag name scheme is important
  * [ ] Push tag: `git push origin v7.2.0`
  * [ ] Watch the PR to see pipeline succeed (can be restarted)
  * [ ] Smoke test docker image `owncloud/ocis@v7.2.0`
    * [ ] Choose any docker-compose example from oCIS repository
    * [ ] Export `OCIS_DOCKER_IMAGE=owncloud/ocis`
    * [ ] Export `OCIS_DOCKER_TAG=7.2.0`
    * [ ] `docker compose up -d`
    * [ ] Confirm oCIS version in browser and start the `upload-download-awesome` test
  * [ ] Remove the **Draft** state and Merge PR\
  This is the PR from step *Prepare Release*
    * [ ] Delete the feature branch
  * [ ] Announce release in *Teams channel: oCIS*

* [ ] Postprocessing and Finalization
  * [ ] Backport `stable-7.2` to the master branch
  * [ ] Create a `docs-stable-7.2` branch
    * [ ] Create orphan branch: `git checkout --orphan docs-stable-7.2`
    * [ ] Initial commit: `git commit --allow-empty -m "initial commit"`
    * [ ] Push it: `git push`
  * [ ] Adjust the `.drone.star` file to write to `docs-stable-7.2`
    * [ ] Find target_branch value in the docs section and change it to `docs-stable-7.2`
    * [ ] Example: https://github.com/owncloud/ocis/blame/56f7645f0b11c9112e15ce46f6effd2fea01d6be/.drone.star#L2249
  * [ ] Add `stable-7.2` to the nightly cron jobs in drone (`Settings` -> `Cron Jobs`)
```

### Steps for Patch Releases Only

{{< hint warning >}}
* A patch branch is prepared, based off the appropriate stable branch, and contains all the changes including `web`, `reva` and `dependencies` bumps. The patch branch is merged into the corresponding stable branch. No release candidates are used because the changes are known.
* Translations
  * For oCIS, all required translations changes are included.
  * For Web, all translation changes have been applied and are part of the web version provided.
* **Only** bug fixes to the patch branch are allowed.
{{< /hint >}}

#### Overview

1. Prepare Release\
Integrate all patches and set the final version
1. Sign-off the Releasing Process
1. Release the new Version
1. Check Forward Porting to Master

#### Details

```
https://owncloud.dev/ocis/releasing_guide/

* [ ] Prepare Release
  * [ ] Create a new feature branch based on `stable-7.2`
  * [ ] Move all changelog items from `unreleased` to `7.2.1_2025-04-01`
  * [ ] Bump oCIS version in `ocis-pkg/version/version.go`
  * [ ] Bump oCIS version in `sonar-project.properties`
  * [ ] Create PR with `[full-ci][k6-test]` against `stable-7.2`
  * [ ] Mark PR as **Draft** to avoid accidentially merging
  * [ ] Get PR approved **requires 2 approvals**, **DO NOT MERGE YET**\
  Info: merging will be done in step *Release the new Version*
  * [ ] Wait for pipeline to be green

* [ ] Get new Release Sign-off (confluence)
  * [ ] **EITHER (preferred):** Find someone who wants the release more than you do, and have them sign-off
  * [ ] **OR (not recommended):** Have the appropriate people sign the *release sign-off* document

* [ ] Release the new Version
  * [ ] Create tag: `git tag -s v7.2.1 -m "Release 7.2.1"`\
  Note the tag name scheme is important
  * [ ] Push tag: `git push origin v7.2.1`
  * [ ] Watch the PR to see pipeline succeed (can be restarted)
  * [ ] Smoke test docker image `owncloud/ocis@v7.2.1`
    * [ ] Choose any docker-compose example from ocis repository
    * [ ] Export `OCIS_DOCKER_IMAGE=owncloud/ocis`
    * [ ] Export `OCIS_DOCKER_TAG=7.2.1`
    * [ ] `docker compose up -d`
    * [ ] Confirm oCIS version in browser and start the `upload-download-awesome` test
  * [ ] Remove the **Draft** state and Merge PR\
  This is the PR from step *Prepare Release*
    * [ ] Delete the feature branch
  * [ ] Announce release in *Teams channel: oCIS*

* [ ] Check if Forward Porting to Master is Required
```

## Update Documentation Related Data

### Envvar Changes

Follow the [Release Process for Envvars](https://owncloud.dev/services/general-info/envvars/new-release-process/) documentation to update the required environment variable data. These changes **must** be present in the appropriate `stable-7.2` branch and should be applied to the respective feature branch during the development cycle. This **avoids manually backporting** to `stable-7.2` but may be required if changes go to master instead. Changes to envvars are typically found in the `docs` folder.

{{< hint warning >}}
The admin docs processes access the data via the branch name, not the tag set. Therefore, any changes applied to the stable branch are accessible to the documentation process without the need of a new version.
{{< /hint >}}

### Prepare Admin Docs

The admin documentation must be prepared for the new release.

Note that this section is only **informational** for oCIS developers and is done by the documentation team. This is why they need to be informed in a timely manner.

* [ ] Create **overview release notes** in the `docs-main` repo based on the changelog\
The source is most likely the unreleased folder during the development cycle
* [ ] For Major and Minor versions:
  * [ ] Add all **documentation changes** including the upgrade guide in the `docs-ocis` repo in master
  * [ ] Create a new release branch in the `docs-ocis` repo based off master
  * [ ] **Update oCIS branches and versions** in the `docs` repo
* [ ] For Patch releases:
  * [ ] **Update oCIS versions** in the `docs` repo
  
