---
title: "Architecture"
date: 2018-05-02T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: architecture.md
---

## Frontend

The oCIS frontend service starts all services that handle incoming HTTP requests:
- *ocdav* for ownCloud flavoured WebDAV
- *ocs* for sharing, user provisioning, capabilities and other OCS API endpoints 
- *datagateway* for up and downloads
- TODO: *ocm*

{{< svg src="extensions/storage/static/frontend.drawio.svg" >}}

### WebDAV

The ocdav service not only handles all WebDAV requests under `(remote.php/)(web)dav` but also some other legacy endpoints like `status.php`:

| endpoint | service | CS3 api | CS3 namespace | description | TODO |
|----------|---------|-------------|------|------|------|
| *ownCloud 10 / current ocis setup:* |||||
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
| *Spaces concept also needs a new endpoint:* |||||
| `(remote.php/)dav/spaces/<spaceid>/<relative_path>` | ocdav | storageregistry & storageprovider | bypass path based namespace and directly talk to the responsible storage provider using a relative path | [spaces concept](https://github.com/owncloud/ocis/pull/1827) needs to point to storage spaces or a global endpoint | allow accessing spaces, listing is done by the graph api |


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


### Sharing

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

### User and Group provisioning

In oc10 users are identified by a username, which cannot change, because it is used as a foreign key in several tables. For oCIS we are internally identifying users by a UUID, while using the username in the WebDAV and OCS APIs for backwards compatability. To distinguish this in the URLs we are using `<username>` instead of `<userid>`. You may have encountered `<userlayout>`, which refers to a template that can be configuted to build several path segments by filling in user properties, e.g. the first two characters of the username or the issuer.

{{< hint warning >}}
Make no mistake, the [OCS Provisioning API](https://doc.owncloud.com/server/developer_manual/core/apis/provisioning-api.html) uses `userid` while it actually is the username, because it is what you use to login. 
{{< /hint >}}

We are currently working on adding [user management through the CS3 API](https://github.com/owncloud/ocis/pull/1930) to handle user and group provisioning (and deprovisioning).
