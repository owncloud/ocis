---
title: "Storage"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

This service provides an oCIS extension that wraps [reva](https://github.com/cs3org/reva/) and adds an opinionated configuration to it.

## Architecture Overview

The below diagram shows the oCIS services and the contained reva services within as dashed boxes. In general:
1. A request comes in at the proxy and is authenticated using oidc.
2. It is forwarded to the oCIS frontend which handles ocs and ocdav requests by talking to the reva gateway using the CS3 API.
3. The gateway acts as a facade to the actual CS3 services: storage providers, user providers, group providers and sharing providers.

{{< svg src="extensions/storage/static/overview.drawio.svg" >}}

The dashed lines in the diagram indicate requests that are made to authenticate requests or lookup the storage provider:
1. After authenticating a request, the proxy may either use the CS3 `userprovider` or the accounts service to fetch the user information that will be minted into the `x-access-token`.
2. The gateway will verify the JWT signature of the `x-access-token` or try to authenticate the request itself, e.g. using a public link token.

{{< hint warning >}}
The bottom part is lighter because we will deprecate it in favor of using only the CS3 user and group providers after moving some account functionality into reva and glauth. The metadata storage is not registered in the reva gateway to seperate metadata necessary for running the service from data that is being served directly.
{{< /hint >}}

## Endpoints and references

In order to reason about the request flow, two aspects in the architecture need to be understood well:
1. The endpoints that are handling requests: what resources are presented at the available URL endpoints?
2. The resource identifiers that are exposed or required: path or id based?

### Frontend

The ocis frontend service starts all services that handle incoming HTTP requests:
- *ocdav* for ownCloud flavoured WebDAV
- *ocs* for sharing, user management, capabilities and other OCS API endpoints 
- *datagateway* for up and downloads
- TODO: *ocm*

{{< svg src="extensions/storage/static/frontend.drawio.svg" >}}

#### WebDAV

The ocdav service not only handles all WebDAV requests under `(remote.php/)(web)dav` but also some other legacy endpoints like `status.php`:

| endpoint | service | CS3 api | CS3 namespace | description | TODO |
|----------|---------|-------------|------|------|------|
| `status.php` | ocdav | - |  - | currently static | should return compiled version and dynamic values |
| `(remote.php/)webdav` | ocdav | storageprovider | `/home` | the old webdav endpoint |  |
| `(remote.php/)dav/files/<username>` | ocdav | storageprovider | `/home` | the new webdav endpoint |  |
| `(remote.php/)dav/meta/<fileid>/v` | ocdav | storageprovider | id based | versions |  |
| `(remote.php/)dav/trash-bin/<username>` | ocdav | recycle | - | trash | should aggregate the trash of storage spaces the user has access to |
| `(remote.php/)dav/public-files/<token>` | ocdav | storageprovider | `/public/<token>` | public links |  |
| `(remote.php/)dav/avatars/<username>` | ocdav | - | - | avatars, hardcoded | look up from user provider and cache |
| *CernBox setup:* |||||
| `(remote.php/)webdav` | ocdav | storageprovider | `/` | |  |
| *Note: existing folder sync pairs in legacy clients will break when moving the user home down in the path hierarchy* |||||
| `(remote.php/)webdav/home` | ocdav | storageprovider | `/home` |  |  |
| `(remote.php/)webdav/users` | ocdav | storageprovider | `/users` |  |  |
| `(remote.php/)dav/files/<username>` | ocdav | storageprovider | `/users/<userlayout>` |  |  |
| *Spaces concept:* |||||
| `(remote.php/)dav/(spaces|global)/<username>/<spaceid>` | ocdav | storageregistry & storageprovider | `spaceid`+`relative path` | spaces concept | allow listing and accessing spaces |

The correct endpoint for a users home storage space in oc10 is `remote.php/dav/files/<username>`. In oc10 All requests at this endpoint use a path based reference that is relative to the users home. In oCIS this can be configured and defaults to `/home` as well. Other API endpoints like ocs and the web UI still expect this to be the users home.

In oc10 we originally had `remote.php/webdav` which would render the current users home storage space. The early versions (pre OC7) would jail all received shares into a `remote.php/webdav/shares` subfolder. The semantics for syncing such a folder are [not trivially predictable](https://github.com/owncloud/core/issues/5349), which is why we made shares [freely mountable](https://github.com/owncloud/core/pull/8026) anywhere in the users home.

The current reva implementation jails shares into a `remote.php/webdav/Shares` folder for performance reasons. Obviously, this brings back the [special semantics for syncing](https://github.com/owncloud/product/issues/7). In the future we will follow [a different solution](https://github.com/owncloud/product/issues/302) and jail the received shares into a dedicated `/shares` space, on the same level as `/home` and `/spaces`. We will add a dedicated [API to list all *storage spaces*](https://github.com/owncloud/ocis/pull/1827) a user has access to and where they are mounted in the users *namespace*.

{{< hint warning >}}
Existing folder sync pairs in legacy clients will break when moving the user home down in the path hierarchy like CernBox did.
For legacy clients the `remote.php/webdav` endpoint will no longer list the users home directly, but instead present the different types of storage spaces:
- `remote.php/webdav/home`: the users home is pushed down into a new `home` *storage space*
- `remote.php/webdav/shares`: all mounted shares will be moved to a new `shares` *storage space*
- `remote.php/webdav/spaces`: other *storage spaces* the user has access to, e.g. group or project drives
{{< /hint >}}

{{< hint warning >}}
An alternative would be to introduce a new `remote.php/dav/spaces` or `remote.php/dav/global` endpoint. However, `remote.php/dav` properly follows the WebDAV RFCs strictly. To ensure that all resources under that namespace are scoped to the user the URL would have to include the principal like `remote.php/dav/spaces/<username>`, a precondition for e.g. WebDAV [RFC5397](https://tools.ietf.org/html/rfc5397). For a history lesson start at [Replace WebDAV with REST
owncloud/core#12504](https://github.com/owncloud/core/issues/12504#issuecomment-65218491) which spawned [Add extra layer in DAV to accomodate for other services like versions, trashbin, etc owncloud/core#12543](https://github.com/owncloud/core/issues/12543)
{{< /hint >}}


#### Sharing

The [OCS Share API](https://doc.owncloud.com/server/developer_manual/core/apis/ocs-share-api.html) endpoint `/ocs/v1.php/apps/files_sharing/api/v1/shares` returns shares, which have their own share id and reference files using a path relative to the users home. They API also lists the numeric storage id as well as the string type `storage_id` (which is confusing ... but yeah) which would allow constructing combined references with a storage spacle id and a path relative to the root of that storage space. The web UI however assumes that it can take the path from the `file_target` and append it to the users home to access it.

{{< hint >}}
The API [already returns the storage id](https://doc.owncloud.com/server/developer_manual/core/apis/ocs-share-api.html#example-request-response-payloads-4) (and numeric id) in addition to the file id:
```
    <storage_id>home::auser</storage_id>
    <storage>993</storage>
    <item_source>3994486</item_source>
    <file_source>3994486</file_source>
    <file_parent>3994485</file_parent>
    <file_target>/Shared/Paris.jpg</file_target>
``` 
[Creating shares only takes the **path** as the argument](https://doc.owncloud.com/server/developer_manual/core/apis/ocs-share-api.html#function-arguments) so creating and navigating shares only needs the path. When you update or delete a share it takes the `share id` not the `file id`.
{{< /hint >}}

The OCS service makes a stat request to the storage provider to get a [ResourceInfo](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceInfo) object. It contains both, a [`ResourceId`](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId) *and* an absolute path. If the *resource* exists a request is sent to the gateway. Depending on the type of share the [Collaboration API](https://cs3org.github.io/cs3apis/#cs3.sharing.collaboration.v1beta1.CollaborationAPI), the [Link API](https://cs3org.github.io/cs3apis/#cs3.sharing.link.v1beta1.LinkAPI) or the [Open Cloud Mesh API](https://cs3org.github.io/cs3apis/#cs3.sharing.ocm.v1beta1.OcmAPI) endpoints are used.

| API | Request | Resource identified by | Grant type | Further arguments |
|-----|---------|------------------------|------------|-------------------|
| Collaboration | [CreateShareRequest](https://cs3org.github.io/cs3apis/#cs3.sharing.collaboration.v1beta1.CreateShareRequest) | [ResourceInfo](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceInfo) | [ShareGrant](https://cs3org.github.io/cs3apis/#cs3.sharing.collaboration.v1beta1.ShareGrant) | - |
| Link | [CreatePublicShareRequest](https://cs3org.github.io/cs3apis/#cs3.sharing.link.v1beta1.CreatePublicShareRequest) | [ResourceInfo](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceInfo) | Link [Grant](https://cs3org.github.io/cs3apis/#cs3.sharing.link.v1beta1.Grant) | We send the public link `name` in the `ArbitraryMetadata` of the `ResourceInfo` |
| Open Cloud Mesh | [CreateOCMShareRequest](https://cs3org.github.io/cs3apis/#cs3.sharing.ocm.v1beta1.CreateOCMShareRequest) | [ResourceId](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId) | OCM [ShareGrant](https://cs3org.github.io/cs3apis/#cs3.sharing.ocm.v1beta1.ShareGrant) | OCM [ProviderInfo](https://cs3org.github.io/cs3apis/#cs3.ocm.provider.v1beta1.ProviderInfo) |


{{< hint >}}
The user and public share provider implementations identify the file using the `ResourceId`. The `ResourceInfo` is passed so the share provider can also store who the owner of the resource is. The *path* is not part of the other API calls, e.g. when listing shares.
The OCM API takes an id based reference on the CS3 api, even if the OCM HTTP endpoint takes a path argument. Why? Does it not need the owner? It only stores the owner of the share, which is always the currently looged in user, when creating a share. Afterwards only the owner can update a share ... so collaborative management of shares is not possible. At least for OCM shares.
{{< /hint >}}

#### User and Group provisioning

In oc10 users are identified by a username, which cannot change, because it is used as a foreign key in several tables. For oCIS we are internally identifying users by a UUID, while using the username in the WebDAV and OCS APIs for backwards compatability. To distinguish this in the URLs we are using `<username>` instead of `<userid>`. You may have encountered `<userlayout>`, which refers to a template that can be configuted to build several path segments by filling in user properties, e.g. the first two characters of the username or the issuer.

{{< hint warning >}}
Make no mistake, the [OCS Provisioning API](https://doc.owncloud.com/server/developer_manual/core/apis/provisioning-api.html) uses `userid` while it actually is the username, because it is what you use to login. 
{{< /hint >}}

We are currently working on adding [user management through the CS3 API](https://github.com/owncloud/ocis/pull/1930) to handle user and group provisioning (and deprovisioning).

## Terminology

Communication is hard. And clear communication is even harder. You may encounter the following terms throughout the documentation, in the code or when talking to other developers. Just keep in mind that whenever you hear or read *storage*, that term needs to be clarified, because on its own it is too vague. PR welcome.

### Resources
A *resource* is a logical concept. Ressources can be of [different types](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceType):
- an actual *file*
- a *container*, e.g. a folder or bucket
- a *symlink*, or
- a *reference* which can point to a resource in another *storage provider*

### References

A [*reference*](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.Reference) is a logical concept. It identifies a *resource* and consists of either
- a *path* based reference, used to identify a *resource* in the *namespace* of a *storage provider*. It must start with a `/`.
- an [*id* based reference](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourceId), uniquely identifying a *resource* in the *namespace* of a *storage provider*. It consists of a `storage provider id` and an `opaque id`. The `storage provider id` must NOT start with a `/`.

{{< hint info >}}
The `/` is important because currently the static *storage registry* uses a map to look up which *storage provider* is responsible for the resource. Paths must be prefixed with `/` so there can be no collisions between paths and storage provider ids in the same map.
{{< /hint >}}

{{< hint ok >}}
#### Alternative 1: root:(fileid|path) references ####

A *reference* is a logical concept. It identifies a *resource* and consists of
a *storage space* and an *opaqueid* or a *path* relative to the *root* of the storage space:
```
<storages_space>:<relative_path>
<storages_space>:<opaqueid>
```
In order to build a global namespace storage spaces can have aliases that the registry resolves into id based references. The ocdav service would have to look up the storage space that is responsible for an absolute path by talking to the registry (or the gateway transparently resolves absolute paths)


{{< /hint >}}

{{< hint warning >}}
#### Alternative 2: reference triple ####
A *reference* is a logical concept. It identifies a *resource* and consists of
a *root* and a *path* relative to that *root*. 
A *root* is a `storage space id` and a `logical uuid` of a resource inside the *storage space*:
```
<storages_space>:<root_id>:<relative_path>
`-----------root---------´ `----path-----´
```
*root* and *path* are both optional. The `storage space id` should be a uuid to prevent references from breaking when a *user* or *storage space* gets renamed. The following are all valid references:

or to clarify
```
<storages_space>:<root>:<relative_path>
```
Where actually only three cases are used:
1. `::<relative_path>` = absolute path
2. `<storages_space>::<relative_path>` = path relative to the root of the space
3. `<storages_space>:<root>:<relative_path>` = path relative to the specified node in the space


| Name | Description |
|------|-------------|
|`::/users/alice/projects/foo` | follow this path in the  global namespace (path based reference)|
|`home-alice::projects/foo` | in the *storage space* `home-alice`, start at the *root*, follow `projects/foo` |
|`home-alice:c3cf23bb-8f47-4719-a150-1d25a1f6fb56:foo` | in the *storage space* `home-alice`, start at the resource identified by the `logical uuid` `c3cf23bb-8f47-4719-a150-1d25a1f6fb56`, follow `foo` |
|`home-alice:56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:` | in the *storage space* `home-alice`, start at the resource identified by the `logical uuid` `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a`, no relativ path to follow |
|`ee1687e5-ac7f-426d-a6c0-03fed91d5f62:56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:` | in the *storage space* `ee1687e5-ac7f-426d-a6c0-03fed91d5f62`, start at the resource identified by the `logical uuid` `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a`, no relativ path to follow (id based reference) |

All those examples reference the same *resource* (if the `logical uuid` of `projects` is 
`c3cf23bb-8f47-4719-a150-1d25a1f6fb56` and the `logical uuid` of `foo` is `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a`).
To move *storage spaces* between *storage providers*, migrate files or restore backups, the `logical uuid` of a file must be writable.

A reference will often start as an absolute/global path, e.g. `::/home/Projects/Foo`. The gateway will look up the storage provider that is responsible for the path

| Name | Description | Who resolves it? |
|------|-------------|-|
| `::/home/Projects/Foo` | the absolute path a client like davfs will use. | The gateway uses the storage registry to look up the responsible storage provider |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62::/Projects/Foo` | the `storage space id` is the same as the `root`, the path becomes relative to the root | the storage provider can use this reference to identify this resource |

Now, the same file is accessed as a share
| Name | Description |
| `::/users/Einstein/Projects/Foo` | `Foo` is the shared folder |
| `ee1687e5-ac7f-426d-a6c0-03fed91d5f62:56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a:` | `56f7ceca-e7f8-4530-9a7a-fe4b7ec8089a` is the id of `Foo`, the path is empty |



{{< /hint >}}


{{< hint >}}
Note: the graph api uses `!` to seperate the `driveId` from a drive internal item `id` (`1383` which obviously is an opaque string no one should mess with):
```json
            "parentReference": {
                "driveId": "c12644a14b0a7750",
                "driveType": "personal",
                "id": "C12644A14B0A7750!1383",
                "name": "Screenshots",
                "path": "/drive/root:/Bilder/Screenshots"
            },
```

For the hot migration of oc10 instances a single *storage provider* with the owncloudsql *storage driver* can be used to handle all existing storages by prefixing the numeric `storage` id from the oc_storages table, e.g. `legacy-home` or `oc10-storages`. Using `$` as a seperator the reference *root* would be `<prefix>$<storage>` or `oc10-storages$234`.

File ids would become `<prefix>$<storage>!<fileid>`.

{{< /hint >}}
### Namespaces
A *namespace* is a set of paths with a common prefix. There is a global namespace when making requests against the CS3 *gateway* and a local namespace for every *storage provider*. Currently, oCIS mounts these *storage providers* out of the box:

| mountpoint | served namespace |
|------------|------------------|
| `/home`    | currently logged in users home |
| `/users`   | all users, used to access collaborative shares |
| `/public`  | publicly shared files, used to access public links |


{{< hint warning >}}
We plan to serve all shares that the current user has access to under `/shares`.
{{< /hint >}}

### Storage Drivers

A *storage driver* implements access to a *storage system*:

It maps the *path* and *id* based CS3 *references* to an appropriate *storage system* specific reference, e.g.:
- eos file ids
- posix inodes or paths
- deconstructed filesystem nodes

### Storage Provider

A *storage provider* manages *resources* identified by a *reference*
by accessing a *storage system* with a *storage driver*.

{{< svg src="extensions/storage/static/storageprovider.drawio.svg" >}}

{{< hint warning >}}
**Proposed Change**
A *storage provider* manages multiple *storage spaces*
by accessing a *storage system* with a *storage driver*.

{{< /hint >}}
{{< svg src="extensions/storage/static/storageprovider-spaces.drawio.svg" >}}
{{< hint warning >}}
By making *storage providers* aware of *storage spaces* we can get rid of the current `enablehome` flag / hack in reva. Furthermore, provisioning a new *storage space* becomes a generic operation, regardless of the need of provisioning a new user home or a new project space.
{{< /hint >}}

### Storage Registries

A *storage registry* manages the global *namespace*:
It is used by the *gateway*
to look up `address` and `port` of the *storage provider*
that should handle a *reference*.

{{< svg src="extensions/storage/static/storageregistry.drawio.svg" >}}

{{< hint warning >}}
**Proposed Change**
A *storage registry* manages the *namespace* for a *user*:
It is used by the *gateway*
to look up `address` and `port` of the *storage provider*
that is currently serving a *storage space*.
{{< /hint >}}
{{< svg src="extensions/storage/static/storageregistry-spaces.drawio.svg" >}}
{{< hint warning >}}
By making *storage registries* aware of *storage spaces* we can query them for a listing of all *storage spaces* a user has access to. Including his home, received shares, project folders or group drives. See [a WIP PR for spaces in the oCIS repo (#1827)](https://github.com/owncloud/ocis/pull/1827) for more info.
{{< /hint >}}

### Storage Spaces
A *storage space* is a logical concept:
It is a tree of *resources*
with a single *owner* (*user* or *group*), 
a *quota* and *permissions*, identified by a `storage space id`.

{{< svg src="extensions/storage/static/storagespace.drawio.svg" >}}

Examples would be every user's home storage space, project storage spaces or group storage spaces. While they all serve different purposes and may or may not have workflows like anti virus scanning enabled, we need a way to identify and manage these subtrees in a generic way. By creating a dedicated concept for them this becomes easier and literally makes the codebase cleaner. A *storage registry* then allows listing the properties of storage spaces, e.g. free space, quota, owner, syncable, root etag, uploed workflow steps, ...

Finally, a logical `storage space id` is not tied to a specific storage provider. If the *storage driver* supports it, we can import existing files including their *file id*, which makes it possible to move storage spaces between storage spaces to implement storage classes, e.g. with or without archival, workflows, on SSDs or HDDs.

