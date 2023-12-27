---
title: "Protocol changes"
date: 2022-05-17T08:46:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/architecture
geekdocFilePath: protocol-changes.md
---

The spaces concept allows clients to look up the space endpoints a user has access to and then do individual sync discoveries. Technically, we introduce an indirection that allows clients to rely on server provided URLs instead of hardcoded `/webdav` or `/dav/files/{username}` paths, that may change over time.

## Space discovery

{{<mermaid class="text-center">}}
%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant Client
    participant Graph
    participant SpaceA
    participant SpaceB
    links Client: {"web": "https://owncloud.dev/clients/web/", "RClone": "https://owncloud.dev/clients/rclone/"}
    link Graph: Documentation @ https://owncloud.dev/extensions/graph/

    Note left of Client:  First, a clients looks<br/>up the spaces a user has access to
    opt space lookup
        Client->>+Graph: GET /me/drives
        Graph-->>-Client: 200 OK JSON list of spaces, say A, B and C,<br/> each with a dedicated webDavURL, etag and quota
    end

    Note left of Client: Then it can do a parallel<br/>sync discovery on spaces<br/>whose etag changed
    par Client to Space A
        Client->>+SpaceA: PROPFIND {webDavURL for Space A}
        SpaceA-->>-Client: 207 Multistatus PROPFIND response
    and Client to Space B
        Client->>+SpaceB: PROPFIND {webDavURL for space B}
        SpaceB-->>-Client: 207 Multistatus PROPFIND response
    end
{{</mermaid>}}

### New /dav/spaces/{spaceid} endpoint with spaceid and a relative path

The ocDAV service is responsible for translating ownCloud flavoured WebDAV into CS3 API calls.

**General view**

A PROPFIND finds its way to a storage provider like this:

{{<mermaid class="text-center">}}
%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant Client
    participant ocDAV
    participant StorageProvider

    Note right of Client: {spaceid} identifies the space<br>{relative/path} is relative to the space root
        Client->>+ocDAV: PROPFIND /dav/space/{spaceid}/{relative/path}
    Note right of ocDAV: translate ownCloud flavoured webdav<br>into CS3 API requests
        ocDAV->>+StorageProvider: ListContainer({spaceid}, path: {relative/path})
        StorageProvider-->>-ocDAV: []ResourceInfo
        ocDAV-->>-Client: 207 Multistatus
{{</mermaid>}}

While the above is a simplification to get an understanding of what needs to go where, there are several places where sharding can happen.

**Proxy can do user based routing**

The ocis proxy authenticates requests and can forward requests to different backends, depending on the logged-in user or cookies. For example multiple ocdav services can be configured to shard users based on username or affiliation.

{{<mermaid class="text-center">}}
%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant Client
    participant proxy
    participant ocDAV1 as ocDAV [a-k]
    participant ocDAV2 as ocDAV [l-z]

    Note right of Client: {spaceid} identifies the space<br>{relative/path} is relative to the space root
        Client->>+proxy: PROPFIND /dav/space/{spaceid}/{relative/path}

    alt username starting with a-k
        proxy->>+ocDAV1: PROPFIND /dav/space/{spaceid}/{relative/path}
    Note right of ocDAV1: translate ownCloud flavoured webdav<br>into CS3 API requests
        ocDAV1-->>-Client: 207 Multistatus
    else username starting with l-z
        proxy->>+ocDAV2: PROPFIND /dav/space/{spaceid}/{relative/path}
        ocDAV2-->>-Client: 207 Multistatus
    end
{{</mermaid>}}

**Gateway can do path or storage provider id based routing**

The reva gateway acts as a facade to multiple storage providers that can be configured with the storage registry:

{{<mermaid class="text-center">}}
%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant ocDAV
    participant Gateway
    participant StorageRegistry
    participant StorageProvider1 as StorageProvider [a-k]
    participant StorageProvider2 as StorageProvider [l-z]

    Note right of ocDAV: translate ownCloud flavoured webdav<br>into CS3 API requests
        ocDAV->>+Gateway: ListContainer({spaceid}, path: {relative/path})
    Note right of Gateway: find address of the storage provider<br>that is responsible for the space
        Gateway->>+StorageRegistry: ListStorageProviders({spaceid})
        StorageRegistry-->>-Gateway: []ProviderInfo
    Note right of Gateway: forward request to<br>correct storage provider
    alt username starting with a-k
        Gateway->>+StorageProvider1: ListContainer({spaceid}, path: {relative/path})
        StorageProvider1-->>-Gateway: []ResourceInfo
    else username starting with l-z
        Gateway->>+StorageProvider2: ListContainer({spaceid}, path: {relative/path})
        StorageProvider2-->>-Gateway: []ResourceInfo
    end
        Gateway-->>-ocDAV: []ResourceInfo
{{</mermaid>}}


