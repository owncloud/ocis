Change: storage providers now default to exposing data servers

The flags that let reva storage providers announce that they expose a data server now defaults to true:

`REVA_STORAGE_HOME_EXPOSE_DATA_SERVER=1`
`REVA_STORAGE_OC_EXPOSE_DATA_SERVER=1`

https://github.com/owncloud/ocis-reva/issues/89