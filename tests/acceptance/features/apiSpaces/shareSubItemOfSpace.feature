@api @skipOnOcV10
Feature: Share a file or folder that is inside a space
  As an user with manager space role
  I want to be able to share the data inside the space


  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "share sub-item" with the default quota using the GraphApi
    And user "Alice" has created a folder "folder" in space "share sub-item"
    And user "Alice" has uploaded a file inside space "share sub-item" with content "some content" to "file.txt"


  Scenario Outline: A manager of the space can share an entity inside project space to another user with role:
    When user "Alice" shares the following entity "<entity>" inside of space "share sub-item" with user "Brian" with role "<role>"
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    When user "Brian" accepts share "/<entity>" offered by user "Alice" using the sharing API
    And as "Brian" <type> "Shares/<entity>" should exist
    Examples:
      | entity   | type   | role   |
      | folder   | folder | viewer |
      | folder   | folder | editor |
      | file.txt | file   | viewer |
      | file.txt | file   | editor |


  Scenario Outline: An user participant of the project space with manager role can share an entity to another user
    Given user "Alice" has shared a space "share sub-item" to user "Brian" with role "manager"
    When user "Brian" shares the following entity "<entity>" inside of space "share sub-item" with user "Bob" with role "<role>"
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    When user "Bob" accepts share "/<entity>" offered by user "Brian" using the sharing API
    And as "Bob" <type> "Shares/<entity>" should exist
    Examples:
      | entity   | type   | role   |
      | folder   | folder | viewer |
      | folder   | folder | editor |
      | file.txt | file   | viewer |
      | file.txt | file   | editor |


  Scenario Outline: An user participant of the project space without space manager role cannot share an entity to another user
    Given user "Alice" has shared a space "share sub-item" to user "Brian" with role "<spaceRole>"
    When user "Brian" shares the following entity "<entity>" inside of space "share sub-item" with user "Bob" with role "editor"
    Then the HTTP status code should be "<statusCode>"
    And the OCS status code should be "<statusCode>"
    And the OCS status message should be "<statusMessage>"
    Examples:
      | entity   | spaceRole | statusCode | statusMessage       |
      | folder   | editor    | 404        | No share permission |
      | file.txt | editor    | 404        | No share permission |
      | file.txt | viewer    | 404        | No share permission |
      | folder   | viewer    | 404        | No share permission |
