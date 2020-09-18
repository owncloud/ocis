# ownCloud Infinite Scale: Pkg

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-pkg/status.svg)](https://cloud.drone.io/owncloud/ocis-pkg)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/27b2cb74a61547329f9b3c56d90bd05c)](https://www.codacy.com/manual/owncloud/ocis-pkg?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-pkg&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-pkg?status.svg)](http://godoc.org/github.com/owncloud/ocis-pkg)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-pkg)](http://goreportcard.com/report/github.com/owncloud/ocis-pkg)

This package defines some boilerplate code that reduces the code duplication within the ownCloud Infinite Scale microservice architecture. It can't be used standalone as the is a pure library. For further information about the available packages please read the source code or take a loog at [GoDoc](http://godoc.org/github.com/owncloud/ocis-pkg).

## Install

Import the required packages within your ownCloud Infinite Scale extensions and you are good to go.

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.13.

```console
git clone https://github.com/owncloud/ocis-pkg.git
cd ocis-pkg

make vet
make staticcheck
make lint
make test
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
