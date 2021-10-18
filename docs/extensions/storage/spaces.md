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

Udating an entity attribute:

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
