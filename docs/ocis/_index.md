---
title: "oCIS - ownCloud Infinite Scale"
date: 2020-02-27T20:35:00+01:00
weight: -10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: _index.md
---

{{< figure class="floatright" src="/media/is.png" width="70%" height="auto" >}}

## ownCloud Infinite Scale

Welcome to oCIS, the modern file-sync and share platform, which is based on our knowledge and experience with the PHP based [ownCloud server](https://owncloud.com/#server).

### oCIS runtime

The oCIS runtime allows us to dynamically manage ocis extensions running in a single process. We use [suture](https://github.com/thejerf/suture) to create a supervisor tree that start all ocis extensions in a dedicated goroutine. As oCIS is a micro-service based platform, individual services can be scaled-out to other nodes to meet your specific performance requirements.

### oCIS extensions

Every ocis extension uses [ocis-pkg](https://github.com/owncloud/ocis/ocis-pkg), which implements the [go-micro](https://go-micro.dev/) interfaces for [servers](https://github.com/asim/go-micro/blob/v3.5.0/server/server.go#L17-L37) to register and [clients](https://github.com/asim/go-micro/blob/v3.5.0/client/client.go#L11-L23) to lookup nodes with a service [registry](https://github.com/asim/go-micro/blob/v3.5.0/registry/registry.go).
We love and stick to [12 Factor](https://12factor.net/), the uniformity of services also allows us to use the same command, logging and configuration mechanism and pass configuration from the oCIS runtime to individual extensions.


### go-micro

While the [go-micro](https://go-micro.dev/) framework provides abstractions as well as implementations for the different components in a micro service architecture it uses a more developer focused runtime philosophy: it is used to download a services from a repo, compile them on the fly and start them as individual processes. For oCIS we decided to use a more admin friendly runtime: you can download a single binary and start the contained ocis extensions with a single `bin/ocis server`. This also makes packaging easier.

We use[ocis-pkg](https://github.com/owncloud/ocis/ocis-pkg) to configure the default implementations for the go-micro [grpc server](https://github.com/asim/go-micro/tree/v3.5.0/plugins/server/grpc), [client](https://github.com/asim/go-micro/tree/v3.5.0/plugins/client/grpc) and [mdns registry](https://github.com/asim/go-micro/blob/v3.5.0/registry/mdns_registry.go), swapping them out as needed, eg. to use the [kubernetes registry plugin](https://github.com/asim/go-micro/tree/v3.5.0/plugins/registry/kubernetes).

### REVA
A lot of services that oCIS is built upon are started using the [REVA](https://reva.link/) runtime. We decided to bundle some of the [CS3 services](https://github.com/cs3org/cs3apis) to logically group them. A [home storage provider](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/storagehome.go#L93-L108), which is dealing with [metadata](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ProviderAPI), and the corresponding [data provider](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/storagehome.go#L109-L123), which is dealing with [up and download](https://cs3org.github.io/cs3apis/#cs3.gateway.v1beta1.FileUploadProtocol), are one example. The [frontend](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/frontend.go) with the [oc flavoured webdav](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/frontend.go#L132-L138), [ocs handlers](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/frontend.go#L139-L148) and a [datagateway](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/frontend.go#L126-L131) are another.

### Protocol driven development
Interacting with oCIS involves a multitude af APIs. The server and all clients rely on [OpenID Connect](https://openid.net/connect/) for authentication. The [embedded konnectd](https://github.com/owncloud/ocis/tree/master/idp) can be replaced with any other OpenID Connect Identity Provider. Clients use the [WebDAV](http://webdav.org/) based [oc sync protocol](https://github.com/cernbox/smashbox/blob/master/protocol/protocol.md) to manage files and folders, [ocs to manage shares](https://doc.owncloud.com/server/developer_manual/core/apis/ocs-share-api.html) and [TUS](https://tus.io/protocols/resumable-upload.html) to upload files in a resumable way. On the server side [REVA](https://reva.link/) is the reference implementation of the [CS3 apis](https://github.com/cs3org/cs3apis) which is defined using [protobuf](https://developers.google.com/protocol-buffers/). We are looking into [the Microsoft Graph API](https://docs.microsoft.com/en-us/graph/api/overview?view=graph-rest-1.0), which is based on [odata](http://docs.oasis-open.org/odata/odata/v4.0/odata-v4.0-part1-protocol.html) as a rest/json implementation for current and future endpoints.

### Acceptance test suite
We run a huge [test suite](https://github.com/owncloud/core/tree/master/tests), which originated in ownCloud 10 and continues to grow. A detailed description can be found in the developer docs for [testing]({{< relref "development/testing.md" >}}).

### Architecture Overview

Running `bin/ocis server` will start the below services, all of which can be scaled and deployed on a single node or in a cloud native environment, as needed.

{{< svg src="ocis/static/architecture-overview.drawio.svg" >}}
