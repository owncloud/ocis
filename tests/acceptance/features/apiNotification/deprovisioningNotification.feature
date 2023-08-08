Feature: Deprovisioning notification
  As a user admin
  I want to inform users about shutting down and deprovisioning the instance
  So they can download and save their data in time

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |


  Scenario: the administrator user creates a deprovisioning notification about shutting down the instance
    When the administrator creates a deprovisioning notification
    Then the HTTP status code should be "200"
    When user "Alice" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "Instance will be shut down and deprovisioned" and the message-details should match
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
          "app": {
            "type": "string",
            "enum": [
              "userlog"
            ]
          },
          "message": {
            "type": "string",
            "enum": [
              "Attention! The instance will be shut down and deprovisioned on 2043-07-04T11:23:12Z. Download all your data before that date as no access past that date is possible."
            ]
          },
          "messageRich": {
            "type": "string",
            "enum": [
              "Attention! The instance will be shut down and deprovisioned on {date}. Download all your data before that date as no access past that date is possible."
            ]
          },
          "messageRichParameters": {
            "type": "object"
          },
          "notification_id": {
            "type": "string",
            "enum": [
              "deprovision"
            ]
          },
          "object_id": {
            "type": "string"
          },
          "object_type": {
            "type": "string",
            "enum": [
              "resource"
            ]
          },
          "subject": {
            "type": "string",
            "enum": [
              "Instance will be shut down and deprovisioned"
            ]
          },
          "subjectRich": {
            "type": "string",
            "enum": [
              "Instance will be shut down and deprovisioned"
            ]
          },
          "user": {
            "type": "string"
          }
        }
      }
      """


  Scenario Outline: non-admin user tries to create a deprovisioning notification
    Given the administrator has assigned the role "<role>" to user "Alice" using the Graph API
    When user "Alice" tries to create a deprovisioning notification
    Then the HTTP status code should be "404"
    And user "Alice" should not have any notification
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario Outline: get a deprovisioning notification in various languages
    Given the administrator has created a deprovisioning notification
    And user "Alice" has switched the system language to "<language>"
    When user "Alice" lists all notifications
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
      | language | subject                                                          | message                                                                                                                                                                                                 |
      | de       | Instanz wird heruntergefahren und außer Betrieb genommen werden. | Achtung! Diese Instanz wird am 2043-07-04T11:23:12Z heruntergefahren und außer Betrieb genommen werden. Laden Sie Ihre Daten vor diesem Tag herunter, da Sie danach nicht mehr darauf zugreifen können. |
      | es       | La instancia se cerrará y se desaprovisionará                    | ¡Atención! La instancia se cerrará y se desaprovisionará el 2043-07-04T11:23:12Z. Descarga todos tus datos antes de esa fecha, puesto que el acceso pasada la fecha no será posible.                    |


  Scenario: deprovisioning notification reappears again even after being marked as read
    Given the administrator has created a deprovisioning notification
    And user "Alice" has deleted all notifications
    When user "Alice" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "Instance will be shut down and deprovisioned" and the message-details should match
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
              "Attention! The instance will be shut down and deprovisioned on 2043-07-04T11:23:12Z. Download all your data before that date as no access past that date is possible."
            ]
          }
        }
      }
      """


  Scenario: the administrator user tries to delete the deprovisioning notification
    Given the administrator has created a deprovisioning notification
    When the administrator deletes the deprovisioning notification
    Then the HTTP status code should be "200"
    And user "Alice" should not have any notification

  
  Scenario Outline: non-admin user cannot delete the deprovisioning notification
    Given the administrator has assigned the role "<role>" to user "Alice" using the Graph API
    When user "Alice" tries to delete the deprovisioning notification
    Then the HTTP status code should be "404"
    And user "Alice" should not have any notification
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | User Light  |
