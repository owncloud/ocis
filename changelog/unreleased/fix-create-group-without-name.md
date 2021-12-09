Bugfix: Disallow creation of a group with empty name via the OCS api

We've fixed the behavior for group creation on the OCS api, where it was
possible to create a group with an empty name. This was is not possible
on oC10 and is therefore also forbidden on oCIS to keep compatibility.
This PR forbids the creation and also ensures the correct status code
for both OCS v1 and OCS v2 apis.

https://github.com/owncloud/ocis/pull/2825
https://github.com/owncloud/ocis/issues/2823
