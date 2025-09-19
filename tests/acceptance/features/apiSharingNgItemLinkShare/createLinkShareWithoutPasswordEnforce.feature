@env-config
Feature: Link share without enforcing password
  https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
    And the following configs have been set:
      | service | config                                       | value |
      | sharing | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD | false |

  @issue-7879
  Scenario Outline: create a link share of a folder
    Given user "Alice" has created folder "folderToShare"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folderToShare      |
      | space           | Personal           |
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

  @issue-7879
  Scenario Outline: create a link share of a file
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | textfile1.txt      |
      | space           | Personal           |
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
      | Secure View      | blocksDownload         |


  Scenario: set password on existing file link share
    Given user "Alice" has uploaded file with content "other data" to "textfile1.txt"
    And user "Alice" has created the following resource link share:
      | resource        | textfile1.txt |
      | space           | Personal      |
      | permissionsRole | View          |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource | textfile1.txt |
      | space    | Personal      |
      | password | %public%      |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword"],
        "properties": {
          "hasPassword": {"const": true}
        }
      }
      """
    And the public should be able to download file "textfile1.txt" from the last link share with password "%public%" and the content should be "other data"

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
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {"hasPassword": {"const": false},
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
      | Secure View      | blocksDownload         |


  Scenario: set password on a existing file link share inside project-space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Alice" has created the following resource link share:
      | resource        | textfile.txt |
      | space           | projectSpace |
      | permissionsRole | View         |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource | textfile.txt |
      | space    | projectSpace |
      | password | %public%     |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword"],
        "properties": {
          "hasPassword": {"const": true}
        }
      }
      """
    And the public should be able to download file "textfile.txt" from the last link share with password "%public%" and the content should be "to share"
