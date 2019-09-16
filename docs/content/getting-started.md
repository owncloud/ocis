---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
anchor: "getting-started"
weight: 20
---

## Installation

So far we are offering two different variants for the installation. You can choose between [Docker](docker) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

### Docker

TBD

### Binaries

TBD

## Configuration

We provide overall three different variants of configuration. The variant based on environment variables and commandline flags are split up into global values and command-specific values.

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below.

#### Global

PHOENIX_LOG_LEVEL
: Set logging level, defaults to `info`

PHOENIX_LOG_COLOR
: Enable colored logging, defaults to `true`

PHOENIX_LOG_PRETTY
: Enable pretty logging, defaults to `true`

#### Server

PHOENIX_DEBUG_ADDR
: Address to bind debug server, defaults to `0.0.0.0:8090`

PHOENIX_DEBUG_TOKEN
: Token to grant metrics access, empty default value

PHOENIX_DEBUG_PPROF
: Enable pprof debugging, defaults to `false`

PHOENIX_HTTP_ADDR
: Address to bind http server, defaults to `0.0.0.0:8080`

PHOENIX_HTTP_ROOT
: Root path for http endpoint, defaults to `/`

PHOENIX_ASSET_PATH
: Path to custom assets, empty default value

PHOENIX_CONFIG_FILE
: Path to phoenix config, empty default value

#### Health

PHOENIX_DEBUG_ADDR
: Address to debug endpoint, defaults to `0.0.0.0:8090`

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below.

#### Global

--log-level
: Set logging level, defaults to `info`

--log-color
: Enable colored logging, defaults to `true`

--log-pretty
: Enable pretty logging, defaults to `true`

#### Server

--debug-addr
: Address to bind debug server, defaults to `0.0.0.0:8090`

--debug-token
: Token to grant metrics access, empty default value

--debug-pprof
: Enable pprof debugging, defaults to `false`

--http-addr
: Address to bind http server, defaults to `0.0.0.0:8080`

--http-root
: Root path for http endpoint, defaults to `/`

--asset-path
: Path to custom assets, empty default value

--config-file
: Path to phoenix config, empty default value

#### Health

--debug-addr
: Address to debug endpoint, defaults to `0.0.0.0:8090`

### Configuration file

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](repo), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/phoenix.yml`, `${HOME}/.ocis/phoenix.yml` or `$(pwd)/config/phoenix.yml`.

## Usage

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formated help output if you execute the binary via `ocis-phoenix --help`.

### Server

The server command is used to start the http and debug server on two addresses within a single process. The http server is serving the general webservice while the debug server is used for health check, readiness check and to server the metrics mentioned below. For further help please execute:

{{< highlight txt >}}
ocis-phoenix server --help
{{< / highlight >}}

### Health

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exist code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
ocis-phoenix health --help
{{< / highlight >}}

## Metrics

This service provides some [Prometheus](prom) metrics through the debug endpoint, you can optionally secure the metrics endpoint by some random token, which got to be configured through one of the flag `--debug-token` or the environment variable `PHOENIX_DEBUG_TOKEN` mentioned above. By default the metrics endpoint is bound to `http://0.0.0.0:8090/metrics`.

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

[docker]: https://www.docker.com/
[repo]: https://github.com/owncloud/ocis-phoenix/tree/master/config
[prom]: https://prometheus.io/
