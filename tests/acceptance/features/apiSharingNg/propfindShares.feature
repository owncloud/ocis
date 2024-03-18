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
    And user "Alice" has sent the following share invitation:
      | resource        | <path>   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And user "Carol" has sent the following share invitation:
      | resource        | <path>   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    When user "Brian" sends PROPFIND request to space "Shares" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Shares" with these key and value pairs:
      | key       | value         |
      | oc:fileid | UUIDof:Shares |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key            | value  |
      | oc:name        | <path> |
      | oc:permissions | SR     |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key            | value   |
      | oc:name        | <path2> |
      | oc:permissions | SR      |
    Examples:
      | path          | path2             |
      | textfile.txt  | textfile (1).txt  |
      | folderToShare | folderToShare (1) |

  @issue-4421
  Scenario Outline: sharee PROPFIND same name shares shared by multiple users using new dav path
    Given using new DAV path
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Alice" has created folder "folderToShare"
    And user "Carol" has uploaded file with content "to share" to "textfile.txt"
    And user "Carol" has created folder "folderToShare"
    And user "Alice" has sent the following share invitation:
      | resource        | <path>   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And user "Carol" has sent the following share invitation:
      | resource        | <path>   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "Shares" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Shares" with these key and value pairs:
      | key       | value         |
      | oc:fileid | UUIDof:Shares |
      | oc:name   | Shares        |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key            | value         |
      | oc:fileid      | UUIDof:<path> |
      | oc:name        | <path>        |
      | oc:permissions | SR            |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key            | value          |
      | oc:fileid      | UUIDof:<path2> |
      | oc:name        | <path2>        |
      | oc:permissions | SR             |
    Examples:
      | path          | path2             |
      | textfile.txt  | textfile (1).txt  |
      | folderToShare | folderToShare (1) |
