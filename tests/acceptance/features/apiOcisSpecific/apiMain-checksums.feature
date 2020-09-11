@api @skipOnOcis-OCIS-Storage
Feature: checksums

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-reva-98
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Uploading a file with checksum should return the checksum in the download header
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    When user "Alice" downloads file "/myChecksumFile.txt" using the WebDAV API
    Then the following headers should not be set
      | header      |
      | OC-Checksum |
    Examples:
      | dav_version |
      | old         |
      | new         |

  @issue-ocis-reva-98
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Copying file with checksum should return the checksum in the download header using new DAV path
    Given using new DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/myChecksumFile.txt" with checksum "MD5:d70b40f177b14b470d1756a3c12b963a"
    When user "Alice" copies file "/myChecksumFile.txt" to "/myChecksumFileCopy.txt" using the WebDAV API
    And user "Alice" downloads file "/myChecksumFileCopy.txt" using the WebDAV API
    Then the following headers should not be set
      | header      |
      | OC-Checksum |

  @issue-ocis-reva-99
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Upload a file where checksum does not match
    Given using <dav_version> DAV path
    When user "Alice" uploads file with checksum "SHA1:f005ba11" and content "Some Text" to "/chksumtst.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Alice" should see the following elements
      | /chksumtst.txt |
    Examples:
      | dav_version |
      | old         |
      | new         |

  @issue-ocis-reva-99
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Uploaded file should have the same checksum when downloaded
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with checksum "SHA1:ce5582148c6f0c1282335b87df5ed4be4b781399" and content "Some Text" to "/chksumtst.txt"
    When user "Alice" downloads file "/chksumtst.txt" using the WebDAV API
    Then the following headers should not be set
      | header      |
      | OC-Checksum |
    Examples:
      | dav_version |
      | old         |
      | new         |
