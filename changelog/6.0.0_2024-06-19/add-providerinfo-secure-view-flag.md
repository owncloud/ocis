Enhancement: add secureview flag when listing apps via http

To allow clients to see which application supports secure view, we add a flag to the http response when the app service name matches a configured secure view app provider. The app can be configured by setting `FRONTEND_APP_HANDLER_SECURE_VIEW_APP_ADDR` to the address of the registered CS3 app provider.

https://github.com/owncloud/ocis/pull/9289
https://github.com/owncloud/ocis/pull/9280
https://github.com/owncloud/ocis/pull/9277
