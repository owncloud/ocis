Feature: tests of the creation extension see https://tus.io/protocols/resumable-upload.html#creation-with-upload
  As a user
  I want to be able to include parts of upload while creating resources
  So that I can provide basic information about the resources to the server

  Background:
    Given user "Alice" has been created with default attributes


  Scenario Outline: creating a new upload resource using creation with upload extension
    Given using <dav-path-version> DAV path
    When user "Alice" creates a new TUS resource with content "uploaded content" on the WebDAV API with these headers:
      | Upload-Length   | 16                              |
      | Tus-Resumable   | 1.0.0                           |
      | Content-Type    | application/offset+octet-stream |
      #    dGVzdC50eHQ= is the base64 encode of test.txt
      | Upload-Metadata | filename dGVzdC50eHQ=           |
      | Tus-Extension   | creation-with-upload            |
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions
      | Tus-Resumable | /1\.0\.0/                       |
      | Location      | /http[s]?:\/\/.*:\d+\/data\/.*/ |
      | Upload-Offset | /\d+/                           |
    And the content of file "/test.txt" for user "Alice" should be "uploaded content"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-10346
  Scenario Outline: creating a new resource and upload data in multiple bytes using creation with upload extension
    Given using <dav-path-version> DAV path
    When user "Alice" creates file "textFile.txt" and uploads content "12345" in the same request using the TUS protocol on the WebDAV API
    Then the HTTP status code should be "201"
    And the following headers should be set
      | header                        | value                                  |
      | Access-Control-Expose-Headers | Tus-Resumable, Upload-Offset, Location |
    And the content of file "/textFile.txt" for user "Alice" should be "12345"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |
