---
title: "Releasing"
date: 2020-03-30T12:09:00+01:00
weight: 60
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: releasing.md
---

{{< toc >}}

## Prepare release
**Task:** Create release issue and add the task list. 

The ocis release process gets mainly coordinated through a release ticket in the ocis repository.
To get started please do so by creating a ticket that mentions the upcoming version number like "Release 0.0.0 Tech Preview" as title.

Usually the first issue comment keeps track of all later tasks and shows a task-list to keep track of them.

```markdown
## Release Tasks
- [ ] Brief marketing heads up
- [ ] Release Web [v0.0.0](https://github.com/owncloud/web/releases/tag/v0.0.0)
- [ ] Pin Web [v0.0.0](https://github.com/owncloud/web/releases/tag/v0.0.0) in ocis
- [ ] Pin Reva to [v0.0.0](https://github.com/cs3org/reva/releases/tag/v0.0.0) in ocis
- [ ] Write preliminary changelog
- [ ] Create release branch `release-0.0.0` (CODEFREEZE)
- [ ] Write release notes
- [ ] Prepare Changelog
- [ ] Create pre release tag `v0.0.0-rc1`
- [ ] Smoke test
- [ ] Performance test
- [ ] If needed, fix issues or write down known issues
- [ ] Create final signed tag `v0.0.0`
- [ ] Check successful CI run on `v0.0.0` tag
- [ ] Create release pr
- [ ] Ping marketing to update all download links
- [ ] Merge release pr to master
```

It's a common practice that every task has a competent who is responsible to keep care for this action item.
Let me explain what each item means and why it's there.

## Brief marketing heads up
**Task:** Ping marketing and coordinate timeline

To get everything well organized and give every party time to prepare, it's required to inform them before the release process starts.

## Pin dependent projects
**Task:** Pin and Update dependencies

Ocis depends on many external dependencies like [web](https://github.com/owncloud/web), [reva](https://github.com/cs3org/reva), ...
To get a stable release we need to take care that we pin those specific versions in the project. For go it's the go.mod file, fore node it's package.json. 
Please do not forget to also add individual lock files.

## Write preliminary changelog
**Task:** Collect changes and add them as a comment in the release ticket

The preliminary changelog is a good starting point to get a superficial overview of what is included in this release.
The easiest way of getting it, is to copy the "unreleased" list from ocis changelog from the master branch.
If you have external packages that should also be mentioned you need to copy their changes also.
Keep in mind, it's ok that the changelog is technical focused.

## Create release branch
**Task:** Create a release branch based on current master

At this point of time on we agreed the features that should be part of the next release.
To not block others work we create a dedicated release branch and do a CODEFREEZE.
This means, from now on for this releae, only hot-/Bug-fixes which are release blockers are allowed.
No new features!  

## Write release notes
**Task:** Update release notes in a separate branch and create a pr to release branch

The release notes are similar to the changelog but have a different audience.
These notes can maybe group many changes into one benefit that is described from the perspective of a user or admin.

## Prepare changelog
**Task:** Move all changelog items from changelog/unrelease to changelog/0.0.0_YEAR-MONTH-DAY

To get the ci able to create a final CHANGELOG.md its required to move all changelog items for this release from changelog/unrelease to changelog/0.0.0_YEAR-MONTH-DAY and commit the changes.
Afterward drone will keep care and creates the final changelog.md

## Create release candidate
**Task:** Create rc tag on the current release branch

After the changelog is there it's a common practice to create a release candidate tag like `v0.0.0-rc1` on the current release branch.
This could be used for testers or future reference.

## Smoke testing
**Task:** Click through the product and document result in the ticket 

Before we release a version we need to keep care that everything works as expected.
This is done manually, we do this to find cases which are maybe not covered by our ci.
After this is done, document the results in the release ticket, even if there are no failures it's a good practice to keep the test-plan as comment for later reference.

## Performance testing
**Task:** Run cdperf and document result in the ticket

At ownCloud, we decided to performance test our product which is even more important to be done right before the release.
It would go beyond the scope of this document to explain every step that is required to do so.
Please refer to [cdperf](https://github.com/owncloud/cdperf) for reference.   

## QA
**Task:** if there release blockers, fix them and repeat above steps

As said above, it's not allowed to introduce new features without stopping the release.
But if you have found any release blockers in your testing steps, it's required to keep care of them and repeat the steps before if required.
Please do not forget to bump the release candidate tag e.g. `v0.0.0-rc2` in this case.

## Prepare Release
**Task:** create release tag on release branch, wait for ci to be green, and create pr to master

After all steps are done and everything is fine it's time to create the fine release tag `v0.0.0`.
Then again wait for the ci if still everything works as expected and then create a release PR `release v0.0.0`.

## Notify Marketing
**Task:** notify marketing that everything is done and all links can be updated

Now it's time to notify marketing so that they can start to update all links, create a blog post and communicate the news.

## Notify Marketing
**Task:** close release pr and bring back to master

As a final step you need to take care that the PR gets back to master.
If there are any unexpected failures while merging you need to take of them. Please check twice if every tag is as it should.