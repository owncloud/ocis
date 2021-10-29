@api @skipOnOcV10
Feature: Upload files into a space
  As a user
  I want to be able to work with project spaces and quota

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Bob" has been created with default attributes and without skeleton files

  Scenario: Alice creates a folder via Graph api and uploads a file
    Given the administrator gives "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Moon" of type "project" with quota "2000" using the GraphApi
    And user "Alice" lists all available spaces via the GraphApi
    And user "Alice" creates a folder "NewFolder" in space "Project Moon" using the WebDav Api
    Then the HTTP status code should be "201"
    And user "Alice" uploads a file inside space "Project Moon" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    When user "Alice" lists the content of the space with the name "Project Moon" using the WebDav Api
    Then the propfind result of the space should contain these entries:
      | NewFolder/      |
      | test.txt         |

  Scenario: Alice creates uploads a file and checks her quota
    Given the administrator gives "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Saturn" of type "project" with quota "2000" using the GraphApi
    And the json responded should contain a space "Project Saturn" with these key and value pairs:
      | key              | value         |
      | driveType        | project       |
      | id               | %space_id%    |
      | name             | Project Saturn|
      | quota@@@total    | 2000          |
    And user "Alice" lists all available spaces via the GraphApi
    And user "Alice" uploads a file inside space "Project Saturn" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    When user "Alice" lists the content of the space with the name "Project Saturn" using the WebDav Api
    Then the propfind result of the space should contain these entries:
      | test.txt         |
    And user "Alice" lists all available spaces via the GraphApi
    And the json responded should contain a space "Project Saturn" with these key and value pairs:
      | key              | value         |
      | driveType        | project       |
      | id               | %space_id%    |
      | name             | Project Saturn|
      | quota@@@state    | normal        |
      | quota@@@total    | 2000          |
      | quota@@@remaining| 1996          |
      | quota@@@used     | 4             |
