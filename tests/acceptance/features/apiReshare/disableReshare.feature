@env-config
Feature: share by disabling re-share
  As a user
  I want to share resources
  So that other users can have access to them but cannot re-share them

  Background:
    Given the config "FRONTEND_ENABLE_RESHARING" has been set to "false"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
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
      | permissions | <expectedPermissions> |
    Examples:
      | dav-path-version | role   | expectedPermissions |
      | old              | editor | 15                  |
      | old              | viewer | 1                   |
      | new              | editor | 15                  |
      | new              | viewer | 1                   |
      | spaces           | editor | 15                  |
      | spaces           | viewer | 1                   |


  Scenario Outline: try to re-share folder
    Given using <dav-path-version> DAV path
    And user "Carol" has been created with default attributes and without skeleton files
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
