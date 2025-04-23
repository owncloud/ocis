@ocm
Feature: delete federated connections
  As a user
  I want to delete federated connections if they are no longer needed

  Background:
    Given user "Alice" has been created with default attributes
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has accepted invitation


  Scenario: federated user deletes the federated connection
    When user "Brian" deletes federated connection with user "Alice" using the Graph API
    Then the HTTP status code should be "200"

  @issue-10216
  Scenario: users should not be able to find federated user after federated user has deleted connection
    Given user "Brian" has deleted federated connection with user "Alice"
    And using server "LOCAL"
    When user "Alice" searches for federated user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
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
    When user "Brian" searches for federated user "Alice" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 0,
            "maxItems": 0
          }
        }
      }
      """

  @issue-10216
  Scenario: federated user should not be able to find federated share after federated user has deleted connection
    Given using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And using server "REMOTE"
    And user "Brian" has deleted federated connection with user "Alice"
    When user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 0,
            "maxItems": 0,
          }
        }
      }
      """

  @issue-10213
  Scenario: federated user should not be able to find federated share after local user has deleted connection
    Given using server "LOCAL"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Alice" has deleted federated connection with user "Brian"
    And using server "REMOTE"
    When user "Brian" lists the shares shared with him without retry using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 0,
            "maxItems": 0,
          }
        }
      }
      """

  @issue-10223
  Scenario: local user tries to delete previously deleted federated connection
    Given using server "LOCAL"
    And user "Alice" has deleted federated connection with user "Brian"
    When user "Alice" tries to delete federated connection with user "Brian" using the Graph API
    Then the HTTP status code should be "404"

  @issue-10223
  Scenario: federated user tries to delete previously deleted federated connection
    Given user "Brian" has deleted federated connection with user "Alice"
    When user "Brian" tries to delete federated connection with user "Alice" using the Graph API
    Then the HTTP status code should be "404"

  @issue-10223
  Scenario Outline: federated user tries to delete previously deleted federated connection with random idp
    Given user "Brian" has deleted federated connection with user "Alice"
    When user "Brian" tries to delete federated connection with user "Alice" and provider "<idp>" using the Graph API
    Then the HTTP status code should be "400"
    Examples:
      | idp            |
      | localhost:9244 |
