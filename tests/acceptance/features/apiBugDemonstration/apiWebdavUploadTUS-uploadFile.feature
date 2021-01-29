@api @issue-ocis-1141
# after fixing all issues delete these Scenarios and use the one from oC10 core
Feature: upload file
  As a user
  I want to be able to upload files
  So that I can store and share files between multiple client systems

  Scenario Outline: upload a file using the resource URL of another user
    Given using <dav_version> DAV path
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created a new TUS resource on the WebDAV API with these headers:
      | Upload-Length   | 5                         |
      | Upload-Metadata | filename dGV4dEZpbGUudHh0 |
    When user "Brian" sends a chunk to the last created TUS Location with offset "0" and data "12345" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Alice" file "/textFile.txt" should exist
    Examples:
      | dav_version |
      | old         |
      | new         |
