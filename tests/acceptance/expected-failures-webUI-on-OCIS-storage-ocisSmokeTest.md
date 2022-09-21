## Scenarios from web tests that are expected to fail on OCIS with OCIS storage
The expected failures in this file are from features in the owncloud/web repo.

Lines that contain a format like "[someSuite.someFeature.feature:n](https://github.com/owncloud/web/path/to/feature)"
are lines that document a specific expected failure. Follow that with a URL to the line in the feature file in GitHub.
Please follow this format for the actual expected failures.

Level-3 headings should be used for the references to the relevant issues. Include the issue title with a link to the issue in GitHub.

Other free text and Markdown formatting can be used elsewhere in the document if needed. But if you want to explain something about the issue, then please post that in the issue itself.

Only the web scenarios tagged ocisSmokeTest are run by default in OCIS CI. This file lists the expected-failures of those ocisSmokeTest scenarios.

### [Copy private link option not available](https://github.com/owncloud/ocis/issues/1409)
- [webUIPrivateLinks/accessingPrivateLinks.feature:9](https://github.com/owncloud/web/blob/master/tests/acceptance/features/webUIPrivateLinks/accessingPrivateLinks.feature#L9)
- [webUIPrivateLinks/accessingPrivateLinks.feature:17](https://github.com/owncloud/web/blob/master/tests/acceptance/features/webUIPrivateLinks/accessingPrivateLinks.feature#L17)

### [impossible to navigate into a folder in the trashbin](https://github.com/owncloud/web/issues/1725)
- [webUITrashbinDelete/trashbinDelete.feature:29](https://github.com/owncloud/web/blob/master/tests/acceptance/features/webUITrashbinDelete/trashbinDelete.feature#L29)
