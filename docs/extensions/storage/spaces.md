---
title: "Spaces"
date: 2020-04-27T18:46:00+01:00
weight: 38
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: spaces.md
---

{{< toc >}}

## Editing a Storage Space

The OData specification allows for a mirage of ways of addressing an entity. We will support addressing a Drive entity by its unique identifier, which is the one the graph-api returns when listing spaces, and its format is:

```json
{
  "id": "1284d238-aa92-42ce-bdc4-0b0000009157!b6e2c9cc-9dbe-42f0-b522-4f2d3e175e9c"
}
```

This is an extract of an element of the list spaces response. An entire object has the following shape:

```json
{
    "driveType": "project",
    "id": "1284d238-aa92-42ce-bdc4-0b0000009157!b6e2c9cc-9dbe-42f0-b522-4f2d3e175e9c",
    "lastModifiedDateTime": "2021-10-07T11:06:43.245418+02:00",
    "name": "marketing",
    "owner": {
        "user": {
            "id": "ddc2004c-0977-11eb-9d3f-a793888cd0f8"
        }
    },
    "quota": {
        "total": 65536
    },
    "root": {
        "id": "1284d238-aa92-42ce-bdc4-0b0000009157!b6e2c9cc-9dbe-42f0-b522-4f2d3e175e9c",
        "webDavUrl": "https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157!b6e2c9cc-9dbe-42f0-b522-4f2d3e175e9c"
    }
}
```

### Updating a space property

Having introduced the above, one can refer to a Drive with the following URL format:

```console
'https://localhost:9200/graph/v1.0/drives/1284d238-aa92-42ce-bdc4-0b0000009157!07c26b3a-9944-4f2b-ab33-b0b326fc7570
```

Updating an entity attribute:

```console
curl -X PATCH 'https://localhost:9200/graph/v1.0/drives/1284d238-aa92-42ce-bdc4-0b0000009157!07c26b3a-9944-4f2b-ab33-b0b326fc7570' -d '{"name":"42"}' -v
```

The previous URL resource path segment (`1284d238-aa92-42ce-bdc4-0b0000009157!07c26b3a-9944-4f2b-ab33-b0b326fc7570`) is parsed and handed over to the storage registry in order to apply the patch changes in the body, in this case update the space name attribute to `42`. Since space names are not unique we only support addressing them by their unique identifiers, any other query would render too ambiguous and explode in complexity.


### Updating a space description

Since every space is the root of a webdav directory, following some conventions we can make use of this to set a default storage description and image. In order to do so, every space is created with a hidden `.space` folder at its root, which can be used to store such data.

```curl
curl -k -X PUT https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157\!07c26b3a-9944-4f2b-ab33-b0b326fc7570/.space/description.md -d "Add a description to your spaces" -u admin:admin
```

Verify the description was updated:

```curl
❯ curl -k https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157\!07c26b3a-9944-4f2b-ab33-b0b326fc7570/.space/description.md -u admin:admin
Add a description to your spaces
```

This feature makes use of the internal storage layout and is completely abstracted from the end user.

### Quotas

Spaces capacity (quota) is independent of the Storage quota. As a Space admin you can set the quota for all users of a space, and as such, there are no limitations and is up to the admin to make a correct use of this.

It is possible to have a space quota greater than the storage quota. A Space may also have "infinite" quota, meaning a single space without quota can occupy the entirety of a disk.

#### Quota Enforcement

Creating a Space with a quota of 10 bytes:

`curl -k -X POST 'https://localhost:9200/graph/v1.0/drives' -u admin:admin -d '{"name":"marketing", "quota": {"total": 10}}' -v`

```console
/var/tmp/ocis/storage/users
├── blobs
├── nodes
│   ├── 627981c2-2a71-4adf-b680-177e245afdda
│   ├── 9541e7c3-8fda-4b49-b697-e7e51457cf5a
│   ├── b5692345-108d-4b80-9747-3a7e9739ad57
│   └── root
│       ├── 118351d7-67a4-4cdf-b495-6093d1e572ed -> ../627981c2-2a71-4adf-b680-177e245afdda
│       └── ddc2004c-0977-11eb-9d3f-a793888cd0f8 -> ../b5692345-108d-4b80-9747-3a7e9739ad57
├── spaces
│   ├── personal
│   │   └── b5692345-108d-4b80-9747-3a7e9739ad57 -> ../../nodes/b5692345-108d-4b80-9747-3a7e9739ad57
│   ├── project
│   │   └── 627981c2-2a71-4adf-b680-177e245afdda -> ../../nodes/627981c2-2a71-4adf-b680-177e245afdda
│   └── share
├── trash
└── uploads
```

Verify the new space has 10 bytes, and none of it is used:

