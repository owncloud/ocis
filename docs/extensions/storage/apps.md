---
title: "Apps"
date: 2018-05-02T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: apps.md
---

oCIS is all about files. But most of the time you want to do something with files that is beyond the basic upload, download and share behavior. Therefore, oCIS has a concept for apps, that can handle specific file types, so called mime types.

## App provider capability

The capabilities endpoint (eg. `https://localhost:9200/ocs/v1.php/cloud/capabilities?format=json`) gives you following capabilities which are relevant for the app provider:

```json
{
  "ocs": {
    "data": {
      "capabilities": {
        "files": {
          "app_providers": [
            {
              "enabled": true,
              "version": "1.0.0",
              "apps_url": "/app/list",
              "open_url": "/app/open"
            }
          ]
        }
      }
    }
  }
}
```

Please note that there might be two or more app providers with different versions. This is not be expected to happen on a regular basis. It was designed for a possible migration period for clients when the app provider needs a breaking change.

## App registry

The app registry is the single point where all apps register themselves and their respective supported mime types.

### Mime type configuration / creation allow list

The apps will register their supported mime types automatically, so that users can open supported files with them.

Administrators can set default applications on a per mime type basis and also allow the creation of new files for certain mime types. This per mime type configuration also features a description, file extension option and an icon.

In order to modify the mime type config you need to set `STORAGE_APP_REGISTRY_MIMETYPES_JSON=.../mimetypes.json` to a valid JSON file with content like this:

```json
[
  {
    "mime_type": "applition/vnd.oasis.opendocument.text",
    "extension": "odt",
    "name": "OpenDocument",
    "description": "OpenDocument text document",
    "icon": "https://some-website.test/opendocument-text-icon.png",
    "default_app": "Collabora",
    "allow_creation": true
  },
  {
    "mime_type": "application/vnd.oasis.opendocument.spreadsheet",
    "extension": "ods",
    "name": "OpenSpreadsheet",
    "description": "OpenDocument spreadsheet document",
    "icon": "",
    "default_app": "Collabora",
    "allow_creation": false
  }
]
```

Fields:

- `mime_type` is the mime type you want to configure
- `extension` is the file extension to be used for new files
- `name` is the name of the file / mime type
- `description` is a human readable description of the file / mime type
- `icon` URL to an icon which should be used for that mime type
- `default_app` name of the default app which opens this mime type when the user doesn't specify one
- `allow_creation` is wether a user should be able to create new file from that mime type (`true` or `false`)

### Listing available apps / mime types

Clients, for example ownCloud Web, need to offer users the available apps to open files and mime types for new file creation. This information can be obtained from this endpoint.

**Endpoint**: specified in the capabilities in `apps_url`, currently `/app/list`

**Method**: HTTP GET

**Authentication**: None

**Request example**:

```bash
curl 'https://ocis.test/app/list'
```

**Response example**:

```json
{
  "mime-types": [
    {
      "mime_type": "application/pdf",
      "ext": "pdf",
      "app_providers": [
        {
          "name": "OnlyOffice",
          "icon": "https://www.pikpng.com/pngl/m/343-3435764_onlyoffice-desktop-editors-onlyoffice-logo-clipart.png"
        }
      ],
      "name": "PDF",
      "description": "PDF document"
    },
    {
      "mime_type": "application/vnd.oasis.opendocument.text",
      "ext": "odt",
      "app_providers": [
        {
          "name": "Collabora",
          "icon": "https://www.collaboraoffice.com/wp-content/uploads/2019/01/CP-icon.png"
        },
        {
          "name": "OnlyOffice",
          "icon": "https://www.pikpng.com/pngl/m/343-3435764_onlyoffice-desktop-editors-onlyoffice-logo-clipart.png"
        }
      ],
      "name": "OpenDocument",
      "icon": "https://some-website.test/opendocument-text-icon.png",
      "description": "OpenDocument text document",
      "allow_creation": true
    },
    {
      "mime_type": "text/markdown",
      "ext": "md",
      "app_providers": [
        {
          "name": "CodiMD",
          "icon": "https://avatars.githubusercontent.com/u/67865462?v=4"
        }
      ],
      "name": "Markdown file",
      "description": "Markdown file",
      "allow_creation": true
    },
    {
      "mime_type": "application/vnd.ms-word.document.macroenabled.12",
      "app_providers": [
        {
          "name": "Collabora",
          "icon": "https://www.collaboraoffice.com/wp-content/uploads/2019/01/CP-icon.png"
        },
        {
          "name": "OnlyOffice",
          "icon": "https://www.pikpng.com/pngl/m/343-3435764_onlyoffice-desktop-editors-onlyoffice-logo-clipart.png"
        }
      ]
    },
    {
      "mime_type": "application/vnd.ms-powerpoint.template.macroenabled.12",
      "app_providers": [
        {
          "name": "Collabora",
          "icon": "https://www.collaboraoffice.com/wp-content/uploads/2019/01/CP-icon.png"
        }
      ]
    }
  ]
}
```

