Enhancement: Add sorting to GraphAPI users and groups

The GraphAPI endpoints for users and groups support ordering now.
User can be ordered by displayName, onPremisesSamAccountName and mail.
Groups can be ordered by displayName.

Example:
https://localhost:9200/graph/v1.0/groups?$orderby=displayName asc

https://github.com/owncloud/ocis/issues/3360
