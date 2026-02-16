---
title: "27. New Share Jail"
date: 2024-02-21T15:19:00+01:00
weight: 27
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0027-new-share-jail.md
---

* Status: draft
* Deciders: [@butonic](https://github.com/butonic), [@rhafer](https://github.com/rhafer), [@dragotin](https://github.com/dragotin)
* Date: 2024-02-21

## Context and Problem Statement

The oCIS share jail is a space that contains all accepted / synced shares of a user. In contrast to a personal or project space that contains actual resources, the share jail space only contains references pointing to shared resources. The root directory only consists of mountpoints that actually represent resources in other spaces. On the WebDAV API clients expect an `oc:fileid` property to identify resources in other API endpoints, eg. the libregraph `/me/sharedWithMe` endpoint. 

Currently, we construct the `oc:fileid` from the pattern `{shareproviderid}${sharespaceid}!{sharemountid}`. `{shareproviderid}`and `{sharespaceid}` are both hardcoded to `a0ca6a90-a365-4782-871e-d44447bbc668`. The `{sharemountid}` itself uses the pattern `{shared-resource-providerid}:{shared-resource-spaceid}:{shareid}`.

Since a resource can be shared to the same user in multiple ways (a group share and a user share) we deduplicate the two shares and only show one mountpoint in the share jail root. This is where this solution starts to fall apart:
* When accepting, mounting or syncing a share we implicitly have to accept all shares
* Each share has a different `{shareid}`, so we currently look up the oldest share and use it to build the `oc:fileid`
* Consequently, when the oldest share is revoked the `oc:fileid` changes.

We need to build the `oc:fileid` from a more stable pattern.

### Shareid

The WebDAV PROPFIND response also contains a `oc:shareid` which currently is derived from the path when the spaceid matches the share jail. The jsoncs3 implementation of the share manager currently is the only one using the `{shared-resource-providerid}:{shared-resource-spaceid}:{shareid}` pattern, where `{shareid}` is a uuid that is generated when creating the share. 

Again, the problem is that a resource can be shared multiple times.

## Decision Drivers <!-- optional -->

* We need to change the `oc:fileid` pattern without breaking clients.
* We need to be able to correlate files from WebDAV and the Graph API.

## Considered Options

* [Share based id](#share-based-id)
* [Resource based id](#resource-based-id)
* [Permission based id](#permission-based-id)
* [Use graph for file metadata](#use-graph-for-file-metadata)

## Decision Outcome

Resource based id: it correctly reflects the semantic meaning of a mount point, by indirectly pointing to the resource, not the share. The permissions on the share have to be checked in the storageprovider itself, anyway. Switching to graph requires more effort and the transition can happen gradually ofter changing the `oc:fileid` pattern in the sharejail.

### Positive Consequences:

* We get rid of mixing share ids with fileids, preventing unexpected `oc:fileid` changes.

### Negative Consequences:

* We need to teach clients about a new share jail space that uses the new `oc:fileid` pattern. They may need to implement a migration strategy to switch from the old share jail space to a new share jail space by replacing the fileid in their internal database. The might be able to just switch over, because the only `oc:fileid` that changes is the one from the mountpoints. The other nodes in the subtree already use the resourceid of the shared resource.
* Clients relying on `oc:shareid` to correlate share jail entries in PROPFIND responses need to either deal with multiple `oc:shareid` as a resource can be shared multiple times, or we deprecate `oc:shareid` and only use the `oc:fileid`. *jfd: Who is using this? why? Please explain and add to the decision drivers above!*
* The graph api also needs to be able to list entries from the new share jail. *jfd: clients could use a filter to ask for the new share jail id*

## Pros and Cons of the Options <!-- optional -->

### Share based id
Follow the pattern `{shareproviderid}${sharespaceid}!{sharemountid}`, where `{sharemountid}` is `{shared-resource-providerid}:{shared-resource-spaceid}:{shareid}`.
Combined patter `{shareproviderid}${sharespaceid}!{shared-resource-providerid}:{shared-resource-spaceid}:{shareid}`.
`{shareproviderid}` and `{sharespaceid}`are both hardcodet to `a0ca6a90-a365-4782-871e-d44447bbc668` to route all id based requests for mountpoints to the share jail space.

+ Good, the `{shared-resource-providerid}` and `{shared-resource-spaceid}` are used to shard the shares per space.
- Bad, `oc:fileid` changes if the oldest received share to a resource is revoked.

### Resource based id
Follow the pattern `{shareproviderid}${sharespaceid}!{shared-resource-providerid}:{shared-resource-spaceid}:{shared-resource-opaqueid}`.
Hardcode `756e6cdf-5630-4b66-9380-55a85188e0f6` as a new `{sharespaceid}` to allow clients to detect the new share jail and change it at their own pace.

+ Good, stable `oc:fileid` that remains the same, regardless of permission changes.
+ Good, clients can detect the new share jail and deal with it on their terms.

### Permission based id
Follow the pattern `{shareproviderid}${sharespaceid}!{shared-resource-providerid}:{shared-resource-spaceid}:{shared-resource-opaqueid}:{permission-id}`.

- Bad, same instability as the share id
- Bad, we don't even have a permission id. We could construct one from the grantee, but this leads nowhere.


### Use graph for file metadata
Instead of using WebDAV to correlate files with shares fully embrace libregraph to manage file metadata. 
Follow the pattern `{shareproviderid}${sharespaceid}!{shared-resource-providerid}:{shared-resource-spaceid}:{shared-resource-opaqueid}`.
WebDAV can be stripped of any ownCloud specific properties and will only be used for file up and download.

- Bad, more effort
+ Good, clean way of representing mountpoints and the shared resource in one `driveItem` that can include the resource based id.
+ Good, pagination, sorting and filtering cleanly specified
+ Good, WebDAV can be stripped down.
+ Good, Clients could get rid of WebDAV client and XML libs as all endpoints use JSON (all OCS endpoins return JSON when appending a `format=json` query parameter)

## Links <!-- optional -->
