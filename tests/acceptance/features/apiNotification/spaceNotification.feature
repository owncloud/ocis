@notification
Feature: Notification
  As a user
  I want to be notified of actions related to space
  So that I can stay updated about the spaces

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "notification checking" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | notification checking |
      | sharee          | Brian                 |
      | shareType       | user                  |
      | permissionsRole | Space Editor          |

  @issue-10937 @email
  Scenario Outline: user gets in-app and mail notifications of space shared
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "Space shared" and the message-details should match
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
          "message": {"const": "Alice Hansen added you to Space notification checking"},
          "messageRich": {"const": "{user} added you to Space {space}"},
          "messageRichParameters": {
            "type": "object",
            "required": ["space","user"],
            "properties": {
              "space": {
                "type": "object",
                "required": ["id","name"],
                "properties": {
                  "id": {"pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"},
                  "name": {"const": "notification checking"}
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
          "object_id": {"pattern": "^%space_id_pattern%$"},
          "object_type": {"const": "storagespace"},
          "subject": {"const": "Space shared"},
          "subjectRich": {"const": "Space shared"},
          "user": {"const": "Alice"}
        }
      }
      """
    And user "Brian" should have received the following email from user "Alice" about the share of project space "notification checking"
      """
      Hello Brian Murphy,

      %displayname% has invited you to join "notification checking".

      Click here to view it: %base_url%/f/%space_id%
      """
    Examples:
      | user-role  |
      | User       |
      | User Light |

  @issue-10937 @email
  Scenario Outline: gets in-app and mail notification of space unshared
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When user "Alice" removes the access of user "Brian" from space "notification checking" using root endpoint of the Graph API
    Then the HTTP status code should be "204"
    And user "Brian" should get a notification with subject "Removed from Space" and message:
      | message                                       |
      | Alice Hansen removed you from Space notification checking |
    And user "Brian" should have received the following email from user "Alice" about the share of project space "notification checking"
      """
      Hello Brian Murphy,

      %displayname% has removed you from "notification checking".

      You might still have access through your other groups or direct membership.

      Click here to check it: %base_url%/f/%space_id%
      """
    Examples:
      | user-role  |
      | User       |
      | User Light |

  @issue-10937
  Scenario Outline: user gets in-app notification of "Space disabled" event (note: no mail notification)
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Alice" has disabled a space "notification checking"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And there should be "2" notifications
    And the JSON response should contain a notification message with the subject "Space disabled" and the message-details should match
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
          "message": {"const": "Alice Hansen disabled Space notification checking"},
          "messageRich": {"const": "{user} disabled Space {space}"},
          "messageRichParameters": {
            "type": "object",
            "required": ["space","user"],
            "properties": {
              "space": {
                "type": "object",
                "required": ["id","name"],
                "properties": {
                  "id": {"pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"},
                  "name": {"const": "notification checking"}
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
          "object_id": {"pattern": "^%space_id_pattern%$"},
          "object_type": {"const": "storagespace"},
          "subject": {"const": "Space disabled"},
          "subjectRich": {"const": "Space disabled"},
          "user": {"const": "Alice"}
        }
      }
      """
    Examples:
      | user-role  |
      | User       |
      | User Light |


  Scenario Outline: get in-app notification about a space share in various languages
    Given user "Brian" has switched the system language to "<language>" using the Graph API
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "<subject>" and the message-details should match
      """
      {
        "type": "object",
        "required": ["message"],
        "properties": {
          "message": {"const": "<message>"}
        }
      }
      """
    Examples:
      | language | subject           | message                                                         |
      | de       | Space freigegeben | Alice Hansen hat Sie zu Space notification checking hinzugefügt |
      | es       | Space compartido  | Alice Hansen te añadió al Space notification checking           |


  Scenario: in-app notifications related to space get deleted when the sharer deletes that resource
    Given user "Alice" has removed the access of user "Brian" from space "notification checking"
    And user "Alice" has disabled a space "notification checking"
    And user "Alice" has deleted a space "notification checking"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the notifications should be empty

  @issue-10937 @email
  Scenario Outline: user doesn't get any notification after being removed from space (note: no mail notification)
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And user "Alice" has removed the access of user "Brian" from space "notification checking"
    And user "Alice" has disabled a space "notification checking"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And there should be "2" notifications
    And user "Brian" should have "2" emails
    But user "Brian" should not have a notification related to space "notification checking" with subject "Space disabled"
    Examples:
      | user-role  |
      | User       |
      | User Light |

  @email
  Scenario: group members get in-app and email notifications when someone shares a project space with the group
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "David" has been created with default attributes
    And group "group1" has been created
    And user "David" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | NewSpace |
      | sharee          | group1   |
      | shareType       | group    |
      | permissionsRole | Manager  |
    Then the HTTP status code should be "200"
    And user "David" should get a notification with subject "Space shared" and message:
      | message                                  |
      | Alice Hansen added you to Space NewSpace |
    And user "Carol" should get a notification with subject "Space shared" and message:
      | message                                  |
      | Alice Hansen added you to Space NewSpace |
    And user "David" should have received the following email from user "Alice" about the share of project space "NewSpace"
      """
      Hello David Lopez,

      %displayname% has invited you to join "NewSpace".

      Click here to view it: %base_url%/f/%space_id%
      """
    And user "Carol" should have received the following email from user "Alice" about the share of project space "NewSpace"
      """
      Hello Carol King,

      %displayname% has invited you to join "NewSpace".

      Click here to view it: %base_url%/f/%space_id%
      """

  @email
  Scenario: group members get in-app and email notifications in their respective languages when someone shares a space with the group
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "David" has been created with default attributes
    And group "group1" has been created
    And user "David" has been added to group "group1"
    And user "Carol" has been added to group "group1"
    And user "David" has switched the system language to "es" using the Graph API
    And user "Carol" has switched the system language to "de" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    When user "Alice" shares a space "new-space" with settings:
      | shareWith | group1 |
      | role      | viewer |
    Then the HTTP status code should be "200"
    And user "David" should get a notification with subject "Space compartido" and message:
      | message                                   |
      | Alice Hansen te añadió al Space new-space|
    And user "Carol" should get a notification with subject "Space freigegeben" and message:
      | message                                             |
      | Alice Hansen hat Sie zu Space new-space hinzugefügt |
    And user "David" should have received the following email from user "Alice" about the share of project space "new-space"
      """
      Hola David Lopez,

      Alice Hansen te ha invitado a unirte a "new-space".

      Click aquí para verlo: %base_url%/f/%space_id%
      """
    And user "Carol" should have received the following email from user "Alice" about the share of project space "new-space"
      """
      Hallo Carol King,

      Alice Hansen hat Sie eingeladen, dem Space "new-space" beizutreten.

      Zum Ansehen hier klicken: %base_url%/f/%space_id%
      """

  @issue-10937 @issue-10882 @email
  Scenario Outline: user gets in-app and email notifications when space membership expires
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When user "Alice" has expired the membership of user "Brian" from space "notification checking"
    Then the HTTP status code should be "200"
    And user "Brian" should get a notification with subject "Membership expired" and message:
      | message                                    |
      | Access to Space notification checking lost |
    And user "Brian" should have received the following email from user "Alice" about the share of project space "notification checking"
      """
      Hello Brian Murphy,

      Your membership of space notification checking has expired at %expiry_date_in_mail%

      Even though this membership has expired you still might have access through other shares and/or space memberships
      """
    Examples:
      | user-role  |
      | User       |
      | User Light |
