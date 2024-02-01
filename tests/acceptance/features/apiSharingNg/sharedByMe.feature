Feature: resources shared by user
  As a user
  I want to get resources shared by me
  So that I can know about what resources are shared with others

  https://owncloud.dev/libre-graph-api/#/me.drive/ListSharedByMe

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario: user lists the shared file from personal space
    Given user "Alice" has uploaded file with content "hello world" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file         |
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
              "type": "string",
              "enum": ["personal"]
            },
            "path": {
              "type": "string",
              "enum": ["/"]
            },
            "name": {
              "type": "string",
              "enum": ["/"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["textfile.txt"]
        },
        "size": {
          "type": "number",
          "enum": [
            11
          ]
        }
      }
    }
    """


  Scenario: user lists the shared file inside of a folder from personal space
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "hello world" to "FolderToShare/textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file                       |
      | resource        | FolderToShare/textfile.txt |
      | space           | Personal                   |
      | sharee          | Brian                      |
      | shareType       | user                       |
      | permissionsRole | Viewer                     |
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
              "type": "string",
              "enum": ["personal"]
            },
            "path": {
              "type": "string",
              "enum": ["/FolderToShare"]
            },
            "name": {
              "type": "string",
              "enum": ["FolderToShare"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["textfile.txt"]
        },
        "size": {
          "type": "number",
          "enum": [
            11
          ]
        }
      }
    }
    """


  Scenario: user lists the shared folder from personal space
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following share invitation:
      | resourceType    | folder        |
      | resource        | FolderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "FolderToShare" with the following data:
    """
    {
      "type": "object",
      "required": [
        "parentReference",
        "permissions",
        "name"
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
              "type": "string",
              "enum": ["personal"]
            },
            "path": {
              "type": "string",
              "enum": ["/"]
            },
            "name": {
              "type": "string",
              "enum": ["/"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["FolderToShare"]
        }
      }
    }
    """


  Scenario: user lists shared resources from personal space
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "hello world" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file         |
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Alice" has sent the following share invitation:
      | resourceType    | folder        |
      | resource        | FolderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
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
              "type": "string",
              "enum": ["personal"]
            },
            "path": {
              "type": "string",
              "enum": ["/"]
            },
            "name": {
              "type": "string",
              "enum": ["/"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["textfile.txt"]
        },
        "size": {
          "type": "number",
          "enum": [
            11
          ]
        }
      }
    }
    """
    And the JSON data of the response should contain resource "FolderToShare" with the following data:
    """
    {
      "type": "object",
      "required": [
        "parentReference",
        "permissions",
        "name"
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
              "type": "string",
              "enum": ["personal"]
            },
            "path": {
              "type": "string",
              "enum": ["/"]
            },
            "name": {
              "type": "string",
              "enum": ["/"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["FolderToShare"]
        }
      }
    }
    """


  Scenario: user lists the shared file from project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "hello world" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file         |
      | resource        | textfile.txt |
      | space           | NewSpace     |
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
              "type": "string",
              "enum": ["project"]
            },
            "path": {
              "type": "string",
              "enum": ["/"]
            },
            "name": {
              "type": "string",
              "enum": ["/"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["textfile.txt"]
        },
        "size": {
          "type": "number",
          "enum": [
            11
          ]
        }
      }
    }
    """


  Scenario: user lists the shared file inside of a folder from project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "hello world" to "FolderToShare/textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file                       |
      | resource        | FolderToShare/textfile.txt |
      | space           | NewSpace                   |
      | sharee          | Brian                      |
      | shareType       | user                       |
      | permissionsRole | Viewer                     |
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
              "type": "string",
              "enum": ["project"]
            },
            "path": {
              "type": "string",
              "enum": ["/FolderToShare"]
            },
            "name": {
              "type": "string",
              "enum": ["FolderToShare"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["textfile.txt"]
        },
        "size": {
          "type": "number",
          "enum": [
            11
          ]
        }
      }
    }
    """


  Scenario: user lists the folder shared from project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has sent the following share invitation:
      | resourceType    | folder        |
      | resource        | FolderToShare |
      | space           | NewSpace      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "FolderToShare" with the following data:
    """
    {
      "type": "object",
      "required": [
        "parentReference",
        "permissions",
        "name"
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
              "type": "string",
              "enum": ["project"]
            },
            "path": {
              "type": "string",
              "enum": ["/"]
            },
            "name": {
              "type": "string",
              "enum": ["/"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["FolderToShare"]
        }
      }
    }
    """


  Scenario: user lists resources shared from project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "hello world" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resourceType    | file         |
      | resource        | textfile.txt |
      | space           | NewSpace     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Alice" has sent the following share invitation:
      | resourceType    | folder        |
      | resource        | FolderToShare |
      | space           | NewSpace      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
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
              "type": "string",
              "enum": ["project"]
            },
            "path": {
              "type": "string",
              "enum": ["/"]
            },
            "name": {
              "type": "string",
              "enum": ["/"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["textfile.txt"]
        },
        "size": {
          "type": "number",
          "enum": [
            11
          ]
        }
      }
    }
    """
    And the JSON data of the response should contain resource "FolderToShare" with the following data:
    """
    {
      "type": "object",
      "required": [
        "parentReference",
        "permissions",
        "name"
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
              "type": "string",
              "enum": ["project"]
            },
            "path": {
              "type": "string",
              "enum": ["/"]
            },
            "name": {
              "type": "string",
              "enum": ["/"]
            },
            "id": {
              "type": "string",
              "pattern": "^%file_id_pattern%$"
            }
          }
        },
        "permissions": {
          "type": "array",
          "items": [
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
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
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
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string",
                      "pattern": "^%role_id_pattern%$"
                    }
                  ]
                }
              }
            }
          ]
        },
        "name": {
          "type": "string",
          "enum": ["FolderToShare"]
        }
      }
    }
    """
