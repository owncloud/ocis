Feature: REPORT request to Shares space
  As a user
  I want to check the share REPORT response
  So that I can make sure that the response contains all the relevant details for shares

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has created folder "/folderMain"
    And user "Alice" has created folder "/folderMain/SubFolder1"
    And user "Alice" has created folder "/folderMain/SubFolder1/subFOLDER2"
    And user "Alice" has sent the following resource share invitation:
      | resource        | /folderMain |
      | space           | Personal    |
      | sharee          | Brian       |
      | shareType       | user        |
      | permissionsRole | Viewer      |


  Scenario Outline: check the REPORT response of the found folder
    Given using <dav-path-version> DAV path
    When user "Brian" searches for "SubFolder1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the following headers should match these regular expressions
      | X-Request-Id | /^[a-zA-Z]+\/[a-zA-Z]+\.feature:\d+(-\d+)?$/ |
    And the "REPORT" response to user "Brian" should contain a mountpoint "folderMain" with these key and value pairs:
      | key               | value                |
      | oc:fileid         | UUIDof:SubFolder1    |
      | oc:file-parent    | UUIDof:folderMain    |
      | oc:shareroot      | /folderMain          |
      | oc:name           | SubFolder1           |
      | d:getcontenttype  | httpd/unix-directory |
      | oc:permissions    | S                    |
      | oc:remote-item-id | UUIDof:folderMain    |
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check the REPORT response of the found file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "Not all those who wander are lost." to "/folderMain/SubFolder1/subFOLDER2/frodo.txt"
    When user "Brian" searches for "frodo.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the following headers should match these regular expressions
      | X-Request-Id | /^[a-zA-Z]+\/[a-zA-Z]+\.feature:\d+(-\d+)?$/ |
    And the "REPORT" response to user "Brian" should contain a mountpoint "folderMain" with these key and value pairs:
      | key                | value                                  |
      | oc:fileid          | UUIDof:SubFolder1/subFOLDER2/frodo.txt |
      | oc:file-parent     | UUIDof:SubFolder1/subFOLDER2           |
      | oc:shareroot       | /folderMain                            |
      | oc:name            | frodo.txt                              |
      | d:getcontenttype   | text/plain                             |
      | oc:permissions     | S                                      |
      | d:getcontentlength | 34                                     |
      | oc:remote-item-id | UUIDof:folderMain                       |
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: search for the shared folder when share is not accepted
    Given user "Brian" has disabled auto-accepting
    And using <dav-path-version> DAV path
    And user "Alice" has created folder "/folderToBrian"
    And user "Alice" has sent the following resource share invitation:
      | resource        | /folderToBrian |
      | space           | Personal       |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | Viewer         |
    When user "Brian" searches for "folderToBrian" using the WebDAV API
    Then the HTTP status code should be "207"
    And the following headers should match these regular expressions
      | X-Request-Id | /^[a-zA-Z]+\/[a-zA-Z]+\.feature:\d+(-\d+)?$/ |
    And the search result should contain "0" entries
    Examples:
      | dav-path-version |
      | old              |
      | new              |
