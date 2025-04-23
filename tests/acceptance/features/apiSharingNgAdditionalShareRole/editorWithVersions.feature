@env-config
Feature: an user shares resources
  As a user
  I want to share resources with Editor With Versions role
  So that users can edit the resource and see the versions

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario: sharee checks version of a file shared with FileEditorWithVersions role
    Given the administrator has enabled the permissions role "File Editor With Versions"
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile.txt              |
      | space           | Personal                  |
      | sharee          | Brian                     |
      | shareType       | user                      |
      | permissionsRole | File Editor With Versions |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": ["createdDateTime", "id", "roles", "grantedToV2", "invitation"],
              "properties": {
                "id": { "pattern": "^%permissions_id_pattern%$" },
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": { "const": "b173329d-cf2e-42f0-a595-ee410645d840" }
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
                          "required": ["displayName", "id", "@libre.graph.userType"],
                          "properties": {
                            "displayName": { "const": "Alice Hansen" },
                            "id": { "pattern": "^%user_id_pattern%$" },
                            "@libre.graph.userType": { "const": "Member" }
                          }
                        }
                      }
                    }
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id", "displayName", "@libre.graph.userType"],
                      "properties": {
                        "id": { "pattern": "^%user_id_pattern%$" },
                        "displayName": { "const": "Brian Murphy" },
                        "@libre.graph.userType": { "const": "Member" }
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
    And user "Brian" has uploaded file with content "updated content" to "Shares/textfile.txt"
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"


  Scenario: sharee checks version of a file inside a folder shared with EditorWithVersions role
    Given the administrator has enabled the permissions role "Editor With Versions"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "to share" to "folderToShare/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | folderToShare        |
      | space           | Personal             |
      | sharee          | Brian                |
      | shareType       | user                 |
      | permissionsRole | Editor With Versions |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": ["createdDateTime", "id", "roles", "grantedToV2", "invitation"],
              "properties": {
                "id": { "pattern": "^%permissions_id_pattern%$" },
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": { "const": "0911d62b-1e3f-4778-8b1b-903b7e4e8476" }
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
                          "required": ["displayName", "id", "@libre.graph.userType"],
                          "properties": {
                            "displayName": { "const": "Alice Hansen" },
                            "id": { "pattern": "^%user_id_pattern%$" },
                            "@libre.graph.userType": { "const": "Member" }
                          }
                        }
                      }
                    }
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id", "displayName", "@libre.graph.userType"],
                      "properties": {
                        "id": { "pattern": "^%user_id_pattern%$" },
                        "displayName": { "const": "Brian Murphy" },
                        "@libre.graph.userType": { "const": "Member" }
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
    And user "Brian" has uploaded file with content "updated content" to "Shares/folderToShare/textfile.txt"
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"


  Scenario: user lists permissions of a file in personal space after enabling FileEditorWithVersions role
    Given the administrator has enabled the permissions role "File Editor With Versions"
    And user "Alice" has uploaded file with content "hello world" to "textfile0.txt"
    When user "Alice" gets permissions list for file "textfile0.txt" of the space "Personal" using the Graph API
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
              "oneOf":[
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 1 },
                    "description": { "const": "View and download." },
                    "displayName": { "const": "Can view" },
                    "id": { "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id" ],
                  "properties": {
                    "@libre.graph.weight": { "const": 2 },
                    "description": { "const": "View, download and edit." },
                    "displayName": { "const": "Can edit without versions" },
                    "id": { "const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 3 },
                    "description": { "const": "View, download, edit and show all invited people, show all versions." },
                    "displayName": { "const": "Can edit" },
                    "id": { "const": "b173329d-cf2e-42f0-a595-ee410645d840" }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists permissions of a file in project space after enabling FileEditorWithVersions role
    Given using spaces DAV path
    And the administrator has enabled the permissions role "File Editor With Versions"
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "textfile0.txt"
    When user "Alice" gets permissions list for file "textfile0.txt" of the space "new-space" using the Graph API
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
                "oneOf":[
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 1 },
                    "description": { "const": "View and download." },
                    "displayName": { "const": "Can view" },
                    "id": { "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 2 },
                    "description": { "const": "View, download and edit." },
                    "displayName": { "const": "Can edit without versions" },
                    "id": { "const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 3 },
                    "description": { "const": "View, download, edit and show all invited people, show all versions." },
                    "displayName": { "const": "Can edit" },
                    "id": { "const": "b173329d-cf2e-42f0-a595-ee410645d840" }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists permissions of a folder in personal space after enabling EditorWithVersions role
    Given the administrator has enabled the permissions role "Editor With Versions"
    And user "Alice" has created folder "folderToShare"
    When user "Alice" gets permissions list for folder "folderToShare" of the space "Personal" using the Graph API
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
              "oneOf":[
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 1 },
                    "description": { "const": "View and download." },
                    "displayName": { "const": "Can view" },
                    "id": { "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 2 },
                    "description": { "const": "View, download and upload." },
                    "displayName": { "const": "Can upload" },
                    "id": { "const": "1c996275-f1c9-4e71-abdf-a42f6495e960" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 3 },
                    "description": { "const": "View, download, upload, edit, add and delete." },
                    "displayName": { "const": "Can edit without versions" },
                    "id": { "const": "fb6c3e19-e378-47e5-b277-9732f9de6e21" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 4 },
                    "description": { "const": "View, download, upload, edit, delete and show all invited people, show all versions." },
                    "displayName": { "const": "Can edit" },
                    "id": { "const": "0911d62b-1e3f-4778-8b1b-903b7e4e8476" }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists permissions of a folder in project space after enabling EditorWithVersions role
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Editor With Versions"
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder" in space "new-space"
    When user "Alice" gets permissions list for folder "folder" of the space "new-space" using the Graph API
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
                "oneOf":[
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 1 },
                    "description": { "const": "View and download." },
                    "displayName": { "const": "Can view" },
                    "id": { "const": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 2 },
                    "description": { "const": "View, download and upload." },
                    "displayName": { "const": "Can upload" },
                    "id": { "const": "1c996275-f1c9-4e71-abdf-a42f6495e960" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 3 },
                    "description": { "const": "View, download, upload, edit, add and delete." },
                    "displayName": { "const": "Can edit without versions" },
                    "id": { "const": "fb6c3e19-e378-47e5-b277-9732f9de6e21" }
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight", "description", "displayName", "id"],
                  "properties": {
                    "@libre.graph.weight": { "const": 4 },
                    "description": { "const": "View, download, upload, edit, delete and show all invited people, show all versions." },
                    "displayName": { "const": "Can edit" },
                    "id": { "const": "0911d62b-1e3f-4778-8b1b-903b7e4e8476" }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: sharee lists the file share shared with FileEditorWithVersions permission role (Personal Space)
    Given the administrator has enabled the permissions role "File Editor With Versions"
    And user "Alice" has uploaded file with content "hello world" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt              |
      | space           | Personal                  |
      | sharee          | Brian                     |
      | shareType       | user                      |
      | permissionsRole | File Editor With Versions |
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
                "@UI.Hidden": { "const": false },
                "@client.synchronize": { "const": true },
                "createdBy": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id"],
                      "properties": {
                        "displayName": { "const": "Alice Hansen" },
                        "id": { "pattern": "^%user_id_pattern%$" }
                      }
                    }
                  }
                },
                "eTag": { "pattern": "%etag_pattern%" },
                "file": {
                  "type": "object",
                  "required": ["mimeType"],
                  "properties": { "mimeType": { "const": "text/plain" } }
                },
                "id": { "pattern": "^%share_id_pattern%$" },
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
                    "parentReference",
                    "permissions",
                    "size"
                  ],
                  "properties": {
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["id", "displayName"],
                          "properties": {
                            "id": { "pattern": "^%user_id_pattern%$" },
                            "displayName": { "const": "Alice Hansen" }
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
                      "required": ["driveId", "driveType"],
                      "properties": {
                        "driveId": { "pattern": "^%file_id_pattern%$" },
                        "driveType": { "const": "personal" }
                      }
                    },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["grantedToV2", "id", "invitation", "roles"],
                        "properties": {
                          "id": { "pattern": "^%permissions_id_pattern%$" },
                          "grantedToV2": {
                            "type": "object",
                            "required": ["user"],
                            "properties": {
                              "user": {
                                "type": "object",
                                "required": ["displayName", "id"],
                                "properties": {
                                  "displayName": { "const": "Brian Murphy" },
                                  "id": { "pattern": "^%user_id_pattern%$" }
                                }
                              }
                            }
                          },
                          "invitation": {
                            "type": "object",
                            "properties": {
                              "invitedBy": {
                                "type": "object",
                                "properties": {
                                  "user": {
                                    "type": "object",
                                    "properties": {
                                      "displayName": { "const": "Alice Hansen" },
                                      "id": { "pattern": "^%user_id_pattern%$" }
                                    },
                                    "required": ["displayName", "id"]
                                  }
                                },
                                "required": ["user"]
                              }
                            },
                            "required": ["invitedBy"]
                          },
                          "roles": {
                            "type": "array",
                            "minItems": 1,
                            "maxItems": 1,
                            "items": { "const": "b173329d-cf2e-42f0-a595-ee410645d840" }
                          }
                        }
                      }
                    }
                  }
                },
                "size": { "const": 11 }
              }
            }
          }
        }
      }
      """


  Scenario: sharee lists the file share shared with FileEditorWithVersions permission role (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "File Editor With Versions"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "testfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | testfile.txt              |
      | space           | new-space                 |
      | sharee          | Brian                     |
      | shareType       | user                      |
      | permissionsRole | File Editor With Versions |
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
              "required": [
                "@UI.Hidden",
                "@client.synchronize",
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
                "@UI.Hidden": { "const": false },
                "@client.synchronize": { "const": true },
                "eTag": { "pattern": "%etag_pattern%" },
                "file": {
                  "type": "object",
                  "required": ["mimeType"],
                  "properties": { "mimeType": { "const": "text/plain" } }
                },
                "id": { "pattern": "^%share_id_pattern%$" },
                "name": { "const": "testfile.txt" },
                "parentReference": {
                  "type": "object",
                  "required": ["driveId", "driveType"],
                  "properties": {
                    "driveId": { "pattern": "^%space_id_pattern%$" },
                    "driveType": { "const": "virtual" },
                    "id": { "pattern": "^%file_id_pattern%$" }
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "eTag",
                    "file",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "parentReference",
                    "permissions",
                    "size"
                  ],
                  "properties": {
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["id", "displayName"],
                          "properties": {
                            "id": { "pattern": "^%user_id_pattern%$" },
                            "displayName": { "const": "Alice Hansen" }
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
                    "name": { "const": "testfile.txt" },
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId", "driveType"],
                      "properties": {
                        "driveId": { "pattern": "^%file_id_pattern%$" },
                        "driveType": { "const": "project" }
                      }
                    },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["grantedToV2", "id", "invitation", "roles"],
                        "properties": {
                          "id": { "pattern": "^%permissions_id_pattern%$" },
                          "grantedToV2": {
                            "type": "object",
                            "required": ["user"],
                            "properties": {
                              "user": {
                                "type": "object",
                                "required": ["displayName", "id"],
                                "properties": {
                                  "displayName": { "const": "Brian Murphy" },
                                  "id": { "pattern": "^%user_id_pattern%$" }
                                }
                              }
                            }
                          },
                          "invitation": {
                            "type": "object",
                            "properties": {
                              "invitedBy": {
                                "type": "object",
                                "properties": {
                                  "user": {
                                    "type": "object",
                                    "properties": {
                                      "displayName": { "const": "Alice Hansen" },
                                      "id": { "pattern": "^%user_id_pattern%$" }
                                    },
                                    "required": ["displayName", "id"]
                                  }
                                },
                                "required": ["user"]
                              }
                            },
                            "required": ["invitedBy"]
                          },
                          "roles": {
                            "type": "array",
                            "minItems": 1,
                            "maxItems": 1,
                            "items": { "const": "b173329d-cf2e-42f0-a595-ee410645d840" }
                          }
                        }
                      }
                    }
                  }
                },
                "size": { "const": 12 }
              }
            }
          }
        }
      }
      """


  Scenario: sharee lists the folder share shared with EditorWithVersions permission role (Personal Space)
    Given the administrator has enabled the permissions role "Editor With Versions"
    And user "Alice" has created folder "folder"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder               |
      | space           | Personal             |
      | sharee          | Brian                |
      | shareType       | user                 |
      | permissionsRole | Editor With Versions |
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
                "@UI.Hidden": { "const": false },
                "@client.synchronize": { "const": true },
                "createdBy": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id"],
                      "properties": {
                        "displayName": { "const": "Alice Hansen" },
                        "id": { "pattern": "^%user_id_pattern%$" }
                      }
                    }
                  }
                },
                "eTag": { "pattern": "%etag_pattern%" },
                "folder": {
                  "const": {}
                },
                "id": { "pattern": "^%share_id_pattern%$" },
                "name": { "const": "folder" },
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
                    "parentReference",
                    "permissions"
                  ],
                  "properties": {
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["id", "displayName"],
                          "properties": {
                            "id": { "pattern": "^%user_id_pattern%$" },
                            "displayName": { "const": "Alice Hansen" }
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
                    "name": { "const": "folder" },
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId", "driveType"],
                      "properties": {
                        "driveId": { "pattern": "^%file_id_pattern%$" },
                        "driveType": { "const": "personal" }
                      }
                    },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["grantedToV2", "id", "invitation", "roles"],
                        "properties": {
                          "id": { "pattern": "^%permissions_id_pattern%$" },
                          "grantedToV2": {
                            "type": "object",
                            "required": ["user"],
                            "properties": {
                              "user": {
                                "type": "object",
                                "properties": {
                                  "displayName": { "const": "Brian Murphy" },
                                  "id": { "pattern": "^%user_id_pattern%$" }
                                },
                                "required": ["displayName", "id"]
                              }
                            }
                          },
                          "invitation": {
                            "type": "object",
                            "properties": {
                              "invitedBy": {
                                "type": "object",
                                "properties": {
                                  "user": {
                                    "type": "object",
                                    "properties": {
                                      "displayName": { "const": "Alice Hansen" },
                                      "id": { "pattern": "^%user_id_pattern%$" }
                                    },
                                    "required": ["displayName", "id"]
                                  }
                                },
                                "required": ["user"]
                              }
                            },
                            "required": ["invitedBy"]
                          },
                          "roles": {
                            "type": "array",
                            "minItems": 1,
                            "maxItems": 1,
                            "items": { "const": "0911d62b-1e3f-4778-8b1b-903b7e4e8476" }
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


  Scenario: sharee lists the folder share shared with EditorWithVersions permission role (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Editor With Versions"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder" in space "new-space"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder               |
      | space           | new-space            |
      | sharee          | Brian                |
      | shareType       | user                 |
      | permissionsRole | Editor With Versions |
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
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "@UI.Hidden",
                "@client.synchronize",
                "eTag",
                "folder",
                "id",
                "lastModifiedDateTime",
                "name",
                "parentReference",
                "remoteItem"
              ],
              "properties": {
                "@UI.Hidden": { "const": false },
                "@client.synchronize": { "const": true },
                "eTag": { "pattern": "%etag_pattern%" },
                "folder": {
                  "const": {}
                },
                "id": { "pattern": "^%share_id_pattern%$" },
                "name": { "const": "folder" },
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
                    "eTag",
                    "folder",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "parentReference",
                    "permissions"
                  ],
                  "properties": {
                    "eTag": { "pattern": "%etag_pattern%" },
                    "file": {
                      "type": "object",
                      "required": ["mimeType"],
                      "properties": {
                        "mimeType": { "const": "text/plain" }
                      }
                    },
                    "id": { "pattern": "^%file_id_pattern%$" },
                    "name": { "const": "folder" },
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId", "driveType"],
                      "properties": {
                        "driveId": { "pattern": "^%file_id_pattern%$" },
                        "driveType": { "const": "project" }
                      }
                    },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["grantedToV2", "id", "invitation", "roles"],
                        "properties": {
                          "id": { "pattern": "^%permissions_id_pattern%$" },
                          "grantedToV2": {
                            "type": "object",
                            "required": ["user"],
                            "properties": {
                              "user": {
                                "type": "object",
                                "properties": {
                                  "displayName": { "const": "Brian Murphy" },
                                  "id": { "pattern": "^%user_id_pattern%$" }
                                },
                                "required": ["displayName", "id"]
                              }
                            }
                          },
                          "invitation": {
                            "type": "object",
                            "properties": {
                              "invitedBy": {
                                "type": "object",
                                "properties": {
                                  "user": {
                                    "type": "object",
                                    "properties": {
                                      "displayName": { "const": "Alice Hansen" },
                                      "id": { "pattern": "^%user_id_pattern%$" }
                                    },
                                    "required": ["displayName", "id"]
                                  }
                                },
                                "required": ["user"]
                              }
                            },
                            "required": ["invitedBy"]
                          },
                          "roles": {
                            "type": "array",
                            "minItems": 1,
                            "maxItems": 1,
                            "items": { "const": "0911d62b-1e3f-4778-8b1b-903b7e4e8476" }
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


  Scenario Outline: sharee checks versions after updating the permission role of a shared file from other roles to FileEditorWithVersions (Personal Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role 'File Editor With Versions'
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt       |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole | File Editor With Versions |
      | space           | Personal                  |
      | resource        | textfile.txt              |
    And user "Brian" has uploaded file with content "updated content" to "Shares/textfile.txt"
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | permissions-role |
      | File Editor      |
      | Viewer           |


  Scenario Outline: sharee tries to check versions after updating the permission role of a shared file from FileEditorWithVersions to other roles (Personal Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role 'File Editor With Versions'
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt              |
      | space           | Personal                  |
      | sharee          | Brian                     |
      | shareType       | user                      |
      | permissionsRole | File Editor With Versions |
    And user "Brian" has uploaded file with content "updated content" to "Shares/textfile.txt"
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole | <permissions-role> |
      | space           | Personal           |
      | resource        | textfile.txt       |
    When user "Brian" tries to get versions of file "textfile.txt" from "Alice"
    Then the HTTP status code should be "403"
    Examples:
      | permissions-role |
      | File Editor      |
      | Viewer           |


  Scenario Outline: sharee checks versions after updating the permission role of the shared folder from other roles to EditorWithVersions (Personal Space)
    Given the administrator has enabled the permissions role 'Editor With Versions'
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "to share" to "folderToShare/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole | Editor With Versions |
      | space           | Personal             |
      | resource        | folderToShare        |
    And user "Brian" has uploaded file with content "updated content" to "Shares/folderToShare/textfile.txt"
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | permissions-role |
      | Editor           |
      | Viewer           |


  Scenario Outline: sharee tries to check versions after updating the permission role of the shared folder from EditorWithVersions to other roles (Personal Space)
    Given the administrator has enabled the permissions role 'Editor With Versions'
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "to share" to "folderToShare/textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare        |
      | space           | Personal             |
      | sharee          | Brian                |
      | shareType       | user                 |
      | permissionsRole | Editor With Versions |
    And user "Brian" has uploaded file with content "updated content" to "Shares/folderToShare/textfile.txt"
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole | <permissions-role> |
      | space           | Personal           |
      | resource        | folderToShare      |
    When user "Brian" tries to get versions of file "folderToShare/textfile.txt" from "Alice"
    Then the HTTP status code should be "403"
    Examples:
      | permissions-role |
      | Editor           |
      | Viewer           |


  Scenario Outline: sharee checks versions after updating the permission role of a shared file from other roles to FileEditorWithVersions (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role 'File Editor With Versions'
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "to share" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt       |
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole | File Editor With Versions |
      | space           | new-space                 |
      | resource        | textfile.txt              |
    And user "Brian" has uploaded file with content "updated content" to "Shares/textfile.txt"
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | permissions-role |
      | File Editor      |
      | Viewer           |


  Scenario Outline: sharee tries to check versions after updating the permission role of a shared file from FileEditorWithVersions to other roles (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role 'File Editor With Versions'
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "to share" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt              |
      | space           | new-space                 |
      | sharee          | Brian                     |
      | shareType       | user                      |
      | permissionsRole | File Editor With Versions |
    And user "Brian" has uploaded file with content "updated content" to "Shares/textfile.txt"
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole | <permissions-role> |
      | space           | new-space          |
      | resource        | textfile.txt       |
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "403"
    Examples:
      | permissions-role |
      | File Editor      |
      | Viewer           |


  Scenario Outline: sharee checks versions after updating the permission role of the shared folder from other roles to EditorWithVersions (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role 'Editor With Versions'
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "to share" to "folderToShare/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare      |
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole | Editor With Versions |
      | space           | new-space            |
      | resource        | folderToShare        |
    And user "Brian" has uploaded file with content "updated content" to "Shares/folderToShare/textfile.txt"
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | permissions-role |
      | Editor           |
      | Viewer           |


  Scenario Outline: sharee tries to check versions after updating the permission role of the shared folder from EditorWithVersions to other roles (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role 'Editor With Versions'
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "to share" to "folderToShare/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare        |
      | space           | new-space            |
      | sharee          | Brian                |
      | shareType       | user                 |
      | permissionsRole | Editor With Versions |
    And user "Brian" has uploaded file with content "updated content" to "Shares/folderToShare/textfile.txt"
    And user "Alice" has updated the last resource share with the following properties:
      | permissionsRole | <permissions-role> |
      | space           | new-space          |
      | resource        | folderToShare      |
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "403"
    Examples:
      | permissions-role |
      | Editor           |
      | Viewer           |
