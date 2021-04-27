---
title: "Releasing"
weight: 70
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/accounts
geekdocFilePath: releasing.md
---

{{< toc >}}

## Requirements

You need a working installation of [the Go programming language](https://golang.org/) installed to build the assets for a working release.

## Releasing

After adding changes to the accounts package within oCIS and testing them locally, you want to update the compiled assets to the oCIS binary. 

To achieve this, you have to run a Go command and add the results to your PR. The preferred way to do this is to run `make generate` in the root 
of the repository and then commit the resulting changes to your branch/PR. See below for a way to _only_ build the accounts extension assets.

### Package Hierarchy

- [ocis](https://githug.com/owncloud/ocis)
    - [ocis-accounts](https://github.com/owncloud/ocis/tree/master/accounts)

#### Updating ocis-accounts

1. Make sure you are inside the [ocis repository](https://github.com/owncloud/ocis) and on your feature branch
2. Change into accounts' asset package folder via `cd accounts/pkg/assets`
3. Inside `accounts/pkg/assets`, run `go generate`. The output should look something like this: `accounts: embed.go - YYY/MM/DD ... to write [./embed.go] from config file ...`
4. Commit your changes, push them and [create a PR](https://github.com/owncloud/ocis/pulls)
