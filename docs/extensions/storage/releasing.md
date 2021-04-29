---
title: "Releasing"
date: 2020-05-22T00:00:00+00:00
weight: 60
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: releasing.md
---

{{< toc >}}

To release a new version of the storage submodule, you have to follow a few simple steps.

## Preparation

1. Before releasing, make sure that reva has been [updated to the desired version]({{< ref "updating" >}})

## Release
1. Check out master
{{< highlight txt >}}
git checkout master
git pull origin master
{{< / highlight >}}
2. Create a new tag (preferably signed) and replace the version number accordingly. Prefix the tag with the submodule `storage/v`.
{{< highlight txt >}}
git tag -s storage/vx.x.x -m "release vx.x.x"
git push origin storage/vx.x.x
{{< / highlight >}}
5. Wait for CI and check that the GitHub release was published.

Congratulations, you just released the storage submodule!
