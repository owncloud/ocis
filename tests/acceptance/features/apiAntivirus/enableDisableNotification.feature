@antivirus @notification
Feature: enable/disable malware related notification
  As a system administrator and user
  I want to get notifications related to malware
  So that I can quickly detect and analyze security threats

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: disable in-app notification for "File rejected" event
    Given using <dav-path-version> DAV path
    When user "Brian" disables notification for the following event using the settings API:
      | event             | File Rejected |
      | notificationTypes | in-app        |
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
                  "extension":{ "const": "ocis-accounts" },
                  "bundle":{ "const": "profile" },
                  "setting":{ "const": "event-postprocessing-step-finished-options" }
                }
              },
              "value":{
                "type": "object",
                "required": ["id","bundleId","settingId","accountUuid","resource","collectionValue"],
                "properties":{
                  "id":{ "pattern":"%uuidv4_pattern%" },
                  "bundleId":{ "pattern":"%uuidv4_pattern%" },
                  "settingId":{ "pattern":"%uuidv4_pattern%" },
                  "accountUuid":{ "pattern":"%uuidv4_pattern%" },
                  "resource":{
                    "type": "object",
                    "required":["type"],
                    "properties": {
                      "type":{ "const": "TYPE_USER" }
                    }
                  },
                  "collectionValue":{
                    "type": "object",
                    "required":["values"],
                    "properties": {
                      "values":{
                        "type": "array",
                        "maxItems": 1,
                        "minItems": 1,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "in-app" },
                                "boolValue":{ "const": false }
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    And user "Brian" has uploaded file "filesForUpload/filesWithVirus/<file-name>" to "<new-file-name>"
    And user "Brian" should not have any notification
    Examples:
      | dav-path-version | file-name     | new-file-name  |
      | old              | eicar.com     | virusFile1.txt |
      | new              | eicar.com     | virusFile1.txt |
      | spaces           | eicar.com     | virusFile1.txt |


  Scenario: disable in-app notification for "File rejected" event (Project space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "newSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | newSpace     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Editor |
    When user "Brian" disables notification for the following event using the settings API:
      | event             | File Rejected |
      | notificationTypes | in-app        |
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
                  "extension":{ "const": "ocis-accounts" },
                  "bundle":{ "const": "profile" },
                  "setting":{ "const": "event-postprocessing-step-finished-options" }
                }
              },
              "value":{
                "type": "object",
                "required": ["id","bundleId","settingId","accountUuid","resource","collectionValue"],
                "properties":{
                  "id":{ "pattern":"%uuidv4_pattern%" },
                  "bundleId":{ "pattern":"%uuidv4_pattern%" },
                  "settingId":{ "pattern":"%uuidv4_pattern%" },
                  "accountUuid":{ "pattern":"%uuidv4_pattern%" },
                  "resource":{
                    "type": "object",
                    "required":["type"],
                    "properties": {
                      "type":{ "const": "TYPE_USER" }
                    }
                  },
                  "collectionValue":{
                    "type": "object",
                    "required":["values"],
                    "properties": {
                      "values":{
                        "type": "array",
                        "maxItems": 1,
                        "minItems": 1,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "in-app" },
                                "boolValue":{ "const": false }
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    And user "Brian" has uploaded a file "filesForUpload/filesWithVirus/eicar.com" to "virusFile.txt" in space "newSpace"
    And user "Brian" should get a notification with subject "Space shared" and message:
      | message                                  |
      | Alice Hansen added you to Space newSpace |
    But user "Brian" should not have a notification related to resource "virusFile.txt" with subject "Virus found"