### Open a file with the app provider

**Endpoint**: specified in the capabilities in `open_url`, currently `/app/open`

**Method**: HTTP POST

**Authentication** (one of them):

- `Authorization` header with OIDC Bearer token for authenticated users or basic auth credentials (if enabled in oCIS)
- `Public-Token` header with public link token for public links
- `X-Access-Token` header with a REVA token for authenticated users

**Query parameters**:

- `file_id` (mandatory): id of the file to be opened
- `app_name` (optional)
  - default (not given): default app for mime type
  - possible values depend on the app providers for a mimetype from the `/app/open` endpoint
- `view_mode` (optional)
  - default (not given): highest possible view mode, depending on the file permissions
  - possible values:
    - `write`: user can edit and download in the opening app
    - `read`: user can view and download from the opening app
    - `view`: user can view in the opening app (download is not possible)

**Request examples**:

```bash
curl -X POST 'https://ocis.test/app/open?file_id=ZmlsZTppZAo='

curl -X POST 'https://ocis.test/app/open?file_id=ZmlsZTppZAo=&app_name=Collabora'

curl -X POST 'https://ocis.test/app/open?file_id=ZmlsZTppZAo=&view_mode=read'

curl -X POST 'https://ocis.test/app/open?file_id=ZmlsZTppZAo=&app_name=Collabora&view_mode=write'
```

**Response examples**:

All apps are expected to be opened in an iframe and the response will give some parameters for that action.

There are apps, which need to be opened in the iframe with a form post. The form post must include all form parameters included in the response. For these apps the response will look like this:

```json
{
  "app_url": "https://.....",
  "method": "POST",
  "form_parameters": {
    "access_token": "eyJ0...",
    "access_token_ttl": "1634300912000",
    "arbitrary_param": "lorem-ipsum"
  }
}
```

There are apps, which need to be opened in the iframe with a GET request. The GET request must have set all headers included in the response. For these apps the response will look like this:

```json
{
  "app_url": "https://...",
  "method": "GET",
  "headers": {
    "access_token": "eyJ0e...",
    "access_token_ttl": "1634300912000",
    "arbitrary_header": "lorem-ipsum"
  }
}
```

**Example responses (error case)**:

- wrong `view_mode`

  ```json
  {
    "code": "SERVER_ERROR",
    "message": "Missing or invalid viewmode argument"
  }
  ```

- unknown `app_name`

  ```json
  {
    "code": "SERVER_ERROR",
    "message": "error searching for app provider"
  }
  ```

- wrong / invalid file id / unauthorized to open the file

  ```json
  {
    "code": "SERVER_ERROR",
    "message": "error statting file"
  }
  ```

## App drivers

App drivers represent apps, if the app is not able to register itself. Currently there is only the CS3org WOPI server app driver.

### CS3org WOPI server app driver

The CS3org WOPI server app driver is included in oCIS by default. It needs at least one WOPI compliant app (eg. Collabora, OnlyOffice or Microsoft Online Online Server) or a CS3org WOPI bridge supported app (CodiMD or Etherpad) and the CS3org WOPI server.

Here is a closer look at the configuration of the actual app provider in a docker-compose example (see also [full example](https://github.com/owncloud/ocis/blob/master/deployments/examples/ocis_wopi/docker-compose.yml)):

```yaml
services:
  ocis:
    image: owncloud/ocis:latest
    ...
    environment:
      ...
      STORAGE_GATEWAY_GRPC_ADDR: 0.0.0.0:9142 # make the REVA gateway accessible to the app drivers

  ocis-appdriver-collabora:
    image: owncloud/ocis:latest
    command: storage-app-provider server # start only the app driver
    environment:
      STORAGE_GATEWAY_ENDPOINT: ocis:9142 # oCIS gateway endpoint
      APP_PROVIDER_BASIC_EXTERNAL_ADDR: ocis-appdriver-collabora:9164 # how oCIS can reach this app driver
      OCIS_JWT_SECRET: ocis-jwt-secret
      APP_PROVIDER_DRIVER: wopi
      APP_PROVIDER_WOPI_DRIVER_APP_NAME: Collabora # will be used as name for this app
      APP_PROVIDER_WOPI_DRIVER_APP_ICON_URI: https://www.collaboraoffice.com/wp-content/uploads/2019/01/CP-icon.png # will be used as icon for this app
      APP_PROVIDER_WOPI_DRIVER_APP_URL: https://collabora.owncloud.test # endpoint of collabora
      APP_PROVIDER_WOPI_DRIVER_INSECURE: false
      APP_PROVIDER_WOPI_DRIVER_IOP_SECRET: wopi-iop-secret
      APP_PROVIDER_WOPI_DRIVER_WOPI_URL: https://wopiserver.owncloud.test # endpoint of the CS3org WOPI server
```
