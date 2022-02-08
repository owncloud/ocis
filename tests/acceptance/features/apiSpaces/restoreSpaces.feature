@api @skipOnOcV10
Feature: Restoring space
  As a manager of space
  I want to be able to restore a disabled space.
  Only manager can restore disabled space
  The restored space must be visible to the other participants without loss of data

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario: An owner can restore a Space via the Graph API
    Given user "Alice" has created a space "restore a space" of type "project" with quota "10"
    And user "Alice" has disabled a space "restore a space"
    When user "Alice" restores a disabled space "restore a space"
    Then the HTTP status code should be "200"


  Scenario: Participants can see the data after the space is restored
    Given user "Alice" has created a space "data exists" of type "project" with quota "10"
    And user "Alice" has created a folder "mainFolder" in space "data exists"
    And user "Alice" has uploaded a file inside space "data exists" with content "example" to "test.txt"
    And user "Alice" has shared a space "data exists" to user "Brian" with role "editor"
    And user "Alice" has shared a space "data exists" to user "Bob" with role "viewer"
    And user "Alice" has disabled a space "data exists"
    When user "Alice" restores a disabled space "data exists"
    Then for user "Alice" the space "data exists" should contain these entries:
      | test.txt         |
      | mainFolder       |
    And for user "Brian" the space "data exists" should contain these entries:
      | test.txt         |
      | mainFolder       |
    And for user "Bob" the space "data exists" should contain these entries:
      | test.txt         |
      | mainFolder       |


  Scenario: Participants can create data in the space after restoring
    Given user "Alice" has created a space "create data in restored space" of type "project" with quota "10"
    And user "Alice" has shared a space "create data in restored space" to user "Brian" with role "editor"
    And user "Alice" has disabled a space "create data in restored space"
    And user "Alice" has restored a disabled space "create data in restored space"
    When user "Brian" creates a folder "mainFolder" in space "create data in restored space" using the WebDav Api
    And user "Brian" uploads a file inside space "create data in restored space" with content "test" to "test.txt" using the WebDAV API
    Then for user "Brian" the space "create data in restored space" should contain these entries:
      | test.txt         |
      | mainFolder       |
    