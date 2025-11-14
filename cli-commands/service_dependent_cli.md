---
title: Service Dependent CLI
date: 2025-11-13T00:00:00+00:00
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/cli-commands/
geekdocFilePath: service_dependent_cli.md
---

This document describes ocis CLI commands that are **embedded in a service**.

{{< toc >}}

## Common Parameters

The ocis package offers a variety of CLI commands for monitoring or repairing ocis installations. Most of these commands have common parameters such as:

* `--help` (or `-h`)\
  Use to print all available options.

* `--basePath` (or `-p`)\
  Needs to point to a storage provider, paths can vary depending on your ocis installation. Example paths are:
  ```bash
  .ocis/storage/users          # bare metal installation
  /var/tmp/ocis/storage/users  # docker installation
  ...
  ```

* `--dry-run`\
  This parameter, when available, defaults to `true` and must explicitly set to `false`.

* `--verbose` (or `-v`)\
  Get a more verbose output.

## List of CLI Commands

For CLI commands that are **embedded in a service**, see the following services:

* [Auth-App]({{< ref "../services/auth-app/" >}})
* [Graph]({{< ref "../services/graph/" >}})
* [Postprocessing]({{< ref "../services/postprocessing/" >}})
* [Storage-Users]({{< ref "../services/storage-users/" >}})
