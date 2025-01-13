@email
Feature: Notification Settings
  As a user
  I want to manage my notification settings
  So that I do not get notified of unimportant events


  Scenario: disable email notification
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has uploaded file with content "some data" to "lorem.txt"
    When user "Brian" disables email notification using the settings API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "object",
            "required": ["identifier","value"],
            "properties": {
              "identifier":{
                "type": "object",
                "required": ["extension","bundle","setting"],
                "properties": {
                  "extension":{
                    "const": "ocis-accounts"
                  },
                  "bundle":{
                    "const": "profile"
                  },
                  "setting":{
                    "const": "disable-email-notifications"
                  }
                }
              },
              "value":{
                "type": "object",
                "required": [
                  "bundleId",
                  "settingId",
                  "accountUuid",
                  "resource",
                  "boolValue"
                ],
                "properties":{
                  "bundleId":{
                    "pattern":"%user_id_pattern%"
                  },
                  "settingId":{
                    "pattern":"%user_id_pattern%"
                  },
                  "accountUuid":{
                    "pattern":"%user_id_pattern%"
                  },
                  "resource":{
                    "type": "object",
                    "required":["type"],
                    "properties": {
                      "type":{
                        "const": "TYPE_USER"
                      }
                    }
                  },
                  "boolValue":{
                    "const": true
                  }
                }
              }
            }
          }
        }
      }
      """
    And user "Alice" has sent the following resource share invitation:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    And user "Brian" should have "0" emails