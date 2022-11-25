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
    When user "Alice" tries to get information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the API response should contain the following information:
      | displayName  | id       | mail              | onPremisesSamAccountName |
      | Brian Murphy | %UUIDv4% | brian@example.org | Brian                    |


  Scenario: non-admin user tries get information of a user
    Given user "Brian" has been created with default attributes and without skeleton files
    When user "Brian" tries to get information of user "Alice" using Graph API
    Then the HTTP status code should be "200"
    And the last response should be an unauthorized response


  Scenario: admin user tries get all user
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
      | David    |
    When user "Alice" tries to get all user using the Graph API
    Then the HTTP status code should be "200"
    And the API response should contain all user with following information:
      | displayName  | id       | mail              | onPremisesSamAccountName |
      | Brian Murphy | %UUIDv4% | brian@example.org | Brian                    |
      | David Lopez  | %UUIDv4% | david@example.org | David                    |
      | Carol King   | %UUIDv4% | carol@example.org | Carol                    |


  Scenario: non-admin user tries get all user
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
      | David    |
    When user "Brian" tries to get all user using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response

