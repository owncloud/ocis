@api
Feature: create folder
  As a user
  I want to be able to create folders
  So that I can organise the files in my file system

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-reva-168 @skipOnOcis-EOS-Storage @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: try to create a folder that already exists
    Given using <dav_version> DAV path
    And user "Alice" has created folder "my-data"
    When user "Alice" creates folder "my-data" using the WebDAV API
    Then the HTTP status code should be "405"
    And as "Alice" folder "my-data" should exist
    And the body of the response should be empty
    Examples:
      | dav_version |
      | old         |
      | new         |

  @issue-ocis-reva-168
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: try to create a folder with a name of an existing file
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with content "uploaded data" to "/my-data.txt"
    When user "Alice" creates folder "my-data.txt" using the WebDAV API
    Then the HTTP status code should be "405"
    And the body of the response should be empty
    And the content of file "/my-data.txt" for user "Alice" should be "uploaded data"
    Examples:
      | dav_version |
      | old         |
      | new         |
