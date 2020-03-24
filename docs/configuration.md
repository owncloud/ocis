---
title: "Configuration"
date: 2020-02-27T20:35:00+01:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: configuration.md
---

{{< toc >}}

## Configuration

oCIS Single Binary is not responsible for configuring extensions. Instead, each extension could either be configured by environment variables, cli flags or config files.

### Configuration using config files

Out of the box extensions will attempt to read configuration details from:

```console
/etc/ocis
$HOME/.ocis
./config
```

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-proxy reads `proxy.json | yaml | toml ...`*.

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below.

#### Global

OCIS_CONFIG_FILE
: Path to config file

OCIS_LOG_LEVEL
: Set logging level, defaults to `info`

OCIS_LOG_COLOR
: Enable colored logging, defaults to `true`

OCIS_LOG_PRETTY
: Enable pretty logging, defaults to `true`

#### Server

OCIS_TRACING_ENABLED
: Enable sending traces

OCIS_TRACING_TYPE
: Tracing backend type,

OCIS_TRACING_ENDPOINT
:Endpoint for the agent

OCIS_TRACING_COLLECTOR
: Endpoint for the collector

OCIS_TRACING_SERVICE
: Service name for tracing"

OCIS_DEBUG_ADDR
: Address to bind debug server, defaults to `0.0.0.0:9010`

OCIS_DEBUG_TOKEN
: Token to grant metrics access, empty default value

OCIS_DEBUG_PPROF
: Enable pprof debugging, defaults to `false`

OCIS_DEBUG_ZPAGES
: Enable zpages debugging, defaults to `false`

OCIS_HTTP_ADDR
: Address to bind http server, defaults to `0.0.0.0:9000`

OCIS_HTTP_ROOT
: Root path for http endpoint, defaults to `/`

OCIS_GRPC_ADDR
: Address to bind grpc server, defaults to `0.0.0.0:9001`

OCIS_SERVICES_ENABLED
: List of enabled services, defaults to `phoenix,konnectd,graph,ocs,webdav,hello`

#### Health

OCIS_DEBUG_ADDR
: Address to debug endpoint, defaults to `0.0.0.0:9010`

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below.

#### Global

--config-file
: Path to config file

--log-level
: Set logging level, defaults to `info`

--log-color
: Enable colored logging, defaults to `true`

--log-pretty
: Enable pretty logging, defaults to `true`

#### Server

--tracing-enabled
: Enable sending traces

--tracing-type
: Tracing backend type,

--tracing-endpoint
:Endpoint for the agent

--tracing-collector
: Endpoint for the collector

--tracing-service
: Service name for tracing"

--debug-addr
: Address to bind debug server, defaults to `0.0.0.0:9010`

--debug-token
: Token to grant metrics access, empty default value

--debug-pprof
: Enable pprof debugging, defaults to `false`

--debug-zpages
: Enable zpages debugging, defaults to `false`

--http-addr
: Address to bind http server, defaults to `0.0.0.0:9000`

--http-root
: Root path for http endpoint, defaults to `/`

--grpc-addr
: Address to bind grpc server, defaults to `0.0.0.0:9001`

--services-enabled
: List of enabled services, defaults to `hello,phoenix,graph,graph-explorer,ocs,webdav,reva-frontend,reva-gateway,reva-users,reva-auth-basic,reva-auth-bearer,reva-sharing,reva-storage-root,reva-storage-home,reva-storage-home-data,reva-storage-oc,reva-storage-oc-data,devldap`

#### Health

--debug-addr
: Address to debug endpoint, defaults to `0.0.0.0:9010`

### Configuration file

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

## Usage

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formated help output if you execute the binary via `ocis --help`.

### Server

The server command is used to start the http and debug server on two addresses within a single process. The http server is serving the general webservice while the debug server is used for health check, readiness check and to server the metrics mentioned below. For further help please execute:

{{< highlight txt >}}
ocis server --help
{{< / highlight >}}

### Health

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exist code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
ocis health --help
{{< / highlight >}}

## Metrics

This service provides some [Prometheus](https://prometheus.io/) metrics through the debug endpoint, you can optionally secure the metrics endpoint by some random token, which got to be configured through one of the flag `--debug-token` or the environment variable `OCIS_DEBUG_TOKEN` mentioned above. By default the metrics endpoint is bound to `http://0.0.0.0:8001/metrics`.

go_gc_duration_seconds
: A summary of the GC invocation durations

go_gc_duration_seconds_sum
: A summary of the GC invocation durations

go_gc_duration_seconds_count
: A summary of the GC invocation durations

go_goroutines
: Number of goroutines that currently exist

go_info
: Information about the Go environment

go_memstats_alloc_bytes
: Number of bytes allocated and still in use

go_memstats_alloc_bytes_total
: Total number of bytes allocated, even if freed

go_memstats_buck_hash_sys_bytes
: Number of bytes used by the profiling bucket hash table

go_memstats_frees_total
: Total number of frees

go_memstats_gc_cpu_fraction
: The fraction of this program's available CPU time used by the GC since the program started

go_memstats_gc_sys_bytes
: Number of bytes used for garbage collection system metadata

go_memstats_heap_alloc_bytes
: Number of heap bytes allocated and still in use

go_memstats_heap_idle_bytes
: Number of heap bytes waiting to be used

go_memstats_heap_inuse_bytes
: Number of heap bytes that are in use

go_memstats_heap_objects
: Number of allocated objects

go_memstats_heap_released_bytes
: Number of heap bytes released to OS

go_memstats_heap_sys_bytes
: Number of heap bytes obtained from system

go_memstats_last_gc_time_seconds
: Number of seconds since 1970 of last garbage collection

go_memstats_lookups_total
: Total number of pointer lookups

go_memstats_mallocs_total
: Total number of mallocs

go_memstats_mcache_inuse_bytes
: Number of bytes in use by mcache structures

go_memstats_mcache_sys_bytes
: Number of bytes used for mcache structures obtained from system

go_memstats_mspan_inuse_bytes
: Number of bytes in use by mspan structures

go_memstats_mspan_sys_bytes
: Number of bytes used for mspan structures obtained from system

go_memstats_next_gc_bytes
: Number of heap bytes when next garbage collection will take place

go_memstats_other_sys_bytes
: Number of bytes used for other system allocations

go_memstats_stack_inuse_bytes
: Number of bytes in use by the stack allocator

go_memstats_stack_sys_bytes
: Number of bytes obtained from system for stack allocator

go_memstats_sys_bytes
: Number of bytes obtained from system

go_threads
: Number of OS threads created

promhttp_metric_handler_requests_in_flight
: Current number of scrapes being served

promhttp_metric_handler_requests_total
: Total number of scrapes by HTTP status code
