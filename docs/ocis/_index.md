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

### oCIS server

The oCIS server implementation follows [Go](https://golang.org/) best practices and is based on the [go-micro](https://go-micro.dev/) framework and [REVA](https://reva.link/). We love and stick to [12 Factor](https://12factor.net/).
oCIS is a micro-service based server, which allows scale-out of individual services to meet your specific performance requirements.
We run a huge [test suite](https://github.com/owncloud/core/tree/master/tests), which originated in ownCloud 10 and continues to grow.

### Architecture Overview


{{< svg src="ocis/static/architecture-overview.drawio.svg" >}}
