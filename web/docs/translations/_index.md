---
title: "Translations"
date: 2025-05-21T00:00:00+00:00
weight: 55
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/translations
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

{{< toc >}}

## Introduction

Translations are essential for the web frontend. Sources are located in subfolders of the packages folder. A nightly sync automatically extracts and synchronizes data, see the [Add Translations](https://owncloud.dev/services/general-info/add-translations/) documentation for more details, but it can also be triggered manually.

## Preparation

There are two important prerequisites to work with Transifex:

* In order to work with local translations, you must have installed the [Transifex CLI Client](https://developers.transifex.com/docs/cli) `tx` which must be accessible via the search path in order to be callable from anywhere.\
{{< hint warning >}}
When using the provided curl command from the link above, tx will install in the current directory. This will most likely cause an issue because there is no search path pointing to it. To avoid this issue, change to a directory included in the executable search path before installation. For example, use the directory `/usr/local/bin`. Otherwise, adapt the PATH environment variable for your shell.
{{< /hint >}}

* Transifex requires a token for the access with `tx`. This token can be generated for free and is connected to your Transifex account. See [API token](https://app.transifex.com/user/settings/api/) for more details. For ease of use, add this token to your shell environment variables `export TX_TOKEN=xxx`. 

* See [Using the client](https://developers.transifex.com/docs/using-the-client) for details on how using `tx`.
