@api @skipOnOcV10
Feature: Disabling and deleting space
  As a manager of space
  I want to be able to disable the space first, then delete it.
  I want to make sure that a disabled spaces isn't accessible by shared users. 

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario: An owner can disable a Space via the Graph API
    Given user "Alice" has created a space "disable a space" of type "project" with quota "10"
    And user "Alice" has shared a space "disable a space" to user "Brian" with role "editor"
    When user "Alice" disables a space "disable a space"
    Then the HTTP status code should be "204"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "disable a space" with these key and value pairs:
      | key  | value           |
      | name | disable a space |
    When user "Brian" lists all available spaces via the GraphApi
    Then the json responded should not contain a space with name "disable a space"


  Scenario: An owner can delete a disabled Space via the webDav API
    Given user "Alice" has created a space "delete a space" of type "project" with quota "10"
    And user "Alice" has disabled a space "delete a space"
    When user "Alice" deletes a space "delete a space"
    Then the HTTP status code should be "204"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should not contain a space with name "delete a space"
    

  Scenario: An owner can disable and delete Space in which files and folders exist via the webDav API
    Given user "Alice" has created a space "delete a space with content" of type "project" with quota "10"
    And user "Alice" has uploaded a file inside space "delete a space with content" with content "test" to "test.txt"
    And user "Alice" has created a folder "MainFolder" in space "delete a space with content"
    When user "Alice" disables a space "delete a space with content"
    Then the HTTP status code should be "204"
    When user "Alice" deletes a space "delete a space with content"
    Then the HTTP status code should be "204"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should not contain a space with name "delete a space with content"


  Scenario: An owner cannot delete a space via the webDav API without first disabling it
    Given user "Alice" has created a space "delete without disabling" of type "project" with quota "10"
    When user "Alice" deletes a space "delete without disabling"
    Then the HTTP status code should be "400"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "delete without disabling" with these key and value pairs:
      | key  | value                    |
      | name | delete without disabling |


  Scenario: An user with editor role cannot disable a Space via the Graph API
    Given user "Alice" has created a space "editor tries to disable a space" of type "project" with quota "10"
    And user "Alice" has shared a space "editor tries to disable a space" to user "Brian" with role "editor"
    When user "Brian" disables a space "editor tries to disable a space"
    Then the HTTP status code should be "403"
    When user "Brian" lists all available spaces via the GraphApi
    Then the json responded should contain a space "editor tries to disable a space" with these key and value pairs:
      | key  | value                           |
      | name | editor tries to disable a space |
    

  Scenario: An user with viewer role cannot disable a Space via the Graph API
    Given user "Alice" has created a space "viewer tries to disable a space" of type "project" with quota "10"
    And user "Alice" has shared a space "viewer tries to disable a space" to user "Brian" with role "viewer"
    When user "Brian" disables a space "viewer tries to disable a space"
    Then the HTTP status code should be "403"
    When user "Brian" lists all available spaces via the GraphApi
    Then the json responded should contain a space "viewer tries to disable a space" with these key and value pairs:
      | key  | value                           |
      | name | viewer tries to disable a space |
