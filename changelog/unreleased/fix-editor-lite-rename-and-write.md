Bugfix: Allow renaming, moving and editing files with the editor-lite role

We've fixed the editor-lite role ("Can edit" in the web UI). Sharees can now
rename a file inside a shared folder, move it and change its contents. Before,
they could only upload new files, and the web UI offered neither the rename nor
the overwrite action.

The role has always granted Move, but the persisted ACE format had no flag of
its own for it: Move shared the "w" flag with InitiateFileUpload and was only
recovered on read by assuming a grant may move whenever it may write, download
and delete. The editor-lite role has no delete, so its Move was dropped as soon
as the grant was written to disk, `Decomposedfs.Move` then refused the rename
and propfind reported the grant as a create-only uploader. Move is now
persisted under its own "m" flag in `pkg/storage/utils/ace/ace.go`; grants
written before that flag existed carry no "m" and keep using the old inference.
On top of that, `RoleFromResourcePermissions` in `pkg/conversions/role.go` did
not treat Move as write, so the WebDAV permissions string lacked "NV" (rename)
and "W" (overwrite) even for a grant that had kept its Move. Move now implies
write for grants that do not carry delete, which leaves the OCS permissions of
the deletable roles unchanged.

Note: the effective permission set of the editor-lite role changed. The CS3
resource permissions returned by `NewEditorLiteRole` are unchanged - Stat,
GetPath, ListContainer, InitiateFileDownload, InitiateFileUpload,
CreateContainer and Move, still no Delete and no ListFileVersions - but a grant
created from it is now stored as "txrwma" instead of "txrwa", reports the OCS
permissions 6 (create + write) instead of 4 (create), and reports the WebDAV
permissions "SNVW" on a shared file and "SNVCK" on a shared folder instead of
"S" and "SCK".

The same "m" flag also restores Move for any other non-deletable grant that
carries it. In particular the Uploader share role is affected: like editor-lite
it is a create grant without delete, so its Move used to be dropped on
write-back and is now preserved. This is why the previously expected-to-fail
"sharee moves a file within a shared folder" scenarios for the Uploader role now
pass, and the matching entries were removed from
`tests/acceptance/expected-failures-API-on-OCIS-storage.md`.

This fix DOES NOT include a migration of stored grants. Shares created with
that role before this update keep their old grant on disk, which carries no "m"
flag. For them, Move stays unset and the WebDAV permissions string stays
unchanged, while the share record itself has always kept Move and therefore
already reports the new OCS permissions. Only shares created with that role
AFTER THE UPDATE get the new permission set. Re-creating an affected share, or
updating its role, re-writes the grant and picks up the fix.

https://github.com/owncloud/ocis/pull/12721
https://github.com/owncloud/reva/pull/689
https://github.com/owncloud/ocis/issues/11977
