---
title: "Debugging"
date: 2020-03-19T08:21:00+01:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: debugging.md
---

{{< toc >}}

## Debugging

As a single binary for easy deployment running `ocis server` just forks itself to start all the services, which makes debugging those processes a little harder.

Ultimately, we want to be able to stop a single service using eg. `ocis kill web` so that you can start the service you want to debug in debug mode. We need to [change the way we fork processes](https://github.com/owncloud/ocis/issues/77) though, otherwise the runtime will automatically restart a service if killed.

### Start ocis

For debugging there are two workflows that work well, depending on your preferences.

#### Use the debug binary and attach to the process as needed

Run the debug binary with `OCIS_LOG_LEVEL=debug bin/ocis-debug server` and then find the service you want to debug using:

```console
# ps ax | grep ocis
12837 pts/1    Sl+    0:00 bin/ocis-debug server
12845 pts/1    Sl     0:00 bin/ocis-debug graph
12847 pts/1    Sl     0:00 bin/ocis-debug reva-auth-bearer
12848 pts/1    Sl     0:00 bin/ocis-debug graph-explorer
12849 pts/1    Sl     0:00 bin/ocis-debug ocs
12850 pts/1    Sl     0:00 bin/ocis-debug reva-storage-oc-data
12863 pts/1    Sl     0:00 bin/ocis-debug webdav
12874 pts/1    Sl     0:00 bin/ocis-debug reva-frontend
12897 pts/1    Sl     0:00 bin/ocis-debug reva-sharing
12905 pts/1    Sl     0:00 bin/ocis-debug reva-gateway
12912 pts/1    Sl     0:00 bin/ocis-debug reva-storage-home
12920 pts/1    Sl     0:00 bin/ocis-debug reva-users
12929 pts/1    Sl     0:00 bin/ocis-debug glauth
12940 pts/1    Sl     0:00 bin/ocis-debug reva-storage-home-data
12948 pts/1    Sl     0:00 bin/ocis-debug konnectd
12952 pts/1    Sl     0:00 bin/ocis-debug proxy
12961 pts/1    Sl     0:00 bin/ocis-debug thumbnails
12971 pts/1    Sl     0:00 bin/ocis-debug reva-storage-oc
12981 pts/1    Sl     0:00 bin/ocis-debug web
12993 pts/1    Sl     0:00 bin/ocis-debug api
12998 pts/1    Sl     0:00 bin/ocis-debug registry
13004 pts/1    Sl     0:00 bin/ocis-debug web
13015 pts/1    Sl     0:00 bin/ocis-debug reva-auth-basic
```

Then you can set a breakpoint in the service you need and attach to the process via processid. To debug the `reva-sharing` service the VS Code `launch.json` would look like this:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "ocis attach",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": 12897
    }
  ]
}
```

#### Start all services independently to replace one of them with a debug process

1. You can use this `./ocis.sh` script to start all services independently, so they don't get restarted by the runtime when you kill them:

```bash
#/bin/sh
LOG_LEVEL="debug"

bin/ocis --log-level=$LOG_LEVEL micro &

bin/ocis --log-level=$LOG_LEVEL glauth &
bin/ocis --log-level=$LOG_LEVEL graph-explorer &
bin/ocis --log-level=$LOG_LEVEL graph &
#bin/ocis --log-level=$LOG_LEVEL hello &
bin/ocis --log-level=$LOG_LEVEL konnectd &
#bin/ocis --log-level=$LOG_LEVEL ocs &
bin/ocis --log-level=$LOG_LEVEL web &
bin/ocis --log-level=$LOG_LEVEL reva-auth-basic &
bin/ocis --log-level=$LOG_LEVEL reva-auth-bearer &
bin/ocis --log-level=$LOG_LEVEL reva-frontend &
bin/ocis --log-level=$LOG_LEVEL reva-gateway &
bin/ocis --log-level=$LOG_LEVEL reva-sharing &
bin/ocis --log-level=$LOG_LEVEL reva-storage-home &
bin/ocis --log-level=$LOG_LEVEL reva-storage-home-data &
bin/ocis --log-level=$LOG_LEVEL reva-storage-oc &
bin/ocis --log-level=$LOG_LEVEL reva-storage-oc-data &
bin/ocis --log-level=$LOG_LEVEL reva-storage-root &
bin/ocis --log-level=$LOG_LEVEL reva-users &
#bin/ocis --log-level=$LOG_LEVEL webdav

