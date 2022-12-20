@api @skipOnOcV10
Feature: Copy test
  check that the Copy response contains all the relevant values

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And using spaces DAV path
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "new-space" with the default quota using the GraphApi
    

  Scenario: check the COPY response headers
    Given user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testfile.txt"
    And user "Alice" has created a folder "new" in space "new-space"
    When user "Alice" copies file "testfile.txt" from space "new-space" to "/new/testfile.txt" inside space "new-space" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the COPY response headers for the space "new-space" should contain these key and value pairs:
      | key                         | value                    |
      | Access-Control-Allow-Origin | *                        |
      | Oc-Fileid                   | UUIDof:/new/testfile.txt |
