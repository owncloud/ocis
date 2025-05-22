@env-config
Feature: trashbin
  As a administrator
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
