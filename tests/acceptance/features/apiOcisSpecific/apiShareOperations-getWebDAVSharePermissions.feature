@api @files_sharing-app-required @issue-ocis-reva-47
Feature: sharing

  Background:
    Given using OCS API version "1"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @issue-ocis-reva-47
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Empty webdav share-permissions for owned file
    Given using <dav-path> DAV path
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    When user "Alice" gets the following properties of file "/tmp.txt" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the single response should contain a property "ocs:share-permissions" with value "5"
    Examples:
      | dav-path |
      | old      |
      | new      |
