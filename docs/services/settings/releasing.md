---
title: "Releasing"
weight: 70
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/settings
geekdocFilePath: releasing.md
---

{{< toc >}}

## Requirements

You need a working installation of [the Go programming language](https://golang.org/).

## Releasing

The settings service doesn't have a dedicated release process. Simply commit your changes, make sure linting and unit tests pass locally and open a pull request.

### Package Hierarchy

- [ocis](https://github.com/owncloud/ocis)
    - [ocis-settings](https://github.com/owncloud/ocis/tree/master/settings)
