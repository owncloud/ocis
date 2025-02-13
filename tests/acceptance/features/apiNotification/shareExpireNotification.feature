@notification @email
Feature: Share Expiry Notification
  As a user
  I want to be notified when share expires
  So that I can stay updated about the share

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario: check share expired in-app and mail notifications for Personal space resource
    Given using SharingNG
    And user "Alice" has uploaded file with content "hello world" to "testfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource           | testfile.txt         |
      | space              | Personal             |
      | sharee             | Brian                |
      | shareType          | user                 |
      | permissionsRole    | Viewer               |
      | expirationDateTime | 2025-07-15T14:00:00Z |
    When user "Alice" expires the last share of resource "testfile.txt" inside of the space "Personal"
    Then the HTTP status code should be "200"
    And user "Brian" should get a notification with subject "Membership expired" and message:
      | message                           |
      | Access to Space Alice Hansen lost |
    And user "Brian" should have "2" emails
    And user "Brian" should have received the following email from user "Alice"
      """
      Hello Brian Murphy,

      Your membership of space %displayname% has expired at %expirationDateTime%

      Even though this membership has expired you still might have access through other shares and/or space memberships
      """
