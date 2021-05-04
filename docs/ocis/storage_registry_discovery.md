---
title: "Storage Registry Discovery"
date: 2021-05-04T14:01:00+01:00
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: storage_registry_discovery.md
---

## Introduction

In order for an oCIS client to access oCIS storage spaces for an End-User, the client needs to know where the oCIS instance is. oCIS uses WebFinger [RFC7033](http://tools.ietf.org/html/rfc7033) to locate the oCIS instance for an End-User.

This discovery is optional. If the client has another way of discovering the OpenID instance, e.g. when looging in with a username a static domain might be configured or the domain in the URL might be used.

For guest accounts that do not have an OIDC issuer or whose IdP is not part of a trusted federation clients may fall back to a local IdP.

## User Input using E-Mail Address Syntax

To find the oCIS instance for the given user input in the form of an e-mail address `joe@example.com`, the WebFinger parameters are as follows:

| WebFinger Parameter | Value |
|-|-|
| `resource` | `acct:joe@example.com` |
| `host` | `example.com` |
| `rel` | http://owncloud.com/specs/ocis/1.0/instance |

Note that in this case, [the `acct:` scheme](http://tools.ietf.org/html/draft-ietf-appsawg-acct-uri-07) is prepended to the identifier.

The client (relying party) would make the following WebFinger request to discover the oCIS instance location (with line wraps within lines for display purposes only):

```
  GET /.well-known/webfinger
    ?resource=acct%3Ajoe%40example.com
    &rel=http%3A%2F%2Fowncloud.com%2Fspecs%2Focis%2F1.0%2Finstance
    HTTP/1.1
  Host: example.com

  HTTP/1.1 200 OK
  Content-Type: application/jrd+json

  {
   "subject": "acct:joe@example.com",
   "links":
    [
     {
      "rel": "http://owncloud.com/specs/ocis/1.0/instance",
      "href": "https://cloud.example.com"
     }
    ]
  }
```

{{< hint >}}
Note: the `example.com` domain is derived from the email.
{{< /hint >}}

{{< hint danger >}}
The `https://cloud.example.com` domain above would point to the ocis instance. 
TODO that ins ocis web ... not the registry ... hmmmm
maybe introduce an ocis provider which then has an `/.well-known/ocis-configuration`, similar to `/.well-known/openid-configuration`?
It would contain
- the ocis domain, e.g. `https://cloud.example.com`
- the web endpoint, e.g. `https://cloud.example.com`
- the registry / drives endpoint, e.g. `https://cloud.example.com/graph/v0.1/drives/me` see [Add draft of adr for spaces API. #1827](https://github.com/owncloud/ocis/pull/1827)


example:
```
HTTP/1.1 200 OK
  Content-Type: application/json

  {
   "instance":        "https://cloud.example.com/",
   "graph_endpoint":  "https://cloud.example.com/graph/",
   "ocis_web_config": "https://cloud.example.com/web/config.json",
   "issuer":          "https://idp.example.com/",
  }
```

`graph_endpoint` is the open-graph-api endpoint that is used to list storage spaces at e.g. `https://cloud.example.com/graph/v0.1/me/drives`.

`ocis_web_config` points ocis web to the config for the instance. Maybe we can add more config in the `/.well-known/ocis-configuration` to replace the config.json? Is this the new status.php? How safe is it to expose all this info ...?

The `issuer` could be used to detect the issuer that is used if no other issuer is found ... might be a fallback_issuer, but actually we may decide to skid the OIDC discovery and rely on this property. Maybe we need it if no IdP is present yet or the `/.well-known/openid-configuration` is not set up / reachable.


{{< /hint >}}

## Obtaining oCIS Provider Configuration Information
Using the `instance` location discovered as described above or by other means, the oCIS Provider's configuration information can be retrieved.

oCIS Providers supporting Discovery MUST make a JSON document available at the path formed by concatenating the string `/.well-known/openid-configuration` to the `instance`. The syntax and semantics of `.well-known` are defined in [RFC5785](http://tools.ietf.org/html/rfc5785) and apply to the `instance` value when it contains no path component. `ocis-configuration` MUST point to a JSON document compliant with this specification and MUST be returned using the `application/json` content type.

### oCIS Provider Configuration Request

An oCIS Provider Configuration Document MUST be queried using an HTTP GET request at the previously specified path.

The client (relying party) would make the following request to the instance https://example.com to obtain its Configuration information, since the Issuer contains no path component:

  GET /.well-known/openid-configuration HTTP/1.1
  Host: example.com
If the Issuer value contains a path component, any terminating / MUST be removed before appending /.well-known/openid-configuration. The RP would make the following request to the Issuer https://example.com/issuer1 to obtain its Configuration information, since the Issuer contains a path component:

  GET /issuer1/.well-known/openid-configuration HTTP/1.1
  Host: example.com
Using path components enables supporting multiple issuers per host. This is required in some multi-tenant hosting configurations. This use of .well-known is for supporting multiple issuers per host; unlike its use in RFC 5785 [RFC5785], it does not provide general information about the host.