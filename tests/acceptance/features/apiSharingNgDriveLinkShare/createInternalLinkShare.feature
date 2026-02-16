@issue-8619
Feature: Internal link share of project spaces
  https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |


  Scenario: create an internal link share of a project-space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using permissions endpoint of the Graph API:
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


  Scenario: try to create an internal link share of a project-space with password using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using permissions endpoint of the Graph API:
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

  @issue-8960
  Scenario Outline: create an internal link share by a member of a project space using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Brian" has been created with default attributes
    And user "Alice" has sent the following space share invitation:
      | space           | projectSpace       |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" creates the following space link share using root endpoint of the Graph API:
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
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @issue-8960
  Scenario Outline: try to create an internal link share by a member of a project space drive with password using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Brian" has been created with default attributes
    And user "Alice" has sent the following space share invitation:
      | space           | projectSpace       |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" creates the following space link share using root endpoint of the Graph API:
      | space           | projectSpace |
      | permissionsRole | Internal     |
      | quickLink       | true         |
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
            "required": ["code","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "message": {"const": "password is redundant for the internal link"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @issue-8960
  Scenario Outline: create an internal link share by a member of a project space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Brian" has been created with default attributes
    And user "Alice" has sent the following space share invitation:
      | space           | projectSpace       |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" creates the following space link share using permissions endpoint of the Graph API:
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
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @issue-8960
  Scenario Outline: try to create an internal link share by a member of a project space drive with password using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "projectSpace" with content "to share" to "textfile.txt"
    And user "Brian" has been created with default attributes
    And user "Alice" has sent the following space share invitation:
      | space           | projectSpace       |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" creates the following space link share using permissions endpoint of the Graph API:
      | space           | projectSpace |
      | permissionsRole | Internal     |
      | quickLink       | true         |
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
            "required": ["code","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "message": {"const": "password is redundant for the internal link"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: try to create an internal link share of a Personal and Shares drives using root endpoint
    When user "Alice" tries to create the following space link share using root endpoint of the Graph API:
      | space           | <drive>  |
      | permissionsRole | Internal |
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
              "message": {"const": "unsupported space type"}
            }
          }
        }
      }
      """
    Examples:
      | drive    |
      | Personal |
      | Shares   |


  Scenario Outline: try to create an internal link share of a Personal and Shares drives with password using root endpoint
    When user "Alice" tries to create the following space link share using root endpoint of the Graph API:
      | space           | <drive>  |
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
              "message": {"const": "unsupported space type"}
            }
          }
        }
      }
      """
    Examples:
      | drive    |
      | Personal |
      | Shares   |

  @issue-11409
  Scenario Outline: try to create an internal link share of a Personal and Shares drives using permissions endpoint
    When user "Alice" tries to create the following space link share using permissions endpoint of the Graph API:
      | space           | <drive>  |
      | permissionsRole | Internal |
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
              "message": {"const": "<message>"}
            }
          }
        }
      }
      """
    Examples:
      | drive    | message                                   |
      | Personal | cannot create link on personal space root |
      | Shares   | cannot create link on shares space root   |


  Scenario Outline: try to create an internal link share of a Personal and Shares drives with password using permissions endpoint
    When user "Alice" tries to create the following space link share using permissions endpoint of the Graph API:
      | space           | <drive>  |
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
              "message": {"const": "<message>"}
            }
          }
        }
      }
      """
    Examples:
      | drive    | message                                     |
      | Personal | password is redundant for the internal link |
      | Shares   | cannot create link on shares space root     |


  Scenario: create an internal quick link share of a project-space using permissions endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using permissions endpoint of the Graph API:
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


  Scenario: create an internal quick link share of a project-space multiple times using root endpoint
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "projectSpace" with the default quota using the Graph API
    When user "Alice" creates the following space link share using root endpoint of the Graph API:
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
    When user "Alice" creates the following space link share using root endpoint of the Graph API:
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
    And the last 2 link shares should have the same id
