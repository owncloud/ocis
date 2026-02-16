---
title: "22. Sharing and Space Management API"
date: 2023-09-08T02:29:00+01:00
weight: 22
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0022-sharing-and-space-management-api.md
---

* Status: accepted
* Deciders: [@JammingBen](https://github.com/JammingBen) [@butonic](https://github.com/butonic) [@theonering](https://github.com/theonering) [@kobergj](https://github.com/kobergj) [@micbar](https://github.com/micbar)
* Date: 2023-08-08

Technical Story: [Public issue](https://github.com/owncloud/ocis/issues/6993)

## Context and Problem Statement

In the early days of the rewrite of ownCloud it was an important goal to keep all important APIs compatible with ownCloud 10. Infinite Scale embraced that goal until version 1.0.0.

After that first release, the focus changed.

Infinite Scale started the spaces feature which brings a whole new set of APIs and concepts. We made the conscious decision to keep the sharing API as it was, live with its shortcomings and create workarounds to support spaces. We have come a long way so far. Now we need to move on. The Web Client has made the decision to drop the support of ownCloud 10 and keep version 7.0 alive for ownCloud 10 to keep the easy migration path intact.

The desktop and mobile client platforms were suffering from poor support from the server and can now move forward with a new API implementation. By using openApi 3 and all the needed tooling around it developing the LibreGraph specification, documentaion and SDKs, we now feel confident to move on.

## Decision Drivers

* The Path based nature of the OCS API lacks spaces support
* The permissions bitmask is no longer working when using sharing roles
* We want to support server announced sharing roles which are different per instance or scope
* We need to get rid of the currently hardcoded sharing roles in our clients
* New sharing roles and permissions are needed to support secure view and other new features
* Space Memberships are not shares and need to have different semantics
* Elevation of permissions in subfolders or full denials should be possible without creating a new share
* Third party integrations need generated SDKs in different languages to speed up the development

## Considered Options

* [New OCS Api Version](#new-ocs-api-version)
* [Sharing via LibreGraph](#sharing-via-libregraph)

## Decision Outcome

Chosen option: "[LibreGraph](#sharing-via-libregraph)"

### Positive Consequences:

* We can create a new clean API which fits the spaces concept
* LibreGraph embraces OData which is a known API pattern
* Sharing will be integrated in the existing SDKs and documentation
* Removing the OCS Api reduces complexity
* Removing the OCS Api makes the clients codebases smaller and removes manually maintained parts of the SDKs
* The extra error handling for the OCS API can be dropped from our clients

### Negative Consequences:

* We need to deprecate and remove the OCS API
* Existing third party integrations need to do some refactoring

## Pros and Cons of the Options

### New OCS Api Version

To overcome the limitations of the OCS 2.0 API we could create a new major version with the spaces concept in mind. This would give us the opportunity to create a new openApi Spec.

* Good, because the workarounds from version 2.0 could be dropped
* Bad, because we would need to deprecate the version 2.0
* Bad, because we would need to maintain a separate specification / repository
* Bad, because it would create the need to use two different SDKs in our clients
* Bad, because we would need to implement query parameters and filters on our own
* Bad, because sharing information could not be included in the spaces API via queries or filters

### Sharing via LibreGraph

Integrate Sharing into the [LibreGraph API](https://github.com/owncloud/libre-graph-api) by using the already existing toolchain and documentation flows.

* Good, because that reduces the number of SDKs
* Good, because it reduces the number of APIs
* Good, because spaces and shares can be used together in queries and filters
* Good, because we would use the existing OData pattern
* Bad, because we need to deprecate the OCS API

## Links <!-- optional -->

* [LibreGraph API](https://github.com/owncloud/libre-graph-api)
* [OData](https://www.odata.org/documentation/)
* [OpenAPI Standard](https://www.openapis.org/)
