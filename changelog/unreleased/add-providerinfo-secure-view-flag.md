Enhancement: add secureview flag when listing apps via http

To allow clients to see which application supports secure view we add a flag to the http response when the app name matches a configured secure view app. The app can be configured by setting `FRONTEND_APP_HANDLER_SECURE_VIEW_APP` to the name of the app registered as a CS3 app provider.

https://github.com/owncloud/ocis/pull/9277
