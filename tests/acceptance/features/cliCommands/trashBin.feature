@env-config
Feature: trashbin
  As an administrator
  I want to manage trash-bin
  So that I can manage handle trashed resources efficiently

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: delete empty trashbin folders
    Given user "Alice" has created the following folders
      | path              |
      | folder-to-delete  |
      | folder-to-restore |
    And user "Alice" has uploaded file with content "test file" to "testfile.txt"
    And user "Alice" has deleted the following resources
      | path              |
      | folder-to-delete  |
      | folder-to-restore |
      | testfile.txt      |
    And user "Alice" has restored the folder with original path "folder-to-restore"
    And user "Alice" has deleted the folder with original path "folder-to-delete" from the trashbin
    And the administrator has stopped the server
    When the administrator deletes the empty trashbin folders using the CLI
    Then the command should be successful


  Scenario: list trashed resource of specific space
    Given user "Brian" has been created with default attributes
    And user "Alice" has created the following folders
      | path              |
      | folder-to-delete  |
      | folder-to-restore |
    And user "Alice" has uploaded file with content "test file" to "testfile.txt"
    And user "Alice" has deleted the following resources
      | path              |
      | folder-to-delete  |
      | folder-to-restore |
      | testfile.txt      |
    And user "Brian" has created folder "BrianFolder"
    And user "Brian" has deleted folder "BrianFolder"
    When the administrator lists all the trashed resources of space "Personal" owned by user "Alice"
    Then the command output should contain "3" trashed resources with the following information:
      | resource           | type   |
      | /folder-to-delete  | folder |
      | /folder-to-restore | folder |
      | /testfile.txt      | file   |


  Scenario: restore all trashed resource at once
    Given user "Brian" has been created with default attributes
    And user "Alice" has created the following folders
      | path              |
      | folder-to-delete  |
      | folder-to-restore |
    And user "Alice" has uploaded file with content "test file" to "testfile.txt"
    And user "Alice" has deleted the following resources
      | path              |
      | folder-to-delete  |
      | folder-to-restore |
      | testfile.txt      |
    And user "Brian" has created folder "BrianFolder"
    And user "Brian" has deleted folder "BrianFolder"
    When the administrator restores all the trashed resources of space "Personal" owned by user "Alice"
    Then the command should be successful
    And there should be no trashed resources of space "Personal" owned by user "Alice"
    And user "Alice" should see the following elements
      | /folder-to-delete  |
      | /folder-to-restore |
      | /testfile.txt      |
    And there should be "1" trashed resources of space "Personal" owned by user "Brian":
      | resource     | type   |
      | /BrianFolder | folder |


  Scenario: restore specific trashed resource at once
    Given user "Brian" has been created with default attributes
    And user "Alice" has created the following folders
      | path              |
      | folder-to-delete  |
      | folder-to-restore |
    And user "Alice" has uploaded file with content "test file" to "testfile.txt"
    And user "Alice" has deleted the following resources
      | path              |
      | folder-to-delete  |
      | folder-to-restore |
      | testfile.txt      |
    When the administrator restores the trashed resources "/folder-to-restore" of space "Personal" owned by user "Alice"
    Then the command should be successful
    And there should be "2" trashed resources of space "Personal" owned by user "Alice":
      | resource          | type   |
      | /testfile.txt     | file   |
      | /folder-to-delete | folder |


  Scenario: purge expired trashed resources of Personal space
    Given the following configs have been set:
      | config                                               | value |
      | STORAGE_USERS_PURGE_TRASH_BIN_PERSONAL_DELETE_BEFORE | 1s    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created the following folders
      | path              |
      | folder-to-delete  |
    And user "Alice" has deleted folder "/folder-to-delete"
    And user "Alice" has created a space "Newspace" with the default quota using the Graph API
    And user "Alice" has created a folder "uploadFolder" in space "Newspace"
    And user "Alice" has uploaded a file inside space "Newspace" with content "hello world" to "text.txt"
    And user "Alice" has removed the file "text.txt" from space "Newspace"
    And the administrator waits for "1" seconds
    When the administrator purges the expired trash resources
    Then there should be no trashed resources of space "Personal" owned by user "Alice"
    And there should be "1" trashed resources of space "Newspace" owned by user "Alice":
      | resource  | type   |
      | /text.txt | file |


  Scenario: purge expired trashed resources of project space
    Given the following configs have been set:
      | config                                              | value |
      | STORAGE_USERS_PURGE_TRASH_BIN_PROJECT_DELETE_BEFORE | 1s    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created the following folders
      | path              |
      | folder-to-delete  |
    And user "Alice" has deleted folder "/folder-to-delete"
    And user "Alice" has created a space "Newspace" with the default quota using the Graph API
    And user "Alice" has created a folder "uploadFolder" in space "Newspace"
    And user "Alice" has uploaded a file inside space "Newspace" with content "hello world" to "text.txt"
    And user "Alice" has removed the file "text.txt" from space "Newspace"
    And the administrator waits for "1" seconds
    When the administrator purges the expired trash resources
    Then there should be no trashed resources of space "Newspace" owned by user "Alice"
    And there should be "1" trashed resources of space "Personal" owned by user "Alice":
      | resource          | type   |
      | /folder-to-delete | folder |
