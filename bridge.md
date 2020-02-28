---
title: "Bridge"
date: 2020-02-27T20:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: bridge.md
---

We are planning to build a bridge from ownCloud 10 to ocis. The idea is to have a reverse proxy infront of ownCloud 10 that will forward requests to ownCloud 10 or ocis-reva, depending on the migration status of the logged in user.

This document is a work in progress of the current setup.

## Current status

Using ocis and the ownCloud 10 openidconnect and graphapi plugins it is possible today to introduce openid connect based authentication to existing instances. That is a prerequisite for migrating to ocis.

## How to do it

### Install the owncloud 10 graphapi app

In an owncloud 10 apps folder
```
$ git clone git@github.com:owncloud/graphapi.git
$ cd graphapi
$ composer install
```

No configuration necessary. You can test with `curl`:
```console
$ curl https://cloud.example.com/index.php/apps/graphapi/v1.0/users -u admin | jq
Enter host password for user 'jfd':
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   694  100   694    0     0   4283      0 --:--:-- --:--:-- --:--:--  4283
{
  "value": [
    {
      "id": "admin",
      "displayName": "admin",
      "mail": null
    },
    {
      "id": "demo",
      "displayName": "Demo",
      "mail": null
    },
    ...

  ],
  "@odata.nextLink": "https://oc.butonic.de/apps/graphapi/v1.0/users?$top=10&$skip=10"
}
```

> Note: The MS graph api actually asks for `Bearer` auth, but in order to check users passwords during an LDAP bind we are exploiting ownClouds authentication implementation that will grant access when `Basic` auth is used. An LDAP Bind you may ask? Read on!

### Start ocis-glauth

We are going to use the above ownCloud 10 and graphapi app to turn it into the datastore for an LDAP proxy. Grab it while it is hot:

In an `ocis` folder
```
$ git clone git@github.com:owncloud/ocis-glauth.git
$ cd ocis-glauth
$ git checkout start-glauth
$ make
```

TODO merge glauth PR https://github.com/owncloud/ocis-glauth/pull/1

You need to point `ocis-glauth` to your owncloud domain:
```console
$ bin/ocis-glauth --log-level debug server --backend-server https://cloud.example.com --backend-basedn dc=example,dc=com
```

`--log-level debug` is only used to generate more verbose output
`--backend-server https://cloud.example.com` is the url to an ownCloud instance with an enabled graphapi app
`--backend-basedn dc=example,dc=com` is used to construct the LDAP dn. The user `admin` will become `cn=admin,dc=example,dc=com`.

You should now be able to list accounts from your ownCloud 10 oc_accounts table using:
```console
$ ldapsearch -x -H ldap://localhost:9125 -b dc=example,dc=com -D "cn=admin,dc=example,dc=com" -W '(objectclass=posixaccount)'
```

Groups should work as well:
```console
$ ldapsearch -x -H ldap://localhost:9125 -b dc=example,dc=com -D "cn=admin,dc=example,dc=com" -W '(objectclass=posixgroup)'
```

> Note: This is currently a readonly implementation and minimal to the usecase of authenticating users with konnectd.

### Start ocis-konnectd

#### Get it

In an `ocis` folder
```
$ git clone git@github.com:owncloud/ocis-konnectd.git
$ cd ocis-konnectd
$ make
```

#### Environment variables

Konnectd needs environment variables to configure the LDAP server:
```console
export LDAP_URI=ldap://192.168.1.173:9125
export LDAP_BINDDN="cn=admin,dc=example,dc=com"
export LDAP_BINDPW="its-a-secret"
export LDAP_BASEDN="dc=example,dc=com"
export LDAP_SCOPE=sub
export LDAP_LOGIN_ATTRIBUTE=uid
export LDAP_EMAIL_ATTRIBUTE=mail
export LDAP_NAME_ATTRIBUTE=givenName
export LDAP_UUID_ATTRIBUTE=uid
export LDAP_UUID_ATTRIBUTE_TYPE=text
export LDAP_FILTER="(objectClass=posixaccount)"
```
Don't forget to use an existing user and the correct password.

### Configure clients

Now we need to configure a client we can later use to configure the ownCloud 10 openidconnect app. In the `assets/identifier-registration.yaml` have:
```yaml
---

# OpenID Connect client registry.
clients:
  - id: oc10-openidconnect-app
    name: openidconnect ownCloud app
    insecure: yes
    application_type: web
    redirect_uris:
      - https://cloud.example.com/apps/openidconnect/redirect
```
You will need the `insecure: yes` if you are using self signed certificates.

Replace the host in the redirect URI with your ownCloud 10 host and port.

#### Run it!

`ocis-konnectd` needs to know
- the issuer, which must be a reachable https endpoint. For testing an ip works. HTTPS is NOT optional.
- the identifier-registration.yaml you created
- a signature key id, otherwise the jwks key has no name, which might cause problems with clients. a random key is ok, but it should change when the actual signing key changes.

On the cli it looks like this
```console
$ bin/ocis-konnectd server -iss https://192.168.1.100:9130 --identifier-registration-conf assets/identifier-registration.yaml --signing-kid gen1-2020-02-27
```

#### Check it is up and running

1. Try getting the configuration:
```console
$ curl https://192.168.1.100:9130/.well-known/openid-configuration
```

2. Check if the login works at https://192.168.1.100:9130/signin/v1/identifier

> Note: If you later get a `Unable to find a key for (algorithm, kid):PS256, )` Error make sure you did set a `--signing-kid` when starting `ocis-konnectd` by checking it is present in https://192.168.1.100:9130/konnect/v1/jwks.json

### Patch owncloud

While the UserSession in ownCloud 10 is currently used to test all available IAuthModule implementations, it immediately logs out the user when an exception occurs. However, existing owncloud 10 instances use the oauth2 app to create Bearer tokens for mobile and desktop clients.

To give the openidconnect app a chance to verify the tokens we need to change the code a bit. See https://github.com/owncloud/core/pull/37043 for a possible solution.

> Note: The PR is hot ... as in *younger than this list of steps*. And it messes with authentication. Use with caution.

### Install the owncloud 10 openidconnect app

In an owncloud 10 apps folder
```
$ git clone git@github.com:owncloud/openidconnect.git
$ cd openidconnect
$ composer install
```

After enabling the app configure it in `config/oidc.config.php`

```php
$CONFIG = [
  'openid-connect' => [
    'provider-url' => 'https://192.168.1.100:9130',
    'client-id' => 'oc10-openidconnect-app',
    'loginButtonName' => 'OpenId Connect @ Konnectd',
  ],
  'debug' => true // if using self signed certificates
];
```

In the above configuration replace
- `provider-url` with the url to your `ocis-konnectd` issuer

> Note: By default the openidconnect app will use the email of the user to match the user from the oidc userinfo endpoint with the ownCloud account. So make sure your users have a unique primary email.

## Next steps

Aside from the above todos these are the next stepo
- get `ocis-phoenix` configured to authenticate against `ocis-konnectd` and use the webdav endpoint from owncloud 10.
- tie it all together behind `ocis-proxy` 
- create an `ocis bridge` command that runs all the ocis services in one step with a properly preconfigured `ocis-konnectd` `identifier-registration.yaml` file for `phoenix` and the owncloud 10 `openidconnect` app, as well as a randomized `--signing-kid`.