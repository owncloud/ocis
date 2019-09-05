# ownCloud Infinite Scale: OCS

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-ocs/status.svg)](https://cloud.drone.io/owncloud/ocis-ocs)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/80a22dcfa3cb4f09ba8f63b386683d16)](https://www.codacy.com/app/owncloud/ocis-ocs?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-ocs&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-ocs?status.svg)](http://godoc.org/github.com/owncloud/ocis-ocs)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-ocs)](http://goreportcard.com/report/github.com/owncloud/ocis-ocs)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-ocs.svg)](http://microbadger.com/images/owncloud/ocis-ocs "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/ocs/).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```console
git clone https://github.com/owncloud/ocis-ocs.git
cd ocis-ocs

make generate build

./bin/ocis-ocs -h
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
