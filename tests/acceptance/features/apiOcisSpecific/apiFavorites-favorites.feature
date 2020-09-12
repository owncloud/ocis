@api
Feature: favorite

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "some data" to "/textfile1.txt"
    And user "Alice" has uploaded file with content "some data" to "/textfile2.txt"
    And user "Alice" has uploaded file with content "some data" to "/textfile3.txt"
    And user "Alice" has uploaded file with content "some data" to "/textfile4.txt"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has uploaded file with content "some data" to "/PARENT/parent.txt"

  @skipOnOcis-OC-Storage @skipOnOcis-OCIS-Storage @issue-ocis-reva-276
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Favorite a folder
    Given using <dav_version> DAV path
    When user "Alice" favorites element "/FOLDER" using the WebDAV API
    Then the HTTP status code should be "500"
    Examples:
      | dav_version |
      | old         |
      | new         |
