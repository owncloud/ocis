@api
Feature: REPORT request to Shares space
  Check that the REPORT response contains all relevant details for Shares

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has created folder "/folderMain"
    And user "Alice" has created folder "/folderMain/SubFolder1"
    And user "Alice" has created folder "/folderMain/SubFolder1/subFOLDER2"
    And user "Alice" has shared entry "/folderMain" with user "Brian" with permissions "17"
    And user "Brian" has accepted share "/folderMain" offered by user "Alice"


  Scenario Outline: Check the REPORT response of the found folder
    Given using <dav_version> DAV path
    When user "Brian" searches for "SubFolder1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "REPORT" response to user "Brian" should contain a mountpoint "folderMain" with these key and value pairs:
      | key              | value                |
      | oc:fileid        | UUIDof:SubFolder1    |
      | oc:file-parent   | UUIDof:folderMain    |
      | oc:shareroot     | /folderMain          |
      | oc:name          | SubFolder1           |
      | d:getcontenttype | httpd/unix-directory |
      | oc:permissions   | S                    |
    Examples:
      | dav_version |
      | old         |
      | new         |


  Scenario Outline: Check the REPORT response of the found file
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with content "Not all those who wander are lost." to "/folderMain/SubFolder1/subFOLDER2/frodo.txt"
    When user "Brian" searches for "frodo.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "REPORT" response to user "Brian" should contain a mountpoint "folderMain" with these key and value pairs:
      | key                | value                                  |
      | oc:fileid          | UUIDof:SubFolder1/subFOLDER2/frodo.txt |
      | oc:file-parent     | UUIDof:SubFolder1/subFOLDER2           |
      | oc:shareroot       | /folderMain                            |
      | oc:name            | frodo.txt                              |
      | d:getcontenttype   | text/plain                             |
      | oc:permissions     | S                                      |
      | d:getcontentlength | 34                                     |
    Examples:
      | dav_version |
      | old         |
      | new         |
