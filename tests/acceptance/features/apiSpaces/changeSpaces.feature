@api @skipOnOcV10
Feature: Change data of space
  As a user with admin rights
  I want to be able to change the data of a created space (increase the quota, change name, etc.)

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api

  Scenario: An admin user can change the name of a Space via the Graph API
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "20"
    When user "Alice" changes the name of the "Project Jupiter" space to "Project Death Star"
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Death Star" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | name             | Project Death Star               |
      | quota@@@total    | 20                               |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |

  Scenario: An admin user can increase the quota of a Space via the Graph API
    Given user "Alice" has created a space "Project Earth" of type "project" with quota "20"
    When user "Alice" changes the quota of the "Project Earth" space to "100"
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Earth" with these key and value pairs:
      | key              | value         |
      | name             | Project Earth |
      | quota@@@total    | 100           |
