@ocm
Feature: an user shares resources usin ScienceMesh application
  As a user
  I want to share resources between different ocis instances

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And using server "REMOTE"
    And user "Brian" has been created with default attributes and without skeleton files

  @issue-9534
  Scenario Outline: local user shares resources to federation user
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    When user "Alice" sends the following resource share invitation to federated user using the Graph API:
      | resource        | <resource>                    |
      | space           | Personal                      |
      | sharee          | Brian                         |
      | shareType       | user                          |
      | permissionsRole | Viewer                        |
      | federatedServer | @federation-ocis-server:10200 |
    Then the HTTP status code should be "200"
    When using server "REMOTE"
    And user "Brian" lists the shares shared with him without retry using the Graph API
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
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "@UI.Hidden",
                "@client.synchronize",
                "createdBy",
                "name"
              ],
              "properties": {
                "@UI.Hidden": {
                  "type": "boolean",
                  "enum": [
                    false
                  ]
                },
                "@client.synchronize": {
                  "type": "boolean",
                  "enum": [
                    false
                  ]
                },
                "createdBy": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "const": "Alice Hansen"
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%federated_user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "name": {
                  "const": "<resource>"
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | resource      |
      | folderToShare |
      | textfile.txt  |

  @issue-9534
  Scenario Outline: federation user shares resource to local user after accepting invitation
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And user "Brian" has created folder "folderToShare"
    And user "Brian" has uploaded file with content "ocm test" to "/textfile.txt"
    When user "Brian" sends the following resource share invitation to federated user using the Graph API:
      | resource        | <resource>        |
      | space           | Personal          |
      | sharee          | Alice             |
      | shareType       | user              |
      | permissionsRole | Viewer            |
      | federatedServer | @ocis-server:9200 |
    Then the HTTP status code should be "200"
    When using server "LOCAL"
    And user "Alice" lists the shares shared with her without retry using the Graph API
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
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "@UI.Hidden",
                "@client.synchronize",
                "createdBy",
                "name"
              ],
              "properties": {
                "@UI.Hidden": {
                  "type": "boolean",
                  "enum": [
                    false
                  ]
                },
                "@client.synchronize": {
                  "type": "boolean",
                  "enum": [
                    false
                  ]
                },
                "createdBy": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "const": "Brian Murphy"
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%federated_user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "name": {
                  "const": "<resource>"
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | resource      |
      | folderToShare |
      | textfile.txt  |
