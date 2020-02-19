---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis-konnectd
geekdocEditPath: edit/master/docs
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

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below.

#### Global

KONNECTD_CONFIG_FILE
: Path to config file, empty default value

KONNECTD_LOG_LEVEL
: Set logging level, defaults to `info`

KONNECTD_LOG_COLOR
: Enable colored logging, defaults to `true`

KONNECTD_LOG_PRETTY
: Enable pretty logging, defaults to `true`

#### Server

KONNECTD_TRACING_ENABLED
: Enable sending traces, defaults to `false`

KONNECTD_TRACING_TYPE
: Tracing backend type, defaults to `jaeger`

KONNECTD_TRACING_ENDPOINT
: Endpoint for the agent, empty default value

KONNECTD_TRACING_COLLECTOR
: Endpoint for the collector, empty default value

KONNECTD_TRACING_SERVICE
: Service name for tracing, defaults to `konnectd`

KONNECTD_DEBUG_ADDR
: Address to bind debug server, defaults to `0.0.0.0:9134`

KONNECTD_DEBUG_TOKEN
: Token to grant metrics access, empty default value

KONNECTD_DEBUG_PPROF
: Enable pprof debugging, defaults to `false`

KONNECTD_DEBUG_ZPAGES
: Enable zpages debugging, defaults to `false`

KONNECTD_HTTP_ADDR
: Address to bind http server, defaults to `0.0.0.0:9130`

KONNECTD_HTTP_ROOT
: Root path of http server, defaults to `/`

KONNECTD_HTTP_NAMESPACE
: Set the base namespace for service discovery, defaults to `com.owncloud.web`

KONNECTD_IDENTITY_MANAGER
: Identity manager (one of ldap,kc,cookie,dummy), defaults to `ldap`

KONNECTD_TRANSPORT_TLS_CERT
: Certificate file for transport encryption, uses a temporary dev-cert if empty

KONNECTD_TRANSPORT_TLS_KEY
: Secret file for transport encryption, uses a temporary dev-cert if empty

KONNECTD_ISS
: OIDC issuer URL, defaults to `https://localhost:9130`

KONNECTD_SIGNING_PRIVATE_KEY
: Full path to PEM encoded private key file (must match the --signing-method algorithm)

KONNECTD_SIGNING_KID
: Value of kid field to use in created tokens (uniquely identifying the signing-private-key), empty default value

KONNECTD_VALIDATION_KEYS_PATH
: Full path to a folder containg PEM encoded private or public key files used for token validaton (file name without extension is used as kid), empty default value

KONNECTD_ENCRYPTION_SECRET
: Full path to a file containing a %d bytes secret key, empty default value

KONNECTD_SIGNING_METHOD
: JWT default signing method, defaults to `PS256`

KONNECTD_URI_BASE_PATH
: Custom base path for URI endpoints, empty default value

KONNECTD_SIGN_IN_URI
: Custom redirection URI to sign-in form, empty default value

KONNECTD_SIGN_OUT_URI
: Custom redirection URI to signed-out goodbye page, empty default value

KONNECTD_ENDPOINT_URI
: Custom authorization endpoint URI, empty default value

KONNECTD_ENDSESSION_ENDPOINT_URI
: Custom endsession endpoint URI, empty default value

KONNECTD_ASSET_PATH
: Path to custom assets, empty default value

KONNECTD_IDENTIFIER_CLIENT_PATH
: Path to the identifier web client base folder, defaults to `/var/tmp/konnectd`

KONNECTD_IDENTIFIER_REGISTRATION_CONF
: Path to a identifier-registration.yaml configuration file, defaults to `./config/identifier-registration.yaml`

KONNECTD_IDENTIFIER_SCOPES_CONF
: Path to a scopes.yaml configuration file, empty default value

KONNECTD_INSECURE
: Disable TLS certificate and hostname validation

KONNECTD_TLS
: Use TLS (disable only if konnectd is behind a TLS-terminating reverse-proxy), defaults to `true`

KONNECTD_TRUSTED_PROXY
: List of trusted proxy IP or IP network(s) (usage: KONNECTD_TRUSTED_PROXY=x.x.x.x y.y.y.y)

KONNECTD_ALLOW_SCOPE
: Allow OAuth 2 scope(s) (usage: KONNECTD_ALLOW_SCOPE=A B C)

KONNECTD_ALLOW_CLIENT_GUESTS
: Allow sign in of client controlled guest users

KONNECTD_ALLOW_DYNAMIC_CLIENT_REGISTRATION
: Allow dynamic OAuth2 client registration

