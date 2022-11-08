@api @skipOnOcV10
Feature: get groups and their members
  As an admin
  I want to be able to get groups
  So that I can see all the groups and their members

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario: admin user lists all the groups
    Given group "tea-lover" has been created
    And group "coffee-lover" has been created
    And group "h2o-lover" has been created
    When user "Alice" gets all the groups using the Graph API
    Then the extra groups returned by the API should be
      | tea-lover    |
      | coffee-lover |
      | h2o-lover    |


  Scenario: normal user cannot get the groups list
    Given user "Brian" has been created with default attributes and without skeleton files
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And group "h2o-lover" has been created
    When user "Brian" gets all the groups using the Graph API
    Then the HTTP status code should be "401"


  Scenario: admin user gets users of a group
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
    And group "tea-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    And user "Carol" has been added to group "tea-lover"
    When user "Alice" gets all the members of group "tea-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the users returned by the API should be
      | Brian |
      | Carol |


  Scenario: normal user tries to get users of a group
    Given user "Brian" has been created with default attributes and without skeleton files
    And group "tea-lover" has been created
    When user "Brian" gets all the members of group "tea-lover" using the Graph API
    Then the HTTP status code should be "401"
