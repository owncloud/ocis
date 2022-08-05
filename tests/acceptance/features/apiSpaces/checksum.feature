@api @skipOnOcV10
Feature: checksums

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @files_sharing-app-required @issue-ocis-reva-196
  Scenario: Sharing a file with checksum should return the checksum in the propfind using new DAV path
    Given the administrator has set the default folder for received shares to "Shares"
    And auto-accept shares has been disabled
    And using spaces DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    And user "Alice" has shared file "/myChecksumFile.txt" with user "Brian"
    And user "Brian" has accepted share "/myChecksumFile.txt" offered by user "Alice"
    When user "Brian" requests the checksum of "/myChecksumFile.txt" in space "Shares Jail" via propfind
    Then the HTTP status code should be "207"
    And the webdav checksum should match "SHA1:3ee962b839762adb0ad8ba6023a4690be478de6f MD5:d70b40f177b14b470d1756a3c12b963a ADLER32:8ae90960"

  @files_sharing-app-required @issue-ocis-reva-196
  Scenario: Modifying a shared file should return correct checksum in the propfind using new DAV path
    Given the administrator has set the default folder for received shares to "Shares"
    And auto-accept shares has been disabled
    And using spaces DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    And user "Alice" has shared file "/myChecksumFile.txt" with user "Brian"
    And user "Brian" has accepted share "/myChecksumFile.txt" offered by user "Alice"
    When user "Brian" uploads file with checksum "SHA1:ce5582148c6f0c1282335b87df5ed4be4b781399" and content "Some Text" to "/myChecksumFile.txt" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "204"
    And as user "Alice" the webdav checksum of "/myChecksumFile.txt" via propfind should match "SHA1:ce5582148c6f0c1282335b87df5ed4be4b781399 MD5:56e57920c3c8c727bfe7a5288cdf61c4 ADLER32:1048035a"
