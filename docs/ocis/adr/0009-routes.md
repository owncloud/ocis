---
title: "9. Global URL"
date: 2021-07-07T14:55:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0009-routes.md
---

* Status: proposed
* Deciders: @refs, @butonic, @micbar, @dragotin, @pmaier1
* Date: 2021-07-07

## Context and Problem Statement

When we speak about routes we have to make a difference between browser routes and internal API calls. Browser routes are interpreted by the web client (owncloud/web) to construct API calls. With this in mind, this is the mapping on ownCloud Web with OC10 and OCIS backend:

|      | Browser URL                                                                                                                                                                 | Internal Resolution                                                                                                                         |
|------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|
| OC10 | `https://host/index.php/apps/files/?dir=/TEST&fileid=5472225`                                                                                                                 | `https://host/remote.php/dav/files/aunger/TEST`                                                                                              |
| OCIS | `https://host/#/files/list/all/TEST`                                                                                                                                         | `https://host/remote.php/webdav/TEST`                                                                                                       |
| OCIS (after this ADR is implemented) | `https://host/#/s/<space_id>/path/to/file?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`                      | `https://host/remote.php/webdav/space/relative/path?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`           |

Note that with an OC10 backend ownCloud's Web format remains unchanged: `https://host/index.html#/files/list/all/TEST` -- still resolves to --> `https://host/remote.php/webdav/TEST`. So here we have to make a distinction and limit the scope of this ADR to "how will a web client deal with the browser url?"

Worth mentioning that on an OC10 backend it seems that `fileid` query parameter takes precedence over the `dir`. In fact if `dir` is invalid but `fileid` isn't, the resolution will succeed, as opposed to if the `fileid` is wrong (doesn't exist) and `dir` correct, resolution will fail altogether.

<spaceid> is composed of `<storage_id>:<node_id>`

## Decision Drivers

* Construct a URL that is not only readable by the user, but it contains all the necessary information to build a query to the backend.

## Considered Options

* Consistent Global URL Format

## Decision Outcome

Chosen option: "Consistent Global URL Format".

### Positive Consequences

* Backwards compatibility with existing bookmarks
* Complete visibility of the tree in the URL
* Unify user facing URL

## Proposed Global URL Format

`https://<host>/#/s/<spaceid>/<relative/path>?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`

`/s` denotes that this is a space url.

## URL Semantics

- The relative path and ID are optional. This URL is valid and points to the root of the space with ID = `b78c2044-5b51-446f-82f6-907a664d089`: `https://example.com/#/s/b78c2044-5b51-446f-82f6-907a664d089`
- The following case is valid and will resolve the correct folder ONLY IF the path exists within the space `https://example.com/#/s/b78c2044-5b51-446f-82f6-907a664d089/path/to/file`
- To improve in the previous example and ensure a more resilient link, adding the query string `id` of the target resource (folder or file) is encouraged to prevent always resolving even if the resource is renamed: `https://example.com/#/s/b78c2044-5b51-446f-82f6-907a664d089/path/to/file?id=ba4c1820-df12-11eb-8dcd-ff21f12c1264:beb78dd6-df12-11eb-a05c-a395505126f6`

With the above explained, let's see some use cases:

### Example 1: UserA shares something from her Home folder with UserB

- open the browser and go to `ocis.com`
- the browser's url changes to: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`. You're now in YOUR home folder / personal space.
- you create a new folder `TEST` and navigate into it
  - the URL now changes to: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607/TEST`
- You share `TEST` with some else
- YOU navigate into `TEST`
  - now the URL would look like: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:3a9305da-df17-11eb-ab99-abe09d93e08a`

As you can see, even if you're the owner of `TEST` and navigate into it, the URL changed due to a new space was created. This ensures that while working in your home folder, copying URL and giving them to the person you share the resource with the receiver can still navigate within the new space.

In short terms, while navigating using the WebUI, the URL has to constantly change whenever we change spaces.

### Example 2: UserA shares something from a Workspace

Assuming we only have one storage provider; a consequence of this, all storage spaces will start with the same storage_id.

- open the browser and go to `ocis.com`
- the browser's url changes to: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`. You're now in YOUR home folder / personal space.
- you have access to a workspace called `foo` (created by an admin)
- navigate into workspace `foo`
  - the URL now changes to: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:d342f9ce-df18-11eb-b319-1b6d9df4bc74`. You are now at the root of the workspace `foo`.
    - because we only have one storage provider, the `space_id` section of the URL only updates the `node_id` part of it.
    - had we had more than one storage provider, the `space_id` would depend on which storage provider contains the storage space.
- you create a folder `TEST`
- you navigate into `TEST`
  - now the URL would look like: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:d342f9ce-df18-11eb-b319-1b6d9df4bc74/TEST`
  - or a more robust url: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:d342f9ce-df18-11eb-b319-1b6d9df4bc74/TEST?id=b78c2044-5b51-446f-82f6-907a664d089c:04f1991c-df19-11eb-9cc7-3b09f04f9ca3`

## Rules for reference resolution

- if path can be resolved, then the path is used as the "locator" (instead of the id)
- if path cannot be resolved (eg. because of misspelling, folder moved, renamed etc.) then the id is used as a "locator"

## Manipulating the path in the URL

Suppose a power user knows where their resources are and wants to navigate only by modifying the request in the webUI. The user goes to the browser and changes:

`https://host/space/relative/path?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`

to

`https://host/space/relative/path/deeper/file/inside?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`

Notice that the ID has not changed. Path based resolution should take precedence over ID resolution. Now we have a `GET` request that the webUI has to adapt to the server's format:

## Considerations

Navigating into a folder that is the root of a space changes the url to reflect that we are now in the root of a space.

## Future Improvements

### Spaces Registry

A big drawback against this idea is that the length of the URL is increased by a lot, rendering them almost unreadable. Introducing a Spaces Registry (SR) would shorten them. Let's see how.

A URL without a SR would look like: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:d342f9ce-df18-11eb-b319-1b6d9df4bc74/TEST?id=b78c2044-5b51-446f-82f6-907a664d089c:04f1991c-df19-11eb-9cc7-3b09f04f9ca3`
The same URL with a SR `https://ocis.com/#/s/workspaceFoo/TEST?id=b78c2044-5b51-446f-82f6-907a664d089c:04f1991c-df19-11eb-9cc7-3b09f04f9ca3`

Space Registry resolution can happen at the client side (i.e: the client keeps a list of space name -> space id [where space id = storageid + nodeid]; the client queries a SR) or server side. Server side is more resilient due to clients can have limited networking; for instance if they are running on a tight intranet.
