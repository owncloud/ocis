Feature: Propfind test
  As a user
  I want to check the PROPFIND response
  So that I can make sure that the response contains all the relevant values

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API


  Scenario: space-admin checks the PROPFIND request of a space
    Given user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testfile.txt"
    When user "Alice" sends PROPFIND request to space "new-space" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And the following headers should match these regular expressions
      | X-Request-Id | /^[a-zA-Z]+\/[a-zA-Z]+\.feature:\d+(-\d+)?$/ |
    And as user "Alice" the PROPFIND response should contain a space "new-space" with these key and value pairs:
      | key            | value                     |
      | oc:fileid      | %file_id_pattern%         |
      | oc:name        | new-space                 |
      | oc:permissions | RDNVCKZP                  |
      | oc:privatelink | %base_url%/f/[0-9a-z-$%]+ |
      | oc:size        | 12                        |


  Scenario Outline: space member with a different role checks the PROPFIND request of a space
    Given user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testfile.txt"
    And user "Alice" has sent the following space share invitation:
      | space           | new-space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | <space-role> |
    When user "Brian" sends PROPFIND request to space "new-space" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And the following headers should match these regular expressions
      | X-Request-Id | /^[a-zA-Z]+\/[a-zA-Z]+\.feature:\d+(-\d+)?$/ |
    And as user "Brian" the PROPFIND response should contain a space "new-space" with these key and value pairs:
      | key            | value                     |
      | oc:fileid      | %file_id_pattern%         |
      | oc:name        | new-space                 |
      | oc:permissions | <oc-permission>           |
      | oc:privatelink | %base_url%/f/[0-9a-z-$%]+ |
      | oc:size        | 12                        |
    Examples: 
      | space-role   | oc-permission |
      | Manager      | RDNVCKZP      |
      | Space Editor | DNVCK         |
      | Space Viewer |               |


  Scenario Outline: space member with a different role checks the PROPFIND request of the folder in the space
    Given user "Alice" has created a folder "folderMain" in space "new-space"
    And user "Alice" has sent the following space share invitation:
      | space           | new-space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | <space-role> |
    When user "Brian" sends PROPFIND request from the space "new-space" to the resource "folderMain" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a mountpoint "folderMain" with these key and value pairs:
      | key            | value             |
      | oc:fileid      | %file_id_pattern% |
      | oc:file-parent | %file_id_pattern% |
      | oc:name        | folderMain        |
      | oc:permissions | <oc-permission>   |
      | oc:size        | 0                 |
    Examples:
      | space-role   | oc-permission |
      | Manager      | RDNVCKZP      |
      | Space Editor | DNVCK         |
      | Space Viewer |               |


  Scenario Outline: space member with a different role checks the PROPFIND request of the sub-folder in the space
    Given user "Alice" has created a folder "folderMain/subFolder1/subFolder2" in space "new-space"
    And user "Alice" has sent the following space share invitation:
      | space           | new-space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | <space-role> |
    When user "Brian" sends PROPFIND request from the space "new-space" to the resource "folderMain/subFolder1/subFolder2" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a mountpoint "subFolder2" with these key and value pairs:
      | key            | value             |
      | oc:fileid      | %file_id_pattern% |
      | oc:file-parent | %file_id_pattern% |
      | oc:name        | subFolder2        |
      | oc:permissions | <oc-permission>   |
      | oc:size        | 0                 |
    Examples:
      | space-role   | oc-permission |
      | Manager      | RDNVCKZP      |
      | Space Editor | DNVCK         |
      | Space Viewer |               |


  Scenario Outline: space member with a different role checks the PROPFIND request of the file in the space
    Given user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testfile.txt"
    And user "Alice" has sent the following space share invitation:
      | space           | new-space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | <space-role> |
    When user "Brian" sends PROPFIND request from the space "new-space" to the resource "testfile.txt" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a mountpoint "testfile.txt" with these key and value pairs:
      | key            | value             |
      | oc:fileid      | %file_id_pattern% |
      | oc:file-parent | %file_id_pattern% |
      | oc:name        | testfile.txt      |
      | oc:permissions | <oc-permission>   |
      | oc:size        | 12                |
    Examples:
      | space-role   | oc-permission |
      | Manager      | RDNVWZP       |
      | Space Editor | DNVW          |
      | Space Viewer |               |
