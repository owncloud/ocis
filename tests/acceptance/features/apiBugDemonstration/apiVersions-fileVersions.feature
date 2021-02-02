@api @files_versions-app-required @skipOnOcis-EOS-Storage @issue-ocis-reva-275

Feature: dav-versions

  Background:
    Given using OCS API version "2"
    And using new DAV path
    And user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-reva-17 @issue-ocis-reva-56
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Upload file and no version is available using various chunking methods
    When user "Alice" uploads file "filesForUpload/davtest.txt" to filenames based on "/davtest.txt" with all mechanisms using the WebDAV API
    Then the version folder of file "/davtest.txt-olddav-regular" for user "Alice" should contain "0" elements
    Then the version folder of file "/davtest.txt-newdav-regular" for user "Alice" should contain "0" elements
    Then the version folder of file "/davtest.txt-olddav-oldchunking" for user "Alice" should contain "0" elements
    And as "Alice" file "/davtest.txt-newdav-newchunking" should not exist

  @skipOnOcis-OC-Storage @issue-ocis-reva-17 @issue-ocis-reva-56
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Upload a file twice and versions are available using various chunking methods
    When user "Alice" uploads file "filesForUpload/davtest.txt" to filenames based on "/davtest.txt" with all mechanisms using the WebDAV API
    And user "Alice" uploads file "filesForUpload/davtest.txt" to filenames based on "/davtest.txt" with all mechanisms using the WebDAV API
    Then the version folder of file "/davtest.txt-olddav-regular" for user "Alice" should contain "1" element
    And the version folder of file "/davtest.txt-newdav-regular" for user "Alice" should contain "1" element
    Then the version folder of file "/davtest.txt-olddav-oldchunking" for user "Alice" should contain "0" element
    And as "Alice" file "/davtest.txt-newdav-newchunking" should not exist

  @files_sharing-app-required
  @issue-ocis-reva-243
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: sharer of a file can see the old version information when the sharee changes the content of the file
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "First content" to "sharefile.txt"
    And user "Alice" has shared file "sharefile.txt" with user "Brian"
    When user "Brian" has uploaded file with content "Second content" to "/sharefile.txt"
    Then the HTTP status code should be "201"
    And the version folder of file "/sharefile.txt" for user "Alice" should contain "0" element
#    And the version folder of file "/sharefile.txt" for user "Alice" should contain "1" element

  @files_sharing-app-required
  @issue-ocis-reva-243
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: sharer of a file can restore the original content of a shared file after the file has been modified by the sharee
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "First content" to "sharefile.txt"
    And user "Alice" has shared file "sharefile.txt" with user "Brian"
    And user "Brian" has uploaded file with content "Second content" to "/sharefile.txt"
    When user "Alice" restores version index "0" of file "/sharefile.txt" using the WebDAV API
#    When user "Alice" restores version index "1" of file "/sharefile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/sharefile.txt" for user "Alice" should be "First content"
    And the content of file "/sharefile.txt" for user "Brian" should be "Second content"
#    And the content of file "/sharefile.txt" for user "Brian" should be "First content"

  @files_sharing-app-required
  @issue-ocis-reva-243 @issue-ocis-reva-386
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Moving a file (with versions) into a shared folder as the sharee and as the sharer
    Given using <dav_version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | change    |
      | shareWith   | Alice     |
    And user "Brian" has uploaded file with content "test data 1" to "/testfile.txt"
    And user "Brian" has uploaded file with content "test data 2" to "/testfile.txt"
    And user "Brian" has uploaded file with content "test data 3" to "/testfile.txt"
    And user "Brian" moves file "/testfile.txt" to "/testshare/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/testshare/testfile.txt" for user "Alice" should be ""
#    And the content of file "/testshare/testfile.txt" for user "Alice" should be "test data 3"
    And the content of file "/testshare/testfile.txt" for user "Brian" should be "test data 3"
    And as "Brian" file "/testfile.txt" should not exist
    And as "Alice" file "/testshare/testfile.txt" should not exist
    And the content of file "/testshare/testfile.txt" for user "Brian" should be "test data 3"
#    And the version folder of file "/testshare/testfile.txt" for user "Alice" should contain "2" elements
#    And the version folder of file "/testshare/testfile.txt" for user "Brian" should contain "2" elements
    Examples:
      | dav_version |
      | old         |
      | new         |

  @files_sharing-app-required
  @issue-ocis-reva-243 @issue-ocis-reva-386
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Moving a file (with versions) out of a shared folder as the sharee and as the sharer
    Given using <dav_version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "test data 1" to "/testshare/testfile.txt"
    And user "Brian" has uploaded file with content "test data 2" to "/testshare/testfile.txt"
    And user "Brian" has uploaded file with content "test data 3" to "/testshare/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | change    |
      | shareWith   | Alice     |
    When user "Brian" moves file "/testshare/testfile.txt" to "/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/testfile.txt" for user "Brian" should be "test data 3"
    And as "Alice" file "/testshare/testfile.txt" should not exist
    And as "Brian" file "/testshare/testfile.txt" should not exist
#    And the version folder of file "/testfile.txt" for user "Brian" should contain "2" elements
    Examples:
      | dav_version |
      | old         |
      | new         |

  @skipOnStorage:ceph @files_primary_s3-issue-161 @files_sharing-app-required
  @issue-ocis-reva-376
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Receiver tries get file versions of shared file from the sharer
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "textfile0" to "textfile0.txt"
    And user "Alice" has uploaded file with content "version 1" to "textfile0.txt"
    And user "Alice" has uploaded file with content "version 2" to "textfile0.txt"
    And user "Alice" has uploaded file with content "version 3" to "textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Brian" tries to get versions of file "textfile0.txt" from "Alice"
    Then the HTTP status code should be "207"
    And the number of versions should be "4"
#    And the number of versions should be "3"
