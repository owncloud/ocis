@issue-10739
Feature: Send a drive invitations
  As the owner of a resource
  I want to be able to send drive invitations to other users
  So that they can have access to it

  https://owncloud.dev/libre-graph-api/#/drives.permissions/Invite

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: send share invitation for project space to user with different roles (permissions endpoint)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | NewSpace           |
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


  Scenario Outline: send share invitation for disabled project space to user with different roles (permissions endpoint)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Admin" has disabled a space "NewSpace"
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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


  Scenario Outline: send share invitation for deleted project space to user with different roles (permissions endpoint)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Admin" has disabled a space "NewSpace"
    And user "Admin" has deleted a space "NewSpace"
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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


  Scenario Outline: send share invitation for project space to group with different roles (permissions endpoint)
    Given using spaces DAV path
    And user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
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
                "grantedToV2",
                "roles",
                "id"
              ],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
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
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: send share invitation for disabled project space to group with different roles (permissions endpoint)
    Given using spaces DAV path
    And user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Admin" has disabled a space "NewSpace"
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
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


  Scenario Outline: send share invitation for deleted project space to group with different roles (permissions endpoint)
    Given using spaces DAV path
    And user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Admin" has disabled a space "NewSpace"
    And user "Admin" has deleted a space "NewSpace"
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
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


  Scenario: send share invitation to user for deleted file
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And we save it into "FILEID"
    And user "Alice" has deleted file "textfile1.txt"
    When user "Alice" sends the following share invitation with file-id "<<FILEID>>" using the Graph API:
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    Then the HTTP status code should be "404"
    And user "Brian" should not have a share "textfile1.txt" shared by user "Alice" from space "Personal"
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


  Scenario: send share invitation to group for deleted file
    Given user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And we save it into "FILEID"
    And user "Alice" has deleted file "textfile1.txt"
    When user "Alice" sends the following share invitation with file-id "<<FILEID>>" using the Graph API:
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Viewer   |
    Then the HTTP status code should be "404"
    And user "Brian" should not have a share "textfile1.txt" shared by user "Alice" from space "Personal"
    And user "Carol" should not have a share "textfile1.txt" shared by user "Alice" from space "Personal"
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

  @issue-8494
  Scenario Outline: try to send share invitation for personal space to user with different roles (permissions endpoint)
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
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

  @issue-8495
  Scenario Outline: try to share Shares space with a user (permissions endpoint)
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | Shares             |
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


  Scenario Outline: invite user to a project space with different roles using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace           |
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


  Scenario Outline: invite group to project space with different roles using root endpoint
    Given using spaces DAV path
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
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
                "grantedToV2",
                "roles"
              ],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
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


  Scenario Outline: try to invite multiple users to project space with different roles using root endpoint
    Given using spaces DAV path
    And user "Carol" has been created with default attributes
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | Brian, Carol       |
      | shareType       | user, user         |
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
                "const": "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
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


  Scenario Outline: try to invite one existing user and one non-existing user at once to project space with different roles using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace            |
      | sharee          | Brian, non-existent |
      | shareType       | user, user          |
      | permissionsRole | <permissions-role>  |
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
                "const": "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
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


  Scenario Outline: try to invite multiple groups at once to project space with different roles using root endpoint
    Given using spaces DAV path
    And user "Carol" has been created with default attributes
    And group "grp1" has been created
    And group "grp2" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp2      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1, grp2         |
      | shareType       | group, group       |
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
                "const": "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
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


  Scenario Outline: try to invite one existing group and one non-existing group to project space with different roles using root endpoint
    Given using spaces DAV path
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1, grp2         |
      | shareType       | group, group       |
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
                "const": "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
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


  Scenario Outline: try to invite user and group at once to project space with different roles using root endpoint
    Given using spaces DAV path
    And user "Carol" has been created with default attributes
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | Carol, grp2        |
      | shareType       | user, group        |
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
                "const": "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
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


  Scenario Outline: try to invite user to personal drive with different roles using root endpoint
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
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
      | Viewer           |
      | File Editor      |
      | Viewer           |
      | Editor           |
      | Uploader         |


  Scenario Outline: try to invite group to personal drive with different roles using root endpoint
    Given group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
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
      | Viewer           |
      | File Editor      |
      | Viewer           |
      | Editor           |
      | Uploader         |


  Scenario Outline: try to invite user to shares drive with different re-sharing permissions using root endpoint
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space             | Shares               |
      | sharee            | Brian                |
      | shareType         | user                 |
      | permissionsAction | <permissions-action> |
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
      | permissions-action |
      | permissions/create |
      | permissions/update |
      | permissions/delete |
      | permissions/deny   |


  Scenario Outline: try to invite group to shares drive with different re-sharing permissions using root endpoint
    Given group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space             | Shares               |
      | sharee            | grp1                 |
      | shareType         | group                |
      | permissionsAction | <permissions-action> |
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
      | permissions-action |
      | permissions/create |
      | permissions/update |
      | permissions/delete |
      | permissions/deny   |


  Scenario Outline: try to send a sharing invitation for the personal drive to an non-existent sharee using root endpoint
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | Personal           |
      | sharee          | non-existent       |
      | shareType       | <sharee-type>      |
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
      | permissions-role | sharee-type |
      | Viewer           | user        |
      | File Editor      | user        |
      | Viewer           | user        |
      | Editor           | user        |
      | Uploader         | user        |
      | Viewer           | group       |
      | File Editor      | group       |
      | Viewer           | group       |
      | Editor           | group       |
      | Uploader         | group       |


  Scenario Outline: try to send a sharing invitation for the personal drive with an empty sharee using root endpoint
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | Personal           |
      | sharee          |                    |
      | shareType       | <sharee-type>      |
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
                "const": "Key: 'DriveItemInvite.Recipients[0].ObjectId' Error:Field validation for 'ObjectId' failed on the 'ne' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | sharee-type |
      | Viewer           | user        |
      | File Editor      | user        |
      | Viewer           | user        |
      | Editor           | user        |
      | Uploader         | user        |
      | Viewer           | group       |
      | File Editor      | group       |
      | Viewer           | group       |
      | Editor           | group       |
      | Uploader         | group       |


  Scenario Outline: try to send a sharing invitation for the project drive to an non-existent sharee using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          | non-existent       |
      | shareType       | <sharee-type>      |
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
                "const": "itemNotFound: not found"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | sharee-type |
      | Space Viewer     | user        |
      | Space Editor     | user        |
      | Manager          | user        |
      | Space Viewer     | group       |
      | Space Editor     | group       |
      | Manager          | group       |


  Scenario Outline: try to send a sharing invitation for the project drive with empty sharee using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace           |
      | sharee          |                    |
      | shareType       | <sharee-type>      |
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
                "const": "Key: 'DriveItemInvite.Recipients[0].ObjectId' Error:Field validation for 'ObjectId' failed on the 'ne' tag"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | sharee-type |
      | Space Viewer     | user        |
      | Space Editor     | user        |
      | Manager          | user        |
      | Space Viewer     | group       |
      | Space Editor     | group       |
      | Manager          | group       |

  @env-config
  Scenario Outline: try to invite to the project space after disabling the role
    Given the administrator has disabled the permissions role "<permissions-role>"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
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
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @issue-9303
  Scenario: try to invite user to project space with permissions role Secure Viewer using root endpoint
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation using root endpoint of the Graph API:
      | space           | NewSpace      |
      | sharee          | Alice         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
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
                "const": "role not applicable to this resource"
              }
            }
          }
        }
      }
      """

  @issue-9303
  Scenario: try to invite user to project space with permissions role Secure Viewer (permissions endpoint)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" sends the following space share invitation using permissions endpoint of the Graph API:
      | space           | NewSpace      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
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
                "const": "role not applicable to this resource"
              }
            }
          }
        }
      }
      """
