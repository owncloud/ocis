@api @skipOnOcV10
Feature: List and create spaces
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
      | key              | value                            |
      | driveType        | personal                         |
      | driveAlias       | personal/alice                   |
      | id               | %space_id%                       |
      | name             | Alice Hansen                     |
      | quota@@@state    | normal                           |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario: An ordinary user can request information about their Space via the Graph API using a filter
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "folder"
    And user "Brian" has shared folder "folder" with user "Alice" with permissions "31"
    And user "Alice" has accepted share "/folder" offered by user "Brian"
    Then the user "Alice" should have a space called "Shares Jail" with these key and value pairs:
      | key       | value       |
      | driveType | virtual     |
      | id        | %space_id%  |
      | name      | Shares Jail |
    When user "Alice" lists all available spaces via the GraphApi with query "$filter=driveType eq 'personal'"
    Then the HTTP status code should be "200"
    And the json responded should contain a space "Alice Hansen" with these key and value pairs:
      | key              | value                            |
      | driveType        | personal                         |
      | id               | %space_id%                       |
      | name             | Alice Hansen                     |
      | quota@@@state    | normal                           |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
    And the json responded should not contain a space with name "Shares Jail"
    And the json responded should only contain spaces of type "personal"


  Scenario: An ordinary user will not see any space when using a filter for project
    Given the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "my project" of type "project" with quota "20"
    When user "Alice" lists all available spaces via the GraphApi with query "$filter=driveType eq 'project'"
    Then the HTTP status code should be "200"
    And the json responded should contain a space "my project" with these key and value pairs:
      | key       | value      |
      | driveType | project    |
      | id        | %space_id% |
      | name      | my project |
    And the json responded should not contain a space with name "Alice Hansen"


  Scenario: An ordinary user can access their Space via the webDav API
    When user "Alice" lists all available spaces via the GraphApi
    And user "Alice" lists the content of the space with the name "Alice Hansen" using the WebDav Api
    Then the HTTP status code should be "207"


  Scenario: An ordinary user cannot create a Space via Graph API
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "401"


  Scenario: An admin user can create a Space via the Graph API with default quota
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "201"
    And the json responded should contain a space "Project Mars" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | driveAlias       | project/project-mars             |
      | name             | Project Mars                     |
      | quota@@@total    | 1000000000                       |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario: An admin user can create a Space via the Graph API with certain quota
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Venus" of type "project" with quota "2000" using the GraphApi
    Then the HTTP status code should be "201"
    And the json responded should contain a space "Project Venus" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | name             | Project Venus                    |
      | quota@@@total    | 2000                             |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario: A user can list his personal space via multiple endpoints
    When user "Alice" lists all available spaces via the GraphApi with query "$filter=driveType eq 'personal'"
    Then the json responded should contain a space "Alice Hansen" owned by "Alice" with these key and value pairs:
      | key               | value                            |
      | driveType         | personal                         |
      | name              | Alice Hansen                     |
      | root@@@webDavUrl  | %base_url%/dav/spaces/%space_id% |
      | owner@@@user@@@id | %user_id%                        |
    When user "Alice" looks up the single space "Alice Hansen" via the GraphApi by using its id
    Then the json responded should contain a space "Alice Hansen" with these key and value pairs:
      | key              | value                            |
      | driveType        | personal                         |
      | name             | Alice Hansen                     |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario: A user can list his created spaces via multiple endpoints
    Given the administrator has given "Alice" the role "Admin" using the settings api
    When user "Alice" creates a space "Project Venus" of type "project" with quota "2000" using the GraphApi
    Then the HTTP status code should be "201"
    And the json responded should contain a space "Project Venus" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | driveAlias       | project/project-venus            |
      | name             | Project Venus                    |
      | quota@@@total    | 2000                             |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
    When user "Alice" looks up the single space "Project Venus" via the GraphApi by using its id
    Then the json responded should contain a space "Project Venus" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | driveAlias       | project/project-venus            |
      | name             | Project Venus                    |
      | quota@@@total    | 2000                             |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
