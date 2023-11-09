Feature: move (rename) file
  As a user
  I want to be able to move and rename files
  So that I can manage my file system

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @skipOnReva
  Scenario Outline: moving a file into a shared folder as the sharee and as the sharer
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | change    |
      | shareWith   | Alice     |
    And user "<mover>" has uploaded file with content "test data" to "/testfile.txt"
    When user "<mover>" moves file "/testfile.txt" to "<destination_folder>/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/Shares/testshare/testfile.txt" for user "Alice" should be "test data"
    And the content of file "/testshare/testfile.txt" for user "Brian" should be "test data"
    And as "<mover>" file "/testfile.txt" should not exist
    Examples:
      | dav-path-version | mover | destination_folder |
      | old              | Alice | /Shares/testshare  |
      | old              | Brian | /testshare         |
      | new              | Alice | /Shares/testshare  |
      | new              | Brian | /testshare         |

  @skipOnReva
  Scenario Outline: moving a file out of a shared folder as the sharer
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "test data" to "/testshare/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | change    |
      | shareWith   | Alice     |
    When user "Brian" moves file "/testshare/testfile.txt" to "/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/testfile.txt" for user "Brian" should be "test data"
    And as "Alice" file "/Shares/testshare/testfile.txt" should not exist
    And as "Brian" file "/testshare/testfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva
  Scenario Outline: moving a file out of a shared folder as the sharee
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "test data" to "/testshare/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | change    |
      | shareWith   | Alice     |
    When user "Alice" moves file "/Shares/testshare/testfile.txt" to "/testfile.txt" using the WebDAV API
    Then the HTTP status code should be "502"
    And as "Alice" file "/Shares/testshare/testfile.txt" should exist
    And as "Brian" file "/testshare/testfile.txt" should exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva
  Scenario Outline: moving a folder into a shared folder as the sharee and as the sharer
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | change    |
      | shareWith   | Alice     |
    And user "<mover>" has created folder "/testsubfolder"
    And user "<mover>" has uploaded file with content "test data" to "/testsubfolder/testfile.txt"
    When user "<mover>" moves folder "/testsubfolder" to "<destination_folder>/testsubfolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/Shares/testshare/testsubfolder/testfile.txt" for user "Alice" should be "test data"
    And the content of file "/testshare/testsubfolder/testfile.txt" for user "Brian" should be "test data"
    And as "<mover>" file "/testsubfolder" should not exist
    Examples:
      | dav-path-version | mover | destination_folder |
      | old              | Alice | /Shares/testshare  |
      | old              | Brian | /testshare         |
      | new              | Alice | /Shares/testshare  |
      | new              | Brian | /testshare         |


  Scenario Outline: moving a folder out of a shared folder as the sharer
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created the following folders
      | path                     |
      | /testshare               |
      | /testshare/testsubfolder |
    And user "Brian" has uploaded file with content "test data" to "/testshare/testsubfolder/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | change    |
      | shareWith   | Alice     |
    When user "Brian" moves folder "/testshare/testsubfolder" to "/testsubfolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/testsubfolder/testfile.txt" for user "Brian" should be "test data"
    And as "Alice" folder "/testshare/testsubfolder" should not exist
    And as "Brian" folder "/testshare/testsubfolder" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva
  Scenario Outline: moving a folder out of a shared folder as the sharee
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created the following folders
      | path                     |
      | /testshare               |
      | /testshare/testsubfolder |
    And user "Brian" has uploaded file with content "test data" to "/testshare/testsubfolder/testfile.txt"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | change    |
      | shareWith   | Alice     |
    When user "Alice" moves folder "/Shares/testshare/testsubfolder" to "/testsubfolder" using the WebDAV API
    Then the HTTP status code should be "502"
    And as "Alice" folder "/Shares/testshare/testsubfolder" should exist
    And as "Brian" folder "/testshare/testsubfolder" should exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva
  Scenario Outline: moving a file to a shared folder with no permissions
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | read      |
      | shareWith   | Alice     |
    When user "Alice" moves file "/textfile0.txt" to "/Shares/testshare/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And user "Alice" should not be able to download file "/Shares/testshare/textfile0.txt"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva
  Scenario Outline: moving a file to overwrite a file in a shared folder with no permissions
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "textfile0.txt"
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "Welcome to ownCloud" to "fileToCopy.txt"
    And user "Brian" has created a share with settings
      | path        | testshare |
      | shareType   | user      |
      | permissions | read      |
      | shareWith   | Alice     |
    And user "Brian" has copied file "/fileToCopy.txt" to "/testshare/overwritethis.txt"
    When user "Alice" moves file "/textfile0.txt" to "/Shares/testshare/overwritethis.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And the content of file "/Shares/testshare/overwritethis.txt" for user "Alice" should be "Welcome to ownCloud"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
