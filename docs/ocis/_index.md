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

### The idea of federated storage

To create a truly federated storage architecture oCIS breaks down the old ownCloud 10 user specific namespace, which is assembled on the server side, and makes the individual parts accessible to clients as storage spaces and storage space registries.

The below diagram shows the core concepts that are the foundation for the new architecture:
- End user devices can fetch the list of *storage spaces* a user has access to, by querying one or multiple *storage space registries*. The list contains a unique endpoint for every *storage space*.
- [*Storage space registries*]({{< ref "../extensions/storage/terminology#storage-space-registries" >}}) manage the list of storage spaces a user has access to. They may subscribe to *storage spaces* in order to receive notifications about changes on behalf of an end users mobile or desktop client.
- [*Storage spaces*]({{< ref "../extensions/storage/terminology#storage-spaces" >}}) represent a collection of files and folders. A users personal files are a *storage space*, a group or project drive is a *storage space*, and even incoming shares are treated and implemented as *storage spaces*. Each with properties like owners, permissions, quota and type.
- [*Storage providers*]({{< ref "../extensions/storage/terminology#storage-providers" >}}) can hold multiple *storage spaces*. At an oCIS instance, there might be a dedicated *storage provider* responsible for users personal storage spaces. There might be multiple, sharing the load or there might be just one, hosting all types of *storage spaces*.

{{< svg src="ocis/static/idea.drawio.svg" >}}

