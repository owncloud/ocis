@env-config @email
Feature: get grouped email notification
  As an administrator
  I want to get email notification of grouped events related to me either daily or weekly
  So that I can stay updated about the events either once a day or once a week


  Background:
    Given using spaces DAV path
    And using SharingNG
    And these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "New-Space" with the default quota using the Graph API
    And user "Alice" has uploaded file with content "some data" to "lorem.txt"

  @issue-11690
  Scenario: get daily grouped email notification
    Given user "Brian" has set the email sending interval to "daily" using the settings API
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
    And user "Alice" has triggered the share expiration notification for file "lorem.txt"
    And user "Alice" has sent the following space share invitation:
      | space           | New-Space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    And user "Alice" has removed the access of user "Brian" from space "New-Space"
    And user "Alice" has sent the following space share invitation:
      | space              | New-Space                |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | Space Viewer             |
      | expirationDateTime | 2042-03-25T23:59:59.000Z |
    And user "Alice" has expired the membership of user "Brian" from space "New-Space"
    When the administrator triggers "daily" email notifications using the CLI
    Then the command should be successful
    And the command output should contain "successfully sent SendEmailsEvent"
    And user "Brian" should have received the following grouped email
      """
      Hi Brian Murphy,

      Alice Hansen has shared "lorem.txt" with you.


      Alice Hansen has unshared 'lorem.txt' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.


      Alice Hansen has shared "lorem.txt" with you.


      Your share to lorem.txt has expired at %expiry_date_in_mail%

      Even though this share has been revoked you still might have access through other shares and/or space memberships.


      Alice Hansen has invited you to join "New-Space".


      Alice Hansen has removed you from "New-Space".

      You might still have access through your other groups or direct membership.


      Alice Hansen has invited you to join "New-Space".


      Your membership of space New-Space has expired at %expiry_date_in_mail%

      Even though this membership has expired you still might have access through other shares and/or space memberships
      """

  @issue-11690
  Scenario: get weekly grouped email notification
    Given user "Brian" has set the email sending interval to "weekly" using the settings API
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
    And user "Alice" has triggered the share expiration notification for file "lorem.txt"
    And user "Alice" has sent the following space share invitation:
      | space           | New-Space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    And user "Alice" has removed the access of user "Brian" from space "New-Space"
    And user "Alice" has sent the following space share invitation:
      | space              | New-Space                |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | Space Viewer             |
      | expirationDateTime | 2042-03-25T23:59:59.000Z |
    And user "Alice" has expired the membership of user "Brian" from space "New-Space"
    When the administrator triggers "weekly" email notifications using the CLI
    Then the command should be successful
    And the command output should contain "successfully sent SendEmailsEvent"
    And user "Brian" should have received the following grouped email
      """
      Hi Brian Murphy,

      Alice Hansen has shared "lorem.txt" with you.


      Alice Hansen has unshared 'lorem.txt' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.


      Alice Hansen has shared "lorem.txt" with you.


      Your share to lorem.txt has expired at %expiry_date_in_mail%

      Even though this share has been revoked you still might have access through other shares and/or space memberships.


      Alice Hansen has invited you to join "New-Space".


      Alice Hansen has removed you from "New-Space".

      You might still have access through your other groups or direct membership.


      Alice Hansen has invited you to join "New-Space".


      Your membership of space New-Space has expired at %expiry_date_in_mail%

      Even though this membership has expired you still might have access through other shares and/or space memberships
      """
