@api @skipOnOcV10
Feature: Remove files, folder
  As a user
  I want to be able to remove files, folders
  Users with the editor role can also remove objects
  Users with the viewer role cannot remove objects

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "delete objects" with the default quota using the GraphApi
    And user "Alice" has created a folder "folderForDeleting/sub1/sub2" in space "delete objects"
    And user "Alice" has uploaded a file inside space "delete objects" with content "some content" to "text.txt"


  Scenario Outline: An user deletes a folder with some subfolders in a Space via the webDav API
    Given user "Alice" has shared a space "delete objects" to user "Brian" with role "<role>"
    When user "<user>" removes the folder "folderForDeleting" from space "delete objects"
    Then the HTTP status code should be "<code>"
    And for user "<user>" the space "delete objects" <shouldOrNotBeInSpace> contain these entries:
      | folderForDeleting |
    And as "<user>" folder "folderForDeleting" <shouldOrNotBeInTrash> exist in the trashbin of the space "delete objects"
    Examples:
      | user  | role    | code | shouldOrNotBeInSpace | shouldOrNotBeInTrash |
      | Alice | manager | 204  | should not           | should               |
      | Brian | manager | 204  | should not           | should               |
      | Brian | editor  | 204  | should not           | should               |
      | Brian | viewer  | 403  | should               | should not           |


  Scenario Outline: An user deletes a subfolder in a Space via the webDav API
    Given user "Alice" has shared a space "delete objects" to user "Brian" with role "<role>"
    When user "<user>" removes the folder "folderForDeleting/sub1" from space "delete objects"
    Then the HTTP status code should be "<code>"
    And for user "<user>" the space "delete objects" should contain these entries:
      | folderForDeleting |
    And for user "<user>" folder "folderForDeleting/" of the space "delete objects" <shouldOrNotBeInSpace> contain these entries:
      | sub1 |
    And as "<user>" folder "sub1" <shouldOrNotBeInTrash> exist in the trashbin of the space "delete objects"
    Examples:
      | user  | role    | code | shouldOrNotBeInSpace | shouldOrNotBeInTrash |
      | Alice | manager | 204  | should not           | should               |
      | Brian | manager | 204  | should not           | should               |
      | Brian | editor  | 204  | should not           | should               |
      | Brian | viewer  | 403  | should               | should not           |


  Scenario Outline: An user deletes a file in a Space via the webDav API
    Given user "Alice" has shared a space "delete objects" to user "Brian" with role "<role>"
    When user "<user>" removes the file "text.txt" from space "delete objects"
    Then the HTTP status code should be "<code>"
    And for user "<user>" the space "delete objects" <shouldOrNotBeInSpace> contain these entries:
      | text.txt |
    And as "<user>" file "text.txt" <shouldOrNotBeInTrash> exist in the trashbin of the space "delete objects"
    And the user "<user>" should have a space called "delete objects" with these key and value pairs:
      | key          | value          |
      | name         | delete objects |
      | quota@@@used | <quotaValue>   |
    Examples:
      | user  | role    | code | shouldOrNotBeInSpace | shouldOrNotBeInTrash | quotaValue |
      | Alice | manager | 204  | should not           | should               | 0          |
      | Brian | manager | 204  | should not           | should               | 0          |
      | Brian | editor  | 204  | should not           | should               | 0          |
      | Brian | viewer  | 403  | should               | should not           | 12         |


  Scenario: An user is unable to delete a Space via the webDav API
    When user "Alice" removes the folder "" from space "delete objects"
    Then the HTTP status code should be "400"
    And the user "Alice" should have a space called "delete objects" with these key and value pairs:
      | key  | value          |
      | name | delete objects |
