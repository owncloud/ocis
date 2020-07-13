---
title: "Releasing"
weight: 40
geekdocRepo: https://github.com/owncloud/ocis-phoenix
geekdocEditPath: edit/master/docs
geekdocFilePath: releasing.md
---

{{< toc >}}

## Releasing

The next generation Web Frontend is shipped as an ocis Extension. The `ocis-phoenix` extension is also embedded in the single binary and part of the `ocis server` command.

To update this package within all the deliveries, we need to update the package in the following chain from the bottom to the top.

### Package Hierarchy

- [ocis](https://githug.com/owncloud/ocis)
    - [ocis-phoenix](https://github.com/owncloud/ocis-phoenix)
      - [ocis-pkg](https://github.com/owncloud/ocis-pkg)
      - [phoenix](https://github.com/owncloud/phoenix)

#### Prerequisites

Before updating the assets, make sure that [Phoenix](https://github.com/owncloud/phoenix) has been released first
and take note of its release tag name.

#### Updating ocis-phoenix

1. Create a branch `release-$version`. in <https://github.com/owncloud/ocis-phoenix>
2. Create a Folder in `changelog` for the release version and date `mkdir $major.$minor.$patchVersion_YYYY-MM-DD`.
3. Move all changelog items from the `changelog/unreleased/` folder to the `$major.$minor.$patchVersion_YYYY-MM-DD` folder.
4. Update the go module `ocis-pkg` to the latest version <https://blog.golang.org/using-go-modules> .
5. Update the phoenix asset by adjusting the value of `PHOENIX_ASSETS_VERSION` at the top of the Makefile and specify the tag name of the latest [Phoenix release](https://github.com/owncloud/phoenix/tags).
6. Run `make generate`.
7. Create a changelog item for the update in the `changelog/$major.$minor.$patchVersion_YYYY-MM-DD` folder.
8. Commit your changes.
9. After merging, wait for the CI to run on the merge commit.
10. Go to "Releases" in GH click "Draft a new Release".
11. Use `v$major.$minor.$patch` as a tag (the `v` prefix is important) and publish it.
12. The tag and the Release artifacts will be created automatically.

#### Next steps

Next steps is usually updating the ocis-phoenix version in the [ocis](https://githug.com/owncloud/ocis) repository.

