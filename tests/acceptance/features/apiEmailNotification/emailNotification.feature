@api @email
Feature: Email notification
  As a user
  I want to get email notification of events related to me
  So that I can stay updated about the events

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario: a user gets an email notification when someone shares a project space
    Given the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "new-space" with the default quota using the GraphApi
    When user "Alice" shares a space "new-space" with settings:
      | shareWith | Brian  |
      | role      | Editor |
    Then the HTTP status code should be "200"
    And user "Brian" should have received the following email from user "Alice" about the share of project space "new-space"
      """
      Hello Brian Murphy,

      %displayname% has invited you to join "new-space".

      Click here to view it: %base_url%/f/%space_id%
      """


  Scenario: a user gets an email notification when someone shares a file
    Given user "Alice" has uploaded file with content "sample text" to "lorem.txt"
    When user "Alice" has shared file "lorem.txt" with user "Brian" with permissions "17"
    Then the HTTP status code should be "200"
    And user "Brian" should have reveived the following email from user "Alice"
      """
      Hello Brian Murphy

      %displayname% has shared "lorem.txt" with you.

      Click here to view it: %base_url%/files/shares/with-me
      """
