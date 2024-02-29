## Scenarios from web tests that are expected to fail on OCIS with OCIS storage

The expected failures in this file are from features in the owncloud/web repo.

Lines that contain a format like "[someSuite.someFeature.feature:n](https://github.com/owncloud/web/path/to/feature)"
are lines that document a specific expected failure. Follow that with a URL to the line in the feature file in GitHub.
Please follow this format for the actual expected failures.

Level-3 headings should be used for the references to the relevant issues. Include the issue title with a link to the issue in GitHub.

Other free text and Markdown formatting can be used elsewhere in the document if needed. But if you want to explain something about the issue, then please post that in the issue itself.

### [Exit page re-appears in loop when logged in user is deleted](https://github.com/owncloud/web/issues/4677)

- [webUILogin/openidLogin.feature:50](https://github.com/owncloud/web/blob/master/tests/acceptance/features/webUILogin/openidLogin.feature#L50)
- [webUILogin/openidLogin.feature:60](https://github.com/owncloud/web/blob/master/tests/acceptance/features/webUILogin/openidLogin.feature#L60)

### [restoring a file deleted from a received shared folder is not possible](https://github.com/owncloud/ocis/issues/1124)

- [webUITrashbinRestore/trashbinRestore.feature:176](https://github.com/owncloud/web/blob/master/tests/acceptance/features/webUITrashbinRestore/trashbinRestore.feature#L176)

### [empty subfolder inside a folder to be uploaded is not created on the server](https://github.com/owncloud/web/issues/6348)

- [webUIUpload/upload.feature:36](https://github.com/owncloud/web/blob/master/tests/acceptance/features/webUIUpload/upload.feature#L36)