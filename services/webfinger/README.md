# Webfinger

The webfinger service provides an RFC7033 WebFinger lookup of ownCloud instances relevant for a given user account via endpoints a the /.well-known/webfinger implementation.

It is based on https://github.com/owncloud/lookup-webfinger-sciebo but also returns localized `titles` in addition to the `href` property.

## OpenID Connect Discovery

Clients can make an unauthenticated `GET https://drive.ocis.test/.well-known/webfinger?resource=https%3A%2F%2Fdrive.ocis.test` request to discover the OpenID Connect Issuer in the `http://openid.net/specs/connect/1.0/issuer` relation:

```json
{
    "subject": "https://drive.ocis.test",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer",
            "href": "https://sso.example.org/cas/oidc/"
        }
    ]
}
```

Here, the `resource` takes the instance domain URI, but an `acct:` URI works as well. 

## Authenticated Instance Discovery

When using OpenID connect to authenticate requests, clients can look up the owncloud instances a user has access to.

*   Authentication is necessary to prevent leaking information about existing users.
*   Basic auth is not supported.

The default configuration will simply return the `OCIS_URL` and direct clients to that domain:

```json
{
    "subject": "acct:einstein@drive.ocis.test",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer",
            "href": "https://sso.example.org/cas/oidc/"
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://abc.drive.example.org",
            "titles": {
                "en": "oCIS Instance"
            }
        }
    ]
}
```

## Configure Different Instances Based on OpenidConnect UserInfo Claims

A more complex example for configuring different instances could look like this:

```yaml
webfinger:
  instances:
  -  claim: email
     regex: einstein@example\.org
     href: "https://{{.preferred_username}}.cloud.ocis.test"
     title: 
       "en": "oCIS Instance for Einstein"
       "de": "oCIS Instanz für Einstein"
     break: true
  -  claim: "email"
     regex: marie@example\.org
     href: "https://{{.preferred_username}}.cloud.ocis.test"
     title: 
       "en": "oCIS Instance for Marie"
       "de": "oCIS Instanz für Marie"
     break: false
  -  claim: "email"
     regex: .+@example\.org
     href: "https://example-org.cloud.ocis.test"
     title:
       "en": "oCIS Instance for example.org"
       "de": "oCIS Instanz für example.org"
     break: true
  -  claim: "email"
     regex: .+@example\.com
     href: "https://example-com.cloud.ocis.test"
     title:
       "en": "oCIS Instance for example.com"
       "de": "oCIS Instanz für example.com"
     break: true
  -  claim: "email"
     regex: .+@.+\..+
     href: "https://cloud.ocis.test"
     title:
       "en": "oCIS Instance"
       "de": "oCIS Instanz"
     break: true
```

Now, an authenticated webfinger request for `acct:me@example.org` (when logged in as marie) would return two instances, based on her `email` claim, the regex matches and break flags:

```json
{
    "subject": "acct:marie@example.org",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer",
            "href": "https://sso.example.org/cas/oidc/"
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://marie.cloud.ocis.test",
            "titles": {
                "en": "oCIS Instance for Marie",
                "de": "oCIS Instanz für Marie"
            }
        },
        {
            "rel": "http://webfinger.owncloud/rel/server-instance",
            "href": "https://xyz.drive.example.org",
            "titles": {
                "en": "oCIS Instance for example.org",
                "de": "oCIS Instanz für example.org"
            }
        }
    ]
}
```
