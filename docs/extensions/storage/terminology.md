---
title: "Terminology"
date: 2018-05-02T00:00:00+00:00
weight: 17
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: terminology.md
---

Communication is hard. And clear communication is even harder. You may encounter the following terms throughout the documentation, in the code or when talking to other developers. Just keep in mind that whenever you hear or read *storage*, that term needs to be clarified, because on its own it is too vague. PR welcome.

## Logical concepts

### Resources
A *resource* is the basic building block that oCIS manages. It can be of [different types](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceType):
- an actual *file*
- a *container*, e.g. a folder or bucket
- a *symlink*, or
- a [*reference*]({{< ref "#references" >}}) which can point to a resource in another [*storage provider*]({{< ref "#storage-providers" >}})

### References
A *reference* identifies a [*resource*]({{< ref "#resources" >}}). A [*CS3 reference*](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Reference) can carry a *path* and a [CS3 *resource id*](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId). The references come in two flavors: absolute and combined.
Absolute references have either the *path* or the *resource id* set:
- An absolute *path* MUST start with a `/`. The *resource id* MUST be empty.
- An absolute *resource id* uniquely identifies a [*resource*]({{< ref "#resources" >}}) and is used as a stable identifier for sharing. The *path* MUST be empty.
Combined references have both, *path* and *resource id* set:
- the *resource id* identifies the root [*resource*]({{< ref "#resources" >}})
- the *path* is relative to that root. It MUST start with `.`

## Technical concepts

### Storage Systems
Every *storage system* has different native capabilities like id and path based lookups, recursive change time propagation, permissions, trash, versions, archival and more.
A [*storage provider*]({{< ref "#storage-providers" >}}) makes the storage system available in the CS3 API by wrapping the capabilities as good as possible using a [*storage driver*]({{< ref "./storagedrivers.md" >}}).
There might be multiple [*storage drivers*]({{< ref "./storagedrivers.md" >}}) for a *storage system*, implementing different tradeoffs to match varying requirements.

### Gateways
A *gateway* acts as a facade to the storage related services. It authenticates and forwards API calls that are publicly accessible.

