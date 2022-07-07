@api @skipOnOcV10
Feature: Upload files into a space
  As an user
  I want to be able to create folders and files in the space

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "Project Ceres" of type "project" with quota "2000"


  Scenario Outline: An user creates a folder in the Space via the Graph API
    And user "Alice" has shared a space "Project Ceres" to user "Brian" with role "<role>"
    When user "Brian" creates a folder "mainFolder" in space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "<code>"
    And for user "Brian" the space "Project Ceres" <shouldOrNot> contain these entries:
      | mainFolder |
    Examples:
      | role    | code | shouldOrNot |
      | manager | 201  | should      |
      | editor  | 201  | should      |
      | viewer  | 403  | should not  |


  Scenario Outline: An user uploads a file in shared Space via the Graph API
    And user "Alice" has shared a space "Project Ceres" to user "Brian" with role "<role>"
    When user "Brian" uploads a file inside space "Project Ceres" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "<code>"
    And for user "Brian" the space "Project Ceres" <shouldOrNot> contain these entries:
      | test.txt |
    And the user "Brian" should have a space called "Project Ceres" with these key and value pairs:
      | key          | value         |
      | name         | Project Ceres |
      | quota@@@used | <usedQuota>   |
    Examples:
      | role    | code | shouldOrNot | usedQuota |
      | manager | 201  | should      | 4         |
      | editor  | 201  | should      | 4         |
      | viewer  | 403  | should not  | 0         |


  Scenario: An user can create subfolders in a Space via the Graph API
    When user "Alice" creates a subfolder "mainFolder/subFolder1/subFolder2" in space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Ceres" should contain these entries:
      | mainFolder |
    And for user "Alice" folder "mainFolder/subFolder1/" of the space "Project Ceres" should contain these entries:
      | subFolder2 |


  Scenario: An user can create a folder and upload a file to a Space
    When user "Alice" creates a folder "NewFolder" in space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "201"
    And user "Alice" uploads a file inside space "Project Ceres" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Ceres" should contain these entries:
      | NewFolder |
      | test.txt  |


  Scenario: An user cannot create a folder or a file in a Space if they do not have permission
    When user "Bob" creates a folder "forAlice" in space "Project Ceres" owned by the user "Alice" using the WebDav Api
    Then the HTTP status code should be "404"
    When user "Bob" uploads a file inside space "Project Ceres" owned by the user "Alice" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "404"
    And for user "Alice" the space "Project Ceres" should not contain these entries:
      | forAlice |
      | test.txt |
