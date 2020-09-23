Bugfix: match the user response to the OC10 format

The user response contained the field `displayname` but for certain responses
the field `display-name` is expected. The field `display-name` was added and
now both fields are returned to the client.

<https://github.com/owncloud/product/issues/181>
<https://github.com/owncloud/ocis/ocs/pull/61>
