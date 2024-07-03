@ocm
Feature: an user shares resources usin ScienceMesh application
  As a user
  I want to share resources between different ocis instances

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And using server "REMOTE"
    And user "Brian" has been created with default attributes and without skeleton files


  Scenario: users shares folder to federation users after receiver accepted invitation
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    Then the HTTP status code should be "200"
    When using server "REMOTE"
    And user "Brian" lists the shares shared with him using the Graph API
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
                "@UI.Hidden":{
                  "type": "boolean",
                  "enum": [false]
                },
                "@client.synchronize":{
                  "type": "boolean",
                  "enum": [true]
                },
                "createdBy": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id"],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "const": "Alice Hansen"
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "name": {
                  "const": "folderToShare"
                }
              }
            }
          }
        }
      }
      """

  Scenario: users shares folder to federation users after accepting invitation
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And user "Brian" has created folder "folderToShare"
    When user "Brian" sends the following resource share invitation using the Graph API:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Alice         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    Then the HTTP status code should be "200"
    When using server "LOCAL"
    And user "Alice" lists the shares shared with her using the Graph API
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
                "@UI.Hidden":{
                  "type": "boolean",
                  "enum": [false]
                },
                "@client.synchronize":{
                  "type": "boolean",
                  "enum": [true]
                },
                "createdBy": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id"],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "const": "Brian Murphy"
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "name": {
                  "const": "folderToShare"
                }
              }
            }
          }
        }
      }
      """