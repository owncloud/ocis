# ownCloud Infinite Scale: Konnectd

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-konnectd/status.svg)](https://cloud.drone.io/owncloud/ocis-konnectd)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/f6f9033737404c9da3ba4738b6501bdb)](https://www.codacy.com/manual/owncloud/ocis-konnectd?utm_source=github.com&utm_medium=referral&utm_content=owncloud/ocis-konnectd&utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-konnectd?status.svg)](http://godoc.org/github.com/owncloud/ocis-konnectd)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-konnectd)](http://goreportcard.com/report/github.com/owncloud/ocis-konnectd)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-konnectd.svg)](http://microbadger.com/images/owncloud/ocis-konnectd "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/konnectd/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/extensions/ocis_konnectd/)

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.13.

```console
git clone https://github.com/owncloud/ocis-konnectd.git
cd ocis-konnectd

make generate build

./bin/ocis-konnectd -h
```

## Security

If you find a security issue please contact [security@owncloud.com](mailto:security@owncloud.com) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2020 ownCloud GmbH <https://owncloud.com>
```
