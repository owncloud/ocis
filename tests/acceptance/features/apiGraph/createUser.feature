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
      | userName    | <userName>    |
      | displayName | <displayName> |
      | email       | <email>       |
      | password    | <password>    |
    Then the HTTP status code should be "<code>"
    And user "<userName>" <shouldOrNot> exist
    Examples:
      | userName                     | displayName      | email               | password                     | code | shouldOrNot |
      | SameDisplayName              | Alice Hansen     | new@example.org     | containsCharacters(*:!;_+-&) | 200  | should     |
      | withoutPassSameEmail         | without pass     | alice@example.org   |                              | 200  | should     |
      | name                         | pass with space  | example@example.org | my pass                      | 200  | should     |
      | nameWithCharacters(*:!;_+-&) | user             | new@example.org     | 123                          | 400  | should not |
      | withoutEmail                 | without email    |                     | 123                          | 200  | should     |
      | Alice                        | same userName    | new@example.org     | 123                          | 400  | should     |
      | name with space              | name with space  | example@example.org | 123                          | 400  | should not |


  Scenario: a user cannot be created with empty name
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName    |              |
      | displayName | emptyName    |
      | email       | @example.org |
      | password    | 123          |
    Then the HTTP status code should be "400"


  Scenario Outline: a user without admin right cannot create a user
    Given the administrator has given "Alice" the role "<role>" using the settings api
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName    | user         |
      | displayName | user         |
      | email       | @example.org |
      | password    | 123          |
    Then the HTTP status code should be "401"
    And user "user" should not exist
    Examples:
      | role        |
      | Space Admin |
      | User        |

  Scenario Outline: only user with admin role can create user
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "<userRole>" using the Graph API
    When the user "Alice" creates a new user using GraphAPI with the following settings:
      | userName    | <userName>    |
      | displayName | <displayName> |
      | email       | <email>       |
      | password    | <password>    |
    Then the HTTP status code should be "<code>"
    And user "<userName>" <shouldOrNot> exist
    Examples:
      | userRole    | userName          | displayName     | email               | password                     | code | shouldOrNot |
      | Admin       | SameDisplayName   | Alice Hansen    | new@example.org     | containsCharacters(*:!;_+-&) | 200  | should      |
      | Space Admin | withPassSameEmail | with pass       | alice@example.org   | $-ad#                        | 401  | should not  |
      | User        | name              | pass with space | example@example.org | my pass                      | 401  | should not  |
      | Guest       | hari kumar        | hary            | hari@example.com    | 123                          | 401  | should not  |
