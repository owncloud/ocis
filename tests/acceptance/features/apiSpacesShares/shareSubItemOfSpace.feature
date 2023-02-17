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
    And using spaces DAV path
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "share sub-item" with the default quota using the GraphApi
    And user "Alice" has created a folder "folder" in space "share sub-item"
    And user "Alice" has uploaded a file inside space "share sub-item" with content "some content" to "file.txt"
    And using new DAV path


  Scenario Outline: A manager of the space can share an entity inside project space to another user with role
    And user "Alice" creates a share inside of space "share sub-item" with settings:
      | path       | <entity>     |
      | shareWith  | Brian        |
      | role       | <role>       |
      | expireDate | <expireDate> |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    When user "Brian" accepts share "/<entity>" offered by user "Alice" using the sharing API
    Then as "Brian" <type> "Shares/<entity>" should exist
    And the information about the last share for user "Brian" should include
      | expiration | <expiration> |
    Examples:
      | entity   | type   | role   | expireDate               | expiration |
      | folder   | folder | viewer |                          |            |
      | folder   | folder | editor | 2042-03-25T23:59:59+0100 | 2042-03-25 |
      | file.txt | file   | viewer |                          |            |
      | file.txt | file   | editor | 2042-03-25T23:59:59+0100 | 2042-03-25 |


  Scenario Outline: A user participant of the project space with manager role can share an entity to another user
    Given user "Alice" has shared a space "share sub-item" with settings:
      | shareWith | Brian   |
      | role      | manager |
    When user "Brian" creates a share inside of space "share sub-item" with settings:
      | path       | <entity>     |
      | shareWith  | Bob          |
      | role       | <role>       |
      | expireDate | <expireDate> |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    When user "Bob" accepts share "/<entity>" offered by user "Brian" using the sharing API
    Then as "Bob" <type> "Shares/<entity>" should exist
    And the information about the last share for user "Brian" should include
      | expiration | <expiration> |
    Examples:
      | entity   | type   | role   | expireDate               | expiration |
      | folder   | folder | viewer | 2042-03-25T23:59:59+0100 | 2042-03-25 |
      | folder   | folder | editor |                          |            |
      | file.txt | file   | viewer | 2042-03-25T23:59:59+0100 | 2042-03-25 |
      | file.txt | file   | editor |                          |            |


  Scenario Outline: A user participant of the project space without space manager role cannot share an entity to another user
    Given user "Alice" has shared a space "share sub-item" with settings:
      | shareWith | Brian       |
      | role      | <spaceRole> |
    When user "Brian" creates a share inside of space "share sub-item" with settings:
      | path      | <entity> |
      | shareWith | Bob      |
      | role      | editor   |
    Then the HTTP status code should be "<statusCode>"
    And the OCS status code should be "<statusCode>"
    And the OCS status message should be "<statusMessage>"
    Examples:
      | entity   | spaceRole | statusCode | statusMessage       |
      | folder   | editor    | 404        | No share permission |
      | file.txt | editor    | 404        | No share permission |
      | file.txt | viewer    | 404        | No share permission |
      | folder   | viewer    | 404        | No share permission |


  Scenario Outline: A user participant of the project space can see the created resources share
    Given user "Alice" has shared a space "share sub-item" with settings:
      | shareWith | Brian       |
      | role      | <spaceRole> |
    When user "Alice" creates a share inside of space "share sub-item" with settings:
      | path      | file.txt |
      | shareWith | Bob      |
      | role      | editor   |
    Then for user "Alice" the space "share sub-item" should contain the last created share of the file "file.txt"
    And for user "Brian" the space "share sub-item" should contain the last created share of the file "file.txt"
    Examples:
      | spaceRole |
      | editor    |
      | viewer    |
      | manager   |


  Scenario: A user shares the folder to the group
    Given group "sales" has been created
    And the administrator has added a user "Brian" to the group "sales" using GraphApi
    When user "Alice" creates a share inside of space "share sub-item" with settings:
      | path       | folder                   |
      | shareWith  | sales                    |
      | shareType  | 1                        |
      | role       | viewer                   |
      | expireDate | 2042-01-01T23:59:59+0100 |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    When user "Brian" accepts share "/folder" offered by user "Alice" using the sharing API
    Then as "Brian" folder "Shares/folder" should exist
    And the information about the last share for user "Brian" should include
      | expiration | 2042-01-01 |


  Scenario: A user changes the expiration date
    Given user "Alice" has created a share inside of space "share sub-item" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
      | expireDate | 2042-01-01T23:59:59+0100 |
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    When user "Alice" changes the last share with settings:
      | expireDate | 2044-01-01T23:59:59.999+01:00 |
    Then the HTTP status code should be "200"
    And the information about the last share for user "Brian" should include
      | expiration | 2044-01-01 |
