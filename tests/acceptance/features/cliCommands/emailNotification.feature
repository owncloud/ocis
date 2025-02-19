@env-config @email @notification @issue-11001
Feature: get grouped email notification
  As a user
  I want to get email notification of grouped events related to me
  So that I can stay updated about the events once a day or a week

  Background:
    Given using spaces DAV path
    And using SharingNG
    And these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "share space" with the default quota using the Graph API
    And user "Alice" has uploaded file with content "some data" to "lorem.txt"


  Scenario: get daily grouped email notification
    Given user "Brian" has set the email sending interval to "daily" using the settings API
    And user "Alice" has sent the following resource share invitation:
      | resource           | lorem.txt                |
      | space              | Personal                 |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | Viewer                   |
      | expirationDateTime | 2042-01-01T23:59:59.000Z |
    And user "Alice" has expired the last share of resource "lorem.txt" inside of the space "Personal"
    When the administrator triggers "daily" email notifications using the CLI
    Then the command should be successful
    And the command output should contain "successfully sent SendEmailsEvent"
    And user "Brian" should have received the following email from user "Alice"
      """
      Hi Brian Murphy,

      Alice Hansen has shared "lorem.txt" with you.


      Your membership of space Alice Hansen has expired at 2025-02-16 00:00:00

      Even though this membership has expired you still might have access through other shares and/or space memberships
      """


  Scenario: get weekly grouped email notification
    Given user "Brian" has set the email sending interval to "weekly" using the settings API
    And user "Alice" has sent the following resource share invitation:
      | resource           | lorem.txt                |
      | space              | Personal                 |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | Viewer                   |
      | expirationDateTime | 2042-01-01T23:59:59.000Z |
    And user "Alice" has expired the last share of resource "lorem.txt" inside of the space "Personal"
    When the administrator triggers "weekly" email notifications using the CLI
    Then the command should be successful
    And the command output should contain "successfully sent SendEmailsEvent"
    And user "Brian" should have received the following email from user "Alice"
      """
      Hi Brian Murphy,

      Alice Hansen has shared "lorem.txt" with you.


      Your membership of space Alice Hansen has expired at 2025-02-16 00:00:00

      Even though this membership has expired you still might have access through other shares and/or space memberships
      """
