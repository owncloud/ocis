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
'https://localhost:9200/graph/v1.0/Drive(1284d238-aa92-42ce-bdc4-0b0000009157!07c26b3a-9944-4f2b-ab33-b0b326fc7570")
```

Updating an entity attribute:

```console
curl -X PATCH 'https://localhost:9200/graph/v1.0/Drive("1284d238-aa92-42ce-bdc4-0b0000009157!07c26b3a-9944-4f2b-ab33-b0b326fc7570)' -d '{"name":"42"}' -v
```

The previous URL resource path segment (`Drive(1284d238-aa92-42ce-bdc4-0b0000009157!07c26b3a-9944-4f2b-ab33-b0b326fc7570)`) is parsed and handed over to the storage registry in order to apply the patch changes in the body, in this case update the space name attribute to `42`. Since space names are not unique we only support addressing them by their unique identifiers, any other query would render too ambiguous and explode in complexity.


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

Spaces capacity (quota) is independent of where the Storage quota. As a Space admin you can set the quota for all users of a space, and as such, there are no limitations and is up to the admin to make a correct use of this.

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
...
"quota": {
    "remaining": 4,
    "total": 10,
    "used": 6
},
...
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
