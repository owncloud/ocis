Feature: Download file in project space
  As a user with different role
  I want to be able to download files
  So that I can have it in my local storage

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "download file" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "download file" with content "some content" to "file.txt"
    And user "Alice" has sent the following space share invitation:
      | space           | download file |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Space Editor  |
    And user "Alice" has sent the following space share invitation:
      | space           | download file |
      | sharee          | Bob           |
      | shareType       | user          |
      | permissionsRole | Space Viewer  |


  Scenario Outline: user downloads a file in the project space
    When user "<user>" downloads the file "file.txt" of the space "download file" using the WebDAV API
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header         | value |
      | Content-Length | 12    |
    Examples:
      | user  |
      | Alice |
      | Brian |
      | Bob   |


  Scenario Outline: users with role manager and editor can download an old version of the file in the project space
    Given user "Alice" has uploaded a file inside space "download file" with content "new content" to "file.txt"
    And user "Alice" has uploaded a file inside space "download file" with content "newest content" to "file.txt"
    When user "<user>" downloads version of the file "file.txt" with the index "1" of the space "download file" using the WebDAV API
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header         | value |
      | Content-Length | 11    |
    When user "<user>" downloads version of the file "file.txt" with the index "2" of the space "download file" using the WebDAV API
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header         | value |
      | Content-Length | 12    |
    Examples:
      | user  |
      | Alice |
      | Brian |


  Scenario: user with role viewer cannot get versions of a file in the project space
    Given user "Alice" has uploaded a file inside space "download file" with content "new content" to "file.txt"
    And user "Alice" has uploaded a file inside space "download file" with content "newest content" to "file.txt"
    When user "Bob" tries to get versions of the file "file.txt" from the space "download file" using the WebDAV API
    Then the HTTP status code should be "403"
    When user "Bob" tries to download version of the file "file.txt" with the index "1" of the space "download file" using the WebDAV API
    Then the HTTP status code should be "403"
