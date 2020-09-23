* * *

title: "Updating reva"
date: 2020-05-22T00:00:00+00:00
weight: 50
geekdocRepo: <https://github.com/owncloud/ocis-reva>
geekdocEditPath: edit/master/docs

## geekdocFilePath: updating.md

{{&lt; toc >}}

## Updating reva

1.  Run `go get github.com/cs3org/reva@master`
2.  Create a changelog entry containing changes that were done in [reva](https://github.com/cs3org/reva/commits/master)
3.  Create a Pull Request to ocis-reva master with those changes
4.  If test issues appear, you might need to [adjust the tests]\({{&lt; ref "testing.md" >}})
5.  After the PR is merged, consider doing a [release of ocis-reva]\({{&lt; ref "releasing.md" >}})
