@api @skipOnOcV10
Feature: Report test
  check that the REPORT response contains all the relevant value

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "find data" with the default quota using the GraphApi
    And user "Alice" has created a folder "folderMain/SubFolder1/subFOLDER2" in space "find data"
    And user "Alice" has uploaded a file inside space "find data" with content "some content" to "folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt"
    And using new DAV path


  Scenario: check the response of the found folder
    Given user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    And user "Brian" has accepted share "/folderMain" offered by user "Alice"
    When user "Brian" searches for "SubFolder1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "REPORT" response to user "Brian" should contain a mountpoint "folderMain" with these key and value pairs:
      | key              | value                |
      | oc:fileid        | UUIDof:SubFolder1    |
      | oc:file-parent   | UUIDof:folderMain    |
      | oc:shareroot     | /folderMain          |
      | oc:name          | SubFolder1           |
      | d:getcontenttype | httpd/unix-directory |
      | oc:permissions   | SR                   |
      | oc:size          | 12                   |


  Scenario: check the response of the found file
    Given user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | editor     |
    And user "Brian" has accepted share "/folderMain" offered by user "Alice"
    When user "Brian" searches for "insideTheFolder.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "REPORT" response to user "Brian" should contain a mountpoint "folderMain" with these key and value pairs:
      | key                | value                                            |
      | oc:fileid          | UUIDof:SubFolder1/subFOLDER2/insideTheFolder.txt |
      | oc:file-parent     | UUIDof:SubFolder1/subFOLDER2                     |
      | oc:shareroot       | /folderMain                                      |
      | oc:name            | insideTheFolder.txt                              |
      | d:getcontenttype   | text/plain                                       |
      | oc:permissions     | SRDNVW                                           |
      | d:getcontentlength | 12                                               |
