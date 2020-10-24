---
title: "Getting Started"
date: 2020-02-27T20:35:00+01:00
weight: -15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Installation

So far we are offering two different variants for the installation. You can choose between [Docker](https://www.docker.com/) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

### Docker

Docker images for ocis are hosted on https://hub.docker.com/r/owncloud/ocis.

The `latest` tag always reflects the current master branch.

```console
docker pull owncloud/ocis
```

#### Dependencies

- Running ocis currently needs a working Redis caching server
- The default storage location in the container is `/var/tmp/reva/data`. You may want to create a volume to persist the files in the primary storage

#### Docker compose

You can use our docker-compose [playground example](https://github.com/owncloud-docker/compose-playground/tree/master/ocis) to run ocis with dependencies with a single command in a docker network.

```console
git clone git@github.com:owncloud-docker/compose-playground.git
cd compose-playground/ocis
docker-compose -f ocis.yml -f ../cache/redis-ocis.yml up
```

### Binaries

The pre-built binaries for different platforms are downloadable at https://download.owncloud.com/ocis/ocis/ . Specific releases are organized in separate folders. They are in sync which every release tag on GitHub. The binaries from the current master branch can be found in https://download.owncloud.com/ocis/ocis/testing/

```console
curl https://download.owncloud.com/ocis/ocis/1.0.0-beta1/ocis-1.0.0-beta1-darwin-amd64 --output ocis
chmod +x ocis
./ocis server
```

#### Dependencies

- Running ocis currently needs a working Redis caching server
- The default promary storage location is `/var/tmp/reva/data`. You can change that value by configuration.

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

## Quickstart for Developers

Following https://github.com/owncloud/ocis#development

```console
git clone https://github.com/owncloud/ocis.git
cd ocis
make generate build
```

Open https://localhost:9200 and login using one of the demo accounts:

```console
einstein:relativity
marie:radioactivity
richard:superfluidity
```

There are admin demo accounts:
```console
moss:vista
admin:admin
```

## Runtime

Included with the ocis binary is embedded a go-micro runtime that is in charge of starting services as a fork of the master process. This provides complete control over the services. Ocis extensions can be added as part of this runtime.

```console
./bin/ocis micro
```

This will currently boot:

```console
com.owncloud.api
com.owncloud.http.broker
com.owncloud.proxy
com.owncloud.registry
com.owncloud.router
com.owncloud.runtime
com.owncloud.web
go.micro.http.broker
```

Further ocis extensions can be added to the runtime via the ocis command like:

```console
./bin/ocis hello
```

Which will register:

```console
com.owncloud.web.hello
com.owncloud.api.hello
```

To the list of available services.

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
