---
title: "Testing"
date: 2018-05-02T00:00:00+00:00
weight: 37
geekdocRepo: https://github.com/owncloud/ocis-reva
geekdocEditPath: edit/master/docs
geekdocFilePath: testing.md
---


## Acceptance tests

We are using the ownCloud 10 acceptance testsuite against ocis. To set this up you need the owncloud 10 core repo, an ldap server that the acceptance tests can use to manage users and the ocis-reva code.

### Getting the tests

All you need to do to get the acceptance tests is check out the core repo:
```
git clone https://github.com/owncloud/core.git
```

### Run an ldap server in a docker container

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

### Run ocis-reva with that ldap server

`ocis-reva` provides multiple subcommands. To configure them all via env vars you can export these environment variables. 

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
```

Then you need to start the ocis-reva services
```
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

### Notes
- rerunning the tests requires wiping the users in the ldap server, otherwise the tests will fail when trying to populate the users
- users are created with usernames like `user0`, the default password is `123456`