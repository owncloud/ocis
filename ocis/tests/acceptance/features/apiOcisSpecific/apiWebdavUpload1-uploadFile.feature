@api
Feature: upload file
  As a user
  I want to be able to upload files
  So that I can store and share files between multiple client systems

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @skipOnOcis-OC-Storage @issue-ocis-reva-265 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: upload a file and check download content
    Given using <dav_version> DAV path
    When user "Alice" uploads file with content "uploaded content" to <file_name> using the WebDAV API
    Then the content of file <file_name> for user "Alice" should be ""
    Examples:
      | dav_version | file_name           |
      | old         | "file ?2.txt"       |
      | new         | "file ?2.txt"       |

  @skipOnOcis-OC-Storage @issue-product-127 @skipOnOcis-OCIS-Storage
  # this scenario passes/fails intermittently on OC storage, so do not run it in CI
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: uploading a file inside a folder changes its etag
    Given using <dav_version> DAV path
    And user "Alice" has created folder "/upload"
    And user "Alice" has stored etag of element "/<element>"
    When user "Alice" uploads file with content "uploaded content" to "/upload/file.txt" using the WebDAV API
    Then the content of file "/upload/file.txt" for user "Alice" should be "uploaded content"
#    And the etag of element "/<element>" of user "Alice" should have changed
    And the etag of element "/<element>" of user "Alice" should not have changed
    Examples:
      | dav_version | element |
      | old         |         |
      | old         | upload  |
      | new         |         |
      | new         | upload  |
