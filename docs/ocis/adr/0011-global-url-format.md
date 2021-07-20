---
title: "11. WebUI URL format"
weight: 11
date: 2021-07-07T14:55:00+01:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0011-global-url-format.md
---

* Status: proposed
* Deciders: @refs, @butonic, @micbar, @dragotin, @hodyroff, @pmaier1, @fschade, @tbsbdr, @kulmann
* Date: 2021-07-07

## Context and Problem Statement

When speaking about URLs we have to make a difference between browser URLs and API URLs. Only browser URLs are visible to end users and will be bookmarked. The currently existing and bookmarked ownCloud 10 URLs look something like this:

```
GET https://demo.owncloud.com/apps/files/?dir=/path/to/resource&fileid=5472225
303 Location: https://demo.owncloud.com/apps/files/?dir=/path/to/resource
```

When the URL contains a `fileid` parameter the server will look up the corresponding `dir`, overwriting whatever was set before the redirect. The `fileid` always takes precedence and the server is responsible for the lookup.

```
GET https://demo.owncloud.com/apps/files/?dir=/path/to/resource
```

The `dir` parameter is then used to make a WebDAV request against the `/dav/files` endpoint of the currently logged in user:

```
PROPFIND https://demo.owncloud.com/remote.php/dav/files/demo/path/to/resource
```

The resulting PROPFIND response is used to render the file listing. All good so far.

For the new ocis web UI we want to clean up the user visible Browser URLs. They currently look like this:

```
https://demo.owncloud.com/#/files/list/all/path/to/resource
```

Currently, there is no `fileid` like parameter in the browser URL, making bookmarks of it fragile (they break when a bookmarked folder is renamed).

The oCIS web UI just takes the path and uses the `/webdav` endpoint of the currently logged in user:

```
PROPFIND https://demo.owncloud.com/remote.php/webdav/path/to/resource
```


With the new ownCloud web client (owncloud/web)

 needs to interpret them to make API calls. With this in mind, this is the current mapping on ownCloud Web with OC10 and OCIS backend:

|      | Browser URL                                                    | API URL                                |
|------|----------------------------------------------------------------|----------------------------------------------------|
| OC10 + classic WebUI | `https://demo.owncloud.com/apps/files/?dir=/path/to/resource&fileid=5472225`  | `https://demo.owncloud.com/remote.php/dav/files/demo/path/to/resource`    |
| OC10 + OCIS WebUI| `https://web.owncloud.com/index.html#/files/list/all/path%2Fto%2Fresource`  | `https://demo.owncloud.com/remote.php/webdav/path/to/resource`    |
| OCIS | `https://demo.owncloud.com/#/files/list/all/path/to/resource`                           | `https://demo.owncloud.com/remote.php/webdav/path/to/resource`              |


On an OC10 backend the `fileid` query parameter takes precedence over the `dir`. In fact if `dir` is invalid but `fileid` isn't, the resolution will succeed, as opposed to if the `fileid` is wrong (doesn't exist) and `dir` correct, resolution will fail altogether with a 404.

