



oCIS is all about files. But most of the time you wan to do something with files. Therefore oCIS has an concept about apps, that can handle specific file types, so called mime types.

App registry
The app registry is the single point where all apps register itself and their supported mime types.

Mime type configuration / allow list
Administrators need to add supported mime types to the mime type configuration, which works like an allow list. Only mime types on this list will be offered by the app registry even if app providers register additional mime types.

In order to modify the mime type config you need to set STORAGE_APP_REGISTRY_MIMETYPES_JSON=.../mimetypes.json to a valid JSON file with a content like this:

{
   "application/vnd.oasis.opendocument.text":{
      "extension":"odt",
      "name":"OpenDocument",
      "description":"OpenDocument text document",
      "icon":"",
      "default_app":"Collabora"
   },
   "application/vnd.oasis.opendocument.spreadsheet":{
      "extension":"ods",
      "name":"OpenSpreadsheet",
      "description":"OpenDocument spreadsheet document",
      "icon":"",
      "default_app":"Collabora"
   }
}
Please add all mime types you would like use with apps. You also can configure, which application should be treated as a default app for a certain mime type by setting the app provider name in default_app.

Listing available mime types / apps
/app/list

{
  "mime-types": [
    {
      "mime_type": "application/vnd.oasis.opendocument.spreadsheet",
      "ext": "ods",
      "app_providers": [
        {
          "name": "Collabora",
          "icon": "https://www.collaboraoffice.com/wp-content/uploads/2019/01/CP-icon.png"
        }
      ],
      "name": "OpenSpreadsheet",
      "description": "OpenDocument spreadsheet document"
    },
    {
      "mime_type": "application/vnd.oasis.opendocument.presentation",
      "ext": "odp",
      "app_providers": [
        {
          "name": "Collabora",
          "icon": "https://www.collaboraoffice.com/wp-content/uploads/2019/01/CP-icon.png"
        }
      ],
      "name": "OpenPresentation",
      "description": "OpenDocument presentation document"
    }
/app/open

App provider / drivers
WOPI app provider with CS3org WOPI server
You can run an app provider next to your regular oCIS (docker-compose example). Aditionally you need a CS3 WOPI server and Collabora Online instances running. Both can be found in our WOPI deployment example.

Here is a closer look at the configuration of the actual app provider:

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
