---
title: "Build ocis"
date: 2020-02-27T20:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: building.md
---

## Build requirements

All required tools besides `go` and `make` are bundled or getting automatically installed within the `GOPATH`. All commands to build this project are part of our `Makefile`.

The installation of Go is out of the scope of this document, please follow the official documentation for [Go](https://golang.org/doc/install), to build this project you have to install Go >= v1.13.

## Get the sources

{{< highlight txt >}}
git clone https://github.com/owncloud/ocis.git
cd ocis
{{< / highlight >}}

## Build the oCIS binary

The oCIS binary source is in the ocis/ocis folder. In this folder you can build the ocis binary:

{{< highlight txt >}}
make generate
make build
{{< / highlight >}}

Finally, you should have the binary within the `bin/` folder now, give it a try with `./bin/ocis -h` to see all available options.

## Build a local ocis docker image

If you are developing on a local branch based on docker / docker-compose setup, here is how to build a new ocis image. In the root folder:

{{< highlight txt >}}
docker build -t owncloud/ocis:dev .
{{< / highlight >}}

Then you can test as usual via

{{< highlight txt >}}
docker run --rm -ti owncloud/ocis:dev
{{< / highlight >}}
