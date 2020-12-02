---
title: "Configuration Hints"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/glauth
geekdocFilePath: configuration-hints.md
---

{{< toc >}}

## Configuration hints

The default setup does not use a fallback backend. It can be enabled by setting the `GLAUTH_FALLBACK_DATASTORE` environment variable.

When using `owncloud` make sure to use the full URL to the [ownCloud 10 graph api app](https://github.com/owncloud/graphapi) endpoint, eg.: `GLAUTH_FALLBACK_SERVERS="https://demo.owncloud.com/apps/graphapi/v1.0"`
