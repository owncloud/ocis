Feature: re-share resources
  As a user
  I cannot to re-share resources
  This feature has been removed from ocis

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And user "Alice" has created folder "test"


  Scenario Outline: share folder with different roles
    Given using <dav-path-version> DAV path
    When user "Alice" creates a share inside of space "Personal" with settings:
      | path      | test   |
      | shareWith | Brian  |
      | role      | <role> |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the fields of the last response to user "Alice" sharing with user "Brian" should include
      | permissions | <expected-permissions> |
    Examples:
      | dav-path-version | role   | expected-permissions |
      | old              | editor | 15                   |
      | old              | viewer | 1                    |
      | new              | editor | 15                   |
      | new              | viewer | 1                    |
      | spaces           | editor | 15                   |
      | spaces           | viewer | 1                    |


  Scenario Outline: try to re-share folder
    Given using <dav-path-version> DAV path
    And user "Alice" has created a share inside of space "Personal" with settings:
      | path      | test   |
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" creates a share inside of space "Shares" with settings:
      | path      | test   |
      | shareWith | Carol  |
      | role      | <role> |
    Then the HTTP status code should be "403"
    And the OCS status code should be "403"
    And the OCS status message should be "No share permission"
    Examples:
      | dav-path-version | role   |
      | old              | editor |
      | old              | viewer |
      | new              | editor |
      | new              | viewer |
      | spaces           | editor |
      | spaces           | viewer |


  Scenario Outline: try to re-share file
    Given user "Alice" has uploaded file with content "other data" to "/textfile1.txt"
    Given using <dav-path-version> DAV path
    And user "Alice" has created a share inside of space "Personal" with settings:
      | path      | textfile1.txt |
      | shareWith | Brian         |
      | role      | <role>        |
    When user "Brian" creates a share inside of space "Shares" with settings:
      | path      | textfile1.txt |
      | shareWith | Carol         |
      | role      | <role>        |
    Then the HTTP status code should be "403"
    And the OCS status code should be "403"
    And the OCS status message should be "No share permission"
    Examples:
      | dav-path-version | role   |
      | old              | editor |
      | old              | viewer |
      | new              | editor |
      | new              | viewer |
      | spaces           | editor |
      | spaces           | viewer |


  Scenario Outline: try to create a link to the shared folder
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has shared folder "/test" with user "Brian" with permissions "all"
    When user "Brian" creates a public link share using the sharing API with settings
      | path        | /Shares/test |
      | permissions | 1            |
      | password    | %public%     |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "<http_status_code>"
    Examples:
      | ocs_api_version | ocs_status_code | http_status_code |
      | 1               | 403             | 200              |
      | 2               | 403             | 403              |
      