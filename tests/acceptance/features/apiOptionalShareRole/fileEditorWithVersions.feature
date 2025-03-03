Feature: an user shares resources
  As a user
  I want to share files with FileEditorListGrantsWithVersions role

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |

  @env-config
  Scenario: sharee checks version of a file shared with FileEditorWithVersions role
    Given the administrator has enabled the permissions role 'File Editor With Versions'
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
    And user "Brian" should have a share "textfile.txt" shared by user "Alice" from space "Personal"
    And user "Brian" has uploaded file with content "updated content" to "Shares/textfile.txt"
    When user "Brian" gets the number of versions of file "textfile.txt" using file-id "<<FILEID>>"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
