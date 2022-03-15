@api @skipOnOcV10
Feature: List spaces
  As a user
  I want to be able to work with personal and project spaces

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: An ordinary user can request information about their Space via the Graph API
    When user "Alice" lists all available spaces via the GraphApi
    Then the HTTP status code should be "200"
    And the json responded should contain a space "Alice Hansen" with these key and value pairs:
      | key              | value        |
      | driveType        | personal     |
      | id               | %space_id%   |
      | name             | Alice Hansen |
      | quota@@@state    | normal       |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
    And the json responded should contain a space "Shares Jail" with these key and value pairs:
      | key              | value        |
      | driveType        | virtual     |
      | id               | %space_id%   |
      | name             | Shares Jail |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario: An ordinary user can request information about their Space via the Graph API using a filter
    When user "Alice" lists all available spaces via the GraphApi with query "$filter=driveType eq 'personal'"
    Then the HTTP status code should be "200"
    And the json responded should contain a space "Alice Hansen" with these key and value pairs:
      | key              | value        |
      | driveType        | personal     |
      | id               | %space_id%   |
      | name             | Alice Hansen |
      | quota@@@state    | normal       |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
    And the json responded should not contain a space with name "Shares Jail"
    And the json responded should only contain spaces of type "personal"


  Scenario: An ordinary user will not see any space when using a filter for project
    When user "Alice" lists all available spaces via the GraphApi with query "$filter=driveType eq 'project'"
    Then the HTTP status code should be "200"
    And the json responded should not contain a space with name "Alice Hansen"
    And the json responded should not contain spaces of type "personal"


  Scenario: An ordinary user can access their Space via the webDav API
    When user "Alice" lists all available spaces via the GraphApi
    And user "Alice" lists the content of the space with the name "Alice Hansen" using the WebDav Api
    Then the HTTP status code should be "207"


  Scenario: A user can list his personal space via multiple endpoints
    When user "Alice" lists all available spaces via the GraphApi with query "$filter=driveType eq 'personal'"
    Then the json responded should contain a space "Alice Hansen" owned by "Alice" with these key and value pairs:
      | key              | value         |
      | driveType        | personal      |
      | name             | Alice Hansen  |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
      | owner@@@user@@@id | %user_id%    |
    When user "Alice" looks up the single space "Alice Hansen" via the GraphApi by using its id
    Then the json responded should contain a space "Alice Hansen" with these key and value pairs:
      | key              | value         |
      | driveType        | personal      |
      | name             | Alice Hansen  |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario: A user can list his created spaces via multiple endpoints
    Given the administrator has given "Alice" the role "Spacemanager" using the settings api
    And user "Alice" has created a space "Project Venus" of type "project" with quota "2000"
    When user "Alice" looks up the single space "Project Venus" via the GraphApi by using its id
    Then the json responded should contain a space "Project Venus" with these key and value pairs:
      | key              | value         |
      | driveType        | project       |
      | name             | Project Venus |
      | quota@@@total    | 2000          |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
