Feature: Create a share link for a resource
  https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |

  Scenario Outline: create a sharing link for a folder
    Given user "Alice" has created folder "folder"
    When user "Alice" creates a share link for a folder "folder" of the space "Personal" using the Graph API with settings:
      | role     | <role>   |
      | password | %public% |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link"
        ],
        "properties": {
          "hasPassword": {
            "type": "boolean",
            "enum": [true]
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
                "type": "string",
                "enum": [""]
              },
              "@libre.graph.quickLink": {
                "type": "boolean",
                "enum": [false]
              },
              "preventsDownload": {
                "type": "boolean",
                "enum": [false]
              },
              "type": {
                "type": "string",
                "enum": ["<role>"]
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%\/s\/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | role           |
      | view           |
      | edit           |
      | upload         |
      | createOnly     |


  Scenario Outline: create a sharing link for a file
    Given user "Alice" has uploaded file with content "other data" to "/textfile1.txt"
    When user "Alice" creates a share link for a file "textfile1.txt" of the space "Personal" using the Graph API with settings:
      | role     | <role>   |
      | password | %public% |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link"
        ],
        "properties": {
          "hasPassword": {
            "type": "boolean",
            "enum": [true]
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
                "type": "string",
                "enum": [""]
              },
              "@libre.graph.quickLink": {
                "type": "boolean",
                "enum": [false]
              },
              "preventsDownload": {
                "type": "boolean",
                "enum": [false]
              },
              "type": {
                "type": "string",
                "enum": ["<role>"]
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%\/s\/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | role           |
      | view           |
      | edit           |