### Old /dav/files/{username} endpoint with username and a path relative to the users home

**PROPFIND request against old webdav endpoints**

To route a PROPFIND request against the old webdav endpoints like `/dav/files/username`, ocdav first has to build a CS3 namespace prefix, e.g. `/users/{{.Id.OpaqueId}}` to the users home.

{{<mermaid class="text-center">}}
%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant Client
    participant ocDAV
    participant Gateway

    opt old /dav/files/{username} endpoint with username and a path relative to the users home
    Note right of Client: translate ownCloud flavoured webdav<br>into CS3 API requests
        Client->>+ocDAV: PROPFIND /dav/files/{username}/{relative/path}
    Note right of ocDAV: translate ownCloud flavoured webdav<br>into CS3 API requests
        ocDAV->>+Gateway: GetUser({username})
        Gateway-->>-ocDAV: User
    Note right of ocDAV: build path prefix to user home
        ocDAV->>+ocDAV: {namespace/prefix} = ApplyLayout({path layout}, User), eg. /users/e/einstein
    Note right of ocDAV: look up the space responsible for a path
        ocDAV->>+Gateway: ListStorageSpaces(path: {namespace/prefix}/{relative/path})
        Gateway-->>-ocDAV: []StorageSpace
    Note right of ocDAV: make actual request with space and relative path
        ocDAV->>+Gateway: ListContainer({spaceid}, path: {relative/path})
        Gateway-->>-ocDAV: []ResourceInfo
        ocDAV-->>-Client: 207 Multistatus
    end
{{</mermaid>}}

**Handling legacy global namespace webdav endpoints**

The reason ocis uses a path based lookup instead of looking up the current users home using the user id and a space type filter is, because there are deployments that use a global namespace at the legacy `/webdav` endpoint. To support these use cases, the gateway allows looking up spaces using their mount path.

{{<mermaid class="text-center">}}
%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant Client
    participant ocDAV
    participant Gateway

    Note right of Client: translate ownCloud flavoured webdav<br>into CS3 API requests
    alt old /dav/files/{username} endpoint with username and a path relative to the users home
        Client->>+ocDAV: PROPFIND /dav/files/{username}/{relative/path}
    Note right of ocDAV: look up {username} in URL path
        ocDAV->>+Gateway: GetUser({username})
        Gateway-->>-ocDAV: User
    Note right of ocDAV:build namespace prefix to user home
        ocDAV->>+ocDAV: {namespace/prefix} = ApplyLayout({namespace layout}, User), eg. /users/e/einstein
    else legacy /webdav/ endpoint with a path relative to the users home
        Client->>+ocDAV: PROPFIND /webdav/{relative/path}
    Note right of ocDAV: use currently logged in user
        ocDAV->>+ocDAV: ContextGetUser()
    Note right of ocDAV: build namespace prefix to user home
        ocDAV->>+ocDAV: {namespace/prefix} = ApplyLayout({namespace layout}, User), eg. /users/e/einstein
    else legacy /webdav/ endpoint with a path relative to a global namespace
        Client->>+ocDAV: PROPFIND /webdav/{relative/path}
    Note right of ocDAV: omit namespace prefix by using empty layout template
        ocDAV->>+ocDAV: {namespace/prefix} = ApplyLayout("/", u), always returns "/"
    end
    Note right of ocDAV: look up the space responsible for a path
        ocDAV->>+Gateway: ListStorageSpaces(path: {namespace/prefix}/{relative/path})
        Gateway-->>-ocDAV: []StorageSpace
    Note right of ocDAV: make actual request with space and relative path
        ocDAV->>+Gateway: ListContainer({spaceid}, path: {relative/path})
        Gateway-->>-ocDAV: []ResourceInfo
        ocDAV-->>-Client: 207 Multistatus
{{</mermaid>}}
