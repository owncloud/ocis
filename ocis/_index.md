---
title: "oCIS"
date: 2020-02-27T20:35:00+01:00
weight: -10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: _index.md
---

{{< figure class="floatright" src="/media/is.png" width="70%" height="auto" >}}

## ownCloud Infinite Scale

Welcome to oCIS! We develop a modern file-sync and share plattform, based on our knowledge and experience with the PHP ownCloud server project.

### oCIS Server

The oCIS server implementation follows go-lang best practices and is based on the [go-micro](https://go-micro.dev/) framework and [REVA](https://reva.link/). We love and stick to [12 Factor](https://12factor.net/). 
oCIS is a micro-service based server, which allows scale-out of individual services to meet your specific performance requirements.
We run a huge test suite, which was originated in ownCloud 10 and continues to grow.

### Architecture Overview

{{< mermaid class="text-center">}}
graph TD
ocis-proxy -->
    ocis-konnectd & ocis-phoenix & ocis-thumbnails & ocis-ocs & ocis-webdav

ocis-phoenix --> ocis-reva-fronted
ocis-reva-fronted --> ocis-reva-gateway
ocis-konnectd --> ocis-glauth


ocis-reva-gateway --> ocis-reva-users
ocis-reva-gateway --> ocis-reva-authbasic
ocis-reva-gateway --> ocis-reva-auth-bearer

ocis-reva-gateway --> ocis-reva-sharing
ocis-reva-gateway --> ocis-reva-storage-home-*
ocis-reva-storage-home-* --> ocis-reva-storage-home-*-data
ocis-reva-sharing --> redis
{{< /mermaid >}}
