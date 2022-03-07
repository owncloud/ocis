Change: Unify file IDs

We changed the file IDs to be consistent across all our APIs (WebDAV, LibreGraph, OCS). We removed the base64 encoding. Now they are formatted like <storageID>!<opaqueID>. They are using a reserved character ``!`` as a URL safe separator.

https://github.com/owncloud/ocis/pull/3185
