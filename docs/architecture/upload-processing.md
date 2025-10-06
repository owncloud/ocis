---
title: "Upload processing"
date: 2022-07-06T12:47:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/architecture
geekdocFilePath: upload-processing.md
---

Uploads are handled by a dedicated service that uses TUS.io for resumable uploads. When all bytes have been transferred the upload is finalized by making the file available in file listings and for download.

The finalization may be asynchronous when mandatory workflow steps are involved.

## Legacy PUT upload

{{<mermaid class="text-center">}}

%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant Client
    participant ocdav
    participant storageprovider
    participant dataprovider

    Client->>+ocdav: PUT /dav/spaces/{spaceid}/newfile.bin
    ocdav->>+storageprovider: InitiateFileUpload
    storageprovider-->>-ocdav: OK, Protocol simple, UploadEndpoint: /data, Token: {jwt}
    Note right of ocdav: The {jwt} contains the internal actual target, eg.: http://localhost:9158/data/simple/91cc9882-db71-4b37-b694-a522850fcee1
    ocdav->>+dataprovider: PUT /data
    Note right of dataprovider: X-Reva-Transfer: {jwt}
    dataprovider-->>-ocdav: 201 Created
    ocdav-->>-Client: 201 Created

{{</mermaid>}}

## TUS upload

{{<mermaid class="text-center">}}

%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant Client
    participant ocdav
    participant storageprovider
    participant datagateway
    participant dataprovider

    Client->>+ocdav: POST /dav/spaces/{spaceid}\nUpload-Metadata: {base64 encoded filename etc}\nTUS-Resumable: 1.0.0
    ocdav->>+storageprovider: InitiateFileUpload
    storageprovider-->>-ocdav: OK, Protocol tus, UploadEndpoint: /data, Token: {jwt}
    Note right of ocdav: The {jwt} contains the internal actual target, eg.:\nhttp://localhost:9158/data/tus/24d893f5-b942-4bc7-9fb0-28f49f980160
    ocdav-->>-Client: 201 Created\nLocation: /data/{jwt}\nTUS-Resumable: 1.0.0

    Client->>+datagateway: PATCH /data/{jwt}\nTUS-Resumable: 1.0.0\nUpload-Offset: 0
    Note over datagateway: unwrap the {jwt} target
    datagateway->>+dataprovider: PATCH /data/tus/24d893f5-b942-4bc7-9fb0-28f49f980160\nX-Reva-Transfer: {jwt}
    Note over dataprovider: storage driver\nhandles request
    dataprovider-->>-datagateway: 204 No Content\nTUS-Resumable: 1.0.0\nUpload-Offset: 363976
    datagateway-->>-Client: 204 No Content\nTUS-Resumable: 1.0.0\nUpload-Offset: 363976

{{</mermaid>}}


## TUS upload with async postprocessing



{{<mermaid class="text-center">}}

%%{init: {"sequence": { "showSequenceNumbers":true, "messageFontFamily":"courier", "messageFontWeight":"normal", "messageFontSize":"11"}}}%%
%% font weight is a css bug: https://github.com/mermaid-js/mermaid/issues/1976
%% edit this diagram by pasting it into eg. https://mermaid.live
sequenceDiagram
    participant Client
    participant ocdav
    participant storageprovider
    participant datagateway
    participant dataprovider
    participant nats
    participant processing

    Client->>+ocdav: POST /dav/spaces/{spaceid}
    Note left of Client: Upload-Metadata: {base64 encoded filename etc}\nTUS-Resumable: 1.0.0
    ocdav->>+storageprovider: InitiateFileUpload
    storageprovider-->>-ocdav: OK, Protocol tus, UploadEndpoint: /data, Token: {jwt}
    Note right of ocdav: The {jwt} contains the internal actual target, eg.: http://localhost:9158/data/tus/24d893f5-b942-4bc7-9fb0-28f49f980160
    ocdav-->>-Client: 201 Created
    Note right of Client: Location: /data/{jwt}
    Note right of Client: TUS-Resumable: 1.0.0

    Client->>+datagateway: PATCH /data/{jwt}
    Note right of datagateway: TUS-Resumable: 1.0.0\nUpload-Offset: 0

    Note over datagateway: unwrap the {jwt} target
    datagateway->>+dataprovider: PATCH /data/tus/24d893f5-b942-4bc7-9fb0-28f49f980160
    Note over dataprovider: X-Reva-Transfer: {jwt}
    Note over dataprovider: storage driver
    Note over dataprovider: handles request
    dataprovider-)nats: emit all-bytes-received event
    nats-)processing: all-bytes-received({uploadid}) event
    Note over dataprovider: TODO: A lot of time may pass here, we could use the `Prefer: respond-async` header to return early with a 202 Accepted status and a Location header to a websocket endpoint
    alt success
        processing-)nats: emit processing-finished({uploadid}) event
        nats-)dataprovider: processing-finished({uploadid}) event
        dataprovider-->>-datagateway: 204 No Content
        Note over datagateway: TUS-Resumable: 1.0.0\nUpload-Offset: 363976
        datagateway-->>-Client: 204 No Content
        Note over Client: TUS-Resumable: 1.0.0\nUpload-Offset: 363976
    else failure
        activate dataprovider
        activate datagateway
        processing-)nats: emit processing-aborted({uploadid}) event
        nats-)dataprovider: processing-aborted({uploadid}) event
        Note over dataprovider: FIXME: What HTTP status code should we report?
        Note over dataprovider: 422 Unprocessable Content is just a proposal
        Note over dataprovider: see https://httpwg.org/specs/rfc9110.html#status.422
        dataprovider-->>-datagateway: 422 Unprocessable Content
        Note over datagateway: TUS-Resumable: 1.0.0\nUpload-Offset: 363976
        datagateway-->>-Client: 422 Unprocessable Content
        Note over Client: TUS-Resumable: 1.0.0\nUpload-Offset: 363976
    end

{{</mermaid>}}


## Async TUS upload with postprocessing
This might be a TUS extension or a misunderstanding on our side of what tus can do for us. Clients should send a `Prefer: respond-async` header to allow the server to return early when postprocessing might take longer. The PATCH requests can then return status `202 Accepted` and a `Location` header to a websocket that clients can use to track the processing / upload progress.

TODO there is a conflict with the TUS.io POST request with the creation extension, as that also returns a `Location` header which carries the upload URL. We would need another header to transport the websocket location. Maybe `Websocket-Location` or `Progress-Location`?
