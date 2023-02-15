@api
Feature: checksums

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  Scenario Outline: Uploading a file with checksum should work
    Given using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @smokeTest @issue-ocis-reva-196
  Scenario Outline: Uploading a file with checksum should return the checksum in the propfind
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    When user "Alice" requests the checksum of "/myChecksumFile.txt" via propfind
    Then the webdav checksum should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f MD5:d70b40f177b14b470d1756a3c12b963a ADLER32:8ae90960"
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @smokeTest @issue-ocis-reva-98
  Scenario Outline: Uploading a file with checksum should return the checksum in the download header
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    When user "Alice" downloads file "/myChecksumFile.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the header checksum should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f"
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-196
  Scenario Outline: Moving a file with checksum should return the checksum in the propfind
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    When user "Alice" moves file "/myChecksumFile.txt" to "/myMovedChecksumFile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as user "Alice" the webdav checksum of "/myMovedChecksumFile.txt" via propfind should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f MD5:d70b40f177b14b470d1756a3c12b963a ADLER32:8ae90960"
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-98
  Scenario Outline: Downloading a file with checksum should return the checksum in the download header
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    And user "Alice" has moved file "/myChecksumFile.txt" to "/myMovedChecksumFile.txt"
    When user "Alice" downloads file "/myMovedChecksumFile.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the header checksum should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f"
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-196
  Scenario Outline: Uploading a chunked file with checksum should return the checksum in the propfind
    Given using <dav_version> DAV path
    And user "Alice" has uploaded chunk file "1" of "3" with "AAAAA" to "/myChecksumFile.txt" with checksum "MD5:45a72715acdd5019c5be30bdbb75233e"
    And user "Alice" has uploaded chunk file "2" of "3" with "BBBBB" to "/myChecksumFile.txt" with checksum "MD5:45a72715acdd5019c5be30bdbb75233e"
    And user "Alice" has uploaded chunk file "3" of "3" with "CCCCC" to "/myChecksumFile.txt" with checksum "MD5:45a72715acdd5019c5be30bdbb75233e"
    When user "Alice" requests the checksum of "/myChecksumFile.txt" via propfind
    Then the HTTP status code should be "207"
    And the webdav checksum should match "SHA1:acfa6b1565f9710d4d497c6035d5c069bd35a8e8 MD5:45a72715acdd5019c5be30bdbb75233e ADLER32:1ecd03df"
    Examples:
      | dav_version |
      | old         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-17
  Scenario Outline: Uploading a chunked file with checksum should return the checksum in the download header
    Given using <dav_version> DAV path
    And user "Alice" has uploaded chunk file "1" of "3" with "AAAAA" to "/myChecksumFile.txt" with checksum "MD5:45a72715acdd5019c5be30bdbb75233e"
    And user "Alice" has uploaded chunk file "2" of "3" with "BBBBB" to "/myChecksumFile.txt" with checksum "MD5:45a72715acdd5019c5be30bdbb75233e"
    And user "Alice" has uploaded chunk file "3" of "3" with "CCCCC" to "/myChecksumFile.txt" with checksum "MD5:45a72715acdd5019c5be30bdbb75233e"
    When user "Alice" downloads file "/myChecksumFile.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the header checksum should match "SHA1:acfa6b1565f9710d4d497c6035d5c069bd35a8e8"
    Examples:
      | dav_version |
      | old         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-14
  Scenario Outline: Moving file with checksum should return the checksum in the download header
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    When user "Alice" moves file "/myChecksumFile.txt" to "/myMovedChecksumFile.txt" using the WebDAV API
    And user "Alice" downloads file "/myMovedChecksumFile.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the header checksum should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f"
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-196
  Scenario Outline: Copying a file with checksum should return the checksum in the propfind using new DAV path
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    When user "Alice" copies file "/myChecksumFile.txt" to "/myChecksumFileCopy.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as user "Alice" the webdav checksum of "/myChecksumFileCopy.txt" via propfind should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f MD5:d70b40f177b14b470d1756a3c12b963a ADLER32:8ae90960"
    Examples:
      | dav_version |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-98
  Scenario Outline: Copying file with checksum should return the checksum in the download header using new DAV path
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    When user "Alice" copies file "/myChecksumFile.txt" to "/myChecksumFileCopy.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the header checksum when user "Alice" downloads file "/myChecksumFileCopy.txt" using the WebDAV API should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f"
    Examples:
      | dav_version |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @files_sharing-app-required @issue-ocis-reva-196
  Scenario Outline: Sharing a file with checksum should return the checksum in the propfind using new DAV path
    Given the administrator has set the default folder for received shares to "Shares"
    And auto-accept shares has been disabled
    And using <dav_version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    And user "Alice" has shared file "/myChecksumFile.txt" with user "Brian"
    And user "Brian" has accepted share "/myChecksumFile.txt" offered by user "Alice"
    When user "Brian" requests the checksum of "/Shares/myChecksumFile.txt" via propfind
    Then the HTTP status code should be "207"
    And the webdav checksum should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f MD5:d70b40f177b14b470d1756a3c12b963a ADLER32:8ae90960"
    Examples:
      | dav_version |
      | new         |

  @files_sharing-app-required @issue-ocis-reva-196 @skipOnOcV10
  Scenario Outline: Modifying a shared file should return correct checksum in the propfind using new DAV path
    Given the administrator has set the default folder for received shares to "Shares"
    And auto-accept shares has been disabled
    And using <dav_version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    And user "Alice" has shared file "/myChecksumFile.txt" with user "Brian"
    And user "Brian" has accepted share "/myChecksumFile.txt" offered by user "Alice"
    When user "Brian" uploads file with checksum "SHA1:ce5582148c6f0c1282335b87df5ed4be4b781399" and content "Some Text" to "/Shares/myChecksumFile.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as user "Alice" the webdav checksum of "/myChecksumFile.txt" via propfind should match "<checksum>"
    Examples:
      | dav_version | checksum                                                                                            |
      | new         | SHA1:ce5582148c6f0c1282335b87df5ed4be4b781399 MD5:56e57920c3c8c727bfe7a5288cdf61c4 ADLER32:1048035a |

  @issue-ocis-reva-99
  Scenario Outline: Upload a file where checksum does not match
    Given using <dav_version> DAV path
    When user "Alice" uploads file with checksum "SHA1:f005ba11" and content "Some Text" to "/chksumtst.txt" using the WebDAV API
    Then the HTTP status code should be "400"
    And user "Alice" should not see the following elements
      | /chksumtst.txt |
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |


  Scenario Outline: Upload a file where checksum does match
    Given using <dav_version> DAV path
    When user "Alice" uploads file with checksum "SHA1:ce5582148c6f0c1282335b87df5ed4be4b781399" and content "Some Text" to "/chksumtst.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-99
  Scenario Outline: Uploaded file should have the same checksum when downloaded
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with checksum "SHA1:ce5582148c6f0c1282335b87df5ed4be4b781399" and content "Some Text" to "/chksumtst.txt"
    When user "Alice" downloads file "/chksumtst.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header      | value                                         |
      | OC-Checksum | SHA1:ce5582148c6f0c1282335b87df5ed4be4b781399 |
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  ## Validation Plugin or Old Endpoint Specific
  @issue-ocis-reva-17
  Scenario Outline: Uploading an old method chunked file with checksum should fail using new DAV path
    Given using <dav_version> DAV path
    When user "Alice" uploads chunk file "1" of "3" with "AAAAA" to "/myChecksumFile.txt" with checksum "MD5:45a72715acdd5019c5be30bdbb75233e" using the WebDAV API
    Then the HTTP status code should be "503"
    And user "Alice" should not see the following elements
      | /myChecksumFile.txt |
    Examples:
      | dav_version |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  ## upload overwriting
  @issue-ocis-reva-196
  Scenario Outline: Uploading a file with MD5 checksum overwriting an existing file
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with content "some data" to "textfile0.txt"
    When user "Alice" uploads file "filesForUpload/textfile.txt" to "/textfile0.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "204"
    And as user "Alice" the webdav checksum of "/textfile0.txt" via propfind should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f MD5:d70b40f177b14b470d1756a3c12b963a ADLER32:8ae90960"
    And the content of file "/textfile0.txt" for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-196
  Scenario Outline: Uploading a file with SHA1 checksum overwriting an existing file
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with content "some data" to "textfile0.txt"
    When user "Alice" uploads file "filesForUpload/textfile.txt" to "/textfile0.txt" with checksum "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f" using the WebDAV API
    Then the HTTP status code should be "204"
    And as user "Alice" the webdav checksum of "/textfile0.txt" via propfind should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f MD5:d70b40f177b14b470d1756a3c12b963a ADLER32:8ae90960"
    And the content of file "/textfile0.txt" for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @skipOnStorage:ceph @skipOnStorage:scality @files_primary_s3-issue-224 @issue-ocis-reva-196
  Scenario Outline: Uploading a file with invalid SHA1 checksum overwriting an existing file
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    When user "Alice" uploads file "filesForUpload/textfile.txt" to "/textfile0.txt" with checksum "SHA1:f005ba11f005ba11f005ba11f005ba11f005ba11" using the WebDAV API
    Then the HTTP status code should be "400"
    And as user "Alice" the webdav checksum of "/textfile0.txt" via propfind should match "SHA1:2052377dec0724bda0d57aeab67fa819278b7f74 MD5:096e350e9ff1339a997a14145f9fc4b9 ADLER32:7d5a0921"
    And the content of file "/textfile0.txt" for user "Alice" should be "ownCloud test text file 0"
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |

  @issue-ocis-reva-214
  Scenario Outline: Uploading a file with checksum should work for file with special characters
    Given using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |
