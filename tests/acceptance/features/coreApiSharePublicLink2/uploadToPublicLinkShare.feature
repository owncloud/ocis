@issue-1276 @issue-1277

Feature: upload to a public link share
  As a user
  I want to create a public link with upload permission
  So that the recipient can upload resources

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "FOLDER"

  @smokeTest @issue-1267
  Scenario: uploading same file to a public upload-only share multiple times via new API
    # The new API does the auto rename in upload-only folders
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER   |
      | permissions | create   |
      | password    | %public% |
    When the public uploads file "test.txt" with password "%public%" and content "test" using the new public WebDAV API
    When the public uploads file "test.txt" with password "%public%" and content "test2" using the new public WebDAV API
    Then the HTTP status code of responses on all endpoints should be "201"
    And the following headers should match these regular expressions
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And the content of file "/FOLDER/test.txt" for user "Alice" should be "test"
    And the content of file "/FOLDER/test (2).txt" for user "Alice" should be "test2"


  Scenario Outline: uploading file to a public upload-only share using public API that was deleted does not work
    Given using <dav-path-version> DAV path
    And user "Alice" has created a public link share with settings
      | path        | FOLDER   |
      | permissions | create   |
      | password    | %public% |
    And user "Alice" has deleted folder "/FOLDER"
    When the public uploads file "test.txt" with password "%public%" and content "test-file" using the new public WebDAV API
    And the HTTP status code should be "404"

    @issue-1268
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @issue-1269
  Scenario: uploading file to a public read-only share folder with public API does not work
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER   |
      | permissions | read     |
      | password    | %public% |
    When the public uploads file "test.txt" with password "%public%" and content "test-file" using the new public WebDAV API
    And the HTTP status code should be "403"


  Scenario: uploading to a public upload-only share with public API
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER   |
      | permissions | create   |
      | password    | %public% |
    When the public uploads file "test.txt" with password "%public%" and content "test-file" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/FOLDER/test.txt" for user "Alice" should be "test-file"
    And the following headers should match these regular expressions
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |


  Scenario: uploading to a public upload-only share with password with public API
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER   |
      | password    | %public% |
      | permissions | create   |
    When the public uploads file "test.txt" with password "%public%" and content "test-file" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/FOLDER/test.txt" for user "Alice" should be "test-file"


  Scenario: uploading to a public read/write share with password with public API
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER   |
      | password    | %public% |
      | permissions | change   |
    When the public uploads file "test.txt" with password "%public%" and content "test-file" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/FOLDER/test.txt" for user "Alice" should be "test-file"


  Scenario: uploading file to a public shared folder with read/write permission when the sharer has insufficient quota does not work with public API
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER   |
      | permissions | change   |
      | password    | %public% |
    And the quota of user "Alice" has been set to "0"
    When the public uploads file "test.txt" with password "%public%" and content "test2" using the new public WebDAV API
    Then the HTTP status code should be "507"

  @issue-1290
  Scenario: uploading file to a public shared folder with upload-only permission when the sharer has insufficient quota does not work with public API
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER   |
      | permissions | create   |
      | password    | %public% |
    And the quota of user "Alice" has been set to "0"
    When the public uploads file "test.txt" with password "%public%" and content "test2" using the new public WebDAV API
    Then the HTTP status code should be "507"

  @smokeTest
  Scenario: uploading to a public upload-write and no edit and no overwrite share with public API
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER          |
      | permissions | uploadwriteonly |
      | password    | %public%        |
    When the public uploads file "test.txt" with password "%public%" and content "test2" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/FOLDER/test.txt" for user "Alice" should be "test2"

  @smokeTest @issue-1267
  Scenario: uploading same file to a public upload-write and no edit and no overwrite share multiple times with new public API
    Given user "Alice" has created a public link share with settings
      | path        | FOLDER          |
      | permissions | uploadwriteonly |
      | password    | %public%        |
    When the public uploads file "test.txt" with password "%public%" and content "test" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    When the public uploads file "test.txt" with password "%public%" and content "test2" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/FOLDER/test.txt" for user "Alice" should be "test"
    And the content of file "/FOLDER/test (2).txt" for user "Alice" should be "test2"
