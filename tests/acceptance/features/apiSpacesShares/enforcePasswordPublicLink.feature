@env-config
Feature: enforce password on public link
  As a user
  I want to enforce passwords on public links shared with upload, edit, or contribute permission
  So that the password is required to access the contents of the link

  Background:
    Given the config "OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD" has been set to "true"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"


  Scenario Outline: create a public link with edit permission without a password when enforce-password is enabled
    Given using OCS API version "<ocs-api-version>"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | /testfile.txt |
      | permissions | 3             |
    Then the HTTP status code should be "<http-code>"
    And the OCS status code should be "400"
    And the OCS status message should be "missing required password"
    Examples:
      | ocs-api-version | http-code |
      | 1               | 200       |
      | 2               | 400       |


  Scenario Outline: update a public link to edit permission without a password
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3 |
    Then the HTTP status code should be "<http-code>"
    And the OCS status code should be "400"
    And the OCS status message should be "missing required password"
    Examples:
      | ocs-api-version | http-code |
      | 1               | 200       |
      | 2               | 400       |


  Scenario Outline: updates a public link to edit permission with a password
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3            |
      | password    | testpassword |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs-code>"
    And the OCS status message should be "OK"
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API without a password
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "wrong pass"
    But the public should be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "testpassword"
    Examples:
      | ocs-api-version | ocs-code |
      | 1               | 100      |
      | 2               | 200      |
