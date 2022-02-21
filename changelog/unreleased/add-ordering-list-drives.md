Enhancement: Add sorting to list Spaces

We added the OData query param "orderBy" for listing spaces. We can now order by Space Name and LastModifiedDateTime.

Example 1: https://localhost:9200/graph/v1.0/me/drives/?$orderby=lastModifiedDateTime desc
Example 2: https://localhost:9200/graph/v1.0/me/drives/?$orderby=name asc

https://github.com/owncloud/ocis/issues/3200
https://github.com/owncloud/ocis/pull/3201
