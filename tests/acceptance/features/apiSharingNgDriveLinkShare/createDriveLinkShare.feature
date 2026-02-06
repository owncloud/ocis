@issue-7879
Feature: Link sharing of project spaces
  https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |


  Scenario Outline: create a link share of a project-space drive using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using permissions endpoint of the Graph API:
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
    When user "Alice" creates the following space link share using permissions endpoint of the Graph API:
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
    When user "Alice" creates the following space link share using permissions endpoint of the Graph API:
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


  Scenario Outline: create a link share of a project-space drive using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using root endpoint of the Graph API:
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
    When user "Alice" creates the following space link share using root endpoint of the Graph API:
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
    When user "Alice" creates the following space link share using root endpoint of the Graph API:
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


  Scenario Outline: try to create a link share of a Personal and Shares drives using permissions endpoint
    When user "Alice" tries to create the following space link share using permissions endpoint of the Graph API:
      | space           | <drive>            |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
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
                "required": ["date","request-id"]
              },
              "message": {"const": "<message>"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | drive    | message                                   |
      | View             | Shares   | cannot create link on shares space root   |
      | Edit             | Shares   | cannot create link on shares space root   |
      | Upload           | Shares   | cannot create link on shares space root   |
      | File Drop        | Shares   | cannot create link on shares space root   |
      | Secure View      | Shares   | cannot create link on shares space root   |
      | View             | Personal | cannot create link on personal space root |
      | Edit             | Personal | cannot create link on personal space root |
      | Upload           | Personal | cannot create link on personal space root |
      | File Drop        | Personal | cannot create link on personal space root |
      | Secure View      | Personal | invalid link type                         |


  Scenario Outline: try to create a link share of a Personal and Shares drives using root endpoint
    When user "Alice" tries to create the following space link share using root endpoint of the Graph API:
      | space           | <drive>            |
      | permissionsRole | <permissions-role> |
      | password        | %public%           |
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
                "required": ["date","request-id"]
              },
              "message": {"const": "unsupported space type"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | drive    |
      | View             | Shares   |
      | Edit             | Shares   |
      | Upload           | Shares   |
      | File Drop        | Shares   |
      | Secure View      | Shares   |
      | View             | Personal |
      | Edit             | Personal |
      | Upload           | Personal |
      | File Drop        | Personal |
      | Secure View      | Personal |

  @env-config
  Scenario: try to create a link share of a project-space with a password that is listed in the Banned-Password-List using permissions endpoint
    Given using spaces DAV path
    And the config "SHARING_PASSWORD_POLICY_BANNED_PASSWORDS_LIST" has been set to path "config/drone/banned-password-list.txt" for "sharing" service
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using permissions endpoint of the Graph API:
      | space           | projectSpace |
      | permissionsRole | View         |
      | password        | 123          |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "message": {"const": "unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety"}
            }
          }
        }
      }
      """

  @env-config
  Scenario: try to create a link share of a project-space with a password that is listed in the Banned-Password-List using root endpoint
    Given using spaces DAV path
    And the config "SHARING_PASSWORD_POLICY_BANNED_PASSWORDS_LIST" has been set to path "config/drone/banned-password-list.txt" for "sharing" service
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using root endpoint of the Graph API:
      | space           | projectSpace |
      | permissionsRole | Edit         |
      | password        | password     |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "message": {"const": "unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety"}
            }
          }
        }
      }
      """

  @env-config
  Scenario Outline: create a link share of a project-space when password is not enforced using permissions endpoint
    Given using spaces DAV path
    And the following configs have been set:
      | service | config                                  | value |
      | sharing | SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD | false |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using permissions endpoint of the Graph API:
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
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

  @env-config
  Scenario Outline: create a link share of a project-space when password is not enforced using root endpoint
    Given using spaces DAV path
    And the following configs have been set:
      | service | config                                  | value |
      | sharing | SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD | false |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using root endpoint of the Graph API:
      | space           | projectSpace       |
      | permissionsRole | <permissions-role> |
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
