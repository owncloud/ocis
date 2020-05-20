---
title: "Releasing"
date: 2018-05-02T00:00:00+00:00
weight: 40
geekdocRepo: https://github.com/owncloud/ocis-thumbnails
geekdocEditPath: edit/master/docs
geekdocFilePath: releasing.md
---

{{< toc >}}

To release a new version of ocis-thumbnails, you have to follow a few simple steps.

## Preperation

1. Create a new branch e.g. `prepare-release-x.x.x` where `x.x.x` is the version you want to release.
2. Checkout the preparation branch.
3. Create a new changelog folder and move the unreleased snippets there.
{{< highlight txt >}}
mkdir changelog/x.x.x_yyyy-MM-dd/ # yyyy-MM-dd is the current date
mv changelog/unreleased/* changelog/x.x.x_yyyy-MM-dd/
{{< / highlight >}}
4. Commit and push the changes
{{< highlight txt >}}
git add --all
git commit -m "prepare release x.x.x"
git push
{{< / highlight >}}
5. Create a pull request to the master branch.

## Release
1. After the preparation branch has been merged update your local master.
{{< highlight txt >}}
git checkout master
git pull
{{< / highlight >}}
2. Create a new tag (preferably signed).
{{< highlight txt >}}
git tag -s vx.x.x -m "release vx.x.x"
git push --tags
{{< / highlight >}}


Congratulations, you just released ocis-thumbnails!
