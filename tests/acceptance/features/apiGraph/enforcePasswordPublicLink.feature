@api @env-config
Feature: enforce password on public link
  As a user
  I want to enforce passwords on public links shared with upload, edit, or contribute permission
  So that the password is required to access the contents in the link

  Background:
    Given the config "OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD" has been set to "true"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"


  Scenario: user tries to update a public link to edit permission without a password when enforce-password is enabled
    Given user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3 |
    Then the OCS status code should be "996"
    And the OCS status message should be "Error sending update request to public link provider: the public share needs to have a password"


  Scenario: user tries to update a public link to edit permission with a password when enforce-password is enabled
    Given user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3    |
      | password    | 1234 |
    Then the OCS status code should be "100"
    And the OCS status message should be "OK"
    When the public accesses the preview of file "/textfile.txt" from the last shared public link using the sharing API
    Then the HTTP status code should be "404"
