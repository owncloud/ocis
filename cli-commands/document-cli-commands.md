---
title: Document CLI Commands
date: 2025-11-13T00:00:00+00:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/cli-commands/
geekdocFilePath: document-cli-commands.md
weight: 20
---

{{< toc >}}

Any CLI command that is added to Infinite Scale must be documented here in the dev docs and the [admin docs](https://doc.owncloud.com/ocis/latest/maintenance/commands/commands.html). Note that the admin docs primarily distinguish between online and offline commands because the structure of the documentation is different. Typically, any command documented in the developer documentation is integrated into the admin documentation and adapted according to the target audience. The description here is for developers; the admin docs are derived from it.

{{< hint info >}}
Note that any CLI command requires documentation. However, it may be decided that a CLI command will not be included in the admin documentation. In such a case, the reasons should be valid.
{{< /hint >}}

## Type of CLI Commands

There are three types of CLI commands that require different documentation locations:

1. Commands that are embedded in a service such as:\
`ocis storage-users uploads`
2. Commands that are service independent such as:\
`ocis trash purge-empty-dirs` or `ocis revisions purge`
3. `curl` commands that can be one of the above.

## Rules

* **Service dependent** CLI commands:\
  Add any CLI command into the repsective `README.md` **of the service**. Use as "template" an existing one such as in `services/storage-users/README.md` or `services/auth-app/README.md`. The content created will be transferred automatically to the service in the [Services]({{< ref "../services/" >}}) section.
* **Service independent** CLI commands:\
  Add any CLI command into the [Service Independent CLI]({{< ref "./service_independent_cli.md" >}}) documentation. See the link for an example how to do so.
