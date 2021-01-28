@api
Feature: upload file
  As a user
  I want to be able to upload files when passwords are stored with the full hash difficulty
  So that I can store and share files securely between multiple client systems

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Scenario Outline: upload a file and check download content
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file with content "uploaded content" to "/upload.txt" using the WebDAV API
    Then the content of file "/upload.txt" for user "Alice" should be "uploaded content"
    Examples:
      | ocs_api_version | dav_version |
      | 1               | old         |
      | 1               | new         |
      | 2               | old         |
      | 2               | new         |
