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
