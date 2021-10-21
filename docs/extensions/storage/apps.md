---
title: "Apps"
date: 2018-05-02T00:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: apps.md
---

oCIS is all about files. But most of the time you want to do something with files. Therefore oCIS has an concept about apps, that can handle specific file types, so called mime types.

## App registry

The app registry is the single point where all apps register itself and their supported mime types.

### Mime type configuration / creation allow list

The apps will register their supported mime types automatically, so that users can open supported files with them.

Administrators can set default applications on a per mimetype basis and also allow the creation of new files for certain mime types. This per mime type configuration also features a description, file extension option and an icon.

In order to modify the mime type config you need to set `STORAGE_APP_REGISTRY_MIMETYPES_JSON=.../mimetypes.json` to a valid JSON file with a content like this:

```json
[
  {
    "mime_type": "application/vnd.oasis.opendocument.text",
    "extension": "odt",
    "name": "OpenDocument",
    "description": "OpenDocument text document",
    "icon": "",
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


### Listing available mime types / apps

#### /app/list

Method: HTTP GET

Authentication: None

Result:

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

#### /app/open

Method: HTTP POST

Authentication (one of them):

- `Authorization` header with OIDC Bearer token for authenticated users or basic auth credentials (if enabled in oCIS)
- `Public-Token` header with public link token for public links
- `X-Access-Token` header with a REVA token for authenticated users

Query parameters:

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

Examples:


``` bash
curl 'https://ocis.test/app/open?file_id=ZmlsZTppZAo='

curl 'https://ocis.test/app/open?file_id=ZmlsZTppZAo=&app_name=Collabora'

curl 'https://ocis.test/app/open?file_id=ZmlsZTppZAo=&view_mode=read'

curl 'https://ocis.test/app/open?file_id=ZmlsZTppZAo=&app_name=Collabora&view_mode=write'
```


Response:

Apps are expected to be opened in Iframes and the response will give some parameters for that action.

Some apps expect to be opened in the Iframe with a form post. The response will look like this:

``` json
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


Some apps expect to be opened in the Iframe with a GET request and need additional headers to be set:

``` json
{
  "app_url": "https://...",
  "method": "GET",
  "headers": {
    "access_token": "eyJ0e...",
    "access_token_ttl": "1634300912000",
    "arbitrary_header": "lorem-ipsum",
  }
}

```

If opening an app fails, you may encounter one of the following errors:

- wrong `view_mode`
``` json
{
  "code": "SERVER_ERROR",
  "message": "Missing or invalid viewmode argument"
}
```

- unknown `app_name`
``` json
{
  "code": "SERVER_ERROR",
  "message": "error searching for app provider"
}
```

- wrong / invalid file id / unauthorized to open the file
``` json
{
  "code": "SERVER_ERROR",
  "message": "error statting file"
}
```

## App provider / drivers

WOPI app provider with CS3org WOPI server
You can run an app provider next to your regular oCIS (docker-compose example). Aditionally you need a CS3 WOPI server and Collabora Online instances running. Both can be found in our WOPI deployment example.

Here is a closer look at the configuration of the actual app provider in a docker-compose example:

```yaml
services:
  ocis: ...

  ocis-appdriver-collabora:
    image: owncloud/ocis:latest
    command: storage-app-provider server
    environment:
      STORAGE_GATEWAY_ENDPOINT: ocis:9142
      APP_PROVIDER_BASIC_EXTERNAL_ADDR: ocis-appdriver-collabora:9164
      OCIS_JWT_SECRET: ocis-jwt-secret
      APP_PROVIDER_DRIVER: wopi
      APP_PROVIDER_WOPI_DRIVER_APP_NAME: Collabora
      APP_PROVIDER_WOPI_DRIVER_APP_ICON_URI: https://www.collaboraoffice.com/wp-content/uploads/2019/01/CP-icon.png
      APP_PROVIDER_WOPI_DRIVER_APP_URL: https://collabora.owncloud.test
      APP_PROVIDER_WOPI_DRIVER_INSECURE: false
      APP_PROVIDER_WOPI_DRIVER_IOP_SECRET: wopi-iop-secret
      APP_PROVIDER_WOPI_DRIVER_WOPI_URL: https://wopiserver.owncloud.test
```
