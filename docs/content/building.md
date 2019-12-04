---
title: "Building"
date: 2018-05-02T00:00:00+00:00
anchor: "building"
weight: 30
---

As this project is built with Go and NodeJS, so you need to install that first. The installation of Go and NodeJS is out of the scope of this document, please follow the official documentation for [Go](golang), [NodeJS](nodejs) and [Yarn](yarn), to build this project you have to install Go >= v1.12. After the installation of the required tools you need to get the sources:

{{< highlight txt >}}
git clone https://github.com/owncloud/ocis-graph.git
cd ocis-graph
{{< / highlight >}}

All required tool besides Go itself and make are bundled or getting automatically installed within the `GOPATH`. All commands to build this project are part of our `Makefile` and respectively our `package.json`.

### Frontend

{{< highlight txt >}}
yarn install
yarn build
{{< / highlight >}}

The above commands will install the required build dependencies and build the whole frontend bundle. This bundle will we embeded into the binary later on.

### Backend

{{< highlight txt >}}
make generate
make build
{{< / highlight >}}

The above commands will embed the frontend bundle into the binary. Finally you should have the binary within the `bin/` folder now, give it a try with `./bin/ocis-graph -h` to see all available options.

[golang]: https://golang.org/doc/install
[nodejs]: https://nodejs.org/en/download/package-manager/
[yarn]: https://yarnpkg.com/lang/en/docs/install/
