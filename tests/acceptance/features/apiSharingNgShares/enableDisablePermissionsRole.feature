@env-config
Feature: enable disable permissions role
  As a user
  I want to enable/disable permissions role on shared resources
  So that I can control the accessibility of shared resources to sharee

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: users list the shares shared with Secure Viewer after the role is disabled (Personal Space)
    Given the administrator has enabled the permissions role "Secure Viewer"
    And user "Alice" has uploaded file with content "some content" to "textfile.txt"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "hello world" to "folderToShare/textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>    |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And user "Brian" has a share "<resource>" synced
    And the administrator has disabled the permissions role "Secure Viewer"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "<resource>" with the following data:
      """
      {
        "type": "object",
        "required": [
          "parentReference",
          "permissions",
          "name"
        ],
        "properties": {
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "@libre.graph.permissions.actions",
                "grantedToV2",
                "id",
                "invitation"
              ],
              "properties": {
                "@libre.graph.permissions.actions": {
                  "const": [
                      "libre.graph/driveItem/path/read",
                      "libre.graph/driveItem/children/read",
                      "libre.graph/driveItem/basic/read"
                  ]
                },
                "roles": { "const": null }
              }
            }
          }
        }
      }
      """
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
                "<resource-type>",
                "id",
                "lastModifiedDateTime",
                "name",
                "parentReference",
                "remoteItem"
              ],
              "properties": {
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "createdBy",
                    "eTag",
                    "<resource-type>",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "parentReference",
                    "permissions"
                  ],
                  "properties": {
                    "name": { "const": "<resource>" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": [
                          "@libre.graph.permissions.actions",
                          "grantedToV2",
                          "id",
                          "invitation"
                        ],
                        "properties": {
                          "@libre.graph.permissions.actions": {
                            "const": [
                                "libre.graph/driveItem/path/read",
                                "libre.graph/driveItem/children/read",
                                "libre.graph/driveItem/basic/read"
                            ]
                          },
                          "roles": { "const": null }
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
    And user "Brian" should have a share "<resource>" shared by user "Alice" from space "Personal"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "<resource>" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "<resource>" with these key and value pairs:
      | key            | value      |
      | oc:name        | <resource> |
      | oc:permissions | SX         |
    And user "Brian" should not be able to download file "<resource-to-download>" from space "Shares"
    Examples:
      | resource      | resource-type | resource-to-download        |
      | textfile.txt  | file          | textfile.txt                |
      | folderToShare | folder        | folderToShare/textfile1.txt |


  Scenario: users list the shares shared with Denied after the role is disabled (Personal Space)
    Given the administrator has enabled the permissions role "Denied"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "hello world" to "folderToShare/textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Denied        |
    And the administrator has disabled the permissions role "Denied"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "folderToShare" with the following data:
      """
      {
        "type": "object",
        "required": [
          "parentReference",
          "permissions",
          "name"
        ],
        "properties": {
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "@libre.graph.permissions.actions",
                "grantedToV2",
                "id",
                "invitation"
              ],
              "properties": {
                "@libre.graph.permissions.actions": {
                  "const": ["none"]
                },
                "roles": { "const": null }
              }
            }
          }
        }
      }
      """
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
            "minItems": 0,
            "maxItems": 0
          }
        }
      }
      """
    And user "Brian" should not have a share "folderToShare" shared by user "Alice" from space "Personal"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "folderToShare" with depth "0" using the WebDAV API
    Then the HTTP status code should be "404"
    And user "Brian" should not be able to download file "folderToShare/textfile1.txt" from space "Shares"


  Scenario Outline: users list the shares shared with Secure Viewer after the role is disabled (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "folderToShare/textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>    |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And user "Brian" has a share "<resource>" synced
    And the administrator has disabled the permissions role "Secure Viewer"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "<resource>" with the following data:
      """
      {
        "type": "object",
        "required": [
          "parentReference",
          "permissions",
          "name"
        ],
        "properties": {
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "@libre.graph.permissions.actions",
                "grantedToV2",
                "id",
                "invitation"
              ],
              "properties": {
                "@libre.graph.permissions.actions": {
                  "const": [
                      "libre.graph/driveItem/path/read",
                      "libre.graph/driveItem/children/read",
                      "libre.graph/driveItem/basic/read"
                  ]
                },
                "roles": { "const": null }
              }
            }
          }
        }
      }
      """
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
                "<resource-type>",
                "id",
                "lastModifiedDateTime",
                "name",
                "parentReference",
                "remoteItem"
              ],
              "properties": {
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "eTag",
                    "<resource-type>",
                    "id",
                    "lastModifiedDateTime",
                    "name",
                    "parentReference",
                    "permissions"
                  ],
                  "properties": {
                    "name": { "const": "<resource>" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": [
                          "@libre.graph.permissions.actions",
                          "grantedToV2",
                          "id",
                          "invitation"
                        ],
                        "properties": {
                          "@libre.graph.permissions.actions": {
                            "const": [
                                "libre.graph/driveItem/path/read",
                                "libre.graph/driveItem/children/read",
                                "libre.graph/driveItem/basic/read"
                            ]
                          },
                          "roles": { "const": null }
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
    And user "Brian" should have a share "<resource>" shared by user "Alice" from space "new-space"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "<resource>" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "<resource>" with these key and value pairs:
      | key            | value      |
      | oc:name        | <resource> |
      | oc:permissions | SX         |
    And user "Brian" should not be able to download file "<resource-to-download>" from space "Shares"
    Examples:
      | resource      | resource-type | resource-to-download        |
      | textfile.txt  | file          | textfile.txt                |
      | folderToShare | folder        | folderToShare/textfile1.txt |


  Scenario: users list the shares shared with Denied after the role is disabled (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "folderToShare/textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Denied        |
    And the administrator has disabled the permissions role "Denied"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "folderToShare" with the following data:
      """
      {
        "type": "object",
        "required": [
          "parentReference",
          "permissions",
          "name"
        ],
        "properties": {
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "@libre.graph.permissions.actions",
                "grantedToV2",
                "id",
                "invitation"
              ],
              "properties": {
                "@libre.graph.permissions.actions": {
                  "const": ["none"]
                },
                "roles": { "const": null }
              }
            }
          }
        }
      }
      """
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
            "minItems": 0,
            "maxItems": 0
          }
        }
      }
      """
    And user "Brian" should not have a share "folderToShare" shared by user "Alice" from space "Personal"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "folderToShare" with depth "0" using the WebDAV API
    Then the HTTP status code should be "404"
    And user "Brian" should not be able to download file "folderToShare/textfile1.txt" from space "Shares"


  Scenario: sharee lists drives after the share role Space Editor Without Versions has been disabled
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Space Editor Without Versions"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "textfile1.txt"
    And user "Alice" has sent the following space share invitation:
      | space           | new-space                     |
      | sharee          | Brian                         |
      | shareType       | user                          |
      | permissionsRole | Space Editor Without Versions |
    And the administrator has disabled the permissions role "Space Editor Without Versions"
    When user "Brian" lists all available spaces via the Graph API
    Then the HTTP status code should be "200"
    And the JSON response should contain space called "new-space" and match
      """
      {
        "type": "object",
        "required": [
          "driveType",
          "driveAlias",
          "name",
          "id",
          "quota",
          "root",
          "webUrl"
        ],
        "properties": {
          "root": {
            "type": "object",
            "required": [
              "eTag",
              "id",
              "permissions",
              "webDavUrl"
            ],
            "properties": {
              "permissions": {
                "type": "array",
                "minItems": 2,
                "maxItems": 2,
                "uniqueItems": true,
                "items": {
                  "oneOf": [
                    {
                      "type": "object",
                      "required": ["grantedToV2", "roles"],
                      "properties": {
                        "grantedToV2": {
                          "type": "object",
                          "required": ["user"],
                          "properties": {
                            "user" : {
                              "type": "object",
                              "required": ["@libre.graph.userType", "displayName", "id"],
                              "properties": {
                                "@libre.graph.userType": { "const": "Member" },
                                "displayName": { "const": "Alice Hansen" },
                                "id": { "pattern": "^%user_id_pattern%$" }
                              }
                            }
                          }
                        },
                        "roles": { "pattern": "^%role_id_pattern%$" }
                      }
                    },
                    {
                      "type": "object",
                      "required": ["@libre.graph.permissions.actions", "grantedToV2"],
                      "properties": {
                        "@libre.graph.permissions.actions": {
                          "const": [
                            "libre.graph/driveItem/children/create",
                            "libre.graph/driveItem/standard/delete",
                            "libre.graph/driveItem/path/read",
                            "libre.graph/driveItem/quota/read",
                            "libre.graph/driveItem/content/read",
                            "libre.graph/driveItem/upload/create",
                            "libre.graph/driveItem/permissions/read",
                            "libre.graph/driveItem/children/read",
                            "libre.graph/driveItem/deleted/read",
                            "libre.graph/driveItem/path/update",
                            "libre.graph/driveItem/deleted/update",
                            "libre.graph/driveItem/basic/read"
                          ]
                        },
                        "grantedToV2": {
                          "type": "object",
                          "required": ["user"],
                          "properties": {
                            "user": {
                              "type": "object",
                              "required": ["@libre.graph.userType", "displayName", "id"],
                              "properties": {
                                "@libre.graph.userType": { "const": "Member" },
                                "displayName": { "const": "Brian Murphy" },
                                "id": { "pattern": "^%user_id_pattern%$" }
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
        }
      }
      """
    When user "Brian" sends PROPFIND request to space "new-space" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a space "new-space" with these key and value pairs:
      | key            | value     |
      | oc:name        | new-space |
      | oc:permissions | DNVCK     |
    And user "Brian" should be able to download file "textfile1.txt" from space "new-space"


  Scenario Outline: update share role of a file to an existing role after assigned share role Secure Viewer is disabled (Personal Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And user "Alice" has uploaded file with content "some content" to "textfile.txt"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "hello world" to "folderToShare/textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>    |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And user "Brian" has a share "<resource>" synced
    And the administrator has disabled the permissions role "Secure Viewer"
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <permissions-role> |
      | space           | Personal           |
      | resource        | <resource>         |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "grantedToV2",
          "id",
          "roles"
        ],
        "properties": {
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "const": "<role-id>"
            }
          }
        }
      }
      """
    And user "Brian" should have a share "<resource>" shared by user "Alice" from space "Personal"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "<resource>" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "<resource>" with these key and value pairs:
      | key            | value         |
      | oc:name        | <resource>    |
      | oc:permissions | <permissions> |
    And user "Brian" should be able to download file "<resource-to-download>" from space "Shares"
    Examples:
      | resource      | permissions-role | permissions | resource-to-download        | role-id                              |
      | textfile.txt  | Viewer           | S           | textfile.txt                | b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5 |
      | folderToShare | Uploader         | SCK         | folderToShare/textfile1.txt | 1c996275-f1c9-4e71-abdf-a42f6495e960 |
      | folderToShare | Editor           | SDNVCK      | folderToShare/textfile1.txt | fb6c3e19-e378-47e5-b277-9732f9de6e21 |
      | folderToShare | Viewer           | S           | folderToShare/textfile1.txt | b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5 |


  Scenario Outline: update share role of a file to an existing role after assigned share role Secure Viewer is disabled (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "folderToShare/textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>    |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And user "Brian" has a share "<resource>" synced
    And the administrator has disabled the permissions role "Secure Viewer"
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <permissions-role> |
      | space           | new-space          |
      | resource        | <resource>         |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "grantedToV2",
          "id",
          "roles"
        ],
        "properties": {
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "const": "<role-id>"
            }
          }
        }
      }
      """
    And user "Brian" should have a share "<resource>" shared by user "Alice" from space "new-space"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "<resource>" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "<resource>" with these key and value pairs:
      | key            | value         |
      | oc:name        | <resource>    |
      | oc:permissions | <permissions> |
    And user "Brian" should be able to download file "<resource-to-download>" from space "Shares"
    Examples:
      | resource      | permissions-role | permissions | resource-to-download        | role-id                              |
      | textfile.txt  | File Editor      | SW          | textfile.txt                | 2d00ce52-1fc2-4dbc-8b95-a73b73395f5a |
      | textfile.txt  | Viewer           | S           | textfile.txt                | b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5 |
      | folderToShare | Uploader         | SCK         | folderToShare/textfile1.txt | 1c996275-f1c9-4e71-abdf-a42f6495e960 |
      | folderToShare | Editor           | SDNVCK      | folderToShare/textfile1.txt | fb6c3e19-e378-47e5-b277-9732f9de6e21 |
      | folderToShare | Viewer           | S           | folderToShare/textfile1.txt | b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5 |


  Scenario Outline: update share role of a project space to an existing role after assigned share role is disabled
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Space Editor Without Versions"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following space share invitation:
      | space           | new-space                     |
      | sharee          | Brian                         |
      | shareType       | user                          |
      | permissionsRole | Space Editor Without Versions |
    And the administrator has disabled the permissions role "Space Editor Without Versions"
    When user "Alice" updates the last drive share with the following using root endpoint of the Graph API:
      | permissionsRole | <permissions-role> |
      | space           | new-space          |
      | shareType       | user               |
      | sharee          | Brian              |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "grantedToV2",
          "id",
          "roles"
        ],
        "properties": {
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "const": "<role-id>"
            }
          }
        }
      }
      """
    When user "Brian" sends PROPFIND request to space "new-space" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a space "new-space" with these key and value pairs:
      | key            | value         |
      | oc:name        | new-space     |
      | oc:permissions | <permissions> |
    And user "Brian" should be able to download file "textfile.txt" from space "new-space"
    Examples:
      | permissions-role | permissions | role-id                              |
      | Space Viewer     |             | a8d5fe5e-96e3-418d-825b-534dbdf22b99 |
      | Space Editor     | DNVCK       | 58c63c02-1d89-4572-916a-870abc5a1b7d |
      | Manager          | RDNVCKZP    | 312c0871-5ef7-4b3a-85b6-0e4074c64049 |


  Scenario: users list the folder shares shared with Editor role after the role is disabled (Personal Space)
    Given using spaces DAV path
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "hello world" to "folderToShare/textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Editor        |
    And user "Brian" has a share "folderToShare" synced
    And the administrator has disabled the permissions role "Editor"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "folderToShare" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference", "permissions", "name"],
        "properties": {
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["@libre.graph.permissions.actions", "grantedToV2", "id", "invitation"],
              "properties": {
                "@libre.graph.permissions.actions": {
                  "const": [
                    "libre.graph/driveItem/children/create",
                    "libre.graph/driveItem/standard/delete",
                    "libre.graph/driveItem/path/read",
                    "libre.graph/driveItem/quota/read",
                    "libre.graph/driveItem/content/read",
                    "libre.graph/driveItem/upload/create",
                    "libre.graph/driveItem/children/read",
                    "libre.graph/driveItem/deleted/read",
                    "libre.graph/driveItem/path/update",
                    "libre.graph/driveItem/deleted/update",
                    "libre.graph/driveItem/basic/read"
                  ]
                },
                "roles": { "const": null }
              }
            }
          }
        }
      }
      """
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
              "required": ["@UI.Hidden", "@client.synchronize", "createdBy", "eTag", "folder", "id", "lastModifiedDateTime", "name", "parentReference", "remoteItem"],
              "properties": {
                "remoteItem": {
                  "type": "object",
                  "required": ["createdBy", "eTag", "folder", "id", "lastModifiedDateTime", "name", "parentReference", "permissions"],
                  "properties": {
                    "name": { "const": "folderToShare" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["@libre.graph.permissions.actions", "grantedToV2", "id", "invitation"],
                        "properties": {
                          "@libre.graph.permissions.actions": {
                            "const": [
                              "libre.graph/driveItem/children/create",
                              "libre.graph/driveItem/standard/delete",
                              "libre.graph/driveItem/path/read",
                              "libre.graph/driveItem/quota/read",
                              "libre.graph/driveItem/content/read",
                              "libre.graph/driveItem/upload/create",
                              "libre.graph/driveItem/children/read",
                              "libre.graph/driveItem/deleted/read",
                              "libre.graph/driveItem/path/update",
                              "libre.graph/driveItem/deleted/update",
                              "libre.graph/driveItem/basic/read"
                            ]
                          },
                          "roles": { "const": null }
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
    And user "Brian" should have a share "folderToShare" shared by user "Alice" from space "Personal"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "folderToShare" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare" with these key and value pairs:
      | key            | value         |
      | oc:name        | folderToShare |
      | oc:permissions | SDNVCK        |
    And user "Brian" should be able to upload file "filesForUpload/davtest.txt" to "Shares/folderToShare/textfile.txt"
    And for user "Alice" the content of the file "folderToShare/textfile.txt" of the space "Personal" should be "Dav-Test"


  Scenario: users list the folder shares shared with Editor role after the role is disabled (Project Space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "folderToShare/textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Editor        |
    And user "Brian" has a share "folderToShare" synced
    And the administrator has disabled the permissions role "Editor"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "folderToShare" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference", "permissions", "name"],
        "properties": {
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["@libre.graph.permissions.actions", "grantedToV2", "id", "invitation"],
              "properties": {
                "@libre.graph.permissions.actions": {
                  "const": [
                    "libre.graph/driveItem/children/create",
                    "libre.graph/driveItem/standard/delete",
                    "libre.graph/driveItem/path/read",
                    "libre.graph/driveItem/quota/read",
                    "libre.graph/driveItem/content/read",
                    "libre.graph/driveItem/upload/create",
                    "libre.graph/driveItem/children/read",
                    "libre.graph/driveItem/deleted/read",
                    "libre.graph/driveItem/path/update",
                    "libre.graph/driveItem/deleted/update",
                    "libre.graph/driveItem/basic/read"
                  ]
                },
                "roles": { "const": null }
              }
            }
          }
        }
      }
      """
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
              "required": ["@UI.Hidden", "@client.synchronize", "eTag", "folder", "id", "lastModifiedDateTime", "name", "parentReference", "remoteItem"],
              "properties": {
                "remoteItem": {
                  "type": "object",
                  "required": ["eTag", "folder", "id", "lastModifiedDateTime", "name", "parentReference", "permissions"],
                  "properties": {
                    "name": { "const": "folderToShare" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["@libre.graph.permissions.actions", "grantedToV2", "id", "invitation"],
                        "properties": {
                          "@libre.graph.permissions.actions": {
                            "const": [
                              "libre.graph/driveItem/children/create",
                              "libre.graph/driveItem/standard/delete",
                              "libre.graph/driveItem/path/read",
                              "libre.graph/driveItem/quota/read",
                              "libre.graph/driveItem/content/read",
                              "libre.graph/driveItem/upload/create",
                              "libre.graph/driveItem/children/read",
                              "libre.graph/driveItem/deleted/read",
                              "libre.graph/driveItem/path/update",
                              "libre.graph/driveItem/deleted/update",
                              "libre.graph/driveItem/basic/read"
                            ]
                          },
                          "roles": { "const": null }
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
    And user "Brian" should have a share "folderToShare" shared by user "Alice" from space "new-space"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "folderToShare" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "folderToShare" with these key and value pairs:
      | key            | value         |
      | oc:name        | folderToShare |
      | oc:permissions | SDNVCK        |
    And user "Brian" should be able to upload file "filesForUpload/davtest.txt" to "Shares/folderToShare/textfile.txt"
    And for user "Alice" the content of the file "folderToShare/textfile.txt" of the space "new-space" should be "Dav-Test"


  Scenario: users list the file shares shared with Editor role after the role is disabled (Personal Space)
    Given using spaces DAV path
    And user "Alice" has uploaded file with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | File Editor  |
    And user "Brian" has a share "textfile.txt" synced
    And the administrator has disabled the permissions role "File Editor"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "textfile.txt" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference", "permissions", "name"],
        "properties": {
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["@libre.graph.permissions.actions", "grantedToV2", "id", "invitation"],
              "properties": {
                "@libre.graph.permissions.actions": {
                  "const": [
                    "libre.graph/driveItem/path/read",
                    "libre.graph/driveItem/quota/read",
                    "libre.graph/driveItem/content/read",
                    "libre.graph/driveItem/upload/create",
                    "libre.graph/driveItem/children/read",
                    "libre.graph/driveItem/deleted/read",
                    "libre.graph/driveItem/deleted/update",
                    "libre.graph/driveItem/basic/read"
                  ]
                },
                "roles": { "const": null }
              }
            }
          }
        }
      }
      """
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
              "required": ["@UI.Hidden", "@client.synchronize", "createdBy", "eTag", "file", "id", "lastModifiedDateTime", "name", "parentReference", "remoteItem", "size"],
              "properties": {
                "remoteItem": {
                  "type": "object",
                  "required": ["createdBy", "eTag", "file", "id", "lastModifiedDateTime", "name", "parentReference", "permissions"],
                  "properties": {
                    "name": { "const": "textfile.txt" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["@libre.graph.permissions.actions", "grantedToV2", "id", "invitation"],
                        "properties": {
                          "@libre.graph.permissions.actions": {
                            "const": [
                              "libre.graph/driveItem/path/read",
                              "libre.graph/driveItem/quota/read",
                              "libre.graph/driveItem/content/read",
                              "libre.graph/driveItem/upload/create",
                              "libre.graph/driveItem/children/read",
                              "libre.graph/driveItem/deleted/read",
                              "libre.graph/driveItem/deleted/update",
                              "libre.graph/driveItem/basic/read"
                            ]
                          },
                          "roles": { "const": null }
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
    And user "Brian" should have a share "textfile.txt" shared by user "Alice" from space "Personal"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "textfile.txt" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "textfile.txt" with these key and value pairs:
      | key            | value        |
      | oc:name        | textfile.txt |
      | oc:permissions | SW           |
    And user "Brian" should be able to upload file "filesForUpload/davtest.txt" to "Shares/textfile.txt"
    And for user "Alice" the content of the file "textfile.txt" of the space "Personal" should be "Dav-Test"


  Scenario: users list the file shares shared with Editor role after the role is disabled (Project Space)
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt |
      | space           | new-space    |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | File Editor  |
    And user "Brian" has a share "textfile.txt" synced
    And the administrator has disabled the permissions role "File Editor"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "textfile.txt" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference", "permissions", "name"],
        "properties": {
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["@libre.graph.permissions.actions", "grantedToV2", "id", "invitation"],
              "properties": {
                "@libre.graph.permissions.actions": {
                  "const": [
                    "libre.graph/driveItem/path/read",
                    "libre.graph/driveItem/quota/read",
                    "libre.graph/driveItem/content/read",
                    "libre.graph/driveItem/upload/create",
                    "libre.graph/driveItem/children/read",
                    "libre.graph/driveItem/deleted/read",
                    "libre.graph/driveItem/deleted/update",
                    "libre.graph/driveItem/basic/read"
                  ]
                },
                "roles": { "const": null }
              }
            }
          }
        }
      }
      """
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
              "required": ["@UI.Hidden", "@client.synchronize", "eTag", "file", "id", "lastModifiedDateTime", "name", "parentReference", "remoteItem", "size"],
              "properties": {
                "remoteItem": {
                  "type": "object",
                  "required": ["eTag", "file", "id", "lastModifiedDateTime", "name", "parentReference", "permissions"],
                  "properties": {
                    "name": { "const": "textfile.txt" },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["@libre.graph.permissions.actions", "grantedToV2", "id", "invitation"],
                        "properties": {
                          "@libre.graph.permissions.actions": {
                            "const": [
                              "libre.graph/driveItem/path/read",
                              "libre.graph/driveItem/quota/read",
                              "libre.graph/driveItem/content/read",
                              "libre.graph/driveItem/upload/create",
                              "libre.graph/driveItem/children/read",
                              "libre.graph/driveItem/deleted/read",
                              "libre.graph/driveItem/deleted/update",
                              "libre.graph/driveItem/basic/read"
                            ]
                          },
                          "roles": { "const": null }
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
    And user "Brian" should have a share "textfile.txt" shared by user "Alice" from space "new-space"
    When user "Brian" sends PROPFIND request from the space "Shares" to the resource "textfile.txt" with depth "0" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" the PROPFIND response should contain a resource "textfile.txt" with these key and value pairs:
      | key            | value        |
      | oc:name        | textfile.txt |
      | oc:permissions | SW           |
    And user "Brian" should be able to upload file "filesForUpload/davtest.txt" to "Shares/textfile.txt"
    And for user "Alice" the content of the file "textfile.txt" of the space "new-space" should be "Dav-Test"
