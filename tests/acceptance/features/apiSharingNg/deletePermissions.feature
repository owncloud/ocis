Feature: Delete access to a drive item
  https://owncloud.dev/libre-graph-api/#/drives.permissions/DeletePermission

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario Outline: user removes access from a resource
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType | <resource-type> |
      | resource     | <path>          |
      | space        | Personal        |
      | sharee       | Brian           |
      | shareType    | user            |
      | role         | <role>          |
    When user "Alice" removes the share permission of user "Brian" from <resource-type> "<path>" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <path> |
    Examples:
      | role        | resource-type | path          |
      | Viewer      | file          | textfile.txt  |
      | File Editor | file          | textfile.txt  |
      | Co Owner    | file          | textfile.txt  |
      | Manager     | file          | textfile.txt  |
      | Viewer      | folder        | FolderToShare |
      | Editor      | folder        | FolderToShare |
      | Co Owner    | folder        | FolderToShare |
      | Uploader    | folder        | FolderToShare |
      | Manager     | folder        | FolderToShare |


  Scenario Outline: user removes access from a file inside of a project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType | <resource-type> |
      | resource     | <path>          |
      | space        | NewSpace        |
      | sharee       | Brian           |
      | shareType    | user            |
      | role         | <role>          |
    When user "Alice" removes the share permission of user "Brian" from <resource-type> "<path>" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <path> |
    Examples:
      | role        | resource-type | path          |
      | Viewer      | file          | textfile.txt  |
      | File Editor | file          | textfile.txt  |
      | Co Owner    | file          | textfile.txt  |
      | Manager     | file          | textfile.txt  |
      | Viewer      | folder        | FolderToShare |
      | Editor      | folder        | FolderToShare |
      | Co Owner    | folder        | FolderToShare |
      | Uploader    | folder        | FolderToShare |
      | Manager     | folder        | FolderToShare |
