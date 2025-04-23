Feature: propfind a shares
  As a user
  I want to check the PROPFIND response
  So that I can make sure that the response contains all the relevant values

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |

  @issue-4421 @issue-9933
  Scenario Outline: sharee PROPFIND same name shares shared by multiple users
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
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "/" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a space "Shares" with these key and value pairs:
      | key       | value             |
      | oc:fileid | %file_id_pattern% |
      | oc:name   | Shares            |
    And as user "Brian" the PROPFIND response should contain a resource "<resource>" with these key and value pairs:
      | key            | value             |
      | oc:fileid      | %file_id_pattern% |
      | oc:name        | <resource>        |
      | oc:permissions | S                 |
    And as user "Brian" the PROPFIND response should contain a resource "<resource-2>" with these key and value pairs:
      | key            | value             |
      | oc:fileid      | %file_id_pattern% |
      | oc:name        | <resource-2>      |
      | oc:permissions | S                 |
    Examples:
      | dav-path-version | resource      | resource-2        |
      | old              | textfile.txt  | textfile (1).txt  |
      | old              | folderToShare | folderToShare (1) |
      | new              | textfile.txt  | textfile (1).txt  |
      | new              | folderToShare | folderToShare (1) |

  @issue-4421 @issue-9933
  Scenario: sharee PROPFIND a share having bracket in the name
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
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare (1)" with these key and value pairs:
      | key            | value              |
      | oc:fileid      | %share_id_pattern% |
      | oc:name        | folderToShare      |
      | oc:permissions | S                  |
    And as user "Brian" the PROPFIND response should contain a resource "textfile.txt" with these key and value pairs:
      | key            | value             |
      | oc:fileid      | %file_id_pattern% |
      | oc:name        | textfile.txt      |
      | oc:permissions |                   |

  @issue-8420
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

  @issue-9933
  Scenario Outline: check file-id of different PROPFIND requests to shared items
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "lorem epsum" to "folderToShare/textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Brian" has a share "folderToShare" synced
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "/" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare" with these key and value pairs:
      | key            | value                               |
      | oc:fileid      | <pattern>                           |
      | oc:file-parent | %self::oc:spaceid%!%uuidv4_pattern% |
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "folderToShare" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare" with these key and value pairs:
      | key            | value                               |
      | oc:fileid      | <pattern>                           |
      | oc:file-parent | %self::oc:spaceid%!%uuidv4_pattern% |
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare/textfile.txt" with these key and value pairs:
      | key            | value                               |
      | oc:fileid      | %file_id_pattern%                   |
      | oc:file-parent | %self::oc:spaceid%!%uuidv4_pattern% |
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "folderToShare/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare/textfile.txt" with these key and value pairs:
      | key            | value                               |
      | oc:fileid      | %file_id_pattern%                   |
      | oc:file-parent | %self::oc:spaceid%!%uuidv4_pattern% |
    Examples:
      | dav-path-version | pattern            |
      | old              | %file_id_pattern%  |
      | new              | %file_id_pattern%  |
      | spaces           | %share_id_pattern% |

  @issue-8510
  Scenario Outline: check share-id from PROPFIND request to shared items
    Given using <dav-path-version> DAV path
    And using SharingNG
    And user "Alice" has uploaded file with content "some content" to "testfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | testfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Brian" has a share "testfile.txt" synced
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "/" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "testfile.txt" with these key and value pairs:
      | key        | value           |
      | oc:shareid | %last_share_id% |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-9463
  Scenario Outline: sharer checks share-types property after sharee is deleted (Personal Space)
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "some content" to "testfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Brian" has a share "<resource>" synced
    And the user "Admin" has deleted a user "Brian"
    When user "Alice" gets the following properties of resource "<resource>" inside space "Personal" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the single response should contain a property "oc:share-types" without a child property "oc:share-type"
    Examples:
      | dav-path-version | resource      |
      | old              | folderToShare |
      | old              | testfile.txt  |
      | new              | folderToShare |
      | new              | testfile.txt  |
      | spaces           | folderToShare |
      | spaces           | testfile.txt  |

  @issue-9463
  Scenario Outline: sharer checks share-types property after sharee is deleted (Project Space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "testfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | NewSpace   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Brian" has a share "<resource>" synced
    And the user "Admin" has deleted a user "Brian"
    When user "Alice" gets the following properties of resource "<resource>" inside space "NewSpace" using the WebDAV API
      | propertyName   |
      | oc:share-types |
    Then the HTTP status code should be "207"
    And the single response should contain a property "oc:share-types" without a child property "oc:share-type"
    Examples:
      | resource      |
      | folderToShare |
      | testfile.txt  |
