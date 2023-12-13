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
            "items": [
              {
                "type": "string",
                "required": [
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
              }
            ]
          },
          "@libre.graph.permissions.roles.allowedValues": {
            "type": "array",
            "items": [
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
                      5
                    ]
                  },
                  "description": {
                    "type": "string",
                    "enum": [
                      "Grants co-owner permissions on a resource"
                    ]
                  },
                  "displayName": {
                    "type": "string",
                    "enum": [
                      "Co Owner"
                    ]
                  },
                  "id": {
                    "type": "string",
                    "enum": [
                      "3a4ba8e9-6a0d-4235-9140-0e7a34007abe"
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
                      6
                    ]
                  },
                  "description": {
                    "type": "string",
                    "enum": [
                      "Grants manager permissions on a resource. Semantically equivalent to co-owner"
                    ]
                  },
                  "displayName": {
                    "type": "string",
                    "enum": [
                      "Manager"
                    ]
                  },
                  "id": {
                    "type": "string",
                    "enum": [
                      "312c0871-5ef7-4b3a-85b6-0e4074c64049"
                    ]
                  }
                }
              }
            ]
          }
        }
      }
      """
