Feature: an user gets the resources shared to them
  As a user
  I want to get resources shared with me
  So that I can know about what resources I have access to

  https://owncloud.dev/libre-graph-api/#/me.drive/ListSharedWithMe

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  Scenario: user gets the resources shared with them
    Given user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Brian" lists the resources shared with him using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
{
  "type": "object",
  "required": [
    "value"
  ],
  "properties": {
    "value": {
      "type": "array",
      "items": {
        "type": "object",
        "required": [
          "name",
          "id",
          "remoteItem"
        ],
        "properties": {
          "id": {
            "type": "string",
            "pattern": "^%share_id_pattern%$"
          },
          "name": {
            "type": "string",
            "enum": [
              "textfile0.txt"
            ]
          },
          "remoteItem": {
            "type": "object",
            "required": [
              "createdBy",
              "file",
              "id",
              "name",
              "shared",
              "size"
            ],
            "properties": {
              "createdBy": {
                "type": "object",
                "required": [
                  "user"
                ],
                "properties": {
                  "user": {
                    "type": "object",
                    "required": [
                      "id",
                      "displayName"
                    ],
                    "properties": {
                      "id": {
                        "type": "string",
                        "pattern": "^%user_id_pattern%$"
                      },
                      "displayName": {
                        "type": "string",
                        "enum": [
                          "Alice Hansen"
                        ]
                      }
                    }
                  }
                }
              },
              "file": {
                "type": "object",
                "required": [
                  "mimeType"
                ],
                "properties": {
                  "mimeType": {
                    "type": "string",
                    "pattern": "text/plain"
                  }
                }
              },
              "id": {
                "type": "string",
                "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"

              },
              "name": {
                "type": "string",
                "enum": [
                  "textfile0.txt"
                ]
              },
              "shared": {
                "type": "object",
                "required": [
                  "sharedBy",
                  "owner"
                ],
                "properties": {
                  "owner": {
                    "type": "object",
                    "required": [
                      "user"
                    ],
                    "properties": {
                      "user": {
                        "type": "object",
                        "required": [
                          "id",
                          "displayName"
                        ],
                        "properties": {
                          "id": {
                            "type": "string",
                            "pattern": "^%user_id_pattern%$"
                          },
                          "displayName": {
                            "type": "string",
                            "enum": [
                              "Alice Hansen"
                            ]
                          }
                        }
                      }
                    }
                  },
                  "sharedBy": {
                    "type": "object",
                    "required": [
                      "user"
                    ],
                    "properties": {
                      "user": {
                        "type": "object",
                        "required": [
                          "id",
                          "displayName"
                        ],
                        "properties": {
                          "id": {
                            "type": "string",
                            "pattern": "^%user_id_pattern%$"
                          },
                          "displayName": {
                            "type": "string",
                            "enum": [
                              "Alice Hansen"
                            ]
                          }
                        }
                      }
                    }
                  }
                }
              },
              "size": {
                "type": "number",
                "enum": [
                  11
                ]
              }
            }
          }
        }
      }
    }
  }
}
    """
