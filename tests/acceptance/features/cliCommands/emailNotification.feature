@env-config @email
Feature: get email notification via CLI command

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario: get daily grouped email notification via CLI command
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "share space" with the default quota using the Graph API
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "some data" to "lorem.txt"
    And user "Brian" has enabled notification for the following events using the settings API:
      | Email sending interval | daily |
    And using SharingNG
    And user "Alice" has sent the following resource share invitation:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    And user "Alice" has removed the access of user "Brian" from resource "lorem.txt" of space "Personal"
    And user "Alice" has sent the following resource share invitation:
      | resource           | lorem.txt                |
      | space              | Personal                 |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | Viewer                   |
      | expirationDateTime | 2042-01-01T23:59:59.000Z |
    And user "Alice" has expired the last share of resource "lorem.txt" inside of the space "Personal"
    And using spaces DAV path
    And user "Alice" has sent the following space share invitation:
      | space           | share space  |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    And user "Alice" has removed the access of user "Brian" from space "share space"
    And user "Alice" has sent the following space share invitation:
      | space              | share space              |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | Space Viewer             |
      | expirationDateTime | 2042-03-25T23:59:59.000Z |
    And user "Alice" has expired the user share of space "share space" for user "Brian"
    When the administrator triggers "daily" email notifications using the CLI
    Then the command should be successful
    And the command output should contain "successfully sent SendEmailsEvent"
    And user "Brian" should have received the following email from user "Alice"
      """
      Hi Brian Murphy,

      %displayname% has shared "lorem.txt" with you.


      Alice Hansen has unshared 'lorem.txt' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.


      Alice Hansen has shared "lorem.txt" with you.


      Alice Hansen has invited you to join "share space".


      Alice Hansen has removed you from "share space".

      You might still have access through your other groups or direct membership.


      Alice Hansen has invited you to join "share space".


      Your membership of space share space has expired at %date-time-pattern%

      Even though this membership has expired you still might have access through other shares and/or space memberships
      """


  Scenario: get weekly grouped email notification via CLI command
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "share space" with the default quota using the Graph API
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "some data" to "lorem.txt"
    And user "Brian" has enabled notification for the following events using the settings API:
      | Email sending interval | weekly |
    And using SharingNG
    And user "Alice" has sent the following resource share invitation:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    And user "Alice" has removed the access of user "Brian" from resource "lorem.txt" of space "Personal"
    And user "Alice" has sent the following resource share invitation:
      | resource           | lorem.txt                |
      | space              | Personal                 |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | Viewer                   |
      | expirationDateTime | 2042-01-01T23:59:59.000Z |
    And user "Alice" has expired the last share of resource "lorem.txt" inside of the space "Personal"
    And using spaces DAV path
    And user "Alice" has sent the following space share invitation:
      | space           | share space  |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    And user "Alice" has removed the access of user "Brian" from space "share space"
    And user "Alice" has sent the following space share invitation:
      | space              | share space              |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | Space Viewer             |
      | expirationDateTime | 2042-03-25T23:59:59.000Z |
    And user "Alice" has expired the user share of space "share space" for user "Brian"
    When the administrator triggers "weekly" email notifications using the CLI
    Then the command should be successful
    And the command output should contain "successfully sent SendEmailsEvent"
    And user "Brian" should have received the following email from user "Alice"
      """
      Hi Brian Murphy,

      Alice Hansen has shared "lorem.txt" with you.


      Alice Hansen has unshared 'lorem.txt' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.


      Alice Hansen has shared "lorem.txt" with you.


      Alice Hansen has invited you to join "share space".


      Alice Hansen has removed you from "share space".

      You might still have access through your other groups or direct membership.


      Alice Hansen has invited you to join "share space".


      Your membership of space share space has expired at %date_time_pattern%

      Even though this membership has expired you still might have access through other shares and/or space memberships
      """
