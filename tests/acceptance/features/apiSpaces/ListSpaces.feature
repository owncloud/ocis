@api @skipOnOcV10
Feature: List and create spaces
  As a user
  I want to be able to work with personal and project spaces to collaborate with individuals and teams

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Scenario: list own spaces
    Given user "Alice" has been created with default attributes and without skeleton files
    When user "Alice" lists all available spaces via the GraphApi
    Then the HTTP status code should be "200"
    When user "Alice" lists the content of the space with the name "Alice Hansen" using the WebDav Api
    Then the HTTP status code should be "207"
    #Then the propfind result of the space should contain these entries:
      #| .space/        |
    Then the propfind result of the space should not contain these entries:
      | file        |
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "401"
    When the administrator gives "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "201"
    Then the json responded should contain these key and value pairs
      | key        |     value        |
      | driveType  |     project      |
      | name       |     Project Mars |
    When user "Alice" creates a space "Project Venus" of type "project" with quota "2000" using the GraphApi
