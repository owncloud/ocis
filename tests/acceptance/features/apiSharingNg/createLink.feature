Feature: Create a share link for a resource
  https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |

  @issue-7879
  Scenario Outline: create a link share of a folder
    Given user "Alice" has created folder "folder"
    When user "Alice" creates the following link share using the Graph API:
      | resourceType | folder   |
      | resource     | folder   |
      | space        | Personal |
      | role         | <role>   |
      | password     | %public% |
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
      | blocksDownload |

  @issue-7879
  Scenario Outline: create a link share of a file
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    When user "Alice" creates the following link share using the Graph API:
      | resourceType | file          |
      | resource     | textfile1.txt |
      | space        | Personal      |
      | role         | <role>        |
      | password     | %public%      |
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
      | blocksDownload |

  @issue-7879
  Scenario Outline: create a link share of a folder with display name and expiry date
    Given user "Alice" has created folder "folder"
    When user "Alice" creates the following link share using the Graph API:
      | resourceType       | folder                   |
      | resource           | folder                   |
      | space              | Personal                 |
      | role               | <role>                   |
      | password           | %public%                 |
      | displayName        | Homework                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link",
          "expirationDateTime"
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
          "expirationDateTime": {
            "type": "string",
            "enum": ["2200-07-15T23:59:59Z"]
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
                "enum": ["Homework"]
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
      | blocksDownload |

  @issue-7879
  Scenario Outline: create a link share of a file with display name and expiry date
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    When user "Alice" creates the following link share using the Graph API:
      | resourceType       | file                     |
      | resource           | textfile1.txt            |
      | space              | Personal                 |
      | role               | <role>                   |
      | password           | %public%                 |
      | displayName        | Homework                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link",
          "expirationDateTime"
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
          "expirationDateTime": {
            "type": "string",
            "enum": ["2200-07-15T23:59:59Z"]
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
                "enum": ["Homework"]
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
      | blocksDownload |

  @env-config
  Scenario Outline: create a link share of a file without password
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD           | false |
    And user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    When user "Alice" creates the following link share using the Graph API:
      | resourceType | file          |
      | resource     | textfile1.txt |
      | space        | Personal      |
      | role         | <role>        |
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
            "enum": [false]
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
      | blocksDownload |
