@ocm @issue-10262
Feature: propfind a federated share
  As a user
  I want to check the PROPFIND response
  So that I can make sure that the response contains all the relevant values

  Background:
    Given user "Alice" has been created with default attributes
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has accepted invitation
    And using server "LOCAL"


  Scenario Outline: sharer checks share-types property of a shared file shared to federated user (Personal Space)
    Given user "Alice" has uploaded file with content "ocm test" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt       |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" gets the following properties of file "textfile.txt" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the response should contain a share-types property with
      | 0 |
    And user "Alice" has removed the access of user "Brian" from resource "textfile.txt" of space "Personal"
    When user "Alice" gets the following properties of file "textfile.txt" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the response should contain an empty property "oc:share-types"
    Examples:
      | permissions-role |
      | Viewer           |
      | File Editor      |


  Scenario Outline: sharer checks share-types property of a shared folder shared to federated user (Personal Space)
    Given user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" gets the following properties of folder "folderToShare" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the response should contain a share-types property with
      | 0 |
    And user "Alice" has removed the access of user "Brian" from resource "folderToShare" of space "Personal"
    When user "Alice" gets the following properties of folder "folderToShare" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the response should contain an empty property "oc:share-types"
    Examples:
      | permissions-role |
      | Viewer           |
      | Editor           |
      | Uploader         |


  Scenario Outline: sharer checks share-types property of a shared file shared to federated user (Project Space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt       |
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" gets the following properties of file "textfile.txt" inside space "new-space" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the response should contain a share-types property with
      | 0 |
    And user "Alice" has removed the access of user "Brian" from resource "textfile.txt" of space "new-space"
    When user "Alice" gets the following properties of file "textfile.txt" inside space "new-space" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the response should contain an empty property "oc:share-types"
    Examples:
      | permissions-role |
      | Viewer           |
      | File Editor      |


  Scenario Outline: sharer checks share-types property of a shared folder shared to federated user (Project Space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare      |
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" gets the following properties of folder "folderToShare" inside space "new-space" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the response should contain a share-types property with
      | 0 |
    And user "Alice" has removed the access of user "Brian" from resource "folderToShare" of space "new-space"
    When user "Alice" gets the following properties of folder "folderToShare" inside space "new-space" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the response should contain an empty property "oc:share-types"
    Examples:
      | permissions-role |
      | Viewer           |
      | Editor           |
      | Uploader         |
