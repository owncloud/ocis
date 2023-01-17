@api @skipOnOcV10
Feature: Tag
  The user can add a tag to resources for sorting and quick search

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "use-tag" with the default quota using the GraphApi
    And user "Alice" has created a folder "folderMain" in space "use-tag"
    And user "Alice" has uploaded a file inside space "use-tag" with content "some content" to "folderMain/insideTheFolder.txt"


  Scenario: Alice creates tags for resources in the project space
    Given user "Alice" has shared a space "use-tag" to user "Brian" with role "viewer"
    When user "Alice" creates the following tags for folder "folderMain" of space "use-tag":
      | tag level#1                    |
      | tag with symbols @^$#^%$@%!_+) |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "use-tag" to the resource "folderMain" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "use-tag" with these key and value pairs:
      | key     | value                                      |
      | oc:tags | tag level#1,tag with symbols @^$#^%$@%!_+) |
    When user "Alice" creates the following tags for file "folderMain/insideTheFolder.txt" of space "use-tag":
      | fileTag |
    Then the HTTP status code should be "200"
    When user "Brian" sends PROPFIND request from the space "use-tag" to the resource "folderMain/insideTheFolder.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "use-tag" with these key and value pairs:
      | key     | value   |
      | oc:tags | fileTag |
    When user "Alice" lists all available tags via the GraphApi
    Then the HTTP status code should be "200"
    And the response should contain following tags:
      | tag level#1                    |
      | tag with symbols @^$#^%$@%!_+) |
      | fileTag                        |
    When user "Alice" lists all available tags via the GraphApi
    Then the HTTP status code should be "200"
    And the response should contain following tags:
      | tag level#1                    |
      | tag with symbols @^$#^%$@%!_+) |
      | fileTag                        |


  Scenario: Alice creates tags for resources in the personal space
    Given user "Alice" has created a folder "folderMain" in space "Alice Hansen"
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "file.txt"
    When user "Alice" creates the following tags for folder "folderMain" of space "Alice Hansen":
      | my tag    |
      | important |
    Then the HTTP status code should be "200"
    When user "Alice" creates the following tags for file "file.txt" of space "Alice Hansen":
      | fileTag                       |
      | tag with symbol @^$#^%$@%!_+) |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "folderMain" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a mountpoint "Alice Hansen" with these key and value pairs:
      | key     | value            |
      | oc:tags | my tag,important |
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "file.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a mountpoint "Alice Hansen" with these key and value pairs:
      | key     | value                                 |
      | oc:tags | fileTag,tag with symbol @^$#^%$@%!_+) |
    When user "Alice" lists all available tags via the GraphApi
    Then the HTTP status code should be "200"
    And the response should contain following tags:
      | my tag                        |
      | important                     |
      | fileTag                       |
      | tag with symbol @^$#^%$@%!_+) |
