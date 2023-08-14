Feature: change role
  As an admin
  I want to change the role of user
  So that I can manage the role of user

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: admin user changes the role of another user with different roles
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the administrator has assigned the role "<userRole>" to user "Brian" using the Graph API
    When user "Alice" changes the role of user "Brian" to role "<newRole>" using the Graph API
    Then the HTTP status code should be "201"
    And user "Brian" should have the role "<newRole>"
    Examples:
      | userRole    | newRole     |
      | Admin       | Admin       |
      | Admin       | Space Admin |
      | Admin       | User        |
      | Admin       | User Light  |
      | Space Admin | Admin       |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | User Light  |
      | User        | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | User Light  |
      | User Light  | Admin       |
      | User Light  | Space Admin |
      | User Light  | User        |
      | User Light  | User Light  |


  Scenario Outline: admin user tries to change his/her own role
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When user "Alice" tries to change the role of user "Alice" to role "<newRole>" using the Graph API
    Then the HTTP status code should be "403"
    And user "Alice" should have the role "Admin"
    Examples:
      | newRole     |
      | Space Admin |
      | User        |
      | User Light  |
      | Admin       |


  Scenario Outline: non-admin cannot change the user role
    Given the administrator has assigned the role "<role>" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes and without skeleton files
    When user "Alice" tries to change the role of user "Alice" to role "Admin" using the Graph API
    Then the HTTP status code should be "401"
    And user "Brian" should have the role "User"
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | User Light  |
