@ocm @issue-enterprise-7075
Feature: delete ocm share
  As a user
  I want to delete federated share
  So that I can remove federated user access to resource

  Background:
    Given using SharingNG
    And using spaces DAV path
    And user "Alice" has been created with default attributes
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has accepted invitation
    And using server "LOCAL"


  Scenario Outline: deleting federated file share should delete share on both instances (Personal Space)
    Given user "Alice" has uploaded file with content "ocm test" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt       |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has removed the access of user "Brian" from resource "textfile.txt" of space "Personal"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": ["value"],
          "properties": {
            "value": {
              "type": "array",
              "minItems": 0,
              "maxItems": 0
            }
          }
        }
      """
    And using server "REMOTE"
    When user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": ["value"],
          "properties": {
            "value": {
              "type": "array",
              "minItems": 0,
              "maxItems": 0
            }
          }
        }
      """
    Examples:
      | permissions-role |
      | Viewer           |
      | File Editor      |


  Scenario Outline: deleting federated folder share should delete share on both instances (Personal Space)
    Given user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has removed the access of user "Brian" from resource "folderToShare" of space "Personal"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": ["value"],
          "properties": {
            "value": {
              "type": "array",
              "minItems": 0,
              "maxItems": 0
            }
          }
        }
      """
    And using server "REMOTE"
    When user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": ["value"],
          "properties": {
            "value": {
              "type": "array",
              "minItems": 0,
              "maxItems": 0
            }
          }
        }
      """
    Examples:
      | permissions-role |
      | Viewer           |
      | Editor           |
      | Uploader         |


  Scenario Outline: deleting federated file share should delete share on both instances (Project Space)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "ocm test" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt       |
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has removed the access of user "Brian" from resource "textfile.txt" of space "new-space"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": ["value"],
          "properties": {
            "value": {
              "type": "array",
              "minItems": 0,
              "maxItems": 0
            }
          }
        }
      """
    And using server "REMOTE"
    When user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": ["value"],
          "properties": {
            "value": {
              "type": "array",
              "minItems": 0,
              "maxItems": 0
            }
          }
        }
      """
    Examples:
      | permissions-role |
      | Viewer           |
      | File Editor      |


  Scenario Outline: deleting federated folder share should delete share on both instances (Project Space)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare      |
      | space           | new-space          |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has removed the access of user "Brian" from resource "folderToShare" of space "new-space"
    When user "Alice" lists the shares shared by her using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": ["value"],
          "properties": {
            "value": {
              "type": "array",
              "minItems": 0,
              "maxItems": 0
            }
          }
        }
      """
    And using server "REMOTE"
    When user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": ["value"],
          "properties": {
            "value": {
              "type": "array",
              "minItems": 0,
              "maxItems": 0
            }
          }
        }
      """
    Examples:
      | permissions-role |
      | Viewer           |
      | Editor           |
      | Uploader         |
