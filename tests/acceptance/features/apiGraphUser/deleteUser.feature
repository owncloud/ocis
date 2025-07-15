Feature: delete user
  As an admin
  I want to be able to delete users
  So that I can remove unnecessary users

  Background:
    Given user "Alice" has been created with default attributes

  @issue-11027
  Scenario Outline: admin user deletes a user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the user "Alice" has created a new user with the following attributes:
      | userName    | <user-name>    |
      | displayName | <display-name> |
      | email       | <email>        |
      | password    | <password>     |
    When the user "Alice" deletes a user "<user-name>" using the Graph API
    Then the HTTP status code should be "204"
    And user "<user-name>" should not exist
    Examples:
      | user-name            | display-name    | email               | password                     |
      | SameDisplayName      | Alice Hansen    | new@example.org     | containsCharacters(*:!;_+-&) |
      | withoutPassSameEmail | without pass    | alice@example.org   |                              |
      | name                 | pass with space | example@example.org | my pass                      |
      | user+name            | Alex Hudson     | demo@example.org    | root                         |


  Scenario: delete a user and specify the user name in different case
    Given user "brand-new-user" has been created with default attributes
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When the user "Alice" deletes a user "Brand-New-User" using the Graph API
    Then the HTTP status code should be "204"
    And user "brand-new-user" should not exist


  Scenario Outline: admin user deletes another user with different role
    Given user "Brian" has been created with default attributes
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When the user "Alice" deletes a user "Brian" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should not exist
    Examples:
      | user-role   |
      | Admin       |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: admin user tries to delete his/her own account
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When the user "Alice" deletes a user "Alice" using the Graph API
    Then the HTTP status code should be "403"
    And user "Alice" should exist


  Scenario Outline: non-admin user tries to delete his/her own account
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When the user "Alice" deletes a user "Alice" using the Graph API
    Then the HTTP status code should be "403"
    And user "Alice" should exist
    Examples:
      | user-role   |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: admin user tries to delete a nonexistent user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When the user "Alice" tries to delete a nonexistent user using the Graph API
    Then the HTTP status code should be "404"


  Scenario Outline: non-admin user tries to delete a nonexistent user
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When the user "Alice" tries to delete a nonexistent user using the Graph API
    Then the HTTP status code should be "403"
    Examples:
      | user-role   |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario Outline: non-admin user tries to delete another user with different role
    Given user "Brian" has been created with default attributes
    And the administrator has assigned the role "<user-role-2>" to user "Brian" using the Graph API
    And the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When the user "Alice" deletes a user "Brian" using the Graph API
    Then the HTTP status code should be "403"
    And user "Brian" should exist
    Examples:
      | user-role   | user-role-2 |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | User Light  |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | User Light  |
      | User        | Admin       |
      | User Light  | Space Admin |
      | User Light  | User        |
      | User Light  | User Light  |
      | User Light  | Admin       |


  Scenario: admin user deletes a disabled user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes
    And the user "Alice" has disabled user "Brian"
    When the user "Alice" deletes a user "Brian" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should not exist


  Scenario Outline: normal user tries to delete a disabled user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes
    And user "Carol" has been created with default attributes
    And the administrator has assigned the role "<user-role-2>" to user "Brian" using the Graph API
    And the administrator has assigned the role "<user-role>" to user "Carol" using the Graph API
    And the user "Alice" has disabled user "Brian"
    When the user "Carol" deletes a user "Brian" using the Graph API
    Then the HTTP status code should be "403"
    And user "Brian" should exist
    Examples:
      | user-role   | user-role-2 |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | User Light  |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | User Light  |
      | User        | Admin       |
      | User Light  | Space Admin |
      | User Light  | User        |
      | User Light  | User Light  |
      | User Light  | Admin       |


  Scenario: personal space is deleted automatically when the user is deleted
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes
    When the user "Alice" deletes a user "Brian" using the Graph API
    Then the HTTP status code should be "204"
    When user "Alice" lists all spaces via the Graph API with query "$filter=driveType eq 'personal'"
    Then the json response should not contain a space with name "Brian Murphy"


  Scenario: accepted share is deleted automatically when the user is deleted
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes
    And user "Brian" has created folder "new"
    And user "Brian" has sent the following resource share invitation:
      | resource        | new      |
      | space           | Personal |
      | sharee          | Alice    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    When the user "Alice" deletes a user "Brian" using the Graph API
    Then the HTTP status code should be "204"
    And as "Alice" folder "Shares/new" should not exist
