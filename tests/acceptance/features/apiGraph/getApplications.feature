@api @skipOnOcV10
Feature: get applications
  As an user
  I want to be able to get applications information with existings roles

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: admin user lists all the groups
    Given the administrator has given "Alice" the role "<role>" using the settings api
    When user "Alice" gets all applications using the Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following applications information:
      | key                        | value                   |
      | displayName                | ownCloud Infinite Scale |
      | id                         | %uuid_v4%               |
    And the user retrieve API response should contain the following app roles:
      | Admin       |
      | Space Admin |
      | User        |
      | Guest       |
    Examples:
      | role        |
      | Admin       |
      | Space Admin |
      | User        |
      | Guest       |
