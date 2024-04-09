Bugfix: Validate conditions for sharing roles by resource type

We improved the validation of the allowed sharing roles for specific resource type
for various sharing related graph API endpoints. This allows e.g. the web client to
restrict the sharing roles presented to the user based on the type of the resource
that is being shared.

https://github.com/owncloud/ocis/pull/8815
https://github.com/owncloud/ocis/issues/8331
