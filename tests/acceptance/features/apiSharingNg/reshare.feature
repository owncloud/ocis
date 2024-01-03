Feature: Reshare a share invitation
  As a user
  I want to be able to reshare the share invitations to other users
  So that they can have access to it

  https://owncloud.dev/libre-graph-api/#/drives.permissions/Invite

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |


  Scenario Outline: reshare a file to a user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType | file          |
      | resource     | textfile1.txt |
      | space        | Personal      |
      | sharee       | Brian         |
      | shareType    | user          |
      | role         | <role>        |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType | file           |
      | resource     | textfile1.txt  |
      | space        | Shares         |
      | sharee       | Carol          |
      | shareType    | user           |
      | role         | <reshare-role> |
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
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
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
                            "Carol King"
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
    Examples:
      | role        | reshare-role |
      | Viewer      | Viewer       |
      | File Editor | Viewer       |
      | File Editor | File Editor  |
      | Co Owner    | Viewer       |
      | Co Owner    | File Editor  |
      | Co Owner    | Co Owner     |
      | Manager     | Viewer       |
      | Manager     | File Editor  |
      | Manager     | Co Owner     |
      | Manager     | Manager      |


  Scenario Outline: reshare a folder to a user with different roles
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following share invitation:
      | resourceType | folder        |
      | resource     | FolderToShare |
      | space        | Personal      |
      | sharee       | Brian         |
      | shareType    | user          |
      | role         | <role>        |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType | folder         |
      | resource     | FolderToShare  |
      | space        | Shares         |
      | sharee       | Carol          |
      | shareType    | user           |
      | role         | <reshare-role> |
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
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
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
                            "Carol King"
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
    Examples:
      | role     | reshare-role |
      | Viewer   | Viewer       |
      | Editor   | Viewer       |
      | Editor   | Editor       |
      | Editor   | Uploader     |
      | Co Owner | Viewer       |
      | Co Owner | Editor       |
      | Co Owner | Co Owner     |
      | Co Owner | Uploader     |
      | Manager  | Viewer       |
      | Manager  | Editor       |
      | Manager  | Co Owner     |
      | Manager  | Uploader     |
      | Manager  | Manager      |


  Scenario Outline: reshare a file inside project space to a user with different roles
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "to share" to "textfile1.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType | file          |
      | resource     | textfile1.txt |
      | space        | NewSpace      |
      | sharee       | Brian         |
      | shareType    | user          |
      | role         | <role>        |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType | file           |
      | resource     | textfile1.txt  |
      | space        | Shares         |
      | sharee       | Carol          |
      | shareType    | user           |
      | role         | <reshare-role> |
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
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
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
                            "Carol King"
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
    Examples:
      | role        | reshare-role |
      | Viewer      | Viewer       |
      | File Editor | Viewer       |
      | File Editor | File Editor  |
      | Co Owner    | Viewer       |
      | Co Owner    | File Editor  |
      | Co Owner    | Co Owner     |
      | Manager     | Viewer       |
      | Manager     | File Editor  |
      | Manager     | Co Owner     |
      | Manager     | Manager      |


  Scenario Outline: reshare a folder inside project space to a user with different roles
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has sent the following share invitation:
      | resourceType | folder        |
      | resource     | FolderToShare |
      | space        | NewSpace      |
      | sharee       | Brian         |
      | shareType    | user          |
      | role         | <role>        |
    When user "Brian" sends the following share invitation using the Graph API:
      | resourceType | folder         |
      | resource     | FolderToShare  |
      | space        | Shares         |
      | sharee       | Carol          |
      | shareType    | user           |
      | role         | <reshare-role> |
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
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
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
                            "Carol King"
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
    Examples:
      | role     | reshare-role |
      | Viewer   | Viewer       |
      | Editor   | Viewer       |
      | Editor   | Editor       |
      | Editor   | Uploader     |
      | Co Owner | Viewer       |
      | Co Owner | Editor       |
      | Co Owner | Co Owner     |
      | Co Owner | Uploader     |
      | Manager  | Viewer       |
      | Manager  | Editor       |
      | Manager  | Co Owner     |
      | Manager  | Uploader     |
      | Manager  | Manager      |
