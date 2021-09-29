---
title: "Bridge"
date: 2020-02-27T20:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: bridge.md
---

{{< toc >}}

We are planning to build a bridge from ownCloud 10 to ocis. The idea is to have a reverse proxy infront of ownCloud 10 that will forward requests to ownCloud 10 or ocis-reva, depending on the migration status of the logged in user.

This document is a work in progress of the current setup.

## Current status

Using ocis and the ownCloud 10 [graphapi app](https://github.com/owncloud/graphapi/) it is possible today to use an existing owncloud 10 instance as a userbackend and storage backend for ocis.

## How to do it

### Install the owncloud 10 graphapi app

In an owncloud 10 apps folder
```
$ git clone git@github.com:owncloud/graphapi.git
$ cd graphapi
$ composer install
```

### Enable the graphapi app

```
occ a:e graphapi
```

No configuration necessary. You can test with `curl`:
```console
$ curl https://cloud.ocis.test/index.php/apps/graphapi/v1.0/users -u admin -s | jq
Enter host password for user 'admin':
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
  "@odata.nextLink": "https://cloud.ocis.test/apps/graphapi/v1.0/users?$top=10&$skip=10"
}
```

{{< hint >}}
The MS graph api actually asks for `Bearer` auth, but in order to check users passwords during an LDAP bind we are exploiting ownClouds authentication implementation that will grant access when `Basic` auth is used. An LDAP Bind you may ask? Read on!
{{< /hint >}}

### Grab ocis!

```
$ git clone git@github.com:owncloud/ocis.git
$ cd ocis
$ make -C ocis build
```
This should give you an `ocis/bin/ocis` binary. Try listing the help with `ocis/bin/ocis --help`.

{{< hint >}}
You can check out a custom branch and build a custom binary which can then be used for the below steps.
{{< /hint >}}

### Start ocis glauth

We are going to use the built binary and ownCloud 10 graphapi app to turn ownCloud 10 into the datastore for an LDAP proxy.

#### configure it

While ocis can be configured using environment variables, eg. for a docker compose setup we are going to use a more traditional config file here.
Create a config file for ocis in either `/etc/ocis`, `$HOME/.ocis` or `./.config`. You can use `.json`, `.yaml` or `.toml`. I will use toml here, because ... reasons.

```toml
[glauth.backend]
datastore = "owncloud"                                   # switch to the owncloud datastore
servers = ["https://cloud.ocis.test/apps/graphapi/v1.0"] # the graph api endpoint to connect to
basedn = "dc=ocis,dc=test"                               # base dn to construct the LDAP dn. The user `admin` will become `cn=admin,dc=ocis,dc=test`
```

{{< hint >}}
There is a bug in the config merging for environment variables, cli flags and config files causing log settings not to be picked up from the config file when specifying `--extensions`. That is why I will
* configure most of the config in a file, 
* adjust logging using `OCIS_LOG_*` environment variables and
* specify which extension to run using `ocis/bin/ocis server --extensions "comma, separated, list, of, extensions"`. 
{{< /hint >}}

#### Run it!

For now, we only start the glauth extension:
```console
$ OCIS_LOG_PRETTY=true OCIS_LOG_COLOR=true ocis/bin/ocis server --extensions "glauth"
```

#### Check it is up and running

You should now be able to list accounts from your ownCloud 10 oc_accounts table using:
```console
$ ldapsearch -x -H ldap://127.0.0.1:9125 -b dc=ocis,dc=test -D "cn=admin,dc=ocis,dc=test" -W '(objectclass=posixaccount)'
```

Groups should work as well:
```console
$ ldapsearch -x -H ldap://127.0.0.1:9125 -b dc=ocis,dc=test -D "cn=admin,dc=ocis,dc=test" -W '(objectclass=posixgroup)'
```

{{< hint >}}
This is currently a readonly implementation and minimal to the usecase of authenticating users with an IDP.
{{< /hint >}}

### Start ocis storage-gateway, storage-authbasic and storage-userprovider

We are going to set up reva to authenticate users against our glauth LDAP proxy. This allows us to log in and use the reva cli. The ocis storage-gateway starts the reva gateway which will authenticate basic auth requests using the storage-authbasic service. Furthermore, users have to be available in the storage-userprovider to retrieve displayname, email address and other user metadata.

To configure LDAP to use our glauth we add this section to the config file:

```toml
[storage.reva.ldap]
idp = "https://ocis.ocis.test"
basedn = "dc=ocis,dc=test" 
binddn = "cn=admin,dc=ocis,dc=test" # an admin user in your oc10
bindpassword = "secret"
userschema = { uid = "uid", displayname = "givenname" } # TODO make glauth return an ownclouduuid and displayname attribute
```

Now we can start all necessary services.

```console
$ OCIS_LOG_PRETTY=true OCIS_LOG_COLOR=true ocis/bin/ocis server --extensions "glauth, storage-gateway, storage-authbasic, storage-userprovider"
```


{{< hint warning >}}
Here I ran out of time. I tried to verify this step with the reva cli:
`cmd/reva/reva -insecure -host localhost:9142`
`login basic`
but it tries to create the user home, which cannot be disabled in a config file: https://github.com/owncloud/ocis/issues/2416#issuecomment-901197053

starting `STORAGE_GATEWAY_DISABLE_HOME_CREATION_ON_LOGIN=true OCIS_LOG_LEVEL=debug OCIS_LOG_PRETTY=true OCIS_LOG_COLOR=true ocis/bin/ocis server --extensions "storage-gateway, storage-authbasic, storage-userprovider"` let me login:

```console
✗ cmd/reva/reva -insecure -host localhost:9142
reva-cli v1.11.0-27-g95b1f2ee (rev-95b1f2ee)
Please use `exit` or `Ctrl-D` to exit this program.
>> login basic
username: jfd
password: OK
>> whoami
id:<idp:"https://ocis.ocis.test" opaque_id:"jfd" type:USER_TYPE_PRIMARY > username:"jfd" mail:"jfd@butonic.de" display_name:"J\303\266rn" uid_number:99 gid_number:99 
>> exit
```

I hope https://github.com/owncloud/ocis/pull/2024 fixes the parsing order of things.

everything below this is outdated

... gotta run
{{< /hint >}}


### Start ocis storage-userprovider

```console
ocis/bin/ocis storage-userprovider --ldap-port 19126 --ldap-user-schema-uid uid --ldap-user-schema-displayName givenName --addr :19144
```

TODO clone `git clone git@github.com:cs3org/cs3apis.git`

query users using [grpcurl](https://github.com/fullstorydev/grpcurl)
```console
grpcurl -import-path ./cs3apis/ -proto ./cs3apis/cs3/identity/user/v1beta1/user_api.proto -plaintext localhost:19144 cs3.identity.user.v1beta1.UserAPI/FindUsers
ERROR:
  Code: Unauthenticated
  Message: auth: core access token not found
```


### Start ocis idp

#### Set environment variables

The built in [libregraph/lico](https://github.com/libregraph/lico) needs environment variables to configure the LDAP server:
```console
export OCIS_URL=https://ocis.ocis.test
export IDP_LDAP_URI=ldap://127.0.0.1:9125
export IDP_LDAP_BASE_DN="dc=ocis,dc=test"
export IDP_LDAP_BIND_DN="cn=admin,dc=ocis,dc=test"
export IDP_LDAP_BIND_PASSWORD="its-a-secret"
export IDP_LDAP_SCOPE=sub
export IDP_LDAP_LOGIN_ATTRIBUTE=uid
export IDP_LDAP_NAME_ATTRIBUTE=givenName
```
Don't forget to use an existing user with admin permissions (only admins are allowed to list all users via the graph api) and the correct password.

{{< hint warning >}}
* TODO: change the default values in glauth & ocis to use an `ownclouduuid` attribute.
* TODO: split `OCIS_URL` and `IDP_ISS` env vars and use `OCIS_URL` to generate the clients in the `identifier-registration.yaml`.
{{< /hint >}}

### Configure clients

When the `identifier-registration.yaml` does not exist it will be generated based on the `OCIS_URL` environment variable.

#### Run it!

You can now bring up `ocis/bin/ocis idp` with:
```console
$ ocis/bin/ocis idp server --iss http://127.0.0.1:9130 --signing-kid gen1-2020-02-27
```

`ocis/bin/ocis idp` needs to know
- `--iss http://127.0.0.1:9130` the issuer, which must be a reachable http endpoint. For testing an ip works. For openid connect HTTPS is NOT optional. This URL is exposed in the `http://127.0.0.1:9130/.well-known/openid-configuration` endpoint and clients need to be able to connect to it, securely. We will change this when introducing the proxy.
- `--signing-kid gen1-2020-02-27` a signature key id, otherwise the jwks key has no name, which might cause problems with clients. a random key is ok, but it should change when the actual signing key changes.

{{< hint warning >}}
* TODO: the port in the `--iss` needs to be changed when hiding the idp behind the proxy
* TODO: the signing keys and encryption keys should be precerated so they are reused between restarts. Otherwise all client sessions will become invalid when restarting the IdP.
{{< /hint >}}


#### Check it is up and running

1. Try getting the configuration:
```console
$ curl http://127.0.0.1:9130/.well-known/openid-configuration
```

2. Check if the login works at http://127.0.0.1:9130/signin/v1/identifier

{{< hint >}}
If you later get a `Unable to find a key for (algorithm, kid):PS256, )` Error make sure you did set a `--signing-kid` when starting `ocis/bin/ocis idp` by checking it is present in http://127.0.0.1:9130/konnect/v1/jwks.json
{{< /hint >}}

### Start ocis proxy


{{< hint >}}
Everything below this hint is outdated. Next steps are roughly:
* directly after glauth start the `ocis storage-userporvider`?
  - how to verify that works?
  - https://github.com/fullstorydev/grpcurl
* start proxy
  - the ocis ipd url can be changed to https
  - when do we hide oc10 behind ocis? -> advanced bridge at the end? for now run it without touching the existing oc10 instance
* start web
  - verify the login works, but how?
    - TODO the login works, but then the capabilities requests will fail ... unless we make the proxy answer them by talking to oc10?

Other ideas:
* the owncloud backend in glauth also works with the user provisioning api ... no changes to a running production instance? db access could be done with a read only account as well...
{{< /hint >}}


### Start ocis-web

#### Run it!

Point `ocis-web` to your owncloud domain and tell it where to find the openid connect issuing authority:
```console
$ bin/web server --web-config-server https://cloud.example.com --oidc-authority https://192.168.1.100:9130 --oidc-metadata-url https://192.168.1.100:9130/.well-known/openid-configuration --oidc-client-id ocis
```

`ocis-web` needs to know
- `--web-config-server https://cloud.example.com` is ownCloud url with webdav and ocs endpoints (oc10 or ocis)
- `--oidc-authority https://192.168.1.100:9130` the openid connect issuing authority, in our case `oidc-idp`, running on port 9130
- `--oidc-metadata-url https://192.168.1.100:9130/.well-known/openid-configuration` the openid connect configuration endpoint, typically the issuer host with `.well-known/openid-configuration`, but there are cases when another endpoint is used, eg. ping identity provides multiple endpoints to separate domains
- `--oidc-client-id ocis` the client id we will register later with `ocis-idp` in the `identifier-registration.yaml`

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
    'client-id' => 'ocis',
    'loginButtonName' => 'OpenId Connect @ Konnectd',
  ],
  'debug' => true, // if using self signed certificates
  // allow the different domains access to the ocs and webdav endpoints:
  'cors.allowed-domains' => [
    'https://cloud.example.com',
    'http://localhost:9100',
  ],
];
```

In the above configuration replace
- `provider-url` with the URL to your `ocis-idp` issuer
- `https://cloud.example.com` with the URL to your ownCloud 10 instance
- `http://localhost:9100` with the URL to your ownCloud Web instance

> Note: By default the openidconnect app will use the email of the user to match the user from the oidc userinfo endpoint with the ownCloud account. So make sure your users have a unique primary email.

## Next steps

Aside from the above todos these are the next steps
- tie it all together behind `ocis-proxy`
- create an `ocis bridge` command that runs all the ocis services in one step with a properly preconfigured `ocis-idp` `identifier-registration.yaml` file for `ownCloud Web` and the owncloud 10 `openidconnect` app, as well as a randomized `--signing-kid`.
