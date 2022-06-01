## Scenarios from OCIS API tests that are expected to fail with OCIS storage
The expected failures in this file are from features in the owncloud/ocis repo.

#### [downloading an archive with invalid path returns HTTP/500](https://github.com/owncloud/ocis/issues/2768)
- [apiArchiver/downloadByPath.feature:69](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L69)

#### [Hardcoded call to /home/..., but /home no longer exists](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature)
- [apiArchiver/downloadByPath.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L26)
- [apiArchiver/downloadByPath.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L27)
- [apiArchiver/downloadByPath.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L44)
- [apiArchiver/downloadByPath.feature:45](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L45)
- [apiArchiver/downloadByPath.feature:48](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L48)
- [apiArchiver/downloadByPath.feature:69](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L69)
- [apiArchiver/downloadByPath.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L74)
- [apiArchiver/downloadByPath.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L132)
- [apiArchiver/downloadByPath.feature:133](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L133)

### Tries to download /Shares/ folder but it cannot be downloaded any more directly
- [apiArchiver/downloadById.feature:134](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L134)
- [apiArchiver/downloadById.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L135)

#### [Resharing is now allowed for viewers and editors](https://github.com/owncloud/ocis/issues/3828)
- [apiSpaces/shareSubItemOfSpace.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/shareSubItemOfSpace.feature#L89)
- [apiSpaces/shareSubItemOfSpace.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/shareSubItemOfSpace.feature#L90)
- [apiSpaces/shareSubItemOfSpace.feature:91](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/shareSubItemOfSpace.feature#L91)
- [apiSpaces/shareSubItemOfSpace.feature:92](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/shareSubItemOfSpace.feature#L92)
- [apiSpaces/shareSubItemOfSpaceViaPublicLink.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/shareSubItemOfSpaceViaPublicLink.feature#L89)
- [apiSpaces/shareSubItemOfSpaceViaPublicLink.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/shareSubItemOfSpaceViaPublicLink.feature#L90)
- [apiSpaces/shareSubItemOfSpaceViaPublicLink.feature:91](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/shareSubItemOfSpaceViaPublicLink.feature#L91)
- [apiSpaces/shareSubItemOfSpaceViaPublicLink.feature:92](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/shareSubItemOfSpaceViaPublicLink.feature#L92)

### Visibility of shares is still to discuss
- [apiSpaces/resharing.feature:37](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/resharing.feature#L37)
- [apiSpaces/resharing.feature:38](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/resharing.feature#L38)
- [apiSpaces/resharing.feature:39](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/resharing.feature#L39)
