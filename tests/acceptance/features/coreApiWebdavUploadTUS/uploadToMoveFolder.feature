Feature: move folders
  As a user
  I want to be able to move and upload files/folders
  So that I can organise my data structure

  Background:
    Given user "Alice" has been created with default attributes

  @issue-10346
  Scenario Outline: uploading file into a moved folder
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/test"
    And user "Alice" has created folder "/test-moved"
    And user "Alice" has moved folder "/test-moved" to "/test/test-moved"
    When user "Alice" uploads file with content "uploaded content" to "/test/test-moved/textfile.txt" using the TUS protocol on the WebDAV API
    Then the HTTP status code should be "204"
    And as "Alice" file "/test/test-moved/textfile.txt" should exist
    And the content of file "/test/test-moved/textfile.txt" for user "Alice" should be "uploaded content"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

