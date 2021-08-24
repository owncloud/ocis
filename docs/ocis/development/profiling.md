---
title: "Profiling"
date: 2021-08-24T12:32:20+01:00
weight: 56
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: profiling.md
---

{{< toc >}}

# 0. Prerequisites

- Go development kit of a [supported version](https://golang.org/doc/devel/release.html#policy).
  Follow [these instructions](http://golang.org/doc/code.html) to install the
  go tool and set up GOPATH.

- Graphviz: http://www.graphviz.org/
  Optional, used to generate graphic visualizations of profiles

The only way to enable the profiler currently is to explicitly select which areas to collect samples for. In order to do this, the following steps have to be followed.

## 1. Clone Reva

`git clone github.com/cs3org/reva`

## 2. Patch reva with the area that you want sampled.

For the purposes of this docs let's use the WebDAV `PROPFIND` path.

```diff
diff --git a/internal/http/services/owncloud/ocdav/propfind.go b/internal/http/services/owncloud/ocdav/propfind.go
index 0e9c99be..f271572f 100644
--- a/internal/http/services/owncloud/ocdav/propfind.go
+++ b/internal/http/services/owncloud/ocdav/propfind.go
@@ -32,6 +32,8 @@ import (
        "strings"
        "time"

+       _ "net/http/pprof"
+
        userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
        rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
        link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
@@ -311,6 +313,12 @@ func requiresExplicitFetching(n *xml.Name) bool {
        return true
 }

+func init() {
+       go func() {
+               http.ListenAndServe(":1234", nil)
+       }()
+}
+
 // from https://github.com/golang/net/blob/e514e69ffb8bc3c76a71ae40de0118d794855992/webdav/xml.go#L178-L205
 func readPropfind(r io.Reader) (pf propfindXML, status int, err error) {
        c := countingReader{r: r}
```

The previous patch will:

1. import `net/http/pprof`, which will register debug handlers in `DefaultServeMux`.
2. define a `init()` function that starts an http server with the previously registered handlers.

With everything running one should have access to http://localhost:1234/debug/pprof/

## 3. Replace reva in oCIS go.mod with local version and build a new binary

```diff
diff --git a/go.mod b/go.mod
index 131d14d7b..9668c38e4 100644
--- a/go.mod
+++ b/go.mod
@@ -78,6 +78,7 @@ require (

 replace (
        github.com/crewjam/saml => github.com/crewjam/saml v0.4.5
+       github.com/cs3org/reva => path/to/your/reva
        go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
        go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
 )
```

Make sure to replace `github.com/cs3org/reva => path/to/your/reva` with the correct location of your reva.

## 4. Build a new ocis binary

From owncloud/ocis root:

```console
$ cd ocis
$ make clean build
```

## 5. Start oCIS server

From owncloud/ocis root:

```console
$ OCIS_LOG_PRETTY=true OCIS_LOG_COLOR=true ocis/bin/ocis server
```

## 6. Run `pprof`

### Install pprof

If `pprof` is not installed make sure to get it; one way of installing it is using the Go tools:

```console
$ go get -u github.com/google/pprof
```

### Collecting samples

Collect 30 seconds of samples:

```console
$ pprof -web http://:1234/debug/pprof/profile\?seconds\=30
```

Once the collection is done a browser tab will open with the result `svg`, looking similar to this:

![img](https://i.imgur.com/vo0EbcX.jpg)

For references on how to interpret this graph, [continue reading here](https://github.com/google/pprof/blob/master/doc/README.md#interpreting-the-callgraph).

## References

- https://medium.com/swlh/go-profile-your-code-like-a-master-1505be38fdba
- https://dave.cheney.net/2013/07/07/introducing-profile-super-simple-profiling-for-go-programs
