@api @issue-ocis-reva-14
Feature: move (rename) folder
  As a user
  I want to be able to move and rename folders
  So that I can quickly manage my file system

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-reva-211
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Renaming a folder to a backslash is allowed
    Given using <dav_version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "\" using the WebDAV API
    Then the HTTP status code should be "201" or "500"
    Examples:
      | dav_version |
      | old         |
      | new         |

  @issue-ocis-reva-211
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Renaming a folder beginning with a backslash is allowed
    Given using <dav_version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "\testshare" using the WebDAV API
    Then the HTTP status code should be "201" or "500"
    Examples:
      | dav_version |
      | old         |
      | new         |

  @issue-ocis-reva-211
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Renaming a folder including a backslash encoded is allowed
    Given using <dav_version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "/hola\hola" using the WebDAV API
    Then the HTTP status code should be "201" or "500"
    Examples:
      | dav_version |
      | old         |
      | new         |
