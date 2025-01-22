@email
Feature: Notification Settings
  As a user
  I want to manage my notification settings
  So that I do not get notified of unimportant events

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has uploaded file with content "some data" to "lorem.txt"


  Scenario: disable email notification
    Given user "Alice" has uploaded file with content "some data" to "lorem.txt"
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
                  "id",
                  "bundleId",
                  "settingId",
                  "accountUuid",
                  "resource"
                ],
                "properties":{
                  "id":{
                    "pattern": "%user_id_pattern%"
                  },
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


  Scenario: disable mail and in-app notification for Share Received event
    When user "Brian" disables notification for the following events using the settings API:
      | Share Received | mail,in-app |
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
                    "const": "event-share-created-options"
                  }
                }
              },
              "value":{
                "type": "object",
                "required": [
                  "id",
                  "bundleId",
                  "settingId",
                  "accountUuid",
                  "resource",
                  "collectionValue"
                ],
                "properties":{
                  "id":{
                    "pattern":"%user_id_pattern%"
                  },
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
                  "collectionValue":{
                    "type": "object",
                    "required":["values"],
                    "properties": {
                      "values":{
                        "type": "array",
                        "maxItems": 2,
                        "minItems": 2,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": [
                                "key",
                                "boolValue"
                              ],
                              "properties": {
                                "key":{
                                  "const": "mail"
                                },
                                "boolValue":{
                                  "const": false
                                }
                              }
                            },
                            {
                              "type": "object",
                              "required": [
                                "key",
                                "boolValue"
                              ],
                              "properties": {
                                "key":{
                                  "const": "in-app"
                                },
                                "boolValue":{
                                  "const": false
                                }
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
    And user "Alice" has sent the following resource share invitation:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    And user "Brian" should have "0" emails
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the notifications should be empty


  Scenario: disable in-app notification for Share Removed event
    Given user "Alice" has sent the following resource share invitation:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    When user "Brian" disables notification for the following events using the settings API:
      | Share Removed | in-app |
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
                    "const": "event-share-removed-options"
                  }
                }
              },
              "value":{
                "type": "object",
                "required": [
                  "id",
                  "bundleId",
                  "settingId",
                  "accountUuid",
                  "resource",
                  "collectionValue"
                ],
                "properties":{
                  "id":{
                    "pattern":"%user_id_pattern%"
                  },
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
                              "required": [
                                "key",
                                "boolValue"
                              ],
                              "properties": {
                                "key":{
                                  "const": "in-app"
                                },
                                "boolValue":{
                                  "const": false
                                }
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
    And user "Alice" has removed the access of user "Brian" from resource "lorem.txt" of space "Personal"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                |
      | Alice Hansen shared lorem.txt with you |
    But user "Brian" should not have a notification related to resource "lorem.txt" with subject "Resource unshared"


  Scenario: disable in-app notification for Share Removed event (Project space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "newSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "newSpace" with content "some content" to "insideSpace.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | insideSpace.txt |
      | space           | newSpace        |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | Viewer          |
    When user "Brian" disables notification for the following events using the settings API:
      | Share Removed | in-app |
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
                    "const": "event-share-removed-options"
                  }
                }
              },
              "value":{
                "type": "object",
                "required": [
                  "id",
                  "bundleId",
                  "settingId",
                  "accountUuid",
                  "resource",
                  "collectionValue"
                ],
                "properties":{
                  "id":{
                    "pattern":"%user_id_pattern%"
                  },
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
                              "required": [
                                "key",
                                "boolValue"
                              ],
                              "properties": {
                                "key":{
                                  "const": "in-app"
                                },
                                "boolValue":{
                                  "const": false
                                }
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
    And user "Alice" has removed the access of user "Brian" from resource "insideSpace.txt" of space "newSpace"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                      |
      | Alice Hansen shared insideSpace.txt with you |
    But user "Brian" should not have a notification related to resource "insideSpace.txt" with subject "Resource unshared"
