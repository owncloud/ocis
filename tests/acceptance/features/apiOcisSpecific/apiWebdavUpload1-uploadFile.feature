@api
Feature: upload file
  As a user
  I want to be able to upload files
  So that I can store and share files between multiple client systems

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @skipOnOcis-OC-Storage @issue-ocis-reva-265
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: upload a file and check download content
    Given using <dav_version> DAV path
    When user "Alice" uploads file with content "uploaded content" to <file_name> using the WebDAV API
    Then the content of file <file_name> for user "Alice" should be ""
    Examples:
      | dav_version | file_name           |
      | old         | "file ?2.txt"       |
      | new         | "file ?2.txt"       |
