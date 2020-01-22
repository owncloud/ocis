# ownCloud Infinite Scale: Graph-Explorer

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-graph-explorer/status.svg)](https://cloud.drone.io/owncloud/ocis-graph-explorer)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/afe89eb0894848c5b67dc0343afd1df9)](https://www.codacy.com/app/owncloud/ocis-graph-explorer?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-graph-explorer&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-graph-explorer?status.svg)](http://godoc.org/github.com/owncloud/ocis-graph-explorer)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-graph-explorer)](http://goreportcard.com/report/github.com/owncloud/ocis-graph-explorer)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-graph-explorer.svg)](http://microbadger.com/images/owncloud/ocis-graph-explorer "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/graph-explorer/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/ocis-graph-explorer/)

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.12. For the frontend it's also required to have [NodeJS](https://nodejs.org/en/download/package-manager/) and [Yarn](https://yarnpkg.com/lang/en/docs/install/) installed.

```console
git clone https://github.com/owncloud/ocis-graph-explorer.git
cd ocis-graph-explorer

make generate build

./bin/ocis-graph-explorer -h
```

## Security

If you find a security issue please contact security@owncloud.com first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2019 ownCloud GmbH <https://owncloud.com>
```
