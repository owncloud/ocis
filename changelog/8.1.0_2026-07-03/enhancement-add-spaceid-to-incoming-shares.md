Enhancement: Add space ID to incoming shares

Added the `spaceId` to the incoming shares. This is aligning the graph API with the WebDAV API where the clients can use `spaceid` property.
This change allows clients to get the space ID directly instead of having to parse the resource ID.

https://github.com/owncloud/ocis/pull/12024
