@notification @email
Feature: Email notification
  As a user
  I want to get email notification of events related to me
  So that I can stay updated about the events

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario: user gets an email notification when someone shares a project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
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


  Scenario: user gets an email notification when someone shares a file
    Given user "Alice" has uploaded file with content "sample text" to "lorem.txt"
    When user "Alice" shares file "lorem.txt" with user "Brian" using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "100"
    And user "Brian" should have received the following email from user "Alice"
      """
      Hello Brian Murphy

      %displayname% has shared "lorem.txt" with you.

      Click here to view it: %base_url%/files/shares/with-me
      """


  Scenario: group members get an email notification when someone shares a project space with the group
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Carol" has been created with default attributes
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    When user "Alice" shares a space "new-space" with settings:
      | shareWith | group1 |
      | shareType | 8      |
      | role      | viewer |
    Then the HTTP status code should be "200"
    And user "Brian" should have received the following email from user "Alice" about the share of project space "new-space"
      """
      Hello Brian Murphy,

      %displayname% has invited you to join "new-space".

      Click here to view it: %base_url%/f/%space_id%
      """
    And user "Carol" should have received the following email from user "Alice" about the share of project space "new-space"
      """
      Hello Carol King,

      %displayname% has invited you to join "new-space".

      Click here to view it: %base_url%/f/%space_id%
      """


  Scenario: group members get an email notification in their respective languages when someone shares a folder with the group
    Given user "Carol" has been created with default attributes
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "Brian" has switched the system language to "es" using the Graph API
    And user "Carol" has switched the system language to "de" using the Graph API
    And user "Alice" has created folder "/HelloWorld"
    When user "Alice" shares folder "HelloWorld" with group "group1" using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "100"
    And user "Brian" should have received the following email from user "Alice"
      """
      Hola Brian Murphy

      %displayname% ha compartido "HelloWorld" contigo.

      Click aquí para verlo: %base_url%/files/shares/with-me
      """
    And user "Carol" should have received the following email from user "Alice"
      """
      Hallo Carol King

      %displayname% hat "HelloWorld" mit Ihnen geteilt.

      Zum Ansehen hier klicken: %base_url%/files/shares/with-me
      """


  Scenario: group members get an email notification in their respective languages when someone shares a file with the group
    Given user "Carol" has been created with default attributes
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "Brian" has switched the system language to "es" using the Graph API
    And user "Carol" has switched the system language to "de" using the Graph API
    And user "Alice" has uploaded file with content "hello world" to "text.txt"
    When user "Alice" shares file "text.txt" with group "group1" using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "100"
    And user "Brian" should have received the following email from user "Alice"
      """
      Hola Brian Murphy

      %displayname% ha compartido "text.txt" contigo.

      Click aquí para verlo: %base_url%/files/shares/with-me
      """
    And user "Carol" should have received the following email from user "Alice"
      """
      Hallo Carol King

      %displayname% hat "text.txt" mit Ihnen geteilt.

      Zum Ansehen hier klicken: %base_url%/files/shares/with-me
      """


  Scenario: group members get an email notification in their respective languages when someone shares a space with the group
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Carol" has been created with default attributes
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "Brian" has switched the system language to "es" using the Graph API
    And user "Carol" has switched the system language to "de" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    When user "Alice" shares a space "new-space" with settings:
      | shareWith | group1 |
      | role      | viewer |
    Then the HTTP status code should be "200"
    And user "Brian" should have received the following email from user "Alice" about the share of project space "new-space"
      """
      Hola Brian Murphy,

      Alice Hansen te ha invitado a unirte a "new-space".

      Click aquí para verlo: %base_url%/f/%space_id%
      """
    And user "Carol" should have received the following email from user "Alice" about the share of project space "new-space"
      """
      Hallo Carol King,

      Alice Hansen hat Sie eingeladen, dem Space "new-space" beizutreten.

      Zum Ansehen hier klicken: %base_url%/f/%space_id%
      """


  Scenario: user gets an email notification when space admin unshares a space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | new-space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    When user "Alice" unshares a space "new-space" to user "Brian"
    Then the HTTP status code should be "200"
    And user "Brian" should have received the following email from user "Alice" about the share of project space "new-space"
      """
      Hello Brian Murphy,

      %displayname% has removed you from "new-space".

      You might still have access through your other groups or direct membership.

      Click here to check it: %base_url%/f/%space_id%
      """

  @issue-10904
  Scenario: user gets an email notification when a folder is unshared (Personal Space)
    Given user "Alice" has created folder "SHARED-FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | SHARED-FOLDER |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Alice" removes the access of user "Brian" from resource "SHARED-FOLDER" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should have received the following email from user "Alice" about the share of project space "SHARED-FOLDER"
      """
      Hello Brian Murphy,

      %displayname% has unshared 'SHARED-FOLDER' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.
      """

  @issue-10904
  Scenario: user gets an email notification when a folder is unshared (Project Space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "shared-space" with the default quota using the Graph API
    And user "Alice" has created a folder "SHARED-FOLDER" in space "shared-space"
    And user "Alice" has sent the following resource share invitation:
      | resource        | SHARED-FOLDER |
      | space           | shared-space  |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Alice" removes the access of user "Brian" from resource "SHARED-FOLDER" of space "shared-space" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should have received the following email from user "Alice" about the share of project space "shared-space"
      """
      Hello Brian Murphy,

      %displayname% has unshared 'SHARED-FOLDER' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.
      """

  @issue-10904
  Scenario: user gets an email notification when a file is unshared (Personal Space)
    Given user "Alice" has uploaded file with content "Sample data" to "file-to-share.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | file-to-share.txt  |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | Viewer             |
    When user "Alice" removes the access of user "Brian" from resource "file-to-share.txt" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should have received the following email from user "Alice" about the share of project space "file-to-share.txt"
      """
      Hello Brian Murphy,

      %displayname% has unshared 'file-to-share.txt' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.
      """

  @issue-10904
  Scenario: user gets an email notification when a file is unshared (Project Space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "shared-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "shared-space" with content "Sample data" to "file-to-share.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | file-to-share.txt |
      | space           | shared-space      |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | Viewer            |
    When user "Alice" removes the access of user "Brian" from resource "file-to-share.txt" of space "shared-space" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should have received the following email from user "Alice" about the share of project space "file-to-share.txt"
      """
      Hello Brian Murphy,

      %displayname% has unshared 'file-to-share.txt' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.
      """

  @env-config
  Scenario: group members get an email notification in default language when someone shares a file with the group
    Given the config "OCIS_DEFAULT_LANGUAGE" has been set to "de"
    And user "Carol" has been created with default attributes
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "Alice" has uploaded file with content "hello world" to "text.txt"
    When user "Alice" shares file "text.txt" with group "group1" using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "100"
    And user "Brian" should have received the following email from user "Alice"
      """
      Hallo Brian Murphy

      %displayname% hat "text.txt" mit Ihnen geteilt.

      Zum Ansehen hier klicken: %base_url%/files/shares/with-me
      """
    And user "Carol" should have received the following email from user "Alice"
      """
      Hallo Carol King

      %displayname% hat "text.txt" mit Ihnen geteilt.

      Zum Ansehen hier klicken: %base_url%/files/shares/with-me
      """

  @issue-9530
  Scenario: user gets an email notification when someone with comma in display name shares a file
    Given the administrator has assigned the role "Admin" to user "Brian" using the Graph API
    And the user "Brian" has created a new user with the following attributes:
      | userName    | Carol             |
      | displayName | Carol, King       |
      | email       | carol@example.com |
      | password    | 1234              |
    And user "Carol" has uploaded file with content "sample text" to "lorem.txt"
    When user "Carol" sends the following resource share invitation using the Graph API:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    Then the HTTP status code should be "200"
    And user "Brian" should have received the following email from user "Carol"
      """
      Hello Brian Murphy

      Carol, King has shared "lorem.txt" with you.

      Click here to view it: %base_url%/files/shares/with-me
      """
