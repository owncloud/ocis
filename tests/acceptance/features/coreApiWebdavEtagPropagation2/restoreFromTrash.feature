Feature: propagation of etags when restoring a file or folder from trash
  As a client app
  I want metadata (etags) of parent folders to change when a file or folder is restored
  So that the client app can know to re-scan and sync the content of the folder(s)

  Background:
    Given user "Alice" has been created with default attributes
    And user "Alice" has created folder "/upload"


  Scenario Outline: restoring a file to its original location changes the etags of all parents
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/upload/sub"
    And user "Alice" has uploaded file with content "uploaded content" to "/upload/sub/file.txt"
    And user "Alice" has deleted file "/upload/sub/file.txt"
    And user "Alice" has stored etag of element "/"
    And user "Alice" has stored etag of element "/upload"
    And user "Alice" has stored etag of element "/upload/sub"
    When user "Alice" restores the file with original path "/upload/sub/file.txt" using the trashbin API
    Then the HTTP status code should be "201"
    And these etags should have changed:
      | user  | path        |
      | Alice | /           |
      | Alice | /upload     |
      | Alice | /upload/sub |
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: restoring a file to an other location changes the etags of all parents
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/upload/sub"
    And user "Alice" has created folder "/restore"
    And user "Alice" has created folder "/restore/sub"
    And user "Alice" has uploaded file with content "uploaded content" to "/upload/sub/file.txt"
    And user "Alice" has deleted file "/upload/sub/file.txt"
    And user "Alice" has stored etag of element "/"
    And user "Alice" has stored etag of element "/restore"
    And user "Alice" has stored etag of element "/restore/sub"
    When user "Alice" restores the file with original path "/upload/sub/file.txt" to "/restore/sub/file.txt" using the trashbin API
    Then the HTTP status code should be "201"
    And these etags should have changed:
      | user  | path         |
      | Alice | /            |
      | Alice | /restore     |
      | Alice | /restore/sub |
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: restoring a folder to its original location changes the etags of all parents
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/upload/sub"
    And user "Alice" has created folder "/upload/sub/toDelete"
    And user "Alice" has deleted folder "/upload/sub/toDelete"
    And user "Alice" has stored etag of element "/"
    And user "Alice" has stored etag of element "/upload"
    And user "Alice" has stored etag of element "/upload/sub"
    When user "Alice" restores the folder with original path "/upload/sub/toDelete" using the trashbin API
    Then the HTTP status code should be "201"
    And these etags should have changed:
      | user  | path        |
      | Alice | /           |
      | Alice | /upload     |
      | Alice | /upload/sub |
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: restoring a folder to an other location changes the etags of all parents
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/upload/sub"
    And user "Alice" has created folder "/upload/sub/toDelete"
    And user "Alice" has deleted folder "/upload/sub/toDelete"
    And user "Alice" has created folder "/restore"
    And user "Alice" has created folder "/restore/sub"
    And user "Alice" has stored etag of element "/"
    And user "Alice" has stored etag of element "/restore"
    And user "Alice" has stored etag of element "/restore/sub"
    When user "Alice" restores the folder with original path "/upload/sub/toDelete" to "/restore/sub/toDelete" using the trashbin API
    Then the HTTP status code should be "201"
    And these etags should have changed:
      | user  | path         |
      | Alice | /            |
      | Alice | /restore     |
      | Alice | /restore/sub |
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |
