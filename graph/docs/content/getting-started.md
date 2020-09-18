---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
anchor: "getting-started"
weight: 20
---

### Installation

So far we are offering two different variants for the installation. You can choose between [Docker](https://www.docker.com/) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

#### Docker

TBD

#### Binaries

TBD

### Configuration

We provide overall three different variants of configuration. The variant based on environment variables and commandline flags are split up into global values and command-specific values.

#### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below.

##### Global

GRAPH_CONFIG_FILE
: Path to config file, empty default value

GRAPH_LOG_LEVEL
: Set logging level, defaults to `info`

GRAPH_LOG_COLOR
: Enable colored logging, defaults to `true`

GRAPH_LOG_PRETTY
: Enable pretty logging, defaults to `true`

##### Server

GRAPH_TRACING_ENABLED
: Enable sending traces, defaults to `false`

GRAPH_TRACING_TYPE
: Tracing backend type, defaults to `jaeger`

GRAPH_TRACING_ENDPOINT
: Endpoint for the agent, empty default value

GRAPH_TRACING_COLLECTOR
: Endpoint for the collector, empty default value

GRAPH_TRACING_SERVICE
: Service name for tracing, defaults to `graph`

GRAPH_DEBUG_ADDR
: Address to bind debug server, defaults to `0.0.0.0:9124`

GRAPH_DEBUG_TOKEN
: Token to grant metrics access, empty default value

GRAPH_DEBUG_PPROF
: Enable pprof debugging, defaults to `false`

GRAPH_DEBUG_ZPAGES
: Enable zpages debugging, defaults to `false`

GRAPH_HTTP_ADDR
: Address to bind http server, defaults to `0.0.0.0:9120`

GRAPH_HTTP_ROOT
: Root path of http server, defaults to `/`

##### Health

GRAPH_DEBUG_ADDR
: Address to debug endpoint, defaults to `0.0.0.0:9124`

#### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below.

##### Global

--config-file
: Path to config file, empty default value

--log-level
: Set logging level, defaults to `info`

--log-color
: Enable colored logging, defaults to `true`

--log-pretty
: Enable pretty logging, defaults to `true`

##### Server

--tracing-enabled
: Enable sending traces, defaults to `false`

--tracing-type
: Tracing backend type, defaults to `jaeger`

--tracing-endpoint
: Endpoint for the agent, empty default value

--tracing-collector
: Endpoint for the collector, empty default value

--tracing-service
: Service name for tracing, defaults to `graph`

--debug-addr
: Address to bind debug server, defaults to `0.0.0.0:9124`

--debug-token
: Token to grant metrics access, empty default value

--debug-pprof
: Enable pprof debugging, defaults to `false`

--debug-zpages
: Enable zpages debugging, defaults to `false`

--http-addr
: Address to bind http server, defaults to `0.0.0.0:9120`

--http-namespace
: Namespace for internal services communication, defaults to `com.owncloud.web`

--http-root
: Root path of http server, defaults to `/`

##### Health

--debug-addr
: Address to debug endpoint, defaults to `0.0.0.0:9124`

#### Configuration file

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis-graph/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/graph.yml`, `${HOME}/.ocis/graph.yml` or `$(pwd)/config/graph.yml`.

### Usage

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formated help output if you execute the binary via `ocis-graph --help`.

#### Server

The server command is used to start the http and debug server on two addresses within a single process. The http server is serving the general webservice while the debug server is used for health check, readiness check and to server the metrics mentioned below. For further help please execute:

{{< highlight txt >}}
ocis-graph server --help
{{< / highlight >}}

#### Health

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exist code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
ocis-graph health --help
{{< / highlight >}}

### Metrics

This service provides some [Prometheus](https://prometheus.io/) metrics through the debug endpoint, you can optionally secure the metrics endpoint by some random token, which got to be configured through one of the flag `--debug-token` or the environment variable `GRAPH_DEBUG_TOKEN` mentioned above. By default the metrics endpoint is bound to `http://0.0.0.0:9124/metrics`.

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
