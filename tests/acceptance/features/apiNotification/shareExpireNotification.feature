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

  @issue-10966
  Scenario: check share expired in-app and mail notification for Project space resource
    Given using spaces DAV path
    And using SharingNG
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "testfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource           | testfile.txt         |
      | space              | NewSpace             |
      | sharee             | Brian                |
      | shareType          | user                 |
      | permissionsRole    | Viewer               |
      | expirationDateTime | 2025-07-15T14:00:00Z |
    When user "Alice" expires the last share of resource "testfile.txt" inside of the space "NewSpace"
    Then the HTTP status code should be "200"
    And user "Brian" should get a notification with subject "Membership expired" and message:
      | message                       |
      | Access to Space NewSpace lost |
    And user "Brian" should have "2" emails
