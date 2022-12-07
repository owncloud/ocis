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
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    |


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
      | displayName  | id        | mail              | onPremisesSamAccountName |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    |
      | David Lopez  | %uuid_v4% | david@example.org | David                    |
      | Carol King   | %uuid_v4% | carol@example.org | Carol                    |


  Scenario: non-admin user tries get all user
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
      | David    |
    When user "Brian" tries to get all user using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response


  Scenario: admin user tries to get drive information of a user
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
    When the user "Alice" tries to get user "Brian" along with his drive information using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName |
      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    |
    And the user retrieve API response should contain the following drive information:
#      | driveType         | personal                         |
#      | driveAlias        | personal/brian                   |
#      | id                | %space_id%                       |
#      | name              | Brian Murphy                     |
      | owner@@@user@@@id | %user_id%                        |
#      | quota@@@state     | normal                           |
#      | root@@@id         | %space_id%                       |
#      | root@@@webDavUrl  | %base_url%/dav/spaces/%space_id% |
#      | webUrl            | %base_url%/f/%space_id%          |


#  Scenario: normal user tries to get hid/her own drive information
#    Given these users have been created with default attributes and without skeleton files:
#      | username |
#      | Brian    |
#    When the user "Brian" tries to get his drive information using Graph API
#    Then the HTTP status code should be "200"
#    And the user retrieve API response should contain the following information:
#      | displayName  | id        | mail              | onPremisesSamAccountName |
#      | Brian Murphy | %uuid_v4% | brian@example.org | Brian                    |
#    And the response should contain the following drive information:
#      | driveType        | personal                         |
#      | driveAlias       | personal/brian                   |
#      | name             | Brian Murphy                     |
