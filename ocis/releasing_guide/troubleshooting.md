---
title: "Troubleshooting"
date: 2020-12-16T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis/releasing_guide
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: troubleshooting.md
---

{{< toc >}}

This document covers some of the issues that can happen during the release process.

Use this detailed copy/paste checklists for the PR when starting a backport.

## Backport Missing Commits

Once the `stable-x.y` branch has been created but work has been merged into master, these changes will not appear in the stable branch unless backporting is used.

{{< hint warning >}}
* Any change applied to `stable-7.2` must be made **before** the tag for this release is created to be part of which is mandatory for any code change.
* There is one exception. Changes applied for documentation purposes are accessed via the branch name by the doc process and not via the tag and therefore accessible to the documentation processes.
{{< /hint >}}

Though any changes should ideally be targeted to the correct branch from the beginning, this may sometimes be impossible, especially for the following folder locations:

- `docs` and
- `deployments/examples/ocis_full`

### Documentation

Changes to documentation must be present in the respective stable branch because the documentation process accesses the data for further processing. Follow these steps to apply all doc related changes from master to the stable branch **after** the tag has been created:

```
https://owncloud.dev/ocis/releasing_guide/

* [ ] Check all relevant commits made to the master branch that are not in stable\
`git log  stable-7.2..master -- docs/`
* [ ] Double check all changes if they apply to stable\
`git diff  stable-7.2..master -- docs/`
* [ ] Create a new feature branch based on `stable-7.2`
* [ ] Cherry pick missing commits
* [ ] Create PR with leading `[docs-only]` against `stable-7.2`
* [ ] Get PR approved **requires 2 approvals** and merge it
* [ ] Delete feature branch
```

### Deployment Example

Changes to the `ocis_full` deployment example are less critical. But as this is referenced by the documentation, changes must be present in the respective stable branch. Follow these steps to apply all deployment example related changes from master to the stable branch **after** the tag has been created:

#### Web Extensions

Update all dependencies and extension related data in the `web-extensions` repo and create for each extension a new release if applicapable.

#### ocis_full

Update the `ocis_full` deployment example as required in master which includes to update the versions of the web extensions used. Test the functionality if the example works.

Follow these steps to apply all `ocis_full` deployment example related changes from master to the stable branch:

```
https://owncloud.dev/ocis/releasing_guide/

* [ ] Check all relevant commits made to the master branch that are not in stable\
`git log  stable-7.2..master -- deployments/examples/ocis_full/`
* [ ] Double check all changes if they apply to stable\
`git diff  stable-7.2..master -- deployments/examples/ocis_full/`
* [ ] Create a new feature branch based on `stable-7.2`
* [ ] Cherry pick missing commits
* [ ] Create PR with leading `[docs-only]` against `stable-7.2`
* [ ] Get PR approved **requires 2 approvals** and merge it
* [ ] Delete feature branch
```
