---
title: Create a New CLI Command
date: 2025-11-13T00:00:00+00:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/cli-commands/
geekdocFilePath: create-new-cli-command.md
weight: 20
---

{{< toc >}}

{{< hint info >}}
Existing commands should be checked for commonly used options whenever a new CLI command is created, regardless of whether it is embedded in a service or not. For an example you can see the [Service Independent CLI]({{< ref "./service_independent_cli.md" >}}) documentation.
{{< /hint >}}

## CLI Embedded in a Service

These commands are usually located in the `<service-name>/pkg/command` subfolder.

## CLI Independent of a Service

These commands are located in the `ocis/pkg/command` subfolder.
