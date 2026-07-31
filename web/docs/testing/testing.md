---
title: 'Running Tests'
date: 2021-07-27T00:00:00+00:00
weight: 60
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/testing
geekdocFilePath: testing.md
---

{{< toc >}}

## Introduction

In order to allow us to make changes quickly, often and with a high level of confidence, we heavily rely on tests within the `web` repository.

All the steps below require you to have the `web` repo cloned locally and dependencies installed.
This can be achieved by running

```shell
$ git clone https://github.com/owncloud/web.git
$ cd web
$ pnpm install
```

## Unit Tests

We have a steadily growing coverage of unit tests. You can run them locally via

```shell
$ pnpm test:unit
$ pnpm -r test:unit
```

You can also specify which package to run the test on, such as: `pnpm --filter @ownclouders/web-pkg test:unit`.
Alternatively, tests can also be run by navigating to the package name and then running `pnpm test:unit`.

### Unit Test File Structure

Our unit tests spec files follow a simple structure:

- fixtures and mocks at the top
- helper functions at the bottom
- tests in between

We usually organize tests with nested `describe` blocks. If you would like to get feedback from the core team about
the structure, scope and goals of your unit tests before actually writing some, we invite you to make a pull request
with only `describe` blocks and nested `it.todo("put your test description here")` lines.

## E2E Tests (Playwright)

