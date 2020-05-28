---
title: "Getting started"
date: 2020-05-28T10:39:00+01:00
weight: 1
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: extensions-framework/getting-started.md
---

{{< toc >}}

## Create own extension
- [Create frontend]({{< ref "create-frontend.md" >}})

## How to load extension into Phoenix
All extensions are loaded when the user enters Phoenix. This is achieved with the help of [RequireJS](https://requirejs.org/).
To load your extensions they need to be registered in the config.json.

```json
{
  "external_apps": [
    {
      "id": "hello",
      "path": "http://localhost:9105/hello.js"
    },
    {
      "id": "myapp",
      "path": "http://localhost:6789/superapp.js"
    }
  ]
}
```

You can take a look at example of a full config.json [here](https://github.com/owncloud/phoenix/blob/master/config.json.sample-ocis).

## Test extension framework locally with ocis-simple
oCIS uses build tags to build different flavors of the binary. To see the extensions framework in action, we are going to reduce the scope a little and use the `simple` tag. Let us begin by creating a dedicated folder:

```sh
mkdir ocis-extension && ocis-extension
```

Following [https://github.com/owncloud/ocis](https://github.com/owncloud/ocis)

```sh
git clone https://github.com/owncloud/ocis.git
cd ocis

TAGS=simple make generate build
bin/ocis server
```

Open the browser at [https://localhost:9200](https://localhost:9200)

1. You land on the login screen. Click login
2. You are redirected to an IDP with a login mask. Use `einstein:relativity` to login (one of the three demo users)
3. You are redirected to [http://localhost:9200/#/hello](http://localhost:9200/#/hello) [the ocis-hello app](https://owncloud.github.io/extensions/ocis_hello/)
4. Replace `World` with something else and submit. You should see `Hello %something else%`
