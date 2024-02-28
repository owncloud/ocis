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
      | resource        | folder            |
      | space           | Personal          |
      | permissionsRole | <permissionsRole> |
      | password        | %public%          |
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
                "enum": ["<permissionsRole>"]
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
      | permissionsRole |
      | view            |
      | edit            |
      | internal        |
      | upload          |
      | createOnly      |
      | blocksDownload  |

  @issue-7879
  Scenario Outline: create a link share of a file
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    When user "Alice" creates the following link share using the Graph API:
      | resource        | textfile1.txt     |
      | space           | Personal          |
      | permissionsRole | <permissionsRole> |
      | password        | %public%          |
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
                "enum": ["<permissionsRole>"]
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
      | permissionsRole |
      | view            |
      | edit            |
      | internal        |
      | blocksDownload  |

  @issue-7879
  Scenario Outline: create a link share of a folder with display name and expiry date
    Given user "Alice" has created folder "folder"
    When user "Alice" creates the following link share using the Graph API:
      | resource           | folder                   |
      | space              | Personal                 |
      | permissionsRole    | <permissionsRole>        |
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
                "enum": ["<permissionsRole>"]
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
      | permissionsRole |
      | view            |
      | edit            |
      | internal        |
      | upload          |
      | createOnly      |
      | blocksDownload  |

  @issue-7879
  Scenario Outline: create a link share of a file with display name and expiry date
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    When user "Alice" creates the following link share using the Graph API:
      | resource           | textfile1.txt            |
      | space              | Personal                 |
      | permissionsRole    | <permissionsRole>        |
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
                "enum": ["<permissionsRole>"]
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
      | permissionsRole |
      | view            |
      | edit            |
      | internal        |
      | blocksDownload  |

  @env-config @issue-7879
  Scenario Outline: create a link share of a file without password
    Given the following configs have been set:
      | config                                       | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD | false |
    And user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    When user "Alice" creates the following link share using the Graph API:
      | resource        | textfile1.txt     |
      | space           | Personal          |
      | permissionsRole | <permissionsRole> |
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
                "enum": ["<permissionsRole>"]
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
      | permissionsRole |
      | view            |
      | edit            |
      | internal        |
      | blocksDownload  |

  @issue-7879
  Scenario Outline: update role of a file's link share
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created the following link share:
      | resource        | textfile1.txt             |
      | space           | Personal                  |
      | permissionsRole | <previousPermissionsRole> |
      | password        | %public%                  |
    When user "Alice" updates the last public link share using the Graph API with
      | resource        | textfile1.txt        |
      | space           | Personal             |
      | permissionsRole | <newPermissionsRole> |
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
                "enum": ["<newPermissionsRole>"]
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
      | previousPermissionsRole | newPermissionsRole |
      | view                    | edit               |
      | view                    | internal           |
      | view                    | blocksDownload     |
      | edit                    | view               |
      | edit                    | blocksDownload     |
      | view                    | internal           |
      | blocksDownload          | edit               |
      | blocksDownload          | blocksDownload     |
      | view                    | internal           |


  Scenario: update expiration date of a file's link share
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created the following link share:
      | resource           | textfile1.txt            |
      | space              | Personal                 |
      | permissionsRole    | view                     |
      | password           | %public%                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    When user "Alice" updates the last public link share using the Graph API with
      | resource           | textfile1.txt            |
      | space              | Personal                 |
      | expirationDateTime | 2201-07-15T14:00:00.000Z |
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
            "enum": ["2201-07-15T23:59:59Z"]
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
                "enum": ["view"]
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

  @issues-8405
  Scenario Outline: remove expiration date of a resource link share
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created folder "folder"
    And user "Alice" has created the following link share:
      | resource           | <path>                   |
      | space              | Personal                 |
      | permissionsRole    | view                     |
      | password           | %public%                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    When user "Alice" updates the last public link share using the Graph API with
      | resource           | <path>   |
      | space              | Personal |
      | expirationDateTime |          |
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
                "enum": ["view"]
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
      | path          |
      | textfile1.txt |
      | folder        |

  @env-config
  Scenario: set password on a file's link share
    Given the following configs have been set:
      | config                                       | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD | false |
    And user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created the following link share:
      | resource        | textfile1.txt |
      | space           | Personal      |
      | permissionsRole | view          |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource     | textfile1.txt |
      | space        | Personal      |
      | password     | %public%      |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword"
        ],
        "properties": {
          "hasPassword": {
            "type": "boolean",
            "enum": [true]
          }
        }
      }
      """
    And the public should be able to download file "textfile1.txt" from the last link share with password "%public%" and the content should be "other data"

  Scenario: update password of a file's link share
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created the following link share:
      | resource        | textfile1.txt |
      | space           | Personal      |
      | permissionsRole | view          |
      | password        | $heLlo*1234*  |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource     | textfile1.txt |
      | space        | Personal      |
      | password     | %public%      |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword"
        ],
        "properties": {
          "hasPassword": {
            "type": "boolean",
            "enum": [true]
          }
        }
      }
      """
    And the public should be able to download file "textfile1.txt" from the last link share with password "%public%" and the content should be "other data"
    And the public download of file "textfile1.txt" from the last link share with password "$heLlo*1234*" should fail with HTTP status code "401" using shareNg


  Scenario Outline: create a file's link share with a password that is listed in the Banned-Password-List
    Given the config "OCIS_PASSWORD_POLICY_BANNED_PASSWORDS_LIST" has been set to path "config/drone/banned-password-list.txt"
    And user "Alice" has uploaded file with content "other data" to "text.txt"
    When user "Alice" creates the following link share using the Graph API:
      | resource        | text.txt          |
      | space           | Personal          |
      | permissionsRole | view              |
      | password        | <banned-password> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "pattern": "invalidRequest"
              },
              "message": {
                "type": "string",
                "enum": [
                  "unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | banned-password |
      | 123             |
      | password        |
      | ownCloud        |


  Scenario Outline: update a file's link share with a password that is listed in the Banned-Password-List
    Given the config "OCIS_PASSWORD_POLICY_BANNED_PASSWORDS_LIST" has been set to path "config/drone/banned-password-list.txt"
    And user "Alice" has uploaded file with content "other data" to "text.txt"
    And user "Alice" has created the following link share:
      | resource        | text.txt |
      | space           | Personal |
      | permissionsRole | view     |
      | password        | %public% |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource        | text.txt          |
      | space           | Personal          |
      | permissionsRole | view              |
      | password        | <banned-password> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "pattern": "invalidRequest"
              },
              "message": {
                "type": "string",
                "enum": [
                  "unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | banned-password |
      | 123             |
      | password        |
      | ownCloud        |
