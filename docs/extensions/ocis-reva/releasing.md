---
title: "Releasing"
date: 2020-05-22T00:00:00+00:00
weight: 60
geekdocRepo: https://github.com/owncloud/ocis-reva
geekdocEditPath: edit/master/docs
geekdocFilePath: releasing.md
---

{{< toc >}}

To release a new version of ocis-reva, you have to follow a few simple steps.

## Preparation

1. Before releasing, make sure that reva has been [updated to the desired version]({{< ref "updating.md" >}})
2. Create a new branch e.g. `release-x.x.x` where `x.x.x` is the version you want to release.
3. Checkout the preparation branch.
4. Create a new changelog folder and move the unreleased snippets there.
{{< highlight txt >}}
mkdir changelog/x.x.x_yyyy-MM-dd/ # yyyy-MM-dd is the current date
mv changelog/unreleased/* changelog/x.x.x_yyyy-MM-dd/
{{< / highlight >}}
5. Commit and push the changes
{{< highlight txt >}}
git add --all
git commit -m "prepare release x.x.x"
git push origin release-x.x.x
{{< / highlight >}}
6. Create a pull request to the master branch.

## Release
1. After the preparation branch has been merged update your local master.
2. [Wait for CI](https://cloud.drone.io/owncloud/ocis-reva) to generate a commit for the changelog update
3. Check out master (or make sure to check out the generated changelog commit in case of subsequent merges)
{{< highlight txt >}}
git checkout master
git pull origin master
{{< / highlight >}}
4. Create a new tag (preferably signed) and replace the version number accordingly.
{{< highlight txt >}}
git tag -s vx.x.x -m "release vx.x.x"
git push origin vx.x.x
{{< / highlight >}}
5. Wait for CI and check that the GitHub release was published.

Congratulations, you just released ocis-reva!
