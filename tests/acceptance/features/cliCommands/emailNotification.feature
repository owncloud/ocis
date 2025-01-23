@env-config @email
Feature: get email notification via CLI command

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: get daily/weekly email notification when someone shares a resource
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "some data" to "lorem.txt"
    And user "Brian" has switched the email sending interval to "daily" using the settings API
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>          |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When the administrator triggers "daily" email notifications using the CLI
    Then the command should be successful
    And the command output should contain "successfully sent SendEmailsEvent"
    Examples:
      | permissions-role | resource       |
      | Viewer           | lorem.txt      |
      | File Editor      | lorem.txt      |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |
