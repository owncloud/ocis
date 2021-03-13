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

You only need to run following command if you have changed protobuf definitions or the frontend part in one of the extensions. Run the command in the root directory of the repository. Otherwise you can skip this step and proceed to build the oCIS binary.
This will usually modify multiple `embed.go` files because we embed the frontend build output in these `embed.go` files and a timestamp will be updated and also minor differences are expected between different Node.js versions.

{{< highlight txt >}}
make generate
{{< / highlight >}}

The next step is to build the actual oCIS binary. Therefore you need to navigate to the subdirectory `ocis` and start the build process.

{{< highlight txt >}}
cd ocis
make build
{{< / highlight >}}

After the build process finished, you can find the binary within the `bin/` folder (in `ocis/bin` relative to the oCIS repository root folder).

Try to run it: `./bin/ocis h`

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
