Feature: Link share for resources
  https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |

  @issue-7879
  Scenario Outline: create a link share of a folder
    Given user "Alice" has created folder "folder"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
      | displayName     | folderLink         |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["createdDateTime","hasPassword","id","link"],
        "properties": {
          "createdDateTime": {"format": "date-time"},
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": "folderLink"},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource           | folder                   |
      | space              | Personal                 |
      | permissionsRole    | <permissions-role>       |
      | password           | %public%                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["createdDateTime","hasPassword","id","link","expirationDateTime"],
        "properties": {
          "createdDateTime": {"format": "date-time"},
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "expirationDateTime": {"const": "2200-07-15T14:00:00Z"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
      | quickLink       | true               |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["createdDateTime","hasPassword","id","link"],
        "properties": {
          "createdDateTime": {"format": "date-time"},
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | permissions-role-value |
      | View             | view                   |
      | Edit             | edit                   |
      | Upload           | upload                 |
      | File Drop        | createOnly             |
      | Secure View      | blocksDownload         |

  @issue-7879
  Scenario Outline: create a link share of a file
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
      | displayName     | fileLink           |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["createdDateTime","hasPassword","id","link"],
        "properties": {
          "createdDateTime": {"format": "date-time"},
          "hasPasword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": "fileLink"},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource           | textfile1.txt            |
      | space              | Personal                 |
      | permissionsRole    | Edit                     |
      | password           | %public%                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["createdDateTime","hasPassword","id","link","expirationDateTime"],
        "properties": {
          "createdDateTime": {"format": "date-time"},
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "expirationDateTime": {"const": "2200-07-15T14:00:00Z"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "edit"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource           | textfile1.txt |
      | space              | Personal      |
      | permissionsRole    | Edit          |
      | password           | %public%      |
      | quickLink          | true          |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["createdDateTime","hasPassword","id","link"],
        "properties": {
          "createdDateTime": {"format": "date-time"},
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "edit"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | permissions-role-value |
      | View             | view                   |
      | Edit             | edit                   |
      | Secure View      | blocksDownload         |

  @issue-7879
  Scenario Outline: create a link share of a folder inside project-space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folderToShare      |
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource           | folderToShare            |
      | space              | projectSpace             |
      | permissionsRole    | <permissions-role>       |
      | password           | %public%                 |
      | displayName        | Homework                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link","expirationDateTime"],
        "properties": {
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "expirationDateTime": {"const": "2200-07-15T14:00:00Z"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": "Homework"},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folderToShare      |
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
      | quickLink       | true               |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | permissions-role-value |
      | View             | view                   |
      | Edit             | edit                   |
      | Upload           | upload                 |
      | File Drop        | createOnly             |
      | Secure View      | blocksDownload         |

  @issue-7879
  Scenario Outline: create a link share of a file inside project-space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt       |
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource           | textfile.txt             |
      | space              | projectSpace             |
      | permissionsRole    | <permissions-role>       |
      | password           | %public%                 |
      | displayName        | Homework                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link","expirationDateTime"],
        "properties": {
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "expirationDateTime": {"const": "2200-07-15T14:00:00Z"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": "Homework"},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile.txt       |
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
      | quickLink       | true               |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": true},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": true},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | permissions-role-value |
      | View             | view                   |
      | Edit             | edit                   |
      | Secure View      | blocksDownload         |

  @env-config
  Scenario: try to create a public link share of a folder with denied permissions role
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | FolderToShare |
      | space           | Personal      |
      | permissionsRole | Denied        |
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
            "required": ["code", "innererror", "message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "innererror": {"required": [ "date", "request-id" ]},
              "message": {"const": "invalid body schema definition"}
            }
          }
        }
      }
      """
