Feature: assign role
  As an admin,
  I want to assign roles to users
  So that I can provide them different authority


  Scenario Outline: only admin user can see all existing roles
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "<userRole>" using the settings api
    When user "Alice" tries to get all existing roles
    Then the HTTP status code should be "<statusCode>"
    Examples:
      | userRole    | statusCode |
      | Admin       | 201        |
      | Space Admin | 201        |
      | User        | 201        |

  @issue-5032
  Scenario Outline: only admin user can see assignments list
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "<userRole>" using the settings api
    When user "Alice" tries to get list of assignment
    Then the HTTP status code should be "<statusCode>"
    Examples:
      | userRole    | statusCode |
      | Admin       | 201        |
      | Space Admin | 401        |
      | User        | 401        |


  Scenario Outline: a user cannot change own role
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "<userRole>" using the settings api
    When user "Alice" changes his own role to "<desiredRole>"
    Then the HTTP status code should be "400"
    And user "Alice" should have the role "<userRole>"
    Examples:
      | userRole    | desiredRole |
      | Admin       | User        |
      | Admin       | Space Admin |
      | Space Admin | Admin       |
      | Space Admin | Space Admin |
      | User        | Admin       |
      | User        | Space Admin |


  Scenario Outline: only admin user can change the role for another user
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And the administrator has given "Alice" the role "<userRole>" using the settings api
    When user "Alice" changes the role "<desiredRole>" for user "Brian"
    Then the HTTP status code should be "<statusCode>"
    And user "Brian" should have the role "<expectedRole>"
    Examples:
      | userRole    | desiredRole | statusCode | expectedRole |
      | Admin       | User        | 201        | User         |
      | Admin       | Space Admin | 201        | Space Admin  |
      | Admin       | Admin       | 201        | Admin        |
      | Space Admin | Admin       | 400        | User         |
      | Space Admin | Space Admin | 400        | User         |
      | User        | Admin       | 400        | User         |
      | User        | Space Admin | 400        | User         |
