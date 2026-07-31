---
title: "Releasing"
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/web
geekdocFilePath: releasing.md
---

{{< toc >}}

## Releasing

The Web Frontend source lives in-tree under `web/` and is built as part of the regular oCIS
generate/release process — there is no separate repository, version pin, or release step to
coordinate anymore.

Running `make generate` (or `make -C services/web ci-node-generate`) builds the web frontend
from `web/` via `pnpm install && pnpm build` and copies the resulting assets into
`services/web/assets/core`, from where they are embedded into the `ocis` binary via Go's
`//go:embed` directive in `services/web/web.go`.

To make changes to the web frontend, edit the source under `web/` directly and open a PR
against this repository, the same as any other oCIS change.
