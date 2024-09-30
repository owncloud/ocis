Feature: propfind a shares
  As a user
  I want to check the PROPFIND response
  So that I can make sure that the response contains all the relevant values

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |

  @issue-4421
  Scenario Outline: sharee PROPFIND same name shares shared by multiple users
    Given using spaces DAV path
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Alice" has created folder "folderToShare"
    And user "Carol" has uploaded file with content "to share" to "textfile.txt"
    And user "Carol" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Brian" has a share "<resource>" synced
    And user "Carol" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Brian" has a share "<resource-2>" synced
    When user "Brian" sends PROPFIND request to space "Shares" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Shares" with these key and value pairs:
      | key       | value         |
      | oc:fileid | UUIDof:Shares |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key            | value      |
      | oc:name        | <resource> |
      | oc:permissions | S          |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key            | value        |
      | oc:name        | <resource-2> |
      | oc:permissions | S            |
    Examples:
      | resource      | resource-2        |
      | textfile.txt  | textfile (1).txt  |
      | folderToShare | folderToShare (1) |

  @issue-4421 @issue-9933 @skip
  Scenario Outline: sharee PROPFIND same name shares shared by multiple users using new dav path
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Alice" has created folder "folderToShare"
    And user "Carol" has uploaded file with content "to share" to "textfile.txt"
    And user "Carol" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Brian" has a share "<resource>" synced
    And user "Carol" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Brian" has a share "<resource-2>" synced
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "Shares" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Shares" with these key and value pairs:
      | key       | value         |
      | oc:fileid | UUIDof:Shares |
      | oc:name   | Shares        |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key            | value             |
      | oc:fileid      | UUIDof:<resource> |
      | oc:name        | <resource>        |
      | oc:permissions | S                 |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key            | value               |
      | oc:fileid      | UUIDof:<resource-2> |
      | oc:name        | <resource-2>        |
      | oc:permissions | S                   |
    Examples:
      | dav-path-version | resource      | resource-2        |
      | old              | textfile.txt  | textfile (1).txt  |
      | old              | folderToShare | folderToShare (1) |
      | new              | textfile.txt  | textfile (1).txt  |
      | new              | folderToShare | folderToShare (1) |

  @issue-4421 @issue-9933 @skip
  Scenario: sharee PROPFIND shares with bracket in the name
    Given using spaces DAV path
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "to share" to "folderToShare/textfile.txt"
    And user "Carol" has created folder "folderToShare"
    And user "Carol" has uploaded file with content "to share" to "folderToShare/textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Brian" has a share "folderToShare" synced
    And user "Carol" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Brian" has a share "folderToShare (1)" synced
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "folderToShare (1)" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "folderToShare (1)" with these key and value pairs:
      | key            | value                    |
      | oc:fileid      | UUIDof:folderToShare (1) |
      | oc:name        | folderToShare            |
      | oc:permissions | S                        |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "folderToShare (1)" with these key and value pairs:
      | key            | value               |
      | oc:fileid      | UUIDof:textfile.txt |
      | oc:name        | textfile.txt        |
      | oc:permissions | S                   |


  Scenario Outline: check file-id from PROPFIND with shared-with-me drive-item-id
    Given using spaces DAV path
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Brian" has a share "<resource>" synced
    When user "Brian" sends PROPFIND request to space "Shares" with depth "1" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the key "oc:fileid" from PROPFIND response should match with shared-with-me drive-item-id of share "<resource>"
    Examples:
      | resource      |
      | textfile1.txt |
      | folderToShare |
