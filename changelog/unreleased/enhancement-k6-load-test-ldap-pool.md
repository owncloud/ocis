Enhancement: Enable the LDAP connection pool in the k6 load test

The k6 load test started oCIS with the default LDAP settings. In this setup
the services graph, users, groups and auth-basic each use one LDAP connection.
All LDAP requests of a service wait for this connection. Under a high load,
the test then measures this queue and not the other parts of oCIS.

The workflow `.github/workflows/k6-load-test.yml` now starts oCIS with these
settings:

- `OCIS_LDAP_POOL_ENABLED=true` uses a pool of LDAP connections.
- `OCIS_LDAP_POOL_SIZE=20` sets 20 connections for each pool.
- `OCIS_LDAP_UPDATE_LAST_SIGNIN_DATE=false` stops the LDAP write at each login.
- `SHARING_USER_JSONCS3_MAX_CONCURRENCY=5` limits the parallel work of the
  share manager.

A new step in the workflow reads the environment of the oCIS process. The step
fails the job when one of the settings is missing. Thus a test does not run with
the wrong settings.

https://github.com/owncloud/ocis/pull/13534
