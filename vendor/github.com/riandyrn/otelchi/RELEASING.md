# Release Process

This document explains how to create releases in this project in each release scenario.

Currently there are 2 release procedures for this project:

- [Latest Release](#latest-release)
- [Non-Latest Release](#non-latest-release)

## Latest Release

This procedure is used when we want to create new release from the latest development edge (latest commit in the `master` branch).

The steps for this procedure are the following:

1. Create a new branch from the `master` branch with the following name format: `pre_release/v{MAJOR}.{MINOR}.{BUILD}`. For example, if we want to release for version `0.6.0`, we will first create a new branch from the `master` called `pre_release/v0.6.0`.
2. Update method `Version()` in `version.go` to return the target version.
3. Update the `CHANGELOG.md` to include all the notable changes.
4. Create a new PR from this branch to `master` with the title: `Release v{MAJOR}.{MINOR}.{BUILD}` (e.g `Release v0.6.0`).
5. At least one maintainer should approve the PR. However if the PR is created by the repo owner, it doesn't need to get approval from other maintainers.
6. Upon approval, the PR will be merged to `master` and the branch will be deleted.
7. Create new release from the `master` branch.
8. Set the title to `Release v{MAJOR}.{MINOR}.{BUILD}` (e.g `Release v0.6.0`).
9. Set the newly release tag using this format: `v{MAJOR}.{MINOR}.{BUILD}` (e.g `v0.6.0`).
10. Set the description of the release to match with the content inside `CHANGELOG.md`.
11. Set the release as the latest release.
12. Publish the release.
13. Done.

## Non-Latest Release

This procedure is used when we need to create fix or patch for the older releases. Consider the following scenario:

1. For example our latest release is version `0.7.1` which has the minimum go version `1.18`.
2. Let say our user got a critical bug in version `0.6.0` which has the minimum go version `1.15`.
3. Due to some constraints, this user cannot upgrade his/her minimum go version.
4. We decided to create fix for this version by releasing `0.6.1`.

In this scenario, the procedure is the following:

1. Create a new branch from the version that we want to patch. 
2. We name the new branch with the increment in the build value. So for example if we want to create patch for `0.6.0`, then we should create new branch with name: `patch_release/v0.6.1`.
3. We create again new branch that will use `patch_release/v0.6.1` as base. Let say `fix/handle-cve-233`.
4. We will push any necessary changes to `fix/handle-cve-233`.
5. Create a new PR that target `patch_release/v0.6.1`.
6. Follow step `2-10` as described in the [Latest Release](#latest-release).
7. Publish the release without setting the release as the latest release.
8. Done.