@env-config
Feature: delete empty trash bin folder via CLI command


  Scenario: delete empty trashbin folders
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | Alice        |
      | displayName | Alice Hansen |
      | password    | %alt1%       |
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
    And user "Alice" has restored the folder with original path "folder-to-restore"
    And user "Alice" has deleted the folder with original path "folder-to-delete" from the trashbin
    And the administrator has stopped the server
    When the administrator deletes the empty trashbin folders using the CLI
    Then the command should be successful