This ADR is limited to the scope of "how will a web client deal with the browser URL?". The API URLs will change with the spaces concept to `https://demo.owncloud.com/dav/spaces/<space_id>/relative/path/to/resource`. The Web UI can look up a space id and the mount path using the `/graph/v1.0/drives` API:
1. TODO for a given resource id as part of the URL the `https://demo.owncloud.com/v1.0/drive/items/123456A14B0A7750!359?$select=parentReference` can be used to retrieve the drive/space:
```
{
    "parentReference": {
        "driveId": "123456a14b0a7750",
        "driveType": "personal",
        "id": "123456A14B0A7750!357",
        "path": "/drive/root:"
    }
}
```
2. TODO to fetch the list of all spaces with their mount points we need an API endpoint that allows clients (not only the web ui) to 'sync' the list of storages a user has access to from the storage registry on the server side. This allows clients to directly talk to a storage provider on another instance, allowing true storage federation. The MS graph api has no notion of mount points, so we will need to add a `mountpath` *(or `mountpoint`? or `alias`?)* to our [`drive` resource properties in the libreGraph spec](https://github.com/owncloud/open-graph-api/blob/dc6da5359eee0345429080b5b59762fd8c57b121/api/openapi-spec/v0.0.yaml#L351-L384). Tracked in https://github.com/owncloud/open-graph-api/issues/6


{{< hint >}}
@jfd: The graph api returns a `path` in the `parentReference`, which is part of the `root` in a `drive` resource. But it contains a value in the namespace of the `graph` endpoint, eg.: `/drive/root:/Bilder` for the `/Bilder` folder in the root of the currently logged in users personal drive/space. Which is again relative to the drive. To give the clients a way to determine the mount point we need to add a new `mountpath/point/alias` property.
{{< /hint >}}

## Decision Drivers

* To reveal relevant context to the user URLs should either carry a path component or a meaningful alias
* To prevent bookmarks from breaking URLs should have an id component that can be used by the system to lookup the resource

## Considered Options

* Existing ownCloud 10 URLs
* ID based URLs
* Path based URLs
* Space based URLs
* Mixed Global URLs
* Configurable path component in URLs

## Decision Outcome

Chosen option: "[option 1]", because [justification. e.g., only option, which meets k.o. criterion decision driver | which resolves force force | … | comes out best (see below)].

### Positive Consequences <!-- optional -->

* [e.g., improvement of quality attribute satisfaction, follow-up decisions required, …]
* …

### Negative Consequences <!-- optional -->

* [e.g., compromising quality attribute, follow-up decisions required, …]
* …

## Pros and Cons of the Options

### Existing OwnCloud 10 URLs

The existing ownCloud 10 URLs look like this

| URL | comment |
|-|-|
| `https://<host>/apps/files/?dir=<path>&fileid=<fileid>` | pattern |
| `https://demo.owncloud.com/apps/files/?dir=/&fileid=18` | root of the currently logged in user |
| `https://demo.owncloud.com/index.php/apps/files/?dir=/path/to/resource&fileid=192` | sub folder `/path/to/resource` |

It contains a path and a `fileid` (which takes precedence). 

* Good, because the `fileid` prevents bookmarks from breaking
* Good, because the `dir` reveals context in the form of a path
* Bad, because the web UI needs to look up the space alias in a registry to build an API request for the `/dav/space` endpoint
* Bad, because URLs still contain a long prefix `(/index.php)/apps/files`
* Bad, because the `fileid` needs to be accompanied by a `storageid` to allow efficient routing in ocis
* Bad, because if not configured properly an additional `/index.php` prefixes the route
* Bad, because powerusers cannot navigate by updating only the path in the URL, as the `fileid` takes precedence. They have to delete the `fileid` to navigate

### ID based URLs

MS OneDrive has URLs like this:

| URL | comment |
|-|-|
| `https://<host>/?id=<fileid>(&cid=<cid>)` | pattern, the `cid` is optional but added automatically |
| `https://onedrive.live.com/?id=root&cid=A12345A14B0A7750` | root of a personal drive |
| `https://onedrive.live.com/?id=A12345A14B0A7750%21359&cid=A12345A14B0A7750` | sub folder in a personal drive |

It contains only IDs but no folder names. The `fileid` is a URL encoded `<cid>!<numericid>`. Very similar to the CS3 `resourceid` which consists of `storageid` and `nodeid`.

* Good, because bookmarks cannot break
* Good, because URLs do not disclose unshared path segments
* Bad, because the web UI needs to look up the space id in a registry to build an API request for the `/dav/space` endpoint
* Bad, because URLs reveal no context to users

### Path based URLs

There is a customized ownCloud instance that uses path only based URLs:

| URL | comment |
|-|-|
| `https://<host>/apps/files/?dir=/&` | root of the currently logged in user |
| `https://demo.owncloud.com/apps/files/?dir=/&` | root of the currently logged in user |
| `https://demo.owncloud.com/apps/files/?dir=/path/to/resource&` | sub folder `/path/to/resource` |

* Good, because the URLs reveal the full path context to users
* Good, because powerusers can navigate by updating the path in the url
* Bad, because the web UI needs to look up the space id in a registry to build an API request for the `/dav/space` endpoint
* Bad, because the bookmarks break when someone renames a folder in the path
* Bad, because there is no id that can be used as a fallback lookup mechanism
* Bad, because URLs might leak too much context (parent folders of shared files)

### Space based URLs

| URL | comment |
|-|-|
| `https://<host>/#/s/<space_id>(/<relative/path>)(?id=<resource_id>)` | the pattern, relative `path` and `resource_id` are optional |
| `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607` | root of a storage space, might be the currently logged in users home |
| `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607/relative/path/to/resource` | sub folder `/relative/path/to/resource` in the storage with id `b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`, works ***only*** if path still exists  |
| `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607/relative/path/to/resource?id=ba4c1820-df12-11eb-8dcd-ff21f12c1264:beb78dd6-df12-11eb-a05c-a395505126f6` | sub folder `/relative/path/to/resource` in the storage with id `b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`, lookup can fall back to the `id` |

{{< hint >}}
* `/#` is used by the current vue router.
* `/s` denotes that this is a space url.
* `<space_id>` and `<resource_id>` both consist of `<storage_id>:<node_id>`, but the `space_id` can be replaced with a shorter id or an alias. See furthor down below.
* `<relative/path>` takes precedence over the `<resource_id>`, both are optional
{{< /hint >}}

* Good, because the web UI does not need to look up the space id in a registry to build an API request for the `/dav/space` endpoint
* Good, because the URLs reveal a relevant path context to users
* Good, because everything after the `#` is not sent to the server, building the webdav request to list the folder is offloaded to the clients
* Good, because powerusers can navigate by updating the path in the url
* Bad, because the current ids are uuid based, leading to very long URLs where the path component nearly vanishes between two very long strings
* Bad, because the `#` in the URL is just a technical requirement
* Bad, because ocis web requires a `/#/files/s` at the root of the route to distinguish the files app from other apps
* Bad, while navigating using the WebUI, the URL has to be updated whenever we change spaces.
* Bad, because the technical `<space_id>` is meaningless to end users

With the above explained, let's see some use cases:

#### Example 1: UserA shares something from her Home folder with UserB

- open the browser and go to `demo.owncloud.com`
- the browser's url changes to: `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`. You're now in YOUR home folder / personal space.
- you create a new folder `/relative/path/to/resource` and navigate into `/relative/path/to`
  - the URL now changes to: `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607/relative/path/to`
- You share `resource` with some else
- You navigate into `/relative/path/to/resource`
  - now the URL would look like: `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:3a9305da-df17-11eb-ab99-abe09d93e08a`

As you can see, even if you're the owner of `/relative/path/to/resource` and navigate into it, the URL changes due to a new space being entered. This ensures that while working in your home folder, copying URLs and giving them to the person you share the resource with, the receiver can still navigate within the new space.

In short terms, while navigating using the WebUI, the URL has to constantly change whenever we change spaces to reflect the most explicit one.

#### Example 2: UserA shares something from a Workspace

Assuming we only have one storage provider; a consequence of this, all storage spaces will start with the same storage_id.

- open the browser and go to `demo.owncloud.com`
- the browser's url changes to: `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607`. You're now in YOUR home folder / personal space.
- you have access to a workspace called `foo` (created by an admin)
- navigate into workspace `foo`
  - the URL now changes to: `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:d342f9ce-df18-11eb-b319-1b6d9df4bc74`. You are now at the root of the workspace `foo`.
    - because we only have one storage provider, the `space_id` section of the URL only updates the `node_id` part of it.
    - had we had more than one storage provider, the `space_id` would depend on which storage provider contains the storage space.
- you create a folder `/relative/path/to/resource`
- you navigate into `/relative/path/to/resource`
  - now the URL would look like: `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:d342f9ce-df18-11eb-b319-1b6d9df4bc74/relative/path/to/resource`
  - or a more robust url: `https://demo.owncloud.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:d342f9ce-df18-11eb-b319-1b6d9df4bc74/relative/path/to/resource?id=b78c2044-5b51-446f-82f6-907a664d089c:04f1991c-df19-11eb-9cc7-3b09f04f9ca3`

#### Spaces Registry

A big drawback against this idea is that the length of the URL is increased by a lot, rendering them almost unreadable. Introducing a Spaces Registry (SR) would shorten them. Let's see how.

A URL without a SR would look like: `https://ocis.com/#/s/b78c2044-5b51-446f-82f6-907a664d089c:d342f9ce-df18-11eb-b319-1b6d9df4bc74/TEST?id=b78c2044-5b51-446f-82f6-907a664d089c:04f1991c-df19-11eb-9cc7-3b09f04f9ca3`
The same URL with a SR `https://ocis.com/#/s/workspaceFoo/TEST?id=b78c2044-5b51-446f-82f6-907a664d089c:04f1991c-df19-11eb-9cc7-3b09f04f9ca3`

Space Registry resolution can happen at the client side (i.e: the client keeps a list of space name -> space id [where space id = storageid + nodeid]; the client queries a SR) or server side. Server side is more resilient due to clients can have limited networking; for instance if they are running on a tight intranet.

### Mixed Global URLs

While ID based space URLs can be made more readable by shortening the IDs they only start to reveal context when an alias is used instead of the space id. These aliases however have to be unique identifiers. These aliases should live in namespaces like `/workspaces/marketing` and `/personal/marketing` to make phishing attacks harder (in this case a user that registered with the username `marketing`). But namespaced aliases is semantically equivalent to ... a path hierarchy.

When every space has a namespaced alias and a relative path we can build a global namespace:

| URL | comment |
|-|-|
| `https://<host>/files</namespaced/alias></relative/path/to/resource>?id=<resource_id>` | the pattern, `/files` might become optional |
| `https://demo.owncloud.com/files/personal/einstein/?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21607` | root of user `einstein` |
| `https://demo.owncloud.com/files/personal/einstein/relative/path/to/resource?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21608` | sub folder `/relative/path/to/resource` |
| `https://demo.owncloud.com/files/shares/einstein/somesharename?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21608` | shared URL for `/relative/path/to/resource` |
| `https://demo.owncloud.com/files/public/kcZVYaXr7oZ66bg/relative/path/to/resource` | sub folder `/relative/path/to/resource` in public link with token `kcZVYaXr7oZ66bg` |
| `https://demo.owncloud.com/files/public/kcZVYaXr7oZ66bg/relative/path/to/resource` | sub folder `/relative/path/to/resource` in public link with token `kcZVYaXr7oZ66bg` |
| `https://demo.owncloud.com/files/personal/einstein/marie is stupid/and richard as well/resource?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21608` | sub folder `marie is stupid/and richard as well/resource` ... something einstein might not want to reveal |
| `https://demo.owncloud.com/files/shares/einstein/resource (2)?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21608` | named link URL for `/marie is stupid/and richard as well/resource`, does not disclose the actual hierarchy, has an appended counter to avaid a collision |
| `https://demo.owncloud.com/files/shares/einstein/mybestfriends?id=b78c2044-5b51-446f-82f6-907a664d089c:194b4a97-597c-4461-ab56-afd4f5a21608` | named link URL for `/marie is stupid/and richard as well/resource`, does not disclose the actual hierarchy, has a custom alias for the share |

`</namespaced/alias></relative/path/to/resource>` is the global path in the CS3 api. The CS3 Storage Registry is responsible by managing the mount points.

In order to be able to copy and paste URLs all resources must be uniquely identifiable:

* Instead of `/home` the URL always has to reflect the user: `/personal/einstein`
* Public links can use `/public/<token>`
* workspaces can use `/workspaces/<alias>` or `/workspaces/<additional>/<classification>/<alias>` where the hierarchy is given by the organization
* experiments can use `/experiments/<alias>`
* research institutes could set up `/papers/<researchgroup>/<alias>`
* trash could be accessed by prefixing the namespace alias with `/trash`? or using `/trash/<space_id>`
* instead of a namespaced alias a storage space id could be used with a generic `/space/<space_id>` namespace

The alias namespace hierarchy and depth can be pre determined by the admin. Even if aliases change the `id` parameter prevents bookmarks from breaking. A user can decide to build a different hierarchy by using his own registry. 

What about shares? Similar to `/home` it must reflect the user: `/shares/einstein` would list all shares *by* einstein for the currently logged in user. The ui needs to apply the same URL rewriting as for space based URLs: when navigating into a share the URL has to switch from `/personal/einstein/relative/path/to/shared/resource` to `/shares/einstein/<unique and potentially namespaced alias for shared resource>`. When more than one `resource` was shared a name collision would occur. To prevent this we can use ids `/shares/einstein/id/<resource_id` or namespaced aliases `/shares/einstein/files/alias`. Similar to the `/trash` prefix we could treat `/shares` as a filter for the shared resources a user has access to, but that would disclose unshared path segments in personal spaces. We could make that a feature and let users create an alias for a shared resource, similar as for public links. Then they can decide if they want to disclose the full path in their personal space (or another workspace) or if they want to use an alias which is then accessed at `/shares/einstein/<alias>`. As a default we could take the alias at creation time from the filename. That way two shares to a resource with the same name, eg.: `/personal/einstein/project AAA/foo` and `/personal/einstein/project BBB/foo` would lead to `/shares/einstein/foo` (a CS3 internal reference to `/personal/einstein/project AAA/foo`) and `/shares/einstein/foo (2)` (a CS3 internal reference to `/personal/einstein/project BBB/foo`). `foo (2)` would keep its name even when `foo` is deleted or renamed. Well an id as the alias might be better then, because users might rename these aliases, which would break URLs if they have been bookmarked. In any case this would make end user more aware of what they share AND it would allow them to choose an arbitrary context for the links they want to send out: personal internal share URLs. 

With these different namespaces the `/files` part in the URL becomes obsolete, because the files application can be registered for multiple namespaces: `/personal`, `/workspaces`, `/shares`, `/trash` ...

* Good, because it contains a global path
* Good, because spaces with namespaced aliases can by bookmarked and copied into mails or chat without disclosing unshared path segments, as the space is supposed to be shared
* Good, because the UI can detect broken paths and notify the user to update his bookmark if the resource could be found by `id`
* Good, because the `/files` part might only be required for `id` only based lookup to let the web ui know which app is responsible for the route
* Good, because it turns shares into deliberately named spaces in `/shares/<owner>/<alias>`
* Bad, because the web UI needs to look up the space alias in a registry to build an API request for the `/dav/space` endpoint


### Configurable path component in URLs

Not every deployment may have the requirement to have the path in the URL. We could use id only based URLs, similar to onedrive and make showing paths configurable.


| URL | comment |
|-|-|
| `https://<host>/files?id=<resource_id>` | default id based navigation |
| `https://<host>/files</namespaced/alias></relative/path/to/resource>?id=<resource_id>` | optional path based navigation with fallback to id |

In contrast to ownCloud 10 path takes precedence and the user is warned when the fileid in his bookmark no longer matches the id on the server: sth. like "The path of the resource has changed, please verify and update your bookmark!"

When a file is selected the filename also becomes part of the URL so individual files can be bookmarked.

If navigation is id based we need to look up the path for the id so we can make a webdav request, or we need to implement the graph drives and driveItem resources.

The URL  `https://<host>/files?id=<resource_id>̀` is sent to the server. It has to look up the correct path and redirect the request, including the the path. But that would make all bookmarks contain tha path again, even if paths were configured to not be part of the URL.

The `/meta/<fileid>` webdav endpoint can be used to look up the path with property `meta-path-for-user`.

For now, we would use path based navigation with URLs like this:

```
https://<host>/files</namespaced/alias></relative/path/to/resource>?id=<resource_id>
```

This means that only the _resource path_ is part of the URL path. Any other parameter, eg. file `id`, `page` or sort order must be given as URL parameters.

- [ ] To make lookup by id possible we need to implement the `/meta/<fileid>` endpoint so the sdk can use it to look up the path. We should not implement a redirect on the ocis server side because the same redirect logic would need to be added to oc10. Having it in ocis web is the right place.

- [ ] The old sharing links and oc10 urls still need to be redirected by ocis/reva as in oc10.

Public links would have the same format: `https://<host>/files?id=<resource_id>` The web UI has to detect if the user is logged in or not and adjust the ui accordingly.

{{< hint warning >}}
Since there is no difference between public and private files a logged in user cannot see the public version of a link unless he logs out.
{{< /hint >}}