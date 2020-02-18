---
title: "Building"
date: 2018-05-02T00:00:00+00:00
anchor: "building"
weight: 30
---

As this project is built with Go, so you need to install that first. The installation of Go is out of the scope of this document, please follow the official documentation for [Go](https://golang.org/doc/install), to build this project you have to install Go >= v1.13. After the installation of the required tools you need to get the sources:

{{< highlight txt >}}
git clone https://github.com/owncloud/ocis.git
cd ocis
{{< / highlight >}}

All required tool besides Go itself and make are bundled or getting automatically installed within the `GOPATH`. All commands to build this project are part of our `Makefile`.

### Simple Ocis

Building the simple ocis binary:

```console
TAGS=simple make build
```

The artifact lives in `/bin/ocis`

The generated simple ocis binary is a subset of the ocis command with a restricted set of services meant for ease up development. The services included are


```
ocis-hello
phoenix
konnectd
devldap
micro's own services
```

### Backend

{{< highlight txt >}}
make generate
make build
{{< / highlight >}}

Finally you should have the binary within the `bin/` folder now, give it a try with `./bin/ocis -h` to see all available options.