bin/ocis --log-level=$LOG_LEVEL proxy &
```

2. Get the list of running processes:

```console
# ps ax | grep ocis
12837 pts/1    Sl+    0:00 bin/ocis-debug server
12845 pts/1    Sl     0:00 bin/ocis-debug graph
12847 pts/1    Sl     0:00 bin/ocis-debug reva-auth-bearer
12848 pts/1    Sl     0:00 bin/ocis-debug graph-explorer
12849 pts/1    Sl     0:00 bin/ocis-debug ocs
12850 pts/1    Sl     0:00 bin/ocis-debug reva-storage-oc-data
12863 pts/1    Sl     0:00 bin/ocis-debug webdav
12874 pts/1    Sl     0:00 bin/ocis-debug reva-frontend
12897 pts/1    Sl     0:00 bin/ocis-debug reva-sharing
12905 pts/1    Sl     0:00 bin/ocis-debug reva-gateway
12912 pts/1    Sl     0:00 bin/ocis-debug reva-storage-home
12920 pts/1    Sl     0:00 bin/ocis-debug reva-users
12929 pts/1    Sl     0:00 bin/ocis-debug glauth
12940 pts/1    Sl     0:00 bin/ocis-debug reva-storage-home-data
12948 pts/1    Sl     0:00 bin/ocis-debug konnectd
12952 pts/1    Sl     0:00 bin/ocis-debug proxy
12961 pts/1    Sl     0:00 bin/ocis-debug thumbnails
12971 pts/1    Sl     0:00 bin/ocis-debug reva-storage-oc
12981 pts/1    Sl     0:00 bin/ocis-debug web
12993 pts/1    Sl     0:00 bin/ocis-debug api
12998 pts/1    Sl     0:00 bin/ocis-debug registry
13004 pts/1    Sl     0:00 bin/ocis-debug web
13015 pts/1    Sl     0:00 bin/ocis-debug reva-auth-basic
```

3. Kill the service you want to start in debug mode:

```console
# kill 17628
```

4. Start the service you are interested in in debug mode. When using make to build the binary there is already a `bin/ocis-debug` binary for you. When running an IDE tell it which service to start by providing the corresponding sub command, eg. `bin\ocis-debug reva-frontend`.

### Gather error messages

We recommend you collect all related information in a single file or in a github issue. Let us start with an error that pops up in the Web UI:

> Error while sharing.
> error sending a grpc stat request

This popped up when I tried to add `marie` as a collaborator in ownCloud Web. That triggers a request to the server which I copied as curl. We can strip a lot of headers and the gist of it is:

```console
# curl 'https://localhost:9200/ocs/v1.php/apps/files_sharing/api/v1/shares' -d 'shareType=0&shareWith=marie&path=%2FNeuer+Ordner&permissions=1' -u einstein:relativity -k -v | xmllint -format -
[... headers ...]
<?xml version="1.0" encoding="UTF-8"?>
<ocs>
  <meta>
    <status>error</status>
    <statuscode>998</statuscode>
    <message>error sending a grpc stat request</message>
  </meta>
