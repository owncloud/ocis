Feature: get applications
  As a user
  I want to be able to get application information with existing roles
  So that I can see which role belongs to what user

  Background:
    Given user "Alice" has been created with default attributes


  Scenario Outline: admin user lists all the groups
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When user "Alice" gets all applications using the Graph API
    Then the HTTP status code should be "200"
    And the user API response should contain the following application information:
      | key         | value                   |
      | displayName | ownCloud Infinite Scale |
      | id          | %uuid_v4%               |
    And the user API response should contain the following app roles:
      | Admin       |
      | Space Admin |
      | User        |
      | User Light  |
    Examples:
      | user-role   |
      | Admin       |
      | Space Admin |
      | User        |
      | User Light  |
