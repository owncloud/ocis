---
title: 'Accessibility'
date: 2025-11-10T00:00:00+00:00
weight: 60
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/testing
geekdocFilePath: accessibility.md
---

{{< toc >}}

## Introduction

Accessibility is a crucial aspect of web development. It ensures that web applications are usable by everyone, including people with disabilities.

## What Tools We Use

### eslint-plugin-vuejs-accessibility

We use [eslint-plugin-vuejs-accessibility](https://github.com/vue-a11y/eslint-plugin-vuejs-accessibility) to quickly catch accessibility issues in the codebase. This plugin is used in the [@ownclouders/eslint-config](https://www.npmjs.com/package/@ownclouders/eslint-config).

### @axe-core/playwright

We use [@axe-core/playwright](https://github.com/dequelabs/axe-core-npm) to automatically test the accessibility of the ownCloud Web client. All tests are run automatically on every PR and on commits to the `master` branch. We are not running dedicated accessibility tests and instead make them part of the E2E tests. This way we do not have to maintain duplicate tests and we can granularly add accessibility tests to the specific steps of the tests. The tests are considered failed if any `serious` or `critical` accessibility violations are found.

#### Running Accessibility Tests

To run the accessibility tests, you can simply run our existing E2E tests using the following command:

```bash
pnpm test:e2e:playwright
```

#### Skipping Accessibility Tests Locally

If you want to skip the accessibility tests, you can add the `SKIP_A11Y_TESTS` environment variable to your command.

```bash
SKIP_A11Y_TESTS=true pnpm test:e2e:playwright
```

#### Skipping Accessibility Tests in CI

If you want to skip the accessibility tests in CI, you can add the `[skip-a11y]` flag into the title of the PR.

#### Accessibility Report

After the tests are run, a JSON accessibility report is generated in the `reports/e2e/a11y-report.json` file. This report contains detailed information about the accessibility violations found in the tests.
