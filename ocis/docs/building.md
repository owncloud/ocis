* * *

title: "Building"
date: 2020-02-27T20:35:00+01:00
weight: 50
geekdocRepo: <https://github.com/owncloud/ocis>
geekdocEditPath: edit/master/docs

## geekdocFilePath: building.md

As this project is built with Go, so you need to install that first. The installation of Go is out of the scope of this document, please follow the official documentation for [Go](https://golang.org/doc/install), to build this project you have to install Go >= v1.13. After the installation of the required tools you need to get the sources:

{{&lt; highlight txt >}}
git clone <https://github.com/owncloud/ocis.git>
cd ocis
{{&lt; / highlight >}}

All required tools besides Go itself and make are bundled or getting automatically installed within the `GOPATH`. All commands to build this project are part of our `Makefile`. To build the `ocis` binary run:

{{&lt; highlight txt >}}
make generate
make build
{{&lt; / highlight >}}

Finally, you should have the binary within the `bin/` folder now, give it a try with `./bin/ocis -h` to see all available options.

## Simple Ocis fo extonsions example

Currently, we are using a go build tag to allow building a more simple set of the binary. It was intended to let extension developers focus on only the necessary services.

{{&lt; hint info >}}
While it the tag based simple build demonstrates how to use ocis as a framework for a micro service architecture, we may change to an approach that uses an explicit command to run only a subset of the services.
{{&lt; / hint >}}

```console
TAGS=simple make build
```

The artifact lives in `/bin/ocis`

The generated simple ocis binary is a subset of the ocis command with a restricted set of services meant for ease up development. The services included are

    ocis-hello
    ocis-phoenix
    ocis-konnectd
    ocis-glauth
    micro's own services
