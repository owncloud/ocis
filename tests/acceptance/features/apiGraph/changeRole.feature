Feature: change role
  As an admin
  I want to change the role of user
  So that I can manage the role of user

  Background:
    Given user "Alice" has been created with default attributes


  Scenario Outline: admin user changes the role of another user with different roles
    Given user "Brian" has been created with default attributes
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When user "Alice" changes the role of user "Brian" to role "<new-user-role>" using the Graph API
    Then the HTTP status code should be "201"
    And user "Brian" should have the role "<new-user-role>"
    Examples:
      | user-role   | new-user-role |
      | Admin       | Admin         |
      | Admin       | Space Admin   |
      | Admin       | User          |
      | Admin       | User Light    |
      | Space Admin | Admin         |
      | Space Admin | Space Admin   |
      | Space Admin | User          |
      | Space Admin | User Light    |
      | User        | Admin         |
      | User        | Space Admin   |
      | User        | User          |
      | User        | User Light    |
      | User Light  | Admin         |
      | User Light  | Space Admin   |
      | User Light  | User          |
      | User Light  | User Light    |


  Scenario Outline: admin user tries to change his/her own role
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When user "Alice" tries to change the role of user "Alice" to role "<new-user-role>" using the Graph API
    Then the HTTP status code should be "403"
    And user "Alice" should have the role "Admin"
    Examples:
      | new-user-role |
      | Space Admin   |
      | User          |
      | User Light    |
      | Admin         |


  Scenario Outline: non-admin cannot change the user role
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes
    When user "Alice" tries to change the role of user "Alice" to role "Admin" using the Graph API
    Then the HTTP status code should be "403"
    And user "Brian" should have the role "User"
    Examples:
      | user-role   |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario Outline: non-admin user tries to change the role of nonexistent user
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes
    When user "Alice" tries to change the role of user "nonexistent" to role "Admin" using the Graph API
    Then the HTTP status code should be "403"
    Examples:
      | user-role   |
      | Space Admin |
      | User        |
      | User Light  |
