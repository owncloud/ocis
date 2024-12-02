@ocm
Feature: an user shares resources using ScienceMesh application
  As a user
  I want to share resources between different ocis instances

  Background:
    Given user "Alice" has been created with default attributes
    And using server "REMOTE"
    And user "Brian" has been created with default attributes

  @issue-9534
  Scenario Outline: local user shares resources to federation user
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    When user "Alice" sends the following resource share invitation to federated user using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "200"
    When using server "REMOTE"
    And user "Brian" lists the shares shared with him without retry using the Graph API
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
                "@UI.Hidden": {
                  "type": "boolean",
                  "enum": [
                    false
                  ]
                },
                "@client.synchronize": {
                  "type": "boolean",
                  "enum": [
                    false
                  ]
                },
                "createdBy": {
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
                          "const": "Alice Hansen"
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%federated_user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "name": {
                  "const": "<resource>"
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | resource      |
      | folderToShare |
      | textfile.txt  |


  Scenario: local user shares resources from project space to federation user
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    When user "Alice" sends the following resource share invitation to federated user using the Graph API:
      | resource        | folderToShare |
      | space           | projectSpace  |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    Then the HTTP status code should be "200"
    When using server "REMOTE"
    And user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [ "value" ],
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
                "@UI.Hidden": {
                  "const": false
                },
                "@client.synchronize": {
                  "const": false
                },
                "createdBy": {
                  "type": "object",
                  "required": [ "user" ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [ "displayName", "id" ],
                      "properties": {
                        "displayName": {
                          "const": "Alice Hansen"
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%federated_user_id_pattern%$"
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

  @issue-9534
  Scenario Outline: federation user shares resource to local user after accepting invitation
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And user "Brian" has created folder "folderToShare"
    And user "Brian" has uploaded file with content "ocm test" to "/textfile.txt"
    When user "Brian" sends the following resource share invitation to federated user using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Alice      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "200"
    When using server "LOCAL"
    And user "Alice" lists the shares shared with her without retry using the Graph API
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
                "@UI.Hidden": {
                  "type": "boolean",
                  "enum": [
                    false
                  ]
                },
                "@client.synchronize": {
                  "type": "boolean",
                  "enum": [
                    false
                  ]
                },
                "createdBy": {
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
                          "pattern": "^%federated_user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "name": {
                  "const": "<resource>"
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | resource      |
      | folderToShare |
      | textfile.txt  |

  @issue-10051
  Scenario Outline: try to add federated user as a member of a project space (permissions endpoint)
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "alice's space" with the default quota using the Graph API
    When user "Alice" tries to send the following space share invitation to federated user using permissions endpoint of the Graph API:
      | space           | alice's space      |
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
            "required": [
              "code",
              "innererror",
              "message"
            ],
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
                "const": "federated user can not become a space member"
              }
            }
          }
        }
      }
      """
    And using server "REMOTE"
    And the user "Brian" should not have a space called "alice's space"
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @issue-10051
  Scenario Outline: try to add federated user as a member of a project space (root endpoint)
    Given using server "LOCAL"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "alice's space" with the default quota using the Graph API
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    When user "Alice" tries to send the following space share invitation to federated user using root endpoint of the Graph API:
      | space           | alice's space      |
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
            "required": [
              "code",
              "innererror",
              "message"
            ],
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
                "const": "federated user can not become a space member"
              }
            }
          }
        }
      }
      """
    And using server "REMOTE"
    And the user "Brian" should not have a space called "alice's space"
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @issue-9908
  Scenario: sharer lists the shares shared to a federated user
    Given using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "textfile.txt" with the following data:
      """
      {
        "type": "object",
        "required": [
          "parentReference",
          "permissions",
          "name",
          "size"
        ],
        "properties": {
          "parentReference": {
            "type": "object",
            "required": [
              "driveId",
              "driveType",
              "path",
              "name",
              "id"
            ],
            "properties": {
              "driveId": {
                "type": "string",
                "pattern": "^%space_id_pattern%$"
              },
              "driveType": {
                "const": "personal"
              },
              "path": {
                "const": "/"
              },
              "name": {
                "const": "/"
              },
              "id": {
                "type": "string",
                "pattern": "^%file_id_pattern%$"
              }
            }
          },
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "@libre.graph.userType": {
                          "const": "Federated"
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%federated_user_id_pattern%$"
                        },
                        "displayName": {
                          "const": "Brian Murphy"
                        }
                      }
                    }
                  }
                },
                "id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
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
          },
          "name": {
            "const": "textfile.txt"
          },
          "size": {
            "const": 8
          }
        }
      }
      """

  @issue-9898
  Scenario: user lists permissions of a resource shared to a federated user
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And using server "LOCAL"
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
                            "@libre.graph.userType": {
                              "type": "string",
                              "const": "Federated"
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%federated_user_id_pattern%$"
                            },
                            "displayName": {
                              "const": "Brian Murphy"
                            }
                          }
                        }
                      }
                    },
                    "id": {
                      "type": "string",
                      "pattern": "^%user_id_pattern%$"
                    },
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
                                "@libre.graph.userType": {
                                  "type": "string",
                                  "const": "Member"
                                },
                                "id": {
                                  "type": "string",
                                  "pattern": "^%user_id_pattern%$"
                                },
                                "displayName": {
                                  "const": "Alice Hansen"
                                }
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
                      "items": {
                        "type": "string",
                        "pattern": "^%role_id_pattern%$"
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

  @issue-10222 @issue-10495
  Scenario: local user lists multiple federation shares
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And user "Brian" has uploaded file "filesForUpload/testavatar.jpg" to "testavatar.jpg"
    And user "Brian" has created folder "folderToShare"
    And user "Brian" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Alice         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Brian" has sent the following resource share invitation to federated user:
      | resource        | testavatar.jpg |
      | space           | Personal       |
      | sharee          | Alice          |
      | shareType       | user           |
      | permissionsRole | Viewer         |
    And using server "LOCAL"
    When user "Alice" lists the shares shared with her using the Graph API
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
            "minItems": 2,
            "maxItems": 2,
            "uniqueItems": true,
            "items": {
              "oneOf":[
                {
                  "type": "object",
                  "required": [
                    "@UI.Hidden",
                    "@client.synchronize",
                    "createdBy",
                    "eTag",
                    "folder",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "parentReference",
                    "remoteItem"
                  ],
                  "properties": {
                    "name": {
                      "const": "folderToShare"
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "@UI.Hidden",
                    "@client.synchronize",
                    "createdBy",
                    "eTag",
                    "file",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "parentReference",
                    "remoteItem",
                    "size"
                  ],
                  "properties": {
                    "name": {
                      "const": "testavatar.jpg"
                    },
                    "file": {
                      "type": "object",
                      "required": [
                        "mimeType"
                      ],
                      "properties": {
                        "mimeType": {
                          "const": "image/jpeg"
                        }
                      }
                    },
                    "remoteItem": {
                      "type": "object",
                      "required": [
                        "createdBy",
                        "eTag",
                        "file",
                        "id",
                        "lastModifiedDateTime",
                        "name",
                        "permissions",
                        "size"
                      ],
                      "properties": {
                        "file": {
                          "type": "object",
                          "required": [
                            "mimeType"
                          ],
                          "properties": {
                            "mimeType": {
                              "const": "image/jpeg"
                            }
                          }
                        }
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

  @issue-10285 @issue-10536 @issue-10305
  Scenario: federation user uploads file to a federated shared folder via TUS
    Given using spaces DAV path
    And using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And user "Alice" has created a folder "FOLDER" in space "Personal"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    When using server "REMOTE"
    And user "Brian" uploads a file with content "lorem" to "file.txt" inside federated share "FOLDER" via TUS using the WebDAV API
    Then for user "Brian" the content of file "file.txt" of federated share "FOLDER" should be "lorem"

  @issue-10285 @issue-10536
  Scenario: local user uploads file to a federated shared folder via TUS
    Given using spaces DAV path
    And using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And user "Brian" has created a folder "FOLDER" in space "Personal"
    And user "Brian" has sent the following resource share invitation to federated user:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Alice    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    When using server "LOCAL"
    And user "Alice" uploads a file with content "lorem" to "file.txt" inside federated share "FOLDER" via TUS using the WebDAV API
    Then for user "Alice" the content of file "file.txt" of federated share "FOLDER" should be "lorem"

  @issue-10495
  Scenario: local user downloads thumbnail preview of a federated shared image
    Given using spaces DAV path
    And using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And user "Brian" has uploaded file "filesForUpload/testavatar.jpg" to "testavatar.jpg"
    And user "Brian" has sent the following resource share invitation to federated user:
      | resource        | testavatar.jpg |
      | space           | Personal       |
      | sharee          | Alice          |
      | shareType       | user           |
      | permissionsRole | Viewer         |
    And using server "LOCAL"
    When user "Alice" downloads the preview of federated share image "testavatar.jpg" with width "32" and height "32" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded image should be "32" pixels wide and "32" pixels high
    And the downloaded preview content should match with "thumbnail.png" fixtures preview content

  @issue-10358
  Scenario: user edits content of a federated share file
    Given using spaces DAV path
    And using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | File Editor  |
    And using server "REMOTE"
    And for user "Brian" the content of file "textfile.txt" of federated share "textfile.txt" should be "ocm test"
    When user "Brian" updates the content of federated share "textfile.txt" with "this is a new content" using the WebDAV API
    Then the HTTP status code should be "204"
    And for user "Brian" the content of file "textfile.txt" of federated share "textfile.txt" should be "this is a new content"
    And using server "LOCAL"
    And for user "Alice" the content of the file "textfile.txt" of the space "Personal" should be "this is a new content"
