---
title: "Continuous Integration"
date: 2020-10-01T20:35:00+01:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: continuous-integration.md
---

{{< toc >}}

## Overview

oCIS uses [GitHub Actions](https://github.com/owncloud/ocis/actions) as its CI system. Pipeline logs are visible directly in pull requests.

## Concepts

Pipelines are defined in `.github/workflows/`. The main acceptance test workflow is `.github/workflows/acceptance-tests.yml`.

Upon running the pipeline, tests run against the branch as-is. Make sure your branch is rebased onto the latest master before running CI to avoid false failures.

## Things done in CI

- static code analysis
- linting
- running UI tests
- running API acceptance tests
- build and release docker images
- build and release binaries
- build and release documentation

## Flags in commit message and PR title

You may add flags to your commit message or PR title to control pipeline runs.

- `[CI SKIP]`: no CI is run on the commit or PR
- `[k6-test]`: enabled the run of the k6 with smoke and performance test

### Knowledge base

- My pipeline fails because some CI related files or commands are missing.

  Make sure your branch is rebased onto the latest master. Workflow files in `.github/workflows/` change over time and a stale branch may reference steps or configs that no longer exist.

- How can I see what a workflow does?

  Workflow definitions are plain YAML in `.github/workflows/`. Open the relevant file directly in the repo — no CLI tool needed. The main acceptance test entry point is `.github/workflows/acceptance-tests.yml`.

- How can I re-run a failed job?

  In the GitHub Actions UI, click **Re-run jobs** → **Re-run failed jobs** on the failed workflow run. You can also trigger a full re-run from `gh`:

  ```bash
  gh run rerun <run-id> --failed --repo owncloud/ocis
  ```
