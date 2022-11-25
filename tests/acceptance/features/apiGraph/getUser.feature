@api @skipOnOcV10
Feature: get users
  As an admin
  I want to be able to retrieve user information
  So that I can see the information

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario: admin user tries get information of a user
    Given user "Brian" has been created with default attributes and without skeleton files
    When user "Alice" tries to get information of user "Brian"
    Then the HTTP status code should be "200"
    And the API response should contain the following information:
      | displayName              | Brian Murphy      |
      | id                       | %UUIDv4%          |
      | mail                     | brian@example.org |
      | onPremisesSamAccountName | Brian             |


  Scenario: non-admin user tries get information of a user
    Given user "Brian" has been created with default attributes and without skeleton files
    When user "Brian" tries to get information of user "Alice"
    Then the HTTP status code should be "200"
    And the last response should be an unauthorized response
