@env-config
Feature: recalculate treesize using cli
  As an administrator
  I want to recalculate the treesize of spaces using cli
  So that I can repair a drifted quota usage

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: administrator recalculates the treesize of a personal space using the CS3 space id
    Given user "Alice" has uploaded file with content "test file" to "testfile.txt"
    And the administrator has stopped the server
    When the administrator recalculates the treesize of space "Personal" of user "Alice" using the CLI
    Then the command should be successful
    And the command output should contain "All treesizes are correct"


  Scenario: administrator recalculates the treesize of a project space using the opaque CS3 space id
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "Project" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project" with content "test file" to "testfile.txt"
    And the administrator has stopped the server
    When the administrator recalculates the treesize of space "Project" of user "Alice" using its opaque CS3 ID with the CLI
    Then the command should be successful
    And the command output should contain "All treesizes are correct"


  Scenario: administrator repairs a drifted treesize and the quota is correct after restart
    Given the config "STORAGE_USERS_FILEMETADATA_CACHE_STORE" has been set to "noop" for "storage-users" service
    And user "Alice" has uploaded file with content "test file" to "testfile.txt"
    And the administrator has stopped the server
    And the administrator has removed the storage link for file "testfile.txt" in space "Personal" of user "Alice"
    And the administrator has started the server
    And as user "Alice" folder "/" should contain a property "d:quota-used-bytes" with value "9"
    And the administrator has stopped the server
    When the administrator repairs the treesize of space "Personal" of user "Alice" using the CLI
    Then the command should be successful
    And the command output should contain "Corrected 1 treesizes"
    When the administrator starts the server
    Then as user "Alice" folder "/" should contain a property "d:quota-used-bytes" with value "0"


  Scenario: recalculating the treesize of a space that does not exist fails
    Given the administrator has stopped the server
    When the administrator recalculates the treesize of space with id "00000000-0000-0000-0000-000000000000" using the CLI
    Then the command should be unsuccessful
    And the command output should not contain "All treesizes are correct"
