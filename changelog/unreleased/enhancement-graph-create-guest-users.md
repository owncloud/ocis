Enhancement: Allow creating guest users through the provisioning API

`POST /graph/v1.0/users` now accepts an explicit `userType` (e.g. `Guest`) on
creation instead of always forcing `Member` and rejecting any supplied type. This
aligns the provisioning API with the user types the system already recognizes
(`Member` and `Guest`).

Backwards compatible: omitting `userType` still defaults to `Member`, so existing
API consumers are unaffected. An unknown `userType` is rejected with a 400. Creating
users remains gated by the admin requirement on the route.

A `Guest` created this way is assigned the lightweight guest role at creation time,
matching the role a `Guest` would otherwise only receive from the proxy on first
login. A `userType` is also normalized to its canonical `Member`/`Guest` casing
before being stored, regardless of how it was cased in the request.

https://github.com/owncloud/ocis/pull/12471