</ocs>
```

{{< hint info >}}
The username and password only work when basic auth is available. Otherwise you have to obtain a bearer token, eg. by grabbing it from the browser.
{{< /hint >}}
{{< hint danger >}}
TODO add ocis cli tool to obtain a bearer token.
{{< /hint >}}

We also have a few interesting log entries:

```
0:43PM INF home/jfd/go/pkg/mod/github.com/cs3org/reva@v0.0.2-0.20200318111623-a2f97d4aa741/internal/grpc/interceptors/log/log.go:69 > unary code=OK end="18/Mar/2020:22:43:40 +0100" from=tcp://[::1]:44078 pid=17836 pkg=rgrpc start="18/Mar/2020:22:43:40 +0100" time_ns=95841 traceid=b4eb9a9f45921f7d3632523ca32a42b0 uri=/cs3.storage.registry.v1beta1.RegistryAPI/GetStorageProvider user-agent=grpc-go/1.26.0
10:43PM ERR home/jfd/go/pkg/mod/github.com/cs3org/reva@v0.0.2-0.20200318111623-a2f97d4aa741/internal/grpc/interceptors/log/log.go:69 > unary code=Unknown end="18/Mar/2020:22:43:40 +0100" from=tcp://[::1]:43910 pid=17836 pkg=rgrpc start="18/Mar/2020:22:43:40 +0100" time_ns=586115 traceid=b4eb9a9f45921f7d3632523ca32a42b0 uri=/cs3.gateway.v1beta1.GatewayAPI/Stat user-agent=grpc-go/1.26.0
10:43PM ERR home/jfd/go/pkg/mod/github.com/cs3org/reva@v0.0.2-0.20200318111623-a2f97d4aa741/internal/http/services/owncloud/ocs/reqres.go:94 > error sending a grpc stat request error="rpc error: code = Unknown desc = gateway: error calling Stat: rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing dial tcp [::1]:9152: connect: connection refused\"" pid=17832 pkg=rhttp traceid=b4eb9a9f45921f7d3632523ca32a42b0
```

{{< hint danger >}}
TODO return the trace id in the response so we can correlate easier. For reva tracked in https://github.com/cs3org/reva/issues/587
{{< /hint >}}

The last line gives us a hint where the log message originated: `.../github.com/cs3org/reva@v0.0.2-0.20200318111623-a2f97d4aa741/internal/http/services/owncloud/ocs/reqres.go:94`. Which looks like this:

```go
89: // WriteOCSResponse handles writing ocs responses in json and xml
90: func WriteOCSResponse(w http.ResponseWriter, r *http.Request, res *Response, err error) {
91: 	var encoded []byte
92:
93: 	if err != nil {
94: 		appctx.GetLogger(r.Context()).Error().Err(err).Msg(res.OCS.Meta.Message)
95:     }
```

Ok, so this seems to be a convenience method that is called from multiple places an also handles errors. Unfortunately, this hides the actual source of the error. We could set a breakpoint in line 94 and reproduce the problem, which can be a lot harder than just clicking the share button or sending a curl request again. So let us see what else the log tells us.

The previous line tells us that a Stat request failed: `uri=/cs3.gateway.v1beta1.GatewayAPI/Stat`. This time the line is written by the grpc log interceptor. What else is there?

The first line tells us that looking up the responsible storage provider seems to have succeeded: `uri=/cs3.storage.registry.v1beta1.RegistryAPI/GetStorageProvider`.

At this point it your familiarity with the codebase starts to become a factor. If you are new you should probably go back to setting a break point on the log line and check the stack trace.

Debug wherever the call trace leads you to ... good luck!

### Managing dependencies and testing changes

You can either run and manage the services independently, or you can update the `go.mod` file and replace dependencies with your local version.

To debug the reva frontend we need to add two replacements:

```
// use the local ocis-reva repo
replace github.com/owncloud/ocis-reva => ../ocis-reva
// also use the local reva repo
replace github.com/cs3org/reva => ../reva
```

{{< hint info >}}
The username and password only work when basic auth is available. Otherwise you have to obtain a bearer token, eg. by grabbing it from the browser.
{{< /hint >}}

Rebuild ocis to make sure the dependency is used. It should be sufficient to just restart the service you want to debug.
