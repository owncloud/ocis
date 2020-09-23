* * *

title: "Testing"
date: 2018-05-02T00:00:00+00:00
weight: 37
geekdocRepo: <https://github.com/owncloud/ocis-reva>
geekdocEditPath: edit/master/docs

## geekdocFilePath: testing.md

## API Acceptance tests

We are using the ownCloud 10 API acceptance testsuite against ocis. To set this up you need the owncloud 10 core repo, a ldap server that the acceptance tests can use to manage users, a redis server for file-versions and the ocis-reva code.

### Getting the tests

All you need to do to get the acceptance tests is check out the core repo:

    git clone https://github.com/owncloud/core.git

### Run a ldap server in a docker container

The ownCloud 10 acceptance tests will need write permission. You can start a suitable ldap server in a docker container with:

    docker run --hostname ldap.my-company.com \
        -e LDAP_TLS_VERIFY_CLIENT=never \
        -e LDAP_DOMAIN=owncloud.com \
        -e LDAP_ORGANISATION=ownCloud \
        -e LDAP_ADMIN_PASSWORD=admin \
        --name docker-slapd \
        -p 127.0.0.1:389:389 \
        -p 636:636 -d osixia/openldap

### Run a redis server in a docker container

File versions need a redis server. Start one with docker by using:

`docker run -e REDIS_DATABASES=1 -p 6379:6379 -d webhippie/redis:latest`

### Run ocis-reva with that ldap server

`ocis-reva` provides multiple subcommands. To configure them all via env vars you can export these environment variables.

    export REVA_USERS_DRIVER=ldap
    export REVA_LDAP_HOSTNAME=localhost
    export REVA_LDAP_PORT=636
    export REVA_LDAP_BASE_DN='dc=owncloud,dc=com'
    export REVA_LDAP_USERFILTER='(&(objectclass=posixAccount)(cn=%s))'
    export REVA_LDAP_GROUPFILTER='(&(objectclass=posixGroup)(cn=%s))'
    export REVA_LDAP_BIND_DN='cn=admin,dc=owncloud,dc=com'
    export REVA_LDAP_BIND_PASSWORD=admin
    export REVA_LDAP_SCHEMA_UID=uid
    export REVA_LDAP_SCHEMA_MAIL=mail
    export REVA_LDAP_SCHEMA_DISPLAYNAME=displayName
    export REVA_LDAP_SCHEMA_CN=cn
    export REVA_FRONTEND_URL=http://localhost:9140 # needed because the proxy is not started
    export REVA_DATAGATEWAY_URL=http://localhost:9140/data # needed because the proxy is not started

Then you need to start the ocis-reva services

    bin/ocis-reva frontend & \
    bin/ocis-reva gateway & \
    bin/ocis-reva auth-basic & \
    bin/ocis-reva auth-bearer & \
    bin/ocis-reva sharing & \
    bin/ocis-reva storage-home & \
    bin/ocis-reva storage-home-data & \
    bin/ocis-reva storage-oc & \
    bin/ocis-reva storage-oc-data & \
    bin/ocis-reva users &

### Run the API acceptance tests

In the ownCloud 10 core repo run

    make test-acceptance-api \
    TEST_SERVER_URL=http://localhost:9140 \
    TEST_EXTERNAL_USER_BACKENDS=true \
    TEST_OCIS=true \
    OCIS_REVA_DATA_ROOT=/var/tmp/reva/ \
    BEHAT_FILTER_TAGS='~@notToImplementOnOCIS&&~@toImplementOnOCIS&&~@preview-extension-required' \
    SKELETON_DIR=apps/testing/data/apiSkeleton

Make sure to adjust the settings `TEST_SERVER_URL`,`OCIS_REVA_DATA_ROOT` and `SKELETON_DIR` according to your environment.

This will run all tests that are relevant to OCIS.

To run a single test add `BEHAT_FEATURE=<feature file>` and specify the path to the feature file and an optional line number. For example: `BEHAT_FEATURE='tests/acceptance/features/apiWebdavUpload1/uploadFile.feature:12'`

