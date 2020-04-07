# ownCloud Infinite Scale: Proxy

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-proxy/status.svg)](https://cloud.drone.io/owncloud/ocis-proxy)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d005a4722c1b463b9b95060479018e99)](https://www.codacy.com/gh/owncloud/ocis-proxy?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-proxy&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-proxy?status.svg)](http://godoc.org/github.com/owncloud/ocis-proxy)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-proxy)](http://goreportcard.com/report/github.com/owncloud/ocis-proxy)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-proxy.svg)](http://microbadger.com/images/owncloud/ocis-proxy "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/proxy/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/ocis-proxy/)
****
## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.13.

```console
git clone https://github.com/owncloud/ocis-proxy.git
cd ocis-proxy

make generate build

./bin/ocis-proxy -h
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
