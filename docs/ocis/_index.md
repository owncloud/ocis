---
title: "Infinite Scale"
date: 2020-02-27T20:35:00+01:00
weight: -10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: _index.md
---

This tool provides a single entrypoint for the whole ownCloud Infinite Scale stack.

{{< mermaid class="text-center">}}
graph TD
ocis-proxy -->
    ocis-konnectd & ocis-phoenix & ocis-thumbnails & ocis-ocs & ocis-webdav

ocis-phoenix --> ocis-reva-fronted
ocis-reva-fronted --> ocis-reva-gateway
ocis-konnectd --> ocis-glauth


ocis-reva-gateway --> ocis-reva-users
ocis-reva-gateway --> ocis-reva-authbasic
ocis-reva-gateway --> ocis-reva-auth-bearer

ocis-reva-gateway --> ocis-reva-sharing
ocis-reva-gateway --> ocis-reva-storage-home-*
ocis-reva-storage-home-* --> ocis-reva-storage-home-*-data
ocis-reva-sharing --> redis
{{< /mermaid >}}
