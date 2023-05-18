---
title: "Testing"
date: 2018-05-02T00:00:00+00:00
weight: 37
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: testing.md
---

{{< toc >}}

To run tests in the test suite you have two options. You may go the easy way and just run the test suite in docker. But for some tasks you could also need to install the test suite natively, which requires a little more setup since PHP and some dependencies need to be installed.

Both ways to run tests with the test suites are described here.

## Testing with test suite in docker

Let's see what is available. Invoke the following command from within the root of the oCIS repository.

```bash
make -C tests/acceptance/docker help
```

Basically we have two sources for feature tests and test suites:

- [oCIS feature test and test suites](https://github.com/owncloud/ocis/tree/master/tests/acceptance/features)
- [tests and test suites transferred from ownCloud, they have prefix coreApi](https://github.com/owncloud/ocis/tree/master/tests/acceptance/features)

At the moment both can be applied to oCIS since the api of oCIS is designed to be compatible with ownCloud.

As a storage backend, we offer oCIS native storage, also called "ocis". This stores files directly on disk. Along with that we also provide `S3` storage driver.

You can invoke two types of test suite runs:

- run a full test suite, which consists of multiple feature tests
- run a single feature or single scenario in a feature

### Run full test suite

#### Local oCIS tests (prefix `api`)

The names of the full test suite make targets have the same naming as in the CI pipeline. See the available local oCIS specific test suites [here](https://github.com/owncloud/ocis/tree/master/tests/acceptance/features). They can be run with `ocis` storage and `S3` storage.

For example, command:

```bash
make -C tests/acceptance/docker localApiTests-apiGraph-ocis
```

runs the same tests as the `localApiTests-apiGraph-ocis` CI pipeline, which runs the oCIS test suite "apiGraph" against the oCIS server with ocis storage.

And command:

```bash
make -C tests/acceptance/docker localApiTests-apiGraph-s3ng
```

runs the oCIS test suite `apiGraph` against the oCIS server with s3 storage.

{{< hint info >}}
While running the tests, oCIS server is started with [ociswrapper](https://github.com/owncloud/ocis/blob/master/tests/ociswrapper/README.md) (i.e. `WITH_WRAPPER=true`) by default. In order to run the tests without ociswrapper, provide `WITH_WRAPPER=false` when running the tests. For example:

```bash
WITH_WRAPPER=false \
BEHAT_FEATURE='tests/acceptance/features/apiGraph/createUser.feature:26' \
make -C tests/acceptance/docker test-ocis-feature-ocis-storage
```

But some test suites that are tagged with `@env-config` require the oCIS server to be run with ociswrapper. So, running those tests require `WITH_WRAPPER=true` (default setting).
{{< /hint >}}

#### Tests transferred from ownCloud core (prefix `coreApi`)

Command `make -C tests/acceptance/docker Core-API-Tests-ocis-storage-3` runs the same tests as the `Core-API-Tests-ocis-storage-3` CI pipeline, which runs the third (out of ten) test suites transferred from the ownCloud core against the oCIS server with ocis storage.

And `make -C tests/acceptance/docker Core-API-Tests-s3ng-storage-3` runs the third (out of ten) test suite transferred from the ownCloud core against the oCIS server with s3 storage.

### Run single feature test

A single feature tests (a feature file) can also be run against the different storage backends. To do that, multiple make targets with the schema test-\<test source\>-feature-\<storage-backend\> are available. To select a single feature you have to add an additional `BEHAT_FEATURE=<path-to-feature-file>` parameter when invoking the make command.

For example;

```bash
BEHAT_FEATURE='tests/acceptance/features/apiGraph/createUser.feature' \
make -C tests/acceptance/docker test-ocis-feature-ocis-storage
```

{{< hint info >}}
`BEHAT_FEATURE` must be pointing to a valid feature file
{{< /hint >}}

And to run a single scenario in a feature, you can do:

```bash
BEHAT_FEATURE='tests/acceptance/features/apiGraph/createUser.feature:26' \
make -C tests/acceptance/docker test-ocis-feature-ocis-storage
```

Similarly, with S3 storage;

```bash
# run a whole feature
BEHAT_FEATURE='tests/acceptance/features/apiGraph/createUser.feature' \
make -C tests/acceptance/docker test-ocis-feature-s3ng-storage

# run a single scenario
BEHAT_FEATURE='tests/acceptance/features/apiGraph/createUser.feature:26' \
make -C tests/acceptance/docker test-ocis-feature-s3ng-storage
```

In the same way, tests transferred from ownCloud core can be run as:

```bash
# run a whole feature
BEHAT_FEATURE='tests/acceptance/features/coreApiAuth/webDavAuth.feature' \
make -C tests/acceptance/docker test-core-feature-ocis-storage

# run a single scenario
BEHAT_FEATURE='tests/acceptance/features/coreApiAuth/webDavAuth.feature:13' \
make -C tests/acceptance/docker test-core-feature-ocis-storage
```

{{< hint info >}}
The tests suites transferred from ownCloud core have `coreApi` prefixed
{{< /hint >}}

### oCIS image to be tested (skip local image build)

By default, the tests will be run against the docker image built from your current working state of the oCIS repository. For some purposes it might also be handy to use an oCIS image from Docker Hub. Therefore, you can provide the optional flag `OCIS_IMAGE_TAG=...` which must contain an available docker tag of the [owncloud/ocis registry on Docker Hub](https://hub.docker.com/r/owncloud/ocis) (e.g. 'latest').

```bash
OCIS_IMAGE_TAG=latest \
make -C tests/acceptance/docker localApiTests-apiGraph-ocis
```

### Test log output

While a test is running or when it is finished, you can attach to the logs generated by the tests.

```bash
make -C tests/acceptance/docker show-test-logs
```

{{< hint info >}}
The log output is opened in `less`. You can navigate up and down with your cursors. By pressing "F" you can follow the latest line of the output.
{{< /hint >}}

### Cleanup

During testing we start a redis and oCIS docker container. These will not be stopped automatically. You can stop them with:

```bash
make -C tests/acceptance/docker clean
```

## Testing with test suite natively installed

We have two sets of tests:

- `test-acceptance-from-core-api` set was transferred from [core](https://github.com/owncloud/core) repository
  The suite name of all tests transferred from the core starts with "core"

- `test-acceptance-api` set was created for oCIS. Mainly for testing spaces features

### Run oCIS

Create an up-to-date oCIS binary by [building oCIS]({{< ref "build" >}})

To start oCIS:

```bash
IDM_ADMIN_PASSWORD=admin \
ocis/bin/ocis init --insecure true

OCIS_INSECURE=true PROXY_ENABLE_BASIC_AUTH=true \
ocis/bin/ocis server
```

`PROXY_ENABLE_BASIC_AUTH` will allow the acceptance tests to make requests against the provisioning api (and other endpoints) using basic auth.

#### Run local oCIS tests (prefix `api`)

```bash
make test-acceptance-api \
TEST_SERVER_URL=https://localhost:9200 \
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
```

#### Run tests transferred from ownCloud core (prefix `coreApi`)

```bash
make test-acceptance-from-core-api \
TEST_SERVER_URL=https://localhost:9200 \
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
```

Make sure to adjust the settings `TEST_SERVER_URL` according to your environment.

To run a single feature add `BEHAT_FEATURE=<feature file>`

example: `BEHAT_SUITE=tests/acceptance/features/apiGraph/createUser.feature`

To run a single test add `BEHAT_FEATURE=<file.feature:(line number)>`

example: `BEHAT_SUITE=tests/acceptance/features/apiGraph/createUser.feature:12`

To run a single suite add `BEHAT_SUITE=<test suite>`

example: `BEHAT_SUITE=apiGraph`

To run tests with a different storage driver set `STORAGE_DRIVER` to the correct value. It can be set to `OCIS` or `OWNCLOUD` and uses `OWNCLOUD` as the default value.

### Use existing tests for BDD

As a lot of scenarios from `test-acceptance-from-core-api` are written for oC10, we can use those tests for Behaviour driven development in oCIS.
Every scenario that does not work in oCIS with "ocis" storage, is listed in `tests/acceptance/expected-failures-API-on-OCIS-storage.md` with a link to the related issue.

Those scenarios are run in the ordinary acceptance test pipeline in CI. The scenarios that fail are checked against the
expected failures. If there are any differences then the CI pipeline fails.

The tests are not currently run in CI with the OWNCLOUD or EOS storage drivers, so there are no expected-failures files for those.

If you want to work on a specific issue

1. locally run each of the tests marked with that issue in the expected failures file.

   E.g.:

   ```bash
   make test-acceptance-from-core-api \
   TEST_SERVER_URL=https://localhost:9200 \
   TEST_OCIS=true \
   TEST_WITH_GRAPH_API=true \
   STORAGE_DRIVER=OCIS \
   BEHAT_FEATURE='tests/acceptance/features/coreApiVersions/fileVersions.feature:147'
   ```

2. the tests will fail, try to understand how and why they are failing
3. fix the code
4. go back to 1. and repeat till the tests are passing.
5. remove those tests from the expected failures file
6. make a PR that has the fixed code, and the relevant lines removed from the expected failures file.

## Running ENV config tests (@env-config)

Test suites tagged with `@env-config` are used to test the environment variables that are used to configure oCIS. These tests are special tests that require the oCIS server to be run using [ociswrapper](https://github.com/owncloud/ocis/blob/master/tests/ociswrapper/README.md).

### Run oCIS with ociswrapper

```bash
# working dir: ocis repo root dir

# init oCIS
IDM_ADMIN_PASSWORD=admin \
ocis/bin/ocis init --insecure true

# build the wrapper
cd tests/ociswrapper
make build

# run oCIS
PROXY_ENABLE_BASIC_AUTH=true \
./bin/ociswrapper serve --bin=../../ocis/bin/ocis
```

### Run the tests

```bash
OCIS_WRAPPER_URL=http://localhost:5200 \
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
TEST_SERVER_URL="https://localhost:9200" \
BEHAT_FEATURE=tests/acceptance/features/apiAsyncUpload/delayPostprocessing.feature \
make test-acceptance-api
```

### Writing new ENV config tests

While writing tests for a new oCIS ENV configuration, please make sure to follow these guidelines:

1. Tag the test suite (or test scenarios) with `@env-config`
2. Use `OcisConfigHelper.php` for helper functions - provides functions to reconfigure the running oCIS instance.
3. Recommended: add the new step implementations in `OcisConfigContext.php`

## Running test suite with email service (@email)

### Setup inbucket

Run the following command to setup inbucket

```bash
docker run -d --name inbucket -p 9000:9000 -p 2500:2500 -p 1100:1100 inbucket/inbucket
```

### Run oCIS

Documentation for environment variables is available [here](https://owncloud.dev/services/notifications/#environment-variables)

```bash
# init oCIS
IDM_ADMIN_PASSWORD=admin \
ocis/bin/ocis init --insecure true

# run oCIS
PROXY_ENABLE_BASIC_AUTH=true \
NOTIFICATIONS_SMTP_HOST=localhost \
NOTIFICATIONS_SMTP_PORT=2500 \
NOTIFICATIONS_SMTP_INSECURE=true \
ocis/bin/ocis server
```

### Run the acceptance test

Run the acceptance test with the following command:

```bash
make test-acceptance-api \
TEST_SERVER_URL="https://localhost:9200" \
TEST_OCIS=true \
TEST_WITH_GRAPH_API=true \
EMAIL_HOST="localhost" \
EMAIL_PORT=9000 \
BEHAT_FEATURE="tests/acceptance/features/apiEmailNotification/emailNotification.feature"
```

## Running tests for parallel deployment

### Setup the parallel deployment environment

Instruction on setup is available [here](https://owncloud.dev/ocis/deployment/oc10_ocis_parallel/#local-setup)

Edit the `.env` file and uncomment this line:

```bash
COMPOSE_FILE=docker-compose.yml:testing/docker-compose-additions.yml
```

Start the docker stack with the following command:

```bash
docker-compose up -d
```

### Getting the test helpers

All the test helpers are located in the core repo.

```bash
git clone https://github.com/owncloud/core.git
```

### Run the acceptance tests

Run the acceptance tests with the following command from the root of the oCIS repository:

```bash
make test-paralleldeployment-api \
TEST_SERVER_URL="https://cloud.owncloud.test" \
TEST_OC10_URL="http://localhost:8080" \
TEST_PARALLEL_DEPLOYMENT=true \
TEST_OCIS=true \
TEST_WITH_LDAP=true \
PATH_TO_CORE="<path_to_core>" \
SKELETON_DIR="<path_to_core>/apps/testing/data/apiSkeleton"
```

Replace `<path_to_core>` with the actual path to the root directory of core repo that you have cloned earlier.

In order to run a single test, use the `BEHAT_FEATURE` environment variable.

```bash
make test-paralleldeployment-api \
... \
BEHAT_FEATURE="tests/parallelDeployAcceptance/features/apiShareManagement/acceptShares.feature"
```
