@env-config @email
Feature: get email notification via CLI command

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: get daily grouped email notification via CLI command
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "some data" to "lorem.txt"
    And user "Brian" has enabled notification for the following events using the settings API:
      | Email sending interval | daily |
    And using SharingNG
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>                  |
      | space           | Personal                    |
      | sharee          | Brian                       |
      | shareType       | user                        |
      | permissionsRole | <permissions-role>          |
      | expirationDateTime | 2042-01-01T23:59:59.000Z |
    And user "Alice" has expired the last share of resource "<resource>" inside of the space "Personal"
#    And user "Alice" has removed the access of user "Brian" from resource "<resource>" of space "Personal"
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
