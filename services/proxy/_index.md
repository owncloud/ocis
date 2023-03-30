---
title: Proxy Service
date: 2023-03-30T08:57:59.991150591Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/proxy
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The proxy service is an API-Gateway for the ownCloud Infinite Scale microservices. Every HTTP request goes through this service. Authentication, logging and other preprocessing of requests also happens here. Mechanisms like request rate limitting or intrusion prevention are **not** included in the proxy service and must be setup in front like with an external reverse proxy.
The proxy service is the only service communicating to the outside and needs therefore usual protections against DDOS, Slow Loris or other attack vectors. All other services are not exposed to the outside, but also need protective measures when it comes to distributed setups like when using container orchestration over various physical servers.

## Table of Contents

* [Authentication](#authentication)
* [Automatic Quota Assignments](#automatic-quota-assignments)
* [Automatic Role Assignments](#automatic-role-assignments)
* [Recommendations for Production Deployments](#recommendations-for-production-deployments)
* [Example Yaml Config](#example-yaml-config)

## Authentication

The following request authentication schemes are implemented:
-   Basic Auth (Only use in development, **never in production** setups!)
-   OpenID Connect
-   Signed URL
-   Public Share Token

## Automatic Quota Assignments

It is possible to automatically assign a specific quota to new users depending on their role.
To do this, you need to configure a mapping between roles defined by their ID and the quota in bytes.
The assignment can only be done via a `yaml` configuration and not via environment variables.
See the following `proxy.yaml` config snippet for a configuration example.
```yaml
role_quotas:
    <role ID1>: <quota1>
    <role ID2>: <quota2>
```

## Automatic Role Assignments

When users login, they do automatically get a role assigned. The automatic role assignment can be
configured in different ways. The `PROXY_ROLE_ASSIGNMENT_DRIVER` environment variable (or the `driver`
setting in the `role_assignment` section of the configuration file select which mechanism to use for
the automatic role assignment.
When set to `default`, all users which do not have a role assigned at the time for the first login will
get the role 'user' assigned. (This is also the default behavior if `PROXY_ROLE_ASSIGNMENT_DRIVER`
is unset.
When `PROXY_ROLE_ASSIGNMENT_DRIVER` is set to `oidc` the role assignment for a user will happen
based on the values of an OpenID Connect Claim of that user. The name of the OpenID Connect Claim to
be used for the role assignment can be configured via the `PROXY_ROLE_ASSIGNMENT_OIDC_CLAIM`
environment variable. It is also possible to define a mapping of claim values to role names defined
in ownCloud Infinite Scale via a `yaml` configuration. See the following `proxy.yaml` snippet for an
example.
```yaml
role_assignment:
    driver: oidc
    oidc_role_mapper:
        role_claim: ocisRoles
        role_mapping:
            admin: myAdminRole
            user: myUserRole
            spaceadmin: mySpaceAdminRole
            guest: myGuestRole
```
This would assign the role `admin` to users with the value `myAdminRole` in the claim `ocisRoles`.
The role `user` to users with the values `myUserRole` in the claims `ocisRoles` and so on.
Claim values that are not mapped to a specific ownCloud Infinite Scale role will be ignored.
Note: An ownCloud Infinite Scale user can only have a single role assigned. If the configured
`role_mapping` and a user's claim values result in multiple possible roles for a user, an error
will be logged and the user will not be able to login.
The default `role_claim` (or `PROXY_ROLE_ASSIGNMENT_OIDC_CLAIM`) is `roles`. The `role_mapping` is:
```yaml
admin: ocisAdmin
user: ocisUser
spaceadmin: ocisSpaceAdmin
guest: ocisGuest
```

## Recommendations for Production Deployments

In a production deployment, you want to have basic authentication (`PROXY_ENABLE_BASIC_AUTH`) disabled which is the default state. You also want to setup a firewall to only allow requests to the proxy service or the reverse proxy if you have one. Requests to the other services should be blocked by the firewall.

## Example Yaml Config

{{< include file="services/_includes/proxy-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/proxy_configvars.md" >}}

