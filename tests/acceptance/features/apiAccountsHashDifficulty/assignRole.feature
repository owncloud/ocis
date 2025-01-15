Feature: assign role
  As an admin,
  I want to assign roles to users
  So that I can provide them different authority


  Scenario Outline: user can see all existing roles
    Given user "Alice" has been created with default attributes
    And the administrator has given "Alice" the role "<user-role>" using the settings api
    When user "Alice" gets all applications using the Graph API
    Then the HTTP status code should be "<http-status-code>"
    Examples:
      | user-role   | http-status-code |
      | Admin       | 200              |
      | Space Admin | 200              |
      | User        | 200              |


  Scenario Outline: only admin user can see assignments list
    Given user "Alice" has been created with default attributes
    And the administrator has given "Alice" the role "<user-role>" using the settings api
    When user "Alice" tries to get the assigned role of user "Alice" using the Graph API
    Then the HTTP status code should be "<http-status-code>"
    Examples:
      | user-role   | http-status-code |
      | Admin       | 200              |
      | Space Admin | 403              |
      | User        | 403              |


  Scenario Outline: a user cannot change own role
    Given user "Alice" has been created with default attributes
    And the administrator has given "Alice" the role "<user-role>" using the settings api
    When user "Alice" tries to change the role of user "Alice" to role "<desired-role>" using the Graph API
    Then the HTTP status code should be "403"
    And user "Alice" should have the role "<user-role>"
    Examples:
      | user-role   | desired-role |
      | Admin       | User         |
      | Admin       | Space Admin  |
      | Space Admin | Admin        |
      | Space Admin | Space Admin  |
      | User        | Admin        |
      | User        | Space Admin  |


  Scenario Outline: only admin user can change the role for another user
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And the administrator has given "Alice" the role "<user-role>" using the settings api
    When user "Alice" changes the role of user "Brian" to role "<desired-role>" using the Graph API
    Then the HTTP status code should be "<http-status-code>"
    And user "Brian" should have the role "<expected-role>"
    Examples:
      | user-role   | desired-role | http-status-code | expected-role |
      | Admin       | User         | 201              | User          |
      | Admin       | Space Admin  | 201              | Space Admin   |
      | Admin       | Admin        | 201              | Admin         |
      | Space Admin | Admin        | 403              | User          |
      | Space Admin | Space Admin  | 403              | User          |
      | User        | Admin        | 403              | User          |
      | User        | Space Admin  | 403              | User          |
