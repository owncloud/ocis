## Scenarios from ownCloud10 core API tests that are expected to fail with OCIS storage while running with the Graph API

The expected failures in this file are from features in the owncloud/core repo.

### File

Basic file management like up and download, move, copy, properties, trash, versions and chunking.

#### [Getting information about a folder overwritten by a file gives 500 error instead of 404](https://github.com/owncloud/ocis/issues/1239)

- [apiWebdavProperties1/copyFile.feature:274](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L274)
- [apiWebdavProperties1/copyFile.feature:275](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L275)
- [apiWebdavProperties1/copyFile.feature:292](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L292)
- [apiWebdavProperties1/copyFile.feature:293](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L293)

#### [Custom dav properties with namespaces are rendered incorrectly](https://github.com/owncloud/ocis/issues/2140)

_ocdav: double-check the webdav property parsing when custom namespaces are used_

- [apiWebdavProperties1/setFileProperties.feature:37](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/setFileProperties.feature#L37)
- [apiWebdavProperties1/setFileProperties.feature:38](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/setFileProperties.feature#L38)
- [apiWebdavProperties1/setFileProperties.feature:43](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/setFileProperties.feature#L43)
- [apiWebdavProperties1/setFileProperties.feature:78](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/setFileProperties.feature#L78)
- [apiWebdavProperties1/setFileProperties.feature:79](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/setFileProperties.feature#L79)
- [apiWebdavProperties1/setFileProperties.feature:84](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/setFileProperties.feature#L84)

#### [Cannot set custom webDav properties](https://github.com/owncloud/product/issues/264)

- [apiWebdavProperties2/getFileProperties.feature:360](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties2/getFileProperties.feature#L360)
- [apiWebdavProperties2/getFileProperties.feature:365](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties2/getFileProperties.feature#L365)
- [apiWebdavProperties2/getFileProperties.feature:370](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties2/getFileProperties.feature#L370)
- [apiWebdavProperties2/getFileProperties.feature:400](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties2/getFileProperties.feature#L400)
- [apiWebdavProperties2/getFileProperties.feature:405](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties2/getFileProperties.feature#L405)
- [apiWebdavProperties2/getFileProperties.feature:410](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties2/getFileProperties.feature#L410)

#### [file versions do not report the version author](https://github.com/owncloud/ocis/issues/2914)

- [apiVersions/fileVersionAuthor.feature:14](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L14)
- [apiVersions/fileVersionAuthor.feature:37](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L37)
- [apiVersions/fileVersionAuthor.feature:58](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L58)
- [apiVersions/fileVersionAuthor.feature:78](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L78)
- [apiVersions/fileVersionAuthor.feature:104](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L104)
- [apiVersions/fileVersionAuthor.feature:129](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L129)
- [apiVersions/fileVersionAuthor.feature:154](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L154)
- [apiVersions/fileVersionAuthor.feature:180](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L180)
- [apiVersions/fileVersionAuthor.feature:223](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionAuthor.feature#L223)

### Sync

Synchronization features like etag propagation, setting mtime and locking files

#### [Uploading an old method chunked file with checksum should fail using new DAV path](https://github.com/owncloud/ocis/issues/2323)

- [apiMain/checksums.feature:369](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/checksums.feature#L369)
- [apiMain/checksums.feature:374](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/checksums.feature#L374)

#### [Webdav LOCK operations](https://github.com/owncloud/ocis/issues/1284)

- [apiWebdavLocks/exclusiveLocks.feature:124](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L124)
- [apiWebdavLocks/exclusiveLocks.feature:125](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L125)
- [apiWebdavLocks/exclusiveLocks.feature:126](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L126)
- [apiWebdavLocks/exclusiveLocks.feature:127](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L127)
- [apiWebdavLocks/exclusiveLocks.feature:132](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L132)
- [apiWebdavLocks/exclusiveLocks.feature:133](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L133)
- [apiWebdavLocks/exclusiveLocks.feature:151](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L151)
- [apiWebdavLocks/exclusiveLocks.feature:152](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L152)
- [apiWebdavLocks/exclusiveLocks.feature:153](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L153)
- [apiWebdavLocks/exclusiveLocks.feature:154](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L154)
- [apiWebdavLocks/exclusiveLocks.feature:159](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L159)
- [apiWebdavLocks/exclusiveLocks.feature:160](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L160)
- [apiWebdavLocks/exclusiveLocks.feature:178](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L178)
- [apiWebdavLocks/exclusiveLocks.feature:179](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L179)
- [apiWebdavLocks/exclusiveLocks.feature:180](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L180)
- [apiWebdavLocks/exclusiveLocks.feature:181](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L181)
- [apiWebdavLocks/exclusiveLocks.feature:186](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L186)
- [apiWebdavLocks/exclusiveLocks.feature:187](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/exclusiveLocks.feature#L187)
- [apiWebdavLocks/requestsWithToken.feature:130](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/requestsWithToken.feature#L130)
- [apiWebdavLocks/requestsWithToken.feature:131](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/requestsWithToken.feature#L131)
- [apiWebdavLocks/requestsWithToken.feature:136](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks/requestsWithToken.feature#L136)
- [apiWebdavLocks3/independentLocks.feature:65](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L65)
- [apiWebdavLocks3/independentLocks.feature:66](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L66)
- [apiWebdavLocks3/independentLocks.feature:67](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L67)
- [apiWebdavLocks3/independentLocks.feature:68](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L68)
- [apiWebdavLocks3/independentLocks.feature:73](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L73)
- [apiWebdavLocks3/independentLocks.feature:74](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L74)
- [apiWebdavLocks3/independentLocks.feature:93](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L93)
- [apiWebdavLocks3/independentLocks.feature:94](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L94)
- [apiWebdavLocks3/independentLocks.feature:95](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L95)
- [apiWebdavLocks3/independentLocks.feature:96](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L96)
- [apiWebdavLocks3/independentLocks.feature:97](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L97)
- [apiWebdavLocks3/independentLocks.feature:98](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L98)
- [apiWebdavLocks3/independentLocks.feature:99](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L99)
- [apiWebdavLocks3/independentLocks.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L100)
- [apiWebdavLocks3/independentLocks.feature:105](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L105)
- [apiWebdavLocks3/independentLocks.feature:106](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L106)
- [apiWebdavLocks3/independentLocks.feature:107](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L107)
- [apiWebdavLocks3/independentLocks.feature:108](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocks.feature#L108)
- [apiWebdavLocks3/independentLocksShareToShares.feature:75](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L75)
- [apiWebdavLocks3/independentLocksShareToShares.feature:76](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L76)
- [apiWebdavLocks3/independentLocksShareToShares.feature:77](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L77)
- [apiWebdavLocks3/independentLocksShareToShares.feature:78](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L78)
- [apiWebdavLocks3/independentLocksShareToShares.feature:83](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L83)
- [apiWebdavLocks3/independentLocksShareToShares.feature:84](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L84)
- [apiWebdavLocks3/independentLocksShareToShares.feature:104](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L104)
- [apiWebdavLocks3/independentLocksShareToShares.feature:105](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L105)
- [apiWebdavLocks3/independentLocksShareToShares.feature:106](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L106)
- [apiWebdavLocks3/independentLocksShareToShares.feature:107](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L107)
- [apiWebdavLocks3/independentLocksShareToShares.feature:112](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L112)
- [apiWebdavLocks3/independentLocksShareToShares.feature:113](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocks3/independentLocksShareToShares.feature#L113)
- [apiWebdavLocksUnlock/unlock.feature:40](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L40)
- [apiWebdavLocksUnlock/unlock.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L41)
- [apiWebdavLocksUnlock/unlock.feature:46](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L46)
- [apiWebdavLocksUnlock/unlock.feature:79](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L79)
- [apiWebdavLocksUnlock/unlock.feature:80](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L80)
- [apiWebdavLocksUnlock/unlock.feature:129](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L129)
- [apiWebdavLocksUnlock/unlock.feature:130](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L130)
- [apiWebdavLocksUnlock/unlock.feature:131](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L131)
- [apiWebdavLocksUnlock/unlock.feature:132](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L132)
- [apiWebdavLocksUnlock/unlock.feature:137](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L137)
- [apiWebdavLocksUnlock/unlock.feature:138](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlock.feature#L138)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:28](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L28)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:29](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L29)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:30](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L30)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:31](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L31)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:44](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L44)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:45](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L45)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:60](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L60)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:61](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L61)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:62](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L62)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:63](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L63)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:68](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L68)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:69](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L69)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:109](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L109)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:110](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L110)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:111](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L111)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:112](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L112)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:125](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L125)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:126](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L126)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:142](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L142)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:143](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L143)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:144](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L144)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:145](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L145)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:158](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L158)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:159](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L159)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:174](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L174)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:175](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L175)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:176](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L176)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:177](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L177)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:182](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L182)
- [apiWebdavLocksUnlock/unlockSharingToShares.feature:183](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavLocksUnlock/unlockSharingToShares.feature#L183)

### Share

File and sync features in a shared scenario

#### [ocs sharing api always returns an empty exact list while searching for a sharee](https://github.com/owncloud/ocis/issues/4265)

- [apiSharees/sharees.feature:350](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L350)
- [apiSharees/sharees.feature:351](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L351)
- [apiSharees/sharees.feature:370](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L370)
- [apiSharees/sharees.feature:371](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L371)
- [apiSharees/sharees.feature:390](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L390)
- [apiSharees/sharees.feature:391](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L391)
- [apiSharees/sharees.feature:410](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L410)
- [apiSharees/sharees.feature:411](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L411)
- [apiSharees/sharees.feature:430](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L430)
- [apiSharees/sharees.feature:431](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharees/sharees.feature#L431)


#### [accepting matching name shared resources from different users/groups sets no serial identifiers on the resource name for the receiver](https://github.com/owncloud/ocis/issues/4289)

- [apiShareManagementToShares/acceptShares.feature:365](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L365)
- [apiShareManagementToShares/acceptShares.feature:401](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L401)
- [apiShareManagementToShares/acceptShares.feature:303](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L303)
- [apiShareManagementToShares/acceptShares.feature:590](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L590)
- [apiShareManagementToShares/acceptShares.feature:655](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L655)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:162](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L162)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:163](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L163)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:202](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L202)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:203](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L203)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:45](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L45)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:46](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L46)

#### [sharing the shares folder to users exits with different status code than in oc10 backend](https://github.com/owncloud/ocis/issues/2215)

- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:765](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L765)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:766](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L766)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:784](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L784)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:785](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L785)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:800](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L800)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:801](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L801)

#### [file_target of an auto-renamed file is not correct directly after sharing](https://github.com/owncloud/core/issues/32322)

- [apiShareManagementToShares/mergeShare.feature:121](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/mergeShare.feature#L121)

#### [File deletion using dav gives unique string in filename in the trashbin](https://github.com/owncloud/product/issues/178)

- [apiShareManagementBasicToShares/deleteShareFromShares.feature:69](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L69)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:83](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L83)

cannot share a folder with create permission

#### [Resource with share permission create is readable for sharee](https://github.com/owncloud/ocis/issues/4524)

- [apiShareManagementBasicToShares/deleteShareFromShares.feature:142](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L142)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:155](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L155)


#### [OCS error message for attempting to access share via share id as an unauthorized user is not informative](https://github.com/owncloud/ocis/issues/1233)

- [apiShareOperationsToShares1/gettingShares.feature:184](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L184)
- [apiShareOperationsToShares1/gettingShares.feature:185](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L185)

#### [Listing shares via ocs API does not show path for parent folders](https://github.com/owncloud/ocis/issues/1231)

- [apiShareOperationsToShares1/gettingShares.feature:222](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L222)
- [apiShareOperationsToShares1/gettingShares.feature:223](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/gettingShares.feature#L223)

#### [Public link enforce permissions](https://github.com/owncloud/ocis/issues/1269)

- [apiSharePublicLink1/accessToPublicLinkShare.feature:10](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/accessToPublicLinkShare.feature#L10)
- [apiSharePublicLink1/accessToPublicLinkShare.feature:21](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/accessToPublicLinkShare.feature#L21)
- [apiSharePublicLink1/accessToPublicLinkShare.feature:31](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/accessToPublicLinkShare.feature#L31)
- [apiSharePublicLink1/accessToPublicLinkShare.feature:47](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/accessToPublicLinkShare.feature#L47)
- [apiSharePublicLink1/createPublicLinkShare.feature:581](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/createPublicLinkShare.feature#L581)
- [apiSharePublicLink1/createPublicLinkShare.feature:602](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/createPublicLinkShare.feature#L602)

#### [download previews of other users file](https://github.com/owncloud/ocis/issues/2071)

- [apiWebdavPreviews/previews.feature:102](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L102)

#### [different error message detail for previews of folder](https://github.com/owncloud/ocis/issues/2064)

- [apiWebdavPreviews/previews.feature:111](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L111)

#### [Requesting a file preview when it is disabled by the administrator](https://github.com/owncloud/ocis/issues/192)

- [apiWebdavPreviews/previews.feature:126](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L126)

#### [Cannot set/unset maximum and minimum preview dimensions](https://github.com/owncloud/ocis/issues/2070)

- [apiWebdavPreviews/previews.feature:134](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L134)
- [apiWebdavPreviews/previews.feature:162](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L162)
- [apiWebdavPreviews/previews.feature:163](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L163)
- [apiWebdavPreviews/previews.feature:164](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L164)
- [apiWebdavPreviews/previews.feature:176](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L176)
- [apiWebdavPreviews/previews.feature:177](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L177)

#### [creating public links with permissions fails](https://github.com/owncloud/product/issues/252)

- [apiSharePublicLink1/changingPublicLinkShare.feature:30](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/changingPublicLinkShare.feature#L30)
- [apiSharePublicLink1/changingPublicLinkShare.feature:51](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/changingPublicLinkShare.feature#L51)
- [apiSharePublicLink1/changingPublicLinkShare.feature:90](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink1/changingPublicLinkShare.feature#L90)

#### [copying a folder within a public link folder to folder with same name as an already existing file overwrites the parent file](https://github.com/owncloud/ocis/issues/1232)

- [apiSharePublicLink2/copyFromPublicLink.feature:63](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink2/copyFromPublicLink.feature#L63)
- [apiSharePublicLink2/copyFromPublicLink.feature:89](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink2/copyFromPublicLink.feature#L89)
- [apiSharePublicLink2/copyFromPublicLink.feature:173](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink2/copyFromPublicLink.feature#L173)
- [apiSharePublicLink2/copyFromPublicLink.feature:174](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink2/copyFromPublicLink.feature#L174)
- [apiSharePublicLink2/copyFromPublicLink.feature:189](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink2/copyFromPublicLink.feature#L189)
- [apiSharePublicLink2/copyFromPublicLink.feature:190](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink2/copyFromPublicLink.feature#L190)
- [apiSharePublicLink3/updatePublicLinkShare.feature:45](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/updatePublicLinkShare.feature#L45)
- [apiSharePublicLink3/updatePublicLinkShare.feature:46](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/updatePublicLinkShare.feature#L46)

#### [Upload-only shares must not overwrite but create a separate file](https://github.com/owncloud/ocis/issues/1267)

- [apiSharePublicLink3/uploadToPublicLinkShare.feature:24](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/uploadToPublicLinkShare.feature#L24)
- [apiSharePublicLink3/uploadToPublicLinkShare.feature:273](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/uploadToPublicLinkShare.feature#L273)

#### [Set quota over settings](https://github.com/owncloud/ocis/issues/1290)

- [apiSharePublicLink3/uploadToPublicLinkShare.feature:160](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/uploadToPublicLinkShare.feature#L160)
- [apiSharePublicLink3/uploadToPublicLinkShare.feature:179](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/uploadToPublicLinkShare.feature#L179)


#### [deleting a file inside a received shared folder is moved to the trash-bin of the sharer not the receiver](https://github.com/owncloud/ocis/issues/1124)

- [apiTrashbin/trashbinSharingToShares.feature:29](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L29)
- [apiTrashbin/trashbinSharingToShares.feature:46](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L46)
- [apiTrashbin/trashbinSharingToShares.feature:51](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L51)
- [apiTrashbin/trashbinSharingToShares.feature:73](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L73)
- [apiTrashbin/trashbinSharingToShares.feature:78](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L78)
- [apiTrashbin/trashbinSharingToShares.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L100)
- [apiTrashbin/trashbinSharingToShares.feature:105](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L105)
- [apiTrashbin/trashbinSharingToShares.feature:128](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L128)
- [apiTrashbin/trashbinSharingToShares.feature:133](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L133)
- [apiTrashbin/trashbinSharingToShares.feature:156](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L156)
- [apiTrashbin/trashbinSharingToShares.feature:161](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L161)
- [apiTrashbin/trashbinSharingToShares.feature:184](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L184)
- [apiTrashbin/trashbinSharingToShares.feature:189](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L189)
- [apiTrashbin/trashbinSharingToShares.feature:212](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L212)
- [apiTrashbin/trashbinSharingToShares.feature:236](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinSharingToShares.feature#L236)

#### [changing user quota gives ocs status 103 / Cannot set quota](https://github.com/owncloud/product/issues/247)

- [apiShareOperationsToShares2/uploadToShare.feature:210](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/uploadToShare.feature#L210)
- [apiShareOperationsToShares2/uploadToShare.feature:211](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/uploadToShare.feature#L211)

#### [not possible to move file into a received folder](https://github.com/owncloud/ocis/issues/764)

- [apiShareOperationsToShares1/changingFilesShare.feature:25](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/changingFilesShare.feature#L25)
- [apiShareOperationsToShares1/changingFilesShare.feature:26](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/changingFilesShare.feature#L26)
- [apiShareOperationsToShares1/changingFilesShare.feature:109](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/changingFilesShare.feature#L109)
- [apiShareOperationsToShares1/changingFilesShare.feature:110](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/changingFilesShare.feature#L110)
- [apiShareOperationsToShares1/changingFilesShare.feature:131](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/changingFilesShare.feature#L131)
- [apiShareOperationsToShares1/changingFilesShare.feature:132](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares1/changingFilesShare.feature#L132)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:538](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L538)
- [apiVersions/fileVersionsSharingToShares.feature:220](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionsSharingToShares.feature#L220)
- [apiVersions/fileVersionsSharingToShares.feature:221](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionsSharingToShares.feature#L221)
- [apiWebdavMove2/moveShareOnOcis.feature:30](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveShareOnOcis.feature#L30)
- [apiWebdavMove2/moveShareOnOcis.feature:32](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveShareOnOcis.feature#L32)
- [apiWebdavMove2/moveShareOnOcis.feature:98](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveShareOnOcis.feature#L98)
- [apiWebdavMove2/moveShareOnOcis.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveShareOnOcis.feature#L100)
- [apiWebdavMove2/moveShareOnOcis.feature:169](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveShareOnOcis.feature#L169)
- [apiWebdavMove2/moveShareOnOcis.feature:170](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveShareOnOcis.feature#L170)

#### [Expiration date for shares is not implemented](https://github.com/owncloud/ocis/issues/1250)

#### Expiration date of user shares

- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:52](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L52)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:53](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L53)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:76](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L76)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:77](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L77)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:102](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L102)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:103](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L103)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:128](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L128)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:129](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L129)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:279](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L279)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:280](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L280)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:301](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L301)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:302](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L302)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:323](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L323)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:324](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L324)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:346](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L346)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:347](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L347)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:363](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L363)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:364](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L364)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:380](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L380)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:381](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L381)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:576](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L576)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:577](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L577)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:599](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L599)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:600](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L600)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:601](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L601)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:602](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L602)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:603](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L603)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:624](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L624)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:625](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L625)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:626](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L626)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:627](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L627)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:628](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L628)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:629](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L629)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:630](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L630)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:631](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L631)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:632](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L632)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:633](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L633)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:634](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L634)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:635](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L635)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:656](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L656)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:657](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L657)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:658](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L658)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:659](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L659)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:660](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L660)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:661](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L661)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:682](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L682)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:683](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L683)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:684](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L684)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:685](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L685)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:686](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L686)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:687](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L687)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:708](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L708)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:709](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L709)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:732](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L732)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:733](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L733)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:756](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L756)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:757](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L757)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:34](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L34)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:35](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L35)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:86](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L86)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:87](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L87)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:143](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L143)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:144](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L144)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:201](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L201)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:202](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L202)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:203](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L203)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:204](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L204)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:287](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L287)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:288](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L288)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:318](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L318)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:319](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L319)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:320](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L320)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:321](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L321)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:379](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L379)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:380](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L380)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:381](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L381)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:382](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L382)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:383](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L383)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:384](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L384)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:413](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L413)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:414](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L414)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:415](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L415)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:416](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L416)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:444](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L444)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:445](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L445)

