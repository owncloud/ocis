Feature: Remove access to a drive item
  https://owncloud.dev/libre-graph-api/#/drives.permissions/DeletePermission

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario Outline: user removes access to resource in the user share
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resource        | <path>            |
      | space           | Personal          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | <permissionsRole> |
    When user "Alice" removes the share permission of user "Brian" from "<path>" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <path> |
    Examples:
      | permissionsRole | path          |
      | Viewer          | textfile.txt  |
      | File Editor     | textfile.txt  |
      | Viewer          | FolderToShare |
      | Editor          | FolderToShare |
      | Uploader        | FolderToShare |


  Scenario Outline: user removes access to resource inside of a project space in the user share
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resource        | <path>            |
      | space           | NewSpace          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | <permissionsRole> |
    When user "Alice" removes the share permission of user "Brian" from "<path>" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <path> |
    Examples:
      | permissionsRole | path          |
      | Viewer          | textfile.txt  |
      | File Editor     | textfile.txt  |
      | Viewer          | FolderToShare |
      | Editor          | FolderToShare |
      | Uploader        | FolderToShare |


  Scenario Outline: user removes access to a resource in a group share
    Given group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has been added to group "group1"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resource        | <path>            |
      | space           | Personal          |
      | sharee          | group1            |
      | shareType       | group             |
      | permissionsRole | <permissionsRole> |
    When user "Alice" removes the share permission of group "group1" from "<path>" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <path> |
    Examples:
      | permissionsRole | path          |
      | Viewer          | textfile.txt  |
      | File Editor     | textfile.txt  |
      | Viewer          | FolderToShare |
      | Editor          | FolderToShare |
      | Uploader        | FolderToShare |


  Scenario Outline: user removes access to a resource inside of a project space in group share
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "textfile.txt"
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has been added to group "group1"
    And user "Alice" has sent the following share invitation:
      | resource        | <path>            |
      | space           | NewSpace          |
      | sharee          | group1            |
      | shareType       | group             |
      | permissionsRole | <permissionsRole> |
    When user "Alice" removes the share permission of group "group1" from "<path>" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <path> |
    Examples:
      | permissionsRole | path          |
      | Viewer          | textfile.txt  |
      | File Editor     | textfile.txt  |
      | Viewer          | FolderToShare |
      | Editor          | FolderToShare |
      | Uploader        | FolderToShare |


  Scenario Outline: user removes access to a folder in link share
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has created the following link share:
      | resource        | FolderToShare |
      | space           | Personal      |
      | permissionsRole | <role>        |
      | password        | %public%      |
    When user "Alice" removes the share permission of link from "FolderToShare" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    Examples:
      | role           |
      | view           |
      | edit           |
      | upload         |
      | createOnly     |
      | blocksDownload |


  Scenario Outline: user removes access to a file in link share
    Given user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile.txt"
    And user "Alice" has created the following link share:
      | resource        | textfile.txt      |
      | space           | Personal          |
      | permissionsRole | <permissionsRole> |
      | password        | %public%          |
    When user "Alice" removes the share permission of link from "textfile.txt" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    Examples:
      | permissionsRole |
      | view            |
      | edit            |
      | blocksDownload  |


  Scenario Outline: user removes access to a folder in project space in link share
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has created the following link share:
      | resource        | FolderToShare |
      | space           | NewSpace      |
      | permissionsRole | <role>        |
      | password        | %public%      |
    When user "Alice" removes the share permission of link from "FolderToShare" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    Examples:
      | role           |
      | view           |
      | edit           |
      | upload         |
      | createOnly     |
      | blocksDownload |


  Scenario Outline: user removes access to a file in project space in link share
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "textfile.txt"
    And user "Alice" has created the following link share:
      | resource        | textfile.txt      |
      | space           | NewSpace          |
      | permissionsRole | <permissionsRole> |
      | password        | %public%          |
    When user "Alice" removes the share permission of link from "textfile.txt" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    Examples:
      | permissionsRole |
      | view            |
      | edit            |
      | blocksDownload  |
