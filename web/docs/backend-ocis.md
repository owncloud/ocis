---
title: "Setup With oCIS"
date: 2020-04-15T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs
geekdocFilePath: backend-ocis.md
---

{{< toc >}}

## Setting up Web

- Clone the [repository](https://github.com/owncloud/web/)
- Initally install all dependencies by running `pnpm install`
- Copy `./config/config.json.sample-ocis` to `./config/config.json` and adjust values if required

## Running Web

- Start bundling web with a watcher by running `pnpm build:w`

## Setting up oCIS

- Setup oCIS by following the [setup instructions](https://owncloud.dev/ocis/getting-started/)
- Start oCIS with local links to your bundled web frontend and config by running `WEB_ASSET_CORE_PATH=../../web/dist WEB_UI_CONFIG_FILE=../../web/dist/config.json OCIS_INSECURE=true IDM_CREATE_DEMO_USERS=true ./bin/ocis server` (and make sure to adjust paths as necessary)

## Start oCIS

- open [https://localhost:9200](https://localhost:9200) and accept the certificate
- when signing in, use one of the [available demo users](https://owncloud.dev/ocis/getting-started/demo-users/)
- whenever code changes are made, you need to manually reload the browser page (no hot reload)

## Running tests

For testing, please refer to the [testing docs]({{< ref "testing/_index.md" >}})

