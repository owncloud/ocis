@api @skipOnOcV10
Feature: Restore files, folder
    As a user with manager and editor role
    I want to be able to restore files, folders
    Users with the viewer role cannot restore objects

    Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
    See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

    Background:
        Given these users have been created with default attributes and without skeleton files:
            | username |
            | Alice    |
            | Brian    |
            | Bob      |
            | Carol    |
        And the administrator has given "Alice" the role "Space Admin" using the settings api
        And user "Alice" creates a space "restore objects" of type "project" with the default quota using the GraphApi
        And user "Alice" has created a folder "newFolder" in space "restore objects"
        And user "Alice" has uploaded a file inside space "restore objects" with content "test" to "newFolder/file.txt"


    Scenario Outline: An user with different role can see deleted objects in trash bin of the space via the webDav API
        Given user "Alice" has shared a space "restore objects" to user "Brian" with role "<role>"
        And user "Alice" has removed the file "newFolder/file.txt" from space "restore objects"
        And user "Alice" has removed the folder "newFolder" from space "restore objects"
        When user "<user>" lists all deleted files in the trash bin of the space "restore objects"
        Then the HTTP status code should be "207"
        And as "<user>" folder "newFolder" should exist in the trashbin of the space "restore objects"
        And as "<user>" file "file.txt" should exist in the trashbin of the space "restore objects"
        Examples:
            | user  | role    |
            | Brian | manager |
            | Brian | editor  |
            | Brian | viewer  |


    Scenario Outline: An user restores a folder with some objects from the trash via the webDav API
        Given user "Alice" has shared a space "restore objects" to user "Brian" with role "<role>"
        And user "Alice" has removed the folder "newFolder" from space "restore objects"
        When user "<user>" restores the folder "newFolder" from the trash of the space "restore objects" to "/newFolder"
        Then the HTTP status code should be "<code>"
        And for user "<user>" the space "restore objects" <shouldOrNotBeInSpace> contain these entries:
            | newFolder |
        And as "<user>" folder "newFolder" <shouldOrNotBeInTrash> exist in the trashbin of the space "restore objects"
        Examples:
            | user  | role    | code | shouldOrNotBeInSpace | shouldOrNotBeInTrash |
            | Alice | manager | 201  | should               | should not           |
            | Brian | manager | 201  | should               | should not           |
            | Brian | editor  | 201  | should               | should not           |
            | Brian | viewer  | 403  | should not           | should               |


    Scenario Outline: An user restores a file from the trash via the webDav API
        Given user "Alice" has shared a space "restore objects" to user "Brian" with role "<role>"
        And user "Alice" has removed the file "newFolder/file.txt" from space "restore objects"
        When user "<user>" restores the file "file.txt" from the trash of the space "restore objects" to "newFolder/file.txt"
        Then the HTTP status code should be "<code>"
        And for user "<user>" folder "newFolder" of the space "restore objects" <shouldOrNotBeInSpace> contain these files:
            | file.txt |
        And as "<user>" file "file.txt" <shouldOrNotBeInTrash> exist in the trashbin of the space "restore objects"
        Examples:
            | user  | role    | code | shouldOrNotBeInSpace | shouldOrNotBeInTrash |
            | Alice | manager | 201  | should               | should not           |
            | Brian | manager | 201  | should               | should not           |
            | Brian | editor  | 201  | should               | should not           |
            | Brian | viewer  | 403  | should not           | should               |
