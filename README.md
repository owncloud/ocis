# Reva: OCS

[![Build Status](https://cloud.drone.io/api/badges/owncloud/reva-ocs/status.svg)](https://cloud.drone.io/owncloud/reva-ocs)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/afe89eb0894848c5b67dc0343afd1df9)](https://www.codacy.com/app/owncloud/reva-ocs?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/reva-ocs&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/reva-ocs?status.svg)](http://godoc.org/github.com/owncloud/reva-ocs)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/reva-ocs)](http://goreportcard.com/report/github.com/owncloud/reva-ocs)
[![](https://images.microbadger.com/badges/image/owncloud/reva-ocs.svg)](http://microbadger.com/images/owncloud/reva-ocs "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/reva/ocs/).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```console
git clone https://github.com/owncloud/reva-ocs.git
cd reva-ocs

make generate build

./bin/reva-ocs -h
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
