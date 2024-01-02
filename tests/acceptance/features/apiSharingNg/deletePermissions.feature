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
      | permissions   |
      | read          |
      | change        |
      | create        |
      | all           |
      | share         |
      | delete        |
      | update        |


  Scenario Outline: user removes access from a file
    Given user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And user "Alice" has shared file "fileToShare.txt" with user "Brian" with permissions "<permissions>"
    When user "Alice" removes the share permission of user "Brian" from file "fileToShare.txt" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | fileToShare.txt |
    Examples:
      | permissions   |
      | read          |
      | change        |
      | all           |
      | share         |
      | update        |
