@api @skipOnOcV10
Feature: Change data of space
  As a user
  I want to be able to change the data of the created space(increase the quota, change name, etc.)

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator gives "Alice" the role "Admin" using the settings api

  Scenario: Alice changes a name of the space via the Graph api, she expects a 204 code and checks that the space name has changed
    When user "Alice" creates a space "Project Jupiter" of type "project" with quota "20" using the GraphApi
    And user "Alice" lists all available spaces via the GraphApi
    And user "Alice" changes the name of the "Project Jupiter" space to "Project Death Star"
    Then the HTTP status code should be "204"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Death Star" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | name             | Project Death Star               |
      | quota@@@total    | 20                               |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |

  Scenario: Alice increases quota of the space via the Graph api, she expects a 204 code and checks that the quota has changed
    When user "Alice" creates a space "Project Earth" of type "project" with quota "20" using the GraphApi
    And user "Alice" lists all available spaces via the GraphApi
    And user "Alice" changes the quota of the "Project Earth" space to "100"
    Then the HTTP status code should be "204"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Earth" with these key and value pairs:
      | key              | value         |
      | name             | Project Earth |
      | quota@@@total    | 100           |
      