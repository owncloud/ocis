Feature: settings api
  As a user,
  I want to use settings api
  So that I can manage user specific settings

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario: disable auto-sync share
    Given user "Alice" has uploaded file with content "lorem epsum" to "textfile.txt"
    When user "Brian" disables the auto-sync share using the settings API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value" : {
            "type": "object",
            "required": ["identifier","value"],
            "properties": {
              "identifier": {
                "type": "object",
                "required": [ "extension", "bundle", "setting"],
                "properties": {
                  "extension": { "const": "ocis-accounts" },
                  "bundle": { "const": "profile" },
                  "setting": { "const": "auto-accept-shares" }
                }
              },
              "value": {
                "type": "object",
                "required": [ "id", "bundleId", "settingId", "accountUuid", "resource", "boolValue" ],
                "properties": {
                  "id": { "pattern": "^%user_id_pattern%$" },
                  "bundleId": { "pattern": "^%user_id_pattern%$" },
                  "settingId": { "pattern": "^%user_id_pattern%$" },
                  "accountUuid": { "pattern": "^%user_id_pattern%$" },
                  "resource": {
                    "type": "object",
                    "required": ["type"],
                    "properties": {
                      "type": { "const": "TYPE_USER" }
                    }
                  },
                  "boolValue": { "const": false }
                }
              }
            }
          }
        }
      }
      """
    And for user "Brian" setting "auto-accept-shares" should have value "false"
    Given user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    Then user "Brian" should have sync disabled for share "textfile.txt"

  @issue-11335
  Scenario: enable auto-sync-shares after disabling it
    Given user "Alice" has uploaded file with content "lorem epsum" to "textfile.txt"
    And user "Brian" has disabled the auto-sync share
    When user "Brian" enables the auto-sync share using the settings API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value" : {
            "type": "object",
            "required": ["identifier","value"],
            "properties": {
              "identifier": {
                "type": "object",
                "required": ["extension", "bundle", "setting"],
                "properties": {
                  "extension": { "const": "ocis-accounts" },
                  "bundle": { "const": "profile" },
                  "setting": { "const": "auto-accept-shares" }
                }
              },
              "value": {
                "type": "object",
                "required": ["id", "bundleId", "settingId", "accountUuid", "resource", "boolValue"],
                "properties": {
                  "id": { "pattern": "^%user_id_pattern%$" },
                  "bundleId": { "pattern": "^%user_id_pattern%$" },
                  "settingId": { "pattern": "^%user_id_pattern%$" },
                  "accountUuid": { "pattern": "^%user_id_pattern%$" },
                  "resource": {
                    "type": "object",
                    "required": ["type"],
                    "properties": {
                      "type": { "const": "TYPE_USER" }
                    }
                  },
                  "boolValue": { "const": true }
                }
              }
            }
          }
        }
      }
      """
    And for user "Brian" setting "auto-accept-shares" should have value "true"
    Given user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    Then user "Brian" should have sync enabled for share "textfile.txt"


  Scenario: assign role to user
    When user "Admin" assigns the role "Admin" to user "Alice" using the settings API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["assignment"],
        "properties": {
          "assignment" : {
            "type": "object",
            "required": [ "id", "accountUuid", "roleId" ],
            "properties": {
              "id": { "pattern": "^%user_id_pattern%$" },
              "accountUuid": { "pattern": "^%user_id_pattern%$" },
              "roleId": { "pattern": "^%user_id_pattern%$" }
            }
          }
        }
      }
      """

  @issue-5032
  Scenario Outline: user lists assignments
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When user "Alice" tries to get list of assignment using the settings API
    Then the HTTP status code should be "<http-status-code>"
    And the setting API response should have the role "<user-role>"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["assignments"],
        "properties": {
          "assignments" : {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "uniqueItems": true,
            "items": {
              "type": "object",
              "required": ["id","accountUuid","roleId"],
              "properties": {
                "id": { "pattern": "^%user_id_pattern%$" },
                "accountUuid": { "pattern": "^%user_id_pattern%$" },
                "roleId": { "pattern": "^%user_id_pattern%$" }
              }
            }
          }
        }
      }
      """
    Examples:
      | user-role   | http-status-code |
      | Admin       | 201              |
      | Space Admin | 401              |
      | User        | 401              |
      | User Light  | 401              |


  Scenario: switch language
    When user "Alice" switches the system language to "de" using the settings API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value" : {
            "type": "object",
            "required": [ "identifier", "value"],
            "properties": {
              "identifier": {
                "type": "object",
                "required": [ "extension", "bundle", "setting"],
                "properties": {
                  "extension": { "const": "ocis-accounts" },
                  "bundle": { "const": "profile" },
                  "setting": { "const": "language" }
                }
              },
              "value": {
                "type": "object",
                "required": [ "id", "bundleId", "settingId", "accountUuid", "resource", "listValue" ],
                "properties": {
                  "id": { "pattern": "^%user_id_pattern%$" },
                  "bundleId": { "pattern": "^%user_id_pattern%$" },
                  "settingId": { "pattern": "^%user_id_pattern%$" },
                  "accountUuid": { "pattern": "^%user_id_pattern%$" },
                  "resource": {
                    "type": "object",
                    "required": ["type"],
                    "properties": {
                      "type": { "const": "TYPE_USER" }
                    }
                  },
                  "listValue": {
                    "type": "object",
                    "required": ["values"],
                    "properties": {
                      "values": {
                        "type": "array",
                        "minItems": 1,
                        "maxItems": 1,
                        "items": {
                          "type": "object",
                          "required": ["stringValue"],
                          "properties": {
                            "stringValue": { "const": "de" }
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
      }
      """
    And for user "Alice" setting "language" should have value "de"

  @issue-5079
  Scenario Outline: user lists existing roles
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When user "Alice" tries to get all existing roles using the settings API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["bundles"],
        "properties": {
          "bundles" : {
            "type": "array",
            "minItems": 4,
            "maxItems": 4,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": ["id", "name", "type", "extension", "displayName", "settings", "resource"],
                  "properties": {
                    "id": { "pattern": "^%user_id_pattern%$" },
                    "name": { "const": "spaceadmin" },
                    "type": { "const": "TYPE_ROLE" },
                    "extension": { "const": "ocis-roles" },
                    "displayName": { "const": "Space Admin" },
                    "resource": {
                      "type": "object",
                      "required": ["type"],
                      "properties": {
                        "type": { "const": "TYPE_SYSTEM" }
                      }
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id", "name", "type", "extension", "displayName", "settings", "resource"],
                  "properties": {
                    "id": { "pattern": "^%user_id_pattern%$" },
                    "name": { "const": "admin"},
                    "type": { "const": "TYPE_ROLE" },
                    "extension": { "const": "ocis-roles" },
                    "displayName": { "const": "Admin" },
                    "resource": {
                      "type": "object",
                      "required": ["type"],
                      "properties": {
                        "type": { "const": "TYPE_SYSTEM" }
                      }
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id", "name", "type", "extension", "displayName", "settings", "resource"],
                  "properties": {
                    "id": { "pattern": "^%user_id_pattern%$" },
                    "name": { "const": "user-light" },
                    "type": { "const": "TYPE_ROLE" },
                    "extension": { "const": "ocis-roles" },
                    "displayName": { "const": "User Light" },
                    "resource": {
                      "type": "object",
                      "required": ["type"],
                      "properties": {
                        "type": { "const": "TYPE_SYSTEM" }
                      }
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id", "name", "type", "extension", "displayName", "settings", "resource"],
                  "properties": {
                    "id": { "pattern": "^%user_id_pattern%$" },
                    "name": { "const": "user" },
                    "type": { "const": "TYPE_ROLE" },
                    "extension": { "const": "ocis-roles" },
                    "displayName": { "const": "User" },
                    "resource": {
                      "type": "object",
                      "required": ["type"],
                      "properties": {
                        "type": { "const": "TYPE_SYSTEM" }
                      }
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """
    Examples:
      | user-role   |
      | Admin       |
      | Space Admin |
      | User        |


  Scenario: CORS headers should be returned when setting CORS domain sending origin header in the settings api
    Given the config "OCIS_CORS_ALLOW_ORIGINS" has been set to "https://aphno.badal" for "settings" service
    When user "Alice" lists values-list with headers using the Settings API
      | header | value               |
      | Origin | https://aphno.badal |
    Then the HTTP status code should be "201"
    And the following headers should be set
      | header                      | value               |
      | Access-Control-Allow-Origin | https://aphno.badal |
