Feature: List a sharing permissions
  https://owncloud.dev/libre-graph-api/#/drives.permissions/ListPermissions

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |


  Scenario: user lists permissions via the Graph API
    Given user "Alice" has created folder "folder"
    When user "Alice" gets permissions list for folder "folder" of the space "Personal" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@libre.graph.permissions.actions.allowedValues",
          "@libre.graph.permissions.roles.allowedValues"
        ],
        "properties": {
          "@libre.graph.permissions.actions.allowedValues": {
            "type": "array",
            "enum": [
              [
                "libre.graph/driveItem/permissions/create",
                "libre.graph/driveItem/children/create",
                "libre.graph/driveItem/standard/delete",
                "libre.graph/driveItem/path/read",
                "libre.graph/driveItem/quota/read",
                "libre.graph/driveItem/content/read",
                "libre.graph/driveItem/upload/create",
                "libre.graph/driveItem/permissions/read",
                "libre.graph/driveItem/children/read",
                "libre.graph/driveItem/versions/read",
                "libre.graph/driveItem/deleted/read",
                "libre.graph/driveItem/path/update",
                "libre.graph/driveItem/permissions/delete",
                "libre.graph/driveItem/deleted/delete",
                "libre.graph/driveItem/versions/update",
                "libre.graph/driveItem/deleted/update",
                "libre.graph/driveItem/basic/read",
                "libre.graph/driveItem/permissions/update",
                "libre.graph/driveItem/permissions/deny"
              ]
            ]
          },
          "@libre.graph.permissions.roles.allowedValues": {
            "type": "array",
            "minItems": 4,
            "maxItems": 4,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "type": "integer",
                      "enum": [
                        1
                      ]
                    },
                    "description": {
                      "type": "string",
                      "enum": [
                        "Allows upload file or folder"
                      ]
                    },
                    "displayName": {
                      "type": "string",
                      "enum": [
                        "Uploader"
                      ]
                    },
                    "id": {
                      "type": "string",
                      "enum": [
                        "1c996275-f1c9-4e71-abdf-a42f6495e960"
                      ]
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "type": "integer",
                      "enum": [
                        2
                      ]
                    },
                    "description": {
                      "type": "string",
                      "enum": [
                        "Allows reading the shared file or folder"
                      ]
                    },
                    "displayName": {
                      "type": "string",
                      "enum": [
                        "Viewer"
                      ]
                    },
                    "id": {
                      "type": "string",
                      "enum": [
                        "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
                      ]
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "type": "integer",
                      "enum": [
                        3
                      ]
                    },
                    "description": {
                      "type": "string",
                      "enum": [
                        "Allows reading and updating file"
                      ]
                    },
                    "displayName": {
                      "type": "string",
                      "enum": [
                        "Editor"
                      ]
                    },
                    "id": {
                      "type": "string",
                      "enum": [
                        "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"
                      ]
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "type": "integer",
                      "enum": [
                        4
                      ]
                    },
                    "description": {
                      "type": "string",
                      "enum": [
                        "Allows creating, reading, updating and deleting the shared file or folder"
                      ]
                    },
                    "displayName": {
                      "type": "string",
                      "enum": [
                        "Editor"
                      ]
                    },
                    "id": {
                      "type": "string",
                      "enum": [
                        "fb6c3e19-e378-47e5-b277-9732f9de6e21"
                      ]
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists permissions of a project space
    Given using spaces DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    When user "Alice" lists the permissions of space "new-space" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@libre.graph.permissions.actions.allowedValues",
          "@libre.graph.permissions.roles.allowedValues"
        ],
        "properties": {
          "@libre.graph.permissions.actions.allowedValues": {
            "const": [
              "libre.graph/driveItem/permissions/create",
              "libre.graph/driveItem/children/create",
              "libre.graph/driveItem/standard/delete",
              "libre.graph/driveItem/path/read",
              "libre.graph/driveItem/quota/read",
              "libre.graph/driveItem/content/read",
              "libre.graph/driveItem/upload/create",
              "libre.graph/driveItem/permissions/read",
              "libre.graph/driveItem/children/read",
              "libre.graph/driveItem/versions/read",
              "libre.graph/driveItem/deleted/read",
              "libre.graph/driveItem/path/update",
              "libre.graph/driveItem/permissions/delete",
              "libre.graph/driveItem/deleted/delete",
              "libre.graph/driveItem/versions/update",
              "libre.graph/driveItem/deleted/update",
              "libre.graph/driveItem/basic/read",
              "libre.graph/driveItem/permissions/update",
              "libre.graph/driveItem/permissions/deny"
            ]
          },
          "@libre.graph.permissions.roles.allowedValues": {
            "type": "array",
            "minItems": 3,
            "maxItems": 3,
            "uniqueItems": true,
            "items": {
              "oneOf":  [
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "const": 1
                    },
                    "description": {
                      "const": "Allows reading the shared space"
                    },
                    "displayName": {
                      "const": "Space Viewer"
                    },
                    "id": {
                      "const": "a8d5fe5e-96e3-418d-825b-534dbdf22b99"
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "const": 2
                    },
                    "description": {
                      "const": "Allows creating, reading, updating and deleting file or folder in the shared space"
                    },
                    "displayName": {
                      "const": "Space Editor"
                    },
                    "id": {
                      "const": "58c63c02-1d89-4572-916a-870abc5a1b7d"
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "const": 3
                    },
                    "description": {
                      "const": "Grants manager permissions on a resource. Semantically equivalent to co-owner"
                    },
                    "displayName": {
                      "const": "Manager"
                    },
                    "id": {
                      "const": "312c0871-5ef7-4b3a-85b6-0e4074c64049"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @issues-8352
  Scenario: sharer lists permissions of a shared project space
    Given using spaces DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has sent the following share invitation:
      | space           | new-space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    And user "Alice" has created the following link share:
      | space           | new-space |
      | permissionsRole | view      |
      | password        | %public%  |
      | resource        | new-space |
    When user "Alice" lists the permissions of space "new-space" using the Graph API
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
          "@libre.graph.permissions.actions.allowedValues": {
            "const": [
              "libre.graph/driveItem/permissions/create",
              "libre.graph/driveItem/children/create",
              "libre.graph/driveItem/standard/delete",
              "libre.graph/driveItem/path/read",
              "libre.graph/driveItem/quota/read",
              "libre.graph/driveItem/content/read",
              "libre.graph/driveItem/upload/create",
              "libre.graph/driveItem/permissions/read",
              "libre.graph/driveItem/children/read",
              "libre.graph/driveItem/versions/read",
              "libre.graph/driveItem/deleted/read",
              "libre.graph/driveItem/path/update",
              "libre.graph/driveItem/permissions/delete",
              "libre.graph/driveItem/deleted/delete",
              "libre.graph/driveItem/versions/update",
              "libre.graph/driveItem/deleted/update",
              "libre.graph/driveItem/basic/read",
              "libre.graph/driveItem/permissions/update",
              "libre.graph/driveItem/permissions/deny"
            ]
          },
          "@libre.graph.permissions.roles.allowedValues": {
            "type": "array",
            "minItems": 3,
            "maxItems": 3,
            "uniqueItems": true,
            "items": {
              "oneOf":  [
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "const": 1
                    },
                    "description": {
                      "const": "Allows reading the shared space"
                    },
                    "displayName": {
                      "const": "Space Viewer"
                    },
                    "id": {
                      "const": "a8d5fe5e-96e3-418d-825b-534dbdf22b99"
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "const": 2
                    },
                    "description": {
                      "const": "Allows creating, reading, updating and deleting file or folder in the shared space"
                    },
                    "displayName": {
                      "const": "Space Editor"
                    },
                    "id": {
                      "const": "58c63c02-1d89-4572-916a-870abc5a1b7d"
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "@libre.graph.weight",
                    "description",
                    "displayName",
                    "id"
                  ],
                  "properties": {
                    "@libre.graph.weight": {
                      "const": 3
                    },
                    "description": {
                      "const": "Grants manager permissions on a resource. Semantically equivalent to co-owner"
                    },
                    "displayName": {
                      "const": "Manager"
                    },
                    "id": {
                      "const": "312c0871-5ef7-4b3a-85b6-0e4074c64049"
                    }
                  }
                }
              ]
            }
          },
          "value": {
            "type": "array",
            "minItems": 3,
            "maxItems": 3,
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
                          "required": ["displayName","id"],
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
                },
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
                          "required": ["displayName","id"],
                          "properties": {
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
                },
                {
                  "type": "object",
                  "required": [
                    "hasPassword",
                    "id",
                    "link"
                  ],
                  "properties": {
                    "hasPassword": {
                      "const": true
                    },
                    "id": {
                      "type": "string",
                      "pattern": "^[a-zA-Z]{15}$"
                    },
                    "link": {
                      "type": "object",
                      "required": [
                        "@libre.graph.displayName",
                        "@libre.graph.quickLink",
                        "preventsDownload",
                        "type",
                        "webUrl"
                      ],
                      "properties": {
                        "@libre.graph.displayName": {
                          "const": ""
                        },
                        "@libre.graph.quickLink": {
                          "const": false
                        },
                        "preventsDownload": {
                          "const": false
                        },
                        "type": {
                          "const": "view"
                        },
                        "webUrl": {
                          "type": "string",
                          "pattern": "^%base_url%\/s\/[a-zA-Z]{15}$"
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