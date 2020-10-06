---
title: "Building"
date: 2018-05-02T00:00:00+00:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/web
geekdocFilePath: building.md
---

## Backend

{{< highlight txt >}}
cd ocis-phoenix
make generate
make build
{{< / highlight >}}

The above commands will download a [Phoenix](https://github.com/owncloud/phoenix) release and embed it into the binary. Finally you should have the binary within the `bin/` folder now, give it a try with `./bin/ocis-phoenix -h` to see all available options.
