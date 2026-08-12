# Scorch Segment API

[![Go Reference](https://pkg.go.dev/badge/github.com/blevesearch/scorch_segment_api/v2.svg)](https://pkg.go.dev/github.com/blevesearch/scorch_segment_api/v2)
[![Tests](https://github.com/blevesearch/scorch_segment_api/actions/workflows/tests.yml/badge.svg?branch=master&event=push)](https://github.com/blevesearch/scorch_segment_api/actions/workflows/tests.yml?query=event%3Apush+branch%3Amaster)

Scorch supports a pluggable Segment interface.

By placing these interfaces in their own, *hopefully* slowly evolving module, it frees up Scorch and the underlying segment to each introduce new major versions without interfering with one another.

With that in mind, we anticipate introducing non-breaking changes only to this module, and keeping the major version at 1.x for some time.
