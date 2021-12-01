@api @skipOnOcV10
Feature: Upload files into a space
  As a user
  I want to be able to work with project spaces and quota

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Bob" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api

  Scenario: Alice creates a folder via the Graph api in space, she expects a 201 code and she checks that folder exists
    Given user "Alice" has created a space "Project Venus" of type "project" with quota "2000" 
    When user "Alice" creates a folder "mainFolder" in space "Project Venus" using the WebDav Api
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Venus" should contain these entries:
      | mainFolder        |

  Scenario: Bob creates a folder via the Graph api in a space, he expects a 404 code and Alice checks that this folder does not exist
    Given user "Alice" has created a space "Project Merkur" of type "project" with quota "2000"
    And user "Bob" creates a folder "forAlice" in space "Project Merkur" owned by the user "Alice" using the WebDav Api
    Then the HTTP status code should be "404"
    And for user "Alice" the space "Project Merkur" should not contain these entries:
      | forAlice        |

  Scenario: Alice creates a folder via Graph api and uploads a file
    Given user "Alice" has created a space "Project Moon" of type "project" with quota "2000" 
    When user "Alice" creates a folder "NewFolder" in space "Project Moon" using the WebDav Api
    Then the HTTP status code should be "201"
    And user "Alice" uploads a file inside space "Project Moon" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Moon" should contain these entries:
      | NewFolder        |
      | test.txt         |

  Scenario: Bob uploads a file via the Graph api in a space, he expects a 404 code and Alice checks that this file does not exist
    Given user "Alice" has created a space "Project Pluto" of type "project" with quota "2000" 
    When user "Bob" uploads a file inside space "Project Pluto" owned by the user "Alice" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "404"
    And for user "Alice" the space "Project Pluto" should not contain these entries:
      | test.txt        |

  Scenario: Alice creates uploads a file and checks her quota
    When user "Alice" creates a space "Project Saturn" of type "project" with quota "2000" using the GraphApi
    And the json responded should contain a space "Project Saturn" with these key and value pairs:
      | key              | value         |
      | driveType        | project       |
      | id               | %space_id%    |
      | name             | Project Saturn|
      | quota@@@total    | 2000          |
    And user "Alice" uploads a file inside space "Project Saturn" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Saturn" should contain these entries:
      | test.txt         |
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Saturn" with these key and value pairs:
      | key              | value         |
      | driveType        | project       |
      | id               | %space_id%    |
      | name             | Project Saturn|
      | quota@@@state    | normal        |
      | quota@@@total    | 2000          |
      | quota@@@remaining| 1996          |
      | quota@@@used     | 4             |
