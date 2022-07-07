## Scenarios from ownCloud10 core API tests that are expected to fail with parallel deployment matrix: oc10-ocis-oc10

The expected failures in this file are from features in the owncloud/core repo.

### [System config cannot be set/unset in oCIS]()

- [apiCapabilities/capabilities.feature:8](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L8) Scenario: Check that the sharing API can be enabled
<!--  -->
- [apiCapabilities/capabilities.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L41)
- [apiCapabilities/capabilities.feature:85](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L85)
- [apiCapabilities/capabilities.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L100)
- [apiCapabilities/capabilities.feature:116](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L116)
- [apiCapabilities/capabilities.feature:127](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L127)
- [apiCapabilities/capabilities.feature:139](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L139)
- [apiCapabilities/capabilities.feature:282](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L282)
<!-- Then step must be in oCIS to pass -->
- [apiCapabilities/capabilities.feature:959](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L959)
- [apiMain/status.feature:5](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/status.feature#L5)

#### [Ability to return error messages in Webdav response bodies](https://github.com/owncloud/ocis/issues/1293)

- [apiAuthOcs/ocsPOSTAuth.feature:10](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPOSTAuth.feature#L10) Scenario: send POST requests to OCS endpoints as normal user with wrong password
- [apiAuthOcs/ocsPUTAuth.feature:10](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPUTAuth.feature#L10) Scenario: send PUT request to OCS endpoints as admin with wrong password

### [PROPFIND for checksums returns multiple duplicate checksums (using oCIS selector)](https://github.com/owncloud/ocis/issues/4092)

- [apiMain/checksums.feature:29](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/checksums.feature#L29)
- [apiMain/checksums.feature:30](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/checksums.feature#L30)
- [apiMain/checksums.feature:46](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/checksums.feature#L46)
- [apiMain/checksums.feature:47](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/checksums.feature#L47)

### [Differences in path property while listing pending shares](https://github.com/owncloud/ocis/issues/4035)

- [apiShareCreateSpecialToShares1/createShareUniqueReceivedNames.feature:15](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareUniqueReceivedNames.feature#L15)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:37](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L37)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:38](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L38)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:67](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L67)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:68](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L68)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:235](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L235)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:236](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L236)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:639](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L639)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:640](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L640)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:67](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L67)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:120](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L120)
- [apiShareManagementToShares/acceptShares.feature:65](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L65)
- [apiShareManagementToShares/acceptShares.feature:101](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L101)
- [apiShareManagementToShares/acceptShares.feature:224](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L224)
- [apiShareManagementToShares/acceptShares.feature:252](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L252)
- [apiShareManagementToShares/mergeShare.feature:16](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/mergeShare.feature#L16)
- [apiShareOperationsToShares1/accessToShare.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/accessToShare.feature#L41)
- [apiShareOperationsToShares1/accessToShare.feature:42](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/accessToShare.feature#L42)
- [apiShareOperationsToShares1/accessToShare.feature:74](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/accessToShare.feature#L74)
- [apiShareOperationsToShares1/accessToShare.feature:75](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/accessToShare.feature#L75)
- [apiShareOperationsToShares1/changingFilesShare.feature:25](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/changingFilesShare.feature#L25)
- [apiShareOperationsToShares1/changingFilesShare.feature:26](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/changingFilesShare.feature#L26)
- [apiShareOperationsToShares1/gettingShares.feature:85](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L85)
- [apiShareOperationsToShares1/gettingShares.feature:86](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L86)
- [apiShareOperationsToShares1/gettingShares.feature:106](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L106)
- [apiShareOperationsToShares1/gettingShares.feature:107](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L107)
- [apiShareOperationsToShares1/gettingShares.feature:138](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L138)
- [apiShareOperationsToShares1/gettingShares.feature:139](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L139)
- [apiShareOperationsToShares1/gettingShares.feature:170](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L170)
- [apiShareOperationsToShares1/gettingShares.feature:171](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L171)
- [apiSharePublicLink1/createPublicLinkShare.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/createPublicLinkShare.feature#L41)
- [apiSharePublicLink1/createPublicLinkShare.feature:42](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/createPublicLinkShare.feature#L42)
- [apiSharePublicLink1/createPublicLinkShare.feature:105](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/createPublicLinkShare.feature#L105)
- [apiSharePublicLink1/createPublicLinkShare.feature:106](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/createPublicLinkShare.feature#L106)
- [apiSharePublicLink1/createPublicLinkShare.feature:241](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/createPublicLinkShare.feature#L241)
- [apiSharePublicLink1/createPublicLinkShare.feature:242](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/createPublicLinkShare.feature#L242)

### [Cannot GET share information with share id (using ocis selector)](https://github.com/owncloud/ocis/issues/4101)

- [apiSharePublicLink2/multilinkSharing.feature:43](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink2/multilinkSharing.feature#L43)
- [apiSharePublicLink2/multilinkSharing.feature:44](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink2/multilinkSharing.feature#L44)
- [apiSharePublicLink3/updatePublicLinkShare.feature:78](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/updatePublicLinkShare.feature#L78)
- [apiSharePublicLink3/updatePublicLinkShare.feature:79](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/updatePublicLinkShare.feature#L79)

### [Shares are auto-accepted with ocis selector]()
