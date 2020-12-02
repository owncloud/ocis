---
title: "Extensions"
date: 2020-02-27T20:35:00+01:00
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: extensions.md
---

{{< toc >}}

## How to build and run ocis-simple

ocis uses build tags to build different flavors of the binary. In order to work on a new extension we are going to reduce the scope a little and use the `simple` tag. Let us begin by creating a dedicated folder:

```console
mkdir ocis-extension-workshop && ocis-extension-workshop
```

Following https://github.com/owncloud/ocis

```console
git clone https://github.com/owncloud/ocis.git
cd ocis

TAGS=simple make generate build
```

*Q: Can you specify which version of phoenix to use?*
*A: No, the phoenix that is used is compiled into the [assets of ocis-phoenix](https://github.com/owncloud/ocis-phoenix/blob/master/pkg/assets/embed.go) which is currently not automatically updated. We'll see how to use a custom phoenix later.*

`bin/ocis server`

Open the browser at http://localhost:9100

1. You land on the login screen. click login
2. You are redirected to an idp at http://localhost:9140/oauth2/auth with a login mask. Use `einstein:relativity`to login (one of the three demo users)
3. You are redirected to http://localhost:9100/#/hello the ocis-hello app
4. Replace `World` with something else and submit. You should see `Hello %something else%`

*Q: One of the required ports is already in use. Ocis seems to be trying to restart the service over and over. What gives?*
*A: Using the ocis binary to start the server will case ocis to keep track of the different services and restart them in case they crash.*

## Hacking ocis-hello

go back to the ocis-extension-workshop folder

```console
cd ..
```

Following https://github.com/owncloud/ocis-hello

```
git clone https://github.com/owncloud/ocis-hello.git
cd ocis-hello

yarn install
# this actually creates the assets
yarn build

# this will compile the assets into the binary
make generate build
```

Two options:
1. run only the necessery services from ocis and ocis-hello independently
2. compile ocis with the updated ocis-hello

### Option 1:
get a list of ocis services:

```console
ps ax | grep ocis
```

Try to kill `ocis hello`

Remember: for now, killing a service will cause ocis to restart it. This is subject to change.

In order to be able to manage the processes ourselves we need to start them independently:

`bin/ocis server` starts the same services as:

```
bin/ocis micro &
bin/ocis phoenix &
bin/ocis hello &
bin/ocis reva &
```

Now we can kill the `ocis hello` and use our custom built ocis-hello binary:

```console
cd ../ocis-hello
bin/ocis-hello server
```

## Hacking phoenix (and ocis-phoenix)

Following https://github.com/owncloud/phoenix we are going to build the current phoenix

```
git clone https://github.com/owncloud/phoenix.git
cd phoenix

yarn install
yarn dist
```

We can tell ocis to use the compiled assets:

Kill `ocis phoenix`, then use the compiled assets when starting phoenix.

```console
cd ../ocis
PHOENIX_ASSET_PATH="`pwd`/../phoenix/dist" bin/ocis phoenix
```

## The ownCloud design system

The [ownCloud design system](https://owncloud.design/) contains a set of ownCloud vue components for phoenix or your own ocis extensions. Please use it for a consistent look and feel.

## External phoenix apps

This is what hello is: copy and extend!

1. Phoenix is configured using the config.json which is served by the phoenix service (either `bin/ocis phoenix` or `bin/ocis-phoenix server`)

2. point ocis phoenix to the web config which you extended with an external app:
`PHOENIX_WEB_CONFIG="`pwd`/../phoenix/config.json" PHOENIX_ASSET_PATH="`pwd`/../phoenix/dist" bin/ocis phoenix`

```json
{
  "server": "http://localhost:9140",
  "theme": "owncloud",
  "version": "0.1.0",
  "openIdConnect": {
    "metadata_url": "http://localhost:9140/.well-known/openid-configuration",
    "authority": "http://localhost:9140",
    "client_id": "phoenix",
    "response_type": "code",
    "scope": "openid profile email"
  },
  "apps": [],
  "external_apps": [
    {
      "id": "hello",
      "path": "http://localhost:9105/hello.js",
      "config": {
        "url": "http://localhost:9105"
      }
    },
    {
      "id": "myapp",
      "path": "http://localhost:6789/superapp.js",
      "config": {
        "backend": "http://someserver:1234",
        "myconfig": "is awesome"
      }
    }
  ]
}
```

## Phoenix extension points

{{< hint info >}}
For an up to date list check out [the phoenix documentation](https://github.com/owncloud/phoenix/issues/2423).
{{< /hint >}}

Several ones available:

### Phoenix core
- App switcher (defined in config.json)
- App container (loads UI of your extension)

### Files app
- File action
- Create new file action
- Sidebar
- Quick access for sidebar inside of file actions (in the file row)

Example of a file action in the `app.js`:
```js
const appInfo = {
  name: 'MarkdownEditor',
  id: 'markdown-editor',
  icon: 'text',
  isFileEditor: true,
  extensions: [{
    extension: 'txt',
    newFileMenu: {
      menuTitle ($gettext) {
        return $gettext('Create new plain text file…')
      }
    }
  },
  {
    extension: 'md',
    newFileMenu: {
      menuTitle ($gettext) {
        return $gettext('Create new mark-down file…')
      }
    }
  }]
}
```

For the side bar have a look at the files app, `defaults.js` & `fileSideBars`

## API driven development

Until now we only had a look at the ui and how the extensions are managed on the cli. But how do apps actually talk to the server?

Short answer: any way you like

Long answer: micro and ocis-hello follow a protocol driven development:

- specify the API using protobuf
- generate client and server code
- evolve based on the protocol

- CS3 api uses protobuf as well and uses GRPC

- ocis uses go-micro, which provides http and grpc gateways
- the gateways and protocols are optional

- owncloud and kopano are looking into a [MS graph](https://developer.microsoft.com/de-de/graph) like api to handle phoenix requests.
  - they might be about user, contacrs, calendars ... which is covered by the graph api
  - we want to integrate with eg. kopano and provide a commen api (file sync and share is covered as well)

- as an example for protobuf take a look at [ocis-hello](https://github.com/owncloud/ocis-hello/tree/master/pkg/proto/v0)
