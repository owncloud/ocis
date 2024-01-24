Feature: Send a sharing invitations
  As the owner of a resource
  I want to be able to send invitations to other users
  So that they can have access to it

  https://owncloud.dev/libre-graph-api/#/drives.permissions/Invite

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: send share invitation to user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Brian" the space Shares should contain these entries:
      | <path> |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
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
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
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
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: send share invitation to group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Brian" the space Shares should contain these entries:
      | <path> |
    And for user "Carol" the space Shares should contain these entries:
      | <path> |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
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
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
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
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: send share invitation for a file to user with different permissions
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType      | file                |
      | resource          | textfile1.txt       |
      | space             | Personal            |
      | sharee            | Brian               |
      | shareType         | user                |
      | permissionsAction | <permissionsAction> |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissionsAction>$"
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
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
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
      | permissionsAction  |
      | permissions/create |
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
      | permissions/update |
      | permissions/delete |
      | deleted/delete     |
      | permissions/deny   |


  Scenario Outline: send share invitation for a folder to user with different permissions
    Given user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType      | folder              |
      | resource          | FolderToShare       |
      | space             | Personal            |
      | sharee            | Brian               |
      | shareType         | user                |
      | permissionsAction | <permissionsAction> |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissionsAction>$"
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
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
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
      | permissionsAction  |
      | permissions/create |
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
      | permissions/update |
      | standard/delete    |
      | permissions/delete |
      | deleted/delete     |
      | permissions/deny   |


  Scenario Outline: send share invitation for a file to group with different permissions
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType      | file                |
      | resource          | textfile1.txt       |
      | space             | Personal            |
      | sharee            | grp1                |
      | shareType         | group               |
      | permissionsAction | <permissionsAction> |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissionsAction>$"
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
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
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
      | permissionsAction  |
      | permissions/create |
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
      | permissions/update |
      | permissions/delete |
      | deleted/delete     |
      | permissions/deny   |


  Scenario Outline: send share invitation for a folder to group with different permissions
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType      | folder              |
      | resource          | FolderToShare       |
      | space             | Personal            |
      | sharee            | grp1                |
      | shareType         | group               |
      | permissionsAction | <permissionsAction> |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissionsAction>$"
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
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
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
      | permissionsAction  |
      | permissions/create |
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
      | permissions/update |
      | standard/delete    |
      | permissions/delete |
      | deleted/delete     |
      | permissions/deny   |


  Scenario Outline: send share invitation with expiration date to user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>          |
      | resource        | <path>                   |
      | space           | Personal                 |
      | sharee          | Brian                    |
      | shareType       | user                     |
      | permissionsRole | <permissions-role>       |
      | expireDate      | 2043-07-15T14:00:00.000Z |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2",
                "expirationDateTime"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
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
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
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
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: send share invitation with expiration date to group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>          |
      | resource        | <path>                   |
      | space           | Personal                 |
      | sharee          | grp1                     |
      | shareType       | group                    |
      | permissionsRole | <permissions-role>       |
      | expireDate      | 2043-07-15T14:00:00.000Z |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2",
                "expirationDateTime"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
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
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
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
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |

  @issue-7962
  Scenario Outline: send share invitation to disabled user
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And the user "Admin" has disabled user "Brian"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
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
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
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
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
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
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: send sharing invitation to a deleted group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And the administrator has deleted group "grp1"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
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
                "type": "string",
                "enum": [
                  "generalException"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "itemNotFound: not found"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: send share invitation to deleted user
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And the user "Admin" has deleted a user "Brian"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
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
                "pattern": "generalException"
              },
              "message": {
                "type": "string",
                "enum": [
                  "itemNotFound: not found"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: try to send sharing invitation to multiple groups
    Given these users have been created with default attributes and without skeleton files:
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
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
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
                "type": "string",
                "enum": [
                  "invalidRequest"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Co Owner         | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |
      | Manager          | folder        | FolderToShare  |


  Scenario Outline: try to send sharing invitation to user and group at once
    Given these users have been created with default attributes and without skeleton files:
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
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
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
                "type": "string",
                "enum": [
                  "invalidRequest"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: send sharing invitation to non-existing group
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
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
                "type": "string",
                "enum": [
                  "generalException"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "itemNotFound: not found"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: send sharing invitation to already shared group
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following share invitation:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType    | <resource-type>    |
      | resource        | <path>             |
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
                "type": "string",
                "enum": [
                  "nameAlreadyExists"
                ]
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
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |


  Scenario Outline: send share invitation to wrong user id with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following share invitation using the Graph API:
      | resourceType    | <resource-type>                      |
      | resource        | <path>                               |
      | space           | Personal                             |
      | shareeId        | a4c0c83e-ae24-4870-93c3-fcaf2a2228f7 |
      | shareType       | user                                 |
      | permissionsRole | <permissions-role>                   |
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
                "pattern": "generalException"
              },
              "message": {
                "type": "string",
                "enum": [
                  "itemNotFound: not found"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource-type | path           |
      | Viewer           | file          | /textfile1.txt |
      | File Editor      | file          | /textfile1.txt |
      | Viewer           | folder        | FolderToShare  |
      | Editor           | folder        | FolderToShare  |
      | Uploader         | folder        | FolderToShare  |
