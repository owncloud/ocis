---
title: "Build"
date: 2020-02-27T20:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: build.md
---

{{< toc >}}

## Build requirements

see [Development - Getting Started]({{< relref "getting-started.md/#requirements">}})

## Get the sources

{{< highlight txt >}}
git clone https://github.com/owncloud/ocis.git
cd ocis
{{< / highlight >}}

## Build the oCIS binary

The oCIS binary source is in the `ocis` folder inside the oCIS repository. In this folder you can build the oCIS binary:

{{< highlight txt >}}
cd ocis
make generate
make build
{{< / highlight >}}

After building you have the binary within the `bin/` folder. Try to run it: `./bin/ocis -h`

## Build a local oCIS docker image

If you are developing and want to run your local changes in a docker or docker-compose setup, you have to build an image locally.

Therefore run following commands in the root of the oCIS repository:

{{< highlight txt >}}
docker build -t owncloud/ocis:dev .
{{< / highlight >}}

Then you can test as usual via

{{< highlight txt >}}
docker run --rm -ti owncloud/ocis:dev
{{< / highlight >}}
