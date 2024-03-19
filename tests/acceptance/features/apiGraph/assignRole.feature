Feature: assign role
  As an admin,
  I want to assign roles to users.
  So that users without an admin role cannot get the list of roles, assignments list and assign roles to users


  Scenario Outline: assign role to the user using graph api
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When the administrator retrieves the assigned role of user "Alice" using the Graph API
    Then the HTTP status code should be "200"
    And the Graph API response should have the role "<user-role>"
    Examples:
      | user-role   |
      | Admin       |
      | Space Admin |
      | User        |
      | User Light  |

  @issue-5032
  Scenario Outline: assign role to the user with graph api and list role with setting api
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When user "Alice" tries to get list of assignment
    Then the HTTP status code should be "<http-status-code>"
    And the setting API response should have the role "<user-role>"
    Examples:
      | user-role   | http-status-code |
      | Admin       | 201              |
      | Space Admin | 401              |
      | User        | 401              |
      | User Light  | 401              |


  Scenario Outline: assign role to the user with setting api and list role with graph api
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "<user-role>" using the settings api
    When the administrator retrieves the assigned role of user "Alice" using the Graph API
    Then the HTTP status code should be "200"
    And the Graph API response should have the role "<user-role>"
    Examples:
      | user-role   |
      | Admin       |
      | Space Admin |
      | User        |
      | User Light  |
