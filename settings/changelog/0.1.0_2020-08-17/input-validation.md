Change: Introduce input validation

We set up input validation, starting with enforcing alphanumeric identifier values and UUID
format on account uuids. As a result, traversal into parent folders is not possible anymore.
We also made sure that get and list requests are side effect free, i.e. not creating any folders.

<https://github.com/owncloud/ocis/settings/pull/22>
<https://github.com/owncloud/ocis/settings/issues/15>
<https://github.com/owncloud/ocis/settings/issues/16>
<https://github.com/owncloud/ocis/settings/issues/19>
