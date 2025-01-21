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
                  "extension":{ "const": "ocis-accounts" },
                  "bundle":{ "const": "profile" },
                  "setting":{ "const": "disable-email-notifications" }
                }
              },
              "value":{
                "type": "object",
                "required": ["id","bundleId","settingId","accountUuid","resource"],
                "properties":{
                  "id":{ "pattern": "%uuidv4_pattern%" },
                  "bundleId":{ "pattern":"%uuidv4_pattern%" },
                  "settingId":{ "pattern":"%uuidv4_pattern%" },
                  "accountUuid":{ "pattern":"%uuidv4_pattern%" },
                  "resource":{
                    "type": "object",
                    "required":["type"],
                    "properties": {
                      "type":{ "const": "TYPE_USER" }
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


  Scenario: disable mail and in-app notification for "Share Received" event
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
                  "extension":{ "const": "ocis-accounts" },
                  "bundle":{ "const": "profile" },
                  "setting":{ "const": "event-share-created-options" }
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
                        "maxItems": 2,
                        "minItems": 2,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "mail" },
                                "boolValue":{ "const": false }
                              }
                            },
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


  Scenario: disable mail and in-app notification for "Share Removed" event
    Given user "Alice" has sent the following resource share invitation:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    When user "Brian" disables notification for the following events using the settings API:
      | Share Removed | mail, in-app |
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
                  "setting":{ "const": "event-share-removed-options" }
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
                        "maxItems": 2,
                        "minItems": 2,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "mail" },
                                "boolValue":{ "const": false }
                              }
                            },
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
    And user "Alice" has removed the access of user "Brian" from resource "lorem.txt" of space "Personal"
    And user "Brian" should have "1" emails
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                |
      | Alice Hansen shared lorem.txt with you |
    But user "Brian" should not have a notification related to resource "lorem.txt" with subject "Resource unshared"


  Scenario: disable mail and in-app notification for "Share Removed" event (Project space)
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
      | Share Removed | mail, in-app |
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
                  "setting":{ "const": "event-share-removed-options" }
                }
              },
              "value":{
                "type": "object",
                "required": ["id","bundleId","settingId","accountUuid","resource","collectionValue"],
                "properties":{
                  "id":{ "pattern":"%user_id_pattern%" },
                  "bundleId":{ "pattern":"%user_id_pattern%" },
                  "settingId":{ "pattern":"%user_id_pattern%" },
                  "accountUuid":{ "pattern":"%user_id_pattern%" },
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
                        "maxItems": 2,
                        "minItems": 2,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "mail" },
                                "boolValue":{ "const": false }
                              }
                            },
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
    And user "Alice" has removed the access of user "Brian" from resource "insideSpace.txt" of space "newSpace"
    And user "Brian" should have "1" emails
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                      |
      | Alice Hansen shared insideSpace.txt with you |
    But user "Brian" should not have a notification related to resource "insideSpace.txt" with subject "Resource unshared"

  @antivirus
  Scenario Outline: disable in-app notification for "File rejected" event
    Given using <dav-path-version> DAV path
    When user "Brian" disables notification for the following events using the settings API:
      | File rejected | in-app |
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

  @antivirus
  Scenario: disable in-app notification for "File rejected" event (Project space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "newSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | newSpace     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Editor |
    When user "Brian" disables notification for the following events using the settings API:
      | File rejected | in-app |
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


  Scenario: disable in-app notification for "Space disabled" event
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | new-space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    When user "Brian" disables notification for the following events using the settings API:
      | Space disabled | in-app |
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
                  "setting":{
                    "const": "event-space-disabled-options"
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
                              "required": [
                                "key",
                                "boolValue"
                              ],
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
    And user "Alice" has disabled a space "new-space"
    And user "Brian" should get a notification with subject "Space shared" and message:
      | message                                   |
      | Alice Hansen added you to Space new-space |
    But user "Brian" should not have a notification related to space "new-space" with subject "Space disabled"

  @issue-10864
  Scenario: disable email notification for user light
    Given the administrator has assigned the role "User Light" to user "Brian" using the Graph API
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
                  "extension":{ "const": "ocis-accounts" },
                  "bundle":{ "const": "profile" },
                  "setting":{ "const": "disable-email-notifications" }
                }
              },
              "value":{
                "type": "object",
                "required": ["id","bundleId","settingId","accountUuid","resource"],
                "properties":{
                  "id":{ "pattern": "%uuidv4_pattern%" },
                  "bundleId":{ "pattern":"%uuidv4_pattern%" },
                  "settingId":{ "pattern":"%uuidv4_pattern%" },
                  "accountUuid":{ "pattern":"%uuidv4_pattern%" },
                  "resource":{
                    "type": "object",
                    "required":["type"],
                    "properties": {
                      "type":{ "const": "TYPE_USER" }
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
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                |
      | Alice Hansen shared lorem.txt with you |
    And user "Brian" should have "0" emails


  Scenario: disable mail and in-app notification for "Added as space member" event
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "newSpace" with the default quota using the Graph API
    When user "Brian" disables notification for the following events using the settings API:
      | Added as space member | mail,in-app |
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
                  "setting":{ "const": "event-space-shared-options" }
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
                        "maxItems": 2,
                        "minItems": 2,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "mail" },
                                "boolValue":{ "const": false }
                              }
                            },
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
    And user "Alice" has sent the following space share invitation:
      | space           | newSpace     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the notifications should be empty
    And user "Brian" should have "0" emails


  Scenario: no email should be received when email sending interval is set to never
    When user "Brian" sets the email sending interval to "never" using the settings API
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
                  "setting":{ "const": "email-sending-interval-options" }
                }
              },
              "value":{
                "type": "object",
                "required": ["id","bundleId","settingId","accountUuid","resource","stringValue"],
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
                  "stringValue":{ "const":"never" }
                }
              }
            }
          }
        }
      }
      """
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | lorem.txt |
      | space           | Personal  |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    Then user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                |
      | Alice Hansen shared lorem.txt with you |
    And user "Brian" should have "0" emails


  Scenario: disable mail and in-app notification for "Removed as space member" event
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "newSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | newSpace     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    When user "Brian" disables notification for the following events using the settings API:
      | Removed as space member | mail,in-app |
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
                  "setting":{ "const": "event-space-unshared-options" }
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
                        "maxItems": 2,
                        "minItems": 2,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "mail" },
                                "boolValue":{ "const": false }
                              }
                            },
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
    And user "Alice" has removed the access of user "Brian" from space "newSpace"
    And user "Brian" should get a notification with subject "Space shared" and message:
      | message                                  |
      | Alice Hansen added you to Space newSpace |
    But user "Brian" should not have a notification related to space "newSpace" with subject "Removed from Space"
    And user "Brian" should have "1" emails
    And user "Brian" should have received the following email from user "Alice" about the share of project space "newSpace"
      """
      Hello Brian Murphy,

      %displayname% has invited you to join "newSpace".

      Click here to view it: %base_url%/f/%space_id%
      """
