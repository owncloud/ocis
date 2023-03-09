@api @skipOnOcV10
Feature: create user
  Only user with admin permissions can create new user

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: the admin creates a user
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       | <userName>    |
      | displayName    | <displayName> |
      | email          | <email>       |
      | password       | <password>    |
      | accountEnabled | <enable>      |
    Then the HTTP status code should be "<code>"
    And user "<userName>" <shouldOrNot> exist
    Examples:
      | userName                     | displayName     | email               | password                     | code | enable | shouldOrNot |
      | SameDisplayName              | Alice Hansen    | new@example.org     | containsCharacters(*:!;_+-&) | 200  | true   | should      |
      | withoutPassSameEmail         | without pass    | alice@example.org   |                              | 200  | true   | should      |
      | name                         | pass with space | example@example.org | my pass                      | 200  | true   | should      |
      | nameWithCharacters(*:!;_+-&) | user            | new@example.org     | 123                          | 400  | true   | should not  |
      | withoutEmail                 | without email   |                     | 123                          | 200  | true   | should      |
      | Alice                        | same userName   | new@example.org     | 123                          | 400  | true   | should      |
      | name with space              | name with space | example@example.org | 123                          | 400  | true   | should not  |
      | createDisabledUser           | disabled user   | example@example.org | 123                          | 200  | false  | should      |


  Scenario: a user cannot be created with empty name
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       |              |
      | displayName    | emptyName    |
      | email          | @example.org |
      | password       | 123          |
      | accountEnabled | true         |
    Then the HTTP status code should be "400"


  Scenario Outline: a user without admin right cannot create a user
    Given the administrator has given "Alice" the role "<role>" using the settings api
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       | user         |
      | displayName    | user         |
      | email          | @example.org |
      | password       | 123          |
      | accountEnabled | true         |
    Then the HTTP status code should be "401"
    And user "user" should not exist
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |


  Scenario: a user cannot be created with the name of the disabled user
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api
    And the user "Alice" has disabled user "Brian" using the Graph API
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       | Brian                 |
      | displayName    | This is another Brian |
      | email          | brian@example.com     |
      | password       | 123                   |
      | accountEnabled | true                  |
    Then the HTTP status code should be "400"


  Scenario: a user can be created with the name of the deleted user
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api
    And the user "Alice" has deleted a user "Brian" using the Graph API
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName       | Brian                 |
      | displayName    | This is another Brian |
      | email          | brian@example.com     |
      | password       | 123                   |
      | accountEnabled | true                  |
    Then the HTTP status code should be "200"
    And user "Brian" should exist
