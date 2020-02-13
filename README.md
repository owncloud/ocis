# ownCloud Infinite Scale: Reva

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-reva/status.svg)](https://cloud.drone.io/owncloud/ocis-reva)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6f1eaaa399294d959ef7b3b10deed41d)](https://www.codacy.com/manual/owncloud/ocis-reva?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-reva&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-reva?status.svg)](http://godoc.org/github.com/owncloud/ocis-reva)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-reva)](http://goreportcard.com/report/github.com/owncloud/ocis-reva)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-reva.svg)](http://microbadger.com/images/owncloud/ocis-reva "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/reva/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/ocis-reva/)

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html).

```console
git clone https://github.com/owncloud/ocis-reva.git
cd ocis-reva

make generate build

./bin/ocis-reva -h
```

To run a demo installation you can use the preconfigured defaults and start all necessary services:
```
bin/ocis-reva frontend & \
bin/ocis-reva gateway & \
bin/ocis-reva users & \
bin/ocis-reva auth-basic & \
bin/ocis-reva auth-bearer & \
bin/ocis-reva sharing & \
bin/ocis-reva storage-root & \
bin/ocis-reva storage-home & \
bin/ocis-reva storage-home-data & \
bin/ocis-reva storage-oc & \
bin/ocis-reva storage-oc-data
```

The root storage serves the available namespaces from disk using the local storage driver. In order to be able to navigate into the `/home` and `/oc` storage providers you have to create these directories:
```
mkdir /var/tmp/reva/root/home
mkdir /var/tmp/reva/root/oc
```

Note: the owncloud storage driver currently requires a redis server running on the local machine.

You should now be able to get a file listing of a users home using
```
curl -X PROPFIND http://localhost:9140/remote.php/dav/files/ -v -u einstein:relativity
```

## Users

The default config uses the demo user backend, which contains three users:
```
einstein:relativity
marie:radioactivty
richard:superfluidity
```

For details on the `json` and `ldap` backends see the [documentation](https://owncloud.github.io/ocis-reva/#users)

## Run oC10 API acceptance tests
1.  start an LDAP server e.g. with docker:
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
2.  start a Redis server e.g. with docker:
    `docker run -e REDIS_DATABASES=1 -p 6379:6379 -d webhippie/redis:latest`

3.  clone oC10 code: `git clone https://github.com/owncloud/core.git`

4.  start reva with `REVA_USERS_DRIVER=ldap`:
    ```
    bin/ocis-reva gateway & \
    REVA_USERS_DRIVER=ldap bin/ocis-reva users & \
    REVA_USERS_DRIVER=ldap bin/ocis-reva auth-basic & \
    bin/ocis-reva auth-bearer & \
    bin/ocis-reva sharing & \
    bin/ocis-reva storage-root & \
    bin/ocis-reva storage-home & \
    bin/ocis-reva storage-home-data & \
    bin/ocis-reva storage-oc & \
    bin/ocis-reva storage-oc-data & \
    bin/ocis-reva frontend
    ```

5.  from inside the oC10 repo run the tests:
    ```
    make test-acceptance-api \
        TEST_SERVER_URL=http://localhost:9140 \
        TEST_EXTERNAL_USER_BACKENDS=true \
        TEST_OCIS=true \
        OCIS_REVA_DATA_ROOT=/var/tmp/reva/ \
        BEHAT_FILTER_TAGS='~@skipOnOcis&&@TestAlsoOnExternalUserBackend&&~@skipOnLDAP'
    ```

    Make sure to adjust the settings `TEST_SERVER_URL` and `OCIS_REVA_DATA_ROOT` according to your environment

    This will run all tests that can work with LDAP and are not skipped on OCIS
    To run a subset of tests, e.g. a single suite, file or tag have a look at the [acceptance tests documentation](https://doc.owncloud.com/server/10.0/developer_manual/core/acceptance-tests.html#running-acceptance-tests-for-a-suite).
    E.g. you can run all tests that are marked with a specific issue:
    ```
    make test-acceptance-api \
        TEST_SERVER_URL=http://localhost:9140 \
        TEST_EXTERNAL_USER_BACKENDS=true \
        TEST_OCIS=true \
        OCIS_REVA_DATA_ROOT=/var/tmp/reva/ \
        BEHAT_FILTER_TAGS='@TestAlsoOnExternalUserBackend&&~@skipOnLDAP&&@issue-ocis-reva-46'
    ```

    Note that the `~@skipOnOcis` tag is removed here, because to fix an issue you want also to run the tests that are skipped in the CI run

## Security

If you find a security issue please contact security@owncloud.com first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2019 ownCloud GmbH <https://owncloud.com>
```
