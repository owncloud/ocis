@env-config
Feature: brute-force protection for password protected public links
  As the maintainer of a password protected public link
  I want repeated wrong password guesses against that link to be rate limited
  So that an attacker cannot brute-force the password regardless of which endpoint they use

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
    And using spaces DAV path
    And the config "STORAGE_PUBLICLINK_BRUTEFORCE_MAXATTEMPTS" has been set to "2" for "storage-publiclink" service
    And using SharingNG


  Scenario: brute-force protection must continue to apply to failed public WebDAV downloads of a password protected public link
    Given user "Alice" has uploaded file with content "some content" to "testfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | testfile.txt |
      | space           | Personal     |
      | permissionsRole | View         |
      | password        | %public%     |
    Then the public should be able to download file "testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "%public%"
    And the public download of file "testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong-password" should fail with HTTP status code "401"
    And the public download of file "testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong-password" should fail with HTTP status code "401"
    And the public download of file "testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong-password" should fail with HTTP status code "429"
    And the public download of file "testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "%public%" should fail with HTTP status code "429"


  Scenario: brute-force protection must apply to failed archiver downloads of a password protected public link
    Given user "Alice" has created folder "public-folder"
    And user "Alice" has uploaded file with content "some content" to "public-folder/testfile.txt"
    And user "Alice" has created the following resource link share:
      | resource        | public-folder |
      | space           | Personal      |
      | permissionsRole | View          |
      | password        | %public%      |
    When the public downloads the archive of the last created public link with password "%public%"
    Then the HTTP status code should be "200"
    When the public downloads the archive of the last created public link with incorrect password "wrong-pw"
    Then the HTTP status code should be "401"
    When the public downloads the archive of the last created public link with incorrect password "wrong-pw"
    Then the HTTP status code should be "401"
    When the public downloads the archive of the last created public link with incorrect password "wrong-pw"
    Then the HTTP status code should be "401"
    When the public downloads the archive of the last created public link with password "%public%"
    Then the HTTP status code should be "401"
