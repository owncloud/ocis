@env-config
Feature: backup consistency
  As a user
  I want to check my data for inconsistencies
  So that I can backup my data

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |


  Scenario: check backup consistency via CLI command
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
    And user "Alice" has created folder "/uploadFolder"
    And user "Carol" has created folder "/uploadFolder"
    And user "Alice" has sent the following resource share invitation:
      | resource        | uploadFolder |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Editor       |
    And user "Carol" has deleted file "/uploadFolder"
    And the administrator has stopped the server
    When the administrator checks the backup consistency using the CLI
    Then the command should be successful
    And the command output should contain "💚 No inconsistency found. The backup in '%storage_path%' seems to be valid."

  @issue-9498
  Scenario: check backup consistency after uploading a file multiple times via CLI command
    Given user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    And the administrator has stopped the server
    When the administrator checks the backup consistency using the CLI
    Then the command should be successful
    And the command output should contain "💚 No inconsistency found. The backup in '%storage_path%' seems to be valid."