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


  Scenario: recalculating the treesize of a space that does not exist fails
    Given the administrator has stopped the server
    When the administrator recalculates the treesize of space with id "00000000-0000-0000-0000-000000000000" using the CLI
    Then the command should be unsuccessful
    And the command output should not contain "All treesizes are correct"
