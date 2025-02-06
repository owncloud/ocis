@ocm @issue-9898
Feature: List a federated sharing permissions
  As a user
  I want to list the permissions for federated share
  So that the federated share is assigned the correct permissions

  Background:
    Given user "Alice" has been created with default attributes
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has accepted invitation


  Scenario: user lists permissions of a resource shared to a federated user
    Given using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    When user "Alice" gets permissions list for file "textfile.txt" of the space "Personal" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@libre.graph.permissions.actions.allowedValues",
          "@libre.graph.permissions.roles.allowedValues",
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "uniqueItems": true,
            "items": {
              "oneOf":[
                {
                  "type": "object",
                  "required": ["grantedToV2","id","roles"],
                  "properties": {
                    "grantedToV2": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["@libre.graph.userType","displayName","id"],
                          "properties": {
                            "@libre.graph.userType": {"const": "Federated"},
                            "id": {"pattern": "^%federated_user_id_pattern%$"},
                            "displayName": {"const": "Brian Murphy"}
                          }
                        }
                      }
                    },
                    "id": {"pattern": "^%user_id_pattern%$"},
                    "invitation": {
                      "type": "object",
                      "required": ["invitedBy"],
                      "properties": {
                        "invitedBy": {
                          "type": "object",
                          "required": ["user"],
                          "properties": {
                            "user": {
                              "type": "object",
                              "required": ["@libre.graph.userType", "displayName", "id"],
                              "properties": {
                                "@libre.graph.userType": {"const": "Member"},
                                "id": {"pattern": "^%user_id_pattern%$"},
                                "displayName": {"const": "Alice Hansen"}
                              }
                            }
                          }
                        }
                      }
                    },
                    "roles": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {"pattern": "^%role_id_pattern%$"}
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists permissions of a project resource shared to a federated user
    Given using server "LOCAL"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | projectSpace  |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Alice" gets permissions list for folder "folderToShare" of the space "projectSpace" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@libre.graph.permissions.actions.allowedValues",
          "@libre.graph.permissions.roles.allowedValues",
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "uniqueItems": true,
            "items": {
              "oneOf":[
                {
                  "type": "object",
                  "required": [
                    "grantedToV2",
                    "id",
                    "roles"
                  ],
                  "properties": {
                    "grantedToV2": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["@libre.graph.userType","displayName","id"],
                          "properties": {
                            "@libre.graph.userType": {"const": "Federated"},
                            "id": {"pattern": "^%federated_user_id_pattern%$"},
                            "displayName": {"const": "Brian Murphy"}
                          }
                        }
                      }
                    },
                    "id": { "pattern": "^%user_id_pattern%$" },
                    "invitation": {
                      "type": "object",
                      "required": ["invitedBy"],
                      "properties": {
                        "invitedBy": {
                          "type": "object",
                          "required": ["user"],
                          "properties": {
                            "user": {
                              "type": "object",
                              "required": ["@libre.graph.userType", "displayName", "id"],
                              "properties": {
                                "@libre.graph.userType": {"const": "Member"},
                                "id": {"pattern": "^%user_id_pattern%$"},
                                "displayName": {"const": "Alice Hansen"}
                              }
                            }
                          }
                        }
                      }
                    },
                    "roles": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {"pattern": "^%role_id_pattern%$"}
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists permissions of a project resource shared to a federated user
    Given using server "LOCAL"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | projectSpace |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    When user "Alice" gets permissions list for file "textfile.txt" of the space "projectSpace" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@libre.graph.permissions.actions.allowedValues",
          "@libre.graph.permissions.roles.allowedValues",
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "uniqueItems": true,
            "items": {
              "oneOf":[
                {
                  "type": "object",
                  "required": ["grantedToV2","id","roles"],
                  "properties": {
                    "grantedToV2": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["@libre.graph.userType","displayName","id"],
                          "properties": {
                            "@libre.graph.userType": {"const": "Federated"},
                            "id": {"pattern": "^%federated_user_id_pattern%$"},
                            "displayName": {"const": "Brian Murphy"}
                          }
                        }
                      }
                    },
                    "id": {"pattern": "^%user_id_pattern%$"},
                    "invitation": {
                      "type": "object",
                      "required": ["invitedBy"],
                      "properties": {
                        "invitedBy": {
                          "type": "object",
                          "required": ["user"],
                          "properties": {
                            "user": {
                              "type": "object",
                              "required": ["@libre.graph.userType", "displayName", "id"],
                              "properties": {
                                "@libre.graph.userType": {"const": "Member"},
                                "id": {"pattern": "^%user_id_pattern%$"},
                                "displayName": {"const": "Alice Hansen"}
                              }
                            }
                          }
                        }
                      }
                    },
                    "roles": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {"pattern": "^%role_id_pattern%$"}
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

