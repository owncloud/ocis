# Reva: WebDAV

[![Build Status](https://cloud.drone.io/api/badges/owncloud/reva-webdav/status.svg)](https://cloud.drone.io/owncloud/reva-webdav)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/afe89eb0894848c5b67dc0343afd1df9)](https://www.codacy.com/app/owncloud/reva-webdav?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/reva-webdav&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/reva-webdav?status.svg)](http://godoc.org/github.com/owncloud/reva-webdav)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/reva-webdav)](http://goreportcard.com/report/github.com/owncloud/reva-webdav)
[![](https://images.microbadger.com/badges/image/owncloud/reva-webdav.svg)](http://microbadger.com/images/owncloud/reva-webdav "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/reva/webdav/).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```console
git clone https://github.com/owncloud/reva-webdav.git
cd reva-webdav

make generate build

./bin/reva-webdav -h
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
