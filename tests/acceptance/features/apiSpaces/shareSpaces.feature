@api @skipOnOcV10
Feature: Share spaces
  As the owner of a space
  I want to be able to add members to a space, and to remove access for them

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "share space" with the default quota using the GraphApi


  Scenario Outline:: A user can share a space to another user
    When user "Alice" shares a space "share space" to user "Brian" with role "<role>"
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    And the user "Brian" should have a space called "share space" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | id        | %space_id%  |
      | name      | share space |
    Examples:
      | role    |
      | manager |
      | editor  |
      | viewer  |


  Scenario: A user can see who has been granted access
    Given user "Alice" has shared a space "share space" to user "Brian" with role "viewer"
    And the user "Alice" should have a space called "share space" granted to "Brian" with these key and value pairs:
      | key                                                | value     |
      | root@@@permissions@@@1@@@grantedTo@@@0@@@user@@@id | %user_id% |
      | root@@@permissions@@@1@@@roles@@@0                 | viewer    |


  Scenario: A user can see a file in a received shared space
    Given user "Alice" has uploaded a file inside space "share space" with content "Test" to "test.txt"
    And user "Alice" has created a folder "Folder Main" in space "share space"
    When user "Alice" shares a space "share space" to user "Brian" with role "viewer"
    Then for user "Brian" the space "share space" should contain these entries:
      | test.txt    |
      | Folder Main |


  Scenario: When a user unshares a space, the space becomes unavailable to the receiver
    Given user "Alice" has shared a space "share space" to user "Brian" with role "viewer"
    And the user "Brian" should have a space called "share space" with these key and value pairs:
      | key       | value       |
      | driveType | project     |
      | id        | %space_id%  |
      | name      | share space |
    When user "Alice" unshares a space "share space" to user "Brian"
    Then the HTTP status code should be "200"
    Then the user "Brian" should not have a space called "share space"


  Scenario: A user can add another user to the space managers to enable him
    Given user "Alice" has uploaded a file inside space "share space" with content "Test" to "test.txt"
    When user "Alice" shares a space "share space" to user "Brian" with role "manager"
    Then the user "Brian" should have a space called "share space" granted to "Brian" with role "manager"
    When user "Brian" shares a space "share space" to user "Bob" with role "viewer"
    Then the user "Bob" should have a space called "share space" granted to "Bob" with role "viewer"
    And for user "Bob" the space "share space" should contain these entries:
      | test.txt |
