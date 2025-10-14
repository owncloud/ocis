@issue-10739
Feature: Send a sharing invitations
  As the owner of a resource
  I want to be able to send invitations to other users
  So that they can have access to it

  https://owncloud.dev/libre-graph-api/#/drives.permissions/Invite

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: send share invitation to user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": [
                "createdDateTime",
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
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
                          "const": "Brian Murphy"
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
    And user "Brian" should have a share "<resource>" synced
    And user "Brian" should have the following resource shares:
      | resource   | permissionsRole    | sharer | space    |
      | <resource> | <permissions-role> | Alice  | Personal |
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation to group with different roles
    Given user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And user "Brian" should have a share "<resource>" shared by user "Alice" from space "Personal"
    And user "Carol" should have a share "<resource>" shared by user "Alice" from space "Personal"
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
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": [
                "createdDateTime",
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%group_id_pattern%$"
                        },
                        "displayName": {
                          "const": "grp1"
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
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation for a file to user with different permissions
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource          | textfile1.txt        |
      | space             | Personal             |
      | sharee            | Brian                |
      | shareType         | user                 |
      | permissionsAction | <permissions-action> |
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
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": [
                "createdDateTime",
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissions-action>$"
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
                          "const": "Brian Murphy"
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
      | permissions-action |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | versions/update    |
      | deleted/update     |
      | deleted/delete     |


  Scenario Outline: send share invitation for a folder to user with different permissions
    Given user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource          | FolderToShare        |
      | space             | Personal             |
      | sharee            | Brian                |
      | shareType         | user                 |
      | permissionsAction | <permissions-action> |
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
                "createdDateTime",
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissions-action>$"
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
                          "const": "Brian Murphy"
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
      | permissions-action |
      | children/create    |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | path/update        |
      | versions/update    |
      | deleted/update     |
      | standard/delete    |
      | deleted/delete     |


  Scenario Outline: send share invitation for a file to group with different permissions
    Given user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource          | textfile1.txt        |
      | space             | Personal             |
      | sharee            | grp1                 |
      | shareType         | group                |
      | permissionsAction | <permissions-action> |
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
                "createdDateTime",
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissions-action>$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
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
                          "const": "grp1"
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
      | permissions-action |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | versions/update    |
      | deleted/update     |
      | deleted/delete     |


  Scenario Outline: send share invitation for a folder to group with different permissions
    Given user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource          | FolderToShare        |
      | space             | Personal             |
      | sharee            | grp1                 |
      | shareType         | group                |
      | permissionsAction | <permissions-action> |
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
                "createdDateTime",
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissions-action>$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
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
                          "const": "grp1"
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
      | permissions-action |
      | children/create    |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | path/update        |
      | versions/update    |
      | deleted/update     |
      | standard/delete    |
      | deleted/delete     |


  Scenario Outline: send share invitation with expiration date to user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource           | <resource>               |
      | space              | Personal                 |
      | sharee             | Brian                    |
      | shareType          | user                     |
      | permissionsRole    | <permissions-role>       |
      | expirationDateTime | 2043-07-15T14:00:00.000Z |
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
                "createdDateTime",
                "id",
                "roles",
                "grantedToV2",
                "expirationDateTime"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
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
                          "const": "Brian Murphy"
                        }
                      }
                    }
                  }
                },
                "expirationDateTime": {
                  "type": "string",
                  "enum": [
                    "2043-07-15T14:00:00Z"
                  ]
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation with expiration date to group with different roles
    Given user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource           | <resource>               |
      | space              | Personal                 |
      | sharee             | grp1                     |
      | shareType          | group                    |
      | permissionsRole    | <permissions-role>       |
      | expirationDateTime | 2043-07-15T14:00:00.000Z |
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
                "createdDateTime",
                "id",
                "roles",
                "grantedToV2",
                "expirationDateTime"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%group_id_pattern%$"
                        },
                        "displayName": {
                          "const": "grp1"
                        }
                      }
                    }
                  }
                },
                "expirationDateTime": {
                  "type": "string",
                  "enum": [
                    "2043-07-15T14:00:00Z"
                  ]
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |

  @issue-7962
  Scenario Outline: send share invitation to disabled user
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And the user "Admin" has disabled user "Brian"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
                "createdDateTime",
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
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
                          "const": "Brian Murphy"
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
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send sharing invitation to a deleted group with different roles
    Given user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And group "grp1" has been deleted
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "itemNotFound: not found"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation to deleted user
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And the user "Admin" has deleted a user "Brian"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
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
                "type": "string",
                "pattern": "invalidRequest"
              },
              "message": {
                "const": "itemNotFound: not found"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: try to send sharing invitation to multiple groups
    Given these users have been created with default attributes:
      | username |
      | Carol    |
      | Bob      |
    And group "grp1" has been created
    And group "grp2" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp2      |
      | Bob      | grp2      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1, grp2         |
      | shareType       | group, group       |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |
      | Manager          | FolderToShare  |


  Scenario Outline: try to send sharing invitation to user and group at once
    Given these users have been created with default attributes:
      | username |
      | Carol    |
      | Bob      |
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1, Bob          |
      | shareType       | group, user        |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send sharing invitation to non-existing group
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | nonExistentGroup   |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "itemNotFound: not found"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send sharing invitation to already shared group
    Given user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "409"
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
                "const": "nameAlreadyExists"
              },
              "message": {
                "type": "string",
                "pattern": "^error creating share: error: already exists:.*$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation to wrong user id
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following resource share invitation using the Graph API:
      | resource        | <resource>                           |
      | space           | Personal                             |
      | shareeId        | a4c0c83e-ae24-4870-93c3-fcaf2a2228f7 |
      | shareType       | user                                 |
      | permissionsRole | Viewer                               |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "itemNotFound: not found"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation with empty user id
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following resource share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | shareeId        |            |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "Key: 'DriveItemInvite.Recipients[0].ObjectId' Error:Field validation for 'ObjectId' failed on the 'ne' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation to user with wrong recipient type
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following resource share invitation using the Graph API:
      | resource        | <resource>     |
      | space           | Personal       |
      | sharee          | Brian          |
      | shareType       | wrongShareType |
      | permissionsRole | Viewer         |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "Key: 'DriveItemInvite.Recipients[0].LibreGraphRecipientType' Error:Field validation for 'LibreGraphRecipientType' failed on the 'oneof' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation to group with wrong recipient type
    Given user "Carol" has been created with default attributes
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    When user "Alice" tries to send the following resource share invitation using the Graph API:
      | resource        | <resource>     |
      | space           | Personal       |
      | sharee          | grp1           |
      | shareType       | wrongShareType |
      | permissionsRole | Viewer         |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "Key: 'DriveItemInvite.Recipients[0].LibreGraphRecipientType' Error:Field validation for 'LibreGraphRecipientType' failed on the 'oneof' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation to user with empty recipient type
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following resource share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       |            |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "Key: 'DriveItemInvite.Recipients[0].LibreGraphRecipientType' Error:Field validation for 'LibreGraphRecipientType' failed on the 'oneof' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation to group with empty recipient type
    Given user "Carol" has been created with default attributes
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    When user "Alice" tries to send the following resource share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | grp1       |
      | shareType       |            |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "Key: 'DriveItemInvite.Recipients[0].LibreGraphRecipientType' Error:Field validation for 'LibreGraphRecipientType' failed on the 'oneof' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: try to share a resource with invalid roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "role not applicable to this resource"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Manager          | /textfile1.txt |
      | Space Viewer     | /textfile1.txt |
      | Space Editor     | /textfile1.txt |
      | Manager          | FolderToShare  |
      | Space Viewer     | FolderToShare  |
      | Space Editor     | FolderToShare  |


  Scenario Outline: try to share a file with invalid roles
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
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
                "const": "invalidRequest"
              },
              "message": {
                "const": "role not applicable to this resource"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Editor           |
      | Uploader         |


  Scenario Outline: send share invitation to already shared user
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    When user "Alice" tries to send the following resource share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "409"
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
                "const": "nameAlreadyExists"
              },
              "message": {
                "type": "string",
                "pattern": "^error creating share: error: already exists:.*$"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation for project space resource to user with different roles (permissions endpoint)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "textfile1.txt"
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And user "Brian" should have a share "<resource>" shared by user "Alice" from space "NewSpace"
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
                "createdDateTime",
                "grantedToV2",
                "roles"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
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
      | permissions-role | resource      |
      | Viewer           | textfile1.txt |
      | File Editor      | textfile1.txt |
      | Viewer           | FolderToShare |
      | Editor           | FolderToShare |
      | Uploader         | FolderToShare |


  Scenario Outline: send share invitation for project space resource to group with different roles (permissions endpoint)
    Given using spaces DAV path
    And user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "textfile1.txt"
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And user "Brian" should have a share "<resource>" shared by user "Alice" from space "NewSpace"
    And user "Carol" should have a share "<resource>" shared by user "Alice" from space "NewSpace"
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
                "createdDateTime",
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "createdDateTime": { "format": "date-time" },
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%group_id_pattern%$"
                        },
                        "displayName": {
                          "const": "grp1"
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
      | permissions-role | resource      |
      | Viewer           | textfile1.txt |
      | File Editor      | textfile1.txt |
      | Viewer           | FolderToShare |
      | Editor           | FolderToShare |
      | Uploader         | FolderToShare |


  Scenario Outline: try to send share invitation with different re-sharing permissions
    Given group "grp1" has been created
    And user "Alice" has created folder "FolderToShare"
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource          | textfile1.txt        |
      | space             | Personal             |
      | sharee            | grp1                 |
      | shareType         | group                |
      | permissionsAction | <permissions-action> |
    Then the HTTP status code should be "400"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource          | FolderToShare        |
      | space             | Personal             |
      | sharee            | Brian                |
      | shareType         | user                 |
      | permissionsAction | <permissions-action> |
    Then the HTTP status code should be "400"
    Examples:
      | permissions-action |
      | permissions/create |
      | permissions/update |
      | permissions/delete |
      | permissions/deny   |


  Scenario: share a file to user and group having same name (Personal space)
    Given user "Carol" has been created with default attributes
    And group "Brian" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Carol    | Brian     |
    And user "Alice" has uploaded file with content "lorem" to "textfile.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    Then the HTTP status code should be "200"
    And user "Brian" should have a share "textfile.txt" shared by user "Alice" from space "Personal"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | group        |
      | permissionsRole | Viewer       |
    Then the HTTP status code should be "200"
    And user "Carol" should have a share "textfile.txt" shared by user "Alice" from space "Personal"


  Scenario: share a file to group containing special characters in name (Personal space)
    Given group "?\?@#%@;" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | ?\?@#%@;  |
    And user "Alice" has uploaded file with content "lorem" to "textfile.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | ?\?@#%@;     |
      | shareType       | group        |
      | permissionsRole | Viewer       |
    Then the HTTP status code should be "200"
    And user "Brian" should have a share "textfile.txt" shared by user "Alice" from space "Personal"


  Scenario: share a file to user and group having same name (Project space)
    Given using spaces DAV path
    And user "Carol" has been created with default attributes
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "lorem" to "textfile.txt"
    And group "Brian" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Carol    | Brian     |
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile.txt |
      | space           | NewSpace     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    Then the HTTP status code should be "200"
    And user "Brian" should have a share "textfile.txt" shared by user "Alice" from space "NewSpace"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile.txt |
      | space           | NewSpace     |
      | sharee          | Brian        |
      | shareType       | group        |
      | permissionsRole | Viewer       |
    Then the HTTP status code should be "200"
    And user "Carol" should have a share "textfile.txt" shared by user "Alice" from space "NewSpace"


  Scenario: share a file to group containing special characters in name (Project space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "lorem" to "textfile.txt"
    And group "?\?@#%@;" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | ?\?@#%@;  |
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile.txt |
      | space           | NewSpace     |
      | sharee          | ?\?@#%@;     |
      | shareType       | group        |
      | permissionsRole | Viewer       |
    Then the HTTP status code should be "200"
    And user "Brian" should have a share "textfile.txt" shared by user "Alice" from space "NewSpace"

  @env-config
  Scenario Outline: resource shared with denied permission role should not be visible when the sharee lists all drives (Personal space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "personal space" to "lorem.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Denied     |
    When user "Brian" lists all spaces via the Graph API
    Then the HTTP status code should be "200"
    And the json response should not contain the following shares:
      | <resource> |
    Examples:
      | resource      |
      | FolderToShare |
      | lorem.txt     |

  @env-config
  Scenario Outline: resource shared with denied permission role should not be visible when the sharee lists all drives (Project space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "lorem" to "lorem.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource> |
      | space           | NewSpace   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Denied     |
    When user "Brian" lists all spaces via the Graph API
    Then the HTTP status code should be "200"
    And the json response should not contain the following shares:
      | <resource> |
    Examples:
      | resource      |
      | FolderToShare |
      | lorem.txt     |

  @env-config
  Scenario Outline: try to share resource after disabling the role (Personal Space)
    Given the administrator has disabled the permissions role "<permissions-role>"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "some content" to "textfile.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
              "code": { "const": "invalidRequest" },
              "innererror": {
                "type": "object",
                "required": ["date", "request-id"]
              },
              "message": { "const": "Key: 'DriveItemInvite.Roles' Error:Field validation for 'Roles' failed on the 'available_role' tag" }
            }
          }
        }
      }
      """
    Examples:
      | resource      | permissions-role |
      | folderToShare | Viewer           |
      | folderToShare | Editor           |
      | folderToShare | Uploader         |
      | textfile.txt  | Viewer           |
      | textfile.txt  | File Editor      |

  @env-config
  Scenario Outline: try to share resource after disabling the role (Project Space)
    Given the administrator has disabled the permissions role "<permissions-role>"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
              "code": { "const": "invalidRequest" },
              "innererror": {
                "type": "object",
                "required": ["date", "request-id"]
              },
              "message": { "const": "Key: 'DriveItemInvite.Roles' Error:Field validation for 'Roles' failed on the 'available_role' tag" }
            }
          }
        }
      }
      """
    Examples:
      | resource      | permissions-role |
      | folderToShare | Viewer           |
      | folderToShare | Editor           |
      | folderToShare | Uploader         |
      | textfile.txt  | Viewer           |
      | textfile.txt  | File Editor      |
