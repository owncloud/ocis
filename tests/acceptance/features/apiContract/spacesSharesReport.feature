Feature: Report test
  As a user
  I want to check the share REPORT response
  So that I can make sure that the response contains all the relevant details for shares

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "find data" with the default quota using the Graph API
    And user "Alice" has created a folder "folderMain/SubFolder1/subFOLDER2" in space "find data"
    And user "Alice" has uploaded a file inside space "find data" with content "some content" to "folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt"
    And using new DAV path


  Scenario: check the response of the found folder
    Given user "Alice" has sent the following resource share invitation:
      | resource        | folderMain |
      | space           | find data  |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
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
      | oc:size           | 12                   |
      | oc:remote-item-id | UUIDof:folderMain    |


  Scenario: check the response of the found file
    Given user "Alice" has sent the following resource share invitation:
      | resource        | folderMain |
      | space           | find data  |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Editor     |
    When user "Brian" searches for "insideTheFolder.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the following headers should match these regular expressions
      | X-Request-Id | /^[a-zA-Z]+\/[a-zA-Z]+\.feature:\d+(-\d+)?$/ |
    And the "REPORT" response to user "Brian" should contain a mountpoint "folderMain" with these key and value pairs:
      | key                | value                                            |
      | oc:fileid          | UUIDof:SubFolder1/subFOLDER2/insideTheFolder.txt |
      | oc:file-parent     | UUIDof:SubFolder1/subFOLDER2                     |
      | oc:shareroot       | /folderMain                                      |
      | oc:name            | insideTheFolder.txt                              |
      | d:getcontenttype   | text/plain                                       |
      | oc:permissions     | SD                                               |
      | d:getcontentlength | 12                                               |
      | oc:remote-item-id  | UUIDof:folderMain                                |


  Scenario: search for the shared folder when the share is not accepted
    Given user "Brian" has disabled auto-accepting
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderMain |
      | space           | find data  |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    When user "Brian" searches for "folderMain" using the WebDAV API
    Then the HTTP status code should be "207"
    And the following headers should match these regular expressions
      | X-Request-Id | /^[a-zA-Z]+\/[a-zA-Z]+\.feature:\d+(-\d+)?$/ |
    And the search result should contain "0" entries
