Bugfix: Use key to get specific trash item

The activitylog and clientlog services now only fetch the specific trash item instead of getting all items in trash and filtering them on their side. This reduces the load on the storage users service because it no longer has to assemble a full trash listing.

https://github.com/owncloud/ocis/pull/9879
