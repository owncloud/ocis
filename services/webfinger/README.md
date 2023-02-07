# Webfinger service

The webfinger service provides an RFC7033 WebFinger lookup of ownCloud instances relevant for a given user account.

It is based on https://github.com/owncloud/lookup-webfinger-sciebo but also returns localized `titles` in addition to the `href` property.

A client can discover various endpoints by making a `GET` request to `https://drive.ocis.test/.well-known/webfinger?resource=acct%3Aeinstein%40cloud.ocis.test`.

He will get a response like
```
{
    "subject": "acct:einstein@drive.ocis.test",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer"
            "href": "https://sso.example.org/cas/oidc/",
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://abc.drive.example.org"
    	    "titles": {
    	      "en": "Readable Instance Name"
    	    }
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://xyz.drive.example.org"
    	    "titles": {
    	      "en": "Readable Other Instance Name"
    	    }
        },
    ]
}
```

As a special form of lookup clients can use the relation `acct:me@{host}` to look up the oCIS specific user id, username, email and displayname.

In this case there are two ocis instances and the client has to ask the user which instance he wants to use.

## TODO
Currently, clients need to make subsequent calls to:
- /status.php to check if the instance is in maintenance mode or if the version is supported
- /config.json to get the available apps for ocis web to determine which routes require authentication
- /themes/owncloud/theme.json for theming info
- /.well-known/openid-configuration, auth2 token and oidc userinfo endpoints to authenticate the user
- /ocs/v1.php/cloud/user to get the username, eg. einstein ... again? it contains the oc10 user id (marie, not the uuid)
- /ocs/v1.php/cloud/capabilities to fetch instance capabilites
- /ocs/v1.php/cloud/users/einstein to fetch the quota which could come from graph and actually is now tied to the spaces, not tu users
- /graph/v1.0/me?%24expand=memberOf to fetch the user id and the groups the user is a member of

We need a way to pass oidc claims from the proxy, which does the authentication to the webfinger service, preferably by minting them into the internal reva token.
- Currently, we use machine auth so we can autoprovision an account if it does not exist. We should use revas oidc auth and, when autoprovisioning is enabled, retry the authentication after provisioning the account. This would allow us to use a `roles` claim to decide which roles to use and eg. a `school` claim to determine a specific instance. We may use https://github.com/PerimeterX/marshmallow to parse the RegisteredClaims and get the custom claims as a separate map.

For now, webfinger can only match users based on a regex and produce a list of instances based on that.

Here are some Ideas which need to be discussed with all client teams in the future:

### status properties

The /.well-known/webfinger enpdoint allows us to not only get rid of some of these calls, e.g. by embedding status.php info:

```
{
    "subject": "acct:einstein@drive.ocis.test",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer"
            "href": "https://idp.ocis.test",
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://abc.drive.ocis.test"
    	    "titles": {
    	      "en": "Readable Instance Name"
    	    }
            "properties": {
                "http://webfinger.owncloud/prop/status/maintenance": "false",
                "http://webfinger.owncloud/prop/status/version": "10.0.11.3"
            }
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://xyz.drive.ocis.test"
    	    "titles": {
    	      "en": "Readable Other Instance Name"
    	    }
            "properties": {
                "http://webfinger.owncloud/prop/status/maintenance": "true",
                "http://webfinger.owncloud/prop/status/version": "10.0.11.2"
            }
        },
    ]
}
```

### Dedicated ocis web endpoint

It also allows us to move some services out of a sharded deployment. We could e.g. introduce a relation for a common ocis web endpoint to not exponse the different instances in the browser bar:
```
{
    "subject": "acct:einstein@drive.ocis.test",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer"
            "href": "https://idp.ocis.test",
        },
        {
            "rel": "http://webfinger.owncloud/rel/web"
            "href": "https://drive.ocis.test",
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://abc.drive.ocis.test"
    	    "titles": {
    	      "en": "Readable Instance Name"
    	    }
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://xyz.drive.ocis.test"
    	    "titles": {
    	      "en": "Readable Other Instance Name"
    	    }
        },
    ]
}
```

### Dedicated ocis web endpoint

We could also omit the `http://webfinger.owncloud/rel/server-instance` relation and go straight for a graph service with e.g. `rel=http://libregraph.org/rel/graph`:
```
{
    "subject": "acct:einstein@drive.ocis.test",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer"
            "href": "https://idp.ocis.test",
        },
        {
            "rel": "http://webfinger.owncloud/rel/web"
            "href": "https://drive.ocis.test",
        },
        {
            "rel": "http://libregraph.org/rel/graph",
            "href": "https://abc.drive.ocis.test/graph/v1.0"
    	    "titles": {
    	      "en": "Readable Instance Name"
    	    }
        },
    ]
}
```

In theory the graph endpoint would allow discovering drives on any domain. But there is a lot more work to be done here.

### Subject properties

We could also embed subject metadata, however since apps like ocis web also need the groups a user is member of a dadicated call to the libregraph api is probably better. In any case, we could return properties for the subject:
```
{
    "subject": "acct:einstein@drive.ocis.test",
    "properties": {
        "http://libregraph.org/prop/user/id": "4c510ada-c86b-4815-8820-42cdf82c3d51",
        "http://libregraph.org/prop/user/onPremisesSamAccountName": "einstein",
        "http://libregraph.org/prop/user/mail": "einstein@example.org",
        "http://libregraph.org/prop/user/displayName": "Albert Einstein",
    }
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer"
            "href": "https://idp.ocis.test",
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://abc.drive.ocis.test"
    	    "titles": {
    	      "en": "Readable Instance Name"
    	    }
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://xyz.drive.ocis.test"
    	    "titles": {
    	      "en": "Readable Other Instance Name"
    	    }
        },
    ]
}
```

# status php

```
{
    "subject": "https://drive.ocis.test",
    "properties": {
        "http://webfinger.owncloud/prop/maintenance": "false",
        "http://webfinger.owncloud/prop/version": "10.11.0.6"
    }
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer"
            "href": "https://idp.ocis.test",
        },
    ]
}
```