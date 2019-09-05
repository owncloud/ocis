# ownCloud Infinit Scale: WebDAV

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-webdav/status.svg)](https://cloud.drone.io/owncloud/ocis-webdav)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/0913fcc866a344b587bb867fcec5b848)](https://www.codacy.com/app/owncloud/ocis-webdav?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-webdav&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-webdav?status.svg)](http://godoc.org/github.com/owncloud/ocis-webdav)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-webdav)](http://goreportcard.com/report/github.com/owncloud/ocis-webdav)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-webdav.svg)](http://microbadger.com/images/owncloud/ocis-webdav "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/webdav/).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```console
git clone https://github.com/owncloud/ocis-webdav.git
cd ocis-webdav

make generate build

./bin/ocis-webdav -h
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
