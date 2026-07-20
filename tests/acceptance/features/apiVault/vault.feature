Feature: vault
  As a user
  I want to store resource in vault storage
  So that vault resources are isolated with regular drive storage

  Background:
    Given using spaces DAV path
    And these users have been created with default attributes:
      | username |
      | Alice    |


  Scenario: user can create folders and files in personal space in vault
    Given user "Alice" has logged in via web UI
    When user "Alice" creates a folder "vaultFolder" in space "Personal" in vault using the WebDav Api
    Then the HTTP status code should be "201"
    When user "Alice" uploads a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | vaultFolder   |
      | vaultFile.txt |


  Scenario: user can create folders and files in project space in vault
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has logged in via web UI
    And user "Alice" has created a space "vault-space" in vault with the default quota using the Graph API
    When user "Alice" creates a folder "vaultFolder" in space "vault-space" in vault using the WebDav Api
    Then the HTTP status code should be "201"
    When user "Alice" uploads a file inside space "vault-space" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "vault-space" in vault should contain these entries:
      | vaultFolder   |
      | vaultFile.txt |


  Scenario: resources in drive and vault are isolated
    Given user "Alice" has logged in via web UI
    And user "Alice" has created a folder "driveFolder" in space "Personal"
    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "driveFile.txt"
    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
    When user "Alice" uploads a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | vaultFolder   |
      | vaultFile.txt |
    And for user "Alice" the space "Personal" should contain these entries:
      | driveFolder   |
      | driveFile.txt |
    And for user "Alice" the space "Personal" in vault should not contain these entries:
      | driveFolder   |
      | driveFile.txt |
    And for user "Alice" the space "Personal" should not contain these entries:
      | vaultFolder   |
      | vaultFile.txt |

  @env-config @keycloak-config
  Scenario: user can set custom auth level names
    Given the administrator has set the Keycloak realm attribute "acr.loa.map" to '{"regular":"1","testing":"2"}'
    And the config "OCIS_MFA_AUTH_LEVEL_NAMES" has been set to "testing"
    And user "Alice" has logged in via web UI
    When user "Alice" uploads a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Alice" should have acr value "testing"


  Scenario: check capabilities endpoint for vault
    Given using OCS API version "2"
    And user "Alice" has logged in via web UI
    When user "Alice" retrieves the vault mode capabilities using the capabilities API
    Then the OCS status code should be "200"
    And the HTTP status code should be "200"
    And the ocs JSON data of the response should match
      """
      {
        "type": "object",
        "required": [ "capabilities" ],
        "properties": {
          "capabilities": {
            "type": "object",
            "required": [
              "core",
              "files",
              "files_sharing",
              "auth",
              "vault"
            ],
            "properties": {
              "files_sharing": {
                "type": "object",
                "required": [
                  "api_enabled",
                  "default_permissions",
                  "public",
                  "resharing",
                  "federation",
                  "group_sharing",
                  "share_with_group_members_only",
                  "share_with_membership_groups_only",
                  "auto_accept_share",
                  "user_enumeration"
                ],
                "properties": {
                  "federation": {
                    "type": "object",
                    "required": [
                      "outgoing",
                      "incoming"
                    ],
                    "properties": {
                      "outgoing": {
                        "const": false
                      },
                      "incoming": {
                        "const": false
                      }
                    }
                  },
                  "public": {
                    "type": "object",
                    "required": [
                      "enabled",
                      "multiple",
                      "upload",
                      "supports_upload_only",
                      "send_mail",
                      "social_share"
                    ],
                    "properties": {
                      "enabled": {
                        "const": false
                      }
                    }
                  }
                }
              },
              "auth": {
                "type": "object",
                "required": [
                  "mfa"
                ],
                "properties": {
                  "mfa": {
                    "type": "object",
                    "required": [
                      "enabled",
                      "levelnames"
                    ],
                    "properties": {
                      "enabled": {
                        "const": true
                      },
                      "levelnames": {
                        "type": "array",
                        "minItems": 1,
                        "maxItems": 1,
                        "items": {
                          "const": "advanced"
                        }
                      }
                    }
                  }
                }
              },
              "vault": {
                "type": "object",
                "required": [
                  "enabled",
                  "vault_storage_provider"
                ],
                "properties": {
                  "enabled": {
                    "const": true
                  },
                  "vault_storage_provider": {
                    "pattern": "%uuidv4_pattern%"
                  }
                }
              }
            }
          }
        }
      }
      """


  Scenario Outline: send share invitation for project space in vault to user with different roles (permissions endpoint)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has logged in via web UI
    And user "Brian" has been created with default attributes
    And user "Brian" has logged in via web UI
    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | storage         | vault              |
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
                "grantedToV2",
                "roles"
              ],
              "properties": {
                "grantedToV2": {
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
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: send share invitation for disabled project space in vault to user with different roles (permissions endpoint)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has logged in via web UI
    And user "Brian" has been created with default attributes
    And user "Brian" has logged in via web UI
    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
    And user "Admin" has disabled a space "new-space" in vault
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | storage         | vault              |
    Then the HTTP status code should be "404"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "const": "itemNotFound"
              },
              "message": {
                "type": "string",
                "pattern": "^stat: error: not found: %user_id_pattern%$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: send share invitation for deleted project space in vault to user with different roles (permissions endpoint)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has logged in via web UI
    And user "Brian" has been created with default attributes
    And user "Brian" has logged in via web UI
    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
    And user "Admin" has disabled a space "new-space" in vault
    And user "Admin" has deleted a space "new-space" in vault
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | storage         | vault              |
    Then the HTTP status code should be "404"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "const": "itemNotFound"
              },
              "message": {
                "const": "stat: error: not found: "
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: try to send share invitation for personal space in vault to user with different roles (permissions endpoint)
    Given user "Alice" has logged in via web UI
    And user "Brian" has been created with default attributes
    And user "Brian" has logged in via web UI
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | storage         | vault              |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code", "innererror", "message"],
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
                "const": "space type is not eligible for sharing"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: try to share Shares space in vault with a user (permissions endpoint)
    Given user "Alice" has logged in via web UI
    And user "Brian" has been created with default attributes
    And user "Brian" has logged in via web UI
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | Shares             |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | storage         | vault              |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code", "innererror", "message"],
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
                "const": "<error-message>"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | error-message                        |
      | Space Viewer     | role not applicable to this resource |
      | Space Editor     | role not applicable to this resource |
      | Manager          | role not applicable to this resource |


  Scenario Outline: invite user to a project space in vault with different roles using root endpoint
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has logged in via web UI
    And user "Brian" has been created with default attributes
    And user "Brian" has logged in via web UI
    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using root endpoint of the Graph API:
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | storage         | vault              |
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
                "grantedToV2",
                "roles"
              ],
              "properties": {
                "grantedToV2": {
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
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: try to invite user to personal drive in vault with different roles using root endpoint
    Given user "Alice" has logged in via web UI
    And user "Brian" has been created with default attributes
    And user "Brian" has logged in via web UI
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | storage         | vault              |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code", "innererror", "message"],
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
                "const": "unsupported space type"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: try to invite user to Shares drive in vault with different roles using root endpoint
    Given user "Alice" has logged in via web UI
    And user "Brian" has been created with default attributes
    And user "Brian" has logged in via web UI
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | Shares             |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
      | storage         | vault              |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code", "innererror", "message"],
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
                "const": "unsupported space type"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |
