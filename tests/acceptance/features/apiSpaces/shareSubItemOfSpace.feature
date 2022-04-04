@api @skipOnOcV10
Feature: Share a file or folder that is inside a space
      As an user with manager space role
      I want to be able to share the data inside the space

      | role        | permissions |
      | viewer      | 1           |
      | uploader    | 4           |
      | contributor | 5           |
      | editor      | 15          |

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Admin" using the settings api
    And user "Alice" has created a space "share sub-item" with the default quota using the GraphApi
    And user "Alice" has created a folder "folder" in space "share sub-item"
    And user "Alice" has uploaded a file inside space "share sub-item" with content "some content" to "folder/file.txt"


  Scenario Outline: An user-owner can share a folder inside project space to another user with role:
    When user "Alice" shares the following entity "folder" inside of space "share sub-item" with user "Brian" with role "<role>"
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    And user "Brian" accepts share "/folder" offered by user "Alice" using the sharing API
    And as "Brian" folder "Shares/folder" should exist
    Examples:
      | role   |
      | viewer |
      | editor |


  Scenario Outline: An user-owner can share a file inside project space to another user with role:
    When user "Alice" shares the following entity "folder/file.txt" inside of space "share sub-item" with user "Brian" with role "<role>"
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    And user "Brian" accepts share "/file.txt" offered by user "Alice" using the sharing API
    Then as "Brian" file "Shares/file.txt" should exist
    Examples:
      | role   |
      | viewer |
      | editor |


  Scenario Outline: An user participant of the project space with manager role can share a folder to another user
    Given user "Alice" has shared a space "share sub-item" to user "Brian" with role "manager"
    When user "Brian" shares the following entity "folder" inside of space "share sub-item" with user "Bob" with role "<role>"
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    And user "Bob" accepts share "/folder" offered by user "Brian" using the sharing API
    And as "Bob" folder "Shares/folder" should exist
    Examples:
      | role   |
      | viewer |
      | editor |


  Scenario Outline: An user participant of the project space with manager role can share a file to another user
    Given user "Alice" has shared a space "share sub-item" to user "Brian" with role "manager"
    When user "Brian" shares the following entity "folder/file.txt" inside of space "share sub-item" with user "Bob" with role "<role>"
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    And user "Bob" accepts share "/file.txt" offered by user "Brian" using the sharing API
    Then as "Bob" file "Shares/file.txt" should exist
    Examples:
      | role   |
      | viewer |
      | editor |


  Scenario Outline: An user participant of the project space without space manager role cannot share an entity to another user
    Given user "Alice" has shared a space "share sub-item" to user "Brian" with role "<spaceRole>"
    When user "Brian" shares the following entity "<entity>" inside of space "share sub-item" with user "Bob" with role "editor"
    Then the HTTP status code should be "<statusCode>"
    And the OCS status code should be "<statusCode>"
    And the OCS status message should be "<statusMessage>"
    Examples:
      | entity          | spaceRole | statusCode | statusMessage       |
      | folder          | editor    | 404        | No share permission |
      | folder/file.txt | editor    | 404        | No share permission |
      | folder/file.txt | viewer    | 404        | No share permission |
      | folder          | viewer    | 404        | No share permission |
