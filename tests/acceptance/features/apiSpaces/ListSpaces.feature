@api @skipOnOcV10
Feature: List and create spaces
  As a user
  I want to be able to work with personal and project spaces to collaborate with individuals and teams

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  Scenario: Alice request her space via the Graph api, she expects a 200 code and the correct data in the response
    When user "Alice" lists all available spaces via the GraphApi
    Then the HTTP status code should be "200"
    And the json responded should contain these key and value pairs
      | key       | value        |
      | driveType | personal     |
      | name      | Alice Hansen |

  Scenario: Alice requests her space via webDav api, she expects a 207 code
    When user "Alice" lists all available spaces via the GraphApi
    And user "Alice" lists the content of the space with the name "Alice Hansen" using the WebDav Api
    Then the HTTP status code should be "207"

  Scenario: Alice tryes to create Space via Graph api without right, she expects a response of 401
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "401"

  Scenario: Alice creates Space via Graph api with defaul quota, she expects a 201 code the correct data and that space exists
    Given the administrator gives "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "201"
    And the json responded should contain these key and value pairs
      | key       | value        |
      | driveType | project      |
      | name      | Project Mars |
      | total     | 1000000000   |
    When user "Alice" lists all available spaces via the GraphApi
    And user "Alice" lists the content of the space with the name "Project Mars" using the WebDav Api
    Then the propfind result of the space should contain these entries:
      | .space/ |

  Scenario: Alice creates Space via Graph api with certain quota, she expects a 201 code and the correct data in the response
    Given the administrator gives "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Venus" of type "project" with quota "2000" using the GraphApi
    Then the HTTP status code should be "201"
    And the json responded should contain these key and value pairs
      | key        | value         |
      | driveType  | project       |
      | name       | Project Venus |
      | total      | 2000          |

  Scenario: Alice creates folder via Graph api in space, she expects a 201 code and she checks that folder exists
    Given the administrator gives "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Venus" of type "project" with quota "2000" using the GraphApi
    And user "Alice" lists all available spaces via the GraphApi
    And user "Alice" creates a folder "mainFolder" in space "Project Venus" using the WebDav Api
    Then the HTTP status code should be "201"
    When user "Alice" lists the content of the space with the name "Project Venus" using the WebDav Api
    Then the propfind result of the space should contain these entries:
      | mainFolder/        |
