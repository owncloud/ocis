---
title: "Tests"
weight: 90
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/settings
geekdocFilePath: tests.md
---

{{< toc >}}

## Requirements

You need a working installation of [the Go programming language](https://golang.org/), [the Node runtime](https://nodejs.org/) and [the Yarn package manager](https://yarnpkg.com/) installed to run the acceptance tests. You may also want to use [Docker](https://www.docker.com/) to start the necessary services in their respective containers.

## Acceptance Tests

Make sure you've cloned the [web frontend repo](https://github.com/owncloud/web/) and the [infinite scale repo](https://github.com/owncloud/ocis/) next to each other. If your file/folder structure is different, you'll have to change the paths below accordingly.

{{< hint info >}}
For now, an IDP configuration file gets generated once and will fail upon changing the oCIS url as done below. To avoid any clashes, remove this file before starting the tests:

```
rm ~/.ocis/idp/identifier-registration.yaml
```
{{< /hint >}}

### In the web repo

#### **Optional:** Build web to test local changes

Install dependencies and bundle the frontend with a watcher by running

```
yarn && yarn build:w
```

If you skip the step above, the currently bundled frontend from the oCIS binary will be used.

#### Dockerized acceptance test services

Start the necessary acceptance test services by using Docker (Compose):

```
docker compose up selenium middleware-ocis vnc
```

### In the oCIS repo

#### **Optional:** Build settings UI to test local changes

Navigate into the settings service via `cd ../settings/` and install dependencies and build the bundled settings UI with a watcher by running

```
yarn && yarn watch
```

#### Start oCIS from binary

Navigate into the oCIS directory inside the oCIS repository and build the oCIS binary by running

```
make clean build
```

Then, start oCIS from the binary via

```
OCIS_URL=https://host.docker.internal:9200 OCIS_INSECURE=true PROXY_ENABLE_BASIC_AUTH=true WEB_UI_CONFIG=../../web/dev/docker/ocis.web.config.json ./bin/ocis server
```

If you've built the web bundle locally in its repository, you also need to reference the bundle output in the above command: `WEB_ASSET_PATH=../../web/dist`

If you've built the settings UI bundle locally, you also need to reference the bundle output in the above command: `SETTINGS_ASSET_PATH=../settings/assets/`

#### Run settings acceptance tests

If you want visual feedback on the test run, visit http://host.docker.internal:6080/ in your browser and connect to the VNC client.

Navigate into the settings service via `cd ../settings/` and start the acceptance tests by running

```
SERVER_HOST=https://host.docker.internal:9200 BACKEND_HOST=https://host.docker.internal:9200 RUN_ON_OCIS=true NODE_TLS_REJECT_UNAUTHORIZED=0 WEB_PATH=../../web WEB_UI_CONFIG=../../web/tests/drone/config-ocis.json MIDDLEWARE_HOST=http://host.docker.internal:3000 ./ui/tests/run-acceptance-test.sh ./ui/tests/acceptance/features/
```