#### Expiration date of group shares

- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:175](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L175)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:176](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L176)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:201](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L201)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:202](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L202)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:229](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L229)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:230](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L230)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:258](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L258)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:259](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L259)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:403](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L403)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:404](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L404)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:427](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L427)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:428](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L428)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:451](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L451)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:452](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L452)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:476](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L476)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:477](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L477)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:497](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L497)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:498](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L498)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:518](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L518)
- [apiShareCreateSpecialToShares1/createShareExpirationDate.feature:519](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares1/createShareExpirationDate.feature#L519)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:60](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L60)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:61](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L61)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:116](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L116)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:117](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L117)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:172](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L172)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:173](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L173)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:232](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L232)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:233](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L233)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:234](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L234)
- [apiShareReshareToShares3/reShareWithExpiryDate.feature:235](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareReshareToShares3/reShareWithExpiryDate.feature#L235)

#### [Cannot move folder/file from one received share to another](https://github.com/owncloud/ocis/issues/2442)

- [apiShareUpdateToShares/updateShare.feature:241](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareUpdateToShares/updateShare.feature#L241)
- [apiShareUpdateToShares/updateShare.feature:195](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareUpdateToShares/updateShare.feature#L195)
- [apiShareManagementToShares/mergeShare.feature:141](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/mergeShare.feature#L141)

#### [Sharing folder and sub-folder with same user but different permission,the permission of sub-folder is not obeyed ](https://github.com/owncloud/ocis/issues/2440)

- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:304](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L304)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:344](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L344)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:470](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L470)
- [apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:510](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L510)

#### [Empty OCS response for a share create request using a disabled user](https://github.com/owncloud/ocis/issues/2212)

- [apiShareCreateSpecialToShares2/createShareWithDisabledUser.feature:20](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareWithDisabledUser.feature#L20)
- [apiShareCreateSpecialToShares2/createShareWithDisabledUser.feature:23](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareCreateSpecialToShares2/createShareWithDisabledUser.feature#L23)


#### [Edit user share response has a "name" field](https://github.com/owncloud/ocis/issues/1225)

- [apiShareUpdateToShares/updateShare.feature:287](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareUpdateToShares/updateShare.feature#L287)
- [apiShareUpdateToShares/updateShare.feature:288](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareUpdateToShares/updateShare.feature#L288)

#### [user can access version metadata of a received share before accepting it](https://github.com/owncloud/ocis/issues/760)

- [apiVersions/fileVersionsSharingToShares.feature:283](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiVersions/fileVersionsSharingToShares.feature#L283)

#### [Share lists deleted user as 'user'](https://github.com/owncloud/ocis/issues/903)

- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:700](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L700)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:701](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L701)

#### [deleting a share with wrong authentication returns OCS status 996 / HTTP 500](https://github.com/owncloud/ocis/issues/1229)

- [apiShareManagementBasicToShares/deleteShareFromShares.feature:256](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L256)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:257](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L257)

### User Management

User and group management features

#### [incorrect ocs(v2) status value when getting info of share that does not exist should be 404, gives 998](https://github.com/owncloud/product/issues/250)

_ocs: api compatibility, return correct status code_

- [apiShareOperationsToShares2/shareAccessByID.feature:48](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L48)
- [apiShareOperationsToShares2/shareAccessByID.feature:49](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L49)
- [apiShareOperationsToShares2/shareAccessByID.feature:50](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L50)
- [apiShareOperationsToShares2/shareAccessByID.feature:51](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L51)
- [apiShareOperationsToShares2/shareAccessByID.feature:52](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L52)
- [apiShareOperationsToShares2/shareAccessByID.feature:53](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L53)
- [apiShareOperationsToShares2/shareAccessByID.feature:54](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L54)
- [apiShareOperationsToShares2/shareAccessByID.feature:55](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L55)

### Other

API, search, favorites, config, capabilities, not existing endpoints, CORS and others

#### [Ability to return error messages in Webdav response bodies](https://github.com/owncloud/ocis/issues/1293)

- [apiAuthOcs/ocsDELETEAuth.feature:8](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsDELETEAuth.feature#L8)
- [apiAuthOcs/ocsGETAuth.feature:8](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsGETAuth.feature#L8)
- [apiAuthOcs/ocsGETAuth.feature:42](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsGETAuth.feature#L42)
- [apiAuthOcs/ocsGETAuth.feature:73](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsGETAuth.feature#L73)
- [apiAuthOcs/ocsGETAuth.feature:104](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsGETAuth.feature#L104)
- [apiAuthOcs/ocsGETAuth.feature:121](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsGETAuth.feature#L121)
- [apiAuthOcs/ocsPOSTAuth.feature:8](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPOSTAuth.feature#L8)
- [apiAuthOcs/ocsPUTAuth.feature:8](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPUTAuth.feature#L8)

#### [sending MKCOL requests to another user's webDav endpoints as normal user gives 404 instead of 403 ](https://github.com/owncloud/ocis/issues/3872)

_ocdav: api compatibility, return correct status code_

- [apiAuthWebDav/webDavMKCOLAuth.feature:52](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavMKCOLAuth.feature#L52)
- [apiAuthWebDav/webDavMKCOLAuth.feature:66](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavMKCOLAuth.feature#L66)

#### [trying to lock file of another user gives http 200](https://github.com/owncloud/ocis/issues/2176)

- [apiAuthWebDav/webDavLOCKAuth.feature:56](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavLOCKAuth.feature#L56)
- [apiAuthWebDav/webDavLOCKAuth.feature:68](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavLOCKAuth.feature#L68)

#### [send (MOVE,COPY) requests to another user's webDav endpoints as normal user gives 400 instead of 403](https://github.com/owncloud/ocis/issues/3882)

_ocdav: api compatibility, return correct status code_

- [apiAuthWebDav/webDavMOVEAuth.feature:55](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavMOVEAuth.feature#L55)
- [apiAuthWebDav/webDavMOVEAuth.feature:64](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavMOVEAuth.feature#L64)
- [apiAuthWebDav/webDavCOPYAuth.feature:59](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavCOPYAuth.feature#L59)
- [apiAuthWebDav/webDavCOPYAuth.feature:68](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavCOPYAuth.feature#L68)

#### [send POST requests to another user's webDav endpoints as normal user](https://github.com/owncloud/ocis/issues/1287)

_ocdav: api compatibility, return correct status code_

- [apiAuthWebDav/webDavPOSTAuth.feature:56](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavPOSTAuth.feature#L56)
- [apiAuthWebDav/webDavPOSTAuth.feature:65](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavPOSTAuth.feature#L65)

#### Another users space literally does not exist because it is not listed as a space for him, 404 seems correct, expects 403

- [apiAuthWebDav/webDavPUTAuth.feature:56](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavPUTAuth.feature#L56)
- [apiAuthWebDav/webDavPUTAuth.feature:68](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavPUTAuth.feature#L68)

#### [Using double slash in URL to access a folder gives 501 and other status codes](https://github.com/owncloud/ocis/issues/1667)

- [apiAuthWebDav/webDavSpecialURLs.feature:13](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L13)
- [apiAuthWebDav/webDavSpecialURLs.feature:24](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L24)
- [apiAuthWebDav/webDavSpecialURLs.feature:55](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L55)
- [apiAuthWebDav/webDavSpecialURLs.feature:66](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L66)
- [apiAuthWebDav/webDavSpecialURLs.feature:76](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L76)
- [apiAuthWebDav/webDavSpecialURLs.feature:88](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L88)
- [apiAuthWebDav/webDavSpecialURLs.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L100)
- [apiAuthWebDav/webDavSpecialURLs.feature:111](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L111)
- [apiAuthWebDav/webDavSpecialURLs.feature:121](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L121)
- [apiAuthWebDav/webDavSpecialURLs.feature:132](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L132)
- [apiAuthWebDav/webDavSpecialURLs.feature:142](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L142)
- [apiAuthWebDav/webDavSpecialURLs.feature:153](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L153)
- [apiAuthWebDav/webDavSpecialURLs.feature:163](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L163)
- [apiAuthWebDav/webDavSpecialURLs.feature:174](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L174)
- [apiAuthWebDav/webDavSpecialURLs.feature:184](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L184)
- [apiAuthWebDav/webDavSpecialURLs.feature:195](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavSpecialURLs.feature#L195)

#### [Difference in response content of status.php and default capabilities](https://github.com/owncloud/ocis/issues/1286)

- [apiCapabilities/capabilitiesWithNormalUser.feature:11](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilitiesWithNormalUser.feature#L11)

#### [[old/new/spaces] In ocis and oc10, REPORT request response differently](https://github.com/owncloud/ocis/issues/4712)

- [apiWebdavOperations/search.feature:208](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L208)
- [apiWebdavOperations/search.feature:209](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L209)
- [apiWebdavOperations/search.feature:214](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L214)
- [apiWebdavOperations/search.feature:240](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L240)
- [apiWebdavOperations/search.feature:241](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L241)
- [apiWebdavOperations/search.feature:246](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L246)

#### [Support for favorites](https://github.com/owncloud/ocis/issues/1228)

- [apiFavorites/favorites.feature:115](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L115)
- [apiFavorites/favorites.feature:116](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L116)
- [apiFavorites/favorites.feature:142](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L142)
- [apiFavorites/favorites.feature:143](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L143)
- [apiFavorites/favorites.feature:264](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L264)
- [apiFavorites/favorites.feature:265](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L265)

And other missing implementation of favorites

- [apiFavorites/favorites.feature:183](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L183)
- [apiFavorites/favorites.feature:184](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L184)
- [apiFavorites/favorites.feature:189](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L189)
- [apiFavorites/favorites.feature:216](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L216)
- [apiFavorites/favorites.feature:217](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L217)
- [apiFavorites/favorites.feature:222](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L222)
- [apiFavorites/favoritesSharingToShares.feature:67](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favoritesSharingToShares.feature#L67)
- [apiFavorites/favoritesSharingToShares.feature:68](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favoritesSharingToShares.feature#L68)


#### [WWW-Authenticate header for unauthenticated requests is not clear](https://github.com/owncloud/ocis/issues/2285)

- [apiWebdavOperations/refuseAccess.feature:22](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/refuseAccess.feature#L22)
- [apiWebdavOperations/refuseAccess.feature:23](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/refuseAccess.feature#L23)

#### [wildcard Access-Control-Allow-Origin](https://github.com/owncloud/ocis/issues/1340)

- [apiAuth/cors.feature:24](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L24)
- [apiAuth/cors.feature:25](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L25)
- [apiAuth/cors.feature:26](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L26)
- [apiAuth/cors.feature:27](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L27)
- [apiAuth/cors.feature:28](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L28)
- [apiAuth/cors.feature:29](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L29)
- [apiAuth/cors.feature:30](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L30)
- [apiAuth/cors.feature:31](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L31)
- [apiAuth/cors.feature:32](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L32)
- [apiAuth/cors.feature:33](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L33)
- [apiAuth/cors.feature:44](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L44)
- [apiAuth/cors.feature:45](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L45)
- [apiAuth/cors.feature:46](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L46)
- [apiAuth/cors.feature:47](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L47)
- [apiAuth/cors.feature:48](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L48)
- [apiAuth/cors.feature:49](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L49)
- [apiAuth/cors.feature:68](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L68)
- [apiAuth/cors.feature:69](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L69)
- [apiAuth/cors.feature:70](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L70)
- [apiAuth/cors.feature:71](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L71)
- [apiAuth/cors.feature:72](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L72)
- [apiAuth/cors.feature:73](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L73)
- [apiAuth/cors.feature:92](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L92)
- [apiAuth/cors.feature:93](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L93)
- [apiAuth/cors.feature:94](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L94)
- [apiAuth/cors.feature:95](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L95)
- [apiAuth/cors.feature:96](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L96)
- [apiAuth/cors.feature:97](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L97)
- [apiAuth/cors.feature:98](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L98)
- [apiAuth/cors.feature:99](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L99)
- [apiAuth/cors.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L100)
- [apiAuth/cors.feature:101](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L101)
- [apiAuth/cors.feature:112](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L112)
- [apiAuth/cors.feature:113](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L113)
- [apiAuth/cors.feature:114](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L114)
- [apiAuth/cors.feature:115](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L115)
- [apiAuth/cors.feature:116](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L116)
- [apiAuth/cors.feature:117](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L117)
- [apiAuth/cors.feature:136](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L136)
- [apiAuth/cors.feature:137](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L137)
- [apiAuth/cors.feature:138](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L138)
- [apiAuth/cors.feature:139](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L139)
- [apiAuth/cors.feature:140](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L140)
- [apiAuth/cors.feature:141](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L141)
- [apiAuth/cors.feature:160](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L160)
- [apiAuth/cors.feature:161](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L161)
- [apiAuth/cors.feature:162](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L162)
- [apiAuth/cors.feature:163](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L163)
- [apiAuth/cors.feature:164](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L164)
- [apiAuth/cors.feature:165](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L165)
- [apiAuth/cors.feature:166](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L166)
- [apiAuth/cors.feature:167](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L167)
- [apiAuth/cors.feature:178](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L178)
- [apiAuth/cors.feature:179](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L179)
- [apiAuth/cors.feature:180](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L180)
- [apiAuth/cors.feature:181](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L181)
- [apiAuth/cors.feature:182](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L182)
- [apiAuth/cors.feature:183](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L183)
- [apiAuth/cors.feature:204](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L204)
- [apiAuth/cors.feature:205](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L205)
- [apiAuth/cors.feature:206](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L206)
- [apiAuth/cors.feature:207](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L207)
- [apiAuth/cors.feature:208](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L208)
- [apiAuth/cors.feature:209](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuth/cors.feature#L209)

#### [App Passwords/Tokens for legacy WebDAV clients](https://github.com/owncloud/ocis/issues/197)

- [apiAuthWebDav/webDavDELETEAuth.feature:135](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavDELETEAuth.feature#L135)
- [apiAuthWebDav/webDavDELETEAuth.feature:149](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavDELETEAuth.feature#L149)
- [apiAuthWebDav/webDavDELETEAuth.feature:161](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavDELETEAuth.feature#L161)
- [apiAuthWebDav/webDavDELETEAuth.feature:175](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthWebDav/webDavDELETEAuth.feature#L175)

#### [Request to edit non-existing user by authorized admin gets unauthorized in http response](https://github.com/owncloud/core/issues/38423)

- [apiAuthOcs/ocsPUTAuth.feature:24](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPUTAuth.feature#L24)

#### [Sharing a same file twice to the same group](https://github.com/owncloud/ocis/issues/1710)

- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:748](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L748)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:749](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L749)

#### [PATCH request for TUS upload with wrong checksum gives incorrect response](https://github.com/owncloud/ocis/issues/1755)

- [apiWebdavUploadTUS/checksums.feature:84](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L84)
- [apiWebdavUploadTUS/checksums.feature:85](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L85)
- [apiWebdavUploadTUS/checksums.feature:86](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L86)
- [apiWebdavUploadTUS/checksums.feature:87](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L87)
- [apiWebdavUploadTUS/checksums.feature:92](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L92)
- [apiWebdavUploadTUS/checksums.feature:93](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L93)
- [apiWebdavUploadTUS/checksums.feature:173](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L173)
- [apiWebdavUploadTUS/checksums.feature:174](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L174)
- [apiWebdavUploadTUS/checksums.feature:179](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L179)
- [apiWebdavUploadTUS/checksums.feature:226](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L226)
- [apiWebdavUploadTUS/checksums.feature:227](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L227)
- [apiWebdavUploadTUS/checksums.feature:228](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L228)
- [apiWebdavUploadTUS/checksums.feature:229](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L229)
- [apiWebdavUploadTUS/checksums.feature:234](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L234)
- [apiWebdavUploadTUS/checksums.feature:235](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L235)
- [apiWebdavUploadTUS/checksums.feature:282](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L282)
- [apiWebdavUploadTUS/checksums.feature:283](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L283)
- [apiWebdavUploadTUS/checksums.feature:284](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L284)
- [apiWebdavUploadTUS/checksums.feature:285](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L285)
- [apiWebdavUploadTUS/checksums.feature:290](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L290)
- [apiWebdavUploadTUS/checksums.feature:291](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/checksums.feature#L291)
- [apiWebdavUploadTUS/optionsRequest.feature:8](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/optionsRequest.feature#L8)
- [apiWebdavUploadTUS/optionsRequest.feature:22](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/optionsRequest.feature#L22)
- [apiWebdavUploadTUS/optionsRequest.feature:35](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/optionsRequest.feature#L35)
- [apiWebdavUploadTUS/optionsRequest.feature:49](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/optionsRequest.feature#L49)
- [apiWebdavUploadTUS/uploadToShare.feature:176](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L176)
- [apiWebdavUploadTUS/uploadToShare.feature:177](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L177)
- [apiWebdavUploadTUS/uploadToShare.feature:195](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L195)
- [apiWebdavUploadTUS/uploadToShare.feature:196](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L196)
- [apiWebdavUploadTUS/uploadToShare.feature:214](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L214)
- [apiWebdavUploadTUS/uploadToShare.feature:215](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L215)
- [apiWebdavUploadTUS/uploadToShare.feature:253](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L253)
- [apiWebdavUploadTUS/uploadToShare.feature:254](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L254)
- [apiWebdavUploadTUS/uploadToShare.feature:295](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L295)
- [apiWebdavUploadTUS/uploadToShare.feature:296](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/uploadToShare.feature#L296)

#### [TUS OPTIONS requests do not reply with TUS headers when invalid password](https://github.com/owncloud/ocis/issues/1012)

- [apiWebdavUploadTUS/optionsRequest.feature:62](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/optionsRequest.feature#L62)
- [apiWebdavUploadTUS/optionsRequest.feature:76](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/optionsRequest.feature#L76)
- [apiWebdavUploadTUS/optionsRequest.feature:89](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/optionsRequest.feature#L89)
- [apiWebdavUploadTUS/optionsRequest.feature:104](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUploadTUS/optionsRequest.feature#L104)

#### [Trying to accept a share with invalid ID gives incorrect OCS and HTTP status](https://github.com/owncloud/ocis/issues/2111)

- [apiShareOperationsToShares2/shareAccessByID.feature:85](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L85)
- [apiShareOperationsToShares2/shareAccessByID.feature:86](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L86)
- [apiShareOperationsToShares2/shareAccessByID.feature:87](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L87)
- [apiShareOperationsToShares2/shareAccessByID.feature:88](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L88)
- [apiShareOperationsToShares2/shareAccessByID.feature:89](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L89)
- [apiShareOperationsToShares2/shareAccessByID.feature:90](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L90)
- [apiShareOperationsToShares2/shareAccessByID.feature:91](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L91)
- [apiShareOperationsToShares2/shareAccessByID.feature:92](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L92)
- [apiShareOperationsToShares2/shareAccessByID.feature:104](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L104)
- [apiShareOperationsToShares2/shareAccessByID.feature:105](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L105)
- [apiShareOperationsToShares2/shareAccessByID.feature:143](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L143)
- [apiShareOperationsToShares2/shareAccessByID.feature:144](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L144)
- [apiShareOperationsToShares2/shareAccessByID.feature:145](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L145)
- [apiShareOperationsToShares2/shareAccessByID.feature:146](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L146)
- [apiShareOperationsToShares2/shareAccessByID.feature:147](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L147)
- [apiShareOperationsToShares2/shareAccessByID.feature:148](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L148)
- [apiShareOperationsToShares2/shareAccessByID.feature:149](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L149)
- [apiShareOperationsToShares2/shareAccessByID.feature:150](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L150)
- [apiShareOperationsToShares2/shareAccessByID.feature:162](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L162)
- [apiShareOperationsToShares2/shareAccessByID.feature:163](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L163)


#### [Shares to deleted group listed in the response](https://github.com/owncloud/ocis/issues/2441)

- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:534](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L534)
- [apiShareManagementBasicToShares/createShareToSharesFolder.feature:535](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/createShareToSharesFolder.feature#L535)

#### [copying the file inside Shares folder returns 404](https://github.com/owncloud/ocis/issues/3874)

- [apiWebdavProperties1/copyFile.feature:410](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L410)
- [apiWebdavProperties1/copyFile.feature:411](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L411)
- [apiWebdavProperties1/copyFile.feature:416](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L416)
- [apiWebdavProperties1/copyFile.feature:436](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L436)
- [apiWebdavProperties1/copyFile.feature:437](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L437)
- [apiWebdavProperties1/copyFile.feature:442](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L442)

### Won't fix

Not everything needs to be implemented for ocis. While the oc10 testsuite covers these things we are not looking at them right now.

- _The `OC-LazyOps` header is [no longer supported by the client](https://github.com/owncloud/client/pull/8398), implementing this is not necessary for a first production release. We plan to have an upload state machine to visualize the state of a file, see https://github.com/owncloud/ocis/issues/214_
- _Blacklisted ignored files are no longer required because ocis can handle `.htaccess` files without security implications introduced by serving user provided files with apache._

#### [uploading with old-chunking does not work](https://github.com/owncloud/ocis/issues/1343)

- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:20](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L20)
- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:21](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L21)
- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:26](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L26)
- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:39](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L39)
- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:40](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L40)
- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:45](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L45)
- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:80](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L80)
- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:81](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L81)
- [apiWebdavUpload1/uploadFileToExcludedDirectory.feature:86](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToExcludedDirectory.feature#L86)

#### [Blacklist files extensions](https://github.com/owncloud/ocis/issues/2177)

- [apiWebdavProperties1/copyFile.feature:120](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L120)
- [apiWebdavProperties1/copyFile.feature:121](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L121)
- [apiWebdavProperties1/copyFile.feature:126](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/copyFile.feature#L126)
- [apiWebdavProperties1/createFileFolder.feature:98](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/createFileFolder.feature#L98)
- [apiWebdavProperties1/createFileFolder.feature:99](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/createFileFolder.feature#L99)
- [apiWebdavProperties1/createFileFolder.feature:104](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavProperties1/createFileFolder.feature#L104)
- [apiWebdavUpload1/uploadFile.feature:181](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFile.feature#L181)
- [apiWebdavUpload1/uploadFile.feature:182](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFile.feature#L182)
- [apiWebdavUpload1/uploadFile.feature:187](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFile.feature#L187)
- [apiWebdavUpload2/uploadFileToBlacklistedNameUsingOldChunking.feature:19](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToBlacklistedNameUsingOldChunking.feature#L19)
- [apiWebdavUpload2/uploadFileToBlacklistedNameUsingOldChunking.feature:35](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToBlacklistedNameUsingOldChunking.feature#L35)
- [apiWebdavUpload2/uploadFileToBlacklistedNameUsingOldChunking.feature:36](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToBlacklistedNameUsingOldChunking.feature#L36)
- [apiWebdavUpload2/uploadFileToBlacklistedNameUsingOldChunking.feature:37](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToBlacklistedNameUsingOldChunking.feature#L37)
- [apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature:13](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature#L13)
- [apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature:20](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature#L20)
- [apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature:38](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature#L38)
- [apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature:39](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature#L39)
- [apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature:40](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload2/uploadFileToExcludedDirectoryUsingOldChunking.feature#L40)
- [apiWebdavMove2/moveFile.feature:287](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFile.feature#L287)
- [apiWebdavMove2/moveFile.feature:288](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFile.feature#L288)
- [apiWebdavMove2/moveFile.feature:293](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFile.feature#L293)

#### [cannot set blacklisted file names](https://github.com/owncloud/product/issues/260)

- [apiWebdavMove1/moveFolderToBlacklistedName.feature:21](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L21)
- [apiWebdavMove1/moveFolderToBlacklistedName.feature:22](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L22)
- [apiWebdavMove1/moveFolderToBlacklistedName.feature:27](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L27)
- [apiWebdavMove1/moveFolderToBlacklistedName.feature:40](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L40)
- [apiWebdavMove1/moveFolderToBlacklistedName.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L41)
- [apiWebdavMove1/moveFolderToBlacklistedName.feature:46](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L46)
- [apiWebdavMove1/moveFolderToBlacklistedName.feature:81](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L81)
- [apiWebdavMove1/moveFolderToBlacklistedName.feature:82](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L82)
- [apiWebdavMove1/moveFolderToBlacklistedName.feature:87](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToBlacklistedName.feature#L87)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:19](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L19)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:20](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L20)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:35](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L35)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:36](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L36)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:74](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L74)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:75](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L75)

#### [cannot set excluded directories](https://github.com/owncloud/product/issues/261)

- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:22](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L22)
- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:23](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L23)
- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:28](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L28)
- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:42](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L42)
- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:43](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L43)
- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:48](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L48)
- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:84](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L84)
- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:85](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L85)
- [apiWebdavMove1/moveFolderToExcludedDirectory.feature:90](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolderToExcludedDirectory.feature#L90)
- [apiWebdavMove2/moveFileToExcludedDirectory.feature:20](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L20)
- [apiWebdavMove2/moveFileToExcludedDirectory.feature:21](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L21)
- [apiWebdavMove2/moveFileToExcludedDirectory.feature:37](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L37)
- [apiWebdavMove2/moveFileToExcludedDirectory.feature:38](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L38)
- [apiWebdavMove2/moveFileToExcludedDirectory.feature:78](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L78)
- [apiWebdavMove2/moveFileToExcludedDirectory.feature:79](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L79)

#### [system configuration options missing](https://github.com/owncloud/ocis/issues/1323)

- [apiWebdavUpload1/uploadFileToBlacklistedName.feature:31](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToBlacklistedName.feature#L31)
- [apiWebdavUpload1/uploadFileToBlacklistedName.feature:32](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToBlacklistedName.feature#L32)
- [apiWebdavUpload1/uploadFileToBlacklistedName.feature:37](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToBlacklistedName.feature#L37)
- [apiWebdavUpload1/uploadFileToBlacklistedName.feature:71](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToBlacklistedName.feature#L71)
- [apiWebdavUpload1/uploadFileToBlacklistedName.feature:70](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToBlacklistedName.feature#L70)
- [apiWebdavUpload1/uploadFileToBlacklistedName.feature:76](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavUpload1/uploadFileToBlacklistedName.feature#L76)

#### [Allow public link sharing only for certain groups feature not implemented](https://github.com/owncloud/ocis/issues/4623)

- [apiSharePublicLink3/allowGroupToCreatePublicLinks.feature:35](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/allowGroupToCreatePublicLinks.feature#L35)
- [apiSharePublicLink3/allowGroupToCreatePublicLinks.feature:92](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiSharePublicLink3/allowGroupToCreatePublicLinks.feature#L92)

#### [Preview of text file with UTF content does not render correctly](https://github.com/owncloud/ocis/issues/2570)

- [apiWebdavPreviews/previews.feature:211](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavPreviews/previews.feature#L211)

#### [Share path in the response is different between share states](https://github.com/owncloud/ocis/issues/2540)

- [apiShareManagementToShares/acceptShares.feature:65](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L65)
- [apiShareManagementToShares/acceptShares.feature:93](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L93)
- [apiShareManagementToShares/acceptShares.feature:224](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L224)
- [apiShareManagementToShares/acceptShares.feature:252](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L252)
- [apiShareManagementToShares/acceptShares.feature:295](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L295)
- [apiShareManagementToShares/acceptShares.feature:334](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L334)
- [apiShareManagementToShares/acceptShares.feature:571](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/acceptShares.feature#L571)
- [apiShareOperationsToShares2/shareAccessByID.feature:124](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L124)
- [apiShareOperationsToShares2/shareAccessByID.feature:125](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareOperationsToShares2/shareAccessByID.feature#L125)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:190](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L190)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:191](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L191)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:192](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L192)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:193](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L193)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:218](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L218)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:219](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L219)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:220](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L220)
- [apiShareManagementBasicToShares/deleteShareFromShares.feature:221](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementBasicToShares/deleteShareFromShares.feature#L221)

#### [Content-type is not multipart/byteranges when downloading file with Range Header](https://github.com/owncloud/ocis/issues/2677)

- [apiWebdavOperations/downloadFile.feature:229](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/downloadFile.feature#L229)
- [apiWebdavOperations/downloadFile.feature:230](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/downloadFile.feature#L230)
- [apiWebdavOperations/downloadFile.feature:235](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/downloadFile.feature#L235)

#### [moveShareInsideAnotherShare behaves differently on oCIS than oC10](https://github.com/owncloud/ocis/issues/3047)

- [apiShareManagementToShares/moveShareInsideAnotherShare.feature:25](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/moveShareInsideAnotherShare.feature#L25)
- [apiShareManagementToShares/moveShareInsideAnotherShare.feature:86](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/moveShareInsideAnotherShare.feature#L86)
- [apiShareManagementToShares/moveShareInsideAnotherShare.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/moveShareInsideAnotherShare.feature#L100)

#### [Renaming resource to banned name is allowed in spaces webdav](https://github.com/owncloud/ocis/issues/3099)

- [apiWebdavMove1/moveFolder.feature:27](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolder.feature#L27)
- [apiWebdavMove1/moveFolder.feature:45](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolder.feature#L45)
- [apiWebdavMove1/moveFolder.feature:63](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove1/moveFolder.feature#L63)
- [apiWebdavMove2/moveFile.feature:224](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFile.feature#L224)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:25](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L25)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L41)
- [apiWebdavMove2/moveFileToBlacklistedName.feature:80](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToBlacklistedName.feature#L80)

#### [REPORT method on spaces returns an incorrect d:href response](https://github.com/owncloud/ocis/issues/3111)

- [apiFavorites/favorites.feature:121](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L121)
- [apiFavorites/favorites.feature:148](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L148)
- [apiFavorites/favorites.feature:270](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L270)

#### [could not create system tag](https://github.com/owncloud/ocis/issues/3092)

- [apiWebdavOperations/search.feature:274](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L274)
- [apiWebdavOperations/search.feature:291](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L291)
- [apiWebdavOperations/search.feature:317](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/search.feature#L317)

#### [Cannot disable the dav propfind depth infinity for resources](https://github.com/owncloud/ocis/issues/3720)

- [apiWebdavOperations/listFiles.feature:383](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/listFiles.feature#L383)
- [apiWebdavOperations/listFiles.feature:384](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/listFiles.feature#L384)
- [apiWebdavOperations/listFiles.feature:389](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/listFiles.feature#L389)
- [apiWebdavOperations/listFiles.feature:408](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/listFiles.feature#L408)
- [apiWebdavOperations/listFiles.feature:413](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/listFiles.feature#L413)
- [apiWebdavOperations/listFiles.feature:427](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/listFiles.feature#L427)
- [apiWebdavOperations/listFiles.feature:428](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/listFiles.feature#L428)
- [apiWebdavOperations/listFiles.feature:433](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/listFiles.feature#L433)

#### [Renaming resource to excluded directory name is allowed in spaces webdav](https://github.com/owncloud/ocis/issues/3102)

- [apiWebdavMove2/moveFileToExcludedDirectory.feature:26](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L26)
- [apiWebdavMove2/moveFileToExcludedDirectory.feature:43](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L43)
- [apiWebdavMove2/moveFileToExcludedDirectory.feature:84](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavMove2/moveFileToExcludedDirectory.feature#L84)

### [graph/users: enable/disable users](https://github.com/owncloud/ocis/issues/3064)

- [apiWebdavOperations/refuseAccess.feature:35](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/refuseAccess.feature#L35)
- [apiWebdavOperations/refuseAccess.feature:36](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/refuseAccess.feature#L36)
- [apiWebdavOperations/refuseAccess.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiWebdavOperations/refuseAccess.feature#L41)

#### [OCS status code zero](https://github.com/owncloud/ocis/issues/3621)

- [apiShareManagementToShares/moveReceivedShare.feature:32](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiShareManagementToShares/moveReceivedShare.feature#L32)

#### [HTTP status code differ while deleting file of another user's trash bin](https://github.com/owncloud/ocis/issues/3544)

- [apiTrashbin/trashbinDelete.feature:109](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L109)

#### [Problem accessing trashbin with personal space id](https://github.com/owncloud/ocis/issues/3639)

- [apiTrashbin/trashbinDelete.feature:35](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L35)
- [apiTrashbin/trashbinDelete.feature:36](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L36)
- [apiTrashbin/trashbinDelete.feature:58](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L58)
- [apiTrashbin/trashbinDelete.feature:85](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L85)
- [apiTrashbin/trashbinDelete.feature:132](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L132)
- [apiTrashbin/trashbinDelete.feature:155](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L155)
- [apiTrashbin/trashbinDelete.feature:180](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L180)
- [apiTrashbin/trashbinDelete.feature:205](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L205)
- [apiTrashbin/trashbinDelete.feature:242](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L242)
- [apiTrashbin/trashbinDelete.feature:279](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L279)
- [apiTrashbin/trashbinDelete.feature:328](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinDelete.feature#L328)
- [apiTrashbin/trashbinFilesFolders.feature:25](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L25)
- [apiTrashbin/trashbinFilesFolders.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L41)
- [apiTrashbin/trashbinFilesFolders.feature:60](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L60)
- [apiTrashbin/trashbinFilesFolders.feature:81](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L81)
- [apiTrashbin/trashbinFilesFolders.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L100)
- [apiTrashbin/trashbinFilesFolders.feature:136](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L136)
- [apiTrashbin/trashbinFilesFolders.feature:159](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L159)
- [apiTrashbin/trashbinFilesFolders.feature:339](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L339)
- [apiTrashbin/trashbinFilesFolders.feature:340](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L340)
- [apiTrashbin/trashbinFilesFolders.feature:341](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L341)
- [apiTrashbin/trashbinFilesFolders.feature:346](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L346)
- [apiTrashbin/trashbinFilesFolders.feature:347](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L347)
- [apiTrashbin/trashbinFilesFolders.feature:348](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L348)
- [apiTrashbin/trashbinFilesFolders.feature:349](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L349)
- [apiTrashbin/trashbinFilesFolders.feature:350](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L350)
- [apiTrashbin/trashbinFilesFolders.feature:351](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L351)
- [apiTrashbin/trashbinFilesFolders.feature:368](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L368)
- [apiTrashbin/trashbinFilesFolders.feature:388](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L388)
- [apiTrashbin/trashbinFilesFolders.feature:442](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L442)
- [apiTrashbin/trashbinFilesFolders.feature:479](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiTrashbin/trashbinFilesFolders.feature#L479)

Note: always have an empty line at the end of this file.
The bash script that processes this file requires that the last line has a newline on the end.
