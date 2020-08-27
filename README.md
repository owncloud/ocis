# ownCloud Infinite Scale

[![Build Status](https://cloud.drone.io/api/badges/owncloud/ocis/status.svg)](https://cloud.drone.io/owncloud/ocis)
[![Gitter chat](https://badges.gitter.im/cs3org/reva.svg)](https://gitter.im/cs3org/reva)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/dc97ddfa167641d8b107e9b618823c71)](https://www.codacy.com/app/owncloud/ocis?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=owncloud/ocis&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/owncloud/ocis?status.svg)](http://godoc.org/github.com/owncloud/ocis)
[![Go Report](http://goreportcard.com/badge/github.com/owncloud/ocis)](http://goreportcard.com/report/github.com/owncloud/ocis)
[![](https://images.microbadger.com/badges/image/owncloud/ocis.svg)](http://microbadger.com/images/owncloud/ocis "Get your own image badge on microbadger.com")

**This project is under heavy development, it's not in a working state yet!**

## Install

You can download prebuilt binaries from the GitHub releases or from our [download mirrors](http://download.owncloud.com/ocis/ocis/). For instructions how to install this on your platform you should take a look at our [documentation](https://owncloud.github.io/ocis/)

## Development

Trigger CI.

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.13.

```console
git clone https://github.com/owncloud/ocis.git
cd ocis

make generate build

./bin/ocis -h
```

## Prerequisites

### Redis server

You will need to start a redis server as a cache. The ownCloud storage driver currently will try to connect to the default port.
A quick way to start one for testing is using this docker instance: `docker run -e REDIS_DATABASES=1 -p 6379:6379 -d webhippie/redis:latest`

### Root storage

To prepare the root storage you should fill it with two folders. They are necessary for resolving the home and ownCloud storages. This is subject to change.

```console
mkdir -p /var/tmp/reva/root/{home,oc}
```

## Quickstart

-   Make sure that the binary was built with the above steps.

-   Now start all services with the following command

    ```console
    ./bin/ocis server
    ```

-   Open [https://localhost:9200](https://localhost:9200)

-   Accept the self-signed certificate (it is regenerated every time the server starts)

-   Login using one of the demo accounts:

    ```console
    einstein:relativity
    marie:radioactivity
    richard:superfluidity
    ```

## Running single extensions

The list of available extensions can be found in the "Extensions" section when running `./bin/ocis`.

For example to run the "phoenix" extension:
```console
./bin/ocis --log-level debug phoenix
```

⚠ do not use the **run** subcommand for running extensions

## Security

If you find a security issue please contact security@owncloud.com first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2020 ownCloud GmbH <https://owncloud.com>
```
