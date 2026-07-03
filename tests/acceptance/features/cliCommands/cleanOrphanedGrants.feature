@env-config
Feature: clean orphaned grants using CLI
  As an administrator
  I want to clean orphaned share-manager grants using the CLI
  So that I can remove orphaned shares in the system

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And the administrator has configured service account credentials
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "orphan-test" with the default quota using the Graph API
    And using spaces DAV path
    And user "Alice" has uploaded a file inside space "orphan-test" with content "test content" to "orphanFile.txt"
    And using new DAV path
    And user "Alice" has sent the following resource share invitation:
      | resource        | orphanFile.txt |
      | space           | orphan-test    |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | Viewer         |
    And the administrator has stopped the server
    And the share grants for space "orphan-test" owned by "Alice" have been orphaned
    And the administrator has started the server


  Scenario: administrator runs clean-orphaned-grants in non-dry-run mode
    When the administrator runs clean-orphaned-grants in non-dry-run mode
    Then the command should be successful
    And the command output should contain "mode: dry run disabled: grants may be changed"
    And the command output should contain "Summary:"
    And the command output should contain "Orphans: 1 candidate(s), 1 removed"
    And the command output should contain "Reverse orphans: 0 candidate(s), 0 removed"


  Scenario: administrator runs clean-orphaned-grants with space-id filter
    When the administrator runs clean-orphaned-grants for space "orphan-test" owned by "Alice"
    Then the command should be successful
    And the command output should contain "scope: limiting scan to space"
    And the command output should contain "1 target space(s)"
    And the command output should contain "Summary:"
    And the command output should contain "Orphans: 1 candidate(s), 1 removed"
    And the command output should contain "Reverse orphans: 0 candidate(s), 0 removed"


  Scenario: administrator runs clean-orphaned-grants with force flag
    When the administrator runs clean-orphaned-grants with force flag
    Then the command should be successful
    And the command output should contain "flags: --force active"
    And the command output should contain "Summary:"
    And the command output should contain "Orphans: 1 candidate(s), 1 removed"
    And the command output should contain "Reverse orphans: 0 candidate(s), 0 removed"


  Scenario: administrator runs clean-orphaned-grants on a space with shares
    When the administrator runs clean-orphaned-grants in non-dry-run mode
    Then the command should be successful
    And the command output should contain "== Primary scan =="
    And the command output should contain "Summary:"
    And the command output should contain "Orphans: 1 candidate(s), 1 removed"
    And the command output should contain "Reverse orphans: 0 candidate(s), 0 removed"
