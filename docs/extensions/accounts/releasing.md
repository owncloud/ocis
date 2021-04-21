---
title: "Releasing"
weight: 70
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/accounts
geekdocFilePath: releasing.md
---

{{< toc >}}

## Releasing

After adding changes to the Accounts package within oCIS and testing them locally, you want to update the compiled assets to the oCIS binary. 

To achieve this, you have run a Go command and add the results to your PR.

### Package Hierarchy

- [ocis](https://githug.com/owncloud/ocis)
    - [ocis-accounts](https://github.com/owncloud/ocis/tree/master/accounts)

#### Updating ocis-accounts

1. Make sure you are inside the [ocis repository](https://github.com/owncloud/ocis) and on your feature branch
2. Change into accounts' asset package folder via `cd accounts/pkg/assets`
3. Inside `accounts/pkg/assets`, run `go generate`. The output should look something like this: `accounts: embed.go - YYY/MM/DD ... to write [./embed.go] from config file ...`
4. Commit your changes, push them and [create a PR](https://github.com/owncloud/ocis/pulls)
