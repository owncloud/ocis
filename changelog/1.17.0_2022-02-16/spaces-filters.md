Enhancement: Add filter by driveType and id to /me/drives

We added two possible filter terms (driveType, id) to the /me/drives endpoint on the graph api. These can be used with the odata query parameter "$filter".
We only support the "eq" operator for now.

https://github.com/owncloud/ocis/pull/2946
