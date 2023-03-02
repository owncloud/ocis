@api
Feature: assign role
  As an admin,
  I want to assign roles to users.
  So that users without an admin role cannot get the list of roles, assignments list and assign roles to users


  Scenario Outline: assign role to the user using graph api
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "<userRole>" to user "Alice" using the Graph API
    When the administrator retrieves the assigned role of user "Alice" using the Graph API
    Then the HTTP status code should be "200"
    And the Graph API response should have the role "<userRole>"
    Examples:
      | userRole    |
      | Admin       |
      | Space Admin |
      | User        |
      | Guest       |

  @issue-5032
  Scenario Outline: assign role to the user with graph api and list role with setting api
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "<userRole>" to user "Alice" using the Graph API
    When user "Alice" tries to get list of assignment
    Then the HTTP status code should be "<statusCode>"
    And the setting API response should have the role "<userRole>"
    Examples:
      | userRole    | statusCode |
      | Admin       | 201        |
      | Space Admin | 401        |
      | User        | 401        |
      | Guest       | 401        |
