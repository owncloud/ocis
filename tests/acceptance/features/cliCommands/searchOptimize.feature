@env-config
Feature: optimize search index via CLI command
  As an administrator
  I want to optimize the search index
  So that I can improve search query performance

  Background:
    Given the administrator has cleaned the search service data


  Scenario: optimize the search index
    Given user "Alice" has been created with default attributes
    And using spaces DAV path
    And user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And the administrator reindexes all spaces using the CLI
    And user "Alice" searches for "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the administrator has stopped the server
    When the administrator optimizes the search index using the CLI
    Then the command should be successful
    And the command output should contain "index optimization complete"
    When the administrator starts the server
    And user "Alice" searches for "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these entries:
      | /textfile.txt |


  Scenario: optimize the search index without any indexed content
    Given the administrator has stopped the server
    When the administrator optimizes the search index using the CLI
    Then the command should be successful
    And the command output should contain "index optimization complete"
