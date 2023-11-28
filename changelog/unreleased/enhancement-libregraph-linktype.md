Enhancement: New public sharing via libregraph link types

We added libregraph link types to the graph service as a new set of api calls to replace the legacy OCS API. We added the endpoint to create links https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink
and the linktype listing to the "SharedByMe" call.

We also added a config to match legacy public links created by the OCS API to the new libregraph link types. This config switch is recommended for instances which were already existing before Infinite Scale 5.0.0 was released.

https://github.com/owncloud/ocis/pull/7834
https://github.com/owncloud/ocis/pull/7743
https://github.com/owncloud/ocis/issues/6993
