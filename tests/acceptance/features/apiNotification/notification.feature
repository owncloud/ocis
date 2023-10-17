Feature: Notification
  As a user
  I want to be notified of various events
  So that I can stay updated about the information

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And user "Alice" has uploaded file with content "other data" to "/textfile1.txt"
    And user "Alice" has created folder "my_data"


  Scenario Outline: user gets a notification of resource sharing
    Given user "Alice" has shared entry "<resource>" with user "Brian"
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
        "app": {
          "type": "string",
          "enum": ["userlog"]
        },
        "message": {
          "type": "string",
          "enum": ["Alice Hansen shared <resource> with you"]
        },
        "messageRich": {
          "type": "string",
          "enum": ["{user} shared {resource} with you"]
        },
        "messageRichParameters": {
          "type": "object",
          "required": [
            "resource",
            "user"
          ],
          "properties": {
            "resource": {
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
                  "enum": ["<resource>"]
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
                  "enum": ["Alice Hansen"]
                },
                "id": {
                  "type": "string",
                  "enim": ["%user_id%"]
                },
                "name": {
                  "type": "string",
                  "enum": ["Alice"]
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
          "enum": ["share"]
        },
        "subject": {
          "type": "string",
          "enum": ["Resource shared"]
        },
        "subjectRich": {
          "type": "string",
          "enum": ["Resource shared"]
        },
        "user": {
          "type": "string",
          "enum": ["Alice"]
        }
      }
    }
    """
    Examples:
      | resource      |
      | textfile1.txt |
      | my_data       |


  Scenario Outline: user gets a notification of unsharing resource
    Given user "Alice" has shared entry "<resource>" with user "Brian"
    And user "Brian" has accepted share "/<resource>" offered by user "Alice"
    And user "Alice" has unshared entity "<resource>" shared to "Brian"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "Resource unshared" and the message-details should match
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
          "enum": ["userlog"]
        },
        "message": {
          "type": "string",
          "enum": ["Alice Hansen unshared <resource> with you"]
        },
        "messageRich": {
          "type": "string",
          "enum": ["{user} unshared {resource} with you"]
        },
        "messageRichParameters": {
          "type": "object",
          "required": [
            "resource",
            "user"
          ],
          "properties": {
            "resource": {
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
                  "enum": ["<resource>"]
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
                  "enum": ["Alice Hansen"]
                },
                "id": {
                  "type": "string",
                  "enim": ["%user_id%"]
                },
                "name": {
                  "type": "string",
                  "enum": ["Alice"]
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
          "enum": ["share"]
        },
        "subject": {
          "type": "string",
          "enum": ["Resource unshared"]
        },
        "subjectRich": {
          "type": "string",
          "enum": ["Resource unshared"]
        },
        "user": {
          "type": "string",
          "enum": ["Alice"]
        }
      }
    }
    """
    Examples:
      | resource      |
      | textfile1.txt |
      | my_data       |


  Scenario Outline: get a notification about a file share in various languages
    Given user "Brian" has switched the system language to "<language>"
    And user "Alice" has shared entry "textfile1.txt" with user "Brian" with permissions "17"
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
      | language | subject            | message                                          |
      | de       | Neue Freigabe      | Alice Hansen hat textfile1.txt mit Ihnen geteilt |
      | es       | Recurso compartido | Alice Hansen compartió textfile1.txt contigo     |

  @env-config
  Scenario: get a notification about a file share in default languages
    Given the config "OCIS_DEFAULT_LANGUAGE" has been set to "de"
    And user "Alice" has shared entry "textfile1.txt" with user "Brian" with permissions "17"
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


  Scenario Outline: notifications related to a resource get deleted when the resource is deleted
    Given user "Alice" has shared entry "<resource>" with user "Brian"
    And user "Brian" has accepted share "/<resource>" offered by user "Alice"
    And user "Alice" has unshared entity "<resource>" shared to "Brian"
    And user "Alice" has deleted entity "/<resource>"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the notifications should be empty
    Examples:
      | resource      |
      | textfile1.txt |
      | my_data       |