KONNECTD_DISABLE_IDENTIFIER_WEBAPP
: Disable built-in identifier-webapp to use a frontend hosted elsewhere. Per default we use the built-in webapp. If set to false --identifier-client-path must be provided, defaults to `true`


#### Health

KONNECTD_DEBUG_ADDR
: Address to debug endpoint, defaults to `0.0.0.0:9134`

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
: Service name for tracing, defaults to `konnectd`

--debug-addr
: Address to bind debug server, defaults to `0.0.0.0:9134`

--debug-token
: Token to grant metrics access, empty default value

--debug-pprof
: Enable pprof debugging, defaults to `false`

--debug-zpages
: Enable zpages debugging, defaults to `false`

--http-addr
: Address to bind http server, defaults to `0.0.0.0:9130`

--http-root
: Root path of http server, defaults to `/`

--http-namespace
: Set the base namespace for service discovery, defaults to `com.owncloud.web`

--identity-manager
: Identity manager (one of ldap,kc,cookie,dummy), defaults to `ldap`

--transport-tls-cert
: Certificate file for transport encryption, uses a temporary dev-cert if empty

--transport-tls-key
: Key file for transport encryption, uses a temporary dev-cert if empty

--iss
: OIDC issuer URL, defaults to `https://localhost:9130`

--signing-private-key
: Full path to PEM encoded private key file (must match the --signing-method algorithm)

--signing-kid
: Value of kid field to use in created tokens (uniquely identifying the signing-private-key), empty default value

--validation-keys-path
: Full path to a folder containg PEM encoded private or public key files used for token validaton (file name without extension is used as kid), empty default value

--encryption-secret
: Full path to a file containing a 32 bytes secret key, empty default value

--signing-method
: JWT default signing method, defaults to `PS256`

--uri-base-path
: Custom base path for URI endpoints, empty default value

--sign-in-uri
: Custom redirection URI to sign-in form, empty default value

--signed-out-uri
: Custom redirection URI to signed-out goodbye page, empty default value

--authorization-endpoint-uri
: Custom authorization endpoint URI, empty default value

--endsession-endpoint-uri
: Custom endsession endpoint URI, empty default value

--asset-path
: Path to custom assets, empty default value

--identifier-client-path
: Path to the identifier web client base folder, defaults to `/var/tmp/konnectd`

--identifier-registration-conf
: Path to a identifier-registration.yaml configuration file, defaults to `./config/identifier-registration.yaml`

--identifier-scopes-conf
: Path to a scopes.yaml configuration file, empty default value

--insecure
: Disable TLS certificate and hostname validation

--tls
: Use TLS (disable only if konnectd is behind a TLS-terminating reverse-proxy), defaults to `true`

--trusted-proxy
: List of trusted proxy IP or IP network (usage: --trusted-proxy x.x.x.x --trusted-proxy y.y.y.y)

--allow-scope
: Allow OAuth 2 scope (usage: --allow-scope a --allow-scope b ...)

--allow-client-guests
: Allow sign in of client controlled guest users

--allow-dynamic-client-registration
: Allow dynamic OAuth2 client registration

--disable-identifier-webapp
:  Disable built-in identifier-webapp to use a frontend hosted elsewhere. Per default we use the built-in webapp. If set to false --identifier-client-path must be provided, defaults to `true`


#### Health

--debug-addr
: Address to debug endpoint, defaults to `0.0.0.0:9134`

### Configuration file

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis-konnectd/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/konnectd.yml`, `${HOME}/.ocis/konnectd.yml` or `$(pwd)/config/konnectd.yml`.

## Usage

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formated help output if you execute the binary via `ocis-konnectd --help`.

### Server

The server command is used to start the http and debug server on two addresses within a single process. The http server is serving the general webservice while the debug server is used for health check, readiness check and to server the metrics mentioned below. For further help please execute:

{{< highlight txt >}}
ocis-konnectd server --help
{{< / highlight >}}

### Health

The health command is used to execute a health check, if the exit code equals zero the service should be up and running, if the exist code is greater than zero the service is not in a healthy state. Generally this command is used within our Docker containers, it could also be used within Kubernetes.

{{< highlight txt >}}
ocis-konnectd health --help
{{< / highlight >}}

## Metrics

This service provides some [Prometheus](https://prometheus.io/) metrics through the debug endpoint, you can optionally secure the metrics endpoint by some random token, which got to be configured through one of the flag `--debug-token` or the environment variable `KONNECTD_DEBUG_TOKEN` mentioned above. By default the metrics endpoint is bound to `http://0.0.0.0:9134/metrics`.

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
