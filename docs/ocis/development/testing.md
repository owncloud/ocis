---
title: "Testing"
date: 2018-05-02T00:00:00+00:00
weight: 37
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: testing.md
---

{{< toc >}}

## Acceptance tests

We are using the ownCloud 10 acceptance test suite against oCIS. To set this up you need the ownCloud 10 core repo, a LDAP server that the acceptance tests can use to manage users, a redis server for file-versions and the oCIS code.

### Getting the tests

All you need to do to get the acceptance tests is check out the core repo:

```
git clone https://github.com/owncloud/core.git
```

### Run a redis server in a docker container

File versions need a redis server. Start one with docker by using:

`docker run -e REDIS_DATABASES=1 -p 6379:6379 -d webhippie/redis:latest`

### Run ocis

To start ocis:

```
PROXY_ENABLE_BASIC_AUTH=true bin/ocis server
```

`PROXY_ENABLE_BASIC_AUTH` will allow the acceptance tests to make requests against the provisioning api (and other endpoints) using basic auth.

### Run the acceptance tests

First we will need to clone the testing app in owncloud which contains the skeleton files required for running the tests.
In the ownCloud 10 core clone the testing app with the following command:

```
git clone https://github.com/owncloud/testing apps/testing
```

Then run the api acceptance tests with the following command from the root of the oCIS repository:
```
make test-acceptance-api \
TEST_SERVER_URL=https://localhost:9200 \
TEST_OCIS=true \
SKELETON_DIR=apps/testing/data/apiSkeleton \
DELETE_USER_DATA_CMD='rm -rf /var/tmp/ocis/storage/users/nodes/root/* /var/tmp/ocis/storage/users/nodes/*-*-*-*' \
BEHAT_FILTER_TAGS='~@notToImplementOnOCIS&&~@toImplementOnOCIS'
```

Make sure to adjust the settings `TEST_SERVER_URL` and `OCIS_REVA_DATA_ROOT` according to your environment.

This will run all tests that are relevant to OCIS.

To run a single test add `BEHAT_FEATURE=<feature file>`

### use existing tests for BDD

As a lot of scenarios are written for oC10, we can use those tests for Behaviour driven development in ocis.
Every scenario that does not work in OCIS with "owncloud" storage, is listed in `tests/acceptance/expected-failures-on-OWNCLOUD-storage.txt` with a link to the related issue.
Every scenario that does not work in OCIS with "ocis" storage, is listed in `tests/acceptance/expected-failures-on-OCIS-storage.txt` with a link to the related issue.

Those scenarios are run in the ordinary acceptance test pipeline in CI. The scenarios that fail are checked against the
expected failures. If there are any differences then the CI pipeline fails.
Similarly, scenarios that do not work in OCIS with EOS storage are listed in `tests/acceptance/expected-failures-on-EOS-storage.txt`.
Additionally, some issues have scenarios that demonstrate the current buggy behaviour in ocis(reva).
Those scenarios are in this ocis repository in `tests/acceptance/features/apiOcisSpecific`.
Have a look into the [documentation](https://doc.owncloud.com/server/developer_manual/testing/acceptance-tests.html#writing-scenarios-for-bugs) to understand why we are writing those tests.

If you want to work on a specific issue

1.  adjust the core commit id to the latest commit in core so that CI will run the latest test code and scenarios from core.
    For that change `coreCommit` in the `config` section:

        config = {
          'apiTests': {
            'coreBranch': 'master',
            'coreCommit': 'a06b1bd5ba8e5244bfaf7fa04f441961e6fb0daa',
            'numberOfParts': 6
          }
        }

2.  locally run each of the tests marked with that issue in the expected failures file.

    E.g.:

    ```
    make test-acceptance-api \
    TEST_SERVER_URL=https://localhost:9200 \
    TEST_OCIS=true \
    DELETE_USER_DATA_CMD='rm -rf /var/tmp/ocis/storage/users/nodes/root/* /var/tmp/ocis/storage/users/nodes/*-*-*-*' \
    BEHAT_FEATURE='tests/acceptance/features/apiComments/comments.feature:123'
    ```

3.  the tests will fail, try to understand how and why they are failing
4.  fix the code
5.  go back to 2. and repeat till the tests are passing.
6.  remove those tests from the expected failures file
7.  run each of the local tests that were demonstrating the **buggy** behavior. They should fail.
8.  delete each of the local tests that were demonstrating the **buggy** behavior.
9.  make a PR that has the fixed code, relevant lines removed from the expected failures file and bug demonstration tests deleted.