### use existing tests for BDD

As a lot of scenarios are written for oC10, we can use those tests for Behaviour driven development in ocis.
Every scenario that does not work in OCIS with OC storage, is listed in `tests/acceptance/expected-failures-on-OC-storage.txt` with a link to the related issue.
Similarly, scenarios that do not work in OCIS with EOS storage are listed in `tests/acceptance/expected-failures-on-EOS-storage.txt`.
Scenarios from the oC10 API acceptance tests are run in the ordinary acceptance test pipeline in CI. The scenarios that fail are checked against the
expected failures. If there are any differences then the CI pipeline fails.

Additionally, some issues have scenarios that demonstrate the current buggy behaviour in ocis(reva).
Those scenarios are in this ocis-reva repository in `tests/acceptance/features/apiOcisSpecific`.
Have a look into the [documentation](https://doc.owncloud.com/server/developer_manual/testing/acceptance-tests.html#writing-scenarios-for-bugs) to understand why we are writing those tests.
Also, ocis behaves partly differently with EOS-Storage and OC-Storage. There are scenarios that do not work in OCIS when run on EOS-storage, but works when on OC-Storage, and vice-versa. For those kind of scenarios, ` @skipOnOcis-EOS-Storage` and `@skipOnOcis-OC-Storage` tags are used. For instance, for a scenario that fails on EOS-Storage but passes on OC-Storage, we use `@skipOnOcis-EOS-Storage` tag to let it run on OC-Storage, where it works as expected, instead of skipping the test completely.

If you want to work on a specific issue

1.  adjust the core commit id to the latest commit in core so that CI will run the latest test code and scenarios from core.
    For that change `coreCommit` in the `config` section:

        config = {
          'apiTests': {
            'coreBranch': 'master',
            'coreCommit': 'a06b1bd5ba8e5244bfaf7fa04f441961e6fb0daa',
            'numberOfParts': 2
          }
        }

2.  locally run each of the tests marked with that issue in the expected failures file:

    E.g.:

        make test-acceptance-api \
        TEST_SERVER_URL=http://localhost:9140 \
        TEST_EXTERNAL_USER_BACKENDS=true \
        TEST_OCIS=true \
        OCIS_REVA_DATA_ROOT=/var/tmp/reva/ \
        BEHAT_FEATURE='tests/acceptance/features/apiComments/comments.feature:123'

3.  the tests will fail, try to understand how and why they are failing

4.  fix the code

5.  go back to 2. and repeat till the tests are passing.

6.  remove those tests from the expected failures file.

7.  run each of the local tests that were demonstrating the **buggy** behavior. They should fail.

8.  delete each of the local tests that were demonstrating the **buggy** behavior.

9.  make a PR that has the fixed code, relevant lines removed from the expected failures file and bug demonstration tests deleted.

    If the changes also affect the `ocis` repository make sure the changes get ported over there.
    That will need the fixed code in `ocis-reva` to be applied to `ocis` along with the test-related changes.

### Notes

-   in a normal case the test-code cleans up users after the test-run, but if a test-run is interrupted (e.g. by CTRL+C) users might have been left on the LDAP server. In that case rerunning the tests requires wiping the users in the ldap server, otherwise the tests will fail when trying to populate the users. This can be done by simply running `docker stop docker-slapd && docker rm docker-slapd` and [restarting the LDAP server container](#run-a-ldap-server-in-a-docker-container)
-   the tests usually create users in the OU `TestUsers` with usernames specified in the feature file. If not defined in the feature file, most users have the password `123456`, defined by `regularUserPassword` in `behat.yml`, but other passwords are also used, see [`\FeatureContext::getPasswordForUser()`](https://github.com/owncloud/core/blob/master/tests/acceptance/features/bootstrap/FeatureContext.php#L386) for mapping and [`\FeatureContext::__construct`](https://github.com/owncloud/core/blob/master/tests/acceptance/features/bootstrap/FeatureContext.php#L1668) for the password definitions.
