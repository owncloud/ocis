@api @skipOnOcV10
Feature: Propfind test
  check that the PROPFIND response contains all the relevant values

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And using spaces DAV path
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "new-space" with the default quota using the GraphApi


  Scenario: check the PROPFIND request of a space
    Given user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testfile.txt"
    When user "Alice" sends PROPFIND request to space "new-space" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "new-space" with these key and value pairs:
      | key                          | value            |
      | oc:fileid                    | UUIDof:new-space |
      | oc:name                      | new-space        |
      | oc:permissions               | SRDNVCKZ         |
      | oc:privatelink               |                  |
      | oc:size                      | 12               |
