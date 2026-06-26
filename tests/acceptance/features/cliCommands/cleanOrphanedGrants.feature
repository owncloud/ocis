@env-config
Feature: clean orphaned grants using CLI
  As an administrator
  I want to clean orphaned share-manager grants using the CLI
  So that I can fix orphaned shares in the system

  Background:
    Given user "Alice" has been created with default attributes
    And the administrator has configured service account credentials


  Scenario: administrator runs clean-orphaned-grants in dry-run mode
    When the administrator runs clean-orphaned-grants in dry-run mode
    Then the command should be successful
    And the command output should contain "== Pre-flight =="
    And the command output should contain "mode: dry run enabled"
    And the command output should contain "== Primary scan =="
    And the command output should contain "Dry run mode: no grants were modified"
    And the command output should contain "== Reverse orphan scan =="


  Scenario: administrator runs clean-orphaned-grants in non-dry-run mode
    When the administrator runs clean-orphaned-grants in non-dry-run mode
    Then the command should be successful
    And the command output should contain "mode: dry run disabled: grants may be changed"


  Scenario: administrator runs clean-orphaned-grants with space-id filter
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project1" with the default quota using the Graph API
    When the administrator runs clean-orphaned-grants for space "project1" owned by "Alice"
    Then the command should be successful
    And the command output should contain "scope: limiting scan to space"


  Scenario: administrator runs clean-orphaned-grants with force flag
    When the administrator runs clean-orphaned-grants with force flag
    Then the command should be successful
    And the command output should contain "flags: --force active"
