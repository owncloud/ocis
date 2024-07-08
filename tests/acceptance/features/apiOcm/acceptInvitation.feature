@ocm
Feature: accepting invitation
  As a user
  I can accept invitations from users of other ocis instances

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And using server "REMOTE"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |


  Scenario: user accepts invitation
    Given using server "LOCAL"
    When "Alice" creates the federation share invitation
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expiration",
          "token"
        ],
        "properties": {
          "expiration": {
            "type": "integer",
            "pattern": "^[0-9]{10}$"
          },
          "token": {
            "type": "string",
            "pattern": "^%fed_invitation_token_pattern%$"
          }
        }
      }
      """
    When using server "REMOTE"
    And "Brian" accepts the last federation share invitation
    Then the HTTP status code should be "200"


  Scenario: user accepts invitation sent with email and description
    Given using server "LOCAL"
    When "Alice" creates the federation share invitation with email "alice@example.com" and description "a share invitation from Alice"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expiration",
          "token"
        ],
        "properties": {
          "expiration": {
            "type": "integer",
            "pattern": "^[0-9]{10}$"
          },
          "token": {
            "type": "string",
            "pattern": "^%fed_invitation_token_pattern%$"
          }
        }
      }
      """
    When using server "REMOTE"
    And "Brian" accepts the last federation share invitation
    Then the HTTP status code should be "200"


  Scenario: two users can accept one invitation
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    When using server "REMOTE"
    And "Brian" accepts the last federation share invitation
    Then the HTTP status code should be "200"
    And "Carol" accepts the last federation share invitation
    And the HTTP status code should be "200"


  Scenario: user tries to accept the invitation twice
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    When using server "REMOTE"
    And "Brian" accepts the last federation share invitation
    Then the HTTP status code should be "200"
    When "Brian" tries to accept the last federation share invitation
    Then the HTTP status code should be "409"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "code",
          "message"
        ],
        "properties": {
          "code": {
            "const": "ALREADY_EXIST"
          },
          "message": {
            "const": "user already known"
          }
        }
      }
      """


  Scenario: users try to accept each other's invitation
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    And "Brian" has created the federation share invitation
    When using server "LOCAL"
    And "Alice" tries to accept the last federation share invitation
    Then the HTTP status code should be "409"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "code",
          "message"
        ],
        "properties": {
          "code": {
            "const": "ALREADY_EXIST"
          },
          "message": {
            "const": "user already known"
          }
        }
      }
      """
