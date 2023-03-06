@api @skipOnOcV10
Feature: Space management
  As a user with space admin permission
  I want to be able to manage all existing project space
  - I can get all project space where I am not member using "graph/v1.0/drives" endpoint
  - I can edit space: change quota, name, description
  - I can enable, disable, delete space

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And using spaces DAV path
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Alice" has created a space "Project" of type "project" with quota "10"


  Scenario: The space admin user can see another project space even if he is not member of the space
    When user "Brian" lists all spaces via the GraphApi with query "$filter=driveType eq 'project'"
    Then the HTTP status code should be "200"
    And the json responded should contain a space "Project" with these key and value pairs:
      | key       | value      |
      | driveType | project    |
      | id        | %space_id% |
      | name      | Project    |
    And the json responded should not contain a space with name "Alice Hansen"


  Scenario: The space admin user can see another personal spaces
    When user "Brian" lists all spaces via the GraphApi with query "$filter=driveType eq 'personal'"
    Then the HTTP status code should be "200"
    And the json responded should contain a space "Alice Hansen" with these key and value pairs:
      | key       | value        |
      | driveType | personal     |
      | id        | %space_id%   |
      | name      | Alice Hansen |
    And the json responded should not contain a space with name "Project"


  Scenario: The user without space admin permissions cannot see another spaces
    When user "Carol" tries to list all spaces via the GraphApi
    Then the HTTP status code should be "200"
    And the json responded should not contain a space with name "Project"
    And the json responded should not contain a space with name "Alice Hansen"


  Scenario: The space admin user changes the quota of the project space
    When user "Brian" changes the quota of the "Project" space to "20" owned by user "Alice"
    Then the HTTP status code should be "200"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key           | value |
      | quota@@@total | 20    |


  Scenario: The user without space admin permissions tries to change the quota of the project space
    When user "Carol" tries to change the quota of the "Project" space to "20" owned by user "Alice"
    Then the HTTP status code should be "401"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key           | value |
      | quota@@@total | 10    |


  Scenario: The space admin user tries to change the quota of the personal space
    When user "Brian" tries to change the quota of the "Alice Hansen" space to "20" owned by user "Alice"
    Then the HTTP status code should be "401"
    And the user "Alice" should have a space called "Alice Hansen" with these key and value pairs:
      | key           | value |
      | quota@@@total | 10    |


  Scenario: The user without space admin permissions tries to change the quota of the personal space
    When user "Carol" tries to change the quota of the "Alice Hansen" space to "20" owned by user "Alice"
    Then the HTTP status code should be "401"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key           | value |
      | quota@@@total | 10    |

  @skipOnStable2.0
  Scenario: The space admin user changes the name of the project space
    When user "Brian" changes the name of the "Project" space to "New Name" owned by user "Alice"
    Then the HTTP status code should be "200"
    And the user "Alice" should have a space called "New Name" with these key and value pairs:
      | key  | value    |
      | name | New Name |


  Scenario: The user without space admin permissions tries to change the name of the project space
    When user "Carol" tries to change the name of the "Project" space to "New Name" owned by user "Alice"
    Then the HTTP status code should be "403"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key  | value   |
      | name | Project |

  @skipOnStable2.0
  Scenario: The space admin user changes the description of the project space
    When user "Brian" changes the description of the "Project" space to "New description" owned by user "Alice"
    Then the HTTP status code should be "200"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key         | value           |
      | description | New description |


  Scenario: The user without space admin permissions tries to change the description of the project space
    Given user "Alice" has changed the description of the "Project" space to "old description"
    When user "Carol" tries to change the description of the "Project" space to "New description" owned by user "Alice"
    Then the HTTP status code should be "403"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key         | value           |
      | description | old description |

  @skipOnStable2.0
  Scenario: The space admin user disables the project space
    When user "Brian" disables a space "Project" owned by user "Alice"
    Then the HTTP status code should be "204"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key                    | value   |
      | name                   | Project |
      | root@@@deleted@@@state | trashed |


  Scenario: The user without space admin permissions tries to disable the project space
    When user "Carol" tries to disable a space "Project" owned by user "Alice"
    Then the HTTP status code should be "403"


  Scenario Outline: The space admin user tries to disable the personal space
    When user "<user>" disables a space "Alice Hansen" owned by user "Alice"
    Then the HTTP status code should be "403"
    Examples:
      | user  |
      | Brian |
      | Carol |

  @skipOnStable2.0
  Scenario: The space admin user deletes the project space
    Given user "Alice" has disabled a space "Project"
    When user "Brian" deletes a space "Project" owned by user "Alice"
    Then the HTTP status code should be "204"
    And the user "Alice" should not have a space called "Project"


  Scenario: The user without space admin permissions tries to delete the project space
    Given user "Alice" has disabled a space "Project"
    When user "Carol" tries to delete a space "Project" owned by user "Alice"
    Then the HTTP status code should be "403"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key                    | value   |
      | name                   | Project |
      | root@@@deleted@@@state | trashed |

  @skipOnStable2.0
  Scenario: The space admin user enables the project space
    Given user "Alice" has disabled a space "Project"
    When user "Brian" restores a disabled space "Project" owned by user "Alice"
    Then the HTTP status code should be "200"


  Scenario: The user without space admin permissions tries to enable the project space
    Given user "Alice" has disabled a space "Project"
    When user "Carol" tries to restore a disabled space "Project" owned by user "Alice"
    Then the HTTP status code should be "404"
    And the user "Alice" should have a space called "Project" with these key and value pairs:
      | key                    | value   |
      | name                   | Project |
      | root@@@deleted@@@state | trashed |
