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


  Scenario: sharee PROPFIND a shares when multiple user shares resources with same name
    Given user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Carol" has uploaded file with content "to share" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Carol" has sent the following share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "/" with depth "1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Shares" with these key and value pairs:
      | key       | value         |
      | oc:fileid | UUIDof:Shares |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key              | value        |
      | oc:name          | textfile.txt |
      | oc:permissions   | SR           |
      | oc:size          | 8            |
      | d:getcontenttype | text/plain   |
    And the "PROPFIND" response to user "Brian" should contain a mountpoint "Shares" with these key and value pairs:
      | key              | value            |
      | oc:name          | textfile (1).txt |
      | oc:permissions   | SR               |
      | oc:size          | 8                |
      | d:getcontenttype | text/plain       |
