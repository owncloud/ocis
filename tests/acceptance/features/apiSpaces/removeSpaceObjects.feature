Feature: Remove files, folder
  As a user
  I want to be able to remove files, folders
  So that I can remove unnecessary objects

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "delete objects" with the default quota using the Graph API
    And user "Alice" has created a folder "folderForDeleting/sub1/sub2" in space "delete objects"
    And user "Alice" has uploaded a file inside space "delete objects" with content "some content" to "text.txt"


  Scenario Outline: user deletes a folder with some subfolders in a space via the webDav API
    Given user "Alice" has sent the following space share invitation:
      | space           | delete objects |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | <space-role>   |
    When user "<user>" removes the folder "folderForDeleting" from space "delete objects"
    Then the HTTP status code should be "<http-status-code>"
    And for user "<user>" the space "delete objects" <should-or-not-be-in-space> contain these entries:
      | folderForDeleting |
    And as "<user>" folder "folderForDeleting" <should-or-not-be-in-trash> exist in the trashbin of the space "delete objects"
    Examples:
      | user  | space-role   | http-status-code | should-or-not-be-in-space | should-or-not-be-in-trash |
      | Alice | Manager      | 204              | should not                | should                    |
      | Brian | Manager      | 204              | should not                | should                    |
      | Brian | Space Editor | 204              | should not                | should                    |
      | Brian | Space Viewer | 403              | should                    | should not                |


  Scenario Outline: user deletes a subfolder in a space via the webDav API
    Given user "Alice" has sent the following space share invitation:
      | space           | delete objects |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | <space-role>   |
    When user "<user>" removes the folder "folderForDeleting/sub1" from space "delete objects"
    Then the HTTP status code should be "<http-status-code>"
    And for user "<user>" the space "delete objects" should contain these entries:
      | folderForDeleting |
    And for user "<user>" folder "folderForDeleting/" of the space "delete objects" <should-or-not-be-in-space> contain these entries:
      | sub1 |
    And as "<user>" folder "sub1" <should-or-not-be-in-trash> exist in the trashbin of the space "delete objects"
    Examples:
      | user  | space-role   | http-status-code | should-or-not-be-in-space | should-or-not-be-in-trash |
      | Alice | Manager      | 204              | should not                | should                    |
      | Brian | Manager      | 204              | should not                | should                    |
      | Brian | Space Editor | 204              | should not                | should                    |
      | Brian | Space Viewer | 403              | should                    | should not                |


  Scenario Outline: user deletes a file in a space via the webDav API
    Given user "Alice" has sent the following space share invitation:
      | space           | delete objects |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | <space-role>   |
    When user "<user>" removes the file "text.txt" from space "delete objects"
    Then the HTTP status code should be "<http-status-code>"
    And for user "<user>" the space "delete objects" <should-or-not-be-in-space> contain these entries:
      | text.txt |
    And as "<user>" file "text.txt" <should-or-not-be-in-trash> exist in the trashbin of the space "delete objects"
    Examples:
      | user  | space-role   | http-status-code | should-or-not-be-in-space | should-or-not-be-in-trash |
      | Alice | Manager      | 204              | should not                | should                    |
      | Brian | Manager      | 204              | should not                | should                    |
      | Brian | Space Editor | 204              | should not                | should                    |
      | Brian | Space Viewer | 403              | should                    | should not                |


  Scenario: try to delete an empty string folder from a space
    When user "Alice" removes the folder "" from space "delete objects"
    Then the HTTP status code should be "405"
