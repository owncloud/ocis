---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/thumbnails
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Installation

So far we are offering two different variants for the installation. You can choose between [Docker](https://www.docker.com/) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

### Docker

TBD

### Binaries

TBD

## Configuration

We provide overall three different variants of configuration. The variant based on environment variables and commandline flags are split up into global values and command-specific values.

### Environment variables

If you prefer to configure the service with environment variables you can see the available variables below.

#### Global

THUMBNAILS_CONFIG_FILE
: Path to config file, empty default value

THUMBNAILS_LOG_LEVEL
: Set logging level, defaults to `info`

THUMBNAILS_LOG_COLOR
: Enable colored logging, defaults to `true`

THUMBNAILS_LOG_PRETTY
: Enable pretty logging, defaults to `true`

#### Server

THUMBNAILS_TRACING_ENABLED
: Enable sending traces, defaults to `false`

THUMBNAILS_TRACING_TYPE
: Tracing backend type, defaults to `jaeger`

THUMBNAILS_TRACING_ENDPOINT
: Endpoint for the agent, empty default value

THUMBNAILS_TRACING_COLLECTOR
: Endpoint for the collector, empty default value

THUMBNAILS_TRACING_SERVICE
: Service name for tracing, defaults to `ocis-thumbnails`

THUMBNAILS_DEBUG_ADDR
: Address to bind debug server, defaults to `0.0.0.0:9114`

THUMBNAILS_DEBUG_TOKEN
: Token to grant metrics access, empty default value

THUMBNAILS_DEBUG_PPROF
: Enable pprof debugging, defaults to `false`

THUMBNAILS_DEBUG_ZPAGES
: Enable zpages debugging, defaults to `false`

THUMBNAILS_GRPC_NAME
: Name of the service, defaults to `thumbnails`

THUMBNAILS_GRPC_ADDR
: Address to bind grpc server, defaults to `0.0.0.0:9185`

THUMBNAILS_GRPC_NAMESPACE
: Set the base namespace for the grpc namespace", defaults to `com.owncloud.api`

THUMBNAILS_FILESYSTEMSTORAGE_ROOT
: Root path of the filesystem storage directory, defaults to `<os tempdir>/ocis-thumbnails/`

THUMBNAILS_WEBDAVSOURCE_BASEURL
: Base url for a webdav api, defaults to `https://localhost:9200/remote.php/webdav/`

THUMBNAILS_RESOLUTIONS
: List of resolutions supported by the service, defaults to `["16x16", "32x32", "64x64", "128x128"]

#### Health

THUMBNAILS_DEBUG_ADDR
: Address to debug endpoint, defaults to `0.0.0.0:9189`

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below.

#### Global

--config-file
: Path to config file, empty default value

--log-level
: Set logging level, defaults to `info`

--log-color
: Enable colored logging, defaults to `true`

--log-pretty
: Enable pretty logging, defaults to `true`

#### Server

--tracing-enabled
: Enable sending traces, defaults to `false`

--tracing-type
: Tracing backend type, defaults to `jaeger`

--tracing-endpoint
: Endpoint for the agent, empty default value

--tracing-collector
: Endpoint for the collector, empty default value

--tracing-service
: Service name for tracing, defaults to `thumbnails`

--debug-addr
: Address to bind debug server, defaults to `0.0.0.0:9189`

--debug-token
: Token to grant metrics access, empty default value

--debug-pprof
: Enable pprof debugging, defaults to `false`

--debug-zpages
: Enable zpages debugging, defaults to `false`

--grpc-name
: Name of the service, defaults to `thumbnails`

--grpc-addr
: Address to bind grpc server, defaults to `0.0.0.0:9185`

--grpc-namespace
: Set the base namespace for the grpc namespace", defaults to `com.owncloud.api`

--filesystemstorage-root
: Root path of the filesystem storage directory, defaults to `<os tempdir>/ocis-thumbnails/`

--webdavsource-baseurl
: Base url for a webdav api, defaults to `https://localhost:9200/remote.php/webdav/`

--thumbnail-resolution
: List of resolutions supported by the service, defaults to `["16x16", "32x32", "64x64", "128x128"]

#### Health

--debug-addr
: Address to debug endpoint, defaults to `0.0.0.0:9189`

### Configuration file

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/thumbnails/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/thumbnails.yml`, `${HOME}/.ocis/thumbnails.yml` or `$(pwd)/config/thumbnails.yml`.

## Usage

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formated help output if you execute the binary via `ocis-thumbnails --help`.

### Server

The server command is used to start the http and debug server on two addresses within a single process. The http server is serving the general webservice while the debug server is used for health check, readiness check and to server the metrics mentioned below. For further help please execute:

{{< highlight txt >}}
{{ Name }} server --help
{{< / highlight >}}

### Health

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exist code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
{{ Name }} health --help
{{< / highlight >}}

### Version
The version command lists the versions of all running instances. For further help please execute:

{{< highlight txt >}}
{{ Name }} version --help
{{< /highlight >}}

## Metrics

This service provides some [Prometheus](https://prometheus.io/) metrics through the debug endpoint, you can optionally secure the metrics endpoint by some random token, which got to be configured through one of the flag `--debug-token` or the environment variable `THUMBNAILS_DEBUG_TOKEN` mentioned above. By default the metrics endpoint is bound to `http://0.0.0.0:9114/metrics`.

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
