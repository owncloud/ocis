---
title: "Releasing"
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/web
geekdocFilePath: releasing.md
---

{{< toc >}}

## Releasing

The next generation Web Frontend is shipped as an oCIS Extension. The `ocis-web` extension is also embedded in the single binary and part of the `ocis server` command.

To update this package within all the deliveries, we need to update the package in the following chain from the bottom to the top.

### Package Hierarchy

- [ocis](https://github.com/owncloud/ocis)
    - [ocis-web](https://github.com/owncloud/ocis/tree/master/web)
      - [ocis-pkg](https://github.com/owncloud/ocis/tree/master/ocis-pkg)
      - [ownCloud Web](https://github.com/owncloud/web)

#### Prerequisites

Before updating the assets, make sure that [ownCloud Web](https://github.com/owncloud/web) has been released first
and take note of its release tag name.

#### Updating ocis-web

1. Create a branch `update-web-$version` in the [ocis repository](https://github.com/owncloud/ocis)
2. Change into web package folder via `cd web`
3. Inside `web/`, update the `Makefile` so that the WEB_ASSETS_VERSION variable references the currently released version of https://github.com/owncloud/web
4. Move to the changelog (`cd ../changelog/`) and add a changelog file to the `unreleased/` folder (You can copy an old web release changelog item as a template)
5. Move to the repo root (`cd ..`)and update the WEB_COMMITID in the `/.drone.env` file to the commit id from the released version (unless the existing commit id is already newer)
6. **Optional:** Test the changes locally by running `cd ocis && go run cmd/ocis/main.go server`, visiting [https://localhost:9200](https://localhost:9200) and confirming everything renders correctly
7. Commit your changes, push them and [create a PR](https://github.com/owncloud/ocis/pulls)
