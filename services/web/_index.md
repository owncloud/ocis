---
title: Web Service
date: 2023-04-11T01:26:08.65456837Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/web
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The web service embeds and serves the static files for the [Infinite Scale Web Client](https://github.com/owncloud/web).  
Note that clients will respond with a connection error if the web service is not available.
The web service also provides a minimal API for branding functionality like changing the logo shown.

## Table of Contents

* [Custom Compiled Web Assets](#custom-compiled-web-assets)
* [Example Yaml Config](#example-yaml-config)

## Custom Compiled Web Assets

If you want to use your custom compiled web client assets instead of the embedded ones, then you can do that by setting the `WEB_ASSET_PATH` variable to point to your compiled files. See [ownCloud Web / Getting Started](https://owncloud.dev/clients/web/getting-started/) and [ownCloud Web / Setup with oCIS](https://owncloud.dev/clients/web/backend-ocis/) for more details.

## Example Yaml Config

{{< include file="services/_includes/web-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/web_configvars.md" >}}

