Feature: Restore files, folder
  As a user with manager and editor role
  I want to be able to restore files, folders
  So that I can get the resources that were accidentally deleted

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "restore objects" with the default quota using the Graph API
    And user "Alice" has created a folder "newFolder" in space "restore objects"
    And user "Alice" has uploaded a file inside space "restore objects" with content "test" to "newFolder/file.txt"


  Scenario Outline: user with different role can see deleted objects in trash bin of the space via the webDav API
    Given user "Alice" has sent the following space share invitation:
      | space           | restore objects |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | <space-role>    |
    And user "Alice" has removed the file "newFolder/file.txt" from space "restore objects"
    And user "Alice" has removed the folder "newFolder" from space "restore objects"
    When user "Brian" lists all deleted files in the trash bin of the space "restore objects"
    Then the HTTP status code should be "207"
    And as "Brian" folder "newFolder" should exist in the trashbin of the space "restore objects"
    And as "Brian" file "file.txt" should exist in the trashbin of the space "restore objects"
    Examples:
      | space-role   |
      | Manager      |
      | Space Editor |
      | Space Viewer |


  Scenario Outline: user can restore a folder with some objects from the trash via the webDav API
    Given user "Alice" has sent the following space share invitation:
      | space           | restore objects |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | <space-role>    |
    And user "Alice" has removed the folder "newFolder" from space "restore objects"
    When user "<user>" restores the folder "newFolder" from the trash of the space "restore objects" to "/newFolder"
    Then the HTTP status code should be "<http-status-code>"
    And for user "<user>" the space "restore objects" <should-or-not-be-in-space> contain these entries:
      | newFolder |
    And as "<user>" folder "newFolder" <should-or-not-be-in-trash> exist in the trashbin of the space "restore objects"
    Examples:
      | user  | space-role   | http-status-code | should-or-not-be-in-space | should-or-not-be-in-trash |
      | Alice | Manager      | 201              | should                    | should not                |
      | Brian | Manager      | 201              | should                    | should not                |
      | Brian | Space Editor | 201              | should                    | should not                |
      | Brian | Space Viewer | 403              | should not                | should                    |


  Scenario Outline: user can restore a file from the trash via the webDav API
    Given user "Alice" has sent the following space share invitation:
      | space           | restore objects |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | <space-role>    |
    And user "Alice" has removed the file "newFolder/file.txt" from space "restore objects"
    When user "<user>" restores the file "file.txt" from the trash of the space "restore objects" to "newFolder/file.txt"
    Then the HTTP status code should be "<http-status-code>"
    And for user "<user>" folder "newFolder" of the space "restore objects" <should-or-not-be-in-space> contain these files:
      | file.txt |
    And as "<user>" file "file.txt" <should-or-not-be-in-trash> exist in the trashbin of the space "restore objects"
    Examples:
      | user  | space-role   | http-status-code | should-or-not-be-in-space | should-or-not-be-in-trash |
      | Alice | Manager      | 201              | should                    | should not                |
      | Brian | Manager      | 201              | should                    | should not                |
      | Brian | Space Editor | 201              | should                    | should not                |
      | Brian | Space Viewer | 403              | should not                | should                    |


  Scenario Outline: only space manager can purge the trash via the webDav API
    Given user "Alice" has sent the following space share invitation:
      | space           | restore objects |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | <space-role>    |
    And the administrator has assigned the role "Space Admin" to user "Brian" using the Graph API
    And user "Alice" has removed the file "newFolder/file.txt" from space "restore objects"
    When user "Brian" deletes the file "file.txt" from the trash of the space "restore objects"
    Then the HTTP status code should be "<http-status-code>"
    And as "Brian" file "file.txt" <should-or-not-be-in-trash> exist in the trashbin of the space "restore objects"
    Examples:
      | space-role   | http-status-code | should-or-not-be-in-trash |
      | Manager      | 204              | should not                |
      | Space Editor | 403              | should                    |
      | Space Viewer | 403              | should                    |


  Scenario Outline: admin user who is not a member of space cannot see its trash bin
    Given user "Alice" has removed the file "newFolder/file.txt" from space "restore objects"
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When user "Brian" with admin permission lists all deleted files in the trash bin of the space "restore objects"
    Then the HTTP status code should be "404"
    Examples:
      | user-role   |
      | Space Admin |
      | Admin       |


  Scenario Outline: admin user without space-manager role cannot purge the trash
    Given user "Alice" has sent the following space share invitation:
      | space           | restore objects |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | Space Editor    |
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Alice" has removed the file "newFolder/file.txt" from space "restore objects"
    When user "Brian" tries to delete the file "file.txt" from the trash of the space "restore objects"
    Then the HTTP status code should be "403"
    And as "Alice" file "file.txt" should exist in the trashbin of the space "restore objects"
    Examples:
      | user-role   |
      | Space Admin |
      | Admin       |

