---
title: "oCIS with clamav"
date: 2024-05-21T14:04:00+01:00
weight: 101
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: ocis_clamav.md
---

{{< toc >}}

## Overview

- oCIS with standard clamav setup

[Find this example on GitHub](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_clamav)

The docker stack contains the following services:

oCIS itself, without any proxy in front of it, keep in mind,
the example is for demonstration purposes only and should not be used in production.

A pre-configured clamav container to virus scan files uploaded to oCIS.

## Server Deployment

The provided docker compose file is for local demonstration purposes only.
It is not recommended to use this setup in production.

## Local setup

`docker-compose up -d`

once all containers are up and running, you can access the oCIS instance at `https://localhost:9200`,
clamav could take some time to start up, so please be patient.
