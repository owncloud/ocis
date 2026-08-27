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


#  Scenario: user can create folders and files in project space in vault
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Alice" has created a space "vault-space" in vault with the default quota using the Graph API
#    When user "Alice" creates a folder "vaultFolder" in space "vault-space" in vault using the WebDav Api
#    Then the HTTP status code should be "201"
#    When user "Alice" uploads a file inside space "vault-space" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And for user "Alice" the space "vault-space" in vault should contain these entries:
#      | vaultFolder   |
#      | vaultFile.txt |
#
#
#  Scenario: resources in drive and vault are isolated
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "driveFolder" in space "Personal"
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "driveFile.txt"
#    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
#    When user "Alice" uploads a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And for user "Alice" the space "Personal" in vault should contain these entries:
#      | vaultFolder   |
#      | vaultFile.txt |
#    And for user "Alice" the space "Personal" should contain these entries:
#      | driveFolder   |
#      | driveFile.txt |
#    And for user "Alice" the space "Personal" in vault should not contain these entries:
#      | driveFolder   |
#      | driveFile.txt |
#    And for user "Alice" the space "Personal" should not contain these entries:
#      | vaultFolder   |
#      | vaultFile.txt |
#
#  @env-config @keycloak-config
#  Scenario: user can set custom auth level names
#    Given the administrator has set the Keycloak realm attribute "acr.loa.map" to '{"regular":"1","testing":"2"}'
#    And the config "OCIS_MFA_AUTH_LEVEL_NAMES" has been set to "testing"
#    And user "Alice" has logged in via web UI
#    When user "Alice" uploads a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And user "Alice" should have a JWT token with an ACR value "testing"
#
#
#  Scenario: check capabilities endpoint for vault
#    Given using OCS API version "2"
#    And user "Alice" has logged in via web UI
#    When user "Alice" retrieves the vault mode capabilities using the capabilities API
#    Then the OCS status code should be "200"
#    And the HTTP status code should be "200"
#    And the ocs JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": [ "capabilities" ],
#        "properties": {
#          "capabilities": {
#            "type": "object",
#            "required": [
#              "core",
#              "files",
#              "files_sharing",
#              "auth",
#              "vault"
#            ],
#            "properties": {
#              "files_sharing": {
#                "type": "object",
#                "required": [
#                  "api_enabled",
#                  "default_permissions",
#                  "public",
#                  "resharing",
#                  "federation",
#                  "group_sharing",
#                  "share_with_group_members_only",
#                  "share_with_membership_groups_only",
#                  "auto_accept_share",
#                  "user_enumeration"
#                ],
#                "properties": {
#                  "federation": {
#                    "type": "object",
#                    "required": [
#                      "outgoing",
#                      "incoming"
#                    ],
#                    "properties": {
#                      "outgoing": {
#                        "const": false
#                      },
#                      "incoming": {
#                        "const": false
#                      }
#                    }
#                  },
#                  "public": {
#                    "type": "object",
#                    "required": [
#                      "enabled",
#                      "multiple",
#                      "upload",
#                      "supports_upload_only",
#                      "send_mail",
#                      "social_share"
#                    ],
#                    "properties": {
#                      "enabled": {
#                        "const": false
#                      }
#                    }
#                  }
#                }
#              },
#              "auth": {
#                "type": "object",
#                "required": [
#                  "mfa"
#                ],
#                "properties": {
#                  "mfa": {
#                    "type": "object",
#                    "required": [
#                      "enabled",
#                      "levelnames"
#                    ],
#                    "properties": {
#                      "enabled": {
#                        "const": true
#                      },
#                      "levelnames": {
#                        "type": "array",
#                        "minItems": 1,
#                        "maxItems": 1,
#                        "items": {
#                          "const": "advanced"
#                        }
#                      }
#                    }
#                  }
#                }
#              },
#              "vault": {
#                "type": "object",
#                "required": [
#                  "enabled",
#                  "vault_storage_provider"
#                ],
#                "properties": {
#                  "enabled": {
#                    "const": true
#                  },
#                  "vault_storage_provider": {
#                    "pattern": "%uuidv4_pattern%"
#                  }
#                }
#              }
#            }
#          }
#        }
#      }
#      """
#
#
#  Scenario: user copies folder from drive to vault
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "driveFolder" in space "Personal"
#    When user "Alice" copies folder "driveFolder" from space "Personal" to "driveFolder" inside space "Personal" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And for user "Alice" the space "Personal" in vault should contain these entries:
#      | driveFolder |
#    And for user "Alice" the space "Personal" should contain these entries:
#      | driveFolder |
#
#
#  Scenario: user copies file from drive to vault
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testfile.txt"
#    When user "Alice" copies file "testfile.txt" from space "Personal" to "testfile.txt" inside space "Personal" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And for user "Alice" the content of the file "testfile.txt" of the space "Personal" in vault should be "some content"
#    And for user "Alice" the space "Personal" should contain these entries:
#      | testfile.txt |
#
#
#  Scenario: user tries to copy folder from vault to drive
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
#    When user "Alice" copies folder "vaultFolder" from space "Personal" in vault to "vaultFolder" inside space "Personal" using the WebDAV API
#    Then the HTTP status code should be "409"
#    And for user "Alice" the space "Personal" should not contain these entries:
#      | vaultFolder |
#    And for user "Alice" the space "Personal" in vault should contain these entries:
#      | vaultFolder |
#
#
#  Scenario: user tries to copy file from vault to drive
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testfile.txt" in vault
#    When user "Alice" copies file "testfile.txt" from space "Personal" in vault to "testfile.txt" inside space "Personal" using the WebDAV API
#    Then the HTTP status code should be "409"
#    And for user "Alice" the space "Personal" should not contain these entries:
#      | testfile.txt |
#    And for user "Alice" the space "Personal" in vault should contain these entries:
#      | testfile.txt |
#
#
#  Scenario: user copies sub-folder from drive to vault
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "driveFolder" in space "Personal"
#    And user "Alice" has created a folder "driveFolder/subFolder" in space "Personal"
#    When user "Alice" copies folder "driveFolder/subFolder" from space "Personal" to "subFolder" inside space "Personal" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And for user "Alice" the space "Personal" in vault should contain these entries:
#      | subFolder |
#    And for user "Alice" folder "driveFolder" of the space "Personal" should contain these entries:
#      | subFolder |
#
#
#  Scenario: user copies file inside folder from drive to vault
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "driveFolder" in space "Personal"
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "driveFolder/testfile.txt"
#    When user "Alice" copies file "driveFolder/testfile.txt" from space "Personal" to "testfile.txt" inside space "Personal" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And for user "Alice" the content of the file "testfile.txt" of the space "Personal" in vault should be "some content"
#    And for user "Alice" folder "driveFolder" of the space "Personal" should contain these entries:
#      | testfile.txt |
#
#
#  Scenario: user copies sub-folder from drive to a folder in vault
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "driveFolder" in space "Personal"
#    And user "Alice" has created a folder "driveFolder/subFolder" in space "Personal"
#    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
#    When user "Alice" copies folder "driveFolder/subFolder" from space "Personal" to "vaultFolder/subFolder" inside space "Personal" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And for user "Alice" folder "vaultFolder" of the space "Personal" in vault should contain these entries:
#      | subFolder |
#    And for user "Alice" folder "driveFolder" of the space "Personal" should contain these entries:
#      | subFolder |
#
#
#  Scenario: user copies file inside folder from drive to a folder in vault
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "driveFolder" in space "Personal"
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "driveFolder/testfile.txt"
#    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
#    When user "Alice" copies file "driveFolder/testfile.txt" from space "Personal" to "vaultFolder/testfile.txt" inside space "Personal" in vault using the WebDAV API
#    Then the HTTP status code should be "201"
#    And for user "Alice" the content of the file "vaultFolder/testfile.txt" of the space "Personal" in vault should be "some content"
#    And for user "Alice" folder "driveFolder" of the space "Personal" should contain these entries:
#      | testfile.txt |
#
#
#  Scenario: user tries to create a public link of a folder inside vault
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
#    When user "Alice" creates the following resource link share using the Graph API:
#      | resource        | vaultFolder |
#      | space           | Personal    |
#      | permissionsRole | View        |
#      | storage         | vault       |
#    Then the HTTP status code should be "400"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": ["error"],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": ["code", "innererror", "message"],
#            "properties": {
#              "code": {
#                "const": "invalidRequest"
#              },
#              "innererror": {
#                "type": "object",
#                "required": [
#                  "date",
#                  "request-id"
#                ]
#              },
#              "message": {
#                "const": "public links are not allowed for vault resources"
#              }
#            }
#          }
#        }
#      }
#      """
#
#
#  Scenario: user tries to create a public link of a file inside vault
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testfile.txt" in vault
#    When user "Alice" creates the following resource link share using the Graph API:
#      | resource        | testfile.txt |
#      | space           | Personal     |
#      | permissionsRole | View         |
#      | storage         | vault        |
#    Then the HTTP status code should be "400"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": ["error"],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": ["code", "innererror", "message"],
#            "properties": {
#              "code": {
#                "const": "invalidRequest"
#              },
#              "innererror": {
#                "type": "object",
#                "required": [
#                  "date",
#                  "request-id"
#                ]
#              },
#              "message": {
#                "const": "public links are not allowed for vault resources"
#              }
#            }
#          }
#        }
#      }
#      """
#
#
#  Scenario: user tries to create a public link of a space root inside vault
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Alice" has created a space "vault-space" in vault with the default quota using the Graph API
#    When user "Alice" tries to create the following space link share using permissions endpoint of the Graph API:
#      | space           | vault-space |
#      | permissionsRole | View        |
#      | storage         | vault       |
#    Then the HTTP status code should be "400"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": ["error"],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": ["code", "innererror", "message"],
#            "properties": {
#              "code": {
#                "const": "invalidRequest"
#              },
#              "innererror": {
#                "type": "object",
#                "required": [
#                  "date",
#                  "request-id"
#                ]
#              },
#              "message": {
#                "const": "public links are not allowed for vault resources"
#              }
#            }
#          }
#        }
#      }
#      """
#
#
#  Scenario Outline: send share invitation for project space in vault to user with different roles (permissions endpoint)
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Brian" has been created with default attributes
#    And user "Brian" has logged in via web UI
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
#      | space           | new-space          |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "200"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": [
#          "value"
#        ],
#        "properties": {
#          "value": {
#            "type": "array",
#            "minItems": 1,
#            "maxItems": 1,
#            "items": {
#              "type": "object",
#              "required": [
#                "grantedToV2",
#                "roles"
#              ],
#              "properties": {
#                "grantedToV2": {
#                  "type": "object",
#                  "required": [
#                    "user"
#                  ],
#                  "properties": {
#                    "user": {
#                      "type": "object",
#                      "required": [
#                        "displayName",
#                        "id"
#                      ],
#                      "properties": {
#                        "displayName": {
#                          "const": "Brian Murphy"
#                        },
#                        "id": {
#                          "type": "string",
#                          "pattern": "^%user_id_pattern%$"
#                        }
#                      }
#                    }
#                  }
#                },
#                "roles": {
#                  "type": "array",
#                  "minItems": 1,
#                  "maxItems": 1,
#                  "items": {
#                    "type": "string",
#                    "pattern": "^%role_id_pattern%$"
#                  }
#                }
#              }
#            }
#          }
#        }
#      }
#      """
#    Examples:
#      | permissions-role |
#      | Space Viewer     |
#      | Space Editor     |
#      | Manager          |
#
#
#  Scenario Outline: send share invitation for disabled project space in vault to user with different roles (permissions endpoint)
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Brian" has been created with default attributes
#    And user "Brian" has logged in via web UI
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    And user "Admin" has disabled a space "new-space" in vault
#    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
#      | space           | new-space          |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "404"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": [
#          "error"
#        ],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": [
#              "code",
#              "message"
#            ],
#            "properties": {
#              "code": {
#                "const": "itemNotFound"
#              },
#              "message": {
#                "type": "string",
#                "pattern": "^stat: error: not found: %user_id_pattern%$"
#              }
#            }
#          }
#        }
#      }
#      """
#    Examples:
#      | permissions-role |
#      | Space Viewer     |
#      | Space Editor     |
#      | Manager          |
#
#
#  Scenario Outline: send share invitation for deleted project space in vault to user with different roles (permissions endpoint)
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Brian" has been created with default attributes
#    And user "Brian" has logged in via web UI
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    And user "Admin" has disabled a space "new-space" in vault
#    And user "Admin" has deleted a space "new-space" in vault
#    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
#      | space           | new-space          |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "404"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": [
#          "error"
#        ],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": [
#              "code",
#              "message"
#            ],
#            "properties": {
#              "code": {
#                "const": "itemNotFound"
#              },
#              "message": {
#                "const": "stat: error: not found: "
#              }
#            }
#          }
#        }
#      }
#      """
#    Examples:
#      | permissions-role |
#      | Space Viewer     |
#      | Space Editor     |
#      | Manager          |
#
#
#  Scenario Outline: try to send share invitation for personal space in vault to user with different roles (permissions endpoint)
#    Given user "Alice" has logged in via web UI
#    And user "Brian" has been created with default attributes
#    And user "Brian" has logged in via web UI
#    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
#      | space           | Personal           |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "400"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": ["error"],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": ["code", "innererror", "message"],
#            "properties": {
#              "code": {
#                "const": "invalidRequest"
#              },
#              "innererror": {
#                "type": "object",
#                "required": [
#                  "date",
#                  "request-id"
#                ]
#              },
#              "message": {
#                "const": "space type is not eligible for sharing"
#              }
#            }
#          }
#        }
#      }
#      """
#    Examples:
#      | permissions-role |
#      | Space Viewer     |
#      | Space Editor     |
#      | Manager          |
#
#
#  Scenario Outline: try to share Shares space in vault with a user (permissions endpoint)
#    Given user "Alice" has logged in via web UI
#    And user "Brian" has been created with default attributes
#    And user "Brian" has logged in via web UI
#    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
#      | space           | Shares             |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "400"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": ["error"],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": ["code", "innererror", "message"],
#            "properties": {
#              "code": {
#                "const": "invalidRequest"
#              },
#              "innererror": {
#                "type": "object",
#                "required": [
#                  "date",
#                  "request-id"
#                ]
#              },
#              "message": {
#                "const": "<error-message>"
#              }
#            }
#          }
#        }
#      }
#      """
#    Examples:
#      | permissions-role | error-message                        |
#      | Space Viewer     | role not applicable to this resource |
#      | Space Editor     | role not applicable to this resource |
#      | Manager          | role not applicable to this resource |
#
#
#  Scenario Outline: invite user to a project space in vault with different roles using root endpoint
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Brian" has been created with default attributes
#    And user "Brian" has logged in via web UI
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    When user "Alice" sends the following space share invitation using root endpoint of the Graph API:
#      | space           | new-space          |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "200"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": [
#          "value"
#        ],
#        "properties": {
#          "value": {
#            "type": "array",
#            "minItems": 1,
#            "maxItems": 1,
#            "items": {
#              "type": "object",
#              "required": [
#                "grantedToV2",
#                "roles"
#              ],
#              "properties": {
#                "grantedToV2": {
#                  "type": "object",
#                  "required": [
#                    "user"
#                  ],
#                  "properties": {
#                    "user": {
#                      "type": "object",
#                      "required": [
#                        "displayName",
#                        "id"
#                      ],
#                      "properties": {
#                        "displayName": {
#                          "type": "string",
#                          "const": "Brian Murphy"
#                        },
#                        "id": {
#                          "type": "string",
#                          "pattern": "^%user_id_pattern%$"
#                        }
#                      }
#                    }
#                  }
#                },
#                "roles": {
#                  "type": "array",
#                  "minItems": 1,
#                  "maxItems": 1,
#                  "items": {
#                    "type": "string",
#                    "pattern": "^%role_id_pattern%$"
#                  }
#                }
#              }
#            }
#          }
#        }
#      }
#      """
#    Examples:
#      | permissions-role |
#      | Space Viewer     |
#      | Space Editor     |
#      | Manager          |
#
#
#  Scenario Outline: try to invite user to personal drive in vault with different roles using root endpoint
#    Given user "Alice" has logged in via web UI
#    And user "Brian" has been created with default attributes
#    And user "Brian" has logged in via web UI
#    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
#      | space           | Personal           |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "400"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": ["error"],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": ["code", "innererror", "message"],
#            "properties": {
#              "code": {
#                "const": "invalidRequest"
#              },
#              "innererror": {
#                "type": "object",
#                "required": [
#                  "date",
#                  "request-id"
#                ]
#              },
#              "message": {
#                "const": "unsupported space type"
#              }
#            }
#          }
#        }
#      }
#      """
#    Examples:
#      | permissions-role |
#      | Space Viewer     |
#      | Space Editor     |
#      | Manager          |
#
#
#  Scenario Outline: try to invite user to Shares drive in vault with different roles using root endpoint
#    Given user "Alice" has logged in via web UI
#    And user "Brian" has been created with default attributes
#    And user "Brian" has logged in via web UI
#    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
#      | space           | Shares             |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "400"
#    And the JSON data of the response should match
#      """
#      {
#        "type": "object",
#        "required": ["error"],
#        "properties": {
#          "error": {
#            "type": "object",
#            "required": ["code", "innererror", "message"],
#            "properties": {
#              "code": {
#                "const": "invalidRequest"
#              },
#              "innererror": {
#                "type": "object",
#                "required": [
#                  "date",
#                  "request-id"
#                ]
#              },
#              "message": {
#                "const": "unsupported space type"
#              }
#            }
#          }
#        }
#      }
#      """
#    Examples:
#      | permissions-role |
#      | Space Viewer     |
#      | Space Editor     |
#      | Manager          |
#
#
#  Scenario: search results for resources in Personal space should be isolated between vault and drive
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "testDriveFolder" in space "Personal"
#    And user "Alice" has created a folder "testVaultFolder" in space "Personal" in vault
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testDriveFile.txt"
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testVaultFile.txt" in vault
#    When user "Alice" searches for "*test*" inside space "Personal" in vault using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "2" entries
#    And the search result of user "Alice" should contain only these entries:
#      | testVaultFolder   |
#      | testVaultFile.txt |
#    When user "Alice" searches for "*test*" inside space "Personal" using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "2" entries
#    And the search result of user "Alice" should contain only these entries:
#      | testDriveFolder   |
#      | testDriveFile.txt |
#
#
#  Scenario: search results for resources inside folder with same name should be isolated between vault and drive
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "newFolder" in space "Personal"
#    And user "Alice" has created a folder "newFolder" in space "Personal" in vault
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "newFolder/testDriveFile.txt"
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "newFolder/testVaultFile.txt" in vault
#    When user "Alice" searches for "*test*" inside folder "newFolder" in space "Personal" in vault using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "1" entries
#    And the search result of user "Alice" should contain only these entries:
#      | newFolder/testVaultFile.txt |
#    When user "Alice" searches for "*test*" inside folder "newFolder" in space "Personal" using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "1" entries
#    And the search result of user "Alice" should contain only these entries:
#      | newFolder/testDriveFile.txt |
#
#
#  Scenario: search result for resources inside project spaces with same name should be isolated between vault and drive
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Alice" has created a space "new-space" with the default quota using the Graph API
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    And user "Alice" has created a folder "testDriveFolder" in space "new-space"
#    And user "Alice" has created a folder "testVaultFolder" in space "new-space" in vault
#    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testDriveFile.txt"
#    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testVaultFile.txt" in vault
#    When user "Alice" searches for "*test*" inside space "new-space" in vault using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "2" entries
#    And the search result of user "Alice" should contain only these entries:
#      | testVaultFolder   |
#      | testVaultFile.txt |
#    When user "Alice" searches for "*test*" inside space "new-space" using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "2" entries
#    And the search result of user "Alice" should contain only these entries:
#      | testDriveFolder   |
#      | testDriveFile.txt |
#
#  @tikaServiceNeeded
#  Scenario: search result by content of file should be isolated between vault and drive
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has uploaded a file inside space "Personal" with content "content of file in drive" to "testDriveFile.txt"
#    And user "Alice" has uploaded a file inside space "Personal" with content "content of file in vault" to "testVaultFile.txt" in vault
#    When user "Alice" searches for "Content:content" in vault using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "1" entries
#    And the search result of user "Alice" should contain only these files:
#      | testVaultFile.txt |
#    When user "Alice" searches for "Content:content" using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "1" entries
#    And the search result of user "Alice" should contain only these files:
#      | testDriveFile.txt |
#
#  @tikaServiceNeeded
#  Scenario: search result by content of file inside project space should be isolated between vault and drive
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Alice" has created a space "new-space" with the default quota using the Graph API
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    And user "Alice" has uploaded a file inside space "new-space" with content "content of file in drive" to "testDriveFile.txt"
#    And user "Alice" has uploaded a file inside space "new-space" with content "content of file in vault" to "testVaultFile.txt" in vault
#    When user "Alice" searches for "Content:content" in vault using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "1" entries
#    And the search result of user "Alice" should contain only these files:
#      | testVaultFile.txt |
#    When user "Alice" searches for "Content:content" using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "1" entries
#    And the search result of user "Alice" should contain only these files:
#      | testDriveFile.txt |
#
#
#  Scenario: search results by resource tags should be isolated between vault and drive
#    Given user "Alice" has logged in via web UI
#    And user "Alice" has created a folder "driveFolder" in space "Personal"
#    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testDriveFile.txt"
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testVaultFile.txt" in vault
#    And user "Alice" has tagged the following files of the space "Personal":
#      | path              | tagName |
#      | testDriveFile.txt | tag1    |
#    And user "Alice" has tagged the following folders of the space "Personal":
#      | path        | tagName |
#      | driveFolder | tag1    |
#    And user "Alice" has tagged the following files of the space "Personal" in vault:
#      | path              | tagName |
#      | testVaultFile.txt | tag1    |
#    And user "Alice" has tagged the following folders of the space "Personal" in vault:
#      | path        | tagName |
#      | vaultFolder | tag1    |
#    When user "Alice" searches for "Tags:tag1" in vault using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "2" entries
#    And the search result of user "Alice" should contain only these files:
#      | testVaultFile.txt |
#      | vaultFolder       |
#    When user "Alice" searches for "Tags:tag1" using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "2" entries
#    And the search result of user "Alice" should contain only these files:
#      | testDriveFile.txt |
#      | driveFolder       |
#
#
#  Scenario: search results by resource tags inside project space should be isolated between vault and drive
#    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Alice" has created a space "new-space" with the default quota using the Graph API
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    And user "Alice" has created a folder "driveFolder" in space "new-space"
#    And user "Alice" has created a folder "vaultFolder" in space "new-space" in vault
#    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testDriveFile.txt"
#    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testVaultFile.txt" in vault
#    And user "Alice" has tagged the following files of the space "new-space":
#      | path              | tagName |
#      | testDriveFile.txt | tag1    |
#    And user "Alice" has tagged the following folders of the space "new-space":
#      | path        | tagName |
#      | driveFolder | tag1    |
#    And user "Alice" has tagged the following files of the space "new-space" in vault:
#      | path              | tagName |
#      | testVaultFile.txt | tag1    |
#    And user "Alice" has tagged the following folders of the space "new-space" in vault:
#      | path        | tagName |
#      | vaultFolder | tag1    |
#    When user "Alice" searches for "Tags:tag1" in vault using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "2" entries
#    And the search result of user "Alice" should contain only these files:
#      | testVaultFile.txt |
#      | vaultFolder       |
#    When user "Alice" searches for "Tags:tag1" using the WebDAV API
#    Then the HTTP status code should be "207"
#    And the search result should contain "2" entries
#    And the search result of user "Alice" should contain only these files:
#      | testDriveFile.txt |
#      | driveFolder       |
#
#
#  Scenario Outline: folder share received from vault and drive personal space should be isolated
#    Given user "Brian" has been created with default attributes
#    And user "Alice" has logged in via web UI
#    And user "Brian" has logged in via web UI
#    And user "Alice" has created a folder "driveFolder" in space "Personal"
#    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
#    And user "Alice" has sent the following resource share invitation:
#      | resource        | driveFolder        |
#      | space           | Personal           |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#    When user "Alice" sends the following resource share invitation using the Graph API:
#      | resource        | vaultFolder        |
#      | space           | Personal           |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "200"
#    And user "Brian" should have a share in vault "vaultFolder" synced
#    And user "Brian" should have the following resource shares:
#      | resource    | permissionsRole    | sharer | space    | storage |
#      | vaultFolder | <permissions-role> | Alice  | Personal | vault   |
#    And user "Brian" should have the following resource shares:
#      | resource    | permissionsRole    | sharer | space    |
#      | driveFolder | <permissions-role> | Alice  | Personal |
#    Examples:
#      | permissions-role |
#      | Viewer           |
#      | Editor           |
#      | Uploader         |
#
#
#  Scenario Outline: file share received from vault and drive personal space should be isolated
#    Given user "Brian" has been created with default attributes
#    And user "Alice" has logged in via web UI
#    And user "Brian" has logged in via web UI
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "driveFile.txt"
#    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault
#    And user "Alice" has sent the following resource share invitation:
#      | resource        | driveFile.txt      |
#      | space           | Personal           |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#    When user "Alice" sends the following resource share invitation using the Graph API:
#      | resource        | vaultFile.txt      |
#      | space           | Personal           |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "200"
#    And user "Brian" should have a share in vault "vaultFile.txt" synced
#    And user "Brian" should have the following resource shares:
#      | resource      | permissionsRole    | sharer | space    | storage |
#      | vaultFile.txt | <permissions-role> | Alice  | Personal | vault   |
#    And user "Brian" should have the following resource shares:
#      | resource      | permissionsRole    | sharer | space    |
#      | driveFile.txt | <permissions-role> | Alice  | Personal |
#    Examples:
#      | permissions-role |
#      | Viewer           |
#      | File Editor      |
#
#
#  Scenario Outline: folder share received from vault and drive project space should be isolated
#    Given user "Brian" has been created with default attributes
#    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Brian" has logged in via web UI
#    And user "Alice" has created a space "new-space" with the default quota using the Graph API
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    And user "Alice" has created a folder "driveFolder" in space "new-space"
#    And user "Alice" has created a folder "vaultFolder" in space "new-space" in vault
#    And user "Alice" has sent the following resource share invitation:
#      | resource        | driveFolder        |
#      | space           | new-space          |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#    When user "Alice" sends the following resource share invitation using the Graph API:
#      | resource        | vaultFolder        |
#      | space           | new-space          |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "200"
#    And user "Brian" should have a share in vault "vaultFolder" synced
#    And user "Brian" should have the following resource shares:
#      | resource    | permissionsRole    | sharer | space     | storage |
#      | vaultFolder | <permissions-role> | Alice  | new-space | vault   |
#    And user "Brian" should have the following resource shares:
#      | resource    | permissionsRole    | sharer | space     |
#      | driveFolder | <permissions-role> | Alice  | new-space |
#    Examples:
#      | permissions-role |
#      | Viewer           |
#      | Editor           |
#      | Uploader         |
#
#
#  Scenario Outline: folder share received from vault and drive project space should be isolated
#    Given user "Brian" has been created with default attributes
#    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And user "Alice" has logged in via web UI
#    And user "Brian" has logged in via web UI
#    And user "Alice" has created a space "new-space" with the default quota using the Graph API
#    And user "Alice" has created a space "new-space" in vault with the default quota using the Graph API
#    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "driveFile.txt"
#    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "vaultFile.txt" in vault
#    And user "Alice" has sent the following resource share invitation:
#      | resource        | driveFile.txt      |
#      | space           | new-space          |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#    When user "Alice" sends the following resource share invitation using the Graph API:
#      | resource        | vaultFile.txt      |
#      | space           | new-space          |
#      | sharee          | Brian              |
#      | shareType       | user               |
#      | permissionsRole | <permissions-role> |
#      | storage         | vault              |
#    Then the HTTP status code should be "200"
#    And user "Brian" should have a share in vault "vaultFile.txt" synced
#    And user "Brian" should have the following resource shares:
#      | resource      | permissionsRole    | sharer | space     | storage |
#      | vaultFile.txt | <permissions-role> | Alice  | new-space | vault   |
#    And user "Brian" should have the following resource shares:
#      | resource      | permissionsRole    | sharer | space     |
#      | driveFile.txt | <permissions-role> | Alice  | new-space |
#    Examples:
#      | permissions-role |
#      | Viewer           |
#      | File Editor      |
