@env-config
Feature: ListGrants role
  As a user
  I want to share resources with listGrants role
  So that sharee can view activities and grants list of shared resources

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: user shares personal space file with ListGrants role
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
              "required": ["createdDateTime","id","roles","grantedToV2"],
              "properties": {
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id","displayName"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
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
    And for user "Brian" file "textfile1.txt" should have the following shares:
      | sharee | shareType | permissionsRole    |
      | Brian  | user      | <permissions-role> |
    And for user "Brian" file "textfile1.txt" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}    |
      | {user} shared {resource} with {sharee} |
    Examples:
      | permissions-role            |
      | Viewer With ListGrants      |
      | File Editor With ListGrants |


  Scenario Outline: user shares personal space folder with ListGrants role
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | FolderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
              "required": ["createdDateTime","id","roles","grantedToV2"],
              "properties": {
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id","displayName"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
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
    And for user "Brian" folder "FolderToShare" should have the following shares:
      | sharee | shareType | permissionsRole    |
      | Brian  | user      | <permissions-role> |
    And for user "Brian" folder "FolderToShare" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}    |
      | {user} shared {resource} with {sharee} |
    Examples:
      | permissions-role       |
      | Viewer With ListGrants |
      | Editor With ListGrants |


  Scenario Outline: user shares project space file with ListGrants role
    Given the administrator has enabled the permissions role "<permissions-role>"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "textfile1.txt"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | textfile1.txt      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
              "required": ["createdDateTime","id","roles","grantedToV2"],
              "properties": {
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id","displayName"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
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
    And for user "Brian" file "textfile1.txt" should have the following shares:
      | sharee | shareType | permissionsRole    |
      | Brian  | user      | <permissions-role> |
    And for user "Brian" file "textfile1.txt" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}    |
      | {user} shared {resource} with {sharee} |
    Examples:
      | permissions-role            |
      | Viewer With ListGrants      |
      | File Editor With ListGrants |


  Scenario Outline: user shares project space folder with ListGrants role
    Given the administrator has enabled the permissions role "<permissions-role>"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    When user "Alice" sends the following resource share invitation using the Graph API:
      | resource        | FolderToShare      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
              "required": ["createdDateTime","id","roles","grantedToV2"],
              "properties": {
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                },
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["id","displayName"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
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
    And for user "Brian" folder "FolderToShare" should have the following shares:
      | sharee | shareType | permissionsRole    |
      | Brian  | user      | <permissions-role> |
    And for user "Brian" folder "FolderToShare" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}    |
      | {user} shared {resource} with {sharee} |
    Examples:
      | permissions-role       |
      | Viewer With ListGrants |
      | Editor With ListGrants |


  Scenario Outline: sharer updates shared file roles to ListGrants roles (Personal space)
    Given the administrator has enabled the permissions role "<new-permissions-role>"
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | Personal               |
      | resource        | textfile1.txt          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["grantedToV2","id","roles"],
        "properties": {
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": "Brian Murphy"},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "id": {"pattern": "^%permissions_id_pattern%$"},
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {"pattern": "^%role_id_pattern%$"}
          }
        }
      }
      """
    And for user "Brian" file "textfile1.txt" should have the following shares:
      | sharee | shareType | permissionsRole        |
      | Brian  | user      | <new-permissions-role> |
    And for user "Brian" file "textfile1.txt" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}    |
      | {user} shared {resource} with {sharee} |
    Examples:
      | permissions-role | new-permissions-role        |
      | Viewer           | Viewer With ListGrants      |
      | File Editor      | Viewer With ListGrants      |
      | Viewer           | File Editor With ListGrants |
      | File Editor      | File Editor With ListGrants |


  Scenario Outline: sharer updates shared folder roles to ListGrants roles (Personal space)
    Given the administrator has enabled the permissions role "<new-permissions-role>"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FolderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | Personal               |
      | resource        | FolderToShare          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["grantedToV2","id","roles"],
        "properties": {
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": "Brian Murphy"},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "id": {"pattern": "^%permissions_id_pattern%$"},
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {"pattern": "^%role_id_pattern%$"}
          }
        }
      }
      """
    And for user "Brian" folder "FolderToShare" should have the following shares:
      | sharee | shareType | permissionsRole        |
      | Brian  | user      | <new-permissions-role> |
    And for user "Brian" folder "FolderToShare" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}    |
      | {user} shared {resource} with {sharee} |
    Examples:
      | permissions-role | new-permissions-role   |
      | Viewer           | Viewer With ListGrants |
      | Editor           | Viewer With ListGrants |
      | Uploader         | Viewer With ListGrants |
      | Viewer           | Editor With ListGrants |
      | Editor           | Editor With ListGrants |
      | Uploader         | Editor With ListGrants |


  Scenario Outline: sharer updates shared file roles to ListGrants roles (Project space)
    Given the administrator has enabled the permissions role "<new-permissions-role>"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile1.txt      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | NewSpace               |
      | resource        | textfile1.txt          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["grantedToV2","id","roles"],
        "properties": {
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": "Brian Murphy"},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "id": {"pattern": "^%permissions_id_pattern%$"},
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {"pattern": "^%role_id_pattern%$"}
          }
        }
      }
      """
    And for user "Brian" file "textfile1.txt" should have the following shares:
      | sharee | shareType | permissionsRole        |
      | Brian  | user      | <new-permissions-role> |
    And for user "Brian" file "textfile1.txt" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}       |
      | {user} shared {resource} with {sharee}    |
      | {user} updated {field} for the {resource} |
    Examples:
      | permissions-role | new-permissions-role        |
      | Viewer           | Viewer With ListGrants      |
      | File Editor      | Viewer With ListGrants      |
      | Viewer           | File Editor With ListGrants |
      | File Editor      | File Editor With ListGrants |


  Scenario Outline: sharer updates shared folder roles to ListGrants roles (Project space)
    Given the administrator has enabled the permissions role "<new-permissions-role>"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FolderToShare      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | NewSpace               |
      | resource        | FolderToShare          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["grantedToV2","id","roles"],
        "properties": {
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": "Brian Murphy"},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "id": {"pattern": "^%permissions_id_pattern%$"},
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {"pattern": "^%role_id_pattern%$"}
          }
        }
      }
      """
    And for user "Brian" folder "FolderToShare" should have the following shares:
      | sharee | shareType | permissionsRole        |
      | Brian  | user      | <new-permissions-role> |
    And for user "Brian" folder "FolderToShare" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}       |
      | {user} shared {resource} with {sharee}    |
      | {user} updated {field} for the {resource} |
    Examples:
      | permissions-role | new-permissions-role   |
      | Viewer           | Viewer With ListGrants |
      | Editor           | Viewer With ListGrants |
      | Uploader         | Viewer With ListGrants |
      | Viewer           | Editor With ListGrants |
      | Editor           | Editor With ListGrants |
      | Uploader         | Editor With ListGrants |


  Scenario Outline: sharer updates shared file roles from ListGrants roles to other roles (Personal space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | Personal               |
      | resource        | textfile1.txt          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["grantedToV2","id","roles"],
        "properties": {
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": "Brian Murphy"},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "id": {"pattern": "^%permissions_id_pattern%$"},
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {"pattern": "^%role_id_pattern%$"}
          }
        }
      }
      """
    And for user "Brian" file "textfile1.txt" should have the following shares:
      | sharee | shareType | permissionsRole        |
      | Brian  | user      | <new-permissions-role> |
    And for user "Brian" file "textfile1.txt" of the space "Shares" should not have any activity
    Examples:
      | permissions-role            | new-permissions-role |
      | Viewer With ListGrants      | Viewer               |
      | Viewer With ListGrants      | File Editor          |
      | File Editor With ListGrants | Viewer               |
      | File Editor With ListGrants | File Editor          |


  Scenario Outline: sharer updates shared folder roles from ListGrants roles to other roles (Personal space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FolderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | Personal               |
      | resource        | FolderToShare          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["grantedToV2","id","roles"],
        "properties": {
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": "Brian Murphy"},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "id": {"pattern": "^%permissions_id_pattern%$"},
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {"pattern": "^%role_id_pattern%$"}
          }
        }
      }
      """
    And for user "Brian" folder "FolderToShare" should have the following shares:
      | sharee | shareType | permissionsRole        |
      | Brian  | user      | <new-permissions-role> |
    And for user "Brian" folder "FolderToShare" of the space "Shares" should not have any activity
    Examples:
      | permissions-role       | new-permissions-role |
      | Viewer With ListGrants | Viewer               |
      | Viewer With ListGrants | Editor               |
      | Viewer With ListGrants | Uploader             |
      | Editor With ListGrants | Viewer               |
      | Editor With ListGrants | Editor               |
      | Editor With ListGrants | Uploader             |


  Scenario Outline: sharer updates shared file roles from ListGrants roles to other roles (Project space)
    Given the administrator has enabled the permissions role "<new-permissions-role>"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "textfile1.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile1.txt      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | NewSpace               |
      | resource        | textfile1.txt          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["grantedToV2","id","roles"],
        "properties": {
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": "Brian Murphy"},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "id": {"pattern": "^%permissions_id_pattern%$"},
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {"pattern": "^%role_id_pattern%$"}
          }
        }
      }
      """
    And for user "Brian" file "textfile1.txt" should have the following shares:
      | sharee | shareType | permissionsRole        |
      | Brian  | user      | <new-permissions-role> |
    And for user "Brian" file "textfile1.txt" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}       |
      | {user} shared {resource} with {sharee}    |
      | {user} updated {field} for the {resource} |
    Examples:
      | new-permissions-role        | permissions-role |
      | Viewer With ListGrants      | Viewer           |
      | Viewer With ListGrants      | File Editor      |
      | File Editor With ListGrants | Viewer           |
      | File Editor With ListGrants | File Editor      |


  Scenario Outline: sharer updates shared folder roles to ListGrants roles to other roles (Project space)
    Given the administrator has enabled the permissions role "<new-permissions-role>"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FolderToShare      |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" updates the last resource share with the following properties using the Graph API:
      | permissionsRole | <new-permissions-role> |
      | space           | NewSpace               |
      | resource        | FolderToShare          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["grantedToV2","id","roles"],
        "properties": {
          "grantedToV2": {
            "type": "object",
            "required": ["user"],
            "properties":{
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": "Brian Murphy"},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "id": {"pattern": "^%permissions_id_pattern%$"},
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {"pattern": "^%role_id_pattern%$"}
          }
        }
      }
      """
    And for user "Brian" folder "FolderToShare" should have the following shares:
      | sharee | shareType | permissionsRole        |
      | Brian  | user      | <new-permissions-role> |
    And for user "Brian" folder "FolderToShare" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}       |
      | {user} shared {resource} with {sharee}    |
      | {user} updated {field} for the {resource} |
    Examples:
      | new-permissions-role   | permissions-role |
      | Viewer With ListGrants | Viewer           |
      | Viewer With ListGrants | Editor           |
      | Viewer With ListGrants | Uploader         |
      | Editor With ListGrants | Viewer           |
      | Editor With ListGrants | Editor           |
      | Editor With ListGrants | Uploader         |


  Scenario Outline: sharer lists shared-by-me (Personal space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Alice" has created folder "folder"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder             |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "folder" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference","permissions","name"],
        "properties": {
          "parentReference": {
            "type": "object",
            "required": ["driveId","driveType","path","name","id"],
            "properties": {
              "driveId": {"pattern": "^%space_id_pattern%$"},
              "driveType": {"const": "personal"},
              "path": {"const": "/"},
              "name": {"const": "/"},
              "id": {"pattern": "^%file_id_pattern%$"}
            }
          },
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["grantedToV2","id","roles"],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName","id"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
                      }
                    }
                  }
                },
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                }
              }
            }
          },
          "name": {"const": "folder"}
        }
      }
      """
    Examples:
      | permissions-role       |
      | Viewer With ListGrants |
      | Editor With ListGrants |


  Scenario Outline: sharer list share shared-by-me when same file is shared with multiple user(Personal space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Carol" has been created with default attributes
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt       |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt       |
      | space           | Personal           |
      | sharee          | Carol              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "textfile.txt" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference","permissions","name","size"],
        "properties": {
          "name": {"const": "textfile.txt"},
          "size": {"const": 8},
          "parentReference": {
            "type": "object",
            "required": ["driveId","driveType","path","name","id"],
            "properties": {
              "driveId": {"pattern": "^%space_id_pattern%$"},
              "driveType": {"const": "personal"},
              "path": {"const": "/"},
              "name": {"const": "/"},
              "id": {"pattern": "^%file_id_pattern%$"}
            }
          },
          "permissions": {
            "type": "array",
            "minItems": 2,
            "maxItems": 2,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": ["createdDateTime","grantedToV2","id","roles","invitation"],
                  "properties": {
                    "id": {"pattern": "^%permissions_id_pattern%$"},
                    "roles": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {"pattern": "^%role_id_pattern%$"}
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
                          "required": ["displayName","id"],
                           "properties": {
                              "id": {"pattern": "^%user_id_pattern%$"},
                              "displayName": {"const": "Brian Murphy"}
                           }
                        }
                      }
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["createdDateTime","grantedToV2","id","roles","invitation"],
                  "properties": {
                    "id": {"pattern": "^%permissions_id_pattern%$"},
                    "roles": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {"pattern": "^%role_id_pattern%$"}
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
                          "required": ["displayName","id"],
                          "properties": {
                            "id": {"pattern": "^%user_id_pattern%$"},
                            "displayName": {"const": "Carol King"}
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
    Examples:
      | permissions-role            |
      | Viewer With ListGrants      |
      | File Editor With ListGrants |


  Scenario Outline: sharee list shared-with-me (Personal space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt       |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
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
              "required": ["@UI.Hidden","@client.synchronize","createdBy","eTag","file",
                "id","lastModifiedDateTime","name","parentReference","remoteItem","size"],
              "properties": {
                "@UI.Hidden": {"const": false},
                "@client.synchronize": {"const": true},
                "createdBy": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id"],
                      "properties": {
                        "displayName": {"const": "Alice Hansen"},
                        "id": {"pattern": "^%user_id_pattern%$"}
                      }
                    }
                  }
                },
                "eTag": {"pattern": "%etag_pattern%"},
                "file": {
                  "type": "object",
                  "required": ["mimeType"],
                  "properties": {
                    "mimeType": {"const": "text/plain"}
                  }
                },
                "id": {"pattern": "^%share_id_pattern%$"},
                "name": {"const": "textfile.txt"},
                "parentReference": {
                  "type": "object",
                  "required": ["driveId","driveType","id"],
                  "properties": {
                    "driveId": {"pattern": "^%space_id_pattern%$"},
                    "driveType": {"const": "virtual"},
                    "id": {"pattern": "^%file_id_pattern%$"}
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": ["createdBy","eTag","file","id","lastModifiedDateTime",
                    "name","parentReference","permissions","size"
                  ],
                  "properties": {
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["id","displayName"],
                          "properties": {
                            "id": {"pattern": "^%user_id_pattern%$"},
                            "displayName": {"const": "Alice Hansen"}
                          }
                        }
                      }
                    },
                    "eTag": {"pattern": "%etag_pattern%"},
                    "file": {
                      "type": "object",
                      "required": ["mimeType"],
                      "properties": {
                        "mimeType": {"const": "text/plain"}
                      }
                    },
                    "id": {"pattern": "^%file_id_pattern%$"},
                    "name": {"const": "textfile.txt"},
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId","driveType"],
                      "properties": {
                        "driveId": {"pattern": "^%file_id_pattern%$"},
                        "driveType": {"const": "personal"}
                      }
                    },
                    "permissions": {
                      "type": "array",
                      "minItems": 1,
                      "maxItems": 1,
                      "items": {
                        "type": "object",
                        "required": ["grantedToV2","id","invitation","roles"],
                        "properties": {
                          "id": {"pattern": "^%permissions_id_pattern%$"},
                          "grantedToV2": {
                            "type": "object",
                            "required": ["user"],
                            "properties": {
                              "user": {
                                "type": "object",
                                "required": ["displayName","id"],
                                "properties": {
                                  "displayName": {"const": "Brian Murphy"},
                                  "id": {"pattern": "^%user_id_pattern%$"}
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
                                      "displayName": {"const": "Alice Hansen"},
                                      "id": {"pattern": "^%user_id_pattern%$"}
                                    },
                                    "required": ["displayName","id"]
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
                            "items": {"pattern": "^%role_id_pattern%$"}
                          }
                        }
                      }
                    }
                  }
                },
                "size": {"const": 8}
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role            |
      | Viewer With ListGrants      |
      | File Editor With ListGrants |


  Scenario Outline: sharee lists same name shares received via user and group invitations (Personal space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And  user "Alice" has created folder "folder"
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder             |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
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
              "required": ["@UI.Hidden","@client.synchronize","createdBy","eTag","folder",
                "id","lastModifiedDateTime","name","parentReference","remoteItem"
              ],
              "properties": {
                "@UI.Hidden":{"const": false},
                "@client.synchronize":{"const": true},
                "createdBy": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName", "id"],
                      "properties": {
                        "displayName": {"const": "Alice Hansen"},
                        "id": {"pattern": "^%user_id_pattern%$"}
                      }
                    }
                  }
                },
                "eTag": {"pattern": "%etag_pattern%"},
                "folder": {"const": {}},
                "id": {"pattern": "^%share_id_pattern%$"},
                "name": {"const": "folder"},
                "parentReference": {
                  "type": "object",
                  "required": ["driveId","driveType","id"],
                  "properties": {
                    "driveId": {"pattern": "^%space_id_pattern%$"},
                    "driveType" : {"const": "virtual"},
                    "id" : {"pattern": "%space_id_pattern%"}
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": ["createdBy","eTag","folder","id","lastModifiedDateTime",
                    "name","parentReference","permissions"
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
                            "id": {"pattern": "^%user_id_pattern%$"},
                            "displayName": {"const": "Alice Hansen"}
                          }
                        }
                      }
                    },
                    "eTag": {"pattern": "%etag_pattern%"},
                    "file": {},
                    "id": {"pattern": "^%file_id_pattern%$"},
                    "name": {"const": "folder"},
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId","driveType"],
                      "properties": {
                        "driveId": {"pattern": "%space_id_pattern%"},
                        "driveType" : {"const": "personal"}
                      }
                    },
                    "permissions": {
                      "type": "array",
                      "minItems": 2,
                      "maxItems": 2,
                      "uniqueItems": true,
                      "items": {
                        "oneOf": [
                          {
                            "type": "object",
                            "required": ["grantedToV2","id","invitation","roles"],
                            "properties": {
                              "grantedToV2": {
                                "type": "object",
                                "required": ["group"],
                                "properties":{
                                  "group": {
                                    "type": "object",
                                    "required": ["displayName","id"],
                                    "properties": {
                                      "displayName": {"const": "grp1"},
                                      "id": {"pattern": "^%user_id_pattern%$"}
                                    }
                                  }
                                }
                              },
                              "id": {"pattern": "^%permissions_id_pattern%$"},
                              "invitation": {
                                "type": "object",
                                "required": ["invitedBy"],
                                "properties": {
                                  "user":{
                                    "type": "object",
                                    "required": ["displayName","id"],
                                    "properties": {
                                      "displayName": {"const": "Alice Hansen"},
                                      "id": {"pattern": "^%user_id_pattern%$"}
                                    }
                                  }
                                }
                              },
                              "roles": {
                                "type": "array",
                                "minItems": 1,
                                "maxItems": 1,
                                "items": {"pattern": "^%role_id_pattern%$"}
                              }
                            }
                          },
                          {
                            "type": "object",
                            "required": ["grantedToV2","id","invitation","roles"],
                            "properties": {
                              "grantedToV2": {
                                "type": "object",
                                "required": ["user"],
                                "properties":{
                                  "user": {
                                    "type": "object",
                                    "required": ["displayName","id"],
                                    "properties": {
                                      "displayName": {"const": "Brian Murphy"},
                                      "id": {"pattern": "^%user_id_pattern%$"}
                                    }
                                  }
                                }
                              },
                              "id": {"pattern": "^%permissions_id_pattern%$"},
                              "invitation": {
                                "type": "object",
                                "required": ["invitedBy"],
                                "properties": {
                                  "user":{
                                    "type": "object",
                                    "required": ["displayName","id"],
                                    "properties": {
                                      "displayName": {"const": "Alice Hansen"},
                                      "id": {"pattern": "^%user_id_pattern%$"}
                                    }
                                  }
                                }
                              },
                              "roles": {
                                "type": "array",
                                "minItems": 1,
                                "maxItems": 1,
                                "items": {"pattern": "^%role_id_pattern%$"}
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
          }
        }
      }
      """
    And for user "Brian" folder "folder" should have the following shares:
      | sharee | shareType | permissionsRole    |
      | Brian  | user      | Viewer             |
      | grp1   | group     | <permissions-role> |
    And for user "Brian" folder "folder" of the space "Shares" should have the following activities:
      | {user} added {resource} to {folder}    |
      | {user} shared {resource} with {sharee} |
    Examples:
      | permissions-role       |
      | Viewer With ListGrants |
      | Editor With ListGrants |


  Scenario Outline: sharer lists the file shared-by-me (Project space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "hello world" to "FolderToShare/textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FolderToShare/textfile.txt |
      | space           | NewSpace                   |
      | sharee          | Brian                      |
      | shareType       | user                       |
      | permissionsRole | <permissions-role>         |
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "textfile.txt" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference","permissions","name","size"],
        "properties": {
          "parentReference": {
            "type": "object",
            "required": ["driveId","driveType","path","name","id"],
            "properties": {
              "driveId": {"pattern": "^%space_id_pattern%$"},
              "driveType": {"const": "project"},
              "path": {"const": "/FolderToShare"},
              "name": {"const": "FolderToShare"},
              "id": {"pattern": "^%file_id_pattern%$"}
            }
          },
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["grantedToV2","id","roles"],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName","id"],
                      "properties": {
                        "id": {"pattern": "^%user_id_pattern%$"},
                        "displayName": {"const": "Brian Murphy"}
                      }
                    }
                  }
                },
                "id": {"pattern": "^%permissions_id_pattern%$"},
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                }
              }
            }
          },
          "name": {"const": "textfile.txt"},
          "size": {"const": 11}
        }
      }
      """
    Examples:
      | permissions-role            |
      | Viewer With ListGrants      |
      | File Editor With ListGrants |


  Scenario: share exists even though sharee has been disabled (Project space)
    Given the following configs have been set:
      | config                       | value                                |
      | GRAPH_SPACES_USERS_CACHE_TTL | 1                                    |
      | GRAPH_AVAILABLE_ROLES        | d5041006-ebb3-4b4a-b6a4-7c180ecfb17d |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt           |
      | space           | NewSpace               |
      | sharee          | Brian                  |
      | shareType       | user                   |
      | permissionsRole | Viewer With ListGrants |
    And the user "Admin" has disabled user "Brian"
    When user "Alice" lists the shares shared by her after clearing user cache using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain resource "textfile.txt" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference","permissions","name"],
        "properties": {
          "parentReference": {
            "type": "object",
            "required": ["driveId","driveType","path","name","id"],
            "properties": {
              "driveType": {"const": "project"}
            }
          },
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["grantedToV2","id","roles"],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName","id"],
                      "properties": {
                        "displayName": {"const": "Brian Murphy"}
                      }
                    }
                  }
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                }
              }
            }
          },
          "name": {"const": "textfile.txt"}
        }
      }
      """


  Scenario Outline: share doesn't exist for disabled space (Project space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt       |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Admin" has disabled a space "NewSpace"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should not contain resource "textfile.txt" with the following data:
      """
      {
        "type": "object",
        "required": ["parentReference","permissions","name"],
        "properties": {
          "parentReference": {
            "type": "object",
            "required": ["driveId","driveType","path","name","id"],
            "properties": {
              "driveType": {"const": "project"}
            }
          },
          "permissions": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["grantedToV2","id","roles"],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": ["user"],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": ["displayName","id"],
                      "properties": {
                        "displayName": {"const": "Brian Murphy"}
                      }
                    }
                  }
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {"pattern": "^%role_id_pattern%$"}
                }
              }
            }
          },
          "name": {"const": "textfile.txt"}
        }
      }
      """
    And user "Brian" should not have a share "textfile.txt" shared by user "Alice" from space "NewSpace"
    Examples:
      | permissions-role            |
      | Viewer With ListGrants      |
      | File Editor With ListGrants |


  Scenario Outline: list activities of folder shared with listGrant roles (Personal space)
    Given the administrator has enabled the permissions role "<permissions-role>"
    And using spaces DAV path
    And using SharingNG
    And user "Alice" has created folder "folder"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder             |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Brian" has a share "folder" synced
    When user "Brian" lists the activities of folder "folder" from space "Shares" using the Graph API
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
              "oneOf": [
                {
                  "type": "object",
                  "required": ["id","template","times"],
                  "properties": {
                    "id": {"pattern": "^%user_id_pattern%$"},
                    "template": {
                      "type": "object",
                      "required": ["message","variables"],
                      "properties": {
                        "message": {"const": "{user} added {resource} to {folder}"},
                        "variables":{
                          "type": "object",
                          "required": ["folder","resource","user"],
                          "properties": {
                            "folder": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties":{
                                "id": {"const": ""},
                                "name": {"const": "shared-with-me"}
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties": {"name": {"const": "folder"}}
                            },
                            "user": {
                              "type": "object",
                              "required": ["id","displayName"],
                              "properties":{"displayName": {"const": "Alice Hansen"}}
                            }
                          }
                        }
                      }
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id","template","times"],
                  "properties": {
                    "id": {"pattern": "^%user_id_pattern%$"},
                    "template": {
                      "type": "object",
                      "required": ["message","variables"],
                      "properties": {
                        "message": {"const": "{user} shared {resource} with {sharee}"},
                        "variables": {
                          "type": "object",
                          "required": ["folder","resource","sharee","user"],
                          "properties": {
                            "resource": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties": {"name": {"const": "folder"}}
                            },
                            "sharee": {
                              "type": "object",
                              "required": ["id","displayName"],
                              "properties": {"displayName": {"const": "Brian"}}
                            },
                            "user": {
                              "type": "object",
                              "required": ["id","displayName"],
                              "properties": {"displayName": {"const": "Alice Hansen"}}
                            }
                          }
                        }
                      }
                    },
                    "times": {
                      "type": "object",
                      "required": ["recordedTime"],
                      "properties": {
                        "recordedTime": {
                          "type": "string",
                          "format": "date-time"
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
    Examples:
      | permissions-role       |
      | Viewer With ListGrants |
      | Editor With ListGrants |


  Scenario: user lists permissions of a folder after enabling 'Viewer With ListGrants' role (Personal space)
    Given the administrator has enabled the permissions role "Viewer With ListGrants"
    And user "Alice" has created folder "folder"
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
          "@libre.graph.permissions.roles.allowedValues": {
            "type": "array",
            "minItems": 4,
            "maxItems": 4,
            "uniqueItems": true,
            "items": {
              "oneOf": [
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
                    "description": {"const": "View, download and show all invited people."},
                    "displayName": {"const": "Can view"},
                    "id": {"const": "d5041006-ebb3-4b4a-b6a4-7c180ecfb17d"}
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "displayName": {"const": "Can upload"}
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "displayName": {"const": "Can edit without versions"}
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists permissions of a file after enabling 'File Editor With ListGrants' role (Project space)
    Given the administrator has enabled the permissions role "File Editor With ListGrants"
    And using spaces DAV path
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
          "@libre.graph.permissions.roles.allowedValues": {
            "type": "array",
            "minItems": 3,
            "maxItems": 3,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 1},
                    "displayName": {"const": "Can view"}
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 2},
                    "displayName": {"const": "Can edit without versions"},
                    "id": {"const": "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"}
                  }
                },
                {
                  "type": "object",
                  "required": ["@libre.graph.weight","description","displayName","id"],
                  "properties": {
                    "@libre.graph.weight": {"const": 3},
                    "description": {"const": "View, download, edit and show all invited people."},
                    "displayName": {"const": "Can edit without versions"},
                    "id": {"const": "c1235aea-d106-42db-8458-7d5610fb0a67"}
                  }
                }
              ]
            }
          }
        }
      }
      """
