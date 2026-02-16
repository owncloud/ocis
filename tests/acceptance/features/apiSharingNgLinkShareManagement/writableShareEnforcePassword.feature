@env-config
Feature: enforce password on writable shares
  As a user
  I want to enforce passwords on public links shared with upload, edit, or contribute permission
  So that the password is required to access the contents of the link


  Background:
    Given the following configs have been set:
      | service | config                                            | value |
      | sharing | SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD           | false |
      | sharing | SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
    And user "Alice" has been created with default attributes
    And using spaces DAV path


  Scenario Outline: user should be able to create public link with permissions view and internal without a password
    Given user "Alice" has created folder "folder"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | View     |
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
              "type": {"const": "view"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
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


  Scenario Outline: user shouldn't be able to create writable public links without a password
    Given user "Alice" has created folder "folder"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
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
              "message": {"const": "password protection is enforced"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Edit             |
      | Upload           |
      | File Drop        |


  Scenario Outline: user should be able to update link share to writable permissions if password is set
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    And using SharingNG
    And user "Alice" has created the following resource link share:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | View     |
      | password        | %public% |
    And user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource        | folder                 |
      | space           | Personal               |
      | permissionsRole | <new-permissions-role> |
    Then the HTTP status code should be "200"
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API without a password
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong pass"
    But the public should be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "%public%"
    Examples:
      | new-permissions-role |
      | Edit                 |
      | Upload               |


  Scenario Outline: user should not be able to update link share to writable permissions if password is not set
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    And user "Alice" has created the following resource link share:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | View     |
    And user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource        | folder                 |
      | space           | Personal               |
      | permissionsRole | <new-permissions-role> |
    Then the HTTP status code should be "400"
    Examples:
      | new-permissions-role |
      | Edit                 |
      | Upload               |
      | File Drop            |
