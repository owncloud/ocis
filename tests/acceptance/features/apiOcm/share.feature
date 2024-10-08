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

  @issue-10051
  Scenario Outline: try to add federated user as a member of a project space (permissions endpoint)
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Brian" using the Graph API
    And user "Brian" has created a space "brian's space" with the default quota using the Graph API
    When user "Brian" tries to send the following space share invitation to federated user using permissions endpoint of the Graph API:
      | space           | brian's space      |
      | sharee          | Alice              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | federatedServer | @ocis-server:9200  |
    Then the HTTP status code should be "403"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "code",
          "message"
        ],
        "properties": {
          "code": {
            "const": "PERMISSION_DENIED"
          },
          "message": {
            "const": "permission denied to create the file"
          }
        }
      }
      """
    And using server "LOCAL"
    And the user "Alice" should not have a space called "brian's space"
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: try to add federated user as a member of a project space (root endpoint)
    Given using server "LOCAL"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "alice's space" with the default quota using the Graph API
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    When user "Alice" tries to send the following space share invitation to federated user using root endpoint of the Graph API:
      | space           | alice's space                 |
      | sharee          | Brian                         |
      | shareType       | user                          |
      | permissionsRole | <permissions-role>            |
      | federatedServer | @federation-ocis-server:10200 |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "innererror",
              "message"
            ],
            "properties": {
              "code": {
                "const": "invalidRequest"
              },
              "innererror": {
                "type": "object",
                "required": [
                  "date",
                  "request-id"
                ]
              },
              "message": {
                "const": "federated user can not become a space member"
              }
            }
          }
        }
      }
      """
    And using server "REMOTE"
    And the user "Brian" should not have a space called "alice's space"
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |
