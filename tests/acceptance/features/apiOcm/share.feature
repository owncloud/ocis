@ocm
Feature: an user shares resources using ScienceMesh application
  As a user
  I want to share resources between different ocis instances

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has accepted invitation

  @issue-9534 @issue-11054
  Scenario Outline: local user shares a folder to federation user
    Given using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    When user "Alice" sends the following resource share invitation to federated user using the Graph API:
      | resource        | folderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    When using server "REMOTE"
    And user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
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
                "eTag",
                "folder",
                "id",
                "lastModifiedDateTime",
                "name",
                "parentReference",
                "remoteItem"
              ],
              "properties": {
                "lastModifiedDateTime": { "format": "date-time" },
                "@UI.Hidden": { "const": false },
                "@client.synchronize": { "const": false },
                "createdBy": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id", "@libre.graph.userType"],
                      "properties": {
                        "displayName": { "const": "Alice Hansen" },
                        "id": { "pattern": "^%federated_user_id_pattern%$" },
                        "@libre.graph.userType": { "const": "Federated" }
                      }
                    }
                  }
                },
                "eTag": { "pattern": "%etag_pattern%" },
                "folder": { "const": {} },
                "id": { "pattern": "^%file_id_pattern%$" },
                "name": { "const": "folderToShare" },
                "parentReference": {
                  "type": "object",
                  "required": ["driveId", "driveType", "id"],
                  "properties": {
                    "driveId": { "pattern": "^%space_id_pattern%$" },
                    "driveType": { "const": "virtual" },
                    "id": { "pattern": "^%file_id_pattern%$" }
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "createdBy",
                    "eTag",
                    "folder",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "permissions"
                  ],
                  "properties": {
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["id", "displayName", "@libre.graph.userType"],
                          "properties": {
                            "id": { "pattern": "^%federated_user_id_pattern%$" },
                            "displayName": { "const": "Alice Hansen" },
                            "@libre.graph.userType": { "const": "Federated" }
                          }
                        }
                      }
                    },
                    "eTag": { "pattern": "%etag_pattern%" },
                    "folder": { "const": {} },
                    "id": { "pattern": "^%federated_file_id_pattern%$" },
                    "name": { "const": "folderToShare" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": [
                          "createdDateTime",
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "createdDateTime": { "format": "date-time" },
                          "grantedToV2": {
                            "type": "object",
                            "required": ["user"],
                            "properties": {
                              "user": {
                                "type": "object",
                                "required": ["displayName", "id", "@libre.graph.userType"],
                                "properties": {
                                  "displayName": { "const": "Brian Murphy" },
                                  "id": { "pattern": "^%user_id_pattern%$" },
                                  "@libre.graph.userType": { "const": "Member" }
                                }
                              }
                            }
                          },
                          "id": { "pattern": "^%uuidv4_pattern%$" },
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
                                    "required": ["displayName", "id", "@libre.graph.userType"],
                                    "properties": {
                                      "displayName": { "const": "Alice Hansen" },
                                      "id": { "pattern": "^%federated_user_id_pattern%$" },
                                      "@libre.graph.userType": { "const": "Federated" }
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
                              "pattern": "^%role_id_pattern%$"
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
      }
      """
    Examples:
      | permissions-role |
      | Viewer           |
      | Editor           |
      | Uploader         |

  @issue-9534 @issue-11054
  Scenario Outline: local user shares a file to federation user
    Given using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "textfile.txt"
    When user "Alice" sends the following resource share invitation to federated user using the Graph API:
      | resource        | textfile.txt       |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    When using server "REMOTE"
    And user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
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
                "eTag",
                "file",
                "id",
                "lastModifiedDateTime",
                "name",
                "parentReference",
                "remoteItem"
              ],
              "properties": {
                "lastModifiedDateTime": { "format": "date-time" },
                "@UI.Hidden": { "const": false },
                "@client.synchronize": { "const": false },
                "createdBy": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id", "@libre.graph.userType"],
                      "properties": {
                        "displayName": { "const": "Alice Hansen" },
                        "id": { "pattern": "^%federated_user_id_pattern%$" },
                        "@libre.graph.userType": { "const": "Federated" }
                      }
                    }
                  }
                },
                "eTag": { "pattern": "%etag_pattern%" },
                "file": {
                  "type": "object",
                  "required": ["mimeType"],
                  "properties": {
                    "mimeType": { "const": "text/plain" }
                  }
                },
                "id": { "pattern": "^%file_id_pattern%$" },
                "name": { "const": "textfile.txt" },
                "parentReference": {
                  "type": "object",
                  "required": ["driveId", "driveType", "id"],
                  "properties": {
                    "driveId": { "pattern": "^%space_id_pattern%$" },
                    "driveType": { "const": "virtual" },
                    "id": { "pattern": "^%file_id_pattern%$" }
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
                    "permissions"
                  ],
                  "properties": {
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["id", "displayName", "@libre.graph.userType"],
                          "properties": {
                            "id": { "pattern": "^%federated_user_id_pattern%$" },
                            "displayName": { "const": "Alice Hansen" },
                            "@libre.graph.userType": { "const": "Federated" }
                          }
                        }
                      }
                    },
                    "eTag": { "pattern": "%etag_pattern%" },
                    "file": {
                      "type": "object",
                      "required": ["mimeType"],
                      "properties": {
                        "mimeType": { "const": "text/plain" }
                      }
                    },
                    "id": { "pattern": "^%federated_file_id_pattern%$" },
                    "name": { "const": "textfile.txt" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": [
                          "createdDateTime",
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "createdDateTime": { "format": "date-time" },
                          "grantedToV2": {
                            "type": "object",
                            "required": ["user"],
                            "properties": {
                              "user": {
                                "type": "object",
                                "required": ["displayName", "id", "@libre.graph.userType"],
                                "properties": {
                                  "displayName": { "const": "Brian Murphy" },
                                  "id": { "pattern": "^%user_id_pattern%$" },
                                  "@libre.graph.userType": { "const": "Member" }
                                }
                              }
                            }
                          },
                          "id": { "pattern": "^%uuidv4_pattern%$" },
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
                                    "required": ["displayName", "id", "@libre.graph.userType"],
                                    "properties": {
                                      "displayName": { "const": "Alice Hansen" },
                                      "id": { "pattern": "^%federated_user_id_pattern%$" },
                                      "@libre.graph.userType": { "const": "Federated" }
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
                              "pattern": "^%role_id_pattern%$"
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
      }
      """
    Examples:
      | permissions-role |
      | Viewer           |
      | File Editor      |

  @issue-11054
  Scenario Outline: local user shares a folder from project space to federation user
    Given using server "LOCAL"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    When user "Alice" sends the following resource share invitation to federated user using the Graph API:
      | resource        | folderToShare      |
      | space           | projectSpace       |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    When using server "REMOTE"
    And user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
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
                "eTag",
                "folder",
                "id",
                "lastModifiedDateTime",
                "name",
                "parentReference",
                "remoteItem"
              ],
              "properties": {
                "lastModifiedDateTime": { "format": "date-time" },
                "@UI.Hidden": { "const": false },
                "@client.synchronize": { "const": false },
                "createdBy": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id", "@libre.graph.userType"],
                      "properties": {
                        "displayName": { "const": "Alice Hansen" },
                        "id": { "pattern": "^%federated_user_id_pattern%$" },
                        "@libre.graph.userType": { "const": "Federated" }
                      }
                    }
                  }
                },
                "eTag": { "pattern": "%etag_pattern%" },
                "folder": { "const": {} },
                "id": { "pattern": "^%file_id_pattern%$" },
                "name": { "const": "folderToShare" },
                "parentReference": {
                  "type": "object",
                  "required": ["driveId", "driveType", "id"],
                  "properties": {
                    "driveId": { "pattern": "^%space_id_pattern%$" },
                    "driveType": { "const": "virtual" },
                    "id": { "pattern": "^%file_id_pattern%$" }
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "createdBy",
                    "eTag",
                    "folder",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "permissions"
                  ],
                  "properties": {
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["id", "displayName", "@libre.graph.userType"],
                          "properties": {
                            "id": { "pattern": "^%federated_user_id_pattern%$" },
                            "displayName": { "const": "Alice Hansen" },
                            "@libre.graph.userType": { "const": "Federated" }
                          }
                        }
                      }
                    },
                    "eTag": { "pattern": "%etag_pattern%" },
                    "folder": { "const": {} },
                    "id": { "pattern": "^%federated_file_id_pattern%$" },
                    "name": { "const": "folderToShare" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": [
                          "createdDateTime",
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "createdDateTime": { "format": "date-time" },
                          "grantedToV2": {
                            "type": "object",
                            "required": ["user"],
                            "properties": {
                              "user": {
                                "type": "object",
                                "required": ["displayName", "id", "@libre.graph.userType"],
                                "properties": {
                                  "displayName": { "const": "Brian Murphy" },
                                  "id": { "pattern": "^%user_id_pattern%$" },
                                  "@libre.graph.userType": { "const": "Member" }
                                }
                              }
                            }
                          },
                          "id": { "pattern": "^%uuidv4_pattern%$" },
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
                                    "required": ["displayName", "id", "@libre.graph.userType"],
                                    "properties": {
                                      "displayName": { "const": "Alice Hansen" },
                                      "id": { "pattern": "^%federated_user_id_pattern%$" },
                                      "@libre.graph.userType": { "const": "Federated" }
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
                              "pattern": "^%role_id_pattern%$"
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
      }
      """
    Examples:
      | permissions-role |
      | Viewer           |
      | Editor           |
      | Uploader         |

  @issue-10051
  Scenario Outline: try to add federated user as a member of a project space (permissions endpoint)
    Given using server "LOCAL"
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
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "alice's space" with the default quota using the Graph API
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
          "eTag",
          "file",
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
          "eTag": {
            "type": "string",
            "pattern": "%etag_pattern%"
          },
          "file": {
            "type": "object",
            "required": ["mimeType"],
            "properties": {
              "mimeType": {
                "const": "text/plain"
              }
            }
          },
          "size": {
            "const": 8
          }
        }
      }
      """

  @issue-10222 @issue-10495
  Scenario: local user lists multiple federation shares
    Given user "Brian" has uploaded file "filesForUpload/testavatar.jpg" to "testavatar.jpg"
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
                    "lastModifiedDateTime": { "format": "date-time" },
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
                    "lastModifiedDateTime": { "format": "date-time" },
                    "name": {
                      "const": "testavatar.jpg"
                    },
                    "file": {
                      "type": "object",
                      "required": ["mimeType"],
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
                          "required": ["mimeType"],
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
    Given using server "LOCAL"
    And user "Alice" has created a folder "FOLDER" in space "Personal"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    When using server "REMOTE"
    And user "Brian" uploads a file with content "lorem" to "file.txt" inside federated share "FOLDER" via TUS using the WebDAV API
    And the HTTP status code should be "204"
    Then for user "Brian" the content of file "file.txt" of federated share "FOLDER" should be "lorem"

  @issue-10285 @issue-10536
  Scenario: local user uploads file to a federated shared folder via TUS
    Given user "Brian" has created a folder "FOLDER" in space "Personal"
    And user "Brian" has sent the following resource share invitation to federated user:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Alice    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    When using server "LOCAL"
    And user "Alice" uploads a file with content "lorem" to "file.txt" inside federated share "FOLDER" via TUS using the WebDAV API
    And the HTTP status code should be "204"
    Then for user "Alice" the content of file "file.txt" of federated share "FOLDER" should be "lorem"

  @issue-10495
  Scenario: local user downloads thumbnail preview of a federated shared image
    Given user "Brian" has uploaded file "filesForUpload/testavatar.jpg" to "testavatar.jpg"
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
    Given using server "LOCAL"
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

  @issue-10488
  Scenario: local user shares a folder copied from an already shared folder to federation user
    Given using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Editor        |
    And user "Alice" has copied folder "folderToShare" to "folderToShareCopy"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShareCopy |
      | space           | Personal          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | Viewer            |
    And using server "REMOTE"
    When user "Brian" lists the shares shared with him using the Graph API
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
                    "lastModifiedDateTime": { "format": "date-time" },
                    "name": {
                      "const": "folderToShare"
                    },
                    "remoteItem": {
                      "type": "object",
                      "required": ["permissions"],
                      "properties": {
                        "permissions": {
                          "type": "array",
                          "minItems": 1,
                          "maxItems": 1,
                          "items": {
                            "type": "object",
                            "properties": {
                              "@libre.graph.permissions.actions": {
                                "type": "null"
                              }
                            }
                          }
                        }
                      }
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
                    "folder",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "parentReference",
                    "remoteItem"
                  ],
                  "properties": {
                    "lastModifiedDateTime": { "format": "date-time" },
                    "name": {
                      "const": "folderToShareCopy"
                    },
                    "remoteItem": {
                      "type": "object",
                      "required": ["permissions"],
                      "properties": {
                        "permissions": {
                          "type": "array",
                          "minItems": 1,
                          "maxItems": 1,
                          "items": {
                            "type": "object",
                            "properties": {
                              "@libre.graph.permissions.actions": {
                                "type": "null"
                              }
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

  @issue-9926 @issue-11022
  Scenario: federated user updates a shared file after sharer has updated the role
    Given using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole    | File Editor              |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
      | space              | Personal                 |
      | resource           | textfile.txt             |
    And using server "REMOTE"
    When user "Brian" updates the content of federated share "textfile.txt" with "this is a new content" using the WebDAV API
    Then the HTTP status code should be "204"
    And using server "LOCAL"
    And the content of file "textfile.txt" for user "Alice" should be "this is a new content"

  @issue-10689
  Scenario: federation user lists all the spaces
    Given using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Editor        |
    And using server "REMOTE"
    When user "Brian" lists all available spaces via the Graph API
    Then the HTTP status code should be "200"
    And the json response should not contain a space with name "folderToShare"

  @issue-10213
  Scenario Outline: local user removes access of federated user from a resource
    Given using server "LOCAL"
    And user "Alice" has created a folder "FOLDER" in space "Personal"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | FOLDER            |
      | space           | Personal          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | <permissionsRole> |
    When user "Alice" removes the access of user "Brian" from resource "FOLDER" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And using server "REMOTE"
    And user "Brian" should not have a federated share "FOLDER" shared by user "Alice" from space "Personal"
    Examples:
      | permissionsRole |
      | Viewer          |
      | Uploader        |
      | Editor          |

  @issue-10272
  Scenario: federated user downloads shared resources as an archive
    Given using server "LOCAL"
    And user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And user "Alice" has created folder "imageFolder"
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "imageFolder/testavatar.png"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | imageFolder |
      | space           | Personal    |
      | sharee          | Brian       |
      | shareType       | user        |
      | permissionsRole | Viewer      |
    And using server "REMOTE"
    When user "Brian" downloads the archive of these items using the resource remoteItemIds
      | textfile.txt |
      | imageFolder  |
    Then the HTTP status code should be "200"
    And the downloaded zip archive should contain these files:
      | name                       | content   |
      | textfile.txt               | some data |
      | imageFolder                |           |
      | imageFolder/testavatar.png |           |

  @issue-11033
  Scenario: external sharee shouldn't be able to the access file when federated share expires
    Given using SharingNG
    And using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | File Editor  |
    When user "Alice" expires the last share of resource "textfile.txt" inside of the space "Personal"
    Then the HTTP status code should be "200"
    And using server "REMOTE"
    And user "Brian" should not have a federated share "textfile.txt" shared by user "Alice" from space "Personal"

  @issue-11033
  Scenario: external sharee shouldn't be able to the access folder when federated share expires
    Given using SharingNG
    And using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Editor        |
    When user "Alice" expires the last share of resource "folderToShare" inside of the space "Personal"
    Then the HTTP status code should be "200"
    And using server "REMOTE"
    And user "Brian" should not have a federated share "folderToShare" shared by user "Alice" from space "Personal"

  @issue-10719
  Scenario: federated user hides the file shared by local user
    Given using server "LOCAL"
    And user "Alice" has uploaded file with content "hello world" to "testfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | testfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And using server "REMOTE"
    When user "Brian" hides the federated share "testfile.txt" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["@UI.Hidden"],
        "properties": {
          "@UI.Hidden": { "const": true }
        }
      }
      """

  @issue-10719
  Scenario: federated user hides the folder shared by local user
    Given using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And using server "REMOTE"
    When user "Brian" hides the federated share "folderToShare" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["@UI.Hidden"],
        "properties": {
          "@UI.Hidden": { "const": true }
        }
      }
      """

  @issue-10582
  Scenario Outline: federation user creates folder inside shared folder (Personal Space)
    Given using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Editor        |
    And using server "REMOTE"
    When user "Brian" creates a folder "newFolder" inside federated share "folderToShare" using the WebDav API
    Then the HTTP status code should be "201"
    And using server "LOCAL"
    And as "Alice" folder "folderToShare/newFolder" should exist
    When user "Alice" requests "<dav-path>" with "PROPFIND" without retrying
    Then the HTTP status code should be "207"
    And as user "Alice" the PROPFIND response should contain a resource "newFolder" with these key and value pairs:
      | key     | value     |
      | oc:name | newFolder |
    Examples:
      | dav-path                            |
      | /webdav/folderToShare               |
      | /dav/files/%username%/folderToShare |
      | /dav/spaces/%spaceid%/folderToShare |

  @issue-10582
  Scenario: federation user creates folder inside shared folder (Project Space)
    Given using server "LOCAL"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | projectSpace  |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Editor        |
    And using server "REMOTE"
    When user "Brian" creates a folder "newFolder" inside federated share "folderToShare" using the WebDav API
    Then the HTTP status code should be "201"
    And using server "LOCAL"
    And for user "Alice" folder "folderToShare" of the space "projectSpace" should contain these entries:
      | newFolder |
    When user "Alice" sends PROPFIND request from the space "projectSpace" to the resource "folderToShare" with depth "1" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Alice" the PROPFIND response should contain a resource "newFolder" with these key and value pairs:
      | key     | value     |
      | oc:name | newFolder |

  @issue-10719
  Scenario: enable sync of a federated shared resource
    Given using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And using server "REMOTE"
    When user "Brian" enables sync of federated share "folderToShare" using the Graph API
    Then the HTTP status code should be "201"
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["@client.synchronize"],
              "properties": {
                "@client.synchronize": { "const": true }
              }
            }
          }
        }
      }
      """

  @issue-10719
  Scenario: enable sync of a federated shared resource when multiple federated shares exist
    Given using server "LOCAL"
    And user "Alice" has created folder "folderOneShare"
    And user "Alice" has created folder "folderTwoShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderOneShare |
      | space           | Personal       |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | Viewer         |
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderTwoShare |
      | space           | Personal       |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | Viewer         |
    And using server "REMOTE"
    When user "Brian" enables sync of federated share "folderOneShare" using the Graph API
    Then the HTTP status code should be "201"
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
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
                  "required": ["@client.synchronize"],
                  "properties": {
                    "@client.synchronize": { "const": true }
                  }
                },
                {
                  "type": "object",
                  "required": ["@client.synchronize"],
                  "properties": {
                    "@client.synchronize": { "const": false }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: local user shares multiple resources concurrently to a single federated user (Personal Space)
    Given using server "LOCAL"
    And user "Alice" has created the following folders
      | path           |
      | folderToShare1 |
      | folderToShare2 |
    And user "Alice" has uploaded file with content "some content" to "textfile1.txt"
    And user "Alice" has uploaded file with content "hello world" to "textfile2.txt"
    When user "Alice" sends the following concurrent resource share invitations to federated user using the Graph API:
      | resource       | space    | sharee | shareType | permissionsRole |
      | folderToShare1 | Personal | Brian  | user      | Viewer          |
      | folderToShare2 | Personal | Brian  | user      | Editor          |
      | textfile1.txt  | Personal | Brian  | user      | Viewer          |
      | textfile2.txt  | Personal | Brian  | user      | File Editor     |
    Then the HTTP status code of responses on each endpoint should be "200, 200, 200, 200" respectively
    And using server "REMOTE"
    And user "Brian" should have the following federated shares:
      | resource       | permissionsRole | sharer |
      | folderToShare1 | Viewer          | Alice  |
      | folderToShare2 | Editor          | Alice  |
      | textfile1.txt  | Viewer          | Alice  |
      | textfile2.txt  | File Editor     | Alice  |


  Scenario: local user shares multiple resources concurrently to a single federated user (Project Space)
    Given using server "LOCAL"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare1" in space "new-space"
    And user "Alice" has created a folder "folderToShare2" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile1.txt"
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "textfile2.txt"
    When user "Alice" sends the following concurrent resource share invitations to federated user using the Graph API:
      | resource       | space     | sharee | shareType | permissionsRole |
      | folderToShare1 | new-space | Brian  | user      | Viewer          |
      | folderToShare2 | new-space | Brian  | user      | Editor          |
      | textfile1.txt  | new-space | Brian  | user      | Viewer          |
      | textfile2.txt  | new-space | Brian  | user      | File Editor     |
    Then the HTTP status code of responses on each endpoint should be "200, 200, 200, 200" respectively
    And using server "REMOTE"
    And user "Brian" should have the following federated shares:
      | resource       | permissionsRole | sharer |
      | folderToShare1 | Viewer          | Alice  |
      | folderToShare2 | Editor          | Alice  |
      | textfile1.txt  | Viewer          | Alice  |
      | textfile2.txt  | File Editor     | Alice  |


  Scenario: local user shares multiple resources form different spaces concurrently to a single federated user
    Given using server "LOCAL"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created folder "folderToShare1"
    And user "Alice" has uploaded file with content "some content" to "textfile1.txt"
    And user "Alice" has created a folder "folderToShare2" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "textfile2.txt"
    When user "Alice" sends the following concurrent resource share invitations to federated user using the Graph API:
      | resource       | space     | sharee | shareType | permissionsRole |
      | folderToShare1 | Personal  | Brian  | user      | Viewer          |
      | folderToShare2 | new-space | Brian  | user      | Editor          |
      | textfile1.txt  | Personal  | Brian  | user      | Viewer          |
      | textfile2.txt  | new-space | Brian  | user      | File Editor     |
    Then the HTTP status code of responses on each endpoint should be "200, 200, 200, 200" respectively
    And using server "REMOTE"
    And user "Brian" should have the following federated shares:
      | resource       | permissionsRole | sharer |
      | folderToShare1 | Viewer          | Alice  |
      | folderToShare2 | Editor          | Alice  |
      | textfile1.txt  | Viewer          | Alice  |
      | textfile2.txt  | File Editor     | Alice  |


  Scenario: local user shares multiple resources concurrently to multiple federated users
    Given user "Carol" has been created with default attributes
    And "Carol" has accepted invitation
    And using server "LOCAL"
    And user "Alice" has created the following folders
      | path           |
      | folderToShare1 |
      | folderToShare2 |
    And user "Alice" has uploaded file with content "some content" to "textfile1.txt"
    And user "Alice" has uploaded file with content "hello world" to "textfile2.txt"
    When user "Alice" sends the following concurrent resource share invitations to federated user using the Graph API:
      | resource       | space    | sharee | shareType | permissionsRole |
      | folderToShare1 | Personal | Brian  | user      | Viewer          |
      | folderToShare2 | Personal | Carol  | user      | Editor          |
      | textfile1.txt  | Personal | Brian  | user      | Viewer          |
      | textfile2.txt  | Personal | Carol  | user      | File Editor     |
    Then the HTTP status code of responses on each endpoint should be "200, 200, 200, 200" respectively
    And using server "REMOTE"
    And user "Brian" should have the following federated shares:
      | resource       | permissionsRole | sharer |
      | folderToShare1 | Viewer          | Alice  |
      | textfile1.txt  | Viewer          | Alice  |
    And user "Carol" should have the following federated shares:
      | resource       | permissionsRole | sharer |
      | folderToShare2 | Editor          | Alice  |
      | textfile2.txt  | File Editor     | Alice  |
