# ownCloud Infinite Scale: Reva

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis-reva/status.svg)](https://cloud.drone.io/owncloud/ocis-reva)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6f1eaaa399294d959ef7b3b10deed41d)](https://www.codacy.com/manual/owncloud/ocis-reva?utm_source=github.com&utm_medium=referral&utm_content=owncloud/ocis-reva&utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis-reva?status.svg)](http://godoc.org/github.com/owncloud/ocis-reva)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis-reva)](http://goreportcard.com/report/github.com/owncloud/ocis-reva)
[![](https://images.microbadger.com/badges/image/owncloud/ocis-reva.svg)](http://microbadger.com/images/owncloud/ocis-reva "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/reva/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/extensions/ocis_reva/)

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html).

```console
git clone https://github.com/owncloud/ocis-reva.git
cd ocis-reva

make generate build

./bin/ocis-reva -h
```

To run a demo installation you can use the preconfigured defaults and start all necessary services:

    export REVA_USERS_DRIVER=demo

    bin/ocis-reva frontend & \
    bin/ocis-reva gateway & \
    bin/ocis-reva users & \
    bin/ocis-reva auth-basic & \
    bin/ocis-reva auth-bearer & \
    bin/ocis-reva sharing & \
    bin/ocis-reva storage-root & \
    bin/ocis-reva storage-home & \
    bin/ocis-reva storage-home-data & \
    bin/ocis-reva storage-oc & \
    bin/ocis-reva storage-oc-data

The root storage serves the available namespaces from disk using the local storage driver. In order to be able to navigate into the `/home` and `/oc` storage providers you have to create these directories:

    mkdir /var/tmp/reva/root/home
    mkdir /var/tmp/reva/root/oc

Note: the owncloud storage driver currently requires a redis server running on the local machine.

You should now be able to get a file listing of a users home using

    curl -X PROPFIND http://localhost:9140/remote.php/dav/files/ -v -u einstein:relativity

## Users

The default config uses the demo user backend, which contains three users:

    einstein:relativity
    marie:radioactivty
    richard:superfluidity

For details on the `json` and `ldap` backends see the [documentation](https://owncloud.github.io/extensions/ocis_reva/users/)

## Security

If you find a security issue please contact [security@owncloud.com](mailto:security@owncloud.com) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2019 ownCloud GmbH <https://owncloud.com>
```
