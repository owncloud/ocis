---
title: "External IdP"
date: 2020-08-25T15:31:23+01:00
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: external-idp.md
---

{{< toc >}}
To set up OCIS to use an external IdP you need to configure the services to use that IdP using environment variables (or config flags) or try a docker compose setup.

# Environment variables

`PHOENIX_OIDC_AUTHORITY="https://idp.owncloud.works"`
`PHOENIX_OIDC_METADATA_URL="https://idp.owncloud.works/.well-known/openid-configuration"`
`PHOENIX_OIDC_CLIENT_ID=?` something that the idp might create

`PROXY_OIDC_ISSUER="https://idp.owncloud.works"`

`REVA_OIDC_ISSUER="https://idp.owncloud.works"` (should not be necessary, as users are authenticated by the proxy middleware)

`GRAPH_OIDC_ENDPOINT="https://idp.owncloud.works"`

# Docker Compose

You can start an IdP using
```
docker-compose -f docker-compose.idp.yml up
```

## Create an Account
You can add accounts to the IdP using:

```
docker-compose-f docker-compose.idp.yml exec ocis-accounts add --preferred-name bob --on-premises-sam-account-name bob --displayname "Bob" --uidnumber 33333 --gidnumber 30000 --password 123456 --mail bob@example.org --enabled
```

`--preferred-name bob` and `--on-premises-sam-account-name bob` should be the username. The duplication is reserved for future usage.

## List Accounts

You can list accounts using:
```
docker-compose-f docker-compose.idp.yml exec ocis-accounts ls
```

## Start ocis

Then ocis can be started using
```
docker-compose -f docker-compose.ocis-external-idp.yml -f docker-compose.idp.yml up
```

Both `.yml` files are necessary to let the ocis proxy connect to the IdP. For this to work you also need to add a hosts entry `konnectd` to point to the IdP.
This is easier when running on properly available domains as otherwise a lot af insecure flags have to be enabled.


# TODO
- add global `OIDC_ISSUER` for all
- document how to run it on localhost using fake domain names/entries in the hosts file and in the docker containers so all services can talk to the idp?

