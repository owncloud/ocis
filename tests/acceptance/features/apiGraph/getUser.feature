@api
Feature: get user information
  As user
  I want to be able to retrieve my own information
  So that I can see my information

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: user gets his/her own information with no group involvement
    When the user "Alice" retrives her information using the Graph API
    Then the HTTP status code should be "200"
    And the api response for user "Alice" should contains the following information:
      | displayName              | Alice Hansen      |
      | id                       | %user_id%         |
      | mail                     | alice@example.org |
      | onPremisesSamAccountName | Alice             |
      | memberOf                 |                   |


  Scenario: user gets his/her own information with group involvement
    Given group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Alice" has been added to group "coffee-lover"
    When the user "Alice" retrives her information using the Graph API
    And the api response for user "Alice" should contains the following information:
      | displayName              | Alice Hansen            |
      | id                       | %user_id%               |
      | onPremisesSamAccountName | Alice                   |
      | mail                     | alice@example.org       |
      | memberOf                 | tea-lover, coffee-lover |



