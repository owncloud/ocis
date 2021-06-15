---
title: "8. Extension Template"
---

* Status: proposed
* Deciders: @c0rby <!-- optional -->
* Date: 2021-06-10

Technical Story: [description | ticket/issue URL] <!-- optional -->

## Context and Problem Statement

We want to accelerate and simplify extension development by removing the necessity to type or copy the boilerplate code. Can we provide a template or a similar mechanism to aid when developing new extensions?


## Decision Drivers <!-- optional -->

* The solution should be easily maintainable.
  * It should always be up-to-date.
* The solution should be easy to use.

## Considered Options

* Use [boilr](https://github.com/tmrts/boilr)
* Create a template git repository.
* Use [ocis-hello](https://github.com/owncloud/ocis-hello/) as a "template"

## Decision Outcome

Chosen option: "[option 1]", because [justification. e.g., only option, which meets k.o. criterion decision driver | which resolves force force | … | comes out best (see below)].

### Positive Consequences: <!-- optional -->

* [e.g., improvement of quality attribute satisfaction, follow-up decisions required, …]
* …

### Negative consequences: <!-- optional -->

* [e.g., compromising quality attribute, follow-up decisions required, …]
* …

## Pros and Cons of the Options <!-- optional -->

### [boilr](https://github.com/tmrts/boilr)

We have a boilr template already. [boilr-ocis-extension](https://github.com/owncloud/boilr-ocis-extension/)
This approach is nice because it provides placeholders which can be filled during the generation of a new extension from the template. It also provides prompts for the placeholder values during generation.

* Good, because with the placeholders it is hard to miss values which should be changed
* Bad, because maintaining is more complex

### Template git repository

Create a git repository with an extension containing the boilerplate code.

* Good, because we can use the usual tools for QA and dependency scanning/updating.
* Good, because it doesn't require any additional tool.

### [ocis-hello](https://github.com/owncloud/ocis-hello/) as a "template"

We have the ocis-hello repository which acts as an example extension containing a grpc and http service and a web UI. It also demonstrates the usage of the settings service.

* Good, because it contains a bit more code than just the plain boilerplate
* Good, because the integration into oCIS is already tested for the Hello extension (eg. with Proxy and Settings). This will ensure, that the example extension is up to date.
* Bad, because if you don't require all features you have to delete stuff

