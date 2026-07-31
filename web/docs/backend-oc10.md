---
title: "Setup With ownCloud Classic"
date: 2020-04-15T00:00:00+00:00
weight: 40
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs
geekdocFilePath: backend-oc10.md
---

{{< toc >}}

## Compatibility

Please note that the usage of Web UI and ownCloud Classic as backend is not recommended starting with version 7.1.0 of the Web UI. Therefore, this section only applies to versions < 7.1.0.

## Prerequisites

Decide on which host and port Web will be served, for example `https://web-host:9100/web-path/`.
In this document, we will refer to the following:
- `<web-url>` as the full URL, for example `https://web-host:9100/web-path/`
- `<web-domain>` as the protocol, domain and port, for example: `https://web-host:9100`

## Setting up ownCloud Classic

Make sure you have [ownCloud Classic](https://owncloud.org/download/#owncloud-server) already installed.

### Adjusting config.php

Add the following entries to config/config.php:

- tell ownCloud where Web is located:
```
'web.baseUrl' => '<web-url>',
```

- add a CORS domain entry for Web in config.php:
```
'cors.allowed-domains' => ['<web-domain>'],
```

### Setting up OAuth2

To connect to the ownCloud server, it is necessary to set it up with OAuth2.

Install and enable the [oauth2 app](https://marketplace.owncloud.com/apps/oauth2):
```bash
% occ market:install oauth2
% occ app:enable oauth2
```

Login as administrator in the ownCloud Classic web interface and go to the "User Authentication" section in the admin settings and add an entry for Web as follows:

- pick an arbitrary name for the client
- set the redirection URI to `<web-url>/oidc-callback.html`
- make sure to take note of the **client identifier** value as it will be needed in the Web configuration later on

### Setting up Web

In the local Web checkout, copy the `config/config.json.sample-oc10` file to `config/config.json` and adjust it accordingly:

- Set the "server" key to the URL of the ownCloud server including path. If the URL contains a path, please also add a **trailing slash** there.
- Set the "clientId" key to the **client identifier** as copied from the "User Authentication" section before.
- Adjust "url" and "authUrl" using the ownCloud server URL as prefix for both
- Optionally adjust "apps" for the list of apps to be loaded. These match the app names inside the "apps" folder.

## Running Web

- if running from source, make sure to [build Web]({{< ref "./building.md" >}}) first
- run by launching a rollup dev server `pnpm serve`
- when working on the Web code, rollup will recompile the code automatically

## Running acceptance tests

For testing, please refer to the [testing docs]({{< ref "testing/_index.md" >}})
