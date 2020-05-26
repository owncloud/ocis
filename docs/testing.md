---
title: "Testing"
date: 2018-05-02T00:00:00+00:00
weight: 37
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: testing.md
---


## Acceptance tests

We are using the ownCloud 10 acceptance testsuite against ocis. To set this up you need the owncloud 10 core repo, a ldap server that the acceptance tests can use to manage users, a redis server for file-versions and the ocis code.

### Getting the tests

All you need to do to get the acceptance tests is check out the core repo:
```
git clone https://github.com/owncloud/core.git
```

### Run a ldap server in a docker container

The ownCloud 10 acceptance tests will need write permission. You can start a suitable ldap server in a docker container with:

```
docker run --hostname ldap.my-company.com \
    -e LDAP_TLS_VERIFY_CLIENT=never \
    -e LDAP_DOMAIN=owncloud.com \
    -e LDAP_ORGANISATION=ownCloud \
    -e LDAP_ADMIN_PASSWORD=admin \
    --name docker-slapd \
    -p 127.0.0.1:389:389 \
    -p 636:636 -d osixia/openldap
```
### Run a redis server in a docker container

File versions need a redis server. Start one with docker by using:

`docker run -e REDIS_DATABASES=1 -p 6379:6379 -d webhippie/redis:latest`

### Run ocis with that ldap server

`ocis` provides multiple subcommands. To configure them all via env vars you can export these environment variables.

```
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
export LDAP_URI=ldap://localhost
export LDAP_BINDDN='cn=admin,dc=owncloud,dc=com'
export LDAP_BINDPW=admin
export LDAP_BASEDN='dc=owncloud,dc=com'
```

Then you need to start ocis
```
bin/ocis server
```

### Run the acceptance tests

In the ownCloud 10 core repo run

```
make test-acceptance-api \
TEST_SERVER_URL=http://localhost:9140 \
TEST_EXTERNAL_USER_BACKENDS=true \
TEST_OCIS=true \
OCIS_REVA_DATA_ROOT=/var/tmp/reva/ \
BEHAT_FILTER_TAGS='~@skipOnOcis&&~@skipOnLDAP&&@TestAlsoOnExternalUserBackend&&~@local_storage'
```

Make sure to adjust the settings `TEST_SERVER_URL` and `OCIS_REVA_DATA_ROOT` according to your environment

This will run all tests that can work with LDAP and are not skipped on ocis

To run a single test add `BEHAT_FEATURE=<feature file>`

### use existing tests for BDD

As a lot of scenarios are written for oC10, we can use those tests for Behaviour driven development in ocis.
Every scenario that does not work in ocis, is tagged with `@skipOnOcis` and additionally should be marked with an issue number e.g. `@issue-ocis-20`.
This tag means that this particular scenario is skipped because of [issue no 20 in the ocis repository](https://github.com/owncloud/ocis/issues/20).
Additionally, some issues have scenarios that demonstrate the current buggy behaviour in ocis(reva) and are skipped on oC10.
Have a look into the [documentation](https://doc.owncloud.com/server/developer_manual/testing/acceptance-tests.html#writing-scenarios-for-bugs) to understand why we are writing those tests.

If you want to work on a specific issue

1.  run the tests marked with that issue tag

    E.g.:
    ```
    make test-acceptance-api \
    TEST_SERVER_URL=http://localhost:9140 \
    TEST_EXTERNAL_USER_BACKENDS=true \
    TEST_OCIS=true \
    OCIS_REVA_DATA_ROOT=/var/tmp/reva/ \
    BEHAT_FILTER_TAGS='~@skipOnOcV10&&~@skipOnLDAP&&@TestAlsoOnExternalUserBackend&&~@local_storage&&@issue-ocis-20'
    ```

    Note that the `~@skipOnOcis` tag is replaced by `~@skipOnOcV10` and the issue tag `@issue-ocis-20` is added.
    We want to run all tests that are skipped in CI because of this particular bug, but we don't want to run the tests
    that demonstrate the current buggy behaviour.

2.  the tests will fail, try to understand how and why they are failing
3.  fix the code
4.  go back to 1. and repeat till the tests are passing.
5.  adjust tests that demonstrate the **buggy** behaviour

    delete the tests in core that are tagged with that particular issue and `@skipOnOcV10`, but be careful because a lot of tests are tagged with multiple issues.
    Only delete tests that demonstrate the buggy behaviour if you fixed all bugs related to that test. If not you might have to adjust the test.
6.  unskip tests that demonstrate the **correct** behaviour

    The `@skipOnOcis` tag should not be needed now, so delete it, but leave the issue tag for future reference.
7.  make a PR to core with the changed tests
8.  make a PR to ocis running the adjusted tests

    To confirm that all tests (old and changed) run fine make a PR to ocis with your code changes and point drone to your branch in core to get the changed tests.
    For that change this line in the `acceptance-tests` section

    `'git clone -b master --depth=1 https://github.com/owncloud/core.git /srv/app/testrunner',`

    to clone your core branch e.g.

    `'git clone -b fixRevaIssue122 --depth=1 https://github.com/owncloud/core.git /srv/app/testrunner',`

9.  merge PRs

    After you have confirmed that the tests pass everywhere merge the core PR and immediately revert the change in 8. and merge the ocis PR

    If the changes also affect the `ocis-reva` repository make sure the changes get ported over there immediately, otherwise the tests will start failing there.


### Notes
- in a normal case the test-code cleans up users after the test-run, but if a test-run is interrupted (e.g. by CTRL+C) users might have been left on the LDAP server. In that case rerunning the tests requires wiping the users in the ldap server, otherwise the tests will fail when trying to populate the users.
- the tests usually create users in the OU `TestUsers` with usernames specified in the feature file. If not defined in the feature file, most users have the password `123456`, defined by `regularUserPassword` in `behat.yml`, but other passwords are also used, see [`\FeatureContext::getPasswordForUser()`](https://github.com/owncloud/core/blob/master/tests/acceptance/features/bootstrap/FeatureContext.php#L386) for mapping and [`\FeatureContext::__construct`](https://github.com/owncloud/core/blob/master/tests/acceptance/features/bootstrap/FeatureContext.php#L1668) for the password definitions.
