Enhancement: The graph endpoints for listing permission works for spaces now

We enhanced the 'graph/v1beta1/drives/{{driveid}}/items/{{itemid}}/permissions' endpoint
to list permission of the space when the 'itemid' refers to a space root.

https://github.com/owncloud/ocis/pull/8642
https://github.com/owncloud/ocis/issues/8352
