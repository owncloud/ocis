# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.7.0] - 2024-04-22

### ⚠️ Notice ⚠️

This release is the last to support Go `1.18`. The next release will require at least Go `1.19`.

### Changed

- Upgrade to `v1.14.0` of `go.opentelemetry.io/otel`. (#38)
- Upgrade to `v1.17.0` of `go.opentelemetry.io/otel/semconv`. (#38)
- Adjust Go version for both `examples/basic` & `examples/multi-services` to `1.18` & `go.opentelemetry.io/otel` to `v1.14.0`. (#38)
- Change `http.server_name` attributes to `net.host.name`, this is because semconv is removing this attribute for http. (#38)

### Removed

- Remove `http.target` attribute on implementation & tests based on [this comment](https://github.com/open-telemetry/opentelemetry-go/blob/v1.17.0/semconv/internal/v2/http.go#L160-L165). (#39)
- Drop support for Go `<1.18`. (#38)

## [0.6.0] - 2024-04-02

### ⚠️ Notice ⚠️

This release is the last to support Go `1.15`. The next release will require at least Go `1.18`.

### Added

- Add `WithTraceIDResponseHeader` option to enable adding trace id into response header. (#36)
- Add multiple go versions test scripts for local and CI pipeline. (#29)
- Add compatibility testing for `ubuntu`, `macos` and `windows`. (#32)
- Add repo essentials docs. (#33)

### Changed

- Upgrade to `v5.0.12` of `go-chi/chi`. (#29)
- Upgrade to `v1.10.0` of `go.opentelemetry.io/otel`. (#29)
- Upgrade to `v1.12.0` of `go.opentelemetry.io/otel/semconv`. (#29)
- Set the required go version for both `examples/basic` & `examples/multi-services` to `1.15`, `go-chi/chi` to `v5.0.12`, & `go.opentelemetry.io/otel` to `v1.10.0` (#35)

## [0.5.2] - 2024-03-25

### Fixed

- Fix empty status code. (#30)

### Changed

- Return `http.StatusOK` (200) as a default `http.status_code` span attribute. (#30)

## [0.5.1] - 2023-02-18

### Fixed

- Fix broken empty routes. (#18)

### Changed

- Upgrade to `v5.0.8` of `go-chi/chi`.

## [0.5.0] - 2022-10-02

### Added

- Add multi services example. (#9)
- Add `WithFilter()` option to ignore tracing in certain endpoints. (#11)

## [0.4.0] - 2022-02-22

### Added

- Add Option `WithRequestMethodInSpanName()` to handle vendor that do not include HTTP request method as mentioned in #6. (#7)
- Refine description for `WithChiRoutes()` option to announce it is possible to override the span name in underlying handler with this option.

### Changed

## [0.3.0] - 2022-01-18

### Fixed

- Fix both `docker-compose.yml` & `Dockerfile` in the example. (#5)

### Added

- Add `WithChiRoutes()` option to make the middleware able to determine full route pattern on span creation. (#5)
- Set all known span attributes on span creation rather than set them after request is being executed. (#5)

## [0.2.1] - 2022-01-08

### Added

- Add build example to CI pipeline. (#2)

### Changed

- Use `ctx.RoutePattern()` to get span name, this is to strip out noisy wildcard pattern. (#1)

## [0.2.0] - 2021-10-18

### Added

- Set service name on tracer provider from code example.

### Changed

- Update dependencies in go.mod
- Upgrade to `v1.0.1` of `go.opentelemetry.io/otel`.
- Upgrade to `v5.0.4` of `go-chi/chi`.
- Update latest test to use `otelmux` format.

### Removed

- Remove `HTTPResponseContentLengthKey`
- Remove `HTTPTargetKey`, since automatically set in `HTTPServerAttributesFromHTTPRequest`

## [0.1.0] - 2021-08-11

This is the first release of otelchi.
It contains instrumentation for trace and depends on:

- otel => `v1.0.0-RC2`
- go-chi/chi => `v5.0.3`

### Added

- Instrumentation for trace.
- CI files.
- Example code for a basic usage.
- Apache-2.0 license.

[Unreleased]: https://github.com/riandyrn/otelchi/compare/v0.7.0...HEAD
[0.7.0]: https://github.com/riandyrn/otelchi/releases/tag/v0.7.0
[0.6.0]: https://github.com/riandyrn/otelchi/releases/tag/v0.6.0
[0.5.2]: https://github.com/riandyrn/otelchi/releases/tag/v0.5.2
[0.5.1]: https://github.com/riandyrn/otelchi/releases/tag/v0.5.1
[0.5.0]: https://github.com/riandyrn/otelchi/releases/tag/v0.5.0
[0.4.0]: https://github.com/riandyrn/otelchi/releases/tag/v0.4.0
[0.3.0]: https://github.com/riandyrn/otelchi/releases/tag/v0.3.0
[0.2.1]: https://github.com/riandyrn/otelchi/releases/tag/v0.2.1
[0.2.0]: https://github.com/riandyrn/otelchi/releases/tag/v0.2.0
[0.1.0]: https://github.com/riandyrn/otelchi/releases/tag/v0.1.0
