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

## Testing With Test Suite in Docker

Let's see what is available. Invoke the following command from within the root of the oCIS repository.

```bash
make -C tests/acceptance/docker help
```

Basically we have two sources for feature tests and test suites:

- [oCIS feature test and test suites](https://github.com/owncloud/ocis/tree/master/tests/acceptance/features)
- [tests and test suites transferred from ownCloud core, they have prefix coreApi](https://github.com/owncloud/ocis/tree/master/tests/acceptance/features)

At the moment, both can be applied to oCIS since the api of oCIS is designed to be compatible with ownCloud.

As a storage backend, we offer oCIS native storage, also called `ocis`. This stores files directly on disk. Along with that we also provide `S3` storage driver.

You can invoke two types of test suite runs:

- run a full test suite, which consists of multiple feature tests
- run a single feature or single scenario in a feature

### Run Full Test Suite

#### Local oCIS Tests (prefix `api`)

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

{{< hint info >}}
To run the tests that require an email server (tests tagged with `@email`), you need to provide `START_EMAIL=true` while running the tests.

```bash
START_EMAIL=true \
BEHAT_FEATURE='tests/acceptance/features/apiEmailNotification/emailNotification.feature' \
make -C tests/acceptance/docker test-ocis-feature-ocis-storage
```

{{< /hint >}}

#### Tests Transferred From ownCloud Core (prefix `coreApi`)

Command `make -C tests/acceptance/docker Core-API-Tests-ocis-storage-3` runs the same tests as the `Core-API-Tests-ocis-storage-3` CI pipeline, which runs the third (out of ten) test suite groups transferred from ownCloud core against the oCIS server with ocis storage.

And `make -C tests/acceptance/docker Core-API-Tests-s3ng-storage-3` runs the third (out of ten) test suite groups transferred from ownCloud core against the oCIS server with s3 storage.

### Run Single Feature Test

The tests for a single feature (a feature file) can also be run against the different storage backends. To do that, multiple make targets with the schema **test-_\<test-source\>_-feature-_\<storage-backend\>_** are available. To select a single feature you have to add an additional `BEHAT_FEATURE=<path-to-feature-file>` parameter when invoking the make command.

For example;

```bash
BEHAT_FEATURE='tests/acceptance/features/apiGraph/createUser.feature' \
make -C tests/acceptance/docker test-ocis-feature-ocis-storage
```

{{< hint info >}}
`BEHAT_FEATURE` must be pointing to a valid feature file
{{< /hint >}}

And to run a single scenario in a feature, you can do:

{{< hint info >}}
A specific scenario from a feature can be run by adding `:<line-number>` at the end of the feature file path. For example, to run the scenario at line 26 of the feature file `apiGraph/createUser.feature`, simply add the line number like this: `apiGraph/createUser.feature:26`. Note that the line numbers mentioned in the examples might not always point to a scenario, so always check the line numbers before running the test.
{{< /hint >}}

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
The test suites transferred from ownCloud core have `coreApi` prefixed
{{< /hint >}}

### oCIS Image to Be Tested (Skip Local Image Build)

By default, the tests will be run against the docker image built from your current working state of the oCIS repository. For some purposes it might also be handy to use an oCIS image from Docker Hub. Therefore, you can provide the optional flag `OCIS_IMAGE_TAG=...` which must contain an available docker tag of the [owncloud/ocis registry on Docker Hub](https://hub.docker.com/r/owncloud/ocis) (e.g. 'latest').

```bash
OCIS_IMAGE_TAG=latest \
make -C tests/acceptance/docker localApiTests-apiGraph-ocis
```

### Test Log Output

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

## Testing With Test Suite Natively Installed

We have two sets of tests:

- `test-acceptance-api` set was created for oCIS. Mainly for testing spaces features.

- `test-acceptance-from-core-api` set was transferred from [ownCloud core](https://github.com/owncloud/core) repository. The suite name of all tests transferred from the ownCloud core repository starts with `coreApi`

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

#### Run Local oCIS Tests (prefix `api`)

```bash
make test-acceptance-api \
TEST_SERVER_URL=https://localhost:9200 \
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
```

#### Run Tests Transferred From ownCloud Core (prefix `coreApi`)

```bash
make test-acceptance-from-core-api \
TEST_SERVER_URL=https://localhost:9200 \
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
```

Useful environment variables:

`TEST_SERVER_URL`: oCIS server url. Please, adjust the server url according to your setup.

`BEHAT_FEATURE`: to run a single feature

{{< hint info >}}
A specific scenario from a feature can be run by adding `:<line-number>` at the end of the feature file path. For example, to run the scenario at line 26 of the feature file `apiGraph/createUser.feature`, simply add the line number like this: `apiGraph/createUser.feature:26`. Note that the line numbers mentioned in the examples might not always point to a scenario, so always check the line numbers before running the test.
{{< /hint >}}

> Example:
>
> BEHAT_FEATURE=tests/acceptance/features/apiGraph/createUser.feature
>
> Or
>
> BEHAT_FEATURE=tests/acceptance/features/apiGraph/createUser.feature:12

`BEHAT_SUITE`: to run a single suite

> Example:
>
> BEHAT_SUITE=apiGraph

`STORAGE_DRIVER`: to run tests with a different user storage driver. Available options are `ocis` (default), `owncloudsql` and `s3ng`

> Example:
>
> STORAGE_DRIVER=owncloudsql

### Use Existing Tests for BDD

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

## Running ENV Config Tests (@env-Config)

Test suites tagged with `@env-config` are used to test the environment variables that are used to configure oCIS. These tests are special tests that require the oCIS server to be run using [ociswrapper](https://github.com/owncloud/ocis/blob/master/tests/ociswrapper/README.md).

### Run oCIS With ociswrapper

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

### Run the Tests

```bash
OCIS_WRAPPER_URL=http://localhost:5200 \
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
TEST_SERVER_URL="https://localhost:9200" \
BEHAT_FEATURE=tests/acceptance/features/apiAsyncUpload/delayPostprocessing.feature \
make test-acceptance-api
```

### Writing New ENV Config Tests

While writing tests for a new oCIS ENV configuration, please make sure to follow these guidelines:

1. Tag the test suite (or test scenarios) with `@env-config`
2. Use `OcisConfigHelper.php` for helper functions - provides functions to reconfigure the running oCIS instance.
3. Recommended: add the new step implementations in `OcisConfigContext.php`

## Running Test Suite With Email Service (@email)

Test suites that are tagged with `@email` require an email service. We use inbucket as the email service in our tests.

### Setup Inbucket

Run the following command to setup inbucket

```bash
docker run -d -p9000:9000 -p2500:2500 --name inbucket inbucket/inbucket
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

### Run the Acceptance Test

Run the acceptance test with the following command:

```bash
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
TEST_SERVER_URL="https://localhost:9200" \
EMAIL_HOST="localhost" \
EMAIL_PORT=9000 \
BEHAT_FEATURE="tests/acceptance/features/apiNotification/emailNotification.feature" \
make test-acceptance-api
```

## Running Tests for Parallel Deployment

### Setup the Parallel Deployment Environment

Instruction on setup is available [here](https://owncloud.dev/ocis/deployment/oc10_ocis_parallel/#local-setup)

Edit the `.env` file and uncomment this line:

```bash
COMPOSE_FILE=docker-compose.yml:testing/docker-compose-additions.yml
```

Start the docker stack with the following command:

```bash
docker-compose up -d
```

### Getting the Test Helpers

All the test helpers are located in the core repo.

```bash
git clone https://github.com/owncloud/core.git
```

### Run the Acceptance Tests

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

## Running Test Suite With Antivirus Service (@antivirus)
Test suites that are tagged with `@antivirus` require antivirus service. The available antivirus and the configuration related to them can be found [here](https://doc.owncloud.com/ocis/next/deployment/services/s-list/antivirus.html). This documentation is only going to use `clamAv` as antivirus.

### Setup clamAV
#### 1. Setup Locally
Run the following command to set up calmAV and clamAV daemon
```bash
sudo apt install clamav clamav-daemon -y
```

Make sure that the  clamAV daemon is up and running

```bash
sudo service clamav-daemon status
```
{{< hint info >}}
The commands are ubuntu specific and may differ according to your system. You can find information related to installation of clamAV in their official documentation [here](https://docs.clamav.net/manual/Installing/Packages.html).
{{< /hint>}}

#### 2. Setup clamAV With Docker
##### a. Create a Volume
For `clamAV` only local sockets can currently be configured we need to create a volume in order to share the socket with `oCIS server`. Run the following command to do so:
```bash
 docker volume create -d local -o device=/your/local/filesystem/path/ -o o=bind -o type=none clamav_vol
```
##### b. Run the Container
Run `clamAV` through docker and bind the path to the socket of clamAV from the image to the pre-created volume
```bash
docker run -v clamav_vol:/var/run/clamav/ owncloudci/clamavd
```
{{< hint info >}}
The path to the socket i.e. `/var/run/clamav/` may differ as per the image you are using. Make sure that you're providing the correct path to the socket if you're using image other than `owncloudci/clamavd`.
{{< /hint>}}

##### b. Change Ownership
Change the ownership of the path of your local filesystem that the volume `clamav_vol` is mounted on. After running `clamav` through docker the ownership of the bound path gets changed. As we need to provide this path to ocis server the ownership should be changed back to $USER or whatever ownership that your server requires.
```bash
 sudo chown -R $USER:$USER /your/local/filesystem/path/
```
{{< hint info >}}
Make sure that `clamAV` is fully up before running this command. The command is ubuntu specific and may differ according to your system.
{{< /hint>}}

{{< hint info >}}
If you want to use the same volume after the container is down. Before running the container once again you need to either remove all the data inside `/your/local/filesystem/path/` or give the ownership back. For instance, it ubuntu it might be `sudo chown -R systemd-network:systemd-journal /your/local/filesystem/path/`  and repeat step 2 and 3`
{{< /hint>}}

### Run oCIS

As `antivirus` service is not enabled by default we need to enable the service while running oCIS server. We also need to enable `async upload` and as virus scan is performed in post-processing step, we need to set it as well. Documentation for environment variables related to antivirus is available [here](https://owncloud.dev/services/antivirus/#environment-variables)

```bash
# run oCIS
PROXY_ENABLE_BASIC_AUTH=true \
ANTIVIRUS_SCANNER_TYPE="clamav" \
ANTIVIRUS_CLAMAV_SOCKET="/var/run/clamav/clamd.ctl" \
POSTPROCESSING_STEPS="virusscan" \
OCIS_ASYNC_UPLOADS=true \
OCIS_ADD_RUN_SERVICES="antivirus"
ocis/bin/ocis server
```
{{< hint info >}}
The value for `ANTIVIRUS_CLAMAV_SOCKET` is an example which needs adaption according your OS. If you are running `clamAv` with docker as per this documentation check the path that you mounted the volume i.e. `/your/local/filesystem/path/` to make sure the socket exists and give the full path to socket i.e. `/your/local/filesystem/path/clamd.sock` to `ANTIVIRUS_CLAMAV_SOCKET`.
{{< /hint>}}

#### Run the Acceptance Test

Run the acceptance test with the following command:

```bash
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
TEST_SERVER_URL="https://localhost:9200" \
BEHAT_FEATURE="tests/acceptance/features/apiAntivirus/antivirus.feature" \
make test-acceptance-api
```
