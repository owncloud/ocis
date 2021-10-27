---
title: "Releasing"
weight: 70
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/accounts
geekdocFilePath: releasing.md
---

{{< toc >}}

## Requirements

You need a working installation of [the Go programming language](https://golang.org/), [the Node runtime](https://nodejs.org/) and [the Yarn package manager](https://yarnpkg.com/) installed to build the assets for a working release.
## Releasing

The accounts service doesn't have a dedicated release process. Simply commit your changes, make sure linting and unit tests pass locally and open a pull request.

### Package Hierarchy

- [ocis](https://github.com/owncloud/ocis)
    - [ocis-accounts](https://github.com/owncloud/ocis/tree/master/accounts)