As an example, Einstein might want to share something with Marie, who has an account at a different identity provider and uses a different storage space registry. The process makes use of [OpenID Connect (OIDC)](https://openid.net/specs/openid-connect-core-1_0.html) for authentication and would look something like this:

To share something with Marie, Einstein would open `https://cloud.zurich.test`. His browser loads oCIS web and presents a login form that uses the [OpenID Connect Discovery](https://openid.net/specs/openid-connect-discovery-1_0.html#EmailSyntax) to look up the OIDC issuer. For `einstein@zurich.test` he will end up at `https://idp.zurich.test`, authenticate and get redirected back to `https://cloud.zurich.test`. Now, oCIS web will use a similar discovery to look up the *storage space registry* for the account, based on the email (or username). He will discover that `https://cloud.zurich.test` is also his *storage registry* that the web UI will use to load the list of *storage spaces* that are available to him.

After locating a folder that he wants to share with Marie he enters her email `marie@paris.test` in the sharing dialog to grant her the editor role. This, in effect, creates a new *storage space* that is registered with the *storage space registry* at `https://cloud.zurich.test`.

Einstein copies the URL in the browser (or an email with the same URL is sent automatically, or the storage registries use a backchannel mechanism). It contains the most specific `storage space id` and a path relative to it: `https://cloud.zurich.test/#/spaces/716199a6-00c0-4fec-93d2-7e00150b1c84/a/rel/path`.

When Marie enters that URL she will be presented with a login form on the `https://cloud.zurich.test` instance, because the share was created on that domain. If `https://cloud.zurich.test` trusts her OpenID Connect identity provider `https://idp.paris.test` she can log in. This time, the *storage space registry* discovery will come up with `https://cloud.paris.test` though. Since that registry is different than the registry tied to `https://cloud.zurich.test` oCIS web can look up the *storage space* `716199a6-00c0-4fec-93d2-7e00150b1c84` and register the WebDAV URL `https://cloud.zurich.test/dav/spaces/716199a6-00c0-4fec-93d2-7e00150b1c84/a/rel/path` in Maries *storage space registry* at `https://cloud.paris.test`. When she accepts that share her clients will be able to sync the new *storage space* at `https://cloud.zurich.test`.

### oCIS microservice runtime

The oCIS runtime allows us to dynamically manage services running in a single process. We use [suture](https://github.com/thejerf/suture) to create a supervisor tree that starts each service in a dedicated goroutine. By default oCIS will start all built-in oCIS extensions in a single process. Individual services can be moved to other nodes to scale-out and meet specific performance requirements. A [go-micro](https://github.com/asim/go-micro/blob/master/registry/registry.go) based registry allows services in multiple nodes to form a distributed microservice architecture.

### oCIS extensions

Every oCIS extension uses [ocis-pkg](https://github.com/owncloud/ocis/tree/master/ocis-pkg), which implements the [go-micro](https://go-micro.dev/) interfaces for [servers](https://github.com/asim/go-micro/blob/v3.5.0/server/server.go#L17-L37) to register and [clients](https://github.com/asim/go-micro/blob/v3.5.0/client/client.go#L11-L23) to lookup nodes with a service [registry](https://github.com/asim/go-micro/blob/v3.5.0/registry/registry.go).
We are following the [12 Factor](https://12factor.net/) methodology with oCIS. The uniformity of services also allows us to use the same command, logging and configuration mechanism. Configurations are forwarded from the 
oCIS runtime to the individual extensions.


### go-micro

While the [go-micro](https://go-micro.dev/) framework provides abstractions as well as implementations for the different components in a microservice architecture, it uses a more developer focused runtime philosophy: It is used to download services from a repo, compile them on the fly and start them as individual processes. For oCIS we decided to use a more admin friendly runtime: You can download a single binary and start the contained oCIS extensions with a single `bin/ocis server`. This also makes packaging easier.

We use [ocis-pkg](https://github.com/owncloud/ocis/tree/master/ocis-pkg) to configure the default implementations for the go-micro [grpc server](https://github.com/asim/go-micro/tree/v3.5.0/plugins/server/grpc), [client](https://github.com/asim/go-micro/tree/v3.5.0/plugins/client/grpc) and [mdns registry](https://github.com/asim/go-micro/blob/v3.5.0/registry/mdns_registry.go), swapping them out as needed, eg. to use the [kubernetes registry plugin](https://github.com/asim/go-micro/tree/v3.5.0/plugins/registry/kubernetes).

### REVA
A lot of embedded services in oCIS are built upon the [REVA](https://reva.link/) runtime. We decided to bundle some of the [CS3 services](https://github.com/cs3org/cs3apis) to logically group them. A [home storage provider](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/storagehome.go#L93-L108), which is dealing with [metadata](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ProviderAPI), and the corresponding [data provider](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/storagehome.go#L109-L123), which is dealing with [up and download](https://cs3org.github.io/cs3apis/#cs3.gateway.v1beta1.FileUploadProtocol), are one example. The [frontend](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/frontend.go) with the [oc flavoured webdav](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/frontend.go#L132-L138), [ocs handlers](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/frontend.go#L139-L148) and a [datagateway](https://github.com/owncloud/ocis/blob/v1.2.0/storage/pkg/command/frontend.go#L126-L131) are another.

### Protocol driven development
Interacting with oCIS involves a multitude af APIs. The server and all clients rely on [OpenID Connect](https://openid.net/connect/) for authentication. The [embedded LibreGraph Connect](https://github.com/owncloud/ocis/tree/master/idp) can be replaced with any other OpenID Connect Identity Provider. Clients use the [WebDAV](http://webdav.org/) based [oc sync protocol](https://github.com/cernbox/smashbox/blob/master/protocol/protocol.md) to manage files and folders, [ocs to manage shares](https://doc.owncloud.com/server/developer_manual/core/apis/ocs-share-api.html) and [TUS](https://tus.io/protocols/resumable-upload.html) to upload files in a resumable way. On the server side [REVA](https://reva.link/) is the reference implementation of the [CS3 apis](https://github.com/cs3org/cs3apis) which is defined using [protobuf](https://developers.google.com/protocol-buffers/). By embedding [glauth](https://github.com/glauth/glauth/), oCIS provides a read-only [LDAP](https://tools.ietf.org/html/rfc2849) interface to make accounts, including guests available to firewalls and other systems. In the future, we are looking into [the Microsoft Graph API](https://docs.microsoft.com/en-us/graph/api/overview?view=graph-rest-1.0), which is based on [odata](http://docs.oasis-open.org/odata/odata/v4.0/odata-v4.0-part1-protocol.html), as a well defined REST/JSON dialect for the existing endpoints.

### Acceptance test suite
We run a huge [test suite](https://github.com/owncloud/core/tree/master/tests), which originated in ownCloud 10 and continues to grow. A detailed description can be found in the developer docs for [testing]({{< ref "development/testing" >}}).

### Architecture Overview

Running `bin/ocis server` will start the below services, all of which can be scaled and deployed on a single node or in a cloud native environment, as needed.

{{< svg src="ocis/static/architecture-overview.drawio.svg" >}}
