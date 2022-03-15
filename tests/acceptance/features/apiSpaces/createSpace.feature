@api @skipOnOcV10
Feature: Create spaces
  As a user with space manager role
  I want to be able to create project spaces

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And the administrator has given "Alice" the role "Spacemanager" using the settings api


  Scenario: An user without spacemanager role cannot create a Space via Graph API
    When user "Brian" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "401"


  Scenario: An user with spacemanager role can create a Space via the Graph API with default quota
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "201"
    And the json responded should contain a space "Project Mars" with these key and value pairs:
      | key              | value        |
      | driveType        | project      |
      | name             | Project Mars |
      | quota@@@total    | 1000000000   |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario: An user with spacemanager role user can create a Space via the Graph API with certain quota
    When user "Alice" creates a space "Project Venus" of type "project" with quota "2000" using the GraphApi
    Then the HTTP status code should be "201"
    And the json responded should contain a space "Project Venus" with these key and value pairs:
      | key              | value         |
      | driveType        | project       |
      | name             | Project Venus |
      | quota@@@total    | 2000          |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
