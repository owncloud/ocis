@ocm
Feature: an user shares resources usin ScienceMesh application
  As a user
  I want to share resources between different ocis instances

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And using server "REMOTE"
    And user "Brian" has been created with default attributes and without skeleton files


  Scenario: user generates invitation
    Given using server "LOCAL"
    When "Alice" generates invitation
    Then the HTTP status code should be "200"
    When using server "REMOTE"
    And "Brian" accepts invitation
    Then the HTTP status code should be "200"
