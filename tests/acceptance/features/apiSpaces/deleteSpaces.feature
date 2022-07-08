@api @skipOnOcV10
Feature: Disabling and deleting space
  As a manager of space
  I want to be able to disable the space first, then delete it.
  I want to make sure that a disabled spaces isn't accessible by shared users.

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "Project Moon" with the default quota using the GraphApi
    And user "Alice" has shared a space "Project Moon" to user "Brian" with role "editor"
    And user "Alice" has shared a space "Project Moon" to user "Bob" with role "viewer"


  Scenario Outline: A space admin user can disable a Space via the Graph API
    When user "Alice" disables a space "Project Moon"
    Then the HTTP status code should be "204"
    And the user "Alice" should have a space called "Project Moon" with these key and value pairs:
      | key                    | value        |
      | name                   | Project Moon |
      | root@@@deleted@@@state | trashed      |
    And the user "<user>" should not have a space called "Project Moon"
    Examples:
      | user  |
      | Brian |
      | Bob   |


  Scenario Outline: An user without space admin role cannot disable a Space via the Graph API
    When user "<user>" disables a space "Project Moon"
    Then the HTTP status code should be "403"
    And the user "<user>" should have a space called "Project Moon" with these key and value pairs:
      | key  | value        |
      | name | Project Moon |
    Examples:
      | user  |
      | Brian |
      | Bob   |


  Scenario: A space manager can delete a disabled Space via the webDav API
    Given user "Alice" has disabled a space "Project Moon"
    When user "Alice" deletes a space "Project Moon"
    Then the HTTP status code should be "204"
    And the user "Alice" should not have a space called "Project Moon"


  Scenario: An space manager can disable and delete Space in which files and folders exist via the webDav API
    Given user "Alice" has uploaded a file inside space "Project Moon" with content "test" to "test.txt"
    And user "Alice" has created a folder "MainFolder" in space "Project Moon"
    When user "Alice" disables a space "Project Moon"
    Then the HTTP status code should be "204"
    When user "Alice" deletes a space "Project Moon"
    Then the HTTP status code should be "204"
    And the user "Alice" should not have a space called "Project Moon"


  Scenario: An space manager cannot delete a space via the webDav API without first disabling it
    When user "Alice" deletes a space "Project Moon"
    Then the HTTP status code should be "400"
    And the user "Alice" should have a space called "Project Moon" with these key and value pairs:
      | key  | value        |
      | name | Project Moon |
