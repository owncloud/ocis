---
title: "Architecture Decisions"
date: 2021-02-10T20:21:00+01:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

oCIS is documenting architecture decisions using [Markdown Architectural Decision Records](https://adr.github.io/madr/) (MADR), following [Documenting Architecture Decisions by Michael Nygard](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions).

{{< toc >}}

To manage the records we use [butonic/adr-tools](https://github.com/butonic/adr-tools), a fork of the original [npryce/adr-tools](https://github.com/npryce/adr-tools), based on [a pull request that should have added MADR support](https://github.com/npryce/adr-tools/pull/43). It also supports a YAML header that is used by our Hugo based doc generation