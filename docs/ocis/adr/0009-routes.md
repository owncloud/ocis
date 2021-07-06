## Problem Statement

When we speak about routes we have to make a difference between browser routes and internal API calls. Browser routes are interpreted by the web client (owncloud/web) to construct API calls<sup>1</sup>. With this in mind, this is the mapping on ownCloud Web:

|      | Browser URL                                                   | Internal Resolution                             |
|------|---------------------------------------------------------------|-------------------------------------------------|
| OCIS | `https://host/#/files/list/all/TEST`                           | `https://host/remote.php/webdav/TEST`           |
| OC10 | `https://host/index.php/apps/files/?dir=/TEST&fileid=5472225`   | `https://host/remote.php/dav/files/aunger/TEST`  |

Note that with an OC10 backend ownCloud's Web format remains unchanged: `https://host/index.html#/files/list/all/TEST` -- still resolves to --> `https://host/remote.php/webdav/TEST`. So here we have to make a distinction and limit the scope of this ADR to "how will a web client deal with browser urls?"<sup>2</sup>

Worth mentioning that on an OC10 backend it seems that `fileid` query parameter takes precedence over the `dir`. In fact if `dir` is invalid but `fileid` isn't, the resolution will succeed, as opposed to if the `fileid` is wrong and `dir` correct, resolution will fail altogether.

## Use private links as routes

First of, let's define what a private link is. A private link is:

> Another way to access a file or folder is via a private link. It’s a handy way of creating a permanent link for yourself or to point others to a file or folder, within a share, more efficiently. To access the private link, in the Sharing Panel for a file or folder, next to its name you’ll see a small link icon (1), as in the screenshot below.

_[source](https://doc.owncloud.com/server/user_manual/files/webgui/sharing.html#using-private-links)_

Private links are preceded by `/f/` to distinguish from `/s/` shares; this is convention.

## Private link path resolution

Let's have a look at the following scenario:

![img](https://i.imgur.com/hy0gSpB.jpeg)

_fig. 1_

We can observe that the private link can still remain the unchanged, this should provide functionality with existing bookmarks. In order for bookmarked private links to work, the migration from OC10 to OCIS should have taken effect, and the recommended SQL storage provider should be in use, this will allow the OCIS backend to resolve the ID's pre-migration.

Let us break down every step of figure 1.

1. `GET https://cloud.ocis.com/index.php/f/5472225`
    - all the web client has to "remember" is the id of the resource
2. `[303] Location=/marketing/path/to/file?id=storageid:resourceid`
    - the server will resolve the file by ID and provide with a URL for the webUI to render of the format: `/space/relative/path?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`
3. `PROPFIND https://host/remote.php/webdav/TEST`
    - To adjust to the new format, this WebDAV URL `MUST` change. If we don't have a namespace we can easily encounter naming collisions with different storage spaces<sup>1</sup>. A proposed WebDAV url format is recommended at this step of the following format: `https://host/remote.php/webdav/space/path/to/file?id=storageid:resourceid` which is provided by the server's original resolution.

When it comes to display the path, in order to avoid leaking parent information because the resource is shared, the rules in the following diagram `MUST` be followed:

 ![img](https://i.imgur.com/bE4xymv.png)

## On the server side

Receiving a GET request to the following resource `GET https://cloud.ocis.com/index.php/f/5472225` will trigger a few hops that `MUST` be cached in order to prevent slow response times. The nature of these requests can be cached because the resources ID are not subject to changes.

The server `MUST` have a way to resolve the `ID=5472225`. The easiest approach that comes to mind is using the SQL storage driver, that provides compatibility when it comes to migrating files from an OC10 to an OCIS backend. The queried ID already exists in the DB, and the storage driver will just pull all the info it needs to construct the URL to set in the `Location` header of the response.

Let us breakdown each section (except the obvious):

`/f/` = private link prefix. This is a convention.
`/space/` = space name. In which storage space does the target file / folder exist.
`/relative/path/` = path of the target file / folder relative to the storage space.
`?id=[...]` = combination of `storage_id` + `:` + `resource_id`

With all the above data we can start building the Location response header.

## Footnotes

- <sup>1</sup> is this a real concern? Need read proof.
