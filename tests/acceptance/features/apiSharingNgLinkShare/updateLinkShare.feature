Feature: Update a link share for a resource
  https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |

  @issue-7879
  Scenario Outline: update role of a file's link share using permissions endpoint
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created the following resource link share:
    | resource        | textfile1.txt      |
    | space           | Personal           |
    | permissionsRole | <permissions-role> |
    | password        | %public%           |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
    | resource        | textfile1.txt          |
    | space           | Personal               |
    | permissionsRole | <new-permissions-role> |
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
            "const": true
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
                "const": ""
              },
              "@libre.graph.quickLink": {
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "<new-permissions-role>"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
    | permissions-role | new-permissions-role |
    | view             | edit                 |
    | view             | blocksDownload       |
    | edit             | view                 |
    | edit             | blocksDownload       |
    | blocksDownload   | edit                 |
    | blocksDownload   | blocksDownload       |

  @issue-8619
  Scenario Outline: update role of a file's to internal link share using permissions endpoint
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created the following resource link share:
    | resource        | textfile1.txt      |
    | space           | Personal           |
    | permissionsRole | <permissions-role> |
    | password        | %public%           |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
    | resource        | textfile1.txt |
    | space           | Personal      |
    | permissionsRole | internal      |
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
            "const": false
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
                "const": ""
              },
              "@libre.graph.quickLink": {
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "internal"
              },
              "webUrl": {
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
    | permissions-role |
    | view             |
    | edit             |
    | blocksDownload   |


  Scenario: update expiration date of a file's link share using permissions endpoint
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
      And user "Alice" has created the following resource link share:
      | resource           | textfile1.txt            |
      | space              | Personal                 |
      | permissionsRole    | view                     |
      | password           | %public%                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
      When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
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
            "const": true
          },
          "id": {
            "pattern": "^[a-zA-Z]{15}$"
          },
          "expirationDateTime": {
            "const": "2201-07-15T23:59:59Z"
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
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "view"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """


  Scenario: update password of a file's link share using permissions endpoint
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created the following resource link share:
      | resource        | textfile1.txt |
      | space           | Personal      |
      | permissionsRole | view          |
      | password        | $heLlo*1234*  |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource | textfile1.txt |
      | space    | Personal      |
      | password | %public%      |
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
            "const": true
          }
        }
      }
      """
    And the public should be able to download file "textfile1.txt" from the last link share with password "%public%" and the content should be "other data"
    And the public download of file "textfile1.txt" from the last link share with password "$heLlo*1234*" should fail with HTTP status code "401" using shareNg


  Scenario Outline: update a file's link share with a password that is listed in the Banned-Password-List using permissions endpoint
    Given the config "OCIS_PASSWORD_POLICY_BANNED_PASSWORDS_LIST" has been set to path "config/drone/banned-password-list.txt"
    And user "Alice" has uploaded file with content "other data" to "text.txt"
    And user "Alice" has created the following resource link share:
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
                "const": "unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety"
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

  @env-config
  Scenario: set password on a existing link share of a folder inside project-space using permissions endpoint
    Given the following configs have been set:
      | config                                       | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD | false |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "folderToShare/textfile.txt"
    And user "Alice" has created the following resource link share:
      | resource        | folderToShare |
      | space           | projectSpace  |
      | permissionsRole | view          |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource | folderToShare |
      | space    | projectSpace  |
      | password | %public%      |
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
            "const": true
          }
        }
      }
      """
    And the public should be able to download file "/textfile.txt" from the last link share with password "%public%" and the content should be "to share"

  @env-config
  Scenario: set password on a existing link share of a project-space drive using root endpoint
    Given the following configs have been set:
      | config                                       | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD | false |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Alice" has created the following space link share:
      | space           | projectSpace |
      | permissionsRole | view         |
    When user "Alice" sets the following password for the last space link share using root endpoint of the Graph API:
      | space    | projectSpace |
      | password | %public%     |
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
            "const": true
          }
        }
      }
      """
    And the public should be able to download file "textfile.txt" from the last link share with password "%public%" and the content should be "to share"


  Scenario: update password on a existing link share of a project-space drive using root endpoint
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Alice" has created the following space link share:
      | space           | projectSpace |
      | permissionsRole | view         |
      | password        | $heLlo*1234* |
    When user "Alice" sets the following password for the last space link share using root endpoint of the Graph API:
      | space    | projectSpace |
      | password | %public%     |
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
            "const": true
          }
        }
      }
      """
    And the public should be able to download file "textfile.txt" from the last link share with password "%public%" and the content should be "to share"
    And the public download of file "textfile.txt" from the last link share with password "$heLlo*1234*" should fail with HTTP status code "401" using shareNg

  @issue-7879
  Scenario Outline: update link share of a project-space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created the following space link share:
      | space              | projectSpace             |
      | permissionsRole    | <permissions-role>       |
      | password           | %public%                 |
      | displayName        | Homework                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | space              | projectSpace           |
      | permissionsRole    | <new-permissions-role> |
      | password           | p@$$w0rD               |
      | expirationDateTime | 2201-07-15T14:00:00Z   |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link",
          "expirationDateTime",
          "createdDateTime"
        ],
        "properties": {
          "hasPassword": {
            "const": true
          },
          "id": {
            "type": "string",
            "pattern": "^[a-zA-Z]{15}$"
          },
          "expirationDateTime": {
            "const": "2201-07-15T23:59:59Z"
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
                "const": "Homework"
              },
              "@libre.graph.quickLink": {
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "<new-permissions-role>"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | new-permissions-role |
      | view             | edit                 |
      | view             | upload               |
      | view             | createOnly           |
      | edit             | view                 |
      | edit             | upload               |
      | edit             | createOnly           |
      | upload           | view                 |
      | upload           | edit                 |
      | upload           | createOnly           |
      | createOnly       | view                 |
      | createOnly       | edit                 |
      | createOnly       | upload               |
      | blocksDownload   | blocksDownload       |


  Scenario Outline: update role of a folder's link share inside project-space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    And user "Alice" has created the following resource link share:
      | resource        | folderToShare      |
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
      | displayName     | Link               |
      | password        | %public%           |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource           | folderToShare          |
      | space              | projectSpace           |
      | permissionsRole    | <new-permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link",
          "createdDateTime"
        ],
        "properties": {
          "hasPassword": {
            "const": true
          },
          "id": {
            "type": "string",
            "pattern": "^[a-zA-Z]{15}$"
          },
          "expirationDateTime": {
            "const": "2201-07-15T23:59:59Z"
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
                "const": "Link"
              },
              "@libre.graph.quickLink": {
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "<new-permissions-role>"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | new-permissions-role |
      | view             | edit                 |
      | view             | upload               |
      | view             | createOnly           |
      | edit             | view                 |
      | edit             | upload               |
      | edit             | createOnly           |
      | upload           | view                 |
      | upload           | edit                 |
      | upload           | createOnly           |
      | createOnly       | view                 |
      | createOnly       | edit                 |
      | createOnly       | upload               |


  Scenario Outline: update role of a file's link share inside a project-space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Alice" has created the following resource link share:
      | resource        | textfile.txt       |
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
      | displayName     | Link               |
      | password        | %public%           |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource           | textfile.txt           |
      | space              | projectSpace           |
      | permissionsRole    | <new-permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link",
          "createdDateTime"
        ],
        "properties": {
          "hasPassword": {
            "const": true
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
                "const": "Link"
              },
              "@libre.graph.quickLink": {
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "<new-permissions-role>"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | new-permissions-role |
      | view             | edit                 |
      | edit             | view                 |


  Scenario Outline: update role of a file's link share to internal inside a project-space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Alice" has created the following resource link share:
      | resource        | textfile.txt       |
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
      | displayName     | Link               |
      | password        | %public%           |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource           | textfile.txt           |
      | space              | projectSpace           |
      | permissionsRole    | <new-permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link",
          "createdDateTime"
        ],
        "properties": {
          "hasPassword": {
            "const": false
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
                "const": "Link"
              },
              "@libre.graph.quickLink": {
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "<new-permissions-role>"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | new-permissions-role |
      | view             | internal             |
      | edit             | internal             |


  Scenario Outline: update role of a folder's link share to internal inside project-space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "projectSpace"
    And user "Alice" has created the following resource link share:
      | resource        | folderToShare      |
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
      | displayName     | Link               |
      | password        | %public%           |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource           | folderToShare          |
      | space              | projectSpace           |
      | permissionsRole    | <new-permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link",
          "createdDateTime"
        ],
        "properties": {
          "hasPassword": {
            "const": false
          },
          "id": {
            "type": "string",
            "pattern": "^[a-zA-Z]{15}$"
          },
          "expirationDateTime": {
            "const": "2201-07-15T23:59:59Z"
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
                "const": "Link"
              },
              "@libre.graph.quickLink": {
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "<new-permissions-role>"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | new-permissions-role |
      | view             | internal             |
      | edit             | internal             |
      | upload           | internal             |
      | createOnly       | internal             |


  Scenario Outline: update link share of a folder inside personal drive using permissions endpoint
    Given user "Alice" has created folder "folder"
    And user "Alice" has created the following resource link share:
      | space           | Personal           |
      | resource        | folder             |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
      | displayName     | Homework           |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | space              | Personal               |
      | resource           | folder                 |
      | permissionsRole    | <new-permissions-role> |
      | password           | p@$$w0rD               |
      | expirationDateTime | 2201-07-15T14:00:00Z   |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "hasPassword",
          "id",
          "link",
          "expirationDateTime",
          "createdDateTime"
        ],
        "properties": {
          "hasPassword": {
            "const": true
          },
          "id": {
            "type": "string",
            "pattern": "^[a-zA-Z]{15}$"
          },
          "expirationDateTime": {
            "const": "2201-07-15T23:59:59Z"
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
                "const": "Homework"
              },
              "@libre.graph.quickLink": {
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "<new-permissions-role>"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | new-permissions-role |
      | view             | edit                 |
      | view             | upload               |
      | view             | createOnly           |
      | edit             | view                 |
      | edit             | upload               |
      | edit             | createOnly           |
      | upload           | view                 |
      | upload           | edit                 |
      | upload           | createOnly           |
      | createOnly       | view                 |
      | createOnly       | edit                 |
      | createOnly       | upload               |
      | blocksDownload   | blocksDownload       |

  @issues-8405
  Scenario Outline: remove expiration date of a resource link share using permissions endpoint
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created folder "folder"
    And user "Alice" has created the following resource link share:
      | resource           | <resource>               |
      | space              | Personal                 |
      | permissionsRole    | view                     |
      | password           | %public%                 |
      | expirationDateTime | 2200-07-15T14:00:00.000Z |
    When user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource           | <resource> |
      | space              | Personal   |
      | expirationDateTime |            |
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
            "const": true
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
                "const": false
              },
              "preventsDownload": {
                "const": false
              },
              "type": {
                "const": "view"
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/s/[a-zA-Z]{15}$"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource      |
      | textfile1.txt |
      | folder        |
