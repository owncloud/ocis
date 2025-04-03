@env-config
Feature: ListGrants role
  As a user
  I want to share resources with listGrants role
  So that sharee can view activities and grants list of shared resources

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: user shares personal space file with ListGrants role
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": ["createdDateTime","id","roles","grantedToV2"],
              "properties": {
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id","displayName"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
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
      | permissions-role            |
      | Viewer With ListGrants      |
      | File Editor With ListGrants |


  Scenario Outline: user shares personal space folder with ListGrants role
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | FolderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": ["createdDateTime","id","roles","grantedToV2"],
              "properties": {
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id","displayName"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
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
      | permissions-role       |
      | Viewer With ListGrants |
      | Editor With ListGrants |


  Scenario Outline: user shares project space file with ListGrants role
    Given the administrator has enabled the permissions role "<permissions-role>"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "textfile1.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile1.txt      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": ["createdDateTime","id","roles","grantedToV2"],
              "properties": {
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id","displayName"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
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
      | permissions-role            |
      | Viewer With ListGrants      |
      | File Editor With ListGrants |


  Scenario Outline: user shares project space folder with ListGrants role
    Given the administrator has enabled the permissions role "<permissions-role>"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | FolderToShare      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": ["createdDateTime","id","roles","grantedToV2"],
              "properties": {
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id","displayName"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
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
      | permissions-role       |
      | Viewer With ListGrants |
      | Editor With ListGrants |