```json
{
    "driveType": "project",
    "id": "1284d238-aa92-42ce-bdc4-0b0000009157!627981c2-2a71-4adf-b680-177e245afdda",
    "lastModifiedDateTime": "2021-10-15T11:16:26.029188+02:00",
    "name": "marketing",
    "owner": {
        "user": {
            "id": "ddc2004c-0977-11eb-9d3f-a793888cd0f8"
        }
    },
    "quota": {
        "remaining": 10,
        "total": 10,
        "used": 0
    },
    "root": {
        "id": "1284d238-aa92-42ce-bdc4-0b0000009157!627981c2-2a71-4adf-b680-177e245afdda",
        "webDavUrl": "https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157!627981c2-2a71-4adf-b680-177e245afdda"
    }
}
```

Upload a 6 bytes file:

`curl -k -X PUT https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157\!627981c2-2a71-4adf-b680-177e245afdda/6bytes.txt -d "012345" -u admin:admin -v`

Query the quota again:

```json
{
  "quota": {
    "remaining": 4,
    "total": 10,
    "used": 6
  }
}
```

Now attempt to upload 5 bytes to the space:

`curl -k -X PUT https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157\!627981c2-2a71-4adf-b680-177e245afdda/5bytes.txt -d "01234" -u admin:admin -v`

The request will fail with `507 Insufficient Storage`:

```
 HTTP/1.1 507 Insufficient Storage
< Access-Control-Allow-Origin: *
< Content-Length: 0
< Content-Security-Policy: default-src 'none';
< Date: Fri, 15 Oct 2021 09:24:46 GMT
< Vary: Origin
< X-Content-Type-Options: nosniff
< X-Download-Options: noopen
< X-Frame-Options: SAMEORIGIN
< X-Permitted-Cross-Domain-Policies: none
< X-Robots-Tag: none
< X-Xss-Protection: 1; mode=block
<
* Connection #0 to host localhost left intact
* Closing connection 0
```

##### Considerations

- If a Space quota is updated to unlimited, the upper limit is the entire available space on disk
{{< hint warning >}}

The current implementation in oCIS might not yet fully reflect this concept. Feel free to add links to ADRs, PRs and Issues in short warning boxes like this.

{{< /hint >}}

## Storage Spaces
A storage *space* is a logical concept. It organizes a set of [*resources*]({{< ref "#resources" >}}) in a hierarchical tree. It has a single *owner* (*user* or *group*),
a *quota*, *permissions* and is identified by a `storage space id`.

{{< svg src="extensions/storage/static/storagespace.drawio.svg" >}}

Examples would be every user's personal storage *space*, project storage *spaces* or group storage *spaces*. While they all serve different purposes and may or may not have workflows like anti virus scanning enabled, we need a way to identify and manage these subtrees in a generic way. By creating a dedicated concept for them this becomes easier and literally makes the codebase cleaner. A storage [*Spaces Registry*]({{< ref "./spacesregistry.md" >}}) then allows listing the capabilities of storage *spaces*, e.g. free space, quota, owner, syncable, root etag, upload workflow steps, ...

Finally, a logical `storage space id` is not tied to a specific [*spaces provider*]({{< ref "./spacesprovider.md" >}}). If the [*storage driver*]({{< ref "./storagedrivers.md" >}}) supports it, we can import existing files including their `file id`, which makes it possible to move storage *spaces* between [*spaces providers*]({{< ref "./spacesprovider.md" >}}) to implement storage classes, e.g. with or without archival, workflows, on SSDs or HDDs.

## Shares
*To be clarified: we are aware that [*storage spaces*]({{< ref "#storage-spaces" >}}) may be too 'heavywheight' for ad hoc sharing with groups. That being said, there is no technical reason why group shares should not be treated like storage [*spaces*]({{< ref "#storage-spaces" >}}) that users can provision themselves. They would share the quota with the users home or personal storage [*space*]({{< ref "#storage-spaces" >}}) and the share initiator would be the sole owner. Technically, the mechanism of treating a share like a new storage [*space*]({{< ref "#storage-spaces" >}}) would be the same. This obviously also extends to user shares and even file individual shares that would be wrapped in a virtual collection. It would also become possible to share collections of arbitrary files in a single storage space, e.g. the ten best pictures from a large album.*

## Notes

We can implement [ListStorageSpaces](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ListStorageSpacesRequest) by either
- iterating over the root of the storage and treating every folder following the `<user_layout>` as a `home` *storage space*,
- iterating over the root of the storage and treating every folder following a new `<project_layout>` as a `project` *storage space*, or
- iterating over the root of the storage and treating every folder following a generic `<layout>` as a *storage space* for a configurable space type, or
- we allow configuring a map of `space type` to `layout` (based on the [CreateStorageSpaceRequest](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.CreateStorageSpaceRequest)) which would allow things like
```
home=/var/lib/ocis/storage/home/{{substr 0 1 .Owner.Username}}/{{.Owner.Username}}
spaces=/spaces/var/lib/ocis/storage/projects/{{.Name}}
```
