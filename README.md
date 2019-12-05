# ownCloud Infinite Scale: Pkg

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-pkg/status.svg)](https://cloud.drone.io/owncloud/ocis-pkg)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6f1eaaa399294d959ef7b3b10deed41d)](https://www.codacy.com/manual/owncloud/ocis-pkg?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis-pkg&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-pkg?status.svg)](http://godoc.org/github.com/owncloud/ocis-pkg)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-pkg)](http://goreportcard.com/report/github.com/owncloud/ocis-pkg)

**This project is under heavy development, it's not in a working state yet!**

## Install

Just import the required packages within your ownCloud Infinite Scale extensions, nothing else should be required to do.

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.12.

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
