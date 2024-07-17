Feature: Download space
  As a user
  I want to download space
  So that I can store it locally


  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "Project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project-space" with content "some data" to "file1.txt"
    And user "Alice" has created a folder ".space" in space "Project-space"
    And user "Alice" has uploaded a file inside space "Project-space" with content "space description" to ".space/readme.md"


  Scenario: user downloads a space
    Given user "Alice" has uploaded a file inside space "Project-space" with content "other data" to "file2.txt"
    When user "Alice" downloads the space "Project-space" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded "zip" archive should contain these files:
      | name             | content           |
      | file1.txt        | some data         |
      | file2.txt        | other data        |
      | .space/readme.md | space description |


  Scenario Outline: user downloads a shared space (shared by others)
    Given user "Alice" has sent the following space share invitation:
      | space           | Project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    When user "Brian" downloads the space "Project-space" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded "zip" archive should contain these files:
      | name             | content           |
      | file1.txt        | some data         |
      | .space/readme.md | space description |
    Examples:
      | space-role   |
      | Manager      |
      | Space Editor |
      | Space Viewer |


  Scenario Outline: admin/space-admin tries to download a space that they do not have access to
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When user "Brian" tries to download the space "Project-space" owned by user "Alice" using the WebDAV API
    Then the HTTP status code should be "404"
    Examples:
      | user-role   |
      | Admin       |
      | Space Admin |


  Scenario: user tries to download disabled space
    Given user "Alice" has disabled a space "Project-space"
    When user "Alice" tries to download the space "Project-space" using the WebDAV API
    Then the HTTP status code should be "404"
