@issue-8619
Feature: Internal link share of resources
  https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |


  Scenario: create an internal link share of resources
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test content" to "textfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | Internal |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt |
      | space           | Personal     |
      | permissionsRole | Internal     |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """


  Scenario: try to create an internal link share of resources with password
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "other data" to "textfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | Internal |
      | password        | %public% |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","innererror","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "innererror": {
                "type": "object",
                "required": ["date","request-id"]
              },
              "message": {"const": "password is redundant for the internal link"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt |
      | space           | Personal     |
      | permissionsRole | Internal     |
      | password        | %public%     |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","innererror","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "innererror": {
                "type": "object",
                "required": ["date","request-id"]
              },
              "message": {"const": "password is redundant for the internal link"}
            }
          }
        }
      }
      """


  Scenario: create an internal link share of a folder inside project-space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folderToShare |
      | space           | projectSpace  |
      | permissionsRole | Internal      |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt |
      | space           | projectSpace |
      | permissionsRole | Internal     |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """


  Scenario: try to create an internal link share of a folder inside project-space with password using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folderToShare |
      | space           | projectSpace  |
      | permissionsRole | Internal      |
      | password        | %public%      |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","innererror","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "innererror": {
                "type": "object",
                "required": ["date","request-id"]
              },
              "message": {"const": "password is redundant for the internal link"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt |
      | space           | projectSpace |
      | permissionsRole | Internal     |
      | password        | %public%     |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","innererror","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "innererror": {
                "type": "object",
                "required": ["date","request-id"]
              },
              "message": {"const": "password is redundant for the internal link"}
            }
          }
        }
      }
      """


  Scenario: create an internal quick link share of a folder
    Given user "Alice" has created folder "folder"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | Internal |
      | quickLink       | true     |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """


  Scenario: create internal quick link share of a resource multiple times
    Given user "Alice" has uploaded file with content "test content" to "textfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt |
      | space           | Personal     |
      | permissionsRole | Internal     |
      | quickLink       | true         |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt |
      | space           | Personal     |
      | permissionsRole | Internal     |
      | quickLink       | true         |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    And the last 2 link shares should have the same id


  Scenario: create an internal quick link share of a folder of a project-space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folderToShare |
      | space           | projectSpace  |
      | permissionsRole | Internal      |
      | quickLink       | true          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """


  Scenario: create internal quick link share of a resource of a project-space multiple times
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt |
      | space           | projectSpace |
      | permissionsRole | Internal     |
      | quickLink       | true         |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt |
      | space           | projectSpace |
      | permissionsRole | Internal     |
      | quickLink       | true         |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "internal"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    And the last 2 link shares should have the same id
