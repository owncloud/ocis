---
title: "Building"
date: 2018-05-02T00:00:00+00:00
anchor: "building"
weight: 30
---

As this project is built with Go, so you need to install that first. The installation of Go is out of the scope of this document, please follow the official documentation for [Go](golang). After the installation of the required tools you need to get the sources:

{{< highlight txt >}}
git clone https://github.com/owncloud/ocis-ocs.git
cd ocis-ocs
{{< / highlight >}}

All required tool besides Go itself and make are bundled or getting automatically installed within the `GOPATH`. All commands to build this project are part of our `Makefile` and respectively our `package.json`.

## Backend

{{< highlight txt >}}
make generate
make build
{{< / highlight >}}

Finally you should have the binary within the `bin/` folder now, give it a try with `./bin/ocis-ocs -h` to see all available options.

[golang]: https://golang.org/doc/install
[nodejs]: https://nodejs.org/en/download/package-manager/
[yarn]: https://yarnpkg.com/lang/en/docs/install/
