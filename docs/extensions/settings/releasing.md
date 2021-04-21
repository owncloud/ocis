---
title: "Releasing"
weight: 70
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/settings
geekdocFilePath: releasing.md
---

{{< toc >}}

## Releasing

After adding changes to the Settings package within oCIS and testing them locally, you want to update the compiled assets to the oCIS binary. 

To achieve this, you have run a Go command and add the results to your PR.

### Package Hierarchy

- [ocis](https://githug.com/owncloud/ocis)
    - [ocis-settings](https://github.com/owncloud/ocis/tree/master/settings)

#### Updating ocis-settings

1. Make sure you are inside the [ocis repository](https://github.com/owncloud/ocis) and on your feature branch
2. Change into settings' asset package folder via `cd settings/pkg/assets`
3. Inside `settings/pkg/assets`, run `go generate`. The output should look something like this: `settings: embed.go - YYY/MM/DD ... to write [./embed.go] from config file ...`
4. Commit your changes, push them and [create a PR](https://github.com/owncloud/ocis/pulls)
