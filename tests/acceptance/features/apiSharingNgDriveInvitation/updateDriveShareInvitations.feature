Feature: Update permission of a share
  As a user
  I want to update drive invitations
  So that I can have more control over my shares and manage it
  https://owncloud.dev/libre-graph-api/#/drives.permissions/UpdatePermission

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: space admin updates role of a member in project space (permissions endpoint)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the space share for user "Brian" with the following using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | NewSpace               |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
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
            "required": [
              "user"
            ],
            "properties":{
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
          "id": {
            "type": "string",
            "pattern": "^u:%user_id_pattern%$"
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
      """
    Examples:
      | permissions-role | new-permissions-role |
      | Space Viewer     | Space Editor         |
      | Space Viewer     | Manager              |
      | Space Editor     | Space Viewer         |
      | Space Editor     | Manager              |
      | Manager          | Space Editor         |
      | Manager          | Space Viewer         |

  @issue-8905
  Scenario Outline: update role of a shared project space to group with different roles using root endpoint
    Given using spaces DAV path
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | grp1                 |
      | shareType          | group                |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-15T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | permissionsRole    | <new-permissions-role> |
      | space              | NewSpace               |
      | shareType          | group                  |
      | sharee             | grp1                   |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
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
            "required": [
              "group"
            ],
            "properties":{
              "group": {
                "type": "object",
                "required": [
                  "displayName",
                  "id"
                ],
                "properties": {
                  "displayName": {
                    "const": "grp1"
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%group_id_pattern%$"
                  }
                }
              }
            }
          },
          "id": {
            "type": "string",
            "pattern": "^g:%group_id_pattern%$"
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
      """
    Examples:
      | permissions-role | new-permissions-role |
      | Space Viewer     | Space Editor         |
      | Space Viewer     | Manager              |
      | Space Editor     | Space Viewer         |
      | Space Editor     | Manager              |
      | Manager          | Space Viewer         |
      | Manager          | Space Editor         |

  @issue-8905
  Scenario Outline: remove expiration date of a shared project space to group using root endpoint
    Given using spaces DAV path
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | grp1                 |
      | shareType          | group                |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-15T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | expirationDateTime |          |
      | space              | NewSpace |
      | shareType          | group    |
      | sharee             | grp1     |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
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
            "required": [
              "group"
            ],
            "properties":{
              "user": {
                "type": "object",
                "required": [
                  "displayName",
                  "id"
                ],
                "properties": {
                  "displayName": {
                    "const": "grp1"
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%user_id_pattern%$"
                  }
                }
              }
            }
          },
          "id": {
            "type": "string",
            "pattern": "^g:%group_id_pattern%$"
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
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Viewer     |
      | Space Editor     |
      | Space Editor     |
      | Manager          |
      | Manager          |

  @issue-8905
  Scenario Outline: update expiration date of a shared project space to group with different roles using root endpoint
    Given using spaces DAV path
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | grp1                 |
      | shareType          | group                |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-14T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | expirationDateTime | 2200-07-15T14:00:00Z |
      | space              | NewSpace             |
      | shareType          | group                |
      | sharee             | grp1                 |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expirationDateTime",
          "grantedToV2",
          "id",
          "roles"
        ],
        "properties": {
          "expirationDateTime": {
            "type": "string",
            "enum": ["2200-07-15T14:00:00Z"]
          },
          "grantedToV2": {
            "type": "object",
            "required": [
              "group"
            ],
            "properties":{
              "group": {
                "type": "object",
                "required": [
                  "displayName",
                  "id"
                ],
                "properties": {
                  "displayName": {
                    "const": "grp1"
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%group_id_pattern%$"
                  }
                }
              }
            }
          },
          "id": {
            "type": "string",
            "pattern": "^g:%group_id_pattern%$"
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
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: remove expiration date and update role of a shared project space at once to group with different roles at once using root endpoint
    Given using spaces DAV path
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | grp1                 |
      | shareType          | group                |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-15T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | expirationDateTime |                        |
      | permissionsRole    | <new-permissions-role> |
      | space              | NewSpace               |
      | shareType          | group                  |
      | sharee             | grp1                   |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
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
            "required": [
              "group"
            ],
            "properties":{
              "user": {
                "type": "object",
                "required": [
                  "displayName",
                  "id"
                ],
                "properties": {
                  "displayName": {
                    "const": "grp1"
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%user_id_pattern%$"
                  }
                }
              }
            }
          },
          "id": {
            "type": "string",
            "pattern": "^g:%group_id_pattern%$"
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
      """
    Examples:
      | permissions-role | new-permissions-role |
      | Space Viewer     | Space Editor         |
      | Space Viewer     | Manager              |
      | Space Editor     | Space Viewer         |
      | Space Editor     | Manager              |
      | Manager          | Space Viewer         |
      | Manager          | Space Editor         |


  Scenario Outline: update expiration date and role of a shared project space at once to group with different roles at once using root endpoint
    Given using spaces DAV path
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | grp1                 |
      | shareType          | group                |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-14T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | expirationDateTime | 2200-07-15T14:00:00Z   |
      | permissionsRole    | <new-permissions-role> |
      | space              | NewSpace               |
      | shareType          | group                  |
      | sharee             | grp1                   |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expirationDateTime",
          "grantedToV2",
          "id",
          "roles"
        ],
        "properties": {
          "expirationDateTime": {
            "type": "string",
            "enum": ["2200-07-15T14:00:00Z"]
          },
          "grantedToV2": {
            "type": "object",
            "required": [
              "group"
            ],
            "properties":{
              "group": {
                "type": "object",
                "required": [
                  "displayName",
                  "id"
                ],
                "properties": {
                  "displayName": {
                    "const": "grp1"
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%group_id_pattern%$"
                  }
                }
              }
            }
          },
          "id": {
            "type": "string",
            "pattern": "^g:%group_id_pattern%$"
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
      """
    Examples:
      | permissions-role | new-permissions-role |
      | Space Viewer     | Space Editor         |
      | Space Viewer     | Manager              |
      | Space Editor     | Space Viewer         |
      | Space Editor     | Manager              |
      | Manager          | Space Viewer         |
      | Manager          | Space Editor         |

  @issue-8905
  Scenario Outline: update role of a shared project space to user with different roles using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | Brian                |
      | shareType          | user                 |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-15T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | permissionsRole    | <new-permissions-role> |
      | space              | NewSpace               |
      | shareType          | user                   |
      | sharee             | Brian                  |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
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
            "required": [
              "user"
            ],
            "properties":{
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
          "id": {
            "type": "string",
            "pattern": "^u:%user_id_pattern%$"
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
      """
    Examples:
      | permissions-role | new-permissions-role |
      | Space Viewer     | Space Editor         |
      | Space Viewer     | Manager              |
      | Space Editor     | Space Viewer         |
      | Space Editor     | Manager              |
      | Manager          | Space Viewer         |
      | Manager          | Space Editor         |

  @issue-8905
  Scenario Outline: remove expiration date of a shared project space to user using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | Brian                |
      | shareType          | user                 |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-15T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | expirationDateTime |                        |
      | space              | NewSpace               |
      | shareType          | user                   |
      | sharee             | Brian                  |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
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
            "required": [
              "user"
            ],
            "properties":{
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
          "id": {
            "type": "string",
            "pattern": "^u:%user_id_pattern%$"
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
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Viewer     |
      | Space Editor     |
      | Space Editor     |
      | Manager          |
      | Manager          |


  Scenario Outline: update expiration date of a share to user with different roles at once using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | Brian                |
      | shareType          | user                 |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-14T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | expirationDateTime | 2200-07-15T14:00:00Z   |
      | permissionsRole    | <new-permissions-role> |
      | space              | NewSpace               |
      | shareType          | user                   |
      | sharee             | Brian                  |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expirationDateTime",
          "grantedToV2",
          "id",
          "roles"
        ],
        "properties": {
          "expirationDateTime": {
            "type": "string",
            "enum": ["2200-07-15T14:00:00Z"]
          },
          "grantedToV2": {
            "type": "object",
            "required": [
              "user"
            ],
            "properties":{
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
          "id": {
            "type": "string",
            "pattern": "^u:%user_id_pattern%$"
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
      """
    Examples:
      | permissions-role | new-permissions-role |
      | Space Viewer     | Space Editor         |
      | Space Viewer     | Manager              |
      | Space Editor     | Space Viewer         |
      | Space Editor     | Manager              |
      | Manager          | Space Viewer         |
      | Manager          | Space Editor         |


  Scenario Outline: remove expiration and update role of a share to user with different roles at once using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | Brian                |
      | shareType          | user                 |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-15T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | expirationDateTime |                        |
      | permissionsRole    | <new-permissions-role> |
      | space              | NewSpace               |
      | shareType          | user                   |
      | sharee             | Brian                  |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
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
            "required": [
              "user"
            ],
            "properties":{
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
          "id": {
            "type": "string",
            "pattern": "^u:%user_id_pattern%$"
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
      """
    Examples:
      | permissions-role | new-permissions-role |
      | Space Viewer     | Space Editor         |
      | Space Viewer     | Manager              |
      | Space Editor     | Space Viewer         |
      | Space Editor     | Manager              |
      | Manager          | Space Viewer         |
      | Manager          | Space Editor         |


  Scenario Outline: update expiration date and role of a share to user with different roles at once using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space              | NewSpace             |
      | sharee             | Brian                |
      | shareType          | user                 |
      | permissionsRole    | <permissions-role>   |
      | expirationDateTime | 2200-07-14T14:00:00Z |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | expirationDateTime | 2200-07-15T14:00:00Z   |
      | permissionsRole    | <new-permissions-role> |
      | space              | NewSpace               |
      | shareType          | user                   |
      | sharee             | Brian                  |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expirationDateTime",
          "grantedToV2",
          "id",
          "roles"
        ],
        "properties": {
          "expirationDateTime": {
            "type": "string",
            "enum": ["2200-07-15T14:00:00Z"]
          },
          "grantedToV2": {
            "type": "object",
            "required": [
              "user"
            ],
            "properties":{
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
          "id": {
            "type": "string",
            "pattern": "^u:%user_id_pattern%$"
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
      """
    Examples:
      | permissions-role | new-permissions-role |
      | Space Viewer     | Space Editor         |
      | Space Viewer     | Manager              |
      | Space Editor     | Space Viewer         |
      | Space Editor     | Manager              |
      | Manager          | Space Viewer         |
      | Manager          | Space Editor         |

  @issue-10768
  Scenario Outline: sharer updates share permissions role of space to Space Editor Without Versions without enabling it
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | permissionsRole | Space Editor Without Versions |
      | space           | new-space                     |
      | shareType       | user                          |
      | sharee          | Brian                         |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@libre.graph.permissions.actions",
          "grantedToV2",
          "id"
        ],
        "properties": {
          "@libre.graph.permissions.actions": {
            "const": [
              "libre.graph/driveItem/children/create",
              "libre.graph/driveItem/standard/delete",
              "libre.graph/driveItem/path/read",
              "libre.graph/driveItem/quota/read",
              "libre.graph/driveItem/content/read",
              "libre.graph/driveItem/upload/create",
              "libre.graph/driveItem/permissions/read",
              "libre.graph/driveItem/children/read",
              "libre.graph/driveItem/deleted/read",
              "libre.graph/driveItem/path/update",
              "libre.graph/driveItem/deleted/update",
              "libre.graph/driveItem/basic/read"
            ]
          },
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
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
                    "pattern": "^%user_id_pattern%$"
                  }
                }
              }
            }
          },
          "id": {
            "pattern": "^u:%user_id_pattern%$"
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |
