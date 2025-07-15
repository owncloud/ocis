@ocm @env-config @issue-10824
Feature: enable/disable permissions role
  As a user
  I want to enable/disable permissions role on shared resources
  So that I can control the accessibility of shared resources to sharee

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And the administrator has enabled the permissions role "Secure Viewer"

  @issue-10822
  Scenario: user accesses federated shared file shared with permissions role Secure Viewer after the role is disabled (Personal Space)
    Given user "Alice" has uploaded file with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt  |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And the administrator has disabled the permissions role "Secure Viewer"
    And using server "REMOTE"
    When user "Brian" sends PROPFIND request to federated share "textfile.txt" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "textfile.txt" with these key and value pairs:
      | key            | value        |
      | oc:name        | textfile.txt |
      | oc:permissions |              |
    And user "Brian" should have a federated share "textfile.txt" shared by user "Alice" from space "Personal"
    And user "Brian" should be able to download federated shared file "textfile.txt"


  Scenario: user accesses federated shared folder shared with permissions role Secure Viewer after the role is disabled (Personal Space)
    Given user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And the administrator has disabled the permissions role "Secure Viewer"
    And using server "REMOTE"
    When user "Brian" sends PROPFIND request to federated share "folderToShare" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare" with these key and value pairs:
      | key            | value         |
      | oc:name        | folderToShare |
      | oc:permissions |               |
    And user "Brian" should have a federated share "folderToShare" shared by user "Alice" from space "Personal"
    And user "Brian" should be able to download archive of federated shared folder "folderToShare"

  @issue-10822
  Scenario: user accesses federated shared file shared with permissions role Secure Viewer after the role is disabled (Project Space)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt  |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And the administrator has disabled the permissions role "Secure Viewer"
    And using server "REMOTE"
    When user "Brian" sends PROPFIND request to federated share "textfile.txt" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "textfile.txt" with these key and value pairs:
      | key            | value        |
      | oc:name        | textfile.txt |
      | oc:permissions |              |
    And user "Brian" should have a federated share "textfile.txt" shared by user "Alice" from space "new-space"
    And user "Brian" should be able to download federated shared file "textfile.txt"


  Scenario: user accesses federated shared folder shared with permissions role Secure Viewer after the role is disabled (Project Space)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And the administrator has disabled the permissions role "Secure Viewer"
    And using server "REMOTE"
    When user "Brian" sends PROPFIND request to federated share "folderToShare" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare" with these key and value pairs:
      | key            | value         |
      | oc:name        | folderToShare |
      | oc:permissions |               |
    And user "Brian" should have a federated share "folderToShare" shared by user "Alice" from space "new-space"
    And user "Brian" should be able to download archive of federated shared folder "folderToShare"
