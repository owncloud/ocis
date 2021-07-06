## Problem Statement

When we speak about routes we have to make a difference between browser routes and internal API calls. Browser routes are interpreted by the web client (owncloud/web) to construct API calls<sup>1</sup>. With this in mind, this is the mapping on ownCloud Web:

|      | Browser URL                                                   | Internal Resolution                             |
|------|---------------------------------------------------------------|-------------------------------------------------|
| OCIS | `https://host/#/files/list/all/TEST`                           | `https://host/remote.php/webdav/TEST`           |
| OC10 | `https://host/index.php/apps/files/?dir=/TEST&fileid=5472225`   | `https://host/remote.php/dav/files/aunger/TEST`  |

Note that with an OC10 backend ownCloud's Web format remains unchanged: `https://host/index.html#/files/list/all/TEST` -- still resolves to --> `https://host/remote.php/webdav/TEST`. So here we have to make a distinction and limit the scope of this ADR to "how will a web client deal with browser urls?"<sup>2</sup>

Worth mentioning that on an OC10 backend it seems that `fileid` query parameter takes precedence over the `dir`. In fact if `dir` is invalid but `fileid` isn't, the resolution will succeed, as opposed to if the `fileid` is wrong and `dir` correct, resolution will fail altogether.

## Proposals

### Use private links as routes

First of, let's define what a private link is. A private link is:

> Another way to access a file or folder is via a private link. It’s a handy way of creating a permanent link for yourself or to point others to a file or folder, within a share, more efficiently. To access the private link, in the Sharing Panel for a file or folder, next to its name you’ll see a small link icon (1), as in the screenshot below.

_[source](https://doc.owncloud.com/server/user_manual/files/webgui/sharing.html#using-private-links)_

Private links are preceded by `/f/` to distinguish from `/s/` shares; this is convention.

### How MUST a bookmarked private link work with OCIS

The following flow chart provides an overview on how this should work.

![img](https://i.imgur.com/bE4xymv.png)

Regardless of the user being logged in or the private link being "public", ownCloud web receives the following URL:

`https://host/f/2748872` part of a more general format: `https://host/f/<resourceid>`

With an OC10 backend the resolution happens on OC10, and the response is a `3XX` with `Location` set to the actual URL, in this case:

`Location /index.php/apps/files/?dir=/TEST`

The proposed solution on `WEB-551` relies on the URL being of the format:

`https://xmpl.com/f/<spacealias>/<relative/path>?id=<b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607>`

or more blunt:

`https://host/f/space/relative/path?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`

Let us breakdown each section (except the obvious):

`/f/` = private link prefix. This is a convention.
`/space/` = space name. In which storage space does the target file / folder exist.
`/relative/path/` = path of the target file / folder relative to the storage space.
`?id=[...]` = combination of `storage_id` + `:` + `resource_id`

With the following information we can uniquely identify any resource within any known storages. This path is conditionally displayed.

`https://host/f/fileid` MUST therefore be expanded to the format: `https://host/f/space/relative/path?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`, as seen in the previous image.

This is achieved in ocis because:

- Admins will use the reva sql storage
- when OC10 + Web receives `https://host/index.php/f/5472225` -- which the internal resolution is -> `https://host/remote.php/dav/files/aunger/TEST` under an OC10 backend.
- when OCIS + Web receives `https://host/index.php/f/5472225` -- OCIS MUST

Now we have to do more hops. We have gaps to fill up here, so let us delegate the responsibility to OCIS. Parting from the assumption that there is a sql storage provider (this means we can query the resource by its old OC10 id) we could:

1. fetch the resource by its old ID = `/id=<>:<oc10_id>/`
2. find out in which storage space the reference exists = `/space/`
3. find out the storage provider that contains the reference = `/id=<storage_provider>:<>/`
4. find out the relative path to the root of the storage space of the reference = `/relative/path/`

Here we have essentially reconstructed all the info that we need that was defined in WEB-551, this can then be added to the resolution response for `GET https://host/index.php/f/5472225`

### How would an existing "general purpose" URL bookmark work with an OCIS backend? (WEB requirement)

A "general purpose bookmark" are just common paths we encounter by browsing on the web-ui.

`https://host/index.php/apps/files/?dir=/TEST&fileid=5472225`

Then at some point in the future the admin migrates to OCIS. What would happen to that bookmark? What would OCIS do if a request comes with such format? Well then OCIS has to transform the encoded information within that URL into something its API understand. As we can see here this is NOT a private link, but a simple URL I got just by opening the browser and navigating through my files.

As we mentioned we want OCIS to be backwards compatible with existing bookmarks, but we're now in a broken state. What must be done by OCIS in order to resolve this? As we saw in the "Browser URL - Internal Resolution" table this URL (received by the web client) needs to be adapted to the OCIS format, and we already have all the information we need. The result is:

`https://host/index.php/apps/files/?dir=/TEST&fileid=5472225` -> `https://host/remote.php/webdav/TEST`

## Sources

1. [Concepting: Use private link as route?](https://jira.owncloud.com/browse/WEB-551)
2. [Translate OC 10 paths to OCIS](https://jira.owncloud.com/browse/OCIS-1765) - fastlane ticket, blocked by (1)

## Assumptions

- <sup>1</sup> please provide input.
- <sup>2</sup> assumption. please provide input on this topic.
- <sup></sup>
