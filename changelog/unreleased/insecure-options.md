Change: Make all insecure options configurable and change the default to false

We had several hard-coded 'insecure' flags. These options are now configurable and default to false. Also we changed all other 'insecure' flags with a previous default of true to false.

In development environments using self signed certs (the default) you now need to set these flags:

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

As an alternative you also can set a single flag, which configures all options together:

```
OCIS_INSECURE=true
```

https://github.com/owncloud/ocis/issues/2700
https://github.com/owncloud/ocis/pull/2745
