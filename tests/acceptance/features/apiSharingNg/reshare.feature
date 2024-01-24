Feature: Reshare a share invitation
  As a user
  I want to be able to reshare the share invitations to other users
  So that they can have access to it

  https://owncloud.dev/libre-graph-api/#/drives.permissions/Invite

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |


  Scenario Outline: reshare a file to a user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file               |
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType    | file                       |
      | resource        | textfile1.txt              |
      | space           | Shares                     |
      | sharee          | Carol                      |
      | shareType       | user                       |
      | permissionsRole | <reshare-permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Carol" the space Shares should contain these entries:
      | textfile1.txt |
    Examples:
      | permissions-role | reshare-permissions-role |
      | Viewer           | Viewer                   |
      | File Editor      | Viewer                   |
      | File Editor      | File Editor              |


  Scenario Outline: reshare a folder to a user with different roles
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following share invitation:
      | resourceType    | folder             |
      | resource        | FolderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType    | folder                     |
      | resource        | FolderToShare              |
      | space           | Shares                     |
      | sharee          | Carol                      |
      | shareType       | user                       |
      | permissionsRole | <reshare-permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Carol" the space Shares should contain these entries:
      | FolderToShare |
    Examples:
      | permissions-role | reshare-permissions-role |
      | Viewer           | Viewer                   |
      | Editor           | Viewer                   |
      | Editor           | Editor                   |
      | Editor           | Uploader                 |


  Scenario Outline: reshare a file inside project space to a user with different roles
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "to share" to "textfile1.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file               |
      | resource        | textfile1.txt      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType    | file                       |
      | resource        | textfile1.txt              |
      | space           | Shares                     |
      | sharee          | Carol                      |
      | shareType       | user                       |
      | permissionsRole | <reshare-permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Carol" the space Shares should contain these entries:
      | textfile1.txt |
    Examples:
      | permissions-role | reshare-permissions-role |
      | Viewer           | Viewer                   |
      | File Editor      | Viewer                   |
      | File Editor      | File Editor              |


  Scenario Outline: reshare a folder inside project space to a user with different roles
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has sent the following share invitation:
      | resourceType    | folder             |
      | resource        | FolderToShare      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType    | folder                     |
      | resource        | FolderToShare              |
      | space           | Shares                     |
      | sharee          | Carol                      |
      | shareType       | user                       |
      | permissionsRole | <reshare-permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Carol" the space Shares should contain these entries:
      | FolderToShare |
    Examples:
      | permissions-role | reshare-permissions-role |
      | Viewer           | Viewer                   |
      | Editor           | Viewer                   |
      | Editor           | Editor                   |
      | Editor           | Uploader                 |


  Scenario Outline: reshare a file to a group with different roles
    Given user "Bob" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Carol    | grp1      |
      | Bob      | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file               |
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType    | file                       |
      | resource        | textfile1.txt              |
      | space           | Shares                     |
      | sharee          | grp1                       |
      | shareType       | group                      |
      | permissionsRole | <reshare-permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Carol" the space Shares should contain these entries:
      | textfile1.txt |
    And for user "Bob" the space Shares should contain these entries:
      | textfile1.txt |
    Examples:
      | permissions-role | reshare-permissions-role |
      | Viewer           | Viewer                   |
      | File Editor      | Viewer                   |
      | File Editor      | File Editor              |


  Scenario Outline: reshare a folder to a group with different roles
    Given user "Bob" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Carol    | grp1      |
      | Bob      | grp1      |
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following share invitation:
      | resourceType    | folder             |
      | resource        | FolderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType    | folder                     |
      | resource        | FolderToShare              |
      | space           | Shares                     |
      | sharee          | grp1                       |
      | shareType       | group                      |
      | permissionsRole | <reshare-permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Carol" the space Shares should contain these entries:
      | FolderToShare |
    And for user "Bob" the space Shares should contain these entries:
      | FolderToShare |
    Examples:
      | permissions-role | reshare-permissions-role |
      | Viewer           | Viewer                   |
      | Editor           | Viewer                   |
      | Editor           | Editor                   |
      | Editor           | Uploader                 |


  Scenario Outline: reshare a file inside project space to a group with different roles
    Given using spaces DAV path
    And user "Bob" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Carol    | grp1      |
      | Bob      | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "to share" to "textfile1.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file               |
      | resource        | textfile1.txt      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType    | file                       |
      | resource        | textfile1.txt              |
      | space           | Shares                     |
      | sharee          | grp1                       |
      | shareType       | group                      |
      | permissionsRole | <reshare-permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Carol" the space Shares should contain these entries:
      | textfile1.txt |
    And for user "Bob" the space Shares should contain these entries:
      | textfile1.txt |
    Examples:
      | permissions-role | reshare-permissions-role |
      | Viewer           | Viewer                   |
      | File Editor      | Viewer                   |
      | File Editor      | File Editor              |


  Scenario Outline: reshare a folder inside project space to a group with different roles
    Given using spaces DAV path
    And user "Bob" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Carol    | grp1      |
      | Bob      | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has sent the following share invitation:
      | resourceType    | folder             |
      | resource        | FolderToShare      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType    | folder                     |
      | resource        | FolderToShare              |
      | space           | Shares                     |
      | sharee          | grp1                       |
      | shareType       | group                      |
      | permissionsRole | <reshare-permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Carol" the space Shares should contain these entries:
      | FolderToShare |
    And for user "Bob" the space Shares should contain these entries:
      | FolderToShare |
    Examples:
      | permissions-role | reshare-permissions-role |
      | Viewer           | Viewer                   |
      | Editor           | Viewer                   |
      | Editor           | Editor                   |
      | Editor           | Uploader                 |
