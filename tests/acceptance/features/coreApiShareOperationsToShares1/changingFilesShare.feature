@skipOnReva @issue-1289 @issue-1328
Feature: sharing
  As a user
  I want to move shares that I received
  So that I can organise them according to my needs

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @smokeTest
  Scenario Outline: moving a file into a share as recipient
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/shared"
    And user "Alice" has shared folder "/shared" with user "Brian"
    And user "Brian" has uploaded file with content "some data" to "/textfile0.txt"
    When user "Brian" moves file "textfile0.txt" to "/Shares/shared/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" file "/Shares/shared/shared_file.txt" should exist
    And as "Alice" file "/shared/shared_file.txt" should exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: move files between shares by same user
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "share1"
    And user "Alice" has created folder "share2"
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has moved file "textfile0.txt" to "share1/textfile0.txt"
    And user "Alice" has shared folder "/share1" with user "Brian"
    And user "Alice" has shared folder "/share2" with user "Brian"
    When user "Brian" moves file "/Shares/share1/textfile0.txt" to "/Shares/share2/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "502"
    And as "Brian" file "/Shares/share1/textfile0.txt" should exist
    And as "Alice" file "share1/textfile0.txt" should exist
    But as "Brian" file "/Shares/share2/textfile0.txt" should not exist
    And as "Alice" file "share2/textfile0.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: move files between shares by same user added by sharee
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "share1"
    And user "Alice" has created folder "share2"
    And user "Brian" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has shared folder "/share1" with user "Brian"
    And user "Alice" has shared folder "/share2" with user "Brian"
    When user "Brian" moves file "textfile0.txt" to "/Shares/share1/shared_file.txt" using the WebDAV API
    And user "Brian" moves file "/Shares/share1/shared_file.txt" to "/Shares/share2/shared_file.txt" using the WebDAV API
    Then the HTTP status code of responses on all endpoints should be "201"
    And as "Brian" file "/Shares/share1/shared_file.txt" should not exist
    And as "Alice" file "share1/shared_file.txt" should not exist
    But as "Brian" file "/Shares/share2/shared_file.txt" should exist
    And as "Alice" file "share2/shared_file.txt" should exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: move files between shares by different users
    Given using <dav-path-version> DAV path
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has created folder "/PARENT"
    And user "Brian" has created folder "/PARENT"
    And user "Alice" has moved file "textfile0.txt" to "PARENT/shared_file.txt"
    And user "Alice" has shared folder "/PARENT" with user "Carol"
    And user "Brian" has shared folder "/PARENT" with user "Carol"
    When user "Carol" moves file "/Shares/PARENT/shared_file.txt" to "/Shares/PARENT (2)/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Carol" file "/Shares/PARENT (2)/shared_file.txt" should exist
    And as "Brian" file "PARENT/shared_file.txt" should exist
    But as "Alice" file "PARENT/shared_file.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: overwrite a received file share
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "this is the old content" to "/textfile1.txt"
    And user "Alice" has shared file "/textfile1.txt" with user "Brian"
    When user "Brian" uploads file with content "this is a new content" to "/Shares/textfile1.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Brian" file "Shares/textfile1.txt" should exist
    And the content of file "Shares/textfile1.txt" for user "Brian" should be "this is a new content"
    And the content of file "textfile1.txt" for user "Alice" should be "this is a new content"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

