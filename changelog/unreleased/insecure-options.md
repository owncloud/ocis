Enhancement: Make insecure options configurable

We had several hard-coded 'insecure' flags. These options are now configurable and default to false. Also we changed all other 'insecure' flags with a previous default of true to false. In development environments using self signed certs (the default) you need to set these flags:

```
PROXY_OIDC_INSECURE=true
STORAGE_FRONTEND_APPPROVIDER_INSECURE=true
STORAGE_FRONTEND_ARCHIVER_INSECURE=true
STORAGE_FRONTEND_OCDAV_INSECURE=true
STORAGE_HOME_DATAPROVIDER_INSECURE=true
STORAGE_METADATA_DATAPROVIDER_INSECURE=true
STORAGE_OIDC_INSECURE=true
STORAGE_USERS_DATAPROVIDER_INSECURE=true
THUMBNAILS_CS3SOURCE_INSECURE=true
THUMBNAILS_WEBDAVSOURCE_INSECURE=true
```

https://github.com/owncloud/ocis/issues/2700
https://github.com/owncloud/ocis/pull/2745
