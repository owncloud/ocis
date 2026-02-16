@notification
Feature: Notification
  As a user
  I want to be notified of various events
  So that I can stay updated about the information

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And user "Alice" has uploaded file with content "other data" to "/textfile1.txt"
    And user "Alice" has created folder "my_data"

  @issue-10937 @email
  Scenario Outline: user gets in-app and mail notifications of resource sharing
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "Resource shared" and the message-details should match
      """
      {
      "type": "object",
        "required": [
          "app",
          "datetime",
          "message",
          "messageRich",
          "messageRichParameters",
          "notification_id",
          "object_id",
          "object_type",
          "subject",
          "subjectRich",
          "user"
        ],
        "properties": {
          "app": {"const": "userlog"},
          "message": {"const": "Alice Hansen shared <resource> with you"},
          "messageRich": {"const": "{user} shared {resource} with you"},
          "messageRichParameters": {
            "type": "object",
            "required": ["resource","user"],
            "properties": {
              "resource": {
                "type": "object",
                "required": ["id","name"],
                "properties": {
                  "id": {"pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"},
                  "name": {"const": "<resource>"}
                }
              },
              "user": {
                "type": "object",
                "required": ["displayname","id","name"],
                "properties": {
                  "displayname": {"const": "Alice Hansen"},
                  "id": {"pattern": "^%user_id_pattern%$"},
                  "name": {"const": "Alice"}
                }
              }
            }
          },
          "notification_id": {"type": "string"},
          "object_id": {"type": "string"},
          "object_type": {"const": "share"},
          "subject": {"const": "Resource shared"},
          "subjectRich": {"const": "Resource shared"},
          "user": {"const": "Alice"}
        }
      }
      """
    And user "Brian" should have received the following email from user "Alice"
      """
      Hello Brian Murphy

      %displayname% has shared "<resource>" with you.

      Click here to view it: %base_url%/files/shares/with-me
      """
    Examples:
      | user-role     | resource      |
      | User          | textfile1.txt |
      | User          | my_data       |
      | User Light    | textfile1.txt |
      | User Light    | my_data       |

  @issue-10937 @email
  Scenario Outline: user gets a notification about a file share in various languages
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Brian" has switched the system language to "<language>" using the <api> API
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile1.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "<subject>" and the message-details should match
      """
      {
        "type": "object",
        "required": [
          "message"
        ],
        "properties": {
          "message": {
            "type": "string",
            "enum": [
              "<message>"
            ]
          }
        }
      }
      """
    Examples:
      | user-role     | language | subject            | message                                          | api      |
      | User          | de       | Neue Freigabe      | Alice Hansen hat textfile1.txt mit Ihnen geteilt | Graph    |
      | User          | de       | Neue Freigabe      | Alice Hansen hat textfile1.txt mit Ihnen geteilt | settings |
      | User          | es       | Recurso compartido | Alice Hansen compartió textfile1.txt contigo     | Graph    |
      | User          | es       | Recurso compartido | Alice Hansen compartió textfile1.txt contigo     | settings |
      | User Light    | de       | Neue Freigabe      | Alice Hansen hat textfile1.txt mit Ihnen geteilt | Graph    |
      | User Light    | de       | Neue Freigabe      | Alice Hansen hat textfile1.txt mit Ihnen geteilt | settings |
      | User Light    | es       | Recurso compartido | Alice Hansen compartió textfile1.txt contigo     | Graph    |
      | User Light    | es       | Recurso compartido | Alice Hansen compartió textfile1.txt contigo     | settings |

  @env-config
  Scenario: get a notification about a file share in default languages
    Given the config "OCIS_DEFAULT_LANGUAGE" has been set to "de" for "settings" service
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile1.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "Neue Freigabe" and the message-details should match
      """
      {
        "type": "object",
        "required": [
          "message"
        ],
        "properties": {
          "message": {
            "type": "string",
            "enum": [
              "Alice Hansen hat textfile1.txt mit Ihnen geteilt"
            ]
          }
        }
      }
      """


  Scenario Outline: in-app notifications related to a resource get deleted when the resource is deleted
    Given user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Alice" has removed the access of user "Brian" from resource "<resource>" of space "Personal"
    And user "Alice" has deleted entity "/<resource>"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the notifications should be empty
    Examples:
      | resource      |
      | textfile1.txt |
      | my_data       |

  @issue-10937 @issue-10966 @email
  Scenario Outline: check share expired in-app and mail notifications for Personal space file
    Given using SharingNG
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Alice" has uploaded file with content "hello world" to "testfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | testfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    When user "Alice" expires the last share of resource "testfile.txt" inside of the space "Personal"
    Then the HTTP status code should be "200"
    And as "Brian" file "Shares/testfile.txt" should not exist
    And user "Brian" should get a notification with subject "Share expired" and message:
      | message                        |
      | Access to testfile.txt expired |
    And user "Brian" should have received the following email from user "Alice"
      """
      Hello Brian Murphy,

      Your share to testfile.txt has expired at %expiry_date_in_mail%

      Even though this share has been revoked you still might have access through other shares and/or space memberships.
      """
    Examples:
      | user-role  |
      | User       |
      | User Light |

  @issue-10937 @issue-10966 @email
  Scenario Outline: check share expired in-app and mail notifications for Personal space folder
    Given using SharingNG
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Alice" expires the last share of resource "folderToShare" inside of the space "Personal"
    Then the HTTP status code should be "200"
    And as "Brian" file "Shares/folderToShare" should not exist
    And user "Brian" should get a notification with subject "Share expired" and message:
      | message                         |
      | Access to folderToShare expired |
    And user "Brian" should have received the following email from user "Alice"
      """
      Hello Brian Murphy,

      Your share to folderToShare has expired at %expiry_date_in_mail%

      Even though this share has been revoked you still might have access through other shares and/or space memberships.
      """
    Examples:
      | user-role  |
      | User       |
      | User Light |

  @issue-10904 @issue-10937 @email
  Scenario Outline: user gets an in-app and mail notifications of unsharing resource
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    When user "Alice" removes the access of user "Brian" from resource "<resource>" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should get a notification with subject "Resource unshared" and message:
      | message                                   |
      | Alice Hansen unshared <resource> with you |
    And user "Brian" should have received the following email from user "Alice" about the share of project space "<resource>"
      """
      Hello Brian Murphy,

      %displayname% has unshared '<resource>' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.
      """
    Examples:
      | user-role  | resource      |
      | User       | textfile1.txt |
      | User       | my_data       |
      | User Light | textfile1.txt |
      | User Light | my_data       |

  @issue-10904 @issue-10937 @email
  Scenario Outline: user gets in-app and mail notifications when a resource is unshared (Project Space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "shared-space" with the default quota using the Graph API
    And user "Alice" has created a folder "SHARED-FOLDER" in space "shared-space"
    And user "Alice" has uploaded a file inside space "shared-space" with content "Sample data" to "file-to-share.txt"
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>   |
      | space           | shared-space |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    When user "Alice" removes the access of user "Brian" from resource "<resource>" of space "shared-space" using the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should get a notification with subject "Resource unshared" and message:
      | message                                   |
      | Alice Hansen unshared <resource> with you |
    And user "Brian" should have received the following email from user "Alice" about the share of project space "<resource>"
      """
      Hello Brian Murphy,

      %displayname% has unshared '<resource>' with you.

      Even though this share has been revoked you still might have access through other shares and/or space memberships.
      """
    Examples:
      | user-role  | resource          |
      | User       | file-to-share.txt |
      | User       | SHARED-FOLDER     |
      | User Light | file-to-share.txt |
      | User Light | SHARED-FOLDER     |

  @issue-9530 @email
  Scenario: user gets an in-app and mail notifications when someone with comma in display name shares a file
    Given the administrator has assigned the role "Admin" to user "Brian" using the Graph API
    And the user "Brian" has created a new user with the following attributes:
      | userName    | David             |
      | displayName | David, Lopez      |
      | email       | david@example.com |
      | password    | 1234              |
    And user "David" has uploaded file with content "sample text" to "lorem.txt"
    When user "David" sends the following resource share invitation using the Graph API:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    Then the HTTP status code should be "200"
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                |
      | David, Lopez shared lorem.txt with you |
    And user "Brian" should have received the following email from user "David"
      """
      Hello Brian Murphy

      David, Lopez has shared "lorem.txt" with you.

      Click here to view it: %base_url%/files/shares/with-me
      """

  @email
  Scenario Outline: group members gets an in-app and mail notifications in their respective languages when someone shares resources with the group
    Given group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "Brian" has switched the system language to "es" using the Graph API
    And user "Carol" has switched the system language to "de" using the Graph API
    And user "Alice" has created folder "HelloWorld"
    And user "Alice" has uploaded file with content "hello world" to "text.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | group1     |
      | shareType       | group      |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "200"
    And user "Brian" should get a notification with subject "Recurso compartido" and message:
      | message                                   |
      | Alice Hansen compartió <resource> contigo |
    And user "Carol" should get a notification with subject "Neue Freigabe" and message:
      | message                                   |
      | Alice Hansen hat <resource> mit Ihnen geteilt |
    And user "Brian" should have received the following email from user "Alice"
      """
      Hola Brian Murphy

      %displayname% ha compartido "<resource>" contigo.

      Click aquí para verlo: %base_url%/files/shares/with-me
      """
    And user "Carol" should have received the following email from user "Alice"
      """
      Hallo Carol King

      %displayname% hat "<resource>" mit Ihnen geteilt.

      Zum Ansehen hier klicken: %base_url%/files/shares/with-me
      """
    Examples:
      | resource   |
      | HelloWorld |
      | text.txt   |

  @env-config @email
  Scenario Outline: group members gets an in-app and mail notifications in default language when someone shares a file with the group
    Given the config "OCIS_DEFAULT_LANGUAGE" has been set to "de" for "notifications" service
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "Alice" has created folder "HelloWorld"
    And user "Alice" has uploaded file with content "hello world" to "text.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | group1     |
      | shareType       | group      |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "200"
    And user "Brian" should get a notification with subject "Neue Freigabe" and message:
      | message                                     |
      | Alice Hansen hat <resource> mit Ihnen geteilt |
    And user "Carol" should get a notification with subject "Neue Freigabe" and message:
      | message                                     |
      | Alice Hansen hat <resource> mit Ihnen geteilt |
    And user "Brian" should have received the following email from user "Alice"
      """
      Hallo Brian Murphy

      %displayname% hat "<resource>" mit Ihnen geteilt.

      Zum Ansehen hier klicken: %base_url%/files/shares/with-me
      """
    And user "Carol" should have received the following email from user "Alice"
      """
      Hallo Carol King

      %displayname% hat "<resource>" mit Ihnen geteilt.

      Zum Ansehen hier klicken: %base_url%/files/shares/with-me
      """
    Examples:
      | resource   |
      | HelloWorld |
      | text.txt   |
