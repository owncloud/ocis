# ownCloud Infinite Scale: SETTINGS

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-settings/status.svg)](https://cloud.drone.io/owncloud/ocis-settings)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4383f209aa904572b415ef5a8f9e379f)](https://www.codacy.com/gh/owncloud/ocis-settings?utm_source=github.com&utm_medium=referral&utm_content=owncloud/ocis-settings&utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-settings?status.svg)](http://godoc.org/github.com/owncloud/ocis-settings)
[![Go Report Card](https://goreportcard.com/badge/github.com/owncloud/ocis-settings)](https://goreportcard.com/report/github.com/owncloud/ocis-settings)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-settings.svg)](https://microbadger.com/images/owncloud/ocis-settings "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/settings/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/extensions/ocis_settings/)

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.13.

```console
git clone https://github.com/owncloud/ocis-settings.git
cd ocis-settings

make generate build

./bin/ocis-settings -h
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
