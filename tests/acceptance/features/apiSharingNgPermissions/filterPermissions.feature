Feature: filter sharing permissions
  As a user
  I want to filter sharing permissions
  So that I can get specific permissions

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: filter permissions of a folder for federated user type (Personal space)
    Given user "Alice" has created folder "folder"
    When user "Alice" lists permissions with following filters for folder "folder" of the space "Personal" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Federated"')) |
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
            "minItems": 2,
            "maxItems": 2,
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
                      "const": 1
                    },
                    "description": {
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload, edit, add and delete."
                    },
                    "displayName": {
                      "const": "Can edit with trashbin"
                    },
                    "id": {
                      "const": "fb6c3e19-e378-47e5-b277-9732f9de6e21"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @env-config
  Scenario: filter lists permissions of a file for federated user type (Personal space)
    Given the administrator has enabled the permissions role "Secure Viewer"
    And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    When user "Alice" lists permissions with following filters for file "textfile0.txt" of the space "Personal" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Federated"')) |
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
            "minItems": 2,
            "maxItems": 2,
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
                      "const": 1
                    },
                    "description": {
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload and edit."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @env-config
  Scenario: filter lists permissions of a folder for Member user type (Personal space)
    Given the administrator has enabled the permissions role "Denied"
    And user "Alice" has created folder "folder"
    When user "Alice" lists permissions with following filters for folder "folder" of the space "Personal" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Member"')) |
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
                      "const": 1
                    },
                    "description": {
                      "const": "Deny all access."
                    },
                    "displayName": {
                      "const": "Cannot access"
                    },
                    "id": {
                      "const": "63e64e19-8d43-42ec-a738-2b6af2610efa"
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
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload, edit and add."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "1c996275-f1c9-4e71-abdf-a42f6495e960"
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
                      "const": 4
                    },
                    "description": {
                      "const": "View, download, upload, edit, add and delete."
                    },
                    "displayName": {
                      "const": "Can edit with trashbin"
                    },
                    "id": {
                      "const": "fb6c3e19-e378-47e5-b277-9732f9de6e21"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @env-config
  Scenario: filter permissions of a file for Member user type (Personal space)
    Given the administrator has enabled the permissions role "Secure Viewer"
    And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    When user "Alice" lists permissions with following filters for file "textfile0.txt" of the space "Personal" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Member"')) |
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
                      "const": 1
                    },
                    "description": {
                      "const": "View only documents, images and PDFs. Watermarks will be applied."
                    },
                    "displayName": {
                      "const": "Can view (secure)"
                    },
                    "id": {
                      "const": "aa97fe03-7980-45ac-9e50-b325749fd7e6"
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
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload and edit."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists allowed role permissions of a folder (Personal space)
    Given user "Alice" has created folder "folderToShare"
    When user "Alice" lists permissions with following filters for folder "folderToShare" of the space "Personal" using the Graph API:
      | $select=@libre.graph.permissions.roles.allowedValues |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
          "required": [
            "@libre.graph.permissions.roles.allowedValues"
          ],
          "properties": {
            "@libre.graph.permissions.actions.allowedValues": { "const": null },
            "@libre.graph.permissions.roles.allowedValues": {
              "type": "array",
              "minItems": 3,
              "maxItems": 3,
              "uniqueItems": true,
              "items": {
                "oneOf":[
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
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload, edit and add."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "1c996275-f1c9-4e71-abdf-a42f6495e960"
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
                      "const": "View, download, upload, edit, add and delete."
                    },
                    "displayName": {
                      "const": "Can edit with trashbin"
                    },
                    "id": {
                      "const": "fb6c3e19-e378-47e5-b277-9732f9de6e21"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @env-config
  Scenario: user lists allowed role permissions of a file (Personal space)
    Given the administrator has enabled the permissions role "Secure Viewer"
    And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    When user "Alice" lists permissions with following filters for file "textfile0.txt" of the space "Personal" using the Graph API:
      | $select=@libre.graph.permissions.roles.allowedValues |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@libre.graph.permissions.roles.allowedValues"
        ],
        "properties": {
          "@libre.graph.permissions.actions.allowedValues": { "const": null },
          "@libre.graph.permissions.roles.allowedValues": {
            "type": "array",
            "minItems": 3,
            "maxItems": 3,
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
                      "const": 1
                    },
                    "description": {
                      "const": "View only documents, images and PDFs. Watermarks will be applied."
                    },
                    "displayName": {
                      "const": "Can view (secure)"
                    },
                    "id": {
                      "const": "aa97fe03-7980-45ac-9e50-b325749fd7e6"
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
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload and edit."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @env-config
  Scenario: filter lists permissions of a file for Member user type (Project space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    When user "Alice" lists permissions with following filters for file "textfile.txt" of the space "new-space" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Member"')) |
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
                      "const": 1
                    },
                    "description": {
                      "const": "View only documents, images and PDFs. Watermarks will be applied."
                    },
                    "displayName": {
                      "const": "Can view (secure)"
                    },
                    "id": {
                      "const": "aa97fe03-7980-45ac-9e50-b325749fd7e6"
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
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload and edit."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @env-config
  Scenario: filter lists permissions of a folder for Member user type (Project space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder" in space "new-space"
    When user "Alice" lists permissions with following filters for folder "folder" of the space "new-space" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Member"')) |
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
                      "const": 1
                    },
                    "description": {
                      "const": "Deny all access."
                    },
                    "displayName": {
                      "const": "Cannot access"
                    },
                    "id": {
                      "const": "63e64e19-8d43-42ec-a738-2b6af2610efa"
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
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload, edit and add."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "1c996275-f1c9-4e71-abdf-a42f6495e960"
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
                      "const": 4
                    },
                    "description": {
                      "const": "View, download, upload, edit, add and delete."
                    },
                    "displayName": {
                      "const": "Can edit with trashbin"
                    },
                    "id": {
                      "const": "fb6c3e19-e378-47e5-b277-9732f9de6e21"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @env-config
  Scenario: user lists allowed role permissions of a folder (Project space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder" in space "new-space"
    When user "Alice" lists permissions with following filters for folder "folder" of the space "new-space" using the Graph API:
      | $select=@libre.graph.permissions.roles.allowedValues |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
          "required": [
            "@libre.graph.permissions.roles.allowedValues"
          ],
          "properties": {
            "@libre.graph.permissions.actions.allowedValues": { "const": null },
            "@libre.graph.permissions.roles.allowedValues": {
              "type": "array",
              "minItems": 4,
              "maxItems": 4,
              "uniqueItems": true,
              "items": {
                "oneOf":[
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
                      "const": "Deny all access."
                    },
                    "displayName": {
                      "const": "Cannot access"
                    },
                    "id": {
                      "const": "63e64e19-8d43-42ec-a738-2b6af2610efa"
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
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload, edit and add."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "1c996275-f1c9-4e71-abdf-a42f6495e960"
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
                      "const": 4
                    },
                    "description": {
                      "const": "View, download, upload, edit, add and delete."
                    },
                    "displayName": {
                      "const": "Can edit with trashbin"
                    },
                    "id": {
                      "const": "fb6c3e19-e378-47e5-b277-9732f9de6e21"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @env-config
  Scenario: user lists allowed role permissions of a file (Project space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    When user "Alice" lists permissions with following filters for file "textfile.txt" of the space "new-space" using the Graph API:
      | $select=@libre.graph.permissions.roles.allowedValues |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "@libre.graph.permissions.roles.allowedValues"
        ],
        "properties": {
          "@libre.graph.permissions.actions.allowedValues": { "const": null },
          "@libre.graph.permissions.roles.allowedValues": {
            "type": "array",
            "minItems": 3,
            "maxItems": 3,
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
                      "const": 1
                    },
                    "description": {
                      "const": "View only documents, images and PDFs. Watermarks will be applied."
                    },
                    "displayName": {
                      "const": "Can view (secure)"
                    },
                    "id": {
                      "const": "aa97fe03-7980-45ac-9e50-b325749fd7e6"
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
                      "const": "View and download."
                    },
                    "displayName": {
                      "const": "Can view"
                    },
                    "id": {
                      "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
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
                      "const": "View, download, upload and edit."
                    },
                    "displayName": {
                      "const": "Can edit"
                    },
                    "id": {
                      "const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @issue-9745 @env-config
  Scenario: user lists allowed file permissions for federated user
    Given the administrator has enabled the permissions role "Secure Viewer"
    And user "Alice" has uploaded file with content "ocm test" to "/textfile.txt"
    When user "Alice" lists permissions with following filters for file "textfile.txt" of the space "Personal" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Federated"')) |
      | $select=@libre.graph.permissions.roles.allowedValue                                                                                 |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
          "required": [
            "@libre.graph.permissions.roles.allowedValues"
          ],
          "properties": {
            "@libre.graph.permissions.roles.allowedValues": {
              "type": "array",
              "minItems": 2,
              "maxItems": 2,
              "uniqueItems": true,
              "items": {
                "oneOf":[
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 1},
                    "description": {"const": "View and download."},
                    "displayName": {"const": "Can view"},
                    "id": {"const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"}
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 2},
                    "description": {"const": "View, download, upload and edit."},
                    "displayName": {"const": "Can edit"},
                    "id": {"const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a" }
                  }
                }
              ]
            }
          }
        }
      }
      """

  @issue-9745 @env-config
  Scenario: user lists allowed folder permissions for federated user
    Given the administrator has enabled the permissions role "Denied"
    And user "Alice" has created folder "folderToShare"
    When user "Alice" lists permissions with following filters for folder "folderToShare" of the space "Personal" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Federated"')) |
      | $select=@libre.graph.permissions.roles.allowedValue                                                                                 |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
          "required": [
            "@libre.graph.permissions.roles.allowedValues"
          ],
          "properties": {
            "@libre.graph.permissions.roles.allowedValues": {
              "type": "array",
              "minItems": 2,
              "maxItems": 2,
              "uniqueItems": true,
              "items": {
                "oneOf":[
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 1},
                    "description": {"const": "View and download."},
                    "displayName": {"const": "Can view"},
                    "id": {"const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"}
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 2},
                    "description": {"const": "View, download, upload, edit, add and delete."},
                    "displayName": {"const": "Can edit with trashbin"},
                    "id": {"const": "fb6c3e19-e378-47e5-b277-9732f9de6e21"}
                  }
                }
              ]
            }
          }
        }
      }
      """

  @issue-9745 @env-config
  Scenario: user lists allowed file permissions for federated user (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    When user "Alice" lists permissions with following filters for file "textfile.txt" of the space "new-space" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Federated"')) |
      | $select=@libre.graph.permissions.roles.allowedValue                                                                                 |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
          "required": [
            "@libre.graph.permissions.roles.allowedValues"
          ],
          "properties": {
            "@libre.graph.permissions.roles.allowedValues": {
              "type": "array",
              "minItems": 2,
              "maxItems": 2,
              "uniqueItems": true,
              "items": {
                "oneOf":[
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 1},
                    "description": {"const": "View and download."},
                    "displayName": {"const": "Can view"},
                    "id": {"const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"}
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 2},
                    "description": {"const": "View, download, upload and edit."},
                    "displayName": {"const": "Can edit"},
                    "id": {"const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"}
                  }
                }
              ]
            }
          }
        }
      }
      """

  @issue-9745 @env-config
  Scenario: user lists allowed folder permissions for federated user (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    When user "Alice" lists permissions with following filters for folder "folderToShare" of the space "new-space" using the Graph API:
      | $filter=@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition,+'@Subject.UserType=="Federated"')) |
      | $select=@libre.graph.permissions.roles.allowedValue                                                                                 |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
          "required": [
            "@libre.graph.permissions.roles.allowedValues"
          ],
          "properties": {
            "@libre.graph.permissions.roles.allowedValues": {
              "type": "array",
              "minItems": 2,
              "maxItems": 2,
              "uniqueItems": true,
              "items": {
                "oneOf":[
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 1},
                    "description": {"const": "View and download."},
                    "displayName": {"const": "Can view"},
                    "id": {"const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"}
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 2},
                    "description": {"const": "View, download, upload, edit, add and delete."},
                    "displayName": {"const": "Can edit with trashbin"},
                    "id": {"const": "fb6c3e19-e378-47e5-b277-9732f9de6e21"}
                  }
                }
              ]
            }
          }
        }
      }
      """
