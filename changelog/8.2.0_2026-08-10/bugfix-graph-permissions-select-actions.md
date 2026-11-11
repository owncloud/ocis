Bugfix: Honor $select=actions on the drive item permissions endpoint

Listing the permissions of a drive item with
`$select=@libre.graph.permissions.actions.allowedValues` still returned the
roles allowedValues as well. The handler only had a projection for the roles
selection (`@libre.graph.permissions.roles.allowedValues`), which drops the
actions, but no symmetric handling for an actions-only selection, so the roles
were always included. The actions selection now drops the roles allowedValues
(and, like the roles selection, skips the share lookup since only the allowed
values are requested).

https://github.com/owncloud/ocis/issues/10816
https://github.com/owncloud/ocis/pull/12419
