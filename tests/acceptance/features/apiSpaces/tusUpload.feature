Feature: upload resources using TUS protocol
  As a user
  I want to be able to upload files
  So that I can store and share files between multiple client systems

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path


  Scenario: upload a file within the set quota to a project space
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "10000"
    When user "Alice" uploads a file with content "uploaded content" to "/upload.txt" via TUS inside of the space "Project Jupiter" using the WebDAV API
    Then for user "Alice" the space "Project Jupiter" should contain these entries:
      | upload.txt |


  Scenario: upload a file bigger than the set quota to a project space
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "10"
    When user "Alice" creates a new TUS resource for the space "Project Jupiter" with content "file content is 24 bytes" using the WebDAV API with these headers:
      | Upload-Length   | 24                              |
      # dXBsb2FkLnR4dA== is the base64 encoded value of filename upload.txt
      | Upload-Metadata | filename dXBsb2FkLnR4dA==       |
      | Content-Type    | application/offset+octet-stream |
      | Tus-Resumable   | 1.0.0                           |
      | Tus-Extension   | creation-with-upload            |
    Then the HTTP status code should be "507"
    And for user "Alice" the space "Project Jupiter" should not contain these entries:
      | upload.txt |


  Scenario: upload the same file after renaming the first one
    Given user "Alice" has uploaded a file with content "uploaded content" to "/upload.txt" via TUS inside of the space "Alice Hansen"
    And user "Alice" has moved file "upload.txt" to "test.txt" in space "Alice Hansen"
    When user "Alice" uploads a file with content "uploaded content" to "/upload.txt" via TUS inside of the space "Alice Hansen" using the WebDAV API
    Then for user "Alice" the space "Alice Hansen" should contain these entries:
      | test.txt   |
      | upload.txt |
