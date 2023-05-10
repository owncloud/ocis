@api
Feature: Notification
  As a user
  I want to be notified of actions related to me
  So that I can stay updated about the information

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And the administrator has given "Alice" the role "Space Admin" using the settings api


  Scenario: user gets a notification of space sharing
    Given user "Alice" has created a space "notification checking" with the default quota using the GraphApi
    And user "Alice" has shared a space "notification checking" with settings:
      | shareWith | Brian  |
      | role      | editor |
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
          "app": {
            "type": "string",
            "enum": [
              "userlog"
            ]
          },
          "message": {
            "type": "string",
            "enum": [
              "Alice Hansen added you to Space notification checking"
            ]
          },
          "messageRich": {
            "type": "string",
            "enum": [
              "{user} added you to Space {space}"
            ]
          },
          "messageRichParameters": {
            "type": "object",
            "required": [
              "space",
              "user"
            ],
            "properties": {
              "space": {
                "type": "object",
                "required": [
                  "id",
                  "name"
                ],
                "properties": {
                  "id": {
                    "type": "string",
                    "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
                  },
                  "name": {
                    "type": "string",
                    "enum": [
                      "notification checking"
                    ]
                  }
                }
              },
              "user": {
                "type": "object",
                "required": [
                  "displayname",
                  "id",
                  "name"
                ],
                "properties": {
                  "displayname": {
                    "type": "string",
                    "enum": [
                      "Alice Hansen"
                    ]
                  },
                  "id": {
                    "type": "string",
                    "enim": [
                      "%user_id%"
                    ]
                  },
                  "name": {
                    "type": "string",
                    "enum": [
                      "Alice"
                    ]
                  }
                }
              }
            }
          },
          "notification_id": {
            "type": "string"
          },
          "object_id": {
            "type": "string"
          },
          "object_type": {
            "type": "string",
            "enum": [
              "storagespace"
            ]
          },
          "subject": {
            "type": "string",
            "enum": [
              "Space shared"
            ]
          },
          "subjectRich": {
            "type": "string",
            "enum": [
              "Space shared"
            ]
          },
          "user": {
            "type": "string",
            "enum": [
              "Alice"
            ]
          }
        }
      }
      """


  Scenario Outline: user gets a notification of space sharing in different languages
    Given user "Alice" has created a space "notification checking" with the default quota using the GraphApi
    And user "Brian" has switched the system language to "<language>"
    And user "Alice" has shared a space "notification checking" with settings:
      | shareWith | Brian  |
      | role      | editor |
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
      | language | subject           | message                                                          |
      | de       | Space freigegeben | Alice Hansen hat Sie zum Space notification checking hinzugef\u00fcgt |
      | en       | Space shared      | Alice Hansen added you to Space notification checking            |
      | es       | Space compartido  | Alice Hansen te añadió al Space notification checking            |
