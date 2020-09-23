# ownCloud Infinite Scale: THUMBNAILS

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-thumbnails/status.svg)](https://cloud.drone.io/owncloud/ocis-thumbnails)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5dcc0c9b2319462dbbe517d90219062c)](https://www.codacy.com/gh/owncloud/ocis-thumbnails?utm_source=github.com&utm_medium=referral&utm_content=owncloud/ocis-thumbnails&utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-thumbnails?status.svg)](http://godoc.org/github.com/owncloud/ocis-thumbnails)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-thumbnails)](http://goreportcard.com/report/github.com/owncloud/ocis-thumbnails)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-thumbnails.svg)](http://microbadger.com/images/owncloud/ocis-thumbnails "Get your own image badge on microbadger.com")
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=owncloud_ocis-thumbnails&metric=alert_status)](https://sonarcloud.io/dashboard?id=owncloud_ocis-thumbnails)

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/thumbnails/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/extensions/ocis_thumbnails/)

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.13.

```console
git clone https://github.com/owncloud/ocis-thumbnails.git
cd ocis-thumbnails

make generate build

./bin/ocis-thumbnails -h
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
