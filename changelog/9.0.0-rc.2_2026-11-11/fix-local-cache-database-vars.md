Bugfix: The global variable OCIS_CACHE_DATABASE was replaced

Some services accepted only the global variable OCIS_CACHE_DATABASE for their
cache database name. The default values for each service is unique that is
why this has not been noticed until now.

We added a local variable for each service. The GATEWAY service now accepts
GATEWAY_PROVIDER_CACHE_DATABASE. The PROXY service now accepts
PROXY_OIDC_USERINFO_CACHE_DATABASE. The SETTINGS service now accepts
SETTINGS_CACHE_DATABASE. The STORAGE-SYSTEM service now accepts
STORAGE_SYSTEM_CACHE_DATABASE. The do not use the global variable OCIS_CACHE_DATABASE 
to avoid the cache collisions.

https://github.com/owncloud/ocis/pull/13159
