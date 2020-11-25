@api @issue-ocis-reva-14
Feature: move (rename) file
  As a user
  I want to be able to move and rename files
  So that I can manage my file system

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "text file 0" to "/textfile0.txt"

  @issue-ocis-reva-211 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: rename a file into an invalid filename
    Given using <dav_version> DAV path
    When user "Alice" moves file "/textfile0.txt" to "/a\\a" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/a\\a" should exist
    Examples:
      | dav_version |
      | old         |
      | new         |

  @issue-ocis-reva-211 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Renaming a file to a path with extension .part is possible
    Given using <dav_version> DAV path
    When user "Alice" moves file "/textfile0.txt" to "/textfile0.part" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/textfile0.part" should exist
    Examples:
      | dav_version |
      | old         |
      | new         |

  @skipOnOcis-OC-Storage @issue-ocis-reva-211 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: renaming to a file with special characters
    When user "Alice" moves file "/textfile0.txt" to "/<renamed_file>" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/<renamed_file>" for user "Alice" should be ""
    Examples:
      | renamed_file  |
      | #oc ab?cd=ef# |