Our end-to-end test suite is built upon the [Playwright Framework](https://github.com/microsoft/playwright),
which makes it easy to write tests, debug them and have them run cross-browser with minimal overhead.

### Preparation

Please make sure you have installed all dependencies and started the server(s) as described in [tooling]({{< ref "tooling.md#development-setup" >}}).

### Prepare Web

Bundle the web frontend with the following command:

```shell
$ pnpm build
```

Our compose setup automatically mounts it into an oCIS backend, respectively. Web also gets recompiled on changes.

### Start Web

Start the web with the following command:

```shell
docker compose up
```

This will start all the services. The ENV variables specific to each services are defined in the `docker-compose.yml` file.

### Run E2E Tests

The following command will run all available e2e tests:

```shell
$ pnpm test:e2e:playwright 'tests/e2e/specs/**/*.spec.ts'
```

### Options

To run a particular test, simply add the spec file and line number to the test command, e.g. `pnpm test:e2e:playwright tests/e2e/specs/admin-settings/users.feature:14`

Various options are available via ENV variables, e.g.

- `BASIC_AUTH=true` use basic authorization for api requests.
- `RETRY=n` to retry failures `n` times
- `SLOW_MO=n` to slow the execution time by `n` milliseconds
- `TIMEOUT=n` to set tests to timeout after `n` milliseconds
- `HEADLESS=bool` to open the browser while the tests run (defaults to true => headless mode)
- `BROWSER=name` to run tests against a specific browser. Defaults to chromium, available are chromium, firefox, webkit, chromium
- `ADMIN_PASSWORD` to set administrator password. By default, the `admin` password is used in the test

For debugging reasons, you may want to record a video or traces of your test run.
Again, you can use the following ENV variables in your command:

- `REPORT_DIR=another/path` to set a directory for your recorded files (defaults to "reports")
- `REPORT_VIDEO=true` to record a video of the test run
- `REPORT_HAR=true` to save request information from the test run
- `REPORT_TRACING=true` to record traces from the test run

To then open e.g. the tracing from the `REPORT_DIR`, run

```shell
$ npx playwright show-trace path/to/file.zip
```

### Lint E2E Test Code

Run the following command to find out the lint issues early in the test codes:

```shell
$ pnpm lint
```

And to fix the lint problems run the following command:

```shell
$ pnpm lint --fix
```

If the lint problems are not fixed by `--fix` option, we have to manually fix the code.

## Analyze the Test Report

After running tests, report artifacts are written under `REPORT_DIR` (defaults to `reports/e2e`).

- Accessibility report: `reports/e2e/a11y-report.json`
- Traces (when `REPORT_TRACING=true`): `reports/e2e/playwright/tracing/*.zip`
- Optional videos/HAR files (when enabled via `REPORT_VIDEO=true` or `REPORT_HAR=true`) are also stored in the report directory.

To inspect a trace file:

```bash
npx playwright show-trace reports/e2e/playwright/tracing/<trace-file>.zip
```

If you want an HTML Playwright report for a run, execute tests with the HTML reporter enabled and then open it:

```bash
pnpm test:e2e:playwright -- --reporter=html
npx playwright show-report
```

## E2E Tests on oCIS With Keycloak

We can run some of the e2e tests on oCIS setup with Keycloak as an external idp. To run tests against locally, please follow the steps below:

### Run oCIS With Keycloak

There's a documentation to serve [oCIS with Keycloak](https://owncloud.dev/ocis/deployment/ocis_keycloak/). Please follow each step to run **oCIS with Keycloak**.

### Run E2E Tests

```bash
KEYCLOAK=true \
BASE_URL_OCIS=ocis.owncloud.test \
pnpm run test:e2e:playwright tests/e2e/specs/journeys
```

Following environment variables come in use while running e2e tests on oCIS with Keycloak:

- `BASE_URL_OCIS` sets oCIS url (e.g.: ocis.owncloud.test)
- `KEYCLOAK_HOST` sets Keycloak url (e.g.: keycloak.owncloud.test)
- `KEYCLOAK=true` runs the tests with Keycloak
- `KEYCLOAK_REALM` sets oCIS realm name used on Keycloak

## E2E Tests With Predefiend Users (`@predefined-users`)

It is possible to run e2e tests with predefined users. This is useful for running tests in a production-like environment.
The following environment variables are used to run the tests with predefined users:

- `PREDEFINED_USERS`: `true`|`false`
- `PREDEFINED_USERS_FILE`: path to a JSON file mapping predefined users

We have to create a JSON file that contains the mapping of the users. JSON file MUST contain the following keys:

```json
{
 "alice": {// map user},
 "brian": {// mapuser},
 "carol": {// mapuser},
}
```

And the user object MUST have the following properties defined:

```json
{
  "id": "<usernmae>",
  "displayName": "<display-name>",
  "password": "<password>",
  "email": "<email>"
}
```

A complete example of a JSON file.

```json
{
  "alice": {
    "id": "einstein",
    "displayName": "Albert Einstein",
    "password": "relativity",
    "email": "einstein@example.org"
  },
  "brian": {
    "id": "marie",
    "displayName": "Marie Skłodowska Curie",
    "password": "radioactivity",
    "email": "marie@example.org"
  },
  "carol": {
    "id": "moss",
    "displayName": "Maurice Moss",
    "password": "vista",
    "email": "moss@example.org"
  }
}
```

The test scenarios that can run with predefined users are marked with the `@predefined-users` tag and can be run with the following command:

```bash
PREDEFINED_USERS=true \
PREDEFINED_USERS_FILE='<path-to>/users.json' \
pnpm test:e2e:playwright tests/e2e/specs/file-action/rename.feature --tags '@predefined-users'
```

**The following tests cannot be run with predefined users:**

All tests which are related to:

- Admin Actions
- Groups

**The tests might show flakiness or fail due to the following reasons:**

- Slower network connection
- Features enabled/disabled
- Running latest tests against an older version of oCIS/Web
- Large file uploads may take longer time

## Usage of `web-packages.txt` In the Test Suite

Test suites may include the `web-packages.txt` file to denote which web packages changes affect the defined test scenarios. This information is used in CI pipelines to determine which test suites to run based on the changed web packages.

The `web-packages.txt` file should be included within the test suite directory as shown below:

```
└── tests/e2e/specs
    └── admin-settings
        ├── users.feature
        └── web-packages.txt
```

And the `web-packages.txt` file should list the dependent web packages, one per line, for example:

NOTE: The package name should start with `web-` in order to be recognized correctly, if not, the line will be ignored.

```
web-app-files
web-app-admin-settings
```
