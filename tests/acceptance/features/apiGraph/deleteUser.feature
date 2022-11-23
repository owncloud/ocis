@api @skipOnOcV10
Feature: delete user
  Only user with admin permissions can delete user

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: the admin user deletes a user
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api
    When the user "Alice" deletes a user "Brian" using GraphAPI
    Then the HTTP status code should be "204"
    And user "Brian" should not exist


  Scenario: the admin user tries to delete self
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When the user "Alice" deletes a user "Alice" using GraphAPI
    Then the HTTP status code should be "403"
    And user "Alice" should exist


  Scenario: the admin user tries to delete non existent user
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When the user "Alice" deletes a user "nonExistentUser" using GraphAPI
    Then the HTTP status code should be "404"


  Scenario Outline: a user without admin right cannot delete a user
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "<role>" using the settings api
    When the user "Alice" deletes a user "Brian" using GraphAPI
    Then the HTTP status code should be "401"
    And user "Brian" should exist
    Examples:
      | role        |
      | Space Admin |
      | User        |
