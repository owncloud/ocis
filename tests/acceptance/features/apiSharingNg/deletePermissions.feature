Feature: Delete access to a drive item
  https://owncloud.dev/libre-graph-api/#/drives.permissions/DeletePermission

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario Outline: user removes access from a folder
    Given user "Alice" has created folder "folderToShare"
    And user "Alice" has shared folder "folderToShare" with user "Brian" with permissions "<permissions>"
    When user "Alice" removes the share permission of user "Brian" from folder "folderToShare" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | folderToShare |
    Examples:
      | permissions |
      | read        |
      | change      |
      | create      |
      | all         |
      | share       |
      | delete      |
      | update      |


  Scenario Outline: user removes access from a file
    Given user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And user "Alice" has shared file "fileToShare.txt" with user "Brian" with permissions "<permissions>"
    When user "Alice" removes the share permission of user "Brian" from file "fileToShare.txt" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | fileToShare.txt |
    Examples:
      | permissions |
      | read        |
      | change      |
      | all         |
      | share       |
      | update      |


  Scenario Outline: user removes access from a folder inside of a project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "NewSpace"
    And user "Alice" has created a share inside of space "NewSpace" with settings:
      | path      | folderToShare |
      | shareWith | Brian         |
      | role      | <role>        |
    When user "Alice" removes the share permission of user "Brian" from folder "folderToShare" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | folderToShare |
    Examples:
      | role   |
      | editor |
      | viewer |


  Scenario Outline: user removes access from a file inside of a project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "file.txt"
    And user "Alice" has created a share inside of space "NewSpace" with settings:
      | path      | file.txt |
      | shareWith | Brian    |
      | role      | <role>   |
    When user "Alice" removes the share permission of user "Brian" from file "file.txt" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | file.txt |
    Examples:
      | role   |
      | editor |
      | viewer |
