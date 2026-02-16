@env-config
Feature: delete space using cli
  As an administrator
  I want to delete spaces of users using cli
  So that I can manage them

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: administrator deletes personal space of users
    Given user "Admin" has disabled personal space of user "Alice"
    When administrator deletes "Personal" space using the CLI
    Then the command should be successful
    And the command output should contain "Purge completed. Purged 1 spaces"


  Scenario: administrator deletes project spaces of users
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project1" with the default quota using the Graph API
    And user "Alice" has created a space "project2" with the default quota using the Graph API
    And user "Alice" has disabled a space "project1"
    And user "Alice" has disabled a space "project2"
    When administrator deletes "Project" space using the CLI
    Then the command should be successful
    And the command output should contain "Purge completed. Purged 2 spaces"


  Scenario: administrator deletes personal space of users using space-id (Personal)
    Given user "Admin" has disabled personal space of user "Alice"
    When administrator deletes "Personal" space of user "Alice" with space-id using the CLI
    Then the command should be successful
    And the command output should contain "Purge completed. Purged 1 spaces"


  Scenario: administrator deletes project space of users using space-id (Project)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project" with the default quota using the Graph API
    And user "Alice" has disabled a space "project"
    When administrator deletes "project" space of user "Alice" with space-id using the CLI
    Then the command should be successful
    And the command output should contain "Purge completed. Purged 1 spaces"


  Scenario: administrator deletes spaces beyond the retention-period
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project" with the default quota using the Graph API
    And user "Alice" has disabled a space "project"
    And the administrator has waited for "5s" seconds
    And user "Alice" has created a space "project1" with the default quota using the Graph API
    And user "Alice" has disabled a space "project1"
    When administrator deletes "Project" space with "2s" retention period using the CLI
    Then the command should be successful
    And the command output should contain "Purge completed. Purged 1 spaces"


  Scenario: administrator tries to delete enabled spaces
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project" with the default quota using the Graph API
    When administrator deletes "project" space of user "Alice" with space-id using the CLI
    Then the command should be successful
    And the command output should contain "Purge completed. Purged 0 spaces"

